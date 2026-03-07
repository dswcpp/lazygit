package skills

import (
	"context"
	"strings"
	"testing"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/repocontext"
	"github.com/dswcpp/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

// mockProvider 用于测试
type mockCommitMsgProvider struct {
	response string
	err      error
}

func (m *mockCommitMsgProvider) Complete(ctx context.Context, messages []provider.Message) (provider.Result, error) {
	if m.err != nil {
		return provider.Result{}, m.err
	}
	return provider.Result{Content: m.response}, nil
}

func (m *mockCommitMsgProvider) CompleteStream(ctx context.Context, messages []provider.Message, onChunk func(string)) error {
	if m.err != nil {
		return m.err
	}
	onChunk(m.response)
	return nil
}

func (m *mockCommitMsgProvider) ModelID() string {
	return "mock-model"
}

func TestCommitMsgSkill_Name(t *testing.T) {
	skill := NewCommitMsgSkill()
	assert.Equal(t, "commit_msg", skill.Name())
}

func TestCommitMsgSkill_EmptyDiff(t *testing.T) {
	skill := NewCommitMsgSkill()
	prov := &mockCommitMsgProvider{response: "feat: add feature"}
	tr := aii18n.NewTranslator(i18n.EnglishTranslationSet())

	input := Input{
		RepoCtx: repocontext.RepoContext{CurrentBranch: "main"},
		Extra:   map[string]any{"diff": ""},
		Tr:      tr,
	}

	_, err := skill.Execute(context.Background(), prov, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestCommitMsgSkill_ValidDiff(t *testing.T) {
	skill := NewCommitMsgSkill()
	prov := &mockCommitMsgProvider{response: "feat(auth): 添加用户登录功能\n\n实现了基于 JWT 的用户认证系统"}
	tr := aii18n.NewTranslator(i18n.EnglishTranslationSet())

	input := Input{
		RepoCtx: repocontext.RepoContext{CurrentBranch: "feature/login"},
		Extra: map[string]any{
			"diff":         "+func Login() { ... }",
			"project_type": "Go",
		},
		Tr: tr,
	}

	output, err := skill.Execute(context.Background(), prov, input)
	assert.NoError(t, err)
	assert.NotEmpty(t, output.Content)
	assert.Contains(t, output.Content, "feat")
}

func TestCommitMsgSkill_EmptyResponse(t *testing.T) {
	skill := NewCommitMsgSkill()
	prov := &mockCommitMsgProvider{response: "   "} // 空白响应
	tr := aii18n.NewTranslator(i18n.EnglishTranslationSet())

	input := Input{
		RepoCtx: repocontext.RepoContext{},
		Extra:   map[string]any{"diff": "+some changes"},
		Tr:      tr,
	}

	_, err := skill.Execute(context.Background(), prov, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestDetectChangeScenario(t *testing.T) {
	tests := []struct {
		name     string
		diff     string
		expected string
	}{
		{
			name:     "文档变更",
			diff:     "diff --git a/README.md b/README.md\n+# New section",
			expected: "docs",
		},
		{
			name:     "测试文件",
			diff:     "diff --git a/user_test.go b/user_test.go\n+func TestLogin(t *testing.T) {}",
			expected: "test",
		},
		{
			name:     "Bug修复",
			diff:     "+// fix: resolve memory leak\n+defer close(ch)",
			expected: "bugfix",
		},
		{
			name:     "重构",
			diff:     "+// refactor: rename function\n+func NewName() {}",
			expected: "refactor",
		},
		{
			name:     "小变更",
			diff:     "+const MAX = 100",
			expected: "small",
		},
		{
			name:     "大变更",
			diff:     strings.Repeat("+line\n", 600),
			expected: "large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectChangeScenario(tt.diff)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildCommitMsgUserPrompt(t *testing.T) {
	tr := aii18n.NewTranslator(i18n.EnglishTranslationSet())
	diff := "+func Login() {}"
	branch := "feature/auth"
	projectType := "Go"
	scenario := "normal"
	safetyNote := ""

	prompt := buildCommitMsgUserPrompt(tr, diff, branch, projectType, scenario, safetyNote)

	assert.Contains(t, prompt, "feature/auth")
	assert.Contains(t, prompt, "Go")
	assert.Contains(t, prompt, "+func Login() {}")
	assert.Contains(t, prompt, "commit message")
}

func TestBuildCommitMsgUserPrompt_WithSafetyNote(t *testing.T) {
	tr := aii18n.NewTranslator(i18n.EnglishTranslationSet())
	diff := "+func Login() {}"
	safetyNote := "注意：diff 已截断"

	prompt := buildCommitMsgUserPrompt(tr, diff, "", "", "normal", safetyNote)

	assert.Contains(t, prompt, safetyNote)
}

func TestBuildCommitMsgUserPrompt_ScenarioHints(t *testing.T) {
	tests := []struct {
		scenario     string
		expectedHint string
	}{
		{"bugfix", "bug fix"},
		{"refactor", "refactor"},
		{"docs", "documentation"},
		{"test", "test"},
		{"large", "commit message"},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			tr := aii18n.NewTranslator(i18n.EnglishTranslationSet())
			prompt := buildCommitMsgUserPrompt(tr, "+code", "", "", tt.scenario, "")
			assert.Contains(t, prompt, tt.expectedHint)
		})
	}
}
