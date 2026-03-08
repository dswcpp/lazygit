package skills

import (
	"context"
	"errors"
	"strings"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

// StashNameSkill generates a descriptive stash message for the current working-tree changes.
//
// Extra keys:
//   - "diff"     (string, optional) the working-tree diff to summarise
//   - "context"  (string, optional) extra hint from the user about what they're doing
type StashNameSkill struct{}

func NewStashNameSkill() Skill { return &StashNameSkill{} }
func (s *StashNameSkill) Name() string { return "stash_name" }

func (s *StashNameSkill) Execute(ctx context.Context, p provider.Provider, input Input) (Output, error) {
	diff := extraStr(input.Extra, "diff")
	userContext := extraStr(input.Extra, "context")

	if strings.TrimSpace(diff) == "" && strings.TrimSpace(userContext) == "" {
		return Output{}, errors.New("both diff and context are empty — cannot generate stash name")
	}

	var messages []provider.Message
	if input.Tr != nil {
		messages = append(messages, provider.Message{
			Role:    provider.RoleSystem,
			Content: input.Tr.SkillStashNameSystemPrompt(),
		})
	}
	messages = append(messages, provider.Message{
		Role:    provider.RoleUser,
		Content: buildStashNamePrompt(diff, userContext),
	})

	result, err := p.Complete(ctx, messages)
	if err != nil {
		return Output{}, err
	}
	name := strings.TrimSpace(result.Content)
	// Strip surrounding quotes that LLMs sometimes add
	name = strings.Trim(name, `"'`)
	name = strings.TrimSpace(name)
	if name == "" {
		return Output{}, errors.New("AI returned empty stash name")
	}
	return Output{Content: name}, nil
}

func buildStashNamePrompt(diff, userContext string) string {
	var sb strings.Builder

	if userContext != "" {
		sb.WriteString("**Current task:** " + userContext + "\n\n")
	}

	if diff != "" {
		// Extract changed file list from unified diff headers
		var changedFiles []string
		for _, line := range strings.Split(diff, "\n") {
			if strings.HasPrefix(line, "+++ b/") {
				changedFiles = append(changedFiles, strings.TrimPrefix(line, "+++ b/"))
			}
		}
		if len(changedFiles) > 0 {
			sb.WriteString("**Changed files:**\n")
			maxFiles := 10
			for i, f := range changedFiles {
				if i >= maxFiles {
					sb.WriteString("  ... and more\n")
					break
				}
				sb.WriteString("  - " + f + "\n")
			}
			sb.WriteString("\n")
		}

		// Include a diff excerpt (first 60 lines) for content signal
		lines := strings.Split(diff, "\n")
		limit := 60
		if len(lines) < limit {
			limit = len(lines)
		}
		excerpt := strings.Join(lines[:limit], "\n")
		sb.WriteString("**Diff excerpt:**\n```diff\n")
		sb.WriteString(excerpt)
		sb.WriteString("\n```\n\n")
	}

	sb.WriteString("Generate a single-line stash message (no quotes, no prefix).\n")
	sb.WriteString("Format: `<type>: <concise description in Simplified Chinese>` (max 60 chars)\n")
	sb.WriteString("Examples: `wip: 登录页面样式调整未完成`, `fix: 修复空指针但测试未通过`\n")

	return sb.String()
}
