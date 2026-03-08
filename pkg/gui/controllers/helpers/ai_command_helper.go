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
		return nil, errors.New(self.c.Tr.AICommandNotEnabled)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err := self.c.AIManager.RunSkill(ctx, "shell_cmd", map[string]any{
		"intent": userIntent,
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, errors.New(self.c.Tr.AICommandGenerationCancelled)
		}
		return nil, self.aiHelper.HandleAIError(err)
	}

	var suggestions []CommandSuggestion
	if err := json.Unmarshal([]byte(output.Content), &suggestions); err != nil {
		return nil, fmt.Errorf(self.c.Tr.AICommandInvalidFormat, err, output.Content)
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
		return "", errors.New(self.c.Tr.AICommandNotEnabled)
	}

	prompt := fmt.Sprintf(self.c.Tr.AICommandExplainPrompt, command)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := self.c.AIManager.Provider().Complete(ctx, []aiprovider.Message{
		{Role: aiprovider.RoleUser, Content: prompt},
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", errors.New(self.c.Tr.AICommandExplainCancelled)
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
	for _, pattern := range self.GetDangerousCommandPatterns() {
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
func (self *AICommandHelper) GetDangerousCommandPatterns() []DangerousPattern {
	return []DangerousPattern{
		{
			Pattern: "git reset --hard",
			Risk:    self.c.Tr.AICommandRiskHardReset,
		},
		{
			Pattern: "git clean -fdx",
			Risk:    self.c.Tr.AICommandRiskCleanFdx,
		},
		{
			Pattern: "git clean -fd",
			Risk:    self.c.Tr.AICommandRiskCleanFd,
		},
		{
			Pattern: "git push --force",
			Risk:    self.c.Tr.AICommandRiskForcePush1,
		},
		{
			Pattern: "git push -f",
			Risk:    self.c.Tr.AICommandRiskForcePush2,
		},
		{
			Pattern: "git reflog expire",
			Risk:    self.c.Tr.AICommandRiskReflogExpire,
		},
		{
			Pattern: "rm -rf",
			Risk:    self.c.Tr.AICommandRiskRmRf,
		},
		{
			Pattern: "git branch -D",
			Risk:    self.c.Tr.AICommandRiskBranchD,
		},
		{
			Pattern: "git filter-branch",
			Risk:    self.c.Tr.AICommandRiskRebaseI,
		},
		{
			Pattern: "git gc --aggressive --prune=now",
			Risk:    self.c.Tr.AICommandRiskGcAggressive,
		},
	}
}

// GetSafterAlternative suggests a safer alternative for dangerous commands.
func (self *AICommandHelper) GetSafterAlternative(command string) string {
	cmdLower := strings.ToLower(command)

	alternatives := map[string]string{
		"git reset --hard": self.c.Tr.AICommandSuggestionHardReset,
		"git clean -fdx":   self.c.Tr.AICommandSuggestionCleanFdx,
		"git clean -fd":    self.c.Tr.AICommandSuggestionCleanFd,
		"git push --force": self.c.Tr.AICommandSuggestionForcePush1,
		"git push -f":      self.c.Tr.AICommandSuggestionForcePush2,
		"git branch -D":    self.c.Tr.AICommandSuggestionBranchD,
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
