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
//   - "diff" (string, optional) diff of all changes (truncated if too large)
type BranchNameSkill struct{}

func NewBranchNameSkill() Skill { return &BranchNameSkill{} }
func (s *BranchNameSkill) Name() string { return "branch_name" }

func (s *BranchNameSkill) Execute(ctx context.Context, p provider.Provider, input Input) (Output, error) {
	diff := extraStr(input.Extra, "diff")

	var sb strings.Builder
	sb.WriteString(input.Tr.SkillBranchNamePromptIntro())

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
		name = "feature/" + name
	}
	return Output{Content: name}, nil
}

func cleanBranchName(raw string) string {
	// Strip surrounding quotes/backticks
	raw = strings.Trim(raw, "\"'`")
	raw = strings.ReplaceAll(raw, " ", "-")
	raw = strings.ToLower(raw)
	for _, ch := range []string{"~", "^", ":", "?", "*", "[", "\\", "..", "@{", "//"} {
		raw = strings.ReplaceAll(raw, ch, "")
	}
	return raw
}
