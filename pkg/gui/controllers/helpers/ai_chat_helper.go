package helpers

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/dswcpp/lazygit/pkg/gui/style"
	"github.com/dswcpp/lazygit/pkg/gui/types"
	"github.com/jesseduffield/gocui"
)

// ChatMessage 聊天消息
type ChatMessage struct {
	Role      string    // "user" | "assistant" | "system"
	Content   string
	Timestamp time.Time
	IsError   bool
}

// AIChatSession 保持 AI 对话的会话状态
type AIChatSession struct {
	c            *HelperCommon
	aiHelper     *AIHelper
	messages     []ChatMessage
	isTyping     bool
	ctx          context.Context
	cancel       context.CancelFunc
	inputHistory []string
	historyIndex int
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
			"你好！我是你的 Git 智能助手\n\n" +
				"我可以帮你：\n" +
				"  • 解答 Git 相关问题\n" +
				"  • 分析当前仓库状态\n" +
				"  • 提供操作建议和最佳实践\n" +
				"  • 生成和解释 Git 命令\n\n" +
				"有什么我可以帮助你的吗？",
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
	if self.c.AI == nil {
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
}

func (s *AIChatSession) addAssistantMessage(content string) {
	s.messages = append(s.messages, ChatMessage{
		Role: "assistant", Content: content, Timestamp: time.Now(),
	})
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
}

// render 渲染所有消息到 AIChat 视图
func (s *AIChatSession) render() {
	aiView := s.c.Views().AIChat
	aiView.Clear()
	for i, msg := range s.messages {
		renderAIChatMessage(aiView, msg)
		if i < len(s.messages)-1 {
			fmt.Fprintln(aiView)
		}
	}
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
			} else if strings.HasPrefix(trimmed, "git ") || strings.HasPrefix(trimmed, "$ ") {
				fmt.Fprintf(view, "  %s\n", style.FgYellow.Sprint(line))
			} else {
				fmt.Fprintf(view, "  %s\n", line)
			}
		}
	}
}

// getAIResponse 异步获取 AI 回复
func (s *AIChatSession) getAIResponse(userMessage string) {
	s.isTyping = true

	// 显示"正在思考"占位
	s.c.GocuiGui().Update(func(*gocui.Gui) error {
		s.messages = append(s.messages, ChatMessage{
			Role:      "assistant",
			Content:   "正在思考...",
			Timestamp: time.Now(),
		})
		s.render()
		return nil
	})

	prompt := s.buildPrompt(userMessage)
	result, err := s.c.AI.Complete(s.ctx, prompt)
	assistantContent := ""
	if err == nil {
		assistantContent = strings.TrimSpace(result.Content)
	}

	s.c.GocuiGui().Update(func(*gocui.Gui) error {
		// 移除"正在思考"消息
		if len(s.messages) > 0 && strings.Contains(s.messages[len(s.messages)-1].Content, "正在思考") {
			s.messages = s.messages[:len(s.messages)-1]
		}
		if err != nil {
			s.addErrorMessage(fmt.Sprintf("抱歉，发生错误：%v", err))
		} else {
			s.addAssistantMessage(assistantContent)
		}
		s.isTyping = false
		s.render()
		return nil
	})
}

func currentShellInfoForChat() string {
	switch runtime.GOOS {
	case "windows":
		return "操作系统: Windows\n推荐 Shell: Git Bash（直接使用 && 连接命令）"
	case "darwin":
		return "操作系统: macOS\n推荐 Shell: zsh/bash"
	default:
		return "操作系统: Linux\n推荐 Shell: bash"
	}
}

func (s *AIChatSession) buildPrompt(userMessage string) string {
	var sb strings.Builder

	sb.WriteString("你是一个专业的 Git 助手，运行在 lazygit 终端界面中。\n")
	sb.WriteString("你的职责是帮助用户理解和使用 Git，回答问题，提供建议。\n\n")
	sb.WriteString("═══ 运行环境 ═══\n")
	sb.WriteString(currentShellInfoForChat())
	sb.WriteString("\n\n命令格式要求：\n")
	sb.WriteString("- git commit -m 的提交信息如果包含空格，必须使用双引号\n")
	sb.WriteString("- 如果需要执行多条命令，请每行一条命令\n\n")
	sb.WriteString("回答要求：简洁、准确、实用，提供具体的命令和步骤。\n\n")

	repoCtx := s.buildGitContext()
	if repoCtx != "" {
		sb.WriteString("═══ 当前仓库状态 ═══\n")
		sb.WriteString(repoCtx)
		sb.WriteString("\n")
	}

	// 最近 6 条对话历史
	historyCount := 0
	start := len(s.messages) - 1
	for start >= 0 && historyCount < 6 {
		msg := s.messages[start]
		if msg.Role != "system" && !msg.IsError && !strings.Contains(msg.Content, "正在思考") {
			historyCount++
		}
		start--
	}
	start++
	if start < len(s.messages)-1 {
		sb.WriteString("═══ 对话历史 ═══\n")
		for i := start; i < len(s.messages); i++ {
			msg := s.messages[i]
			if msg.Role == "system" || msg.IsError || strings.Contains(msg.Content, "正在思考") {
				continue
			}
			if msg.Role == "user" {
				sb.WriteString(fmt.Sprintf("用户: %s\n", msg.Content))
			} else {
				sb.WriteString(fmt.Sprintf("助手: %s\n", msg.Content))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("═══ 当前问题 ═══\n")
	sb.WriteString(fmt.Sprintf("用户: %s\n\nassistant: ", userMessage))
	return sb.String()
}

func (s *AIChatSession) buildGitContext() string {
	var sb strings.Builder
	branch := s.c.Model().CheckedOutBranch
	if branch != "" {
		sb.WriteString(fmt.Sprintf("分支: %s\n", branch))
	}
	files := s.c.Model().Files
	if len(files) > 0 {
		sb.WriteString(fmt.Sprintf("变更文件: %d 个\n", len(files)))
		limit := len(files)
		if limit > 5 {
			limit = 5
		}
		for i := 0; i < limit; i++ {
			sb.WriteString(fmt.Sprintf("  %s %s\n", files[i].ShortStatus, files[i].Path))
		}
		if len(files) > 5 {
			sb.WriteString(fmt.Sprintf("  ... 还有 %d 个文件\n", len(files)-5))
		}
	}
	commits := s.c.Model().Commits
	if len(commits) > 0 {
		sb.WriteString(fmt.Sprintf("最近提交: %s - %s\n", commits[0].ShortHash(), commits[0].Name))
	}
	return sb.String()
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
		s.addSystemMessage("已停止生成")
		s.render()
		s.c.Toast("已停止 AI 生成")
	}
	return nil
}
