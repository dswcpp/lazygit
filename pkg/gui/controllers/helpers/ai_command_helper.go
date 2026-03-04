package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"
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
	if self.c.AI == nil {
		return nil, errors.New("AI 功能未启用")
	}

	// Build repository context
	repoContext := self.aiHelper.buildGitContext()

	// Detect OS for shell-specific command format
	osShellHint := func() string {
		switch runtime.GOOS {
		case "windows":
			return "运行环境: Windows + Git Bash\n命令格式: 直接用 && 连接多个命令，禁止使用 cmd /c 或 ^&^& 转义"
		case "darwin":
			return "运行环境: macOS + zsh/bash\n命令格式: 直接用 && 连接多个命令"
		default:
			return "运行环境: Linux + bash\n命令格式: 直接用 && 连接多个命令"
		}
	}()

	// Create detailed prompt for structured JSON output
	prompt := fmt.Sprintf(`你是一个 Git 命令专家。根据用户意图生成精确的 shell 命令。

%s

规则：
1. 输出 JSON 数组，每个元素包含：
   - "command": 完整的可执行命令（必须符合上述运行环境的格式）
   - "explanation": 命令的中文解释（简洁，1-2 句话）
   - "risk_level": 风险等级（必须是 "safe", "medium", "dangerous" 之一）
   - "alternatives": 替代命令（可选，如果有更安全的方案）

2. 优先使用安全的命令
3. 对于危险操作，提供更安全的替代方案
4. 返回 1-3 个命令建议，按推荐程度排序

风险等级定义：
- safe: 不会丢失数据或更改历史（如 git status, git log, git stash）
- medium: 可能影响工作区但可恢复（如 git reset --soft, git stash pop）
- dangerous: 可能永久丢失数据（如 git reset --hard, git clean -fdx, git push --force）

仓库状态：
%s

用户意图：%s

示例输出：
[
  {
    "command": "git commit -m \"feat: add feature\"",
    "explanation": "提交当前暂存的更改",
    "risk_level": "safe"
  },
  {
    "command": "git reset --soft HEAD~1",
    "explanation": "撤销最后一次提交但保留更改",
    "risk_level": "medium",
    "alternatives": "git commit --amend (如果只是想修改提交消息)"
  }
]

输出（仅 JSON，不要其他内容）：`,
		osShellHint,
		repoContext,
		userIntent,
	)

	// Call AI with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := self.c.AI.Complete(ctx, prompt)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, errors.New("AI 命令生成已取消")
		}
		return nil, self.aiHelper.HandleAIError(err)
	}

	// Clean response (remove markdown code blocks if present)
	content := strings.TrimSpace(result.Content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// Parse JSON
	var suggestions []CommandSuggestion
	if err := json.Unmarshal([]byte(content), &suggestions); err != nil {
		return nil, fmt.Errorf("AI 返回的格式无效: %v\n响应内容: %s", err, content)
	}

	// Validate risk levels
	for i := range suggestions {
		if !isValidRiskLevel(suggestions[i].RiskLevel) {
			suggestions[i].RiskLevel = "medium" // Default to medium if invalid
		}
	}

	return suggestions, nil
}

// ExplainCommand uses AI to explain what a shell command does.
func (self *AICommandHelper) ExplainCommand(command string) (string, error) {
	if self.c.AI == nil {
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

	result, err := self.c.AI.Complete(ctx, prompt)
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
