package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// openAIProvider implements Provider using the OpenAI chat completions API.
// Compatible with OpenAI, DeepSeek, Ollama (/v1), and other OpenAI-compatible services.
type openAIProvider struct {
	endpoint       string
	apiKey         string
	model          string
	maxTokens      int
	enableThinking bool
	client         *http.Client
	streamClient   *http.Client // no total timeout so SSE streams can run indefinitely
}

// thinkingConfig maps to DeepSeek's {"thinking": {"type": "enabled"}} request field.
type thinkingConfig struct {
	Type string `json:"type"`
}

type chatRequest struct {
	Model     string          `json:"model"`
	Messages  []chatMessage   `json:"messages"`
	MaxTokens int             `json:"max_tokens,omitempty"`
	// Thinking enables thinking mode for models that support it via parameter
	// (e.g. deepseek-chat). Omitted for models with native reasoning (deepseek-reasoner).
	Thinking  *thinkingConfig `json:"thinking,omitempty"`
}

type streamRequest struct {
	Model     string          `json:"model"`
	Messages  []chatMessage   `json:"messages"`
	MaxTokens int             `json:"max_tokens,omitempty"`
	Stream    bool            `json:"stream"`
	Thinking  *thinkingConfig `json:"thinking,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	// ReasoningContent carries the thinking chain in assistant messages.
	// Must be included when replying within the same turn during tool calls.
	// Omit (omitempty) when starting a new turn to save bandwidth.
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type streamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func newOpenAIProvider(endpoint, apiKey, model string, maxTokens, timeoutSecs int, enableThinking bool) *openAIProvider {
	return &openAIProvider{
		endpoint:       strings.TrimRight(endpoint, "/"),
		apiKey:         apiKey,
		model:          model,
		maxTokens:      maxTokens,
		enableThinking: enableThinking,
		client:         &http.Client{Timeout: time.Duration(timeoutSecs) * time.Second},
		streamClient:   &http.Client{Timeout: 0}, // no total timeout; rely on context cancellation
	}
}

// isReasonerModel reports whether the model has native reasoning built in
// (e.g. deepseek-reasoner). These models do not need the thinking parameter.
func isReasonerModel(model string) bool {
	return strings.Contains(strings.ToLower(model), "reasoner")
}

func (p *openAIProvider) Complete(ctx context.Context, prompt string) (Result, error) {
	req := chatRequest{
		Model:     p.model,
		Messages:  []chatMessage{{Role: "user", Content: prompt}},
		MaxTokens: p.maxTokens,
	}

	// Pass the thinking parameter only for non-reasoner models that need it.
	// deepseek-reasoner always thinks natively; the parameter is unnecessary there.
	if p.enableThinking && !isReasonerModel(p.model) {
		req.Thinking = &thinkingConfig{Type: "enabled"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return Result{}, fmt.Errorf("AI: failed to marshal request: %w", err)
	}

	url := p.endpoint + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Result{}, fmt.Errorf("AI: failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return Result{}, fmt.Errorf("AI: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, fmt.Errorf("AI: failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("AI: unexpected status %d: %s", resp.StatusCode, string(respBytes))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return Result{}, fmt.Errorf("AI: failed to parse response: %w", err)
	}

	if chatResp.Error != nil {
		return Result{}, fmt.Errorf("AI: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return Result{}, fmt.Errorf("AI: empty response from model")
	}

	msg := chatResp.Choices[0].Message
	return Result{
		Content:          strings.TrimSpace(msg.Content),
		ReasoningContent: msg.ReasoningContent,
	}, nil
}

// CompleteStream sends a prompt and streams the response via SSE, calling onChunk
// for each content fragment received. Uses context for cancellation.
func (p *openAIProvider) CompleteStream(ctx context.Context, prompt string, onChunk func(string)) error {
	req := streamRequest{
		Model:     p.model,
		Messages:  []chatMessage{{Role: "user", Content: prompt}},
		MaxTokens: p.maxTokens,
		Stream:    true,
	}

	if p.enableThinking && !isReasonerModel(p.model) {
		req.Thinking = &thinkingConfig{Type: "enabled"}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("AI: failed to marshal request: %w", err)
	}

	url := p.endpoint + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("AI: failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	if p.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.streamClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("AI: stream request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("AI: unexpected status %d: %s", resp.StatusCode, string(respBytes))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk streamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue // skip malformed chunks
		}
		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				onChunk(content)
			}
		}
	}

	return scanner.Err()
}
