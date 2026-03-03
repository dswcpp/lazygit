package helpers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dswcpp/lazygit/pkg/gui/types"
)

// SuggestionType represents the type of suggestion.
type SuggestionType string

const (
	SuggestionTypeAI         SuggestionType = "ai"
	SuggestionTypeTemplate   SuggestionType = "template"
	SuggestionTypeCompletion SuggestionType = "completion"
	SuggestionTypeHistory    SuggestionType = "history"
)

// EnhancedSuggestion represents a unified suggestion from any source.
type EnhancedSuggestion struct {
	Command     string
	Description string
	Type        SuggestionType
	Icon        string
	Score       float64
	RiskLevel   string // "safe", "medium", "dangerous"
	Metadata    map[string]interface{}
}

// EnhancedShellCommandHelper integrates all shell command enhancement features.
type EnhancedShellCommandHelper struct {
	c *HelperCommon

	// Sub-systems
	templateEngine   *TemplateEngine
	completionEngine *CompletionEngine
	aiCommandHelper  *AICommandHelper
	aiHelper         *AIHelper

	// State
	aiMode          bool
	lastInput       string
	cachedAISuggest []CommandSuggestion
}

// NewEnhancedShellCommandHelper creates a new enhanced shell command helper.
func NewEnhancedShellCommandHelper(c *HelperCommon, aiHelper *AIHelper) *EnhancedShellCommandHelper {
	return &EnhancedShellCommandHelper{
		c:                c,
		templateEngine:   NewTemplateEngine(),
		completionEngine: NewCompletionEngine(c),
		aiCommandHelper:  NewAICommandHelper(c, aiHelper),
		aiHelper:         aiHelper,
		aiMode:           false,
	}
}

// GetSuggestions returns intelligent suggestions for the given input.
// Combines AI, templates, completions, and history.
func (self *EnhancedShellCommandHelper) GetSuggestions(input string) []EnhancedSuggestion {
	var allSuggestions []EnhancedSuggestion

	// If in AI mode, prioritize AI suggestions
	if self.aiMode {
		aiSuggestions := self.getAISuggestions(input)
		allSuggestions = append(allSuggestions, aiSuggestions...)
	}

	// Always include template suggestions (unless AI mode and already has suggestions)
	if !self.aiMode || len(allSuggestions) == 0 {
		templateSuggestions := self.getTemplateSuggestions(input)
		allSuggestions = append(allSuggestions, templateSuggestions...)
	}

	// Completion suggestions
	completionSuggestions := self.getCompletionSuggestions(input)
	allSuggestions = append(allSuggestions, completionSuggestions...)

	// History suggestions
	historySuggestions := self.getHistorySuggestions(input)
	allSuggestions = append(allSuggestions, historySuggestions...)

	// Rank and sort suggestions
	rankedSuggestions := self.rankSuggestions(input, allSuggestions)

	// Limit to top suggestions
	const maxSuggestions = 20
	if len(rankedSuggestions) > maxSuggestions {
		rankedSuggestions = rankedSuggestions[:maxSuggestions]
	}

	return rankedSuggestions
}

// getAISuggestions gets suggestions from AI (cached if same input).
func (self *EnhancedShellCommandHelper) getAISuggestions(input string) []EnhancedSuggestion {
	// Return cached if same input
	if input == self.lastInput && len(self.cachedAISuggest) > 0 {
		return self.convertAISuggestions(self.cachedAISuggest)
	}

	// Generate new AI suggestions (this should be called asynchronously in real implementation)
	// For now, return empty to avoid blocking
	return nil
}

// GenerateAISuggestionsAsync generates AI suggestions asynchronously.
func (self *EnhancedShellCommandHelper) GenerateAISuggestionsAsync(input string, callback func([]EnhancedSuggestion)) {
	go func() {
		suggestions, err := self.aiCommandHelper.GenerateShellCommand(input)
		if err != nil {
			// Log error but don't block UI
			self.c.Log.Errorf("Failed to generate AI suggestions: %v", err)
			callback(nil)
			return
		}

		self.cachedAISuggest = suggestions
		self.lastInput = input
		callback(self.convertAISuggestions(suggestions))
	}()
}

