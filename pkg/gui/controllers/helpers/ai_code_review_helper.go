package helpers

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

type AICodeReviewHelper struct {
	c             *HelperCommon
	loadingHelper *LoadingHelper
}

func NewAICodeReviewHelper(c *HelperCommon, loadingHelper *LoadingHelper) *AICodeReviewHelper {
	return &AICodeReviewHelper{c: c, loadingHelper: loadingHelper}
}

// ReviewDiff asks the user to confirm, then streams an AI code review for the
// given diff into the command log (Extras) panel.
//
// Flow:
//  1. Confirmation dialog: "Review file X?"
//  2. User confirms → centered loading overlay: "AI reviewing, please wait..."
//  3. First SSE chunk arrives → overlay closes; Extras panel header + content stream in.
//  4. Error before first chunk → overlay closes; error toast is shown.
func (self *AICodeReviewHelper) ReviewDiff(filePath string, diff string) error {
	if self.c.AI == nil {
		return errors.New(self.c.Tr.AINotEnabled)
	}

	if diff == "" {
		return errors.New(self.c.Tr.AICodeReviewNoDiff)
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AICodeReviewConfirmTitle,
		Prompt: fmt.Sprintf(self.c.Tr.AICodeReviewConfirmPrompt, filePath),
		HandleConfirm: func() error {
			return self.startReview(filePath, diff)
		},
	})
	return nil
}

// startReview shows the AI code review popup and launches the streaming review.
// Must be called from the UI thread (inside a Confirm HandleConfirm callback).
func (self *AICodeReviewHelper) startReview(filePath, diff string) error {
	lang := detectLanguage(filePath)
	prompt := buildCodeReviewPrompt(filePath, lang, diff)

	// Prepare the popup view before pushing the context.
	aiView := self.c.Views().AICodeReview
	aiView.Clear()
	aiView.Autoscroll = true
	aiView.Title = fmt.Sprintf(" %s: %s ", self.c.Tr.AICodeReviewTitle, filePath)

	// Push the AI code review context to show the floating popup.
	self.c.Context().Push(self.c.Contexts().AICodeReview, types.OnFocusOpts{})

	// WithCenteredLoadingStatus runs the callback on a worker goroutine and
	// hides the overlay when the callback returns.
	self.loadingHelper.WithCenteredLoadingStatus(self.c.Tr.AICodeReviewStatus, func(_ gocui.Task) error {
		// firstChunk is closed exactly once: when the first response chunk
		// arrives (or when the request errors out). Closing it causes the
		// loading overlay to disappear.
		firstChunk := make(chan struct{})
		var once sync.Once
		signalFirst := func() { once.Do(func() { close(firstChunk) }) }

		// Streaming goroutine: runs independently after the overlay closes.
		// All UI writes use OnUIThreadSync (gocui.UpdateAsync) so that events
		// are enqueued directly from this single goroutine in order, avoiding
		// the race condition caused by OnUIThread spawning a new goroutine per
		// chunk which can arrive at the UI event queue out of order.
		go func() {
			err := self.c.AI.CompleteStream(context.Background(), prompt, func(chunk string) {
				signalFirst()
				self.c.OnUIThreadSync(func() error {
					fmt.Fprint(self.c.Views().AICodeReview, chunk)
					return nil
				})
			})

			if err != nil {
				signalFirst()
				self.c.OnUIThread(func() error {
					self.c.Toast("AI code review failed: " + err.Error())
					return nil
				})
			}
		}()

		// Block here until the first chunk arrives → overlay hides.
		<-firstChunk
		return nil
	})

	return nil
}

// buildCodeReviewPrompt constructs a structured, language-aware code review prompt.
func buildCodeReviewPrompt(filePath, lang, diff string) string {
	roleDesc := "You are a senior software engineer"
	if lang != "" {
		roleDesc = fmt.Sprintf("You are a senior %s engineer", lang)
	}

	langSection := ""
	if guidelines := languageGuidelines(lang); guidelines != "" {
		langSection = fmt.Sprintf("\n## Language-Specific Checks (%s)\n%s\n", lang, guidelines)
	}

	return roleDesc + " performing a thorough code review of a git diff.\n\n" +
		"**File:** " + filePath + "\n\n" +
		"## Output Language\n" +
		"You MUST write the entire review in **Simplified Chinese (简体中文)**.\n" +
		"All section headings, finding descriptions, and the verdict must be in Chinese.\n" +
		"Code snippets remain in the original programming language.\n\n" +
		"## Review Rules\n" +
		"1. Focus on **added lines** (lines starting with '+'); use unchanged context lines only to understand intent.\n" +
		"2. Do NOT comment on removed lines (starting with '-') unless a deletion creates a new problem.\n" +
		"3. Be specific: quote the relevant code snippet, never say \"line N\".\n" +
		"4. Provide a corrected snippet whenever you suggest a change.\n" +
		"5. Skip purely cosmetic nits unless they meaningfully affect readability.\n" +
		langSection +
		"\n## Severity Levels\n" +
		"- CRITICAL : Bug, crash, security vulnerability, data corruption, or incorrect logic that will cause failures.\n" +
		"- MAJOR    : Significant performance issue, resource leak, incorrect API/contract usage, missing error handling.\n" +
		"- MINOR    : Edge case not handled, suboptimal but working code, missing validation that could matter.\n" +
		"- NIT      : Style, naming, micro-optimisation, optional improvement. Only include if genuinely useful.\n" +
		"\n## Required Output Format\n" +
		"\n### 摘要\n" +
		"用一句话描述本次改动的目的，以及从整体看是否正确。\n" +
		"\n### 发现的问题\n" +
		"对每个问题，严格使用以下格式输出：\n\n" +
		"[严重等级] 类别 - 简短标题\n" +
		"代码：   <有问题的代码片段，尽量单行>\n" +
		"问题：   <说明存在什么问题以及影响>\n" +
		"修复建议：<具体的修正代码或可操作的改进说明>\n\n" +
		"可选类别：正确性、安全性、性能、错误处理、并发、资源管理、API 使用、可读性、可维护性、测试覆盖。\n\n" +
		"若没有发现问题，写：「无问题」\n" +
		"\n### 评审结论\n" +
		"- 无问题：LGTM — 用一句话确认此改动可以合入。\n" +
		"- 有问题：逐条列出 CRITICAL/MAJOR 必须修复项，并汇总 MINOR/NIT 可选优化项。\n\n" +
		"---\n\n" +
		"## Diff\n" +
		"```diff\n" + diff + "\n```"
}

