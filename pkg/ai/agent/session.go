package agent

import (
	"fmt"
	"strings"
	"time"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

// Session manages the conversation state for one agent dialogue.
// DEPRECATED for TwoPhaseAgent: All state now lives in GraphState.
// Session is kept for backward compatibility with the legacy Agent implementation.
type Session struct {
	systemPrompt     string
	UIMessages       []UIMessage
	providerMessages []provider.Message

	// 流式输出状态
	streamingMessageIndex int             // 正在流式输出的消息索引（-1 表示无流式消息）
	streamingBuffer       strings.Builder // 流式消息缓冲区
}

// NewSession creates a new Session with the given system prompt.
func NewSession(systemPrompt string) *Session {
	return &Session{
		systemPrompt:          systemPrompt,
		streamingMessageIndex: -1,
	}
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

// StartStreamingMessage 开始一条新的流式 Assistant 消息。
func (s *Session) StartStreamingMessage() int {
	s.streamingBuffer.Reset()
	s.UIMessages = append(s.UIMessages, UIMessage{
		Kind: KindAssistant, Content: "", Timestamp: time.Now(),
	})
	s.streamingMessageIndex = len(s.UIMessages) - 1
	return s.streamingMessageIndex
}

// AppendToStreamingMessage 向当前流式消息追加内容。
func (s *Session) AppendToStreamingMessage(chunk string) {
	if s.streamingMessageIndex < 0 || s.streamingMessageIndex >= len(s.UIMessages) {
		return
	}
	s.streamingBuffer.WriteString(chunk)
	s.UIMessages[s.streamingMessageIndex].Content = s.streamingBuffer.String()
}

// FinishStreamingMessage 完成流式消息，将其加入 provider 消息历史。
func (s *Session) FinishStreamingMessage() {
	if s.streamingMessageIndex < 0 {
		return
	}
	content := s.streamingBuffer.String()
	if content != "" {
		s.providerMessages = append(s.providerMessages, provider.Message{
			Role: provider.RoleAssistant, Content: content,
		})
	}
	s.streamingMessageIndex = -1
	s.streamingBuffer.Reset()
}

// IsStreaming 返回当前是否有流式消息正在进行。
func (s *Session) IsStreaming() bool {
	return s.streamingMessageIndex >= 0
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
func (s *Session) AddToolResult(result tools.ToolResult, toolName string) {
	s.UIMessages = append(s.UIMessages, UIMessage{
		Kind:        KindToolResult,
		Content:     result.Output,
		ToolName:    toolName,
		ToolSuccess: result.Success,
		Timestamp:   time.Now(),
	})
	status := "成功"
	if !result.Success {
		status = "失败"
	}
	s.providerMessages = append(s.providerMessages, provider.Message{
		Role:    provider.RoleUser,
		Content: fmt.Sprintf("[工具结果 %s - %s]\n%s", toolName, status, result.Output),
	})
}

// ProviderMessages returns the full message history ready for the provider.
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

// AddPlanUIMessage adds a KindPlan UIMessage for the given plan.
func (s *Session) AddPlanUIMessage(plan *ExecutionPlan) {
	var sb strings.Builder
	sb.WriteString(plan.Summary)
	for _, step := range plan.Steps {
		sb.WriteString(fmt.Sprintf("\n  %s. %s", step.ID, step.Description))
	}
	s.UIMessages = append(s.UIMessages, UIMessage{
		Kind:      KindPlan,
		Content:   sb.String(),
		Timestamp: time.Now(),
	})
}

// AddStepUpdate records a step execution event in UIMessages.
func (s *Session) AddStepUpdate(step *PlanStep, status StepStatus, result, errMsg string) {
	content := fmt.Sprintf("[%s] %s", status, step.Description)
	if result != "" {
		content += "\n" + result
	}
	if errMsg != "" {
		content += "\n错误: " + errMsg
	}
	s.UIMessages = append(s.UIMessages, UIMessage{
		Kind:        KindStepUpdate,
		Content:     content,
		ToolName:    step.ToolName,
		ToolSuccess: status == StepDone,
		Timestamp:   time.Now(),
	})
}

// Reset clears all conversation history but preserves the system prompt.
func (s *Session) Reset() {
	s.UIMessages = nil
	s.providerMessages = nil
	s.streamingMessageIndex = -1
	s.streamingBuffer.Reset()
}

// Summary returns a compact text summary of the session for debugging.
func (s *Session) Summary() string {
	var sb strings.Builder
	for _, m := range s.UIMessages {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", m.Kind, m.Content))
	}
	return sb.String()
}
