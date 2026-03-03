package context

import (
	"context"

	"github.com/dswcpp/lazygit/pkg/gui/types"
)

type AICodeReviewContext struct {
	*SimpleContext
	// Zoomed toggles between the normal size (90% screen) and maximised (full screen minus margin).
	Zoomed bool
	// CancelFunc cancels the ongoing AI request. Set by AI helpers during streaming.
	CancelFunc context.CancelFunc
}

func NewAICodeReviewContext(c *ContextCommon) *AICodeReviewContext {
	return &AICodeReviewContext{
		SimpleContext: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
			View:                  c.Views().AICodeReview,
			WindowName:            "aiCodeReview",
			Key:                   AI_CODE_REVIEW_CONTEXT_KEY,
			Kind:                  types.TEMPORARY_POPUP,
			Focusable:             true,
			HasUncontrolledBounds: true,
		})),
	}
}
