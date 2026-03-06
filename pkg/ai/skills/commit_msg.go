package skills

import (
	"context"
	"errors"
	"fmt"
	"strings"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
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

	systemPrompt := input.Tr.SkillCommitMsgSystemPrompt()
	userPrompt := buildCommitMsgUserPrompt(input.Tr, diff, input.RepoCtx.CurrentBranch, projectType, scenario, safetyNote)

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

func buildCommitMsgUserPrompt(tr *aii18n.Translator, diff, branch, projectType, scenario, safetyNote string) string {
	var sb strings.Builder

	sb.WriteString(tr.SkillRepoBackground())
	if branch != "" {
		sb.WriteString(tr.SkillCurrentBranch(branch))
	}
	if projectType != "" && projectType != "Mixed" {
		sb.WriteString(fmt.Sprintf("项目类型: %s\n", projectType))
	}
	sb.WriteString("\n")

	sb.WriteString(tr.SkillCodeChangesSection())
	sb.WriteString("```diff\n")
	sb.WriteString(diff)
	if safetyNote != "" {
		sb.WriteString("\n" + safetyNote)
	}
	sb.WriteString("\n```\n\n")

	sb.WriteString(tr.SkillOutputRules())
	sb.WriteString(tr.SkillFormatExample())
	sb.WriteString(tr.SkillTypeList())
	sb.WriteString(tr.SkillSubjectRules())
	sb.WriteString(tr.SkillScopeOptional())
	sb.WriteString(tr.SkillBodyRequired())

	switch scenario {
	case "bugfix":
		sb.WriteString(tr.SkillScenarioBugfix())
	case "refactor":
		sb.WriteString(tr.SkillScenarioRefactor())
	case "docs":
		sb.WriteString(tr.SkillScenarioDocs())
	case "test":
		sb.WriteString(tr.SkillScenarioTest())
	default:
		sb.WriteString(tr.SkillScenarioDefault())
	}

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