// convertAISuggestions converts AI suggestions to enhanced suggestions.
func (self *EnhancedShellCommandHelper) convertAISuggestions(suggestions []CommandSuggestion) []EnhancedSuggestion {
	var result []EnhancedSuggestion
	for i, sugg := range suggestions {
		result = append(result, EnhancedSuggestion{
			Command:     sugg.Command,
			Description: sugg.Explanation,
			Type:        SuggestionTypeAI,
			Icon:        GetRiskLevelIcon(sugg.RiskLevel),
			Score:       100.0 - float64(i)*10.0, // Higher score for earlier suggestions
			RiskLevel:   sugg.RiskLevel,
			Metadata: map[string]interface{}{
				"alternatives": sugg.Alternatives,
			},
		})
	}
	return result
}

// getTemplateSuggestions gets suggestions from command templates.
func (self *EnhancedShellCommandHelper) getTemplateSuggestions(input string) []EnhancedSuggestion {
	templates := self.templateEngine.Search(input)

	var result []EnhancedSuggestion
	for _, tmpl := range templates {
		// Determine risk level from template command
		riskLevel := self.aiCommandHelper.GetRiskLevel(tmpl.Command)

		result = append(result, EnhancedSuggestion{
			Command:     tmpl.Command,
			Description: tmpl.Description,
			Type:        SuggestionTypeTemplate,
			Icon:        tmpl.Icon,
			Score:       float64(tmpl.Priority) * 5.0,
			RiskLevel:   riskLevel,
			Metadata: map[string]interface{}{
				"category":     tmpl.Category,
				"placeholders": tmpl.Placeholders,
			},
		})
	}

	return result
}

// getCompletionSuggestions gets suggestions from completion engine.
func (self *EnhancedShellCommandHelper) getCompletionSuggestions(input string) []EnhancedSuggestion {
	completions := self.completionEngine.Complete(input)

	var result []EnhancedSuggestion
	for _, comp := range completions {
		// Build full command from completion
		command := self.buildCommandFromCompletion(input, comp.Text)

		result = append(result, EnhancedSuggestion{
			Command:     command,
			Description: comp.Description,
			Type:        SuggestionTypeCompletion,
			Icon:        getIconForCompletionType(comp.Type),
			Score:       float64(comp.Priority) * 4.0,
			RiskLevel:   "safe", // Completions are generally safe
			Metadata: map[string]interface{}{
				"completion_type": comp.Type,
			},
		})
	}

	return result
}

// buildCommandFromCompletion builds a full command from partial input and completion.
func (self *EnhancedShellCommandHelper) buildCommandFromCompletion(input string, completion string) string {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return completion
	}

	// Replace last part with completion
	parts[len(parts)-1] = completion
	return strings.Join(parts, " ")
}

// getHistorySuggestions gets suggestions from command history.
func (self *EnhancedShellCommandHelper) getHistorySuggestions(input string) []EnhancedSuggestion {
	history := self.c.GetAppState().ShellCommandsHistory

	var result []EnhancedSuggestion
	for _, cmd := range history {
		// Filter by input
		if input != "" && !strings.Contains(strings.ToLower(cmd), strings.ToLower(input)) {
			continue
		}

		// Determine risk level
		riskLevel := self.aiCommandHelper.GetRiskLevel(cmd)

		result = append(result, EnhancedSuggestion{
			Command:     cmd,
			Description: "历史命令",
			Type:        SuggestionTypeHistory,
			Icon:        "📜",
			Score:       30.0, // Lower score than templates/completions
			RiskLevel:   riskLevel,
		})
	}

	return result
}

