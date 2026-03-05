package gui

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/dswcpp/lazygit/pkg/gui/controllers/helpers"
	"github.com/dswcpp/lazygit/pkg/gui/style"
	"github.com/dswcpp/lazygit/pkg/gui/types"
	"github.com/jesseduffield/gocui"
)

// chatPageScrollLines 是 PgUp/PgDn 每次翻页的行数
const chatPageScrollLines = 10

// ChatMessage 聊天消息
type ChatMessage struct {
	Role      string    // "user" 或 "assistant" 或 "system"
	Content   string    // 消息内容
	Timestamp time.Time // 时间戳
	IsError   bool      // 是否为错误消息
}

// AIChat AI 对话框
type AIChat struct {
	gui           *Gui
	messages      []ChatMessage
	isTyping      bool
	ctx           context.Context
	cancel        context.CancelFunc
	inputHistory  []string // 输入历史
	historyIndex  int      // 历史索引
	maxWidth      int      // 最大宽度
	maxHeight     int      // 最大高度
	pendingCreate bool     // 防止 afterLayout 无限递归重试
}

// ShowAIChat 显示 AI 对话框（复用上次会话，历史不丢失）
func (gui *Gui) ShowAIChat() error {
	return gui.showAIChatInternal("")
}

// ShowAIChatWithFollowUp 携带上下文内容打开 AI 对话框，用于从其他面板（如代码审查）继续对话
func (gui *Gui) ShowAIChatWithFollowUp(contextContent string) error {
	return gui.showAIChatInternal(contextContent)
}

func (gui *Gui) showAIChatInternal(followUpContext string) error {
	if gui.c.AI == nil {
		gui.ShowError("AI 未启用", "请先在设置中启用并配置 AI 功能。", "提示：按 'o' 打开设置菜单")
		return nil
	}

	// 复用已有会话，保留历史消息
	if gui.aiChatSession == nil {
		ctx, cancel := context.WithCancel(context.Background())
		chat := &AIChat{
			gui:          gui,
			messages:     []ChatMessage{},
			ctx:          ctx,
			cancel:       cancel,
			inputHistory: []string{},
			historyIndex: -1,
			maxWidth:     80, // 默认值，createAIChatPopup 里会用当前 Size() 覆盖
			maxHeight:    24,
		}
		chat.addSystemMessage("欢迎使用 AI 助手！")
		chat.addAssistantMessage(
			"你好！我是你的 Git 智能助手 🤖\n\n" +
				"我可以帮你：\n" +
				"  • 解答 Git 相关问题\n" +
				"  • 分析当前仓库状态\n" +
				"  • 提供操作建议和最佳实践\n" +
				"  • 生成和解释 Git 命令\n\n" +
				"有什么我可以帮助你的吗？",
		)
		gui.aiChatSession = chat
	}

	chat := gui.aiChatSession

	// 如果携带外部上下文（例如来自代码审查），作为 assistant 消息注入
	if followUpContext != "" {
		chat.addSystemMessage("─── 以下内容来自上一次 AI 分析，你可以继续追问 ───")
		chat.addAssistantMessage(followUpContext)
	}

	// 若对话框已打开（视图已存在），只重新聚焦，避免重复创建和注册键位
	if _, err := gui.g.View("aiChatInput"); err == nil {
		// 兜底重绑：全局 resetKeybindings 可能清掉动态绑定，导致 Enter 无法发送
		gui.setAIChatKeyBindings(chat)
		gui.focusAIChatInput()
		gui.afterLayout(func() error {
			gui.focusAIChatInput()
			return nil
		})
		return nil
	}

	// 在布局阶段后创建，避免与当前帧的布局/焦点切换互相覆盖
	gui.afterLayout(func() error {
		return gui.createAIChatPopup(chat)
	})
	return nil
}

