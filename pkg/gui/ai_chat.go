package gui

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
