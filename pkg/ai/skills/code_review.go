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
//   - "focus"     (string, optional) review focus: "security" | "performance" | "correctness"
type CodeReviewSkill struct{}

func NewCodeReviewSkill() Skill { return &CodeReviewSkill{} }
func (s *CodeReviewSkill) Name() string { return "code_review" }

func (s *CodeReviewSkill) Execute(ctx context.Context, p provider.Provider, input Input) (Output, error) {
	diff := extraStr(input.Extra, "diff")
	if strings.TrimSpace(diff) == "" {
		return Output{}, errors.New("diff is empty")
	}
	filePath := extraStr(input.Extra, "file_path")
	focus := extraStr(input.Extra, "focus")
	lang := detectLang(filePath)

	var messages []provider.Message
	if input.Tr != nil {
		messages = append(messages, provider.Message{
			Role:    provider.RoleSystem,
			Content: input.Tr.SkillCodeReviewSystemPrompt(),
		})
	}
	messages = append(messages, provider.Message{
		Role:    provider.RoleUser,
		Content: buildReviewPrompt(filePath, lang, focus, diff),
	})

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

// buildReviewPrompt builds the user-role review request.
// The role description ("You are a senior engineer…") lives in the system message;
// this function only provides context and instructions specific to this review.
func buildReviewPrompt(filePath, lang, focus, diff string) string {
	var sb strings.Builder

	if filePath != "" {
		sb.WriteString(fmt.Sprintf("**File:** `%s`", filePath))
		if lang != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", lang))
		}
		sb.WriteString("\n\n")
	}

	sb.WriteString("## Review Principles\n")
	sb.WriteString("- **Focus on added lines** (starting with `+`); context lines are for reference only.\n")
	sb.WriteString("- **Be conservative**: only report issues you are certain exist.\n")
	sb.WriteString("- **No false positives**: do not flag correct idiomatic code.\n\n")

	if focusSection := buildFocusSection(focus); focusSection != "" {
		sb.WriteString(focusSection)
	}

	if guidelines := langGuidelines(lang); guidelines != "" {
		hint := lang
		sb.WriteString(fmt.Sprintf("## %s-Specific Checks\n%s\n", hint, guidelines))
	}

	sb.WriteString("## Severity Levels\n")
	sb.WriteString("- **CRITICAL**: Crashes, data corruption, security vulnerabilities.\n")
	sb.WriteString("- **MAJOR**: Resource leaks, missing error handling, API misuse.\n")
	sb.WriteString("- **MINOR**: Edge cases, robustness improvements.\n")
	sb.WriteString("- **NIT**: Style issues affecting readability.\n\n")

	sb.WriteString("## Output Format\n")
	sb.WriteString("Use Simplified Chinese. Code snippets stay in the original language.\n\n")
	sb.WriteString("### Summary\nOne sentence: purpose of the change and overall assessment.\n\n")
	sb.WriteString("### Issues\n")
	sb.WriteString("For each issue:\n```\n**[LEVEL] Category — Title**\nCode: `<snippet>`\nProblem: <description>\nFix: <suggestion>\n```\n")
	sb.WriteString("If no issues found, write: 无问题\n\n")
	sb.WriteString("### Conclusion\nLGTM / or list must-fix items.\n\n")

	sb.WriteString("---\n\n## Diff\n```diff\n")
	sb.WriteString(diff)
	sb.WriteString("\n```")

	return sb.String()
}

func buildFocusSection(focus string) string {
	switch strings.ToLower(focus) {
	case "security":
		return "## Focus: Security\n" +
			"Pay special attention to: injection attacks, authentication bypasses, " +
			"insecure deserialization, sensitive data exposure, improper access control.\n\n"
	case "performance":
		return "## Focus: Performance\n" +
			"Pay special attention to: unnecessary allocations, N+1 queries, " +
			"blocking calls in hot paths, inefficient algorithms, missing caching.\n\n"
	case "correctness":
		return "## Focus: Correctness\n" +
			"Pay special attention to: logic errors, off-by-one errors, " +
			"unhandled edge cases, incorrect error propagation, race conditions.\n\n"
	default:
		return ""
	}
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
		return "- Errors must be handled or explicitly ignored with `_` and a comment\n" +
			"- Goroutines must have a clear owner and exit path; check for leaks\n" +
			"- Shared state accessed from multiple goroutines needs synchronisation\n" +
			"- Prefer `errors.Is`/`errors.As` over string comparison for error checks\n"
	case "TypeScript", "TypeScript/React":
		return "- Avoid `any`; use precise types or `unknown` with type guards\n" +
			"- Every `async` function needs error handling (`try/catch` or `.catch`)\n" +
			"- Avoid non-null assertions (`!`) without a clear justification\n"
	case "JavaScript", "JavaScript/React":
		return "- Prefer `const`/`let` over `var`\n" +
			"- Every `async` function needs error handling\n" +
			"- Watch for implicit type coercions (use `===` not `==`)\n"
	case "Python":
		return "- Use context managers (`with`) for resources (files, locks, connections)\n" +
			"- Catch specific exceptions, not bare `except:`\n" +
			"- Mutable default arguments (`def f(x=[])`) are a common bug\n"
	case "Rust":
		return "- `Result`/`Option` must be handled; avoid unnecessary `unwrap()`\n" +
			"- Justify every `unsafe` block with a safety comment\n" +
			"- Watch for unintentional clones in performance-sensitive paths\n"
	case "Java":
		return "- Resources (`Closeable`) must be closed in `finally` or via try-with-resources\n" +
			"- Check for `NullPointerException` risks; prefer `Optional` in new code\n" +
			"- Checked exceptions should not be silently swallowed\n"
	case "C", "C++":
		return "- Check for buffer overflows and out-of-bounds array access\n" +
			"- Every heap allocation needs a corresponding free/delete; check for leaks\n" +
			"- Pointer arithmetic must stay within valid bounds\n" +
			"- (C++) Prefer RAII and smart pointers over raw `new`/`delete`\n"
	case "Shell":
		return "- Quote all variable expansions to prevent word-splitting (`\"$var\"`)\n" +
			"- Use `set -euo pipefail` at the top of scripts\n" +
			"- External input must be sanitised before use in commands\n"
	case "SQL":
		return "- Parameterise all user-supplied values; never interpolate into query strings\n" +
			"- Verify indexes exist for columns used in `WHERE`/`JOIN` of large tables\n" +
			"- Check that transactions are properly committed or rolled back\n"
	default:
		return ""
	}
}