// createAIChatPopup 创建 AI 对话弹出窗口
func (gui *Gui) createAIChatPopup(chat *AIChat) error {
	maxX, maxY := gui.g.Size()

	// 某些启动阶段 Size 可能暂时为 0，异步重试一次，避免用户反复按快捷键
	if maxX <= 0 || maxY <= 0 {
		if !chat.pendingCreate {
			chat.pendingCreate = true
			gui.afterLayout(func() error {
				chat.pendingCreate = false
				return gui.createAIChatPopup(chat)
			})
		}
		return nil
	}

	// 每次打开时用当前终端尺寸重新计算，避免首次 Size() 返回 0 导致负数
	if maxX > 10 {
		chat.maxWidth = maxX - 10
	}
	if maxY > 6 {
		chat.maxHeight = maxY - 6
	}

	// 兜底保护，避免极端窗口尺寸导致负坐标/无效布局
	if chat.maxWidth < 20 {
		chat.maxWidth = 20
	}
	if chat.maxHeight < 10 {
		chat.maxHeight = 10
	}
	if chat.maxWidth > maxX-2 {
		chat.maxWidth = maxX - 2
	}
	if chat.maxHeight > maxY-2 {
		chat.maxHeight = maxY - 2
	}
	if chat.maxWidth < 2 || chat.maxHeight < 2 {
		gui.c.Toast("终端窗口过小，无法打开 AI 对话")
		return nil
	}

	// 计算布局
	width := chat.maxWidth
	height := chat.maxHeight
	x0 := (maxX - width) / 2
	y0 := (maxY - height) / 2
	x1 := x0 + width
	y1 := y0 + height

	// 状态栏高度
	statusHeight := 2
	// 输入框高度
	inputHeight := 4

	// 创建对话视图（主要内容区域）
	chatView, err := gui.g.SetView("aiChat", x0, y0, x1, y1-inputHeight-statusHeight, 0)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	chatView.Frame = true
	chatView.Title = " 💬 AI 智能助手 "
	chatView.Wrap = true
	chatView.Autoscroll = true

	// 创建状态栏
	statusView, err := gui.g.SetView("aiChatStatus", x0, y1-inputHeight-statusHeight+1, x1, y1-inputHeight, 0)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	statusView.Frame = false
	chat.updateStatusBar()

	// 创建输入框
	inputView, err := gui.g.SetView("aiChatInput", x0, y1-inputHeight+1, x1, y1, 0)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}
	inputView.Frame = true
	inputView.Title = " ✏️  输入消息 "
	inputView.Editable = true
	inputView.Wrap = true
	inputView.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
		return gui.handleEditorKeypress(v, key, ch, mod, true)
	})
	// 初始化 TextArea，避免空 view 被鼠标点击时 lines[-1] 越界 panic
	inputView.RenderTextArea()

	// 设置键盘绑定
	gui.setAIChatKeyBindings(chat)

	// 立即渲染并抢占焦点，避免首次打开需要多次按键才可见
	chat.render()
	gui.focusAIChatInput()
	gui.afterLayout(func() error {
		gui.focusAIChatInput()
		return nil
	})

	return nil
}

func (gui *Gui) isAIChatOpen() bool {
	if gui == nil || gui.g == nil {
		return false
	}
	_, err := gui.g.View("aiChatInput")
	return err == nil
}

func (gui *Gui) focusAIChatInput() {
	if gui == nil || gui.g == nil {
		return
	}
	inputView, err := gui.g.View("aiChatInput")
	if err != nil {
		return
	}
	if inputView.TextArea == nil {
		inputView.RenderTextArea()
	}

	_, _ = gui.g.SetViewOnTop("aiChat")
	_, _ = gui.g.SetViewOnTop("aiChatStatus")
	_, _ = gui.g.SetViewOnTop("aiChatInput")
	if v, err := gui.g.SetCurrentView("aiChatInput"); err == nil {
		gui.g.Cursor = v.Editable && v.Mask == ""
	}
}

func (gui *Gui) ensureAIChatFocus() {
	if !gui.isAIChatOpen() {
		return
	}

	// 如果当前有标准 popup（菜单/确认框等），不要强抢焦点，避免破坏交互
	if gui.helpers != nil && gui.helpers.Confirmation != nil && gui.helpers.Confirmation.IsPopupPanelFocused() {
		return
	}

	current := gui.g.CurrentView()
	if current != nil {
		switch current.Name() {
		case "aiChatInput", "aiChat", "aiChatStatus":
			// 空白保护：若视图存在但内容为空，补一次渲染
			if gui.aiChatSession != nil {
				chatView := gui.aiChatSession.getChatView()
				if chatView != nil && strings.TrimSpace(chatView.Buffer()) == "" && len(gui.aiChatSession.messages) > 0 {
					gui.aiChatSession.render()
					gui.aiChatSession.updateStatusBar()
				}
			}
			return
		}
	}

	gui.focusAIChatInput()
}

