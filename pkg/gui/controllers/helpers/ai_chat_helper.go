package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dswcpp/lazygit/pkg/ai/agent"
	"github.com/dswcpp/lazygit/pkg/gui/style"
	"github.com/dswcpp/lazygit/pkg/gui/types"
	"github.com/jesseduffield/gocui"
)

// ChatMessage 聊天消息
type ChatMessage struct {
	Role      string    // "user" | "assistant" | "system" | "action"
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
	messages       []ChatMessage   // welcome + user messages (pre-agent)
	agentSession   *agent.Session  // current agent session; nil when idle
	isTyping       bool
	ctx            context.Context
	cancel         context.CancelFunc
	inputHistory   []string
	historyIndex   int
	scrollToBottom bool // 新消息到来时置 true，render 后重置，允许用户自由向上滚动
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
		self.session = &AIChatSession{
			c:            self.c,
			aiHelper:     self.aiHelper,
			messages:     []ChatMessage{},
			ctx:          ctx,
			cancel:       cancel,
			inputHistory: []string{},
			historyIndex: -1,
		}
		self.session.addSystemMessage("欢迎使用 AI 助手！")
		self.session.addAssistantMessage(
			"你好！我是你的 Git Agent\n\n" +
				"我可以直接帮你操作仓库，例如：\n" +
				"  • 「帮我提交当前修改」\n" +
				"  • 「创建一个 feature/login 分支」\n" +
				"  • 「查看最近的提交记录」\n" +
				"  • 「把这些改动 stash 起来，切到 main 分支」\n\n" +
				"说你想做什么，我来执行。",
		)
	}
	return self.session
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
		self.c.Alert("AI 未启用", "请先在设置中启用并配置 AI 功能。\n提示：按 'o' 打开设置菜单")
		return nil
	}

	session := self.GetOrCreateSession()

	if followUpContext != "" {
		session.addSystemMessage("─── 以下内容来自上一次 AI 分析，你可以继续追问 ───")
		session.addAssistantMessage(followUpContext)
	}

	// 准备视图
	aiView := self.c.Views().AIChat
	aiView.Clear()
	aiView.Title = " AI Chat "
	aiView.Wrap = true
	aiView.Autoscroll = true

	// 渲染已有消息
	session.render()

	// 推入上下文（显示弹窗）
	self.c.Context().Push(self.c.Contexts().AIChat, types.OnFocusOpts{})
	return nil
}

// SendMessage 发送一条消息给 AI
func (self *AIChatHelper) SendMessage(content string) error {
	session := self.GetOrCreateSession()

	session.inputHistory = append(session.inputHistory, content)
	session.historyIndex = len(session.inputHistory)

	session.addUserMessage(content)
	session.render()

	go session.getAIResponse(content)
	return nil
}

// CopyLastResponse 复制最后一条 AI 回复到剪贴板
func (self *AIChatHelper) CopyLastResponse() error {
	if self.session == nil {
		self.c.Toast("没有可复制的内容")
		return nil
	}
	return self.session.copyLastResponse()
}

// ExecuteLastCommands 提取并执行最后一条 AI 回复中的命令
func (self *AIChatHelper) ExecuteLastCommands() error {
	if self.session == nil {
		self.c.Toast("没有可执行的 AI 回复")
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
	s.messages = append(s.messages, ChatMessage{
		Role: "user", Content: content, Timestamp: time.Now(),
	})
	s.scrollToBottom = true
}

func (s *AIChatSession) addAssistantMessage(content string) {
	s.messages = append(s.messages, ChatMessage{
		Role: "assistant", Content: content, Timestamp: time.Now(),
	})
	s.scrollToBottom = true
}

func (s *AIChatSession) addSystemMessage(content string) {
	s.messages = append(s.messages, ChatMessage{
		Role: "system", Content: content, Timestamp: time.Now(),
	})
}

func (s *AIChatSession) addErrorMessage(content string) {
	s.messages = append(s.messages, ChatMessage{
		Role: "assistant", Content: content, Timestamp: time.Now(), IsError: true,
	})
	s.scrollToBottom = true
}

// render 渲染所有消息到 AIChat 视图（包含 Agent 消息）
func (s *AIChatSession) render() {
	aiView := s.c.Views().AIChat
	aiView.Clear()

	agentMsgs := s.agentUIMessages()
	total := len(s.messages) + len(agentMsgs)

	for i, msg := range s.messages {
		renderAIChatMessage(aiView, msg)
		if i < total-1 {
			fmt.Fprintln(aiView)
		}
	}
	for i, msg := range agentMsgs {
		renderAIChatMessage(aiView, msg)
		if len(s.messages)+i < total-1 {
			fmt.Fprintln(aiView)
		}
	}

	if s.isTyping && len(agentMsgs) == 0 {
		if total > 0 {
			fmt.Fprintln(aiView)
		}
		fmt.Fprintf(aiView, "  %s\n", style.FgYellow.Sprint("正在思考..."))
	}

	// 仅在有新消息时才滚动到底部；用户手动向上滚动后不会被打断
	if s.scrollToBottom {
		aiView.ScrollDown(9999)
		s.scrollToBottom = false
	}
}

// flushAgentSession moves all messages from the current agent session into s.messages
// so they survive across turns, then clears agentSession.
func (s *AIChatSession) flushAgentSession() {
	if s.agentSession == nil {
		return
	}
	s.messages = append(s.messages, s.agentUIMessages()...)
	s.agentSession = nil
}

// agentUIMessages converts the current agent session's UIMessages to ChatMessages for rendering.
// KindUser messages are skipped as the user message is already in s.messages.
func (s *AIChatSession) agentUIMessages() []ChatMessage {
	if s.agentSession == nil {
		return nil
	}
	msgs := s.agentSession.UIMessages
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
		}
	}
	return result
}

