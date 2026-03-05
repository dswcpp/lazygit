package provider

import "context"

// Role identifies who authored a message in a conversation.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message is a single turn in a conversation.
type Message struct {
	Role    Role
	Content string
}

// Result holds the AI response, including an optional reasoning chain.
type Result struct {
	Content          string
	ReasoningContent string
}

// Provider is the interface for AI backend implementations.
// All methods accept a full conversation history to support multi-turn dialogue.
type Provider interface {
	// Complete sends the conversation and returns the full response.
	Complete(ctx context.Context, messages []Message) (Result, error)
	// CompleteStream sends the conversation and streams the response chunk by chunk.
	CompleteStream(ctx context.Context, messages []Message, onChunk func(string)) error
	// ModelID returns the model identifier string.
	ModelID() string
}
