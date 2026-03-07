package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
	"github.com/dswcpp/lazygit/pkg/ai/agent"
	"github.com/dswcpp/lazygit/pkg/gui/style"
	"github.com/dswcpp/lazygit/pkg/gui/types"
	"github.com/jesseduffield/gocui"
)

// ChatMessage 聊天消息
type ChatMessage struct {
	Role      string // "user" | "assistant" | "system" | "action"
	Content   string
	Timestamp time.Time
	IsError   bool
	// Action 相关（仅 Role=="action" 时有效）
	ActionSuccess bool
	ActionType    string
}

// AIChatSession 保持 AI 对话的会话状态
type AIChatSession struct {
	c              *HelperCommon
	aiHelper       *AIHelper
	tr             *aii18n.Translator
	messages       []ChatMessage        // 已完成轮次的历史消息
	twoPhaseAgent  *agent.TwoPhaseAgent // 当前两阶段 Agent；在 PhaseWaitingConfirm 时保持存活
	isTyping       bool
	ctx            context.Context
	cancel         context.CancelFunc
	inputHistory   []string
	historyIndex   int
	scrollToBottom bool // 新消息到来时置 true，render 后重置，允许用户自由向上滚动
	statusLabel    string
	statusDetail   string
	logoFrame      int  // 动态 logo 帧索引
	isAnimating    bool // 标记动画是否正在运行
	mu             sync.Mutex // 保护并发访问
}

// AI Chat 动态 logo 字符序列
var aiChatLogoFrames = []string{
	"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
}

// agentSession 返回当前 TwoPhaseAgent 的会话（可能为 nil）。
func (s *AIChatSession) agentSession() *agent.Session {
	if s.twoPhaseAgent == nil {
		return nil
	}
	return s.twoPhaseAgent.Session()
}

// AIChatHelper 管理 AI 聊天弹窗
type AIChatHelper struct {
	c        *HelperCommon
	aiHelper *AIHelper
	session  *AIChatSession
}

func NewAIChatHelper(c *HelperCommon, aiHelper *AIHelper) *AIChatHelper {
	return &AIChatHelper{c: c, aiHelper: aiHelper}
}

// GetOrCreateSession 获取或创建聊天会话（保留历史）
func (self *AIChatHelper) GetOrCreateSession() *AIChatSession {
	if self.session == nil {
		ctx, cancel := context.WithCancel(context.Background())
		tr := aii18n.NewTranslator(self.c.Tr)
		self.session = &AIChatSession{
			c:            self.c,
			aiHelper:     self.aiHelper,
			tr:           tr,
			messages:     []ChatMessage{},
			ctx:          ctx,
			cancel:       cancel,
			inputHistory: []string{},
			historyIndex: -1,
			statusLabel:  self.c.Tr.AIIdle,
			statusDetail: self.c.Tr.AIChatCanInputNext,
		}
		self.session.addSystemMessage(tr.ChatWelcomeSystem())
		self.session.addAssistantMessage(tr.ChatWelcomeMessage())
	}
	return self.session
}

// CloseSession 关闭并清理当前会话
func (self *AIChatHelper) CloseSession() {
	if self.session != nil {
		self.session.mu.Lock()
		defer self.session.mu.Unlock()

		// 取消 context，停止所有 goroutine
		if self.session.cancel != nil {
			self.session.cancel()
		}

		// 清空引用，帮助 GC
		self.session = nil
	}
}

// ShowChat 打开 AI 聊天弹窗
func (self *AIChatHelper) ShowChat() error {
	return self.showChatInternal("")
}

// ShowChatWithContext 携带上下文内容打开 AI 聊天（用于从其他面板继续对话）
func (self *AIChatHelper) ShowChatWithContext(contextContent string) error {
	return self.showChatInternal(contextContent)
}