// renderAIChatMessage 渲染单条消息（lazygit 极简风格，无 emoji）
func renderAIChatMessage(view *gocui.View, msg ChatMessage) {
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
		if strings.Contains(msg.Content, "正在思考") {
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
		// 操作执行结果
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
	}
}

// getAIResponse 使用 Agent ReAct 循环异步处理用户消息。
// 必须在 goroutine 中调用（非 UI 线程）。
func (s *AIChatSession) getAIResponse(userMessage string) {
	s.isTyping = true

	mgr := s.c.AIManager
	if mgr == nil {
		s.c.GocuiGui().Update(func(*gocui.Gui) error {
			s.addErrorMessage("AI 未初始化，请先配置 AI 功能。")
			s.isTyping = false
			s.render()
			return nil
		})
		return
	}

	a := mgr.NewAgent("", agent.AutoApproveAll())
	s.agentSession = a.Session()

	repoCtx := mgr.RepoContext()

	onUpdate := func() {
		s.scrollToBottom = true
		s.c.GocuiGui().Update(func(*gocui.Gui) error {
			s.render()
			return nil
		})
	}

	err := a.Run(s.ctx, userMessage, repoCtx, onUpdate)

	s.c.GocuiGui().Update(func(*gocui.Gui) error {
		s.isTyping = false
		s.flushAgentSession()
		if err != nil && !errors.Is(err, context.Canceled) {
			s.addErrorMessage(fmt.Sprintf("AI 请求失败: %v", err))
		}
		s.render()
		return nil
	})
}


func (s *AIChatSession) copyLastResponse() error {
	for i := len(s.messages) - 1; i >= 0; i-- {
		msg := s.messages[i]
		if msg.Role == "assistant" && !msg.IsError && !strings.Contains(msg.Content, "正在思考") {
			if err := s.c.OS().CopyToClipboard(msg.Content); err != nil {
				s.c.Toast("复制失败")
				return err
			}
			s.c.Toast("已复制到剪贴板")
			return nil
		}
	}
	s.c.Toast("没有可复制的内容")
	return nil
}

func (s *AIChatSession) executeLastResponseCommands() error {
	for i := len(s.messages) - 1; i >= 0; i-- {
		msg := s.messages[i]
		if msg.Role != "assistant" || msg.IsError || strings.Contains(msg.Content, "正在思考") {
			continue
		}
		cmds := ExtractCommandsFromMessage(msg.Content)
		if len(cmds) == 0 {
			s.c.Toast("上一条 AI 回复中未找到可执行命令")
			return nil
		}
		return s.aiHelper.ConfirmAndSilentExecute(cmds)
	}
	s.c.Toast("没有可执行的 AI 回复")
	return nil
}

func (s *AIChatSession) clearHistory() error {
	s.c.Confirm(types.ConfirmOpts{
		Title:  "确认清空",
		Prompt: "确定要清空对话历史吗？此操作不可撤销。",
		HandleConfirm: func() error {
			s.messages = []ChatMessage{}
			s.agentSession = nil
			s.addSystemMessage("对话历史已清空")
			s.addAssistantMessage("有什么我可以帮助你的吗？")
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
		s.addSystemMessage("已停止生成")
		s.render()
		s.c.Toast("已停止 AI 生成")
	}
	return nil
}
