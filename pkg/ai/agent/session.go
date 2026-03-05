package agent

import (
	"fmt"
	"strings"
	"time"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

// MessageKind classifies a message for UI rendering purposes.
type MessageKind string

const (
	KindSystem    MessageKind = "system"
	KindUser      MessageKind = "user"
	KindAssistant MessageKind = "assistant"
	KindToolCall  MessageKind = "tool_call"   // agent called a tool
	KindToolResult MessageKind = "tool_result" // tool execution result
	KindError     MessageKind = "error"
)

// UIMessage is a displayable record stored in the session.
// It is separate from provider.Message to allow richer UI metadata.
type UIMessage struct {
	Kind          MessageKind
	Content       string
	Timestamp     time.Time
	// ToolName is set for KindToolCall and KindToolResult.
	ToolName      string
	// ToolSuccess is set for KindToolResult.
	ToolSuccess   bool
}

// Session manages the full conversation state for one agent dialogue.
// It maintains two parallel views:
//   - UIMessages: for rendering to the user (includes tool call details, errors…)
//   - providerMessages: the conversation history sent to the LLM provider
type Session struct {
	systemPrompt     string
	UIMessages       []UIMessage
	providerMessages []provider.Message
}

// NewSession creates a new Session with the given system prompt.
func NewSession(systemPrompt string) *Session {
	return &Session{systemPrompt: systemPrompt}
}

// AddUserMessage records a user turn in both views.
func (s *Session) AddUserMessage(content string) {
	s.UIMessages = append(s.UIMessages, UIMessage{
		Kind: KindUser, Content: content, Timestamp: time.Now(),
	})
	s.providerMessages = append(s.providerMessages, provider.Message{
		Role: provider.RoleUser, Content: content,
	})
}

// AddAssistantMessage records an assistant turn.
func (s *Session) AddAssistantMessage(content string) {
	if content == "" {
		return
	}
	s.UIMessages = append(s.UIMessages, UIMessage{
		Kind: KindAssistant, Content: content, Timestamp: time.Now(),
	})
	s.providerMessages = append(s.providerMessages, provider.Message{
		Role: provider.RoleAssistant, Content: content,
	})
}

// AddSystemNote adds a system-level note visible to the user but not sent to the LLM.
func (s *Session) AddSystemNote(content string) {
	s.UIMessages = append(s.UIMessages, UIMessage{
		Kind: KindSystem, Content: content, Timestamp: time.Now(),
	})
}

// AddErrorMessage records an error in the UI view only.
func (s *Session) AddErrorMessage(content string) {
	s.UIMessages = append(s.UIMessages, UIMessage{
		Kind: KindError, Content: content, Timestamp: time.Now(),
	})
}

// AddToolCall records a tool invocation in the UI view.
func (s *Session) AddToolCall(call tools.ToolCall) {
	s.UIMessages = append(s.UIMessages, UIMessage{
		Kind:      KindToolCall,
		Content:   fmt.Sprintf("调用工具: %s %v", call.Name, call.Params),
		ToolName:  call.Name,
		Timestamp: time.Now(),
	})
}

// AddToolResult records a tool result in both views.
// The provider sees tool results as a user-role message so it can reason about them.
func (s *Session) AddToolResult(result tools.ToolResult, toolName string) {
	s.UIMessages = append(s.UIMessages, UIMessage{
		Kind:        KindToolResult,
		Content:     result.Output,
		ToolName:    toolName,
		ToolSuccess: result.Success,
		Timestamp:   time.Now(),
	})
	// Feed result back to the LLM as a user message with structured context.
	status := "成功"
	if !result.Success {
		status = "失败"
	}
	s.providerMessages = append(s.providerMessages, provider.Message{
		Role:    provider.RoleUser,
		Content: fmt.Sprintf("[工具结果 %s - %s]\n%s", toolName, status, result.Output),
	})
}

// ProviderMessages returns the full message history ready for the provider,
// with the system prompt prepended as the first message.
func (s *Session) ProviderMessages() []provider.Message {
	if s.systemPrompt == "" {
		return s.providerMessages
	}
	out := make([]provider.Message, 0, len(s.providerMessages)+1)
	out = append(out, provider.Message{Role: provider.RoleSystem, Content: s.systemPrompt})
	out = append(out, s.providerMessages...)
	return out
}

// LastAssistantContent returns the content of the most recent assistant message.
func (s *Session) LastAssistantContent() string {
	for i := len(s.UIMessages) - 1; i >= 0; i-- {
		if s.UIMessages[i].Kind == KindAssistant {
			return s.UIMessages[i].Content
		}
	}
	return ""
}

// Reset clears all conversation history but preserves the system prompt.
func (s *Session) Reset() {
	s.UIMessages = nil
	s.providerMessages = nil
}

// Summary returns a compact text summary of the session for debugging.
func (s *Session) Summary() string {
	var sb strings.Builder
	for _, m := range s.UIMessages {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", m.Kind, m.Content))
	}
	return sb.String()
}
