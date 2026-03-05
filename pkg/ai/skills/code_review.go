package skills

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

// CodeReviewSkill reviews a git diff and returns structured feedback.
// Extra keys:
//   - "diff"      (string, required) the diff to review
//   - "file_path" (string, optional) file path for language detection
type CodeReviewSkill struct{}

func NewCodeReviewSkill() Skill { return &CodeReviewSkill{} }
func (s *CodeReviewSkill) Name() string { return "code_review" }

func (s *CodeReviewSkill) Execute(ctx context.Context, p provider.Provider, input Input) (Output, error) {
	diff := extraStr(input.Extra, "diff")
	if strings.TrimSpace(diff) == "" {
		return Output{}, errors.New("diff is empty")
	}
	filePath := extraStr(input.Extra, "file_path")
	lang := detectLang(filePath)

	messages := []provider.Message{
		{Role: provider.RoleUser, Content: buildReviewPrompt(filePath, lang, diff)},
	}

	result, err := p.Complete(ctx, messages)
	if err != nil {
		return Output{}, err
	}
	content := strings.TrimSpace(result.Content)
	if content == "" {
		return Output{}, errors.New("AI returned empty review")
	}
	return Output{Content: content}, nil
}

// buildReviewPrompt mirrors the existing buildCodeReviewPrompt in ai_code_review_helper.go.
func buildReviewPrompt(filePath, lang, diff string) string {
	langHint := ""
	if lang != "" {
		langHint = fmt.Sprintf(" (%s)", lang)
	}
	langSection := ""
	if guidelines := langGuidelines(lang); guidelines != "" {
		langSection = fmt.Sprintf("\n## Language-Specific Checks%s\n%s\n", langHint, guidelines)
	}

	return fmt.Sprintf(
		"You are a senior software engineer conducting a code review on the following git diff.\n\n"+
			"**File:** %s\n\n"+
			"## Core Principles\n"+
			"- **Conservative review**: Only report issues you are **certain** exist.\n"+
			"- **Focus on new lines**: Focus on lines starting with `+`.\n"+
			"- **Reject false positives**: Do not flag correct idiomatic code.\n"+
			"%s"+
			"\n## Severity Levels\n"+
			"- **CRITICAL**: Crashes, data corruption, security vulnerabilities.\n"+
			"- **MAJOR**: Resource leaks, missing error handling, API usage errors.\n"+
			"- **MINOR**: Edge cases, robustness improvements.\n"+
			"- **NIT**: Style issues that affect readability.\n"+
			"\n## Output Format (Simplified Chinese, code snippets in original language)\n\n"+
			"### Summary\nOne sentence explaining the purpose and overall correctness.\n\n"+
			"### Issue List\n"+
			"**[Level] Category — Title**\n"+
			"Code: `<snippet>`\nIssue: <description>\nSuggestion: <fix>\n\n"+
			"If no issues: No issues\n\n"+
			"### Conclusion\nLGTM or list must-fix items.\n\n"+
			"---\n\n## Diff\n```diff\n%s\n```",
		filePath, langSection, diff,
	)
}

func detectLang(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	langs := map[string]string{
		".go": "Go", ".ts": "TypeScript", ".tsx": "TypeScript/React",
		".js": "JavaScript", ".jsx": "JavaScript/React", ".py": "Python",
		".rs": "Rust", ".java": "Java", ".c": "C", ".h": "C",
		".cpp": "C++", ".cc": "C++", ".hpp": "C++",
		".rb": "Ruby", ".php": "PHP", ".swift": "Swift",
		".kt": "Kotlin", ".cs": "C#", ".sh": "Shell", ".bash": "Shell",
		".yaml": "YAML", ".yml": "YAML", ".json": "JSON", ".sql": "SQL",
	}
	return langs[ext]
}

func langGuidelines(lang string) string {
	switch lang {
	case "Go":
		return "- Check for proper error handling (errors must not be ignored)\n" +
			"- Check for goroutine leaks and proper channel usage\n" +
			"- Check for race conditions in concurrent code\n"
	case "TypeScript", "TypeScript/React":
		return "- Check for proper TypeScript type safety (avoid `any`)\n" +
			"- Check for proper async/await error handling\n"
	case "Python":
		return "- Check for proper exception handling\n" +
			"- Check for resource management (context managers)\n"
	case "Rust":
		return "- Check for proper Result/Option handling\n" +
			"- Check for unsafe code blocks\n"
	default:
		return ""
	}
}
