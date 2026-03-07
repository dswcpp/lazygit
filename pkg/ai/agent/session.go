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
	KindSystem     MessageKind = "system"
	KindUser       MessageKind = "user"
	KindAssistant  MessageKind = "assistant"
	KindToolCall   MessageKind = "tool_call"    // agent called a tool
	KindToolResult MessageKind = "tool_result"  // tool execution result
	KindError      MessageKind = "error"
	KindPlan       MessageKind = "plan"         // 完整执行计划（TwoPhaseAgent 阶段一输出）
	KindStepUpdate MessageKind = "step_update"  // 单步执行状态更新（TwoPhaseAgent 阶段二）
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

// Session manages the conversation state for one agent dialogue.
// It maintains two parallel views:
//   - UIMessages: for rendering to the user (includes tool call details, errors…)
//   - providerMessages: the conversation history sent to the LLM provider
//
// Session is intentionally free of agent control-flow state (Phase, Plan).
// Those live in GraphState, which TwoPhaseAgent owns exclusively.
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
		streamingMessageIndex: -1, // 初始无流式消息
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
// 返回消息索引，用于后续追加内容。
func (s *Session) StartStreamingMessage() int {
	s.streamingBuffer.Reset()
	s.UIMessages = append(s.UIMessages, UIMessage{
		Kind: KindAssistant, Content: "", Timestamp: time.Now(),
	})
	s.streamingMessageIndex = len(s.UIMessages) - 1
	return s.streamingMessageIndex
}

// AppendToStreamingMessage 向当前流式消息追加内容。
// 必须先调用 StartStreamingMessage。
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

// AddPlanUIMessage adds a KindPlan UIMessage for the given plan.
// The plan itself is stored in GraphState — Session only records the UI event.
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
//
// Responsibility boundary: Session is pure UI — it does NOT mutate step fields.
// The caller (agent code in two_phase_agent.go) is responsible for setting
// step.Status, step.Result and step.Error before calling this method.
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
// Agent control-flow state (Phase, Plan) is managed by GraphState.Reset().
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
