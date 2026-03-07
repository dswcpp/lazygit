package skills

import (
	"context"
	"errors"
	"fmt"
	"strings"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
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

	userPrompt := buildPRDescUserPrompt(input.Tr, input.RepoCtx, fromBranch, toBranch, diff)

	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: input.Tr.SkillPRDescSystemPrompt()},
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

func buildPRDescUserPrompt(tr *aii18n.Translator, ctx repocontext.RepoContext, fromBranch, toBranch, diff string) string {
	var sb strings.Builder

	sb.WriteString(tr.SkillPRDescBranchInfo(fromBranch, toBranch))

	if len(ctx.RecentCommits) > 0 {
		sb.WriteString(tr.SkillPRDescCommitHistory())
		for _, c := range ctx.RecentCommits {
			sb.WriteString(fmt.Sprintf("- %s %s (%s)\n", c.ShortHash, c.Message, c.Author))
		}
		sb.WriteString("\n")
	}

	if diff != "" {
		sb.WriteString(tr.SkillPRDescCodeChangesSection())
		sb.WriteString(diff)
		sb.WriteString("\n```\n\n")
	}

	sb.WriteString(tr.SkillPRDescGeneratePrompt())
	sb.WriteString(tr.SkillPRDescSummarySection())
	sb.WriteString(tr.SkillPRDescChangesSection())
	if hasBreakingChanges(ctx.RecentCommits, diff) {
		sb.WriteString(tr.SkillPRDescBreakingSection())
	}
	sb.WriteString(tr.SkillPRDescTestingSection())
	sb.WriteString(tr.SkillPRDescChecklistSection())

	return sb.String()
}

// hasBreakingChanges returns true when the commits or diff contain conventional
// breaking-change indicators ("BREAKING CHANGE" footer or "!" in type).
func hasBreakingChanges(commits []repocontext.CommitSummary, diff string) bool {
	for _, c := range commits {
		msg := c.Message
		if strings.Contains(msg, "BREAKING CHANGE") ||
			strings.Contains(msg, "BREAKING-CHANGE") ||
			breakingTypeRe(msg) {
			return true
		}
	}
	return strings.Contains(diff, "BREAKING CHANGE")
}

// breakingTypeRe checks for conventional commit "type!:" pattern.
func breakingTypeRe(msg string) bool {
	for _, t := range []string{"feat!", "fix!", "refactor!", "perf!", "chore!"} {
		if strings.HasPrefix(msg, t) {
			return true
		}
	}
	return false
}
