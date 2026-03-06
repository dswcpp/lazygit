package skills

import (
	"context"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/repocontext"
)

// Input is the standard argument passed to every Skill.
type Input struct {
	// RepoCtx is the current repository snapshot.
	RepoCtx repocontext.RepoContext
	// Extra holds skill-specific parameters (e.g. "diff", "from_branch", "command").
	Extra map[string]any
	// Tr is the translator for i18n
	Tr *aii18n.Translator
}

// Output is the result returned by a Skill.
type Output struct {
	// Content is the primary text result (commit message, branch name, review text…)
	Content string
	// Data holds optional structured results parsed from the model response.
	Data map[string]any
}

// Skill is the interface for a single, stateless AI capability.
// Each skill owns its own prompt construction and response parsing logic.
type Skill interface {
	Name() string
	Execute(ctx context.Context, p provider.Provider, input Input) (Output, error)
}

// extraStr is a helper to extract a string from Input.Extra.
func extraStr(extra map[string]any, key string) string {
	if v, ok := extra[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
