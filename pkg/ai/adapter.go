package ai

import (
	"context"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

// legacyProviderAdapter wraps a provider.Provider (new []Message interface)
// to satisfy the old Provider interface (single prompt string).
// This lets existing *Client code continue working without modification
// during the incremental migration to Manager.
type legacyProviderAdapter struct {
	p provider.Provider
}

func (a *legacyProviderAdapter) Complete(ctx context.Context, prompt string) (Result, error) {
	msgs := []provider.Message{{Role: provider.RoleUser, Content: prompt}}
	res, err := a.p.Complete(ctx, msgs)
	if err != nil {
		return Result{}, err
	}
	return Result{Content: res.Content, ReasoningContent: res.ReasoningContent}, nil
}

func (a *legacyProviderAdapter) CompleteStream(ctx context.Context, prompt string, onChunk func(string)) error {
	msgs := []provider.Message{{Role: provider.RoleUser, Content: prompt}}
	return a.p.CompleteStream(ctx, msgs, onChunk)
}
