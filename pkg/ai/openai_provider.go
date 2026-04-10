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

const (
	wireAPIChat      = "chat"
	wireAPIResponses = "responses"
)

// openAIProvider implements Provider using OpenAI-compatible HTTP APIs.
// It supports both the Chat Completions API and the Responses API.
type openAIProvider struct {
	endpoint       string
	apiKey         string
	model          string
	maxTokens      int
	enableThinking bool
	customHeaders  map[string]string
	wireAPI        string
	client         *http.Client
	streamClient   *http.Client // no total timeout so SSE streams can run indefinitely
}

type thinkingConfig struct {
	Type string `json:"type"`
}

type chatRequest struct {
	Model     string          `json:"model"`
	Messages  []chatMessage   `json:"messages"`
	MaxTokens int             `json:"max_tokens,omitempty"`
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
	Role             string `json:"role"`
	Content          string `json:"content"`
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

type responsesRequest struct {
	Model           string `json:"model"`
	Input           string `json:"input"`
	MaxOutputTokens int    `json:"max_output_tokens,omitempty"`
	Stream          bool   `json:"stream,omitempty"`
}

type responsesResponse struct {
	OutputText string `json:"output_text"`
	Output     []struct {
		Type    string `json:"type"`
		Role    string `json:"role"`
		Content []struct {
			Type    string `json:"type"`
			Text    string `json:"text,omitempty"`
			Refusal string `json:"refusal,omitempty"`
		} `json:"content"`
	} `json:"output"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type responsesStreamEvent struct {
	Type  string `json:"type"`
	Delta string `json:"delta,omitempty"`
	Text  string `json:"text,omitempty"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func newOpenAIProvider(endpoint, apiKey, model string, maxTokens, timeoutSecs int, enableThinking bool, customHeaders map[string]string, wireAPI string) *openAIProvider {
	return &openAIProvider{
		endpoint:       strings.TrimRight(endpoint, "/"),
		apiKey:         apiKey,
		model:          model,
		maxTokens:      maxTokens,
		enableThinking: enableThinking,
		customHeaders:  customHeaders,
		wireAPI:        normalizeWireAPI(wireAPI),
		client:         &http.Client{Timeout: time.Duration(timeoutSecs) * time.Second},
		streamClient:   &http.Client{Timeout: 0}, // no total timeout; rely on context cancellation
	}
}

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

func normalizeWireAPI(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case wireAPIResponses:
		return wireAPIResponses
	default:
		return wireAPIChat
	}
}

func isReasonerModel(model string) bool {
	return strings.Contains(strings.ToLower(model), "reasoner")
}

func (p *openAIProvider) Complete(ctx context.Context, prompt string) (Result, error) {
	if p.wireAPI == wireAPIResponses {
		return p.completeResponses(ctx, prompt)
	}
	return p.completeChat(ctx, prompt)
}

func (p *openAIProvider) completeChat(ctx context.Context, prompt string) (Result, error) {
	req := chatRequest{
		Model:     p.model,
		Messages:  []chatMessage{{Role: "user", Content: prompt}},
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

func (p *openAIProvider) completeResponses(ctx context.Context, prompt string) (Result, error) {
	req := responsesRequest{
		Model:           p.model,
		Input:           prompt,
		MaxOutputTokens: p.maxTokens,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return Result{}, fmt.Errorf("AI: failed to marshal request: %w", err)
	}

	url := p.endpoint + "/responses"
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
		if shouldRetryResponsesAsStream(resp.StatusCode, respBytes) {
			return p.completeResponsesViaStream(ctx, prompt)
		}
		return Result{}, fmt.Errorf("AI: unexpected status %d: %s", resp.StatusCode, string(respBytes))
	}

	return parseResponsesResult(respBytes)
}

func parseResponsesResult(respBytes []byte) (Result, error) {
	var responsesResp responsesResponse
	if err := json.Unmarshal(respBytes, &responsesResp); err != nil {
		return Result{}, fmt.Errorf("AI: failed to parse response: %w", err)
	}

	if responsesResp.Error != nil {
		return Result{}, fmt.Errorf("AI: %s", responsesResp.Error.Message)
	}

	content := strings.TrimSpace(responsesResp.OutputText)
	if content == "" {
		content = strings.TrimSpace(extractResponsesOutputText(responsesResp.Output))
	}
	if content == "" && len(responsesResp.Output) == 0 {
		return Result{}, fmt.Errorf("AI: empty response from model")
	}

	return Result{Content: content}, nil
}

func shouldRetryResponsesAsStream(statusCode int, respBytes []byte) bool {
	if statusCode != http.StatusBadRequest {
		return false
	}
	return strings.Contains(strings.ToLower(string(respBytes)), "stream must be set to true")
}

func (p *openAIProvider) completeResponsesViaStream(ctx context.Context, prompt string) (Result, error) {
	var content strings.Builder
	if err := p.completeResponsesStream(ctx, prompt, func(chunk string) {
		content.WriteString(chunk)
	}); err != nil {
		return Result{}, err
	}

	if content.Len() == 0 {
		return Result{}, fmt.Errorf("AI: empty response from model")
	}

	return Result{Content: strings.TrimSpace(content.String())}, nil
}

func extractResponsesOutputText(output []struct {
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type    string `json:"type"`
		Text    string `json:"text,omitempty"`
		Refusal string `json:"refusal,omitempty"`
	} `json:"content"`
}) string {
	var parts []string
	for _, item := range output {
		if item.Role != "" && item.Role != "assistant" {
			continue
		}
		for _, content := range item.Content {
			switch content.Type {
			case "output_text", "text":
				if text := strings.TrimSpace(content.Text); text != "" {
					parts = append(parts, text)
				}
			case "refusal":
				if refusal := strings.TrimSpace(content.Refusal); refusal != "" {
					parts = append(parts, refusal)
				}
			}
		}
	}
	return strings.Join(parts, "\n")
}

// CompleteStream sends a prompt and streams the response via SSE, calling onChunk
// for each content fragment received. Uses context for cancellation.
func (p *openAIProvider) CompleteStream(ctx context.Context, prompt string, onChunk func(string)) error {
	if p.wireAPI == wireAPIResponses {
		return p.completeResponsesStream(ctx, prompt, onChunk)
	}
	return p.completeChatStream(ctx, prompt, onChunk)
}

func (p *openAIProvider) completeChatStream(ctx context.Context, prompt string, onChunk func(string)) error {
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

func (p *openAIProvider) completeResponsesStream(ctx context.Context, prompt string, onChunk func(string)) error {
	req := responsesRequest{
		Model:           p.model,
		Input:           prompt,
		MaxOutputTokens: p.maxTokens,
		Stream:          true,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("AI: failed to marshal request: %w", err)
	}

	url := p.endpoint + "/responses"
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

	if !strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "text/event-stream") {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("AI: failed to read response: %w", err)
		}
		result, err := parseResponsesResult(respBytes)
		if err != nil {
			return err
		}
		if result.Content != "" {
			onChunk(result.Content)
		}
		return nil
	}

	scanner := bufio.NewScanner(resp.Body)
	sawDelta := false
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var event responsesStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		switch event.Type {
		case "response.output_text.delta":
			if event.Delta != "" {
				sawDelta = true
				onChunk(event.Delta)
			}
		case "response.output_text.done":
			if !sawDelta && event.Text != "" {
				onChunk(event.Text)
			}
		case "error", "response.error", "response.failed":
			if event.Error != nil && event.Error.Message != "" {
				return fmt.Errorf("AI: %s", event.Error.Message)
			}
			if event.Text != "" {
				return fmt.Errorf("AI: %s", event.Text)
			}
		}
	}

	return scanner.Err()
}
