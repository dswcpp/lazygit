package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// anthropicProvider implements Provider using the official Anthropic Go SDK.
type anthropicProvider struct {
	client    anthropic.Client
	model     string
	maxTokens int64
	timeout   time.Duration
}

// NewAnthropicProvider creates an Anthropic provider.
func NewAnthropicProvider(apiKey, model, baseURL string, maxTokens, timeoutSecs int, customHeaders map[string]string) Provider {
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	for k, v := range customHeaders {
		opts = append(opts, option.WithHeader(k, v))
	}
	return &anthropicProvider{
		client:    anthropic.NewClient(opts...),
		model:     model,
		maxTokens: int64(maxTokens),
		timeout:   time.Duration(timeoutSecs) * time.Second,
	}
}

func (p *anthropicProvider) ModelID() string { return p.model }

// toAnthropicMessages converts []Message to Anthropic MessageParams.
// System messages are returned separately (Anthropic puts them in a dedicated field).
func toAnthropicMessages(messages []Message) (system []anthropic.TextBlockParam, turns []anthropic.MessageParam) {
	for _, m := range messages {
		switch m.Role {
		case RoleSystem:
			system = append(system, anthropic.TextBlockParam{Text: m.Content})
		case RoleUser:
			turns = append(turns, anthropic.NewUserMessage(anthropic.NewTextBlock(m.Content)))
		case RoleAssistant:
			turns = append(turns, anthropic.NewAssistantMessage(anthropic.NewTextBlock(m.Content)))
		}
	}
	return system, turns
}

func (p *anthropicProvider) Complete(ctx context.Context, messages []Message) (Result, error) {
	system, turns := toAnthropicMessages(messages)
	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(p.model),
		MaxTokens: p.maxTokens,
		Messages:  turns,
	}
	if len(system) > 0 {
		params.System = system
	}

	msg, err := p.client.Messages.New(ctx, params, option.WithRequestTimeout(p.timeout))
	if err != nil {
		return Result{}, fmt.Errorf("anthropic: %w", err)
	}

	var sb strings.Builder
	for _, block := range msg.Content {
		if block.Type == "text" {
			sb.WriteString(block.Text)
		}
	}
	return Result{Content: sb.String()}, nil
}

func (p *anthropicProvider) CompleteStream(ctx context.Context, messages []Message, onChunk func(string)) error {
	system, turns := toAnthropicMessages(messages)
	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(p.model),
		MaxTokens: p.maxTokens,
		Messages:  turns,
	}
	if len(system) > 0 {
		params.System = system
	}

	stream := p.client.Messages.NewStreaming(ctx, params)
	defer stream.Close()

	for stream.Next() {
		event := stream.Current()
		if event.Type == "content_block_delta" && event.Delta.Type == "text_delta" {
			onChunk(event.Delta.Text)
		}
	}
	if err := stream.Err(); err != nil {
		return fmt.Errorf("anthropic stream: %w", err)
	}
	return nil
}