func (self *AIChatHelper) showChatInternal(followUpContext string) error {
	if self.c.AIManager == nil {
		self.c.Alert(self.c.Tr.AIChatNotEnabled, aii18n.NewTranslator(self.c.Tr).ChatConfigPrompt())
		return nil
	}

	session := self.GetOrCreateSession()

	if followUpContext != "" {
		session.addSystemMessage(session.tr.ChatPreviousContext())
		session.addAssistantMessage(followUpContext)
	}

	// 历史视图设置
	aiView := self.c.Views().AIChat
	aiView.Clear()
	aiView.Wrap = true
	aiView.Autoscroll = false
	aiView.Visible = true

	// 输入条设置：清空上次内容，显示可见
	inputView := self.c.Views().AIChatInput
	ResetAIChatInputView(inputView)
	inputView.Visible = true

	// 渲染已有消息
	session.scrollToBottom = true
	session.render()

	// 推入上下文（显示弹窗），焦点给输入条（gocui 会把光标渲染到 editable 视图）
	self.c.Context().Push(self.c.Contexts().AIChat, types.OnFocusOpts{})
	_, _ = self.c.GocuiGui().SetCurrentView(inputView.Name())

	// 启动标题动画（仅在未运行时启动）
	if !session.isAnimating {
		session.isAnimating = true
		go session.animateTitle()
	}

	return nil
}

// SendMessage 发送一条消息给 AI
func (self *AIChatHelper) SendMessage(content string) error {
	session := self.GetOrCreateSession()

	session.inputHistory = append(session.inputHistory, content)
	session.historyIndex = len(session.inputHistory)
	session.setStatus(self.c.Tr.AIThinking, self.c.Tr.AIChatGeneratingPlan)

	session.addUserMessage(content)
	session.render()

	go session.getAIResponse(content)
	return nil
}

// CopyLastResponse 复制最后一条 AI 回复到剪贴板
func (self *AIChatHelper) CopyLastResponse() error {
	if self.session == nil {
		self.c.Toast(aii18n.NewTranslator(self.c.Tr).ChatNoContentToCopy())
		return nil
	}
	return self.session.copyLastResponse()
}

// ExecuteLastCommands 提取并执行最后一条 AI 回复中的命令
func (self *AIChatHelper) ExecuteLastCommands() error {
	if self.session == nil {
		self.c.Toast(aii18n.NewTranslator(self.c.Tr).ChatNoExecutableReply())
		return nil
	}
	return self.session.executeLastResponseCommands()
}

// ClearHistory 清空对话历史
func (self *AIChatHelper) ClearHistory() error {
	if self.session == nil {
		return nil
	}
	return self.session.clearHistory()
}

// StopGeneration 停止当前 AI 生成
func (self *AIChatHelper) StopGeneration() error {
	if self.session == nil {
		return nil
	}
	return self.session.stopGeneration()
}

// --- AIChatSession 内部方法 ---

func (s *AIChatSession) addUserMessage(content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, ChatMessage{
		Role: "user", Content: content, Timestamp: time.Now(),
	})
	s.scrollToBottom = true
}

func (s *AIChatSession) addAssistantMessage(content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, ChatMessage{
		Role: "assistant", Content: content, Timestamp: time.Now(),
	})
	s.scrollToBottom = true
}

func (s *AIChatSession) addSystemMessage(content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, ChatMessage{
		Role: "system", Content: content, Timestamp: time.Now(),
	})
}

func (s *AIChatSession) addErrorMessage(content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, ChatMessage{
		Role: "assistant", Content: content, Timestamp: time.Now(), IsError: true,
	})
	s.scrollToBottom = true
}

