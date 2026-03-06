package skills

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

// CommandSuggestion is one AI-generated shell command with metadata.
type CommandSuggestion struct {
	Command      string `json:"command"`
	Explanation  string `json:"explanation"`
	RiskLevel    string `json:"risk_level"` // "safe" | "medium" | "dangerous"
	Alternatives string `json:"alternatives,omitempty"`
}

// ShellCmdSkill translates a natural-language intent into shell commands.
// Extra keys:
//   - "intent" (string, required) natural language description of what to do
type ShellCmdSkill struct{}

func NewShellCmdSkill() Skill { return &ShellCmdSkill{} }
func (s *ShellCmdSkill) Name() string { return "shell_cmd" }

func (s *ShellCmdSkill) Execute(ctx context.Context, p provider.Provider, input Input) (Output, error) {
	intent := extraStr(input.Extra, "intent")
	if strings.TrimSpace(intent) == "" {
		return Output{}, errors.New("intent is empty")
	}

	osHint := shellHint(input.Tr)
	repoSummary := input.RepoCtx.CompactString(input.Tr)

	prompt := input.Tr.SkillShellCmdSystemPrompt() +
		input.Tr.SkillShellCmdRuntime(osHint) +
		input.Tr.SkillShellCmdRepoStatus(repoSummary) +
		input.Tr.SkillShellCmdUserIntent(intent) +
		input.Tr.SkillShellCmdOutputFormat() +
		input.Tr.SkillShellCmdCommandField() +
		input.Tr.SkillShellCmdExplanationField() +
		input.Tr.SkillShellCmdRiskLevelField() +
		input.Tr.SkillShellCmdAlternativesField() +
		input.Tr.SkillShellCmdOutputNote()

	messages := []provider.Message{
		{Role: provider.RoleUser, Content: prompt},
	}

	result, err := p.Complete(ctx, messages)
	if err != nil {
		return Output{}, err
	}

	suggestions, err := parseCommandSuggestions(result.Content)
	if err != nil {
		return Output{}, fmt.Errorf("parse suggestions: %w", err)
	}

	// Re-serialise as canonical JSON for the caller
	data, _ := json.Marshal(suggestions)
	return Output{
		Content: string(data),
		Data:    map[string]any{"suggestions": suggestions},
	}, nil
}

func shellHint(tr *aii18n.Translator) string {
	switch runtime.GOOS {
	case "windows":
		return tr.SkillShellCmdWindowsHint()
	case "darwin":
		return tr.SkillShellCmdMacOSHint()
	default:
		return tr.SkillShellCmdLinuxHint()
	}
}

func parseCommandSuggestions(raw string) ([]CommandSuggestion, error) {
	content := strings.TrimSpace(raw)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var suggestions []CommandSuggestion
	if err := json.Unmarshal([]byte(content), &suggestions); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w (response: %s)", err, content)
	}
	for i := range suggestions {
		if suggestions[i].RiskLevel != "safe" &&
			suggestions[i].RiskLevel != "medium" &&
			suggestions[i].RiskLevel != "dangerous" {
			suggestions[i].RiskLevel = "medium"
		}
	}
	return suggestions, nil
}
