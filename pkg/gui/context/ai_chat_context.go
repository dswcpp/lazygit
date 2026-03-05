package context

import (
	"context"

	"github.com/dswcpp/lazygit/pkg/gui/types"
)

// AIChatContext 是 AI 聊天浮窗的上下文，遵循 lazygit Context 系统。
// 与 AICodeReviewContext 相同模式：TEMPORARY_POPUP + HasUncontrolledBounds。
type AIChatContext struct {
	*SimpleContext
	// Zoomed 切换默认尺寸（90% 屏幕）与全屏（96%）
	Zoomed bool
	// CancelFunc 取消正在进行的 AI 请求
	CancelFunc context.CancelFunc
}

func NewAIChatContext(c *ContextCommon) *AIChatContext {
	return &AIChatContext{
		SimpleContext: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
			View:                  c.Views().AIChat,
			WindowName:            "aiChat",
			Key:                   AI_CHAT_CONTEXT_KEY,
			Kind:                  types.TEMPORARY_POPUP,
			Focusable:             true,
			HasUncontrolledBounds: true,
		})),
	}
}
