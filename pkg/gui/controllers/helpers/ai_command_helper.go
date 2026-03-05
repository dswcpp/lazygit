package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	aiprovider "github.com/dswcpp/lazygit/pkg/ai/provider"
)

// CommandSuggestion represents an AI-generated command suggestion.
type CommandSuggestion struct {
	Command      string `json:"command"`
	Explanation  string `json:"explanation"`
	RiskLevel    string `json:"risk_level"` // "safe", "medium", "dangerous"
	Alternatives string `json:"alternatives,omitempty"`
}

// AICommandHelper provides AI-powered command generation and safety checking.
type AICommandHelper struct {
	aiHelper *AIHelper
	c        *HelperCommon
}

// NewAICommandHelper creates a new AI command helper.
func NewAICommandHelper(c *HelperCommon, aiHelper *AIHelper) *AICommandHelper {
	return &AICommandHelper{
		aiHelper: aiHelper,
		c:        c,
	}
}

// GenerateShellCommand uses AI to convert natural language into shell commands.
// Returns multiple suggestions with risk levels and explanations.
func (self *AICommandHelper) GenerateShellCommand(userIntent string) ([]CommandSuggestion, error) {
	if self.c.AIManager == nil {
		return nil, errors.New("AI 功能未启用")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err := self.c.AIManager.RunSkill(ctx, "shell_cmd", map[string]any{
		"intent": userIntent,
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, errors.New("AI 命令生成已取消")
		}
		return nil, self.aiHelper.HandleAIError(err)
	}

	var suggestions []CommandSuggestion
	if err := json.Unmarshal([]byte(output.Content), &suggestions); err != nil {
		return nil, fmt.Errorf("AI 返回的格式无效: %v\n响应内容: %s", err, output.Content)
	}
	for i := range suggestions {
		if !isValidRiskLevel(suggestions[i].RiskLevel) {
			suggestions[i].RiskLevel = "medium"
		}
	}
	return suggestions, nil
}

// ExplainCommand uses AI to explain what a shell command does.
func (self *AICommandHelper) ExplainCommand(command string) (string, error) {
	if self.c.AIManager == nil {
		return "", errors.New("AI 功能未启用")
	}

	prompt := fmt.Sprintf(`解释这个 shell 命令会做什么，用简洁的中文说明：

命令：%s

请说明：
1. 这个命令的作用（1 行）
2. 会产生什么影响（1-2 行）
3. 是否有风险（如果有，说明风险点，1 行）
4. 建议或注意事项（可选，1 行）

保持回答简洁（总共 3-5 行）。`,
		command,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := self.c.AIManager.Provider().Complete(ctx, []aiprovider.Message{
		{Role: aiprovider.RoleUser, Content: prompt},
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", errors.New("命令解释已取消")
		}
		return "", self.aiHelper.HandleAIError(err)
	}

	return strings.TrimSpace(result.Content), nil
}

// CheckCommandSafety checks if a command is potentially dangerous.
// Returns true if safe, false with reason if dangerous.
func (self *AICommandHelper) CheckCommandSafety(command string) (bool, string) {
	cmdLower := strings.ToLower(command)

	// Check against known dangerous patterns
	for _, pattern := range GetDangerousCommandPatterns() {
		if strings.Contains(cmdLower, pattern.Pattern) {
			return false, pattern.Risk
		}
	}

	return true, ""
}

// GetRiskLevel determines the risk level of a command.
func (self *AICommandHelper) GetRiskLevel(command string) string {
	cmdLower := strings.ToLower(command)

	// Dangerous commands
	dangerousPatterns := []string{
		"reset --hard",
		"clean -fdx",
		"clean -fd",
		"push --force",
		"push -f",
		"reflog expire",
		"rm -rf",
		"branch -D",
		"tag -d",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(cmdLower, pattern) {
			return "dangerous"
		}
	}

	// Medium risk commands
	mediumPatterns := []string{
		"reset",
		"push --force-with-lease",
		"rebase",
		"cherry-pick",
		"clean",
		"stash pop",
		"branch -d",
	}

	for _, pattern := range mediumPatterns {
		if strings.Contains(cmdLower, pattern) {
			return "medium"
		}
	}

	// Safe by default
	return "safe"
}

// DangerousPattern represents a dangerous command pattern and its risk.
type DangerousPattern struct {
	Pattern string
	Risk    string
}

// GetDangerousCommandPatterns returns a list of dangerous command patterns.
func GetDangerousCommandPatterns() []DangerousPattern {
	return []DangerousPattern{
		{
			Pattern: "git reset --hard",
			Risk:    "⚠️ 将丢失所有未提交的更改",
		},
		{
			Pattern: "git clean -fdx",
			Risk:    "⚠️ 将删除所有未跟踪和忽略的文件（包括 .gitignore 中的文件）",
		},
		{
			Pattern: "git clean -fd",
			Risk:    "⚠️ 将删除所有未跟踪的文件和目录",
		},
		{
			Pattern: "git push --force",
			Risk:    "⚠️ 可能覆盖远程分支历史，影响其他协作者",
		},
		{
			Pattern: "git push -f",
			Risk:    "⚠️ 可能覆盖远程分支历史，影响其他协作者",
		},
		{
			Pattern: "git reflog expire",
			Risk:    "⚠️ 将永久删除 reflog 记录，无法恢复",
		},
		{
			Pattern: "rm -rf",
			Risk:    "⚠️ 危险：递归删除文件，可能删除重要数据",
		},
		{
			Pattern: "git branch -D",
			Risk:    "⚠️ 强制删除分支，即使分支未合并",
		},
		{
			Pattern: "git filter-branch",
			Risk:    "⚠️ 重写历史，可能造成协作问题",
		},
		{
			Pattern: "git gc --aggressive --prune=now",
			Risk:    "⚠️ 激进的垃圾回收，可能删除最近的对象",
		},
	}
}

// GetSafterAlternative suggests a safer alternative for dangerous commands.
func (self *AICommandHelper) GetSafterAlternative(command string) string {
	cmdLower := strings.ToLower(command)

	alternatives := map[string]string{
		"git reset --hard": "建议：先使用 'git stash' 保存更改，或使用 'git reset --soft' 保留更改",
		"git clean -fdx":   "建议：先使用 'git clean -fdn' 预览将被删除的文件",
		"git clean -fd":    "建议：先使用 'git clean -fdn' 预览将被删除的文件",
		"git push --force": "建议：使用 'git push --force-with-lease' 进行安全的强制推送",
		"git push -f":      "建议：使用 'git push --force-with-lease' 进行安全的强制推送",
		"git branch -D":    "建议：先检查分支是否已合并，使用 'git branch -d' 进行安全删除",
	}

	for pattern, alternative := range alternatives {
		if strings.Contains(cmdLower, pattern) {
			return alternative
		}
	}

	return ""
}

// isValidRiskLevel checks if a risk level is valid.
func isValidRiskLevel(level string) bool {
	validLevels := []string{"safe", "medium", "dangerous"}
	for _, valid := range validLevels {
		if level == valid {
			return true
		}
	}
	return false
}

// GetRiskLevelIcon returns an icon for the risk level.
func GetRiskLevelIcon(riskLevel string) string {
	icons := map[string]string{
		"safe":      "✅",
		"medium":    "⚠️",
		"dangerous": "❌",
	}

	if icon, ok := icons[riskLevel]; ok {
		return icon
	}
	return "❓"
}

// GetRiskLevelColor returns a color for the risk level (for UI styling).
func GetRiskLevelColor(riskLevel string) string {
	colors := map[string]string{
		"safe":      "green",
		"medium":    "yellow",
		"dangerous": "red",
	}

	if color, ok := colors[riskLevel]; ok {
		return color
	}
	return "white"
}
