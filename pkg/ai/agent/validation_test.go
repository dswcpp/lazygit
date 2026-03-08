package agent

import (
	"context"
	"testing"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
	"github.com/stretchr/testify/assert"
)

// mockTool 用于测试
type mockTool struct {
	name   string
	params map[string]tools.ParamSchema
	perm   tools.PermissionLevel
}

func (m *mockTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        m.name,
		Description: "mock tool",
		Params:      m.params,
		Permission:  m.perm,
	}
}

func (m *mockTool) Execute(ctx context.Context, call tools.ToolCall) tools.ToolResult {
	return tools.ToolResult{CallID: call.ID, Success: true, Output: "ok"}
}

// mockProvider 用于测试
type mockProvider struct{}

func (m *mockProvider) Complete(ctx context.Context, messages []provider.Message) (provider.Result, error) {
	return provider.Result{Content: "mock response"}, nil
}

func (m *mockProvider) CompleteStream(ctx context.Context, messages []provider.Message, onChunk func(string)) error {
	return nil
}

func (m *mockProvider) ModelID() string {
	return "mock"
}

func TestValidatePlan_ValidPlan(t *testing.T) {
	// 创建注册表
	registry := tools.NewRegistry()
	registry.Register(&mockTool{
		name: "test_tool",
		params: map[string]tools.ParamSchema{
			"message": {Type: "string", Required: true},
		},
		perm: tools.PermWriteLocal,
	})

	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		provider:     &mockProvider{},
	}

	plan := &ExecutionPlan{
		Summary: "Test plan",
		Steps: []*PlanStep{
			{
				ID:          "1",
				Description: "Test step",
				ToolName:    "test_tool",
				Params:      map[string]any{"message": "hello"},
			},
		},
	}

	errors := agent.validatePlan(plan)
	assert.Empty(t, errors, "Valid plan should have no errors")
}

func TestValidatePlan_UnknownTool(t *testing.T) {
	registry := tools.NewRegistry()
	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		provider:     &mockProvider{},
	}

	plan := &ExecutionPlan{
		Summary: "Test plan",
		Steps: []*PlanStep{
			{
				ID:          "1",
				Description: "Test step",
				ToolName:    "unknown_tool",
				Params:      map[string]any{},
			},
		},
	}

	errors := agent.validatePlan(plan)
	assert.Len(t, errors, 1, "Should have one error")
	assert.Contains(t, errors[0], "未知工具")
	assert.Contains(t, errors[0], "unknown_tool")
}

func TestValidatePlan_MissingRequiredParam(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&mockTool{
		name: "test_tool",
		params: map[string]tools.ParamSchema{
			"message": {Type: "string", Required: true},
		},
		perm: tools.PermWriteLocal,
	})

	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		provider:     &mockProvider{},
	}

	plan := &ExecutionPlan{
		Summary: "Test plan",
		Steps: []*PlanStep{
			{
				ID:          "1",
				Description: "Test step",
				ToolName:    "test_tool",
				Params:      map[string]any{}, // 缺少 message 参数
			},
		},
	}

	errors := agent.validatePlan(plan)
	assert.Len(t, errors, 1, "Should have one error")
	assert.Contains(t, errors[0], "缺少必需参数")
	assert.Contains(t, errors[0], "message")
}

func TestValidatePlan_WrongParamType(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&mockTool{
		name: "test_tool",
		params: map[string]tools.ParamSchema{
			"count": {Type: "int", Required: true},
		},
		perm: tools.PermWriteLocal,
	})

	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		provider:     &mockProvider{},
	}

	plan := &ExecutionPlan{
		Summary: "Test plan",
		Steps: []*PlanStep{
			{
				ID:          "1",
				Description: "Test step",
				ToolName:    "test_tool",
				Params:      map[string]any{"count": "not a number"}, // 类型错误
			},
		},
	}

	errors := agent.validatePlan(plan)
	assert.Len(t, errors, 1, "Should have one error")
	assert.Contains(t, errors[0], "类型错误")
}

func TestValidatePlan_EmptyStringParam(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&mockTool{
		name: "test_tool",
		params: map[string]tools.ParamSchema{
			"message": {Type: "string", Required: true},
		},
		perm: tools.PermWriteLocal,
	})

	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		provider:     &mockProvider{},
	}

	plan := &ExecutionPlan{
		Summary: "Test plan",
		Steps: []*PlanStep{
			{
				ID:          "1",
				Description: "Test step",
				ToolName:    "test_tool",
				Params:      map[string]any{"message": "   "}, // 空字符串
			},
		},
	}

	errors := agent.validatePlan(plan)
	assert.Len(t, errors, 1, "Should have one error")
	assert.Contains(t, errors[0], "不能为空")
}

func TestValidatePlan_WithAlias(t *testing.T) {
	registry := tools.NewRegistry()
	registry.Register(&mockTool{
		name:   "stage_all",
		params: map[string]tools.ParamSchema{},
		perm:   tools.PermWriteLocal,
	})

	agent := &TwoPhaseAgent{
		fullRegistry: registry,
		provider:     &mockProvider{},
	}

	plan := &ExecutionPlan{
		Summary: "Test plan",
		Steps: []*PlanStep{
			{
				ID:          "1",
				Description: "Test step",
				ToolName:    "add", // 使用别名
				Params:      map[string]any{},
			},
		},
	}

	errors := agent.validatePlan(plan)
	assert.Empty(t, errors, "Should handle alias correctly")
}

func TestValidateParamType(t *testing.T) {
	tests := []struct {
		name         string
		paramName    string
		value        any
		expectedType string
		wantError    bool
	}{
		{"string valid", "msg", "hello", "string", false},
		{"string invalid", "msg", 123, "string", true},
		{"int valid - int", "count", 10, "int", false},
		{"int valid - float64", "count", float64(10), "int", false},
		{"int invalid", "count", "ten", "int", true},
		{"bool valid", "flag", true, "bool", false},
		{"bool invalid", "flag", "true", "bool", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateParamType(tt.paramName, tt.value, tt.expectedType)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToolCallHistory_PreventInfiniteLoop(t *testing.T) {
	agent := &TwoPhaseAgent{
		state: GraphState{ToolCallHistory: make(map[string]int)},
	}

	// 模拟重复调用
	callKey := "get_status:{}"

	for i := 1; i <= 5; i++ {
		agent.state.ToolCallHistory[callKey]++

		if agent.state.ToolCallHistory[callKey] > 3 {
			// 应该被阻止
			assert.Greater(t, agent.state.ToolCallHistory[callKey], 3,
				"Should detect repeated calls")
			break
		}
	}

	assert.Equal(t, 4, agent.state.ToolCallHistory[callKey],
		"Should stop after 3 successful calls")
}
