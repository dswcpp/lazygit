package gui

import (
	"strings"

	chathelpers "github.com/dswcpp/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/gocui"
)

// AI 聊天功能现已迁移到 Context 系统，通过 AIChatContext + AIChatController 管理。
// 相关实现位于：
//   - pkg/gui/context/ai_chat_context.go        (Context 定义)
//   - pkg/gui/controllers/ai_chat_controller.go (键盘绑定)
//   - pkg/gui/controllers/helpers/ai_chat_helper.go (会话管理、渲染、AI 调用)

// ShowAIChat 显示 AI 聊天弹窗（保留历史会话）
func (gui *Gui) ShowAIChat() error {
	return gui.helpers.AIChat.ShowChat()
}

// ShowAIChatWithFollowUp 携带上下文内容打开 AI 聊天，用于从其他面板继续对话
func (gui *Gui) ShowAIChatWithFollowUp(contextContent string) error {
	return gui.helpers.AIChat.ShowChatWithContext(contextContent)
}

// aiChatInputEditor 处理 AIChatInput 视图的键盘输入。
//   - Enter：发送消息，清空输入框
//   - Esc / Ctrl+C：关闭聊天弹窗
//   - 其他键：交给 gocui 默认编辑器（字符插入、Backspace、方向键等）
func (gui *Gui) aiChatInputEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	switch key {
	case gocui.KeyEnter:
		content := strings.TrimSpace(v.Buffer())
		if content == "" {
			return true
		}
		chathelpers.ResetAIChatInputView(v)
		// 在 UI goroutine 外发送（helpers.AIChat.SendMessage 内部会用 goroutine）
		_ = gui.helpers.AIChat.SendMessage(content)
		return true

	case gocui.KeyTab:
		_, _ = gui.g.SetCurrentView(gui.Views.AIChat.Name())
		return true

	case gocui.KeyEsc, gocui.KeyCtrlC:
		gui.c.Context().Pop()
		return true

	default:
		return gocui.DefaultEditor.Edit(v, key, ch, mod)
	}
}