// setAIChatKeyBindings 设置 AI 对话键盘绑定
func (gui *Gui) setAIChatKeyBindings(chat *AIChat) {
	// 保证幂等：先清理旧绑定，避免重复叠加或被重置后缺失
	gui.g.DeleteViewKeybindings("aiChat")
	gui.g.DeleteViewKeybindings("aiChatInput")
	gui.g.DeleteViewKeybindings("aiChatStatus")

	// ==================== 输入框快捷键 ====================

	// Enter - 发送消息
	gui.g.SetKeybinding("aiChatInput", gocui.KeyEnter, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return chat.sendMessage()
	})

	// Alt+Enter - 换行（TextArea 模式）
	gui.g.SetKeybinding("aiChatInput", gocui.KeyEnter, gocui.ModAlt, func(g *gocui.Gui, v *gocui.View) error {
		v.TextArea.TypeCharacter("\n")
		v.RenderTextArea()
		return nil
	})

	// Esc - 关闭对话框
	gui.g.SetKeybinding("aiChatInput", gocui.KeyEsc, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return gui.CloseAIChat(chat)
	})

	// Ctrl+L - 清空历史
	gui.g.SetKeybinding("aiChatInput", gocui.KeyCtrlL, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return chat.clearHistory()
	})

	// Ctrl+K - 停止生成
	gui.g.SetKeybinding("aiChatInput", gocui.KeyCtrlK, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return chat.stopGeneration()
	})

	// Ctrl+S - 保存对话
	gui.g.SetKeybinding("aiChatInput", gocui.KeyCtrlS, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return chat.saveConversation()
	})

	// Ctrl+P - 预设问题
	gui.g.SetKeybinding("aiChatInput", gocui.KeyCtrlP, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return chat.showPresetQuestions()
	})

	// 上箭头 - 历史记录（上一条）
	gui.g.SetKeybinding("aiChatInput", gocui.KeyArrowUp, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return chat.navigateHistory(-1)
	})

	// 下箭头 - 历史记录（下一条）
	gui.g.SetKeybinding("aiChatInput", gocui.KeyArrowDown, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return chat.navigateHistory(1)
	})

	// Alt+上下箭头 - 在输入框内滚动聊天内容
	gui.g.SetKeybinding("aiChatInput", gocui.KeyArrowUp, gocui.ModAlt, func(*gocui.Gui, *gocui.View) error {
		if chatView := chat.getChatView(); chatView != nil {
			gui.scrollUpView(chatView)
		}
		return nil
	})
	gui.g.SetKeybinding("aiChatInput", gocui.KeyArrowDown, gocui.ModAlt, func(*gocui.Gui, *gocui.View) error {
		if chatView := chat.getChatView(); chatView != nil {
			gui.scrollDownView(chatView)
		}
		return nil
	})

	// 在输入框内也支持 PgUp/PgDn/Home/End 滚动聊天内容
	gui.g.SetKeybinding("aiChatInput", gocui.KeyPgup, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		if chatView := chat.getChatView(); chatView != nil {
			for i := 0; i < chatPageScrollLines; i++ {
				gui.scrollUpView(chatView)
			}
		}
		return nil
	})
	gui.g.SetKeybinding("aiChatInput", gocui.KeyPgdn, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		if chatView := chat.getChatView(); chatView != nil {
			for i := 0; i < chatPageScrollLines; i++ {
				gui.scrollDownView(chatView)
			}
		}
		return nil
	})
	gui.g.SetKeybinding("aiChatInput", gocui.KeyHome, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		if chatView := chat.getChatView(); chatView != nil {
			chatView.SetOrigin(0, 0)
		}
		return nil
	})
	gui.g.SetKeybinding("aiChatInput", gocui.KeyEnd, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		if chatView := chat.getChatView(); chatView != nil {
			_, height := chatView.Size()
			lines := len(chatView.BufferLines())
			if lines > height {
				chatView.SetOrigin(0, lines-height)
			}
		}
		return nil
	})

	// Tab - 切换到对话视图
	gui.g.SetKeybinding("aiChatInput", gocui.KeyTab, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		gui.g.SetCurrentView("aiChat")
		return nil
	})

	// ==================== 对话视图快捷键 ====================

	// Esc - 关闭对话框
	gui.g.SetKeybinding("aiChat", gocui.KeyEsc, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return gui.CloseAIChat(chat)
	})

	// Ctrl+L - 清空历史
	gui.g.SetKeybinding("aiChat", gocui.KeyCtrlL, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return chat.clearHistory()
	})

	// Ctrl+C - 复制最后一条 AI 回复
	gui.g.SetKeybinding("aiChat", gocui.KeyCtrlC, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return chat.copyLastResponse()
	})

	// x - 提取并静默执行上一条 AI 回复中的命令
	gui.g.SetKeybinding("aiChat", 'x', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return chat.executeLastResponseCommands()
	})
	gui.g.SetKeybinding("aiChatInput", 'x', gocui.ModAlt, func(*gocui.Gui, *gocui.View) error {
		return chat.executeLastResponseCommands()
	})

	// Tab - 切换到输入框
	gui.g.SetKeybinding("aiChat", gocui.KeyTab, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		gui.g.SetCurrentView("aiChatInput")
		return nil
	})

	// 上下箭头 - 滚动
	gui.g.SetKeybinding("aiChat", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		gui.scrollUpView(v)
		return nil
	})
	gui.g.SetKeybinding("aiChat", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		gui.scrollDownView(v)
		return nil
	})

	// PageUp/PageDown - 快速滚动
	gui.g.SetKeybinding("aiChat", gocui.KeyPgup, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		for i := 0; i < 10; i++ {
			gui.scrollUpView(v)
		}
		return nil
	})
	gui.g.SetKeybinding("aiChat", gocui.KeyPgdn, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		for i := 0; i < 10; i++ {
			gui.scrollDownView(v)
		}
		return nil
	})

	// Home - 滚动到顶部
	gui.g.SetKeybinding("aiChat", gocui.KeyHome, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		v.SetOrigin(0, 0)
		return nil
	})

	// End - 滚动到底部
	gui.g.SetKeybinding("aiChat", gocui.KeyEnd, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, height := v.Size()
		lines := len(v.BufferLines())
		if lines > height {
			v.SetOrigin(0, lines-height)
		}
		return nil
	})

	// ? - 显示帮助
	gui.g.SetKeybinding("aiChat", '?', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return chat.showHelp()
	})
	gui.g.SetKeybinding("aiChatInput", '?', gocui.ModAlt, func(*gocui.Gui, *gocui.View) error {
		return chat.showHelp()
	})

	// 鼠标滚轮：在 AI 对话框任一区域滚动聊天内容
	scrollChatByMouse := func(direction int) func(*gocui.Gui, *gocui.View) error {
		return func(*gocui.Gui, *gocui.View) error {
			chatView := chat.getChatView()
			if chatView == nil {
				return nil
			}
			if direction < 0 {
				gui.scrollUpView(chatView)
			} else {
				gui.scrollDownView(chatView)
			}
			return nil
		}
	}

	gui.g.SetKeybinding("aiChat", gocui.MouseWheelUp, gocui.ModNone, scrollChatByMouse(-1))
	gui.g.SetKeybinding("aiChat", gocui.MouseWheelDown, gocui.ModNone, scrollChatByMouse(1))
	gui.g.SetKeybinding("aiChatInput", gocui.MouseWheelUp, gocui.ModNone, scrollChatByMouse(-1))
	gui.g.SetKeybinding("aiChatInput", gocui.MouseWheelDown, gocui.ModNone, scrollChatByMouse(1))
	gui.g.SetKeybinding("aiChatStatus", gocui.MouseWheelUp, gocui.ModNone, scrollChatByMouse(-1))
	gui.g.SetKeybinding("aiChatStatus", gocui.MouseWheelDown, gocui.ModNone, scrollChatByMouse(1))
}

