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
	KindToolCall   MessageKind = "tool_call"
	KindToolResult MessageKind = "tool_result"
	KindError      MessageKind = "error"
	KindPlan       MessageKind = "plan"
	KindStepUpdate MessageKind = "step_update"
)

// UIMessage is a displayable record for UI rendering.
type UIMessage struct {
	Kind        MessageKind
	Content     string
	Timestamp   time.Time
	ToolName    string // set for KindToolCall and KindToolResult
	ToolSuccess bool   // set for KindToolResult
}

// GraphState holds all state needed to describe, checkpoint, and resume a
// TwoPhaseAgent execution. This is the single source of truth for both
// control-flow state AND UI state.
//
// LangGraph analogy: this is the typed State object that flows through every
// node of the graph. All state mutations happen through node return values.
type GraphState struct {
	// ── Control Flow State ──────────────────────────────────────────────────

	// Phase is the current execution phase.
	Phase AgentPhase

	// Plan is the execution plan produced at the end of the planning phase.
	// Nil until the LLM has committed to a concrete plan.
	Plan *ExecutionPlan

	// PlanMessages is the LLM conversation history for the planning phase.
	// It carries a dedicated system prompt and may be pruned independently
	// of ProviderMessages.
	PlanMessages []provider.Message

	// ToolCallHistory tracks the number of times each unique (tool, params)
	// combination has been invoked during planning, used to break infinite loops.
	ToolCallHistory map[string]int

	// PlanStepCount is the number of LLM calls made so far in the current plan
	// loop. Compared against TwoPhaseAgent.maxPlanSteps to prevent runaway loops.
	PlanStepCount int

	// EmptyResponseCount tracks consecutive LLM responses that contained neither
	// a plan block nor tool calls. Three in a row aborts the planning loop.
	EmptyResponseCount int

	// PendingToolCalls holds tool calls produced by nodePlan and consumed by
	// nodeCallTools before control returns to nodePlan.
	PendingToolCalls []tools.ToolCall

	// ExecStepIndex is the index of the next PlanStep to execute.
	// nodeExecuteStep increments this after each successful step dispatch.
	ExecStepIndex int

	// ResumeFrom is set by nodeWaitHuman when the graph suspends at a human
	// interrupt point. A non-empty value means there is a pending checkpoint;
	// the next Send() call will resume the graph from this node.
	//
	// LangGraph analogy: this is the "thread checkpoint" that enables
	// human-in-the-loop resumption without re-running the planning phase.
	ResumeFrom NodeID

	// HumanInput is the user message provided to resume a suspended graph.
	// nodeHandleConfirmation reads it and clears it after routing.
	HumanInput string

	// ── UI State (LangGraph compliance: all state in GraphState) ────────────

	// UIMessages is the displayable conversation history for UI rendering.
	// Nodes append to this slice through return values (immutable pattern).
	UIMessages []UIMessage

	// ProviderMessages is the conversation history sent to the LLM provider.
	// Separate from UIMessages to allow different formatting.
	ProviderMessages []provider.Message

	// StreamingMessageIndex tracks the index of the message being streamed.
	// -1 means no streaming in progress.
	StreamingMessageIndex int

	// StreamingBuffer accumulates chunks during streaming.
	StreamingBuffer strings.Builder
}

// Reset clears all state for a fresh run.
func (s *GraphState) Reset() {
	// Control flow state
	s.Phase = PhasePlanning
	s.Plan = nil
	s.PlanMessages = nil
	s.ToolCallHistory = make(map[string]int)
	s.PlanStepCount = 0
	s.EmptyResponseCount = 0
	s.PendingToolCalls = nil
	s.ExecStepIndex = 0
	s.ResumeFrom = ""
	s.HumanInput = ""

	// UI state
	s.UIMessages = nil
	s.ProviderMessages = nil
	s.StreamingMessageIndex = -1
	s.StreamingBuffer.Reset()
}

// AppendPlanMessage appends msg to PlanMessages, trimming history to max entries
// while always preserving the first entry (the system prompt).
// Returns a new state with the updated PlanMessages (immutable).
func (s GraphState) AppendPlanMessage(msg provider.Message, max int) GraphState {
	// Deep copy PlanMessages
	newMessages := make([]provider.Message, len(s.PlanMessages), len(s.PlanMessages)+1)
	copy(newMessages, s.PlanMessages)
	newMessages = append(newMessages, msg)

	// Trim if necessary
	if len(newMessages) > max {
		newMessages = append(
			newMessages[:1],
			newMessages[len(newMessages)-max+1:]...,
		)
	}

	s.PlanMessages = newMessages
	return s
}

// WithToolCallIncrement returns a new state with the tool call count incremented.
// Pure function: deep copies the map to ensure immutability.
func (s GraphState) WithToolCallIncrement(callKey string) GraphState {
	// Deep copy the map
	newHistory := make(map[string]int, len(s.ToolCallHistory)+1)
	for k, v := range s.ToolCallHistory {
		newHistory[k] = v
	}
	newHistory[callKey]++
	s.ToolCallHistory = newHistory
	return s
}

