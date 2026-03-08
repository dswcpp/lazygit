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

	// System message: role + output format contract (stable, model-facing instructions)
	systemMsg := input.Tr.SkillShellCmdSystemPrompt() +
		input.Tr.SkillShellCmdOutputFormat() +
		input.Tr.SkillShellCmdCommandField() +
		input.Tr.SkillShellCmdExplanationField() +
		input.Tr.SkillShellCmdRiskLevelField() +
		input.Tr.SkillShellCmdAlternativesField() +
		input.Tr.SkillShellCmdOutputNote()

	// User message: runtime context + the actual request
	userMsg := input.Tr.SkillShellCmdRuntime(osHint) +
		input.Tr.SkillShellCmdRepoStatus(repoSummary) +
		input.Tr.SkillShellCmdUserIntent(intent)

	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: systemMsg},
		{Role: provider.RoleUser, Content: userMsg},
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
	content := extractJSONContent(raw)

	// Primary: expect a JSON array.
	var suggestions []CommandSuggestion
	if err := json.Unmarshal([]byte(content), &suggestions); err == nil {
		return normaliseRiskLevels(suggestions), nil
	}

	// Fallback: AI occasionally returns a single JSON object instead of an array.
	var single CommandSuggestion
	if err := json.Unmarshal([]byte(content), &single); err == nil && single.Command != "" {
		return normaliseRiskLevels([]CommandSuggestion{single}), nil
	}

	return nil, fmt.Errorf("invalid JSON response: %s", content)
}

// extractJSONContent strips markdown code fences from AI output.
func extractJSONContent(raw string) string {
	s := strings.TrimSpace(raw)
	// Strip optional ```json or ``` fence
	for _, prefix := range []string{"```json", "```"} {
		if strings.HasPrefix(s, prefix) {
			s = strings.TrimPrefix(s, prefix)
			s = strings.TrimSuffix(strings.TrimSpace(s), "```")
			return strings.TrimSpace(s)
		}
	}
	return s
}

func normaliseRiskLevels(suggestions []CommandSuggestion) []CommandSuggestion {
	for i := range suggestions {
		switch suggestions[i].RiskLevel {
		case "safe", "medium", "dangerous":
			// valid
		default:
			suggestions[i].RiskLevel = "medium"
		}
	}
	return suggestions
}