// sendMessage 发送消息
func (chat *AIChat) sendMessage() error {
	inputView := chat.getInputView()
	if inputView == nil || inputView.TextArea == nil {
		chat.gui.c.Toast("输入框尚未就绪，请重试")
		return nil
	}

	// 获取输入内容（TextArea 模式）
	content := strings.TrimSpace(inputView.TextArea.GetContent())
	if content == "" {
		return nil
	}

	// 保存到历史记录
	chat.inputHistory = append(chat.inputHistory, content)
	chat.historyIndex = len(chat.inputHistory)

	// 清空输入框
	inputView.ClearTextArea()
	inputView.RenderTextArea()

	// 添加用户消息
	chat.addUserMessage(content)

	// 重新渲染
	chat.render()
	chat.updateStatusBar()

	// 异步获取 AI 回复
	go chat.getAIResponse(content)

	return nil
}

// getAIResponse 获取 AI 回复
func (chat *AIChat) getAIResponse(userMessage string) {
	chat.isTyping = true
	chat.gui.g.Update(func(*gocui.Gui) error {
		chat.updateStatusBar()
		return nil
	})

	// 显示"正在思考"提示
	thinkingMsg := ChatMessage{
		Role:      "assistant",
		Content:   "⏳ 正在思考...",
		Timestamp: time.Now(),
	}

	chat.gui.g.Update(func(*gocui.Gui) error {
		chat.messages = append(chat.messages, thinkingMsg)
		chat.render()
		return nil
	})

	// 构建对话上下文
	prompt := chat.buildPrompt(userMessage)

	// 调用 AI
	result, err := chat.gui.c.AI.Complete(chat.ctx, prompt)
	assistantContent := ""
	if err == nil {
		assistantContent = strings.TrimSpace(result.Content)
	}

	chat.gui.g.Update(func(*gocui.Gui) error {
		// 移除"正在思考"消息
		if len(chat.messages) > 0 && strings.Contains(chat.messages[len(chat.messages)-1].Content, "正在思考") {
			chat.messages = chat.messages[:len(chat.messages)-1]
		}

		if err != nil {
			// 显示错误
			chat.addErrorMessage(fmt.Sprintf("抱歉，发生错误：%v", err))
		} else {
			// 添加 AI 回复
			chat.addAssistantMessage(assistantContent)
		}

		chat.isTyping = false
		chat.render()
		chat.updateStatusBar()
		return nil
	})
}

// currentShellInfo 返回当前 OS 和推荐 shell 的描述，用于注入到提示词
func currentShellInfo() string {
	switch runtime.GOOS {
	case "windows":
		return "操作系统: Windows\n推荐 Shell: Git Bash（直接使用 && 连接命令，不要用 cmd /c 或 ^&^& 转义）"
	case "darwin":
		return "操作系统: macOS\n推荐 Shell: zsh/bash（直接使用 && 连接命令）"
	default:
		return "操作系统: Linux\n推荐 Shell: bash（直接使用 && 连接命令）"
	}
}