// WithPhase returns a new state with the phase updated.
func (s GraphState) WithPhase(phase AgentPhase) GraphState {
	s.Phase = phase
	return s
}

// WithPlan returns a new state with the plan updated.
func (s GraphState) WithPlan(plan *ExecutionPlan) GraphState {
	s.Plan = plan
	return s
}

// WithPlanMessages returns a new state with plan messages replaced.
func (s GraphState) WithPlanMessages(messages []provider.Message) GraphState {
	s.PlanMessages = messages
	return s
}

// WithResumeFrom returns a new state with the resume node updated.
func (s GraphState) WithResumeFrom(nodeID NodeID) GraphState {
	s.ResumeFrom = nodeID
	return s
}

// WithHumanInput returns a new state with the human input updated.
func (s GraphState) WithHumanInput(input string) GraphState {
	s.HumanInput = input
	return s
}

// WithExecStepIndex returns a new state with the execution step index updated.
func (s GraphState) WithExecStepIndex(index int) GraphState {
	s.ExecStepIndex = index
	return s
}

// WithEmptyResponseCount returns a new state with the empty response count updated.
func (s GraphState) WithEmptyResponseCount(count int) GraphState {
	s.EmptyResponseCount = count
	return s
}

// WithPendingToolCalls returns a new state with pending tool calls updated.
func (s GraphState) WithPendingToolCalls(calls []tools.ToolCall) GraphState {
	s.PendingToolCalls = calls
	return s
}

// WithPlanStepCount returns a new state with the plan step count updated.
func (s GraphState) WithPlanStepCount(count int) GraphState {
	s.PlanStepCount = count
	return s
}

// ── UI State Helpers (pure functions returning updated state) ───────────────

// WithUIMessage returns a new state with the UIMessage appended.
// Pure function: does not mutate the receiver (deep copy of slices).
func (s GraphState) WithUIMessage(msg UIMessage) GraphState {
	newMessages := make([]UIMessage, len(s.UIMessages), len(s.UIMessages)+1)
	copy(newMessages, s.UIMessages)
	s.UIMessages = append(newMessages, msg)
	return s
}

// WithProviderMessage returns a new state with the provider message appended.
func (s GraphState) WithProviderMessage(msg provider.Message) GraphState {
	newMessages := make([]provider.Message, len(s.ProviderMessages), len(s.ProviderMessages)+1)
	copy(newMessages, s.ProviderMessages)
	s.ProviderMessages = append(newMessages, msg)
	return s
}

// WithUserMessage adds both UI and provider messages for a user turn.
func (s GraphState) WithUserMessage(content string) GraphState {
	// Deep copy UIMessages
	newUIMessages := make([]UIMessage, len(s.UIMessages), len(s.UIMessages)+1)
	copy(newUIMessages, s.UIMessages)
	s.UIMessages = append(newUIMessages, UIMessage{
		Kind:      KindUser,
		Content:   content,
		Timestamp: time.Now(),
	})

	// Deep copy ProviderMessages
	newProviderMessages := make([]provider.Message, len(s.ProviderMessages), len(s.ProviderMessages)+1)
	copy(newProviderMessages, s.ProviderMessages)
	s.ProviderMessages = append(newProviderMessages, provider.Message{
		Role:    provider.RoleUser,
		Content: content,
	})
	return s
}

// WithAssistantMessage adds both UI and provider messages for an assistant turn.
func (s GraphState) WithAssistantMessage(content string) GraphState {
	if content == "" {
		return s
	}

	// Deep copy UIMessages
	newUIMessages := make([]UIMessage, len(s.UIMessages), len(s.UIMessages)+1)
	copy(newUIMessages, s.UIMessages)
	s.UIMessages = append(newUIMessages, UIMessage{
		Kind:      KindAssistant,
		Content:   content,
		Timestamp: time.Now(),
	})

	// Deep copy ProviderMessages
	newProviderMessages := make([]provider.Message, len(s.ProviderMessages), len(s.ProviderMessages)+1)
	copy(newProviderMessages, s.ProviderMessages)
	s.ProviderMessages = append(newProviderMessages, provider.Message{
		Role:    provider.RoleAssistant,
		Content: content,
	})
	return s
}

// WithSystemNote adds a system note (UI only, not sent to LLM).
func (s GraphState) WithSystemNote(content string) GraphState {
	newMessages := make([]UIMessage, len(s.UIMessages), len(s.UIMessages)+1)
	copy(newMessages, s.UIMessages)
	s.UIMessages = append(newMessages, UIMessage{
		Kind:      KindSystem,
		Content:   content,
		Timestamp: time.Now(),
	})
	return s
}

// WithToolCall records a tool invocation in the UI.
func (s GraphState) WithToolCall(call tools.ToolCall) GraphState {
	newMessages := make([]UIMessage, len(s.UIMessages), len(s.UIMessages)+1)
	copy(newMessages, s.UIMessages)
	s.UIMessages = append(newMessages, UIMessage{
		Kind:      KindToolCall,
		Content:   formatToolCall(call),
		ToolName:  call.Name,
		Timestamp: time.Now(),
	})
	return s
}

