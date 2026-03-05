package skills

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

// CommitMsgSkill generates a Conventional Commits message from a staged diff.
// Extra keys:
//   - "diff"         (string, required) staged diff text
//   - "project_type" (string, optional) detected language/project type
//   - "safety_note"  (string, optional) appended when diff was truncated
type CommitMsgSkill struct{}

func NewCommitMsgSkill() Skill { return &CommitMsgSkill{} }
func (s *CommitMsgSkill) Name() string { return "commit_msg" }

func (s *CommitMsgSkill) Execute(ctx context.Context, p provider.Provider, input Input) (Output, error) {
	diff := extraStr(input.Extra, "diff")
	if strings.TrimSpace(diff) == "" {
		return Output{}, errors.New("staged diff is empty")
	}
	projectType := extraStr(input.Extra, "project_type")
	if projectType == "" {
		projectType = "Mixed"
	}
	safetyNote := extraStr(input.Extra, "safety_note")
	scenario := detectChangeScenario(diff)

	systemPrompt := commitMsgSystemPrompt()
	userPrompt := buildCommitMsgUserPrompt(diff, input.RepoCtx.CurrentBranch, projectType, scenario, safetyNote)

	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: systemPrompt},
		{Role: provider.RoleUser, Content: userPrompt},
	}

	result, err := p.Complete(ctx, messages)
	if err != nil {
		return Output{}, err
	}
	content := strings.TrimSpace(result.Content)
	if content == "" {
		return Output{}, errors.New("AI returned empty commit message")
	}
	return Output{Content: content}, nil
}

// ── prompt builders ────────────────────────────────────────────────────────

func commitMsgSystemPrompt() string {
	return `你是一名经验丰富的软件工程师，专门负责撰写高质量的 Git 提交信息。
遵循 Conventional Commits 规范（https://www.conventionalcommits.org/）。
只输出提交信息本身，不要任何额外说明、前缀或引号。
提交信息（subject 和 body）必须使用中文。`
}

func buildCommitMsgUserPrompt(diff, branch, projectType, scenario, safetyNote string) string {
	var sb strings.Builder

	sb.WriteString("## 仓库背景\n")
	if branch != "" {
		sb.WriteString(fmt.Sprintf("当前分支: %s\n", branch))
	}
	if projectType != "" && projectType != "Mixed" {
		sb.WriteString(fmt.Sprintf("项目类型: %s\n", projectType))
	}
	sb.WriteString("\n")

	sb.WriteString("## 代码变更\n")
	sb.WriteString("```diff\n")
	sb.WriteString(diff)
	if safetyNote != "" {
		sb.WriteString("\n" + safetyNote)
	}
	sb.WriteString("\n```\n\n")

	sb.WriteString("## 输出规则\n")
	sb.WriteString("- 格式:\n")
	sb.WriteString("  ```\n")
	sb.WriteString("  <type>(<scope>): <subject>\n")
	sb.WriteString("  \n")
	sb.WriteString("  <body>\n")
	sb.WriteString("  ```\n")
	sb.WriteString("- type: feat | fix | refactor | docs | test | chore | perf | style | ci | revert\n")
	sb.WriteString("- subject: 中文，动词开头，祈使句，不超过 72 字符\n")
	sb.WriteString("- scope 可省略\n")
	sb.WriteString("- body: 必须包含，与 subject 之间空一行，用中文说明本次变更的原因和主要内容（1-4 行）\n\n")

	switch scenario {
	case "bugfix":
		sb.WriteString("场景提示: 这是一个 bug 修复，优先使用 fix 类型。\n")
	case "refactor":
		sb.WriteString("场景提示: 这是重构，优先使用 refactor 类型。\n")
	case "docs":
		sb.WriteString("场景提示: 这是文档更新，使用 docs 类型。\n")
	case "test":
		sb.WriteString("场景提示: 这是测试相关变更，使用 test 类型。\n")
	case "large":
		sb.WriteString("场景提示: 变更较大，body 须逐点列举主要变更，subject 保持简洁。\n")
	}

	sb.WriteString("\n请直接输出提交信息：")
	return sb.String()
}

func detectChangeScenario(diff string) string {
	lower := strings.ToLower(diff)
	switch {
	case strings.Contains(diff, ".md") || strings.Contains(diff, "README"):
		return "docs"
	case strings.Contains(diff, "_test.") || strings.Contains(diff, ".test.") || strings.Contains(diff, "/test/"):
		return "test"
	case strings.Contains(lower, "fix") || strings.Contains(lower, "bug") || strings.Contains(lower, "error"):
		return "bugfix"
	case strings.Contains(lower, "refactor") || strings.Contains(lower, "rename") || strings.Contains(lower, "move"):
		return "refactor"
	default:
		lines := strings.Split(diff, "\n")
		changes := 0
		for _, l := range lines {
			if strings.HasPrefix(l, "+") || strings.HasPrefix(l, "-") {
				changes++
			}
		}
		if changes > 500 {
			return "large"
		}
		if changes < 50 {
			return "small"
		}
		return "normal"
	}
}
