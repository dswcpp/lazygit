package agent

import (
	"context"
	"strings"
	"testing"
	"time"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
	"github.com/dswcpp/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

// mockTranslator 创建一个用于测试的 Translator
func mockTranslator() *aii18n.Translator {
	return aii18n.NewTranslator(&i18n.TranslationSet{
		AIAgentStepTimeout:          "⏱️ 步骤超时: %s\n\n步骤: %s",
		AIAgentCriticalStepFailed:   "关键步骤失败: %s\n原因: %s",
		AIAgentPossibleReasons:      "\n\n💡 可能的原因：",
		AITwoPhaseAgentMaxStepsExceeded: "已达到最大规划步骤数 (%d)，请简化任务或分步执行",
	})
}

// slowTool 模拟一个执行缓慢的工具
type slowTool struct {
	name     string
	duration time.Duration
}

func (s *slowTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        s.name,
		Description: "slow tool for testing",
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermWriteLocal,
	}
}

func (s *slowTool) Execute(ctx context.Context, call tools.ToolCall) tools.ToolResult {
	select {
	case <-time.After(s.duration):
		return tools.ToolResult{CallID: call.ID, Success: true, Output: "completed"}
	case <-ctx.Done():
		return tools.ToolResult{CallID: call.ID, Success: false, Output: "cancelled"}
	}
}

// errorTool 模拟一个总是失败的工具
type errorTool struct {
	name      string
	errorMsg  string
	perm      tools.PermissionLevel
}

func (e *errorTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        e.name,
		Description: "error tool for testing",
		Params:      map[string]tools.ParamSchema{},
		Permission:  e.perm,
	}
}

func (e *errorTool) Execute(ctx context.Context, call tools.ToolCall) tools.ToolResult {
	return tools.ToolResult{
		CallID:  call.ID,
		Success: false,
		Output:  e.errorMsg,
	}
}

func TestExecuteWithTimeout(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&slowTool{
		name:     "slow_operation",
		duration: 2 * time.Second, // 执行需要 2 秒
	})

	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		provider:     &mockProvider{},
		session:      NewSession("test"),
		stepTimeout:  500 * time.Millisecond, // 超时设置为 500ms
		tr:           mockTranslator(),
		state:        GraphState{ToolCallHistory: make(map[string]int)},
	}

	plan := &ExecutionPlan{
		Summary: "Test timeout",
		Steps: []*PlanStep{
			{
				ID:          "1",
				Description: "Slow operation",
				ToolName:    "slow_operation",
				Params:      map[string]any{},
				Critical:    false,
				Status:      StepPending,
			},
		},
	}

	ctx := context.Background()
	agent.state.Plan = plan
	agent.session.AddPlanUIMessage(plan)
	err := agent.execute(ctx, nil)

	// 应该成功完成（非关键步骤超时不会导致整体失败）
	assert.NoError(t, err)

	// 检查步骤状态（从 state 中获取）
	assert.Equal(t, StepFailed, agent.state.Plan.Steps[0].Status)
	assert.Contains(t, agent.state.Plan.Steps[0].Error, "超时")
}

func TestExecuteWithTimeoutCriticalStep(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&slowTool{
		name:     "slow_operation",
		duration: 2 * time.Second,
	})

	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		provider:     &mockProvider{},
		session:      NewSession("test"),
		stepTimeout:  500 * time.Millisecond,
		tr:           mockTranslator(),
		state:        GraphState{ToolCallHistory: make(map[string]int)},
	}

	plan := &ExecutionPlan{
		Summary: "Test timeout with critical step",
		Steps: []*PlanStep{
			{
				ID:          "1",
				Description: "Critical slow operation",
				ToolName:    "slow_operation",
				Params:      map[string]any{},
				Critical:    true, // 关键步骤
				Status:      StepPending,
			},
		},
	}

	ctx := context.Background()
	agent.state.Plan = plan
	agent.session.AddPlanUIMessage(plan)
	err := agent.execute(ctx, nil)

	// 关键步骤超时应该返回错误
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "关键步骤超时")
}

