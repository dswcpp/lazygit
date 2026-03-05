package skills

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"

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

	osHint := shellHint()
	repoSummary := input.RepoCtx.CompactString()

	prompt := fmt.Sprintf(
		"你是一个 Git 命令专家。根据用户意图生成精确的 shell 命令。\n\n"+
			"运行环境: %s\n\n"+
			"仓库状态:\n%s\n\n"+
			"用户意图: %s\n\n"+
			"输出 JSON 数组，每个元素包含:\n"+
			"- command: 完整可执行命令\n"+
			"- explanation: 中文解释（1-2 句）\n"+
			"- risk_level: \"safe\" | \"medium\" | \"dangerous\"\n"+
			"- alternatives: 替代命令（可选）\n\n"+
			"返回 1-3 个建议，按推荐度排序。只输出 JSON，不要其他内容。",
		osHint, repoSummary, intent,
	)

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

func shellHint() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows + Git Bash，用 && 连接命令"
	case "darwin":
		return "macOS + zsh/bash，用 && 连接命令"
	default:
		return "Linux + bash，用 && 连接命令"
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
