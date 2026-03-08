package agent

// Integration tests for the full graph flow:
//   startPlan → nodePlan → nodeWaitHuman (graph suspends)
//   → resume → nodeHandleConfirmation → nodeExecuteStep → nodeDone
//
// These tests exercise the graph as a whole, verifying that state threads
// correctly through every node and that the interrupt/resume mechanism works.

import (
	"context"
	"testing"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
	"github.com/dswcpp/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

// ── helpers ──────────────────────────────────────────────────────────────────

// scriptedProvider returns pre-scripted responses in order.
// Each Send() or Complete() call consumes the next response.
type scriptedProvider struct {
	responses []string
	idx       int
}

func (s *scriptedProvider) Complete(_ context.Context, _ []provider.Message) (provider.Result, error) {
	if s.idx >= len(s.responses) {
		return provider.Result{Content: ""}, nil
	}
	resp := s.responses[s.idx]
	s.idx++
	return provider.Result{Content: resp}, nil
}

func (s *scriptedProvider) CompleteStream(_ context.Context, _ []provider.Message, onChunk func(string)) error {
	if s.idx >= len(s.responses) {
		return nil
	}
	resp := s.responses[s.idx]
	s.idx++
	onChunk(resp)
	return nil
}

func (s *scriptedProvider) ModelID() string { return "scripted" }

// planResponse returns a valid ```plan block as an LLM response string.
func planResponse(toolName string) string {
	return "I have analysed the repo.\n\n" +
		"```plan\n" +
		`{"summary":"Test plan","steps":[{"id":"1","description":"Run tool","tool":"` +
		toolName +
		`","params":{},"critical":false}]}` + "\n" +
		"```\n\n" +
		"Please confirm to execute."
}

// criticalPlanResponse returns a plan whose single step is critical.
func criticalPlanResponse(toolName string) string {
	return "```plan\n" +
		`{"summary":"Critical plan","steps":[{"id":"1","description":"Critical step","tool":"` +
		toolName +
		`","params":{},"critical":true}]}` + "\n" +
		"```"
}

func newTestAgent(p provider.Provider, registry *tools.Registry, tr *aii18n.Translator) *TwoPhaseAgent {
	return &TwoPhaseAgent{
		provider:     p,
		fullRegistry: registry,
		readRegistry: tools.NewRegistry(), // empty read-only registry for planning
		session:      NewSession("test"),
		tr:           tr,
		maxPlanSteps: 10,
		stepTimeout:  defaultStepTimeout,
		state: GraphState{
			Phase:           PhasePlanning,
			ToolCallHistory: make(map[string]int),
		},
	}
}

func testTranslator() *aii18n.Translator {
	return aii18n.NewTranslator(&i18n.TranslationSet{
		AIAgentStepTimeout:                  "⏱️ step timeout: %s\n\nstep: %s",
		AIAgentCriticalStepFailed:           "critical step failed: %s\nreason: %s",
		AIAgentPossibleReasons:              "\n\n💡 possible reasons:",
		AITwoPhaseAgentMaxStepsExceeded:     "max steps exceeded (%d)",
		AITwoPhaseAgentExecutionCancelled:   "execution cancelled",
		AITwoPhaseAgentUserFeedbackPrompt:   "user feedback: %s",
		AITwoPhaseAgentEmptyResponseError:   "empty response error",
		AITwoPhaseAgentContinueAnalysis:     "continue analysis",
		AITwoPhaseAgentSystemPromptIntro:    "system prompt",
		AITwoPhaseAgentExecuting:            "already executing",
		AITwoPhaseAgentToolCallWarning:      "tool %s called %d times",
		AITwoPhaseAgentSystemPrefix:         "[system] ",
		AITwoPhaseAgentToolResultPrefix:     "[%s result] %s",
		AITwoPhaseAgentPlanErrorsIntro:      "plan errors:\n",
		AITwoPhaseAgentPlanRegeneratePrompt: "please fix the plan",
		AITwoPhaseAgentPlanValidationFailed: "plan validation failed",
		AIAgentToolNotAllowedInPlanning:     "tool not allowed in planning: %s",
	})
}

// ── tests ─────────────────────────────────────────────────────────────────────

// TestGraph_PlanSuspendsAtWaitHuman verifies that after the LLM returns a valid
// plan, the graph suspends (nodeWaitHuman), sets Phase=PhaseWaitingConfirm,
// and records ResumeFrom=NodeHandleConfirmation.
func TestGraph_PlanSuspendsAtWaitHuman(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&mockTool{name: "stage_all", params: map[string]tools.ParamSchema{}, perm: tools.PermWriteLocal})

	p := &scriptedProvider{responses: []string{planResponse("stage_all")}}
	a := newTestAgent(p, registry, testTranslator())

	// seed initial plan messages (normally done by startPlan)
	a.state.PlanMessages = []provider.Message{
		{Role: provider.RoleSystem, Content: "sys"},
		{Role: provider.RoleUser, Content: "user request"},
	}

	err := a.planLoop(context.Background(), nil)
	assert.NoError(t, err)

	assert.Equal(t, PhaseWaitingConfirm, a.state.Phase, "should suspend in WaitingConfirm")
	assert.Equal(t, NodeHandleConfirmation, a.state.ResumeFrom, "should record resume checkpoint")
	assert.NotNil(t, a.state.Plan, "should have produced a plan")
	assert.Equal(t, "Test plan", a.state.Plan.Summary)
}

