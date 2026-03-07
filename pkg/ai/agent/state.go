package agent

import (
	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

// GraphState holds all control-flow state needed to describe, checkpoint, and
// resume a TwoPhaseAgent execution. It is intentionally separate from Session
// so that UI-display concerns (UIMessages, streaming) stay isolated from the
// agent's execution logic.
//
// LangGraph analogy: this is the typed State object that flows through every
// node of the graph and is the single source of truth for the agent's progress.
type GraphState struct {
	// Phase is the current execution phase.
	Phase AgentPhase

	// Plan is the execution plan produced at the end of the planning phase.
	// Nil until the LLM has committed to a concrete plan.
	Plan *ExecutionPlan

	// PlanMessages is the LLM conversation history for the planning phase.
	// It carries a dedicated system prompt and may be pruned independently
	// of Session.providerMessages (which serves the basic ReAct Agent).
	PlanMessages []provider.Message

	// ToolCallHistory tracks the number of times each unique (tool, params)
	// combination has been invoked during planning, used to break infinite loops.
	ToolCallHistory map[string]int

	// ── Planning-loop cursors (reset at the start of each plan run) ─────────

	// PlanStepCount is the number of LLM calls made so far in the current plan
	// loop. Compared against TwoPhaseAgent.maxPlanSteps to prevent runaway loops.
	PlanStepCount int

	// EmptyResponseCount tracks consecutive LLM responses that contained neither
	// a plan block nor tool calls. Three in a row aborts the planning loop.
	EmptyResponseCount int

	// PendingToolCalls holds tool calls produced by nodePlan and consumed by
	// nodeCallTools before control returns to nodePlan.
	PendingToolCalls []tools.ToolCall

	// ── Execution cursor ────────────────────────────────────────────────────

	// ExecStepIndex is the index of the next PlanStep to execute.
	// nodeExecuteStep increments this after each successful step dispatch.
	ExecStepIndex int

	// ── Interrupt / Resume ──────────────────────────────────────────────────

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
}

// Reset clears all planning state for a fresh run.
func (s *GraphState) Reset() {
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
}

// AppendPlanMessage appends msg to PlanMessages, trimming history to max entries
// while always preserving the first entry (the system prompt).
func (s *GraphState) AppendPlanMessage(msg provider.Message, max int) {
	s.PlanMessages = append(s.PlanMessages, msg)
	if len(s.PlanMessages) > max {
		s.PlanMessages = append(
			s.PlanMessages[:1],
			s.PlanMessages[len(s.PlanMessages)-max+1:]...,
		)
	}
}