// render 渲染所有消息到 AIChat 视图（包含 Agent 消息）
func (s *AIChatSession) render() {
	s.mu.Lock()
	defer s.mu.Unlock()
	aiView := s.c.Views().AIChat
	aiView.Clear()

	agentMsgs := s.agentUIMessages()
	total := len(s.messages) + len(agentMsgs)

	for i, msg := range s.messages {
		s.renderMessage(aiView, msg)
		if i < total-1 {
			fmt.Fprintln(aiView)
		}
	}
	for i, msg := range agentMsgs {
		s.renderMessage(aiView, msg)
		if len(s.messages)+i < total-1 {
			fmt.Fprintln(aiView)
		}
	}

	// 流式输出指示器：如果当前有流式消息正在进行，显示光标
	sess := s.agentSession()
	if sess != nil && sess.IsStreaming() {
		fmt.Fprintf(aiView, "%s", style.FgYellow.Sprint("▋")) // 闪烁光标效果
	}

	if s.isTyping && len(agentMsgs) == 0 && (sess == nil || !sess.IsStreaming()) {
		if total > 0 {
			fmt.Fprintln(aiView)
		}
		fmt.Fprintf(aiView, "  %s\n", style.FgYellow.Sprint(s.c.Tr.AIThinkingInProgress))
	}

	status, detail := s.deriveStatus()
	if total > 0 || s.isTyping {
		fmt.Fprintln(aiView)
	}
	s.renderStatus(aiView, status, detail)

	// PhaseWaitingConfirm：在底部显示输入提示（不打断滚动）
	if s.twoPhaseAgent != nil && s.twoPhaseAgent.Phase() == agent.PhaseWaitingConfirm && !s.isTyping {
		fmt.Fprintf(aiView, "  %s\n",
			style.FgYellow.Sprint(s.tr.ChatInputPrompt()))
	}

	// 仅在有新消息时才滚动到底部；用户手动向上滚动后不会被打断
	applyAIChatAutoScroll(aiView, &s.scrollToBottom)
}

// flushAgentSession 把当前 Agent 轮次的消息合并到 s.messages（持久化），
// 然后根据阶段决定是否清除 twoPhaseAgent：
//   - PhaseWaitingConfirm → 保留（下一条消息还需要它）
//   - 其他阶段 → 清除（本轮交互完成）
func (s *AIChatSession) flushAgentSession() {
	if s.twoPhaseAgent == nil {
		return
	}
	s.messages = append(s.messages, s.agentUIMessages()...)
	if s.twoPhaseAgent.Phase() != agent.PhaseWaitingConfirm {
		s.twoPhaseAgent = nil
	}
}

// agentUIMessages 把当前 Agent 会话的 UIMessages 转换为 ChatMessage 列表供渲染。
// KindUser 消息跳过（已在 s.messages 中显示，避免重复）。
func (s *AIChatSession) agentUIMessages() []ChatMessage {
	sess := s.agentSession()
	if sess == nil {
		return nil
	}
	msgs := sess.UIMessages
	result := make([]ChatMessage, 0, len(msgs))
	for _, m := range msgs {
		switch m.Kind {
		case agent.KindUser:
			// Already shown in s.messages — skip to avoid duplication.
		case agent.KindAssistant:
			result = append(result, ChatMessage{Role: "assistant", Content: m.Content, Timestamp: m.Timestamp})
		case agent.KindToolCall:
			result = append(result, ChatMessage{Role: "action", Content: m.Content, Timestamp: m.Timestamp, ActionType: m.ToolName})
		case agent.KindToolResult:
			result = append(result, ChatMessage{
				Role: "action", Content: m.Content, Timestamp: m.Timestamp,
				ActionSuccess: m.ToolSuccess, ActionType: m.ToolName,
			})
		case agent.KindSystem:
			result = append(result, ChatMessage{Role: "system", Content: m.Content, Timestamp: m.Timestamp})
		case agent.KindError:
			result = append(result, ChatMessage{Role: "assistant", Content: m.Content, Timestamp: m.Timestamp, IsError: true})
		case agent.KindPlan:
			// 执行计划：用 "plan" 角色渲染（带边框和步骤列表）
			result = append(result, ChatMessage{Role: "plan", Content: m.Content, Timestamp: m.Timestamp})
		case agent.KindStepUpdate:
			// 单步执行状态更新：用 "step" 角色渲染
			result = append(result, ChatMessage{
				Role: "step", Content: m.Content, Timestamp: m.Timestamp,
				ActionSuccess: m.ToolSuccess, ActionType: m.ToolName,
			})
		}
	}
	return result
}

