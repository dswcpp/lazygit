package skills

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/repocontext"
)

// PRDescSkill generates a pull request description from commit history and diff.
// Extra keys:
//   - "diff"        (string) diff between base branch and HEAD (may be truncated)
//   - "from_branch" (string) source branch name
//   - "to_branch"   (string) target branch name
type PRDescSkill struct{}

func NewPRDescSkill() Skill { return &PRDescSkill{} }
func (s *PRDescSkill) Name() string { return "pr_desc" }

func (s *PRDescSkill) Execute(ctx context.Context, p provider.Provider, input Input) (Output, error) {
	diff := extraStr(input.Extra, "diff")
	fromBranch := extraStr(input.Extra, "from_branch")
	toBranch := extraStr(input.Extra, "to_branch")
	if toBranch == "" {
		toBranch = "main"
	}

	userPrompt := buildPRDescUserPrompt(input.RepoCtx, fromBranch, toBranch, diff)

	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: prDescSystemPrompt()},
		{Role: provider.RoleUser, Content: userPrompt},
	}

	result, err := p.Complete(ctx, messages)
	if err != nil {
		return Output{}, err
	}
	content := strings.TrimSpace(result.Content)
	if content == "" {
		return Output{}, errors.New("AI returned empty PR description")
	}
	return Output{Content: content}, nil
}

func prDescSystemPrompt() string {
	return `你是一名高级软件工程师，负责撰写清晰、专业的 Pull Request 描述。
用 Markdown 格式输出，包含 Summary、Changes 和 Testing 三个部分。
使用简洁的中文。`
}

func buildPRDescUserPrompt(ctx repocontext.RepoContext, fromBranch, toBranch, diff string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## 分支信息\n从 `%s` 合并到 `%s`\n\n", fromBranch, toBranch))

	if len(ctx.RecentCommits) > 0 {
		sb.WriteString("## 提交历史\n")
		for _, c := range ctx.RecentCommits {
			sb.WriteString(fmt.Sprintf("- %s %s (%s)\n", c.ShortHash, c.Message, c.Author))
		}
		sb.WriteString("\n")
	}

	if diff != "" {
		sb.WriteString("## 代码变更\n```diff\n")
		sb.WriteString(diff)
		sb.WriteString("\n```\n\n")
	}

	sb.WriteString("## 请生成包含以下部分的 PR 描述\n")
	sb.WriteString("### Summary\n一句话说明 PR 的目的。\n\n")
	sb.WriteString("### Changes\n- 列出主要变更（3-5 条）\n\n")
	sb.WriteString("### Testing\n- 说明如何验证这些变更\n")

	return sb.String()
}
