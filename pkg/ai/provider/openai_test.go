package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestOpenAIProviderCompleteUsesResponsesAPI(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/responses" {
			t.Fatalf("expected /responses, got %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("unexpected auth header: %q", got)
		}

		var req struct {
			Model           string `json:"model"`
			MaxOutputTokens int    `json:"max_output_tokens"`
			Input           []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"input"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		if req.Model != "gpt-5.4" {
			t.Fatalf("unexpected model: %s", req.Model)
		}
		if req.MaxOutputTokens != 128 {
			t.Fatalf("unexpected max_output_tokens: %d", req.MaxOutputTokens)
		}
		if len(req.Input) != 2 || req.Input[0].Role != "system" || req.Input[1].Role != "user" {
			t.Fatalf("unexpected input payload: %+v", req.Input)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"output_text":"OK","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"OK"}]}]}`))
	}))
	defer server.Close()

	provider := NewOpenAIProvider(server.URL, "test-key", "gpt-5.4", 128, 30, false, nil, "responses")
	result, err := provider.Complete(context.Background(), []Message{
		{Role: RoleSystem, Content: "sys"},
		{Role: RoleUser, Content: "hello"},
	})
	if err != nil {
		t.Fatalf("Complete returned error: %v", err)
	}
	if result.Content != "OK" {
		t.Fatalf("unexpected content: %q", result.Content)
	}
}

func TestOpenAIProviderCompleteStreamUsesResponsesAPI(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/responses" {
			t.Fatalf("expected /responses, got %s", r.URL.Path)
		}

		var req struct {
			Stream bool `json:"stream"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if !req.Stream {
			t.Fatalf("expected stream=true request")
		}

		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("event: response.output_text.delta\n"))
		_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"O\"}\n\n"))
		_, _ = w.Write([]byte("event: response.output_text.delta\n"))
		_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"K\"}\n\n"))
		_, _ = w.Write([]byte("event: response.completed\n"))
		_, _ = w.Write([]byte("data: {\"type\":\"response.completed\"}\n\n"))
	}))
	defer server.Close()

	provider := NewOpenAIProvider(server.URL, "test-key", "gpt-5.4", 128, 30, false, nil, "responses")
	var out strings.Builder
	err := provider.CompleteStream(context.Background(), []Message{
		{Role: RoleUser, Content: "hello"},
	}, func(chunk string) {
		out.WriteString(chunk)
	})
	if err != nil {
		t.Fatalf("CompleteStream returned error: %v", err)
	}
	if got := out.String(); got != "OK" {
		t.Fatalf("unexpected streamed content: %q", got)
	}
}

func TestOpenAIProviderDefaultsToChatCompletions(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("expected /chat/completions, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"chat-ok"}}]}`))
	}))
	defer server.Close()

	provider := NewOpenAIProvider(server.URL, "test-key", "gpt-4o-mini", 64, 30, false, nil, "")
	result, err := provider.Complete(context.Background(), []Message{
		{Role: RoleUser, Content: "hello"},
	})
	if err != nil {
		t.Fatalf("Complete returned error: %v", err)
	}
	if result.Content != "chat-ok" {
		t.Fatalf("unexpected content: %q", result.Content)
	}
}

func TestOpenAIProviderCompleteRetriesResponsesAsStreamWhenRequired(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/responses" {
			t.Fatalf("expected /responses, got %s", r.URL.Path)
		}

		var req struct {
			Stream bool `json:"stream"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		mu.Lock()
		requestCount++
		current := requestCount
		mu.Unlock()

		if current == 1 {
			if req.Stream {
				t.Fatalf("expected first request to be non-streaming")
			}
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"detail":"Stream must be set to true"}`))
			return
		}

		if !req.Stream {
			t.Fatalf("expected retry request to be streaming")
		}

		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"O\"}\n\n"))
		_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"K\"}\n\n"))
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer server.Close()

	provider := NewOpenAIProvider(server.URL, "test-key", "gpt-5.4", 128, 30, false, nil, "responses")
	result, err := provider.Complete(context.Background(), []Message{
		{Role: RoleUser, Content: "hello"},
	})
	if err != nil {
		t.Fatalf("Complete returned error: %v", err)
	}
	if result.Content != "OK" {
		t.Fatalf("unexpected content: %q", result.Content)
	}

	mu.Lock()
	defer mu.Unlock()
	if requestCount != 2 {
		t.Fatalf("expected 2 requests, got %d", requestCount)
	}
}