// renderMessage 渲染单条消息（lazygit 极简风格，无 emoji）
func (s *AIChatSession) renderMessage(view *gocui.View, msg ChatMessage) {
	timeStr := msg.Timestamp.Format("15:04")
	w, _ := view.Size()
	if w < 20 {
		w = 60
	}

	switch msg.Role {
	case "system":
		// 系统消息：简洁分隔线风格
		inner := fmt.Sprintf("  %s  ", msg.Content)
		fillLen := w - len([]rune(inner)) - 4
		if fillLen < 0 {
			fillLen = 0
		}
		fmt.Fprintf(view, "── %s%s\n",
			style.FgDefault.Sprint(inner),
			style.FgDefault.Sprint(strings.Repeat("─", fillLen)),
		)

	case "user":
		// 用户消息标头
		header := fmt.Sprintf("─── You  %s ", timeStr)
		fillLen := w - len([]rune(header)) - 1
		if fillLen < 0 {
			fillLen = 0
		}
		fmt.Fprintf(view, "%s%s\n",
			style.FgCyan.Sprint(header),
			style.FgCyan.Sprint(strings.Repeat("─", fillLen)),
		)
		// 消息内容（缩进两格）
		for _, line := range strings.Split(msg.Content, "\n") {
			fmt.Fprintf(view, "  %s\n", line)
		}

	case "assistant":
		if strings.Contains(msg.Content, s.c.Tr.AIThinkingInProgress) {
			fmt.Fprintf(view, "  %s\n", style.FgYellow.Sprint(msg.Content))
			return
		}
		// AI 回复标头
		var headerColor style.TextStyle
		var headerLabel string
		if msg.IsError {
			headerColor = style.FgRed
			headerLabel = fmt.Sprintf("─── Error  %s ", timeStr)
		} else {
			headerColor = style.FgGreen
			headerLabel = fmt.Sprintf("─── AI  %s ", timeStr)
		}
		fillLen := w - len([]rune(headerLabel)) - 1
		if fillLen < 0 {
			fillLen = 0
		}
		fmt.Fprintf(view, "%s%s\n",
			headerColor.Sprint(headerLabel),
			headerColor.Sprint(strings.Repeat("─", fillLen)),
		)
		// 消息内容
		for _, line := range strings.Split(msg.Content, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "```") {
				fmt.Fprintf(view, "  %s\n", style.FgMagenta.Sprint(line))
			} else {
				fmt.Fprintf(view, "  %s\n", line)
			}
		}

	case "action":
		// 工具调用 / 结果
		var indicator string
		var lineColor style.TextStyle
		if msg.ActionSuccess {
			indicator = "✓"
			lineColor = style.FgGreen
		} else {
			indicator = "✗"
			lineColor = style.FgRed
		}
		header := fmt.Sprintf("  [%s %s]", indicator, msg.ActionType)
		fmt.Fprintf(view, "%s\n", lineColor.Sprint(header))
		for _, line := range strings.Split(strings.TrimSpace(msg.Content), "\n") {
			if line != "" {
				fmt.Fprintf(view, "    %s\n", style.FgDefault.Sprint(line))
			}
		}

	case "plan":
		// 执行计划：带边框，突出显示
		lines := strings.Split(msg.Content, "\n")
		summary := ""
		stepLines := []string{}
		for i, l := range lines {
			if i == 0 {
				summary = l
			} else if strings.TrimSpace(l) != "" {
				stepLines = append(stepLines, l)
			}
		}
		fmt.Fprintf(view, "  %s\n", style.FgYellow.Sprint("┌─ "+s.tr.ChatExecutionPlan()+" "+timeStr+" "+"─"))
		fmt.Fprintf(view, "  %s %s\n", style.FgYellow.Sprint("│"), summary)
		for _, sl := range stepLines {
			fmt.Fprintf(view, "  %s %s\n", style.FgYellow.Sprint("│"), style.FgDefault.Sprint(sl))
		}
		fmt.Fprintf(view, "  %s\n", style.FgYellow.Sprint("└"+"─────────────────────────────────────────────────────"))

	case "step":
		// 单步执行状态：一行简洁输出
		var indicator string
		var lineColor style.TextStyle
		if msg.ActionSuccess {
			indicator = "✓"
			lineColor = style.FgGreen
		} else {
			indicator = "✗"
			lineColor = style.FgRed
		}
		fmt.Fprintf(view, "  %s\n", lineColor.Sprint(fmt.Sprintf("%s %s", indicator, msg.Content)))
	}
}