func TestFormatToolNotFoundError(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&mockTool{
		name: "stage_all",
		params: map[string]tools.ParamSchema{},
		perm: tools.PermWriteLocal,
	})
	registry.Register(&mockTool{
		name: "stage_file",
		params: map[string]tools.ParamSchema{},
		perm: tools.PermWriteLocal,
	})

	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		provider:     &mockProvider{},
		tr:           mockTranslator(),
	}

	// 测试相似工具建议
	errMsg := agent.formatToolNotFoundError("stag", "stag")
	assert.Contains(t, errMsg, "未知工具")
	assert.Contains(t, errMsg, "stage_all") // 应该建议相似的工具
}

func TestFormatTimeoutError(t *testing.T) {
	agent := &TwoPhaseAgent{
		tr: mockTranslator(),
	}

	step := &PlanStep{
		ID:          "1",
		Description: "Test operation",
		ToolName:    "test_tool",
	}

	errMsg := agent.formatTimeoutError(step, 30*time.Second)
	assert.Contains(t, errMsg, "超时")
	assert.Contains(t, errMsg, "30s")
	assert.Contains(t, errMsg, "可能的原因")
	assert.Contains(t, errMsg, "建议")
}

func TestFormatExecutionError(t *testing.T) {
	agent := &TwoPhaseAgent{
		tr: mockTranslator(),
	}

	tests := []struct {
		name        string
		toolName    string
		rawError    string
		expectInMsg []string
	}{
		{
			name:        "commit without staged files",
			toolName:    "commit",
			rawError:    "nothing to commit",
			expectInMsg: []string{"执行失败", "暂存文件"},
		},
		{
			name:        "checkout with uncommitted changes",
			toolName:    "checkout",
			rawError:    "uncommitted changes would be overwritten",
			expectInMsg: []string{"执行失败", "提交或暂存"},
		},
		{
			name:        "push rejected",
			toolName:    "push",
			rawError:    "rejected: non-fast-forward",
			expectInMsg: []string{"执行失败", "拉取远程更新"},
		},
		{
			name:        "permission denied",
			toolName:    "any_tool",
			rawError:    "permission denied",
			expectInMsg: []string{"执行失败", "权限"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := &PlanStep{
				ID:          "1",
				Description: "Test step",
				ToolName:    tt.toolName,
			}

			errMsg := agent.formatExecutionError(step, tt.rawError)

			for _, expected := range tt.expectInMsg {
				assert.Contains(t, errMsg, expected,
					"Error message should contain: %s", expected)
			}
		})
	}
}

func TestExecuteWithNonCriticalFailure(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&errorTool{
		name:     "failing_tool",
		errorMsg: "operation failed",
		perm:     tools.PermWriteLocal,
	})
	registry.Register(&mockTool{
		name: "success_tool",
		params: map[string]tools.ParamSchema{},
		perm: tools.PermWriteLocal,
	})

	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		provider:     &mockProvider{},
		session:      NewSession("test"),
		stepTimeout:  30 * time.Second,
		tr:           mockTranslator(),
		state:        GraphState{ToolCallHistory: make(map[string]int)},
	}

	plan := &ExecutionPlan{
		Summary: "Test non-critical failure",
		Steps: []*PlanStep{
			{
				ID:          "1",
				Description: "Non-critical failing step",
				ToolName:    "failing_tool",
				Params:      map[string]any{},
				Critical:    false, // 非关键步骤
				Status:      StepPending,
			},
			{
				ID:          "2",
				Description: "Success step",
				ToolName:    "success_tool",
				Params:      map[string]any{},
				Critical:    false,
				Status:      StepPending,
			},
		},
	}

	ctx := context.Background()
	agent.state.Plan = plan
	agent.session.AddPlanUIMessage(plan)
	err := agent.execute(ctx, nil)

	// 非关键步骤失败不应该导致整体失败
	assert.NoError(t, err)

	// 第一步应该失败（从 state 中获取）
	assert.Equal(t, StepFailed, agent.state.Plan.Steps[0].Status)
	assert.Contains(t, agent.state.Plan.Steps[0].Error, "operation failed")

	// 第二步应该成功
	assert.Equal(t, StepDone, agent.state.Plan.Steps[1].Status)
}

