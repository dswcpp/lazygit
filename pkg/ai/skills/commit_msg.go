package skills

import (
	"context"
	"errors"
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
		sb.WriteString(tr.SkillCommitMsgProjectType(projectType))
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
	case "large":
		sb.WriteString(tr.SkillScenarioLarge())
	default:
		sb.WriteString(tr.SkillScenarioDefault())
	}

	return sb.String()
}

// detectChangeScenario classifies the diff by analysing changed file paths
// (from both "diff --git a/x b/x" and "+++ b/x" unified-diff headers) and
// falling back to content keywords and size heuristics when no headers exist.
func detectChangeScenario(diff string) string {
	var changedFiles []string
	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "+++ b/"):
			changedFiles = append(changedFiles, strings.TrimPrefix(line, "+++ b/"))
		case strings.HasPrefix(line, "diff --git "):
			// "diff --git a/README.md b/README.md" → extract "README.md"
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				changedFiles = append(changedFiles, strings.TrimPrefix(parts[3], "b/"))
			}
		}
	}

	if len(changedFiles) > 0 {
		docsCount, testCount := 0, 0
		for _, f := range changedFiles {
			lower := strings.ToLower(f)
			if isDocFile(lower) {
				docsCount++
			} else if isTestFile(lower) {
				testCount++
			}
		}
		total := len(changedFiles)
		if docsCount == total {
			return "docs"
		}
		if testCount == total {
			return "test"
		}
		// File names are reliable signals; scan for fix/refactor keywords.
		for _, f := range changedFiles {
			lower := strings.ToLower(f)
			if strings.Contains(lower, "fix") || strings.Contains(lower, "bug") || strings.Contains(lower, "patch") {
				return "bugfix"
			}
			if strings.Contains(lower, "refactor") {
				return "refactor"
			}
		}
	}

	// Count changed lines for size heuristics.
	added, removed := 0, 0
	for _, line := range strings.Split(diff, "\n") {
		if len(line) > 0 && line[0] == '+' && !strings.HasPrefix(line, "+++") {
			added++
		} else if len(line) > 0 && line[0] == '-' && !strings.HasPrefix(line, "---") {
			removed++
		}
	}
	if added+removed > 500 {
		return "large"
	}

	// Content-based keyword detection for diffs that lack file headers.
	if len(changedFiles) == 0 {
		lower := strings.ToLower(diff)
		if strings.Contains(lower, "fix") || strings.Contains(lower, "bug") {
			return "bugfix"
		}
		if strings.Contains(lower, "refactor") {
			return "refactor"
		}
	}

	if added+removed <= 5 {
		return "small"
	}
	return "normal"
}

func isDocFile(lower string) bool {
	return strings.HasSuffix(lower, ".md") ||
		strings.HasSuffix(lower, ".txt") ||
		strings.HasSuffix(lower, ".rst") ||
		strings.HasSuffix(lower, ".adoc") ||
		strings.Contains(lower, "readme") ||
		strings.HasPrefix(lower, "docs/") ||
		strings.HasPrefix(lower, "doc/")
}

func isTestFile(lower string) bool {
	return strings.HasSuffix(lower, "_test.go") ||
		strings.Contains(lower, "_test.") ||
		strings.Contains(lower, ".test.") ||
		strings.Contains(lower, ".spec.") ||
		strings.HasPrefix(lower, "test/") ||
		strings.HasPrefix(lower, "tests/") ||
		strings.HasPrefix(lower, "__tests__/")
}