// getAIResponse 通过 TwoPhaseAgent 处理用户消息（必须在 goroutine 中调用）。
//
// 行为由当前阶段决定：
//   - 无 Agent / PhaseDone / PhaseCancelled → 创建新 Agent，开始规划
//   - PhaseWaitingConfirm → 复用现有 Agent，处理 Y/N/补充说明
//   - PhaseExecuting → 忽略（执行中）
func (s *AIChatSession) getAIResponse(userMessage string) {
	s.isTyping = true

	mgr := s.c.AIManager
	if mgr == nil {
		s.c.GocuiGui().Update(func(*gocui.Gui) error {
			s.addErrorMessage(s.tr.ChatNotInitialized())
			s.isTyping = false
			s.render()
			return nil
		})
		return
	}

	// 需要新建 Agent 的条件：当前无 Agent，或上一轮已结束
	needNew := s.twoPhaseAgent == nil ||
		s.twoPhaseAgent.Phase() == agent.PhaseDone ||
		s.twoPhaseAgent.Phase() == agent.PhaseCancelled

	if needNew {
		// 上一轮结束时先 flush 历史
		if s.twoPhaseAgent != nil {
			s.c.GocuiGui().Update(func(*gocui.Gui) error {
				s.flushAgentSession()
				return nil
			})
		}
		s.twoPhaseAgent = mgr.NewTwoPhaseAgent(mgr.DefaultSkillTools())
	}

	a := s.twoPhaseAgent
	repoCtx := mgr.RepoContext()

	onUpdate := func() {
		s.scrollToBottom = true
		s.syncStatusFromCurrentState()
		s.c.GocuiGui().Update(func(*gocui.Gui) error {
			s.render()
			return nil
		})
	}

	err := a.Send(s.ctx, userMessage, repoCtx, onUpdate)

	s.c.GocuiGui().Update(func(*gocui.Gui) error {
		s.isTyping = false
		currentPhase := a.Phase()
		// 执行完成或取消后立即 flush；等待确认时保留 Agent
		if currentPhase == agent.PhaseDone || currentPhase == agent.PhaseCancelled {
			s.flushAgentSession()
		}
		if err != nil && !errors.Is(err, context.Canceled) {
			s.setStatus(s.c.Tr.AIFailed, s.tr.ChatRequestFailed())
			s.addErrorMessage(fmt.Sprintf("AI 请求失败: %v", err))
			s.flushAgentSession()
		} else if errors.Is(err, context.Canceled) {
			s.setStatus(s.c.Tr.AICancelled, s.tr.ChatGenerationStopped())
		} else {
			s.setTerminalStatusForPhase(currentPhase)
		}
		s.syncStatusFromCurrentState()
		s.render()
		return nil
	})
}

func (s *AIChatSession) copyLastResponse() error {
	for i := len(s.messages) - 1; i >= 0; i-- {
		msg := s.messages[i]
		if msg.Role == "assistant" && !msg.IsError && !strings.Contains(msg.Content, s.c.Tr.AIThinkingInProgress) {
			if err := s.c.OS().CopyToClipboard(msg.Content); err != nil {
				s.c.Toast(s.tr.ChatCopyFailed())
				return err
			}
			s.c.Toast(s.tr.ChatCopiedToClipboard())
			return nil
		}
	}
	s.c.Toast(s.tr.ChatNoContentToCopy())
	return nil
}