// detectLanguage infers a human-readable language name from the file extension.
func detectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return "Go"
	case ".ts":
		return "TypeScript"
	case ".tsx":
		return "TypeScript/React"
	case ".js":
		return "JavaScript"
	case ".jsx":
		return "JavaScript/React"
	case ".py":
		return "Python"
	case ".rs":
		return "Rust"
	case ".java":
		return "Java"
	case ".c", ".h":
		return "C"
	case ".cpp", ".cc", ".cxx", ".hpp":
		return "C++"
	case ".rb":
		return "Ruby"
	case ".php":
		return "PHP"
	case ".swift":
		return "Swift"
	case ".kt", ".kts":
		return "Kotlin"
	case ".cs":
		return "C#"
	case ".sh", ".bash":
		return "Shell"
	case ".yaml", ".yml":
		return "YAML"
	case ".json":
		return "JSON"
	case ".sql":
		return "SQL"
	default:
		return ""
	}
}

// languageGuidelines returns a short checklist of common pitfalls for the given language.
// Returns empty string for unknown languages.
func languageGuidelines(lang string) string {
	switch lang {
	case "Go":
		return `- Every error return must be checked; unused errors are bugs.
- Goroutine leaks: ensure goroutines started here are always terminated.
- Context propagation: long-running calls should accept and respect context.Context.
- Defer correctness: deferred calls run in LIFO order; watch for deferred mutations in loops.
- Interface bloat: prefer small, focused interfaces (io.Reader, io.Writer pattern).
- Exported identifiers must have doc comments.`

	case "TypeScript", "TypeScript/React":
		return `- Avoid 'any'; use proper types or generics instead.
- Check for null/undefined: prefer optional chaining (?.) and nullish coalescing (??) over loose checks.
- Async/await: every Promise must be awaited or its rejection handled.
- Side effects in useEffect (React): verify dependency arrays are complete and correct.
- Never mutate state or props directly; always return new objects/arrays.
- Sensitive data must not be logged or exposed to the client.`

	case "JavaScript", "JavaScript/React":
		return `- Unhandled promise rejections: every .then() needs a .catch() or use async/await with try/catch.
- Avoid var; use const by default, let only when reassignment is needed.
- Strict equality: use === and !== instead of == and !=.
- Side effects in useEffect (React): verify dependency arrays are complete and correct.
- Never mutate state or props directly.`

	case "Python":
		return `- Mutable default arguments (def f(x=[])) cause shared state bugs; use None and assign inside.
- Broad exception clauses (except Exception or bare except) hide bugs; catch specific types.
- Resource management: use 'with' statements for files, connections, and locks.
- Type hints: new functions should include parameter and return type annotations.
- Avoid wildcard imports (from x import *); they pollute the namespace.`

	case "Rust":
		return `- Unnecessary clones or copies may indicate a design issue with ownership.
- unwrap()/expect() in production paths: replace with proper error propagation (?).
- Lifetimes: ensure references do not outlive the data they point to.
- Unsafe blocks must be justified with a safety comment explaining the invariants.
- Check for integer overflow in arithmetic; use checked_*/saturating_* in critical paths.`

	case "Java":
		return `- NullPointerException risk: use Optional<T> or @NonNull/@Nullable annotations.
- Resources (streams, connections) must be closed; prefer try-with-resources.
- equals()/hashCode() must be overridden together and consistently.
- Thread safety: shared mutable state needs synchronisation; prefer immutable objects.
- Checked exceptions: do not silently swallow them with an empty catch block.`

	case "C", "C++":
		return `- Memory management: every malloc/new must have a matching free/delete; prefer RAII.
- Buffer bounds: array accesses must be validated; use std::array or std::vector (C++).
- Integer overflow and sign-extension errors in arithmetic.
- Uninitialized variables; always initialise before use.
- Thread safety: data races on shared state; use mutexes or atomics.`

	case "Shell":
		return `- Quote all variable expansions ("$var") to prevent word splitting and globbing.
- Check exit codes: use 'set -e' or explicit checks after each critical command.
- Avoid parsing ls output; use globs or find instead.
- Command injection risk: never pass unsanitised user input to eval or shell expansion.
- Use [[ ]] instead of [ ] for conditionals in bash.`

	case "SQL":
		return `- Parameterised queries only; string concatenation with user input is SQL injection.
- Missing index on columns used in WHERE/JOIN predicates can cause full table scans.
- Transactions: multi-step mutations should be wrapped in a transaction.
- NULL semantics: comparisons with NULL require IS NULL / IS NOT NULL, not = NULL.`

	default:
		return ""
	}
}
