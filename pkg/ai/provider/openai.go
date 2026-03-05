package provider

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
	customHeaders  map[string]string
	client         *http.Client
	streamClient   *http.Client
}

type thinkingConfig struct {
	Type string `json:"type"`
}

type openAIChatMessage struct {
	Role             string          `json:"role"`
	Content          string          `json:"content"`
	ReasoningContent string          `json:"reasoning_content,omitempty"`
}

type openAIChatRequest struct {
	Model     string               `json:"model"`
	Messages  []openAIChatMessage  `json:"messages"`
	MaxTokens int                  `json:"max_tokens,omitempty"`
	Thinking  *thinkingConfig      `json:"thinking,omitempty"`
}

type openAIStreamRequest struct {
	Model     string               `json:"model"`
	Messages  []openAIChatMessage  `json:"messages"`
	MaxTokens int                  `json:"max_tokens,omitempty"`
	Stream    bool                 `json:"stream"`
	Thinking  *thinkingConfig      `json:"thinking,omitempty"`
}

type openAIChatResponse struct {
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

type openAIStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

// NewOpenAIProvider creates an OpenAI-compatible provider.
func NewOpenAIProvider(endpoint, apiKey, model string, maxTokens, timeoutSecs int, enableThinking bool, customHeaders map[string]string) Provider {
	return &openAIProvider{
		endpoint:       strings.TrimRight(endpoint, "/"),
		apiKey:         apiKey,
		model:          model,
		maxTokens:      maxTokens,
		enableThinking: enableThinking,
		customHeaders:  customHeaders,
		client:         &http.Client{Timeout: time.Duration(timeoutSecs) * time.Second},
		streamClient:   &http.Client{Timeout: 0},
	}
}

func (p *openAIProvider) ModelID() string { return p.model }

func (p *openAIProvider) applyHeaders(req *http.Request, acceptSSE bool) {
	req.Header.Set("Content-Type", "application/json")
	if acceptSSE {
		req.Header.Set("Accept", "text/event-stream")
	}
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}
	for k, v := range p.customHeaders {
		req.Header.Set(k, v)
	}
}

func isReasonerModel(model string) bool {
	return strings.Contains(strings.ToLower(model), "reasoner")
}

func toOpenAIMessages(messages []Message) []openAIChatMessage {
	out := make([]openAIChatMessage, 0, len(messages))
	for _, m := range messages {
		out = append(out, openAIChatMessage{Role: string(m.Role), Content: m.Content})
	}
	return out
}

func (p *openAIProvider) Complete(ctx context.Context, messages []Message) (Result, error) {
	req := openAIChatRequest{
		Model:     p.model,
		Messages:  toOpenAIMessages(messages),
		MaxTokens: p.maxTokens,
	}
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
	p.applyHeaders(httpReq, false)

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

	var chatResp openAIChatResponse
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

func (p *openAIProvider) CompleteStream(ctx context.Context, messages []Message, onChunk func(string)) error {
	req := openAIStreamRequest{
		Model:     p.model,
		Messages:  toOpenAIMessages(messages),
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
	p.applyHeaders(httpReq, true)

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
		var chunk openAIStreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) > 0 {
			if content := chunk.Choices[0].Delta.Content; content != "" {
				onChunk(content)
			}
		}
	}
	return scanner.Err()
}