// TestGraph_ConfirmResumesExecution verifies the Y path:
// resume("y") → nodeHandleConfirmation → nodeExecuteStep × N → nodeDone.
func TestGraph_ConfirmResumesExecution(t *testing.T) {
	registry := tools.NewRegistry()
	executed := false
	registry.Register(&mockToolFn{
		name: "stage_all",
		perm: tools.PermWriteLocal,
		fn: func(ctx context.Context, call tools.ToolCall) tools.ToolResult {
			executed = true
			return tools.ToolResult{CallID: call.ID, Success: true, Output: "staged"}
		},
	})

	a := newTestAgent(&scriptedProvider{}, registry, testTranslator())
	// Pre-set state as if planLoop just suspended.
	a.state.Phase = PhaseWaitingConfirm
	a.state.ResumeFrom = NodeHandleConfirmation
	a.state.Plan = &ExecutionPlan{
		Summary: "staged plan",
		Steps: []*PlanStep{
			{ID: "1", Description: "stage", ToolName: "stage_all", Params: map[string]any{}, Critical: false, Status: StepPending},
		},
	}
	a.session.AddPlanUIMessage(a.state.Plan)

	err := a.resume(context.Background(), "y", nil)
	assert.NoError(t, err)

	assert.True(t, executed, "tool should have been executed")
	assert.Equal(t, PhaseDone, a.state.Phase)
	assert.Empty(t, a.state.ResumeFrom, "checkpoint should be cleared after done")
	assert.Equal(t, StepDone, a.state.Plan.Steps[0].Status)
}

// TestGraph_DenyEndsWithCancelled verifies the N path:
// resume("n") → nodeHandleConfirmation → NodeEnd with PhaseCancelled.
func TestGraph_DenyEndsWithCancelled(t *testing.T) {
	registry := tools.NewRegistry()
	executed := false
	registry.Register(&mockToolFn{
		name: "stage_all",
		perm: tools.PermWriteLocal,
		fn: func(_ context.Context, call tools.ToolCall) tools.ToolResult {
			executed = true
			return tools.ToolResult{CallID: call.ID, Success: true, Output: "staged"}
		},
	})

	a := newTestAgent(&scriptedProvider{}, registry, testTranslator())
	a.state.Phase = PhaseWaitingConfirm
	a.state.ResumeFrom = NodeHandleConfirmation
	a.state.Plan = &ExecutionPlan{
		Summary: "plan to cancel",
		Steps:   []*PlanStep{{ID: "1", ToolName: "stage_all", Params: map[string]any{}, Status: StepPending}},
	}
	a.session.AddPlanUIMessage(a.state.Plan)

	err := a.resume(context.Background(), "n", nil)
	assert.NoError(t, err)

	assert.False(t, executed, "tool must NOT execute after denial")
	assert.Equal(t, PhaseCancelled, a.state.Phase)
	assert.Empty(t, a.state.ResumeFrom)
}