// rankSuggestions ranks and sorts suggestions by relevance.
func (self *EnhancedShellCommandHelper) rankSuggestions(input string, suggestions []EnhancedSuggestion) []EnhancedSuggestion {
	inputLower := strings.ToLower(input)

	for i := range suggestions {
		score := suggestions[i].Score

		// Boost for exact prefix match
		if strings.HasPrefix(strings.ToLower(suggestions[i].Command), inputLower) {
			score += 50.0
		}

		// Boost for contains match
		if strings.Contains(strings.ToLower(suggestions[i].Command), inputLower) {
			score += 20.0
		}

		// Boost for type priority
		score += getTypePriority(suggestions[i].Type)

		// Penalty for dangerous commands (unless explicitly searching for them)
		if suggestions[i].RiskLevel == "dangerous" {
			score -= 10.0
		}

		suggestions[i].Score = score
	}

	// Sort by score (descending)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	return suggestions
}

// getTypePriority returns priority boost for suggestion type.
func getTypePriority(suggestionType SuggestionType) float64 {
	priorities := map[SuggestionType]float64{
		SuggestionTypeAI:         100.0, // AI suggestions highest priority
		SuggestionTypeTemplate:   80.0,
		SuggestionTypeCompletion: 60.0,
		SuggestionTypeHistory:    40.0,
	}
	return priorities[suggestionType]
}

// getIconForCompletionType returns an icon for completion type.
func getIconForCompletionType(completionType string) string {
	icons := map[string]string{
		"command": "⚡",
		"flag":    "🚩",
		"branch":  "🌿",
		"file":    "📄",
		"tag":     "🏷️",
		"remote":  "🌐",
		"commit":  "📝",
		"ref":     "🔗",
	}

	if icon, ok := icons[completionType]; ok {
		return icon
	}
	return "💡"
}

// ToggleAIMode toggles AI mode on/off.
func (self *EnhancedShellCommandHelper) ToggleAIMode() bool {
	self.aiMode = !self.aiMode
	return self.aiMode
}

// IsAIMode returns whether AI mode is active.
func (self *EnhancedShellCommandHelper) IsAIMode() bool {
	return self.aiMode
}

// ExplainCommand explains a command using AI.
func (self *EnhancedShellCommandHelper) ExplainCommand(command string) (string, error) {
	return self.aiCommandHelper.ExplainCommand(command)
}

// CheckCommandSafety checks if a command is safe to execute.
func (self *EnhancedShellCommandHelper) CheckCommandSafety(command string) (bool, string) {
	return self.aiCommandHelper.CheckCommandSafety(command)
}

// GetSafterAlternative suggests a safer alternative for dangerous commands.
func (self *EnhancedShellCommandHelper) GetSafterAlternative(command string) string {
	return self.aiCommandHelper.GetSafterAlternative(command)
}

// FormatSuggestionForDisplay formats a suggestion for display in the UI.
func (self *EnhancedShellCommandHelper) FormatSuggestionForDisplay(suggestion EnhancedSuggestion) string {
	riskIcon := GetRiskLevelIcon(suggestion.RiskLevel)
	return fmt.Sprintf("%s %s %s - %s",
		suggestion.Icon,
		riskIcon,
		suggestion.Command,
		suggestion.Description,
	)
}

// GroupSuggestionsByType groups suggestions by their type for categorized display.
func (self *EnhancedShellCommandHelper) GroupSuggestionsByType(suggestions []EnhancedSuggestion) map[SuggestionType][]EnhancedSuggestion {
	groups := make(map[SuggestionType][]EnhancedSuggestion)

	for _, sugg := range suggestions {
		groups[sugg.Type] = append(groups[sugg.Type], sugg)
	}

	return groups
}

// ConvertToTypesSuggestions converts EnhancedSuggestion to types.Suggestion for UI.
func (self *EnhancedShellCommandHelper) ConvertToTypesSuggestions(suggestions []EnhancedSuggestion) []*types.Suggestion {
	var result []*types.Suggestion

	for _, sugg := range suggestions {
		label := self.FormatSuggestionForDisplay(sugg)
		result = append(result, &types.Suggestion{
			Value: sugg.Command,
			Label: label,
		})
	}

	return result
}
