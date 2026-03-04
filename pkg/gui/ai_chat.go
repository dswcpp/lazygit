package gui

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/gui/style"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

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
	chatView      *gocui.View
	inputView     *gocui.View
	statusView    *gocui.View
	messages      []ChatMessage
	isTyping      bool
	ctx           context.Context
	cancel        context.CancelFunc
	inputHistory  []string // 输入历史
	historyIndex  int      // 历史索引
	maxWidth      int      // 最大宽度
	maxHeight     int      // 最大高度
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
		maxX, maxY := gui.g.Size()
		chat := &AIChat{
			gui:          gui,
			messages:     []ChatMessage{},
			ctx:          ctx,
			cancel:       cancel,
			inputHistory: []string{},
			historyIndex: -1,
			maxWidth:     maxX - 10,
			maxHeight:    maxY - 6,
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

	// 如果携带外部上下文（例如来自代码审查），作为 assistant 消息注入，并提示用户继续提问
	if followUpContext != "" {
		chat.addSystemMessage("─── 以下内容来自上一次 AI 分析，你可以继续追问 ───")
		chat.addAssistantMessage(followUpContext)
	}

	return gui.createAIChatPopup(chat)
}

// createAIChatPopup 创建 AI 对话弹出窗口
func (gui *Gui) createAIChatPopup(chat *AIChat) error {
	maxX, maxY := gui.g.Size()

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
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	chatView.Frame = true
	chatView.Title = " 💬 AI 智能助手 "
	chatView.Wrap = true
	chatView.Autoscroll = true
	chat.chatView = chatView

	// 创建状态栏
	statusView, err := gui.g.SetView("aiChatStatus", x0, y1-inputHeight-statusHeight+1, x1, y1-inputHeight, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	statusView.Frame = false
	chat.statusView = statusView
	chat.updateStatusBar()

	// 创建输入框
	inputView, err := gui.g.SetView("aiChatInput", x0, y1-inputHeight+1, x1, y1, 0)
	if err != nil && err != gocui.ErrUnknownView {
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
	chat.inputView = inputView

	// 渲染初始内容
	chat.render()

	// 设置焦点和层级
	gui.g.SetViewOnTop("aiChat")
	gui.g.SetViewOnTop("aiChatStatus")
	gui.g.SetViewOnTop("aiChatInput")
	gui.g.SetCurrentView("aiChatInput")

	// 设置键盘绑定
	gui.setAIChatKeyBindings(chat)

	return nil
}

// setAIChatKeyBindings 设置 AI 对话键盘绑定
func (gui *Gui) setAIChatKeyBindings(chat *AIChat) {
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
}

// sendMessage 发送消息
func (chat *AIChat) sendMessage() error {
	// 获取输入内容（TextArea 模式）
	content := strings.TrimSpace(chat.inputView.TextArea.GetContent())
	if content == "" {
		return nil
	}

	// 保存到历史记录
	chat.inputHistory = append(chat.inputHistory, content)
	chat.historyIndex = len(chat.inputHistory)

	// 清空输入框
	chat.inputView.ClearTextArea()
	chat.inputView.RenderTextArea()

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
			chat.addAssistantMessage(strings.TrimSpace(result.Content))
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

// Continue in next part...
// 继续 ai_chat_v2.go 的实现

// render 渲染对话内容
func (chat *AIChat) render() {
	chat.chatView.Clear()

	// 渲染标题分隔线
	fmt.Fprintln(chat.chatView, style.FgCyan.Sprint(strings.Repeat("─", chat.maxWidth-4)))
	fmt.Fprintln(chat.chatView)

	// 渲染所有消息
	for i, msg := range chat.messages {
		chat.renderMessage(msg)

		// 消息之间添加分隔线（除了最后一条）
		if i < len(chat.messages)-1 {
			fmt.Fprintln(chat.chatView)
		}
	}
}

// renderMessage 渲染单条消息
func (chat *AIChat) renderMessage(msg ChatMessage) {
	timeStr := msg.Timestamp.Format("15:04:05")

	switch msg.Role {
	case "system":
		// 系统消息 - 居中，灰色
		fmt.Fprintf(chat.chatView, "%s\n",
			style.FgDefault.Sprint("┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈"))
		fmt.Fprintf(chat.chatView, "  %s  %s\n",
			style.FgYellow.Sprint("ℹ"),
			style.FgDefault.Sprint(msg.Content))
		fmt.Fprintf(chat.chatView, "%s\n",
			style.FgDefault.Sprint("┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈"))

	case "user":
		// 用户消息 - 右对齐风格，青色
		fmt.Fprintf(chat.chatView, "%s %s\n",
			style.FgCyan.SetBold().Sprint("👤 你"),
			style.FgDefault.Sprint(timeStr))

		// 消息内容，带边框
		lines := strings.Split(msg.Content, "\n")
		fmt.Fprintln(chat.chatView, style.FgCyan.Sprint("╭─────────────────────────────────────────────────────────"))
		for _, line := range lines {
			if line == "" {
				fmt.Fprintln(chat.chatView, style.FgCyan.Sprint("│"))
			} else {
				fmt.Fprintf(chat.chatView, "%s %s\n",
					style.FgCyan.Sprint("│"),
					line)
			}
		}
		fmt.Fprintln(chat.chatView, style.FgCyan.Sprint("╰─────────────────────────────────────────────────────────"))

	case "assistant":
		// AI 消息 - 左对齐风格，绿色
		if msg.IsError {
			// 错误消息 - 红色
			fmt.Fprintf(chat.chatView, "%s %s\n",
				style.FgRed.SetBold().Sprint("❌ AI"),
				style.FgDefault.Sprint(timeStr))

			lines := strings.Split(msg.Content, "\n")
			fmt.Fprintln(chat.chatView, style.FgRed.Sprint("╭─────────────────────────────────────────────────────────"))
			for _, line := range lines {
				if line == "" {
					fmt.Fprintln(chat.chatView, style.FgRed.Sprint("│"))
				} else {
					fmt.Fprintf(chat.chatView, "%s %s\n",
						style.FgRed.Sprint("│"),
						line)
				}
			}
			fmt.Fprintln(chat.chatView, style.FgRed.Sprint("╰─────────────────────────────────────────────────────────"))
		} else if strings.Contains(msg.Content, "正在思考") {
			// 思考中消息 - 黄色，带动画效果
			fmt.Fprintf(chat.chatView, "%s %s\n",
				style.FgYellow.SetBold().Sprint("🤖 AI"),
				style.FgDefault.Sprint(timeStr))
			fmt.Fprintf(chat.chatView, "  %s\n", style.FgYellow.Sprint(msg.Content))
		} else {
			// 正常 AI 回复 - 绿色
			fmt.Fprintf(chat.chatView, "%s %s\n",
				style.FgGreen.SetBold().Sprint("🤖 AI"),
				style.FgDefault.Sprint(timeStr))

			lines := strings.Split(msg.Content, "\n")
			fmt.Fprintln(chat.chatView, style.FgGreen.Sprint("╭─────────────────────────────────────────────────────────"))
			for _, line := range lines {
				if line == "" {
					fmt.Fprintln(chat.chatView, style.FgGreen.Sprint("│"))
				} else {
					// 检测代码块
					if strings.HasPrefix(strings.TrimSpace(line), "```") {
						fmt.Fprintf(chat.chatView, "%s %s\n",
							style.FgGreen.Sprint("│"),
							style.FgMagenta.Sprint(line))
					} else if strings.HasPrefix(strings.TrimSpace(line), "git ") ||
						strings.HasPrefix(strings.TrimSpace(line), "$ ") {
						// 高亮命令
						fmt.Fprintf(chat.chatView, "%s %s\n",
							style.FgGreen.Sprint("│"),
							style.FgYellow.Sprint(line))
					} else if strings.HasPrefix(strings.TrimSpace(line), "•") ||
						strings.HasPrefix(strings.TrimSpace(line), "-") ||
						strings.HasPrefix(strings.TrimSpace(line), "*") {
						// 高亮列表项
						fmt.Fprintf(chat.chatView, "%s %s\n",
							style.FgGreen.Sprint("│"),
							style.FgCyan.Sprint(line))
					} else {
						fmt.Fprintf(chat.chatView, "%s %s\n",
							style.FgGreen.Sprint("│"),
							line)
					}
				}
			}
			fmt.Fprintln(chat.chatView, style.FgGreen.Sprint("╰─────────────────────────────────────────────────────────"))
		}
	}
}

// updateStatusBar 更新状态栏
func (chat *AIChat) updateStatusBar() {
	chat.statusView.Clear()

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
	fmt.Fprintf(chat.statusView, " %s\n", statusLine)
	fmt.Fprintln(chat.statusView, style.FgDefault.Sprint(strings.Repeat("─", chat.maxWidth-4)))
}

// navigateHistory 导航输入历史
func (chat *AIChat) navigateHistory(direction int) error {
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
	chat.inputView.ClearTextArea()
	if chat.historyIndex < len(chat.inputHistory) {
		for _, ch := range chat.inputHistory[chat.historyIndex] {
			chat.inputView.TextArea.TypeCharacter(string(ch))
		}
	}
	chat.inputView.RenderTextArea()

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
		category string
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
					chat.inputView.Clear()
					chat.inputView.SetCursor(0, 0)
					fmt.Fprint(chat.inputView, q)
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

	// 删除视图（会话历史保留在 gui.aiChatSession 中）
	gui.g.DeleteView("aiChat")
	gui.g.DeleteView("aiChatInput")
	gui.g.DeleteView("aiChatStatus")

	return nil
}