func TestExecuteWithCriticalFailure(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&errorTool{
		name:     "failing_tool",
		errorMsg: "critical operation failed",
		perm:     tools.PermWriteLocal,
	})
	registry.Register(&mockTool{
		name: "success_tool",
		params: map[string]tools.ParamSchema{},
		perm: tools.PermWriteLocal,
	})

	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		provider:     &mockProvider{},
		session:      NewSession("test"),
		stepTimeout:  30 * time.Second,
		tr:           mockTranslator(),
		state:        GraphState{ToolCallHistory: make(map[string]int)},
	}

	plan := &ExecutionPlan{
		Summary: "Test critical failure",
		Steps: []*PlanStep{
			{
				ID:          "1",
				Description: "Critical failing step",
				ToolName:    "failing_tool",
				Params:      map[string]any{},
				Critical:    true, // 关键步骤
				Status:      StepPending,
			},
			{
				ID:          "2",
				Description: "Success step (should not execute)",
				ToolName:    "success_tool",
				Params:      map[string]any{},
				Critical:    false,
				Status:      StepPending,
			},
		},
	}

	ctx := context.Background()
	agent.state.Plan = plan
	agent.session.AddPlanUIMessage(plan)
	err := agent.execute(ctx, nil)

	// 关键步骤失败应该返回错误
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "关键步骤失败")

	// 第一步应该失败（从 state 中获取）
	assert.Equal(t, StepFailed, agent.state.Plan.Steps[0].Status)

	// 第二步不应该执行（仍然是 Pending）
	assert.Equal(t, StepPending, agent.state.Plan.Steps[1].Status)
}

func TestFindSimilarTools(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&mockTool{name: "stage_all", params: map[string]tools.ParamSchema{}, perm: tools.PermWriteLocal})
	registry.Register(&mockTool{name: "stage_file", params: map[string]tools.ParamSchema{}, perm: tools.PermWriteLocal})
	registry.Register(&mockTool{name: "unstage_all", params: map[string]tools.ParamSchema{}, perm: tools.PermWriteLocal})
	registry.Register(&mockTool{name: "commit", params: map[string]tools.ParamSchema{}, perm: tools.PermWriteLocal})

	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		tr:           mockTranslator(),
	}

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "prefix match",
			input:    "stage",
			expected: []string{"stage_all", "stage_file"},
		},
		{
			name:     "exact match",
			input:    "commit",
			expected: []string{"commit"},
		},
		{
			name:     "partial match",
			input:    "all",
			expected: []string{"stage_all", "unstage_all"},
		},
		{
			name:     "no match",
			input:    "xyz",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			similar := agent.findSimilarTools(tt.input)

			// 检查是否包含所有期望的工具
			for _, exp := range tt.expected {
				assert.Contains(t, similar, exp)
			}

			// 如果没有期望的工具，确保返回为空
			if len(tt.expected) == 0 {
				assert.Empty(t, similar)
			}
		})
	}
}

func TestGetRecoverySuggestions(t *testing.T) {
	agent := &TwoPhaseAgent{
		tr: mockTranslator(),
	}

	tests := []struct {
		name        string
		toolName    string
		errorMsg    string
		expectInSuggestions []string
	}{
		{
			name:        "commit nothing to commit",
			toolName:    "commit",
			errorMsg:    "nothing to commit, working tree clean",
			expectInSuggestions: []string{"暂存文件"},
		},
		{
			name:        "checkout uncommitted changes",
			toolName:    "checkout",
			errorMsg:    "error: Your local changes would be overwritten by checkout",
			expectInSuggestions: []string{"提交", "暂存"},
		},
		{
			name:        "push rejected",
			toolName:    "push",
			errorMsg:    "error: failed to push some refs to remote",
			expectInSuggestions: []string{"拉取", "更新"},
		},
		{
			name:        "merge conflict",
			toolName:    "merge",
			errorMsg:    "CONFLICT (content): Merge conflict in file.txt",
			expectInSuggestions: []string{"解决冲突"},
		},
		{
			name:        "permission denied",
			toolName:    "any_tool",
			errorMsg:    "permission denied",
			expectInSuggestions: []string{"权限"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := agent.getRecoverySuggestions(tt.toolName, tt.errorMsg)

			assert.NotEmpty(t, suggestions, "Should provide at least one suggestion")

			// 打印实际的建议，用于调试
			t.Logf("Suggestions for %s: %v", tt.name, suggestions)

			// 检查是否包含期望的建议
			found := false
			for _, suggestion := range suggestions {
				for _, expected := range tt.expectInSuggestions {
					if strings.Contains(suggestion, expected) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			assert.True(t, found, "Should contain expected suggestion keywords. Got: %v, Expected one of: %v", suggestions, tt.expectInSuggestions)
		})
	}
}
