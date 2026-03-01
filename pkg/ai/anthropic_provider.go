package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// anthropicProvider implements the Provider interface using the official
// Anthropic Go SDK. It supports claude-* model families.
type anthropicProvider struct {
	client    anthropic.Client
	model     string
	maxTokens int64
	timeout   time.Duration
}

func newAnthropicProvider(apiKey, model, baseURL string, maxTokens, timeoutSecs int) *anthropicProvider {
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	return &anthropicProvider{
		client:    anthropic.NewClient(opts...),
		model:     model,
		maxTokens: int64(maxTokens),
		timeout:   time.Duration(timeoutSecs) * time.Second,
	}
}

// Complete sends a prompt and waits for the full response.
func (p *anthropicProvider) Complete(ctx context.Context, prompt string) (Result, error) {
	msg, err := p.client.Messages.New(ctx,
		anthropic.MessageNewParams{
			Model:     anthropic.Model(p.model),
			MaxTokens: p.maxTokens,
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
			},
		},
		option.WithRequestTimeout(p.timeout),
	)
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

// CompleteStream sends a prompt and streams the response via onChunk callbacks.
func (p *anthropicProvider) CompleteStream(ctx context.Context, prompt string, onChunk func(string)) error {
	stream := p.client.Messages.NewStreaming(ctx,
		anthropic.MessageNewParams{
			Model:     anthropic.Model(p.model),
			MaxTokens: p.maxTokens,
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
			},
		},
	)
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