// TestGraph_FeedbackTriggersReplan verifies the replan path:
// resume("some feedback") → nodeHandleConfirmation → NodePlan → (LLM returns new plan)
// → nodeWaitHuman (second suspend).
func TestGraph_FeedbackTriggersReplan(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&mockTool{name: "stage_all", params: map[string]tools.ParamSchema{}, perm: tools.PermWriteLocal})

	// Second call returns the revised plan.
	p := &scriptedProvider{responses: []string{planResponse("stage_all")}}
	a := newTestAgent(p, registry, testTranslator())
	a.state.Phase = PhaseWaitingConfirm
	a.state.ResumeFrom = NodeHandleConfirmation
	a.state.Plan = &ExecutionPlan{
		Summary: "original plan",
		Steps:   []*PlanStep{{ID: "1", ToolName: "stage_all", Params: map[string]any{}, Status: StepPending}},
	}
	// Seed plan messages so the LLM call inside nodePlan has something to work with.
	a.state.PlanMessages = []provider.Message{
		{Role: provider.RoleSystem, Content: "sys"},
		{Role: provider.RoleUser, Content: "stage everything"},
	}
	a.session.AddPlanUIMessage(a.state.Plan)

	err := a.resume(context.Background(), "please add a commit step", nil)
	assert.NoError(t, err)

	// After feedback → replan, the graph should suspend again at WaitHuman.
	assert.Equal(t, PhaseWaitingConfirm, a.state.Phase, "should suspend again after replan")
	assert.Equal(t, NodeHandleConfirmation, a.state.ResumeFrom)
	assert.NotNil(t, a.state.Plan)
}

// TestGraph_NilPlanInExecutionReturnsError verifies nodeExecuteStep returns an
// explicit error when the plan is nil (defensive nil check).
func TestGraph_NilPlanInExecutionReturnsError(t *testing.T) {
	a := newTestAgent(&scriptedProvider{}, tools.NewRegistry(), testTranslator())
	// Deliberately leave Plan nil while entering execution.
	a.state.Phase = PhaseExecuting
	a.state.Plan = nil

	err := a.execute(context.Background(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil plan")
}

// TestGraph_CheckpointerSaveOnSuspend verifies that nodeWaitHuman saves the
// checkpoint with the updated state (Phase=WaitingConfirm, ResumeFrom set).
func TestGraph_CheckpointerSaveOnSuspend(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&mockTool{name: "stage_all", params: map[string]tools.ParamSchema{}, perm: tools.PermWriteLocal})

	c := NewMemoryCheckpointer()
	p := &scriptedProvider{responses: []string{planResponse("stage_all")}}
	a := newTestAgent(p, registry, testTranslator())
	a.SetCheckpointer(c, "thread-1")
	a.state.PlanMessages = []provider.Message{
		{Role: provider.RoleSystem, Content: "sys"},
		{Role: provider.RoleUser, Content: "do stuff"},
	}

	_ = a.planLoop(context.Background(), nil)

	saved, ok := c.Load("thread-1")
	assert.True(t, ok, "checkpoint should have been saved on suspend")
	assert.Equal(t, PhaseWaitingConfirm, saved.Phase)
	assert.Equal(t, NodeHandleConfirmation, saved.ResumeFrom)
}

// TestGraph_CheckpointerClearedOnDone verifies that the checkpoint is removed
// after successful execution (nodeDone).
func TestGraph_CheckpointerClearedOnDone(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&mockTool{name: "stage_all", params: map[string]tools.ParamSchema{}, perm: tools.PermWriteLocal})

	c := NewMemoryCheckpointer()
	a := newTestAgent(&scriptedProvider{}, registry, testTranslator())
	a.SetCheckpointer(c, "thread-2")
	// Pre-save a checkpoint as if planLoop had suspended.
	_ = c.Save("thread-2", GraphState{Phase: PhaseWaitingConfirm, ResumeFrom: NodeHandleConfirmation})

	a.state.Phase = PhaseWaitingConfirm
	a.state.ResumeFrom = NodeHandleConfirmation
	a.state.Plan = &ExecutionPlan{
		Summary: "plan",
		Steps:   []*PlanStep{{ID: "1", ToolName: "stage_all", Params: map[string]any{}, Status: StepPending}},
	}
	a.state.ToolCallHistory = make(map[string]int)
	a.session.AddPlanUIMessage(a.state.Plan)

	_ = a.resume(context.Background(), "y", nil)

	_, ok := c.Load("thread-2")
	assert.False(t, ok, "checkpoint should be cleared after done")
}

