package skills

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

// ReleaseNotesSkill generates a changelog / release notes section from
// the repository's recent commit history.
//
// Extra keys:
//   - "from_tag" (string, optional) previous tag or version, e.g. "v1.0.0"
//   - "to_tag"   (string, optional) new tag or version,      e.g. "v1.1.0"
type ReleaseNotesSkill struct{}

func NewReleaseNotesSkill() Skill { return &ReleaseNotesSkill{} }
func (s *ReleaseNotesSkill) Name() string { return "release_notes" }

func (s *ReleaseNotesSkill) Execute(ctx context.Context, p provider.Provider, input Input) (Output, error) {
	if len(input.RepoCtx.RecentCommits) == 0 {
		return Output{}, errors.New("no commits available in repository context")
	}

	fromTag := extraStr(input.Extra, "from_tag")
	toTag := extraStr(input.Extra, "to_tag")

	var messages []provider.Message
	if input.Tr != nil {
		messages = append(messages, provider.Message{
			Role:    provider.RoleSystem,
			Content: input.Tr.SkillReleaseNotesSystemPrompt(),
		})
	}
	messages = append(messages, provider.Message{
		Role:    provider.RoleUser,
		Content: buildReleaseNotesPrompt(input, fromTag, toTag),
	})

	result, err := p.Complete(ctx, messages)
	if err != nil {
		return Output{}, err
	}
	content := strings.TrimSpace(result.Content)
	if content == "" {
		return Output{}, errors.New("AI returned empty release notes")
	}
	return Output{Content: content}, nil
}

func buildReleaseNotesPrompt(input Input, fromTag, toTag string) string {
	var sb strings.Builder

	// Version range header
	from := fromTag
	if from == "" {
		from = "previous release"
	}
	to := toTag
	if to == "" {
		to = "HEAD"
	}
	sb.WriteString(fmt.Sprintf("## Version range: %s → %s\n\n", from, to))

	// Raw commit list
	sb.WriteString("## Commits\n")
	for _, c := range input.RepoCtx.RecentCommits {
		sb.WriteString(fmt.Sprintf("- `%s` %s (%s)\n", c.ShortHash, c.Message, c.Author))
	}
	sb.WriteString("\n")

	// Output instructions
	sb.WriteString("## Instructions\n")
	sb.WriteString("Generate professional release notes in Markdown.\n\n")
	sb.WriteString("Rules:\n")
	sb.WriteString("1. Group entries under these headings (omit empty groups):\n")
	sb.WriteString("   - **✨ 新功能** (feat)\n")
	sb.WriteString("   - **🐛 Bug 修复** (fix)\n")
	sb.WriteString("   - **⚡ 性能优化** (perf)\n")
	sb.WriteString("   - **♻️ 重构** (refactor)\n")
	sb.WriteString("   - **📝 文档** (docs)\n")
	sb.WriteString("   - **🔧 其他** (chore/ci/style/test)\n")
	sb.WriteString("2. Rewrite each entry as a user-facing description — don't just copy the commit subject.\n")
	sb.WriteString("3. If any commit signals a breaking change (`BREAKING CHANGE` or `type!:`), add a **⚠️ Breaking Changes** section at the very top.\n")
	sb.WriteString("4. Skip trivial chore/ci/style commits unless they're meaningful to users.\n")
	sb.WriteString("5. Write all descriptions in Simplified Chinese.\n")

	return sb.String()
}