func (s *AIChatSession) executeLastResponseCommands() error {
	for i := len(s.messages) - 1; i >= 0; i-- {
		msg := s.messages[i]
		if msg.Role != "assistant" || msg.IsError || strings.Contains(msg.Content, s.c.Tr.AIThinkingInProgress) {
			continue
		}
		cmds := ExtractCommandsFromMessage(msg.Content)
		if len(cmds) == 0 {
			s.c.Toast(s.tr.ChatNoCommandsFound())
			return nil
		}
		return s.aiHelper.ConfirmAndSilentExecute(cmds)
	}
	s.c.Toast(s.tr.ChatNoExecutableReply())
	return nil
}

func (s *AIChatSession) clearHistory() error {
	s.c.Confirm(types.ConfirmOpts{
		Title:  s.tr.ChatClearHistoryTitle(),
		Prompt: s.tr.ChatClearHistoryPrompt(),
		HandleConfirm: func() error {
			s.messages = []ChatMessage{}
			s.twoPhaseAgent = nil
			s.setStatus(s.c.Tr.AIIdle, s.c.Tr.AIChatCanInputNext)
			s.addSystemMessage(s.tr.ChatHistoryCleared())
			s.addAssistantMessage(s.tr.ChatHowCanIHelp())
			s.render()
			return nil
		},
	})
	return nil
}

func (s *AIChatSession) stopGeneration() error {
	if s.isTyping {
		if s.cancel != nil {
			s.cancel()
			ctx, cancel := context.WithCancel(context.Background())
			s.ctx = ctx
			s.cancel = cancel
		}
		s.isTyping = false
		s.flushAgentSession()
		s.setStatus(s.c.Tr.AICancelled, s.tr.ChatGenerationStopped())
		s.addSystemMessage(s.tr.ChatStoppedGeneration())
		s.render()
		s.c.Toast(s.tr.ChatStoppedGeneration())
	}
	return nil
}

func applyAIChatAutoScroll(view *gocui.View, scrollToBottom *bool) {
	if view == nil || scrollToBottom == nil || !*scrollToBottom {
		return
	}

	view.ScrollDown(9999)
	*scrollToBottom = false
}

func ResetAIChatInputView(view *gocui.View) {
	if view == nil {
		return
	}

	// 先清空视图内容
	view.Clear()
	view.SetCursor(0, 0)
	view.SetOrigin(0, 0)

	// 重新初始化 TextArea（如果存在）
	if view.TextArea != nil {
		view.TextArea = &gocui.TextArea{}
		view.RenderTextArea()
	}
}

func (s *AIChatSession) setStatus(status, detail string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statusLabel = status
	s.statusDetail = detail
}

func (s *AIChatSession) setTerminalStatusForPhase(phase agent.AgentPhase) {
	switch phase {
	case agent.PhaseDone:
		s.setStatus(s.tr.ChatCompleted(), s.c.Tr.AIChatCanInputNext)
	case agent.PhaseCancelled:
		s.setStatus(s.c.Tr.AICancelled, s.c.Tr.AIChatCanInputNext)
	}
}

func (s *AIChatSession) syncStatusFromCurrentState() {
	status, detail := s.deriveStatus()
	s.mu.Lock()
	s.statusLabel = status
	s.statusDetail = detail
	s.mu.Unlock()
}

func (s *AIChatSession) getStatusPresentation() (string, string) {
	s.syncStatusFromCurrentState()
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.statusLabel, s.statusDetail
}

