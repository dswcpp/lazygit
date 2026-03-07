package skills

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

// BranchNameSkill suggests a branch name based on staged/unstaged changes.
// Extra keys:
//   - "diff"        (string, optional) diff of all changes (truncated if too large)
//   - "description" (string, optional) natural-language intent for the branch (e.g. from planning LLM)
type BranchNameSkill struct{}

func NewBranchNameSkill() Skill { return &BranchNameSkill{} }
func (s *BranchNameSkill) Name() string { return "branch_name" }

func (s *BranchNameSkill) Execute(ctx context.Context, p provider.Provider, input Input) (Output, error) {
	diff := extraStr(input.Extra, "diff")
	description := extraStr(input.Extra, "description")

	var sb strings.Builder
	sb.WriteString(input.Tr.SkillBranchNamePromptIntro())

	if description != "" {
		sb.WriteString(input.Tr.SkillBranchNameDescriptionHint(description))
	}

	// Files summary from context
	staged := []string{}
	unstaged := []string{}
	for _, f := range input.RepoCtx.Files {
		if f.HasStaged {
			staged = append(staged, f.Path)
		} else if f.HasUnstaged {
			unstaged = append(unstaged, f.Path)
		}
	}

	if len(staged) > 0 {
		sb.WriteString(input.Tr.SkillBranchNameStagedFiles())
		for i, f := range staged {
			if i >= 15 {
				sb.WriteString(input.Tr.SkillBranchNameMoreFiles(len(staged) - 15))
				break
			}
			sb.WriteString(fmt.Sprintf("  - %s\n", f))
		}
	}
	if len(unstaged) > 0 {
		sb.WriteString(input.Tr.SkillBranchNameUnstagedFiles())
		for i, f := range unstaged {
			if i >= 15 {
				sb.WriteString(input.Tr.SkillBranchNameMoreFiles(len(unstaged) - 15))
				break
			}
			sb.WriteString(fmt.Sprintf("  - %s\n", f))
		}
	}

	if diff != "" {
		sb.WriteString(input.Tr.SkillBranchNameDiffSummaryTitle())
		sb.WriteString(diff)
		sb.WriteString("\n```\n")
	}

	sb.WriteString(input.Tr.SkillBranchNameRules())
	sb.WriteString(input.Tr.SkillBranchNameFormatRule())
	sb.WriteString(input.Tr.SkillBranchNameTypeRule())
	sb.WriteString(input.Tr.SkillBranchNameDescRule())
	sb.WriteString(input.Tr.SkillBranchNameOutputRule())

	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: input.Tr.SkillBranchNameSystemPrompt()},
		{Role: provider.RoleUser, Content: sb.String()},
	}

	result, err := p.Complete(ctx, messages)
	if err != nil {
		return Output{}, err
	}

	name := cleanBranchName(strings.TrimSpace(result.Content))
	if name == "" {
		return Output{}, errors.New("AI returned empty branch name")
	}
	if !strings.Contains(name, "/") {
		name = inferBranchType(description, diff) + "/" + name
	}
	return Output{Content: name}, nil
}

// inferBranchType picks a branch type prefix from description/diff signals.
// Used as fallback when the AI omits the type prefix.
func inferBranchType(description, diff string) string {
	lower := strings.ToLower(description + " " + diff)
	switch {
	case strings.Contains(lower, "fix") || strings.Contains(lower, "bug") || strings.Contains(lower, "patch"):
		return "fix"
	case strings.Contains(lower, "refactor") || strings.Contains(lower, "cleanup") || strings.Contains(lower, "clean up"):
		return "refactor"
	case strings.Contains(lower, "doc") || strings.Contains(lower, "readme") || strings.Contains(lower, "changelog"):
		return "docs"
	case strings.Contains(lower, "test") || strings.Contains(lower, "spec"):
		return "test"
	case strings.Contains(lower, "chore") || strings.Contains(lower, "dep") || strings.Contains(lower, "upgrade") || strings.Contains(lower, "bump"):
		return "chore"
	default:
		return "feature"
	}
}

func cleanBranchName(raw string) string {
	// Strip surrounding quotes/backticks
	raw = strings.Trim(raw, "\"'`")
	raw = strings.ReplaceAll(raw, " ", "-")
	raw = strings.ToLower(raw)
	// Remove all Git-invalid characters
	for _, ch := range []string{"~", "^", ":", "?", "*", "[", "\\", "..", "@{", "//", "{", "}", "!", "@", "#", "$", "%", "&", "+"} {
		raw = strings.ReplaceAll(raw, ch, "")
	}
	// Remove leading dots (Git forbids branch names starting with '.')
	raw = strings.TrimLeft(raw, ".")
	// Remove trailing '.lock' suffix (reserved by Git)
	raw = strings.TrimSuffix(raw, ".lock")
	// Collapse consecutive hyphens
	for strings.Contains(raw, "--") {
		raw = strings.ReplaceAll(raw, "--", "-")
	}
	// Trim leading/trailing hyphens from each path component
	parts := strings.Split(raw, "/")
	for i, part := range parts {
		parts[i] = strings.Trim(part, "-")
	}
	return strings.Join(parts, "/")
}