// WithToolResult records a tool result in both UI and provider messages.
func (s GraphState) WithToolResult(result tools.ToolResult, toolName string) GraphState {
	// Deep copy UIMessages
	newUIMessages := make([]UIMessage, len(s.UIMessages), len(s.UIMessages)+1)
	copy(newUIMessages, s.UIMessages)
	s.UIMessages = append(newUIMessages, UIMessage{
		Kind:        KindToolResult,
		Content:     result.Output,
		ToolName:    toolName,
		ToolSuccess: result.Success,
		Timestamp:   time.Now(),
	})

	// Deep copy ProviderMessages
	status := "成功"
	if !result.Success {
		status = "失败"
	}
	newProviderMessages := make([]provider.Message, len(s.ProviderMessages), len(s.ProviderMessages)+1)
	copy(newProviderMessages, s.ProviderMessages)
	s.ProviderMessages = append(newProviderMessages, provider.Message{
		Role:    provider.RoleUser,
		Content: formatToolResult(toolName, status, result.Output),
	})
	return s
}

// WithPlanUIMessage adds a KindPlan UIMessage for the given plan.
func (s GraphState) WithPlanUIMessage(plan *ExecutionPlan) GraphState {
	newMessages := make([]UIMessage, len(s.UIMessages), len(s.UIMessages)+1)
	copy(newMessages, s.UIMessages)
	s.UIMessages = append(newMessages, UIMessage{
		Kind:      KindPlan,
		Content:   formatPlan(plan),
		Timestamp: time.Now(),
	})
	return s
}

// WithStepUpdate records a step execution event in UIMessages.
func (s GraphState) WithStepUpdate(step *PlanStep, status StepStatus, result, errMsg string) GraphState {
	content := formatStepUpdate(step, status, result, errMsg)
	newMessages := make([]UIMessage, len(s.UIMessages), len(s.UIMessages)+1)
	copy(newMessages, s.UIMessages)
	s.UIMessages = append(newMessages, UIMessage{
		Kind:        KindStepUpdate,
		Content:     content,
		ToolName:    step.ToolName,
		ToolSuccess: status == StepDone,
		Timestamp:   time.Now(),
	})
	return s
}

// StartStreaming initializes streaming state and adds an empty assistant message.
func (s GraphState) StartStreaming() GraphState {
	s.StreamingBuffer.Reset()
	newMessages := make([]UIMessage, len(s.UIMessages), len(s.UIMessages)+1)
	copy(newMessages, s.UIMessages)
	s.UIMessages = append(newMessages, UIMessage{
		Kind:      KindAssistant,
		Content:   "",
		Timestamp: time.Now(),
	})
	s.StreamingMessageIndex = len(s.UIMessages) - 1
	return s
}

// AppendStreamingChunk appends a chunk to the streaming buffer and updates the message.
// Note: This method mutates UIMessages[StreamingMessageIndex] in place for performance.
// This is acceptable because streaming is a transient state within a single node execution.
func (s GraphState) AppendStreamingChunk(chunk string) GraphState {
	if s.StreamingMessageIndex < 0 || s.StreamingMessageIndex >= len(s.UIMessages) {
		return s
	}
	s.StreamingBuffer.WriteString(chunk)
	s.UIMessages[s.StreamingMessageIndex].Content = s.StreamingBuffer.String()
	return s
}

// FinishStreaming completes the streaming message and adds it to provider messages.
func (s GraphState) FinishStreaming() GraphState {
	if s.StreamingMessageIndex < 0 {
		return s
	}
	content := s.StreamingBuffer.String()
	if content != "" {
		newProviderMessages := make([]provider.Message, len(s.ProviderMessages), len(s.ProviderMessages)+1)
		copy(newProviderMessages, s.ProviderMessages)
		s.ProviderMessages = append(newProviderMessages, provider.Message{
			Role:    provider.RoleAssistant,
			Content: content,
		})
	}
	s.StreamingMessageIndex = -1
	s.StreamingBuffer.Reset()
	return s
}

// ── Formatting helpers ───────────────────────────────────────────────────────

func formatToolCall(call tools.ToolCall) string {
	return "调用工具: " + call.Name + " " + formatParams(call.Params)
}

func formatParams(params map[string]any) string {
	// Simple formatting; could use json.Marshal for complex params
	var parts []string
	for k, v := range params {
		parts = append(parts, k+"="+toString(v))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func formatToolResult(toolName, status, output string) string {
	return "[工具结果 " + toolName + " - " + status + "]\n" + output
}

func formatPlan(plan *ExecutionPlan) string {
	var sb strings.Builder
	sb.WriteString(plan.Summary)
	for _, step := range plan.Steps {
		sb.WriteString("\n  ")
		sb.WriteString(step.ID)
		sb.WriteString(". ")
		sb.WriteString(step.Description)
	}
	return sb.String()
}

func formatStepUpdate(step *PlanStep, status StepStatus, result, errMsg string) string {
	content := "[" + status.String() + "] " + step.Description
	if result != "" {
		content += "\n" + result
	}
	if errMsg != "" {
		content += "\n错误: " + errMsg
	}
	return content
}