func (s *AIChatSession) deriveStatus() (string, string) {
	if s.twoPhaseAgent != nil {
		sess := s.twoPhaseAgent.Session() // needed for latestAgentActionDetail
		switch s.twoPhaseAgent.Phase() {
		case agent.PhaseWaitingConfirm:
			return s.tr.ChatWaitingConfirm(), s.tr.ChatConfirmPrompt()
		case agent.PhaseExecuting:
			return s.c.Tr.AIExecuting, latestAgentActionDetail(sess, s.tr.ChatExecutingPlan())
		case agent.PhasePlanning:
			if s.isTyping {
				return s.c.Tr.AIThinking, latestAgentActionDetail(sess, s.c.Tr.AIChatGeneratingPlan)
			}
		case agent.PhaseDone:
			return s.tr.ChatCompleted(), s.c.Tr.AIChatCanInputNext
		case agent.PhaseCancelled:
			return s.c.Tr.AICancelled, s.c.Tr.AIChatCanInputNext
		}
	}

	if s.isTyping {
		return s.c.Tr.AIThinking, s.tr.ChatGeneratingReply()
	}

	if s.statusLabel != "" {
		return s.statusLabel, s.statusDetail
	}

	return s.c.Tr.AIIdle, s.c.Tr.AIChatCanInputNext
}

func latestAgentActionDetail(sess *agent.Session, fallback string) string {
	if sess == nil {
		return fallback
	}

	for i := len(sess.UIMessages) - 1; i >= 0; i-- {
		msg := sess.UIMessages[i]
		switch msg.Kind {
		case agent.KindStepUpdate:
			return firstNonEmptyLine(msg.Content, fallback)
		case agent.KindToolCall:
			if msg.ToolName != "" {
				// 使用全局 Translator 实例或从 session 获取
				return fmt.Sprintf("%s %s", "正在调用", msg.ToolName)
			}
		case agent.KindToolResult:
			if msg.ToolName != "" {
				if msg.ToolSuccess {
					return fmt.Sprintf("%s %s", "已完成工具", msg.ToolName)
				}
				return fmt.Sprintf("%s %s 执行失败", "工具", msg.ToolName)
			}
		case agent.KindPlan:
			return "执行计划已生成，等待确认"
		case agent.KindError, agent.KindSystem:
			return firstNonEmptyLine(msg.Content, fallback)
		}
	}

	return fallback
}

func firstNonEmptyLine(content string, fallback string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}
	return fallback
}

func (s *AIChatSession) renderStatus(view *gocui.View, status string, detail string) {
	statusColor := style.FgGreen
	switch status {
	case s.c.Tr.AIThinking, s.tr.ChatWaitingConfirm():
		statusColor = style.FgYellow
	case s.c.Tr.AIExecuting:
		statusColor = style.FgCyan
	case s.c.Tr.AICancelled:
		statusColor = style.FgMagenta
	case s.c.Tr.AIFailed:
		statusColor = style.FgRed
	}

	fmt.Fprintf(view, "  %s %s\n", statusColor.Sprint(s.tr.ChatStatusLabel()), statusColor.Sprint(status))
	if detail != "" {
		fmt.Fprintf(view, "  %s %s\n", style.FgDefault.Sprint(s.tr.ChatActionLabel()), style.FgDefault.Sprint(detail))
	}
}

// animateTitle 动画更新标题中的 logo
func (s *AIChatSession) animateTitle() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	defer func() { s.isAnimating = false }() // 退出时重置标志

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.c.GocuiGui().Update(func(*gocui.Gui) error {
				aiView := s.c.Views().AIChat
				if aiView == nil || !aiView.Visible {
					return nil
				}

				// 更新 logo 帧
				s.logoFrame = (s.logoFrame + 1) % len(aiChatLogoFrames)
				logo := aiChatLogoFrames[s.logoFrame]

				// 根据状态选择颜色
				var coloredLogo string
				if s.isTyping {
					coloredLogo = style.FgYellow.Sprint(logo)
				} else {
					coloredLogo = style.FgCyan.Sprint(logo)
				}

				aiView.Title = fmt.Sprintf(" %s AI Chat ", coloredLogo)
				return nil
			})
		}
	}
}
