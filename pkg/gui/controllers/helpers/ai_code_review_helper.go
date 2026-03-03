package helpers

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

type AICodeReviewHelper struct {
	c             *HelperCommon
	loadingHelper *LoadingHelper
	aiHelper      *AIHelper
}

func NewAICodeReviewHelper(c *HelperCommon, loadingHelper *LoadingHelper, aiHelper *AIHelper) *AICodeReviewHelper {
	return &AICodeReviewHelper{c: c, loadingHelper: loadingHelper, aiHelper: aiHelper}
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
		// Show first-time wizard instead of error
		return self.aiHelper.ShowFirstTimeWizard()
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

	// Spinner frames for progress indicator
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerFrame := 0
	aiView.Title = fmt.Sprintf(" %s %s: %s ", spinner[0], self.c.Tr.AICodeReviewTitle, filePath)

	// Create cancellable context for the AI request
	ctx, cancel := context.WithCancel(context.Background())

	// Store cancel function in the context so it can be called by Esc key handler
	self.c.Contexts().AICodeReview.CancelFunc = cancel

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

		// Start spinner animation in background
		spinnerDone := make(chan struct{})
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-spinnerDone:
					return
				case <-ticker.C:
					spinnerFrame = (spinnerFrame + 1) % len(spinner)
					self.c.OnUIThreadSync(func() error {
						aiView.Title = fmt.Sprintf(" %s %s: %s ", spinner[spinnerFrame], self.c.Tr.AICodeReviewTitle, filePath)
						return nil
					})
				}
			}
		}()

		// Streaming goroutine: runs independently after the overlay closes.
		// All UI writes use OnUIThreadSync (gocui.UpdateAsync) so that events
		// are enqueued directly from this single goroutine in order, avoiding
		// the race condition caused by OnUIThread spawning a new goroutine per
		// chunk which can arrive at the UI event queue out of order.
		go func() {
			defer func() {
				// Stop spinner
				close(spinnerDone)
				// Clear cancel function when stream completes
				self.c.Contexts().AICodeReview.CancelFunc = nil
				// Update title to show completion
				self.c.OnUIThreadSync(func() error {
					aiView.Title = fmt.Sprintf(" %s: %s ", self.c.Tr.AICodeReviewTitle, filePath)
					return nil
				})
			}()

			err := self.c.AI.CompleteStream(ctx, prompt, func(chunk string) {
				signalFirst()
				self.c.OnUIThreadSync(func() error {
					fmt.Fprint(self.c.Views().AICodeReview, chunk)
					return nil
				})
			})

			if err != nil {
				signalFirst()
				self.c.OnUIThread(func() error {
					// Check if the error is due to cancellation
					if errors.Is(err, context.Canceled) {
						self.c.Toast("AI 代码审查已取消")
						return nil
					}
					// Use friendly error handling from AIHelper
					friendlyErr := self.aiHelper.HandleAIError(err)
					self.c.Toast(friendlyErr.Error())
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
	langHint := ""
	if lang != "" {
		langHint = "（" + lang + "）"
	}

	langSection := ""
	if guidelines := languageGuidelines(lang); guidelines != "" {
		langSection = "\n## 语言特定检查要点" + langHint + "\n" + guidelines + "\n"
	}

	return "你是一名资深软件工程师，正在对以下 git diff 进行代码评审。\n\n" +
		"**文件：** " + filePath + "\n\n" +
		"## 核心原则\n" +
		"- **保守评审**：只报告你**确定**存在问题的地方。不确定时，宁可不报，不要猜测。\n" +
		"- **尊重上下文限制**：你只能看到 diff，看不到完整文件。如果某个问题需要全文上下文才能判断（如某个错误是否已在别处处理），请跳过，不要假设。\n" +
		"- **聚焦新增行**：重点审查以 `+` 开头的新增行；`-` 删除行和上下文行仅用于理解意图，不要对其发表意见。\n" +
		"- **拒绝假阳性**：不要把正确的惯用写法当成问题；不要因为代码\"不是你会写的方式\"就标记为问题。\n" +
		langSection +
		"\n## 严重等级（仅在确认存在时使用）\n" +
		"- **CRITICAL**：会导致崩溃、数据损坏、安全漏洞或明确错误逻辑的 bug。\n" +
		"- **MAJOR**：资源泄漏、明确的错误处理缺失（diff 中可见）、API 使用错误。\n" +
		"- **MINOR**：可能出问题的边界情况、可以更健壮但当前仍能工作的代码。\n" +
		"- **NIT**：纯风格问题，只在确实影响可读性时才报告。\n" +
		"\n## 输出格式（用简体中文输出，代码片段保持原语言）\n\n" +
		"### 摘要\n" +
		"一句话说明本次改动的目的，以及整体是否正确。\n\n" +
		"### 问题列表\n" +
		"每个问题使用以下格式，问题之间空一行：\n\n" +
		"**[等级] 类别 — 标题**\n" +
		"代码：`<有问题的代码片段>`\n" +
		"问题：<问题描述及影响>\n" +
		"建议：<具体修正方案或代码>\n\n" +
		"若无问题，直接写：无问题\n\n" +
		"### 结论\n" +
		"无问题时：LGTM，一句话说明可以合入。\n" +
		"有问题时：列出必须修复的 CRITICAL/MAJOR 项；MINOR/NIT 可一句话汇总。\n\n" +
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
