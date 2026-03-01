package ai

import "context"

// Result holds the AI response, including the optional reasoning chain.
type Result struct {
	// Content is the final answer from the model.
	Content string
	// ReasoningContent is the thinking chain produced by reasoning models
	// (e.g. deepseek-reasoner). Empty for standard chat models.
	ReasoningContent string
}

// Provider is the interface for AI backend implementations.
type Provider interface {
	// Complete sends a prompt and returns the AI result.
	Complete(ctx context.Context, prompt string) (Result, error)
}