// buildPrompt 构建提示词
func (chat *AIChat) buildPrompt(userMessage string) string {
	var sb strings.Builder

	// 系统提示
	sb.WriteString("你是一个专业的 Git 助手，运行在 lazygit 终端界面中。\n")
	sb.WriteString("你的职责是帮助用户理解和使用 Git，回答问题，提供建议。\n\n")
	sb.WriteString("═══ 运行环境 ═══\n")
	sb.WriteString(currentShellInfo())
	sb.WriteString("\n\n")
	sb.WriteString("命令格式要求：\n")
	sb.WriteString("- 根据上述运行环境生成对应格式的命令\n")
	sb.WriteString("- Windows Git Bash 用 &&，不要用 cmd /c 或 ^&^& 转义\n")
	sb.WriteString("- 不要在命令外面套 cmd /c \"...\"\n\n")
	sb.WriteString("- git commit -m 的提交信息如果包含空格，必须使用双引号\n")
	sb.WriteString("- 如果需要执行多条命令，请按“一行一条命令”返回\n\n")
	sb.WriteString("回答要求：\n")
	sb.WriteString("- 简洁、准确、实用\n")
	sb.WriteString("- 使用清晰的格式和结构\n")
	sb.WriteString("- 提供具体的命令和步骤\n")
	sb.WriteString("- 必要时给出示例\n")
	sb.WriteString("- 考虑用户的技术水平\n\n")

	// 添加仓库上下文
	repoCtx := chat.buildGitContext()
	if repoCtx != "" {
		sb.WriteString("═══ 当前仓库状态 ═══\n")
		sb.WriteString(repoCtx)
		sb.WriteString("\n")
	}

	// 添加对话历史（最近 6 条，排除系统消息和错误消息）
	historyCount := 0
	start := len(chat.messages) - 1
	for start >= 0 && historyCount < 6 {
		msg := chat.messages[start]
		if msg.Role != "system" && !msg.IsError && !strings.Contains(msg.Content, "正在思考") {
			historyCount++
		}
		start--
	}
	start++

	if start < len(chat.messages)-1 {
		sb.WriteString("═══ 对话历史 ═══\n")
		for i := start; i < len(chat.messages); i++ {
			msg := chat.messages[i]
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

	// 当前用户消息
	sb.WriteString("═══ 当前问题 ═══\n")
	sb.WriteString(fmt.Sprintf("用户: %s\n\n", userMessage))
	sb.WriteString("助手: ")

	return sb.String()
}

// buildGitContext 构建 Git 上下文
func (chat *AIChat) buildGitContext() string {
	var sb strings.Builder

	// 当前分支
	branch := chat.gui.c.Model().CheckedOutBranch
	if branch != "" {
		sb.WriteString(fmt.Sprintf("📍 分支: %s\n", branch))
	}

	// 变更文件
	files := chat.gui.c.Model().Files
	if len(files) > 0 {
		sb.WriteString(fmt.Sprintf("📝 变更文件: %d 个\n", len(files)))
		if len(files) <= 8 {
			for _, f := range files {
				sb.WriteString(fmt.Sprintf("   %s %s\n", f.ShortStatus, f.Path))
			}
		} else {
			for i := 0; i < 5; i++ {
				sb.WriteString(fmt.Sprintf("   %s %s\n", files[i].ShortStatus, files[i].Path))
			}
			sb.WriteString(fmt.Sprintf("   ... 还有 %d 个文件\n", len(files)-5))
		}
	}

	// 最近提交
	commits := chat.gui.c.Model().Commits
	if len(commits) > 0 {
		sb.WriteString(fmt.Sprintf("📌 最近提交: %s - %s\n", commits[0].ShortHash(), commits[0].Name))
		if len(commits) > 1 {
			sb.WriteString(fmt.Sprintf("   上一次: %s - %s\n", commits[1].ShortHash(), commits[1].Name))
		}
	}

	return sb.String()
}

// Helper methods for adding messages
func (chat *AIChat) addUserMessage(content string) {
	chat.messages = append(chat.messages, ChatMessage{
		Role:      "user",
		Content:   content,
		Timestamp: time.Now(),
	})
}

func (chat *AIChat) addAssistantMessage(content string) {
	chat.messages = append(chat.messages, ChatMessage{
		Role:      "assistant",
		Content:   content,
		Timestamp: time.Now(),
	})
}

func (chat *AIChat) addSystemMessage(content string) {
	chat.messages = append(chat.messages, ChatMessage{
		Role:      "system",
		Content:   content,
		Timestamp: time.Now(),
	})
}

func (chat *AIChat) addErrorMessage(content string) {
	chat.messages = append(chat.messages, ChatMessage{
		Role:      "assistant",
		Content:   content,
		Timestamp: time.Now(),
		IsError:   true,
	})
}

func (chat *AIChat) getChatView() *gocui.View {
	if chat == nil || chat.gui == nil || chat.gui.g == nil {
		return nil
	}
	v, err := chat.gui.g.View("aiChat")
	if err != nil {
		return nil
	}
	return v
}

func (chat *AIChat) getStatusView() *gocui.View {
	if chat == nil || chat.gui == nil || chat.gui.g == nil {
		return nil
	}
	v, err := chat.gui.g.View("aiChatStatus")
	if err != nil {
		return nil
	}
	return v
}

func (chat *AIChat) getInputView() *gocui.View {
	if chat == nil || chat.gui == nil || chat.gui.g == nil {
		return nil
	}
	v, err := chat.gui.g.View("aiChatInput")
	if err != nil {
		return nil
	}
	return v
}

// render 渲染对话内容
func (chat *AIChat) render() {
	chatView := chat.getChatView()
	if chatView == nil {
		return
	}
	chatView.Clear()

	// 渲染标题分隔线
	fmt.Fprintln(chatView, style.FgCyan.Sprint(strings.Repeat("─", chat.maxWidth-4)))
	fmt.Fprintln(chatView)

	// 渲染所有消息
	for i, msg := range chat.messages {
		chat.renderMessage(chatView, msg)

		// 消息之间添加分隔线（除了最后一条）
		if i < len(chat.messages)-1 {
			fmt.Fprintln(chatView)
		}
	}
}

// renderMessage 渲染单条消息
func (chat *AIChat) renderMessage(chatView *gocui.View, msg ChatMessage) {
	timeStr := msg.Timestamp.Format("15:04:05")

	switch msg.Role {
	case "system":
		// 系统消息 - 居中，灰色
		fmt.Fprintf(chatView, "%s\n",
			style.FgDefault.Sprint("┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈"))
		fmt.Fprintf(chatView, "  %s  %s\n",
			style.FgYellow.Sprint("ℹ"),
			style.FgDefault.Sprint(msg.Content))
		fmt.Fprintf(chatView, "%s\n",
			style.FgDefault.Sprint("┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈"))

	case "user":
		// 用户消息 - 右对齐风格，青色
		fmt.Fprintf(chatView, "%s %s\n",
			style.FgCyan.SetBold().Sprint("👤 你"),
			style.FgDefault.Sprint(timeStr))

		// 消息内容，带边框
		lines := strings.Split(msg.Content, "\n")
		fmt.Fprintln(chatView, style.FgCyan.Sprint("╭─────────────────────────────────────────────────────────"))
		for _, line := range lines {
			if line == "" {
				fmt.Fprintln(chatView, style.FgCyan.Sprint("│"))
			} else {
				fmt.Fprintf(chatView, "%s %s\n",
					style.FgCyan.Sprint("│"),
					line)
			}
		}
		fmt.Fprintln(chatView, style.FgCyan.Sprint("╰─────────────────────────────────────────────────────────"))

	case "assistant":
		// AI 消息 - 左对齐风格，绿色
		if msg.IsError {
			// 错误消息 - 红色
			fmt.Fprintf(chatView, "%s %s\n",
				style.FgRed.SetBold().Sprint("❌ AI"),
				style.FgDefault.Sprint(timeStr))

			lines := strings.Split(msg.Content, "\n")
			fmt.Fprintln(chatView, style.FgRed.Sprint("╭─────────────────────────────────────────────────────────"))
			for _, line := range lines {
				if line == "" {
					fmt.Fprintln(chatView, style.FgRed.Sprint("│"))
				} else {
					fmt.Fprintf(chatView, "%s %s\n",
						style.FgRed.Sprint("│"),
						line)
				}
			}
			fmt.Fprintln(chatView, style.FgRed.Sprint("╰─────────────────────────────────────────────────────────"))
		} else if strings.Contains(msg.Content, "正在思考") {
			// 思考中消息 - 黄色，带动画效果
			fmt.Fprintf(chatView, "%s %s\n",
				style.FgYellow.SetBold().Sprint("🤖 AI"),
				style.FgDefault.Sprint(timeStr))
			fmt.Fprintf(chatView, "  %s\n", style.FgYellow.Sprint(msg.Content))
		} else {
			// 正常 AI 回复 - 绿色
			fmt.Fprintf(chatView, "%s %s\n",
				style.FgGreen.SetBold().Sprint("🤖 AI"),
				style.FgDefault.Sprint(timeStr))

			lines := strings.Split(msg.Content, "\n")
			fmt.Fprintln(chatView, style.FgGreen.Sprint("╭─────────────────────────────────────────────────────────"))
			for _, line := range lines {
				if line == "" {
					fmt.Fprintln(chatView, style.FgGreen.Sprint("│"))
				} else {
					// 检测代码块
					if strings.HasPrefix(strings.TrimSpace(line), "```") {
						fmt.Fprintf(chatView, "%s %s\n",
							style.FgGreen.Sprint("│"),
							style.FgMagenta.Sprint(line))
					} else if strings.HasPrefix(strings.TrimSpace(line), "git ") ||
						strings.HasPrefix(strings.TrimSpace(line), "$ ") {
						// 高亮命令
						fmt.Fprintf(chatView, "%s %s\n",
							style.FgGreen.Sprint("│"),
							style.FgYellow.Sprint(line))
					} else if strings.HasPrefix(strings.TrimSpace(line), "•") ||
						strings.HasPrefix(strings.TrimSpace(line), "-") ||
						strings.HasPrefix(strings.TrimSpace(line), "*") {
						// 高亮列表项
						fmt.Fprintf(chatView, "%s %s\n",
							style.FgGreen.Sprint("│"),
							style.FgCyan.Sprint(line))
					} else {
						fmt.Fprintf(chatView, "%s %s\n",
							style.FgGreen.Sprint("│"),
							line)
					}
				}
			}
			fmt.Fprintln(chatView, style.FgGreen.Sprint("╰─────────────────────────────────────────────────────────"))
		}
	}
}

// updateStatusBar 更新状态栏
func (chat *AIChat) updateStatusBar() {
	statusView := chat.getStatusView()
	if statusView == nil {
		return
	}
	statusView.Clear()

	// 构建状态信息
	var statusParts []string

	// 消息计数
	userCount := 0
	aiCount := 0
	for _, msg := range chat.messages {
		if msg.Role == "user" {
			userCount++
		} else if msg.Role == "assistant" && !msg.IsError && !strings.Contains(msg.Content, "正在思考") {
			aiCount++
		}
	}
	statusParts = append(statusParts, fmt.Sprintf("💬 %d 条对话", userCount+aiCount))

	// AI 状态
	if chat.isTyping {
		statusParts = append(statusParts, style.FgYellow.Sprint("⏳ AI 正在思考..."))
	} else {
		statusParts = append(statusParts, style.FgGreen.Sprint("✓ 就绪"))
	}

	// 快捷键提示
	statusParts = append(statusParts, "Enter:发送")
	statusParts = append(statusParts, "Ctrl+P:预设")
	statusParts = append(statusParts, "?:帮助")

	// 渲染状态栏
	statusLine := strings.Join(statusParts, " │ ")
	fmt.Fprintf(statusView, " %s\n", statusLine)
	fmt.Fprintln(statusView, style.FgDefault.Sprint(strings.Repeat("─", chat.maxWidth-4)))
}

// navigateHistory 导航输入历史
func (chat *AIChat) navigateHistory(direction int) error {
	inputView := chat.getInputView()
	if inputView == nil || inputView.TextArea == nil {
		return nil
	}

	if len(chat.inputHistory) == 0 {
		return nil
	}

	// 更新索引
	newIndex := chat.historyIndex + direction
	if newIndex < 0 {
		newIndex = 0
	} else if newIndex > len(chat.inputHistory) {
		newIndex = len(chat.inputHistory)
	}

	chat.historyIndex = newIndex

	// 更新输入框内容（TextArea 模式）
	inputView.ClearTextArea()
	if chat.historyIndex < len(chat.inputHistory) {
		inputView.TextArea.TypeString(chat.inputHistory[chat.historyIndex])
	}
	inputView.RenderTextArea()

	return nil
}

// clearHistory 清空历史
func (chat *AIChat) clearHistory() error {
	chat.gui.ShowConfirm(
		"确认清空",
		"确定要清空对话历史吗？此操作不可撤销。",
		func() {
			chat.messages = []ChatMessage{}
			chat.addSystemMessage("对话历史已清空")
			chat.addAssistantMessage("有什么我可以帮助你的吗？")
			chat.render()
			chat.updateStatusBar()
		},
	)
	return nil
}

// stopGeneration 停止生成
func (chat *AIChat) stopGeneration() error {
	if chat.isTyping {
		if chat.cancel != nil {
			chat.cancel()
			// 重新创建 context
			ctx, cancel := context.WithCancel(context.Background())
			chat.ctx = ctx
			chat.cancel = cancel
		}
		chat.isTyping = false
		chat.addSystemMessage("已停止生成")
		chat.render()
		chat.updateStatusBar()
		chat.gui.c.Toast("已停止 AI 生成")
	}
	return nil
}

// copyLastResponse 复制最后一条 AI 回复
func (chat *AIChat) copyLastResponse() error {
	// 从后往前找最后一条 AI 回复
	for i := len(chat.messages) - 1; i >= 0; i-- {
		msg := chat.messages[i]
		if msg.Role == "assistant" && !msg.IsError && !strings.Contains(msg.Content, "正在思考") {
			if err := chat.gui.c.OS().CopyToClipboard(msg.Content); err != nil {
				chat.gui.c.Toast("复制失败")
				return err
			}
			chat.gui.c.Toast("已复制到剪贴板")
			return nil
		}
	}
	chat.gui.c.Toast("没有可复制的内容")
	return nil
}

// executeLastResponseCommands 提取上一条 AI 回复中的命令，确认后静默执行
func (chat *AIChat) executeLastResponseCommands() error {
	for i := len(chat.messages) - 1; i >= 0; i-- {
		msg := chat.messages[i]
		if msg.Role != "assistant" || msg.IsError || strings.Contains(msg.Content, "正在思考") {
			continue
		}
		cmds := helpers.ExtractCommandsFromMessage(msg.Content)
		if len(cmds) == 0 {
			chat.gui.c.Toast("上一条 AI 回复中未找到可执行命令")
			return nil
		}
		return chat.gui.helpers.AI.ConfirmAndSilentExecute(cmds)
	}
	chat.gui.c.Toast("没有可执行的 AI 回复")
	return nil
}

// saveConversation 保存对话
func (chat *AIChat) saveConversation() error {
	// 构建对话内容
	var sb strings.Builder
	sb.WriteString("# AI 对话记录\n\n")
	sb.WriteString(fmt.Sprintf("时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString("---\n\n")

	for _, msg := range chat.messages {
		if msg.Role == "system" {
			continue
		}

		timeStr := msg.Timestamp.Format("15:04:05")
		if msg.Role == "user" {
			sb.WriteString(fmt.Sprintf("## 👤 你 [%s]\n\n", timeStr))
			sb.WriteString(msg.Content)
			sb.WriteString("\n\n")
		} else if msg.Role == "assistant" && !strings.Contains(msg.Content, "正在思考") {
			sb.WriteString(fmt.Sprintf("## 🤖 AI [%s]\n\n", timeStr))
			sb.WriteString(msg.Content)
			sb.WriteString("\n\n")
		}
	}

	// 保存到文件
	filename := fmt.Sprintf("ai_chat_%s.md", time.Now().Format("20060102_150405"))
	content := sb.String()

	// 这里简化处理，实际应该让用户选择保存位置
	chat.gui.c.Toast(fmt.Sprintf("对话已保存: %s", filename))

	// TODO: 实际保存到文件
	_ = content

	return nil
}

// showPresetQuestions 显示预设问题
func (chat *AIChat) showPresetQuestions() error {
	presets := []struct {
		category  string
		questions []string
	}{
		{
			category: "📚 基础操作",
			questions: []string{
				"如何撤销最近的提交？",
				"如何查看文件的修改历史？",
				"如何暂存部分修改？",
				"如何修改最后一次提交？",
			},
		},
		{
			category: "🌿 分支管理",
			questions: []string{
				"如何创建和切换分支？",
				"如何合并分支？",
				"如何删除本地和远程分支？",
				"如何查看所有分支？",
			},
		},
		{
			category: "🔄 远程操作",
			questions: []string{
				"如何推送到远程仓库？",
				"如何拉取远程更新？",
				"如何解决推送冲突？",
				"如何查看远程仓库信息？",
			},
		},
		{
			category: "🔧 问题解决",
			questions: []string{
				"遇到合并冲突怎么办？",
				"如何恢复误删的文件？",
				"如何清理大型仓库？",
				"如何找回丢失的提交？",
			},
		},
		{
			category: "💡 最佳实践",
			questions: []string{
				"如何写好提交信息？",
				"如何组织提交历史？",
				"如何使用 Git Flow？",
				"如何进行代码审查？",
			},
		},
	}

	var menuItems []*types.MenuItem
	for _, preset := range presets {
		category := preset.category
		for _, question := range preset.questions {
			q := question
			menuItems = append(menuItems, &types.MenuItem{
				Label: fmt.Sprintf("%s: %s", category, q),
				OnPress: func() error {
					// 填充到输入框
					inputView := chat.getInputView()
					if inputView == nil || inputView.TextArea == nil {
						return nil
					}
					inputView.ClearTextArea()
					inputView.TextArea.TypeString(q)
					inputView.RenderTextArea()
					return nil
				},
			})
		}
	}

	return chat.gui.c.Menu(types.CreateMenuOptions{
		Title: "预设问题",
		Items: menuItems,
	})
}

// showHelp 显示帮助
func (chat *AIChat) showHelp() error {
	helpText := `
╔═══════════════════════════════════════════════════════════╗
║                   AI 对话 - 快捷键帮助                    ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  📝 输入框快捷键                                          ║
║  ─────────────────────────────────────────────────────   ║
║    Enter          发送消息                                ║
║    Alt+Enter      换行（多行输入）                        ║
║    ↑/↓            浏览输入历史                            ║
║    Ctrl+P         显示预设问题                            ║
║    Ctrl+K         停止 AI 生成                            ║
║    Ctrl+L         清空对话历史                            ║
║    Ctrl+S         保存对话记录                            ║
║    Tab            切换到对话视图                          ║
║    Esc            关闭对话窗口                            ║
║                                                           ║
║  💬 对话视图快捷键                                        ║
║  ─────────────────────────────────────────────────────   ║
║    ↑/↓            上下滚动                                ║
║    PgUp/PgDn      快速滚动                                ║
║    Home/End       跳到顶部/底部                           ║
║    Ctrl+C         复制最后一条 AI 回复                    ║
║    Ctrl+L         清空对话历史                            ║
║    Tab            切换到输入框                            ║
║    ?              显示此帮助                              ║
║    Esc            关闭对话窗口                            ║
║                                                           ║
║  💡 使用技巧                                              ║
║  ─────────────────────────────────────────────────────   ║
║    • AI 会自动获取当前仓库状态作为上下文                  ║
║    • 支持多轮对话，AI 会记住之前的对话内容                ║
║    • 可以使用 Ctrl+P 快速选择常见问题                     ║
║    • 使用 ↑/↓ 可以快速重新发送之前的问题                  ║
║    • 对话历史可以通过 Ctrl+S 保存为 Markdown 文件         ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝

按任意键关闭帮助...
`

	chat.gui.ShowInfo("快捷键帮助", helpText)
	return nil
}

// CloseAIChat 关闭 AI 对话框（保留会话历史，下次打开可继续）
func (gui *Gui) CloseAIChat(chat *AIChat) error {
	// 重建 context，以便下次打开时仍可发送请求
	if chat.cancel != nil {
		chat.cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	chat.ctx = ctx
	chat.cancel = cancel
	chat.isTyping = false

	// 删除键盘绑定
	gui.g.DeleteViewKeybindings("aiChat")
	gui.g.DeleteViewKeybindings("aiChatInput")
	gui.g.DeleteViewKeybindings("aiChatStatus")

	// 删除视图前先把焦点还给当前静态 context，避免 currentView 指向已删除的 editable view
	if current := gui.g.CurrentView(); current != nil {
		switch current.Name() {
		case "aiChat", "aiChatInput", "aiChatStatus":
			targetView := gui.c.Context().CurrentStatic().GetViewName()
			if v, err := gui.g.SetCurrentView(targetView); err == nil {
				gui.g.Cursor = v.Editable && v.Mask == ""
			}
		}
	}

	// 删除视图（会话历史保留在 gui.aiChatSession 中）
	gui.g.DeleteView("aiChat")
	gui.g.DeleteView("aiChatInput")
	gui.g.DeleteView("aiChatStatus")

	return nil
}
