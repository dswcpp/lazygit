package gittools

import (
	"context"
	"testing"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
	"github.com/dswcpp/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

// mockProvider 模拟AI provider用于测试
type mockProvider struct {
	response string
	err      error
}

func (m *mockProvider) Complete(ctx context.Context, messages []provider.Message) (provider.Result, error) {
	if m.err != nil {
		return provider.Result{}, m.err
	}
	return provider.Result{Content: m.response}, nil
}

func (m *mockProvider) CompleteStream(ctx context.Context, messages []provider.Message, onChunk func(string)) error {
	if m.err != nil {
		return m.err
	}
	onChunk(m.response)
	return nil
}

func (m *mockProvider) ModelID() string {
	return "mock-model"
}

func TestAnalyzeChangesTool_NoChanges(t *testing.T) {
	deps := &Deps{
		GetFiles: func() []*models.File { return []*models.File{} },
	}
	prov := &mockProvider{response: "分析结果"}
	tool := NewAnalyzeChangesTool(deps, prov)

	result := tool.Execute(context.Background(), tools.ToolCall{ID: "test"})

	assert.True(t, result.Success)
	assert.Contains(t, result.Output, "工作区干净")
}

func TestAnalyzeChangesTool_Schema(t *testing.T) {
	deps := &Deps{}
	prov := &mockProvider{}
	tool := NewAnalyzeChangesTool(deps, prov)

	schema := tool.Schema()

	assert.Equal(t, "analyze_changes", schema.Name)
	assert.Contains(t, schema.Description, "智能分析")
	assert.Equal(t, tools.PermReadOnly, schema.Permission)
	assert.Contains(t, schema.Params, "staged")
	assert.Contains(t, schema.Params, "focus")
}

func TestAnalyzeChangesTool_BuildSummary(t *testing.T) {
	deps := &Deps{}
	prov := &mockProvider{}
	tool := NewAnalyzeChangesTool(deps, prov).(*AnalyzeChangesTool)

	analyses := []fileAnalysis{
		{Path: "file1.go", Summary: "添加了新功能", LinesDiff: 50},
		{Path: "file2.go", Summary: "修复了bug", LinesDiff: 20},
		{Path: "file3.go", Error: "分析失败"},
	}

	summary := tool.buildSummary(analyses, "")

	assert.Contains(t, summary, "变更分析报告")
	assert.Contains(t, summary, "file1.go")
	assert.Contains(t, summary, "file2.go")
	assert.Contains(t, summary, "file3.go")
	assert.Contains(t, summary, "成功分析 2")
	assert.Contains(t, summary, "失败 1")
	assert.Contains(t, summary, "约 70 行")
}

func TestAnalyzeChangesTool_BuildAnalysisPrompt(t *testing.T) {
	deps := &Deps{}
	prov := &mockProvider{}
	tool := NewAnalyzeChangesTool(deps, prov).(*AnalyzeChangesTool)

	prompt := tool.buildAnalysisPrompt("test.go", "diff content", "安全问题")

	assert.Contains(t, prompt, "test.go")
	assert.Contains(t, prompt, "diff content")
	assert.Contains(t, prompt, "安全问题")
}
