package skills

import (
	"context"
	"errors"
	"strings"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

// ExplainDiffSkill explains what a diff does in plain language.
// Unlike CodeReviewSkill it describes rather than judges — no severity levels,
// no issue list; just a clear explanation of intent and impact.
//
// Extra keys:
//   - "diff"      (string, required) the diff to explain
//   - "file_path" (string, optional) file path for language context
//   - "audience"  (string, optional) "developer" | "reviewer" | "beginner"
type ExplainDiffSkill struct{}

func NewExplainDiffSkill() Skill { return &ExplainDiffSkill{} }
func (s *ExplainDiffSkill) Name() string { return "explain_diff" }

func (s *ExplainDiffSkill) Execute(ctx context.Context, p provider.Provider, input Input) (Output, error) {
	diff := extraStr(input.Extra, "diff")
	if strings.TrimSpace(diff) == "" {
		return Output{}, errors.New("diff is empty")
	}
	filePath := extraStr(input.Extra, "file_path")
	audience := extraStr(input.Extra, "audience")

	var messages []provider.Message
	if input.Tr != nil {
		messages = append(messages, provider.Message{
			Role:    provider.RoleSystem,
			Content: input.Tr.SkillExplainDiffSystemPrompt(),
		})
	}
	messages = append(messages, provider.Message{
		Role:    provider.RoleUser,
		Content: buildExplainDiffPrompt(filePath, audience, diff),
	})

	result, err := p.Complete(ctx, messages)
	if err != nil {
		return Output{}, err
	}
	content := strings.TrimSpace(result.Content)
	if content == "" {
		return Output{}, errors.New("AI returned empty explanation")
	}
	return Output{Content: content}, nil
}

func buildExplainDiffPrompt(filePath, audience, diff string) string {
	var sb strings.Builder

	if filePath != "" {
		sb.WriteString("**File:** `" + filePath + "`\n\n")
	}

	switch strings.ToLower(audience) {
	case "reviewer":
		sb.WriteString("Audience: code reviewer — focus on intent and potential impact on other modules.\n\n")
	case "beginner":
		sb.WriteString("Audience: developer less familiar with the codebase — be extra clear about what and why.\n\n")
	default:
		sb.WriteString("Audience: developer — be technical but concise.\n\n")
	}

	sb.WriteString("Explain this diff in three short sections:\n\n")
	sb.WriteString("1. **What changed** — describe the code changes clearly.\n")
	sb.WriteString("2. **Why it matters** — what problem does it solve or what feature does it add?\n")
	sb.WriteString("3. **Notable effects** — any impact on behaviour, performance, or API surface?\n\n")
	sb.WriteString("Keep each section to 2–4 sentences. Use Simplified Chinese.\n\n")
	sb.WriteString("```diff\n")
	sb.WriteString(diff)
	sb.WriteString("\n```")

	return sb.String()
}