// TestGraph_CheckpointerClearedOnCancel verifies that the checkpoint is removed
// after the user denies execution.
func TestGraph_CheckpointerClearedOnCancel(t *testing.T) {
	registry := tools.NewRegistry()
	c := NewMemoryCheckpointer()
	a := newTestAgent(&scriptedProvider{}, registry, testTranslator())
	a.SetCheckpointer(c, "thread-3")
	_ = c.Save("thread-3", GraphState{Phase: PhaseWaitingConfirm, ResumeFrom: NodeHandleConfirmation})

	a.state.Phase = PhaseWaitingConfirm
	a.state.ResumeFrom = NodeHandleConfirmation
	a.state.Plan = &ExecutionPlan{Steps: []*PlanStep{}}
	a.state.ToolCallHistory = make(map[string]int)
	a.session.AddPlanUIMessage(a.state.Plan)

	_ = a.resume(context.Background(), "no", nil)

	_, ok := c.Load("thread-3")
	assert.False(t, ok, "checkpoint should be cleared after cancel")
}

// TestGraph_StateThroughNodes verifies that state properly threads through all
// nodes: modifications in one node are visible in the next.
func TestGraph_StateThroughNodes(t *testing.T) {
	registry := tools.NewRegistry()
	callOrder := []string{}

	registry.Register(&mockToolFn{
		name: "stage_all",
		perm: tools.PermWriteLocal,
		fn: func(_ context.Context, call tools.ToolCall) tools.ToolResult {
			callOrder = append(callOrder, "stage_all")
			return tools.ToolResult{CallID: call.ID, Success: true, Output: "ok"}
		},
	})
	registry.Register(&mockToolFn{
		name: "commit",
		perm: tools.PermWriteLocal,
		fn: func(_ context.Context, call tools.ToolCall) tools.ToolResult {
			callOrder = append(callOrder, "commit")
			return tools.ToolResult{CallID: call.ID, Success: true, Output: "committed"}
		},
	})

	a := newTestAgent(&scriptedProvider{}, registry, testTranslator())
	a.state.Phase = PhaseWaitingConfirm
	a.state.ResumeFrom = NodeHandleConfirmation
	a.state.Plan = &ExecutionPlan{
		Summary: "multi-step",
		Steps: []*PlanStep{
			{ID: "1", ToolName: "stage_all", Params: map[string]any{}, Status: StepPending},
			{ID: "2", ToolName: "commit", Params: map[string]any{"message": "msg"}, Status: StepPending},
		},
	}
	a.state.ToolCallHistory = make(map[string]int)
	a.session.AddPlanUIMessage(a.state.Plan)

	err := a.resume(context.Background(), "yes", nil)
	assert.NoError(t, err)

	assert.Equal(t, []string{"stage_all", "commit"}, callOrder, "steps should execute in order")
	assert.Equal(t, StepDone, a.state.Plan.Steps[0].Status)
	assert.Equal(t, StepDone, a.state.Plan.Steps[1].Status)
	assert.Equal(t, PhaseDone, a.state.Phase)
}

// ── mock helpers ──────────────────────────────────────────────────────────────

// mockToolFn is a tool whose Execute function is provided at construction.
type mockToolFn struct {
	name string
	perm tools.PermissionLevel
	fn   func(ctx context.Context, call tools.ToolCall) tools.ToolResult
}

func (m *mockToolFn) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        m.name,
		Description: "mock tool fn",
		Params:      map[string]tools.ParamSchema{},
		Permission:  m.perm,
	}
}

func (m *mockToolFn) Execute(ctx context.Context, call tools.ToolCall) tools.ToolResult {
	return m.fn(ctx, call)
}
