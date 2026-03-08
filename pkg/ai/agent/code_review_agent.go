package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

const (
	// Diff大小限制
	MaxDiffLines = 1000
	MaxDiffBytes = 100 * 1024 // 100KB
)

// CodeReviewAgent 代码评审Agent（支持交互式追问和检查点）
// 特性：
// 1. 基础评审 - 流式输出
// 2. 交互式追问 - Ask方法
// 3. 检查点支持 - 中断恢复
// 4. 批量评审 - ConversationID
type CodeReviewAgent struct {
	provider     provider.Provider
	tr           *aii18n.Translator
	state        CodeReviewState
	mu           sync.Mutex
	checkpointer CodeReviewCheckpointer
	threadID     string
}

// NewCodeReviewAgent 创建代码评审Agent
func NewCodeReviewAgent(p provider.Provider, tr *aii18n.Translator) *CodeReviewAgent {
	return &CodeReviewAgent{
		provider: p,
		tr:       tr,
	}
}

// SetCheckpointer 设置检查点器（支持中断恢复）
func (a *CodeReviewAgent) SetCheckpointer(c CodeReviewCheckpointer, threadID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.checkpointer = c
	a.threadID = threadID
	// 尝试恢复状态
	if saved, ok := c.Load(threadID); ok {
		if saved.ResumeFrom != "" {
			a.state = saved
		}
	}
}

// ReviewWithCallback 支持流式回调的评审方法（线程安全）
func (a *CodeReviewAgent) ReviewWithCallback(
	ctx context.Context,
	filePath string,
	diff string,
	focus string,
	onChunk func(string),
) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 验证输入
	if err := a.validateDiff(diff); err != nil {
		return err
	}

	// 初始化状态
	a.state = CodeReviewState{
		Phase:     PhaseReviewInit,
		FilePath:  filePath,
		Diff:      diff,
		Language:  detectLanguage(filePath),
		Focus:     focus,
		StartTime: time.Now(),
	}

	// 执行评审
	newState, err := a.executeReview(ctx, a.state, onChunk)
	if err != nil {
		a.state = newState
		return err
	}

	a.state = newState

	// 保存检查点（如果配置了）
	a.saveCheckpoint()

	return nil
}

// Ask 追问（交互式评审）
func (a *CodeReviewAgent) Ask(
	ctx context.Context,
	question string,
	onChunk func(string),
) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 检查状态
	if a.state.Phase != PhaseReviewDone && a.state.Phase != PhaseReviewWaiting {
		return fmt.Errorf("cannot ask question in phase: %s", a.state.Phase)
	}

	// 更新状态
	a.state = a.state.WithUserQuestion(question).WithPhase(PhaseReviewInteractive)

	// 追加用户问题到消息历史
	a.state = a.state.AppendMessage(provider.Message{
		Role:    provider.RoleUser,
		Content: question,
	})

	// 流式调用
	var buffer strings.Builder
	err := a.provider.CompleteStream(ctx, a.state.Messages, func(chunk string) {
		buffer.WriteString(chunk)
		if onChunk != nil {
			onChunk(chunk)
		}
	})

	if err != nil {
		a.state = a.state.WithError(err.Error())
		return err
	}

	// 追加AI回复
	answer := buffer.String()
	a.state = a.state.AppendMessage(provider.Message{
		Role:    provider.RoleAssistant,
		Content: answer,
	})

	// 回到等待状态
	a.state = a.state.WithUserQuestion("").WithPhase(PhaseReviewWaiting)

	// 保存检查点
	a.saveCheckpoint()

	return nil
}

// GetState 获取当前状态（线程安全）
func (a *CodeReviewAgent) GetState() CodeReviewState {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state
}

// Phase 返回当前阶段
func (a *CodeReviewAgent) Phase() ReviewPhase {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state.Phase
}

// CanAsk 是否可以追问
func (a *CodeReviewAgent) CanAsk() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state.Phase == PhaseReviewDone || a.state.Phase == PhaseReviewWaiting
}

// validateDiff 验证diff大小
func (a *CodeReviewAgent) validateDiff(diff string) error {
	if strings.TrimSpace(diff) == "" {
		return fmt.Errorf("diff is empty")
	}

	lines := strings.Count(diff, "\n")
	if lines > MaxDiffLines {
		return fmt.Errorf("diff too large: %d lines (max %d)", lines, MaxDiffLines)
	}

	if len(diff) > MaxDiffBytes {
		return fmt.Errorf("diff too large: %d bytes (max %d)", len(diff), MaxDiffBytes)
	}

	return nil
}

// saveCheckpoint 保存检查点
func (a *CodeReviewAgent) saveCheckpoint() {
	if a.checkpointer != nil {
		_ = a.checkpointer.Save(a.threadID, a.state)
	}
}

// clearCheckpoint 清除检查点
func (a *CodeReviewAgent) clearCheckpoint() {
	if a.checkpointer != nil {
		a.checkpointer.Clear(a.threadID)
	}
}

// executeReview 执行评审（纯函数）
func (a *CodeReviewAgent) executeReview(
	ctx context.Context,
	state CodeReviewState,
	onChunk func(string),
) (CodeReviewState, error) {
	state = state.WithPhase(PhaseReviewing)

	prompt := a.buildReviewPrompt(state.FilePath, state.Language, state.Focus, state.Diff)
	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: a.tr.SkillCodeReviewSystemPrompt()},
		{Role: provider.RoleUser, Content: prompt},
	}
	state = state.WithMessages(messages)

	var buffer strings.Builder
	err := a.provider.CompleteStream(ctx, messages, func(chunk string) {
		buffer.WriteString(chunk)
		if onChunk != nil {
			onChunk(chunk)
		}
	})

	if err != nil {
		return state.WithError(err.Error()), err
	}

	result := buffer.String()
	// 完成后进入等待状态（支持追问）
	state = state.WithResult(result).WithPhase(PhaseReviewWaiting)

	return state, nil
}

// buildReviewPrompt 构建评审prompt
func (a *CodeReviewAgent) buildReviewPrompt(filePath, lang, focus, diff string) string {
	langHint := ""
	if lang != "" {
		langHint = fmt.Sprintf(" (%s)", lang)
	}

	focusSection := ""
	if focus != "" {
		focusSection = buildFocusSection(focus)
	}

	langSection := ""
	if guidelines := languageGuidelines(lang); guidelines != "" {
		langSection = fmt.Sprintf("\n## Language-Specific Checks%s\n%s\n", langHint, guidelines)
	}

	return fmt.Sprintf(
		"%s"+
			"**File:** %s\n\n"+
			"## Core Principles\n"+
			"- **Conservative review**: Only report issues you are **certain** exist. When uncertain, prefer not to report rather than guess.\n"+
			"- **Respect context limitations**: You can only see the diff, not the complete file. If an issue requires full file context to judge, skip it.\n"+
			"- **Focus on new lines**: Focus on reviewing new lines starting with `+`; `-` deleted lines and context lines are only for understanding intent.\n"+
			"- **Reject false positives**: Do not flag correct idiomatic code as issues.\n"+
			"%s"+
			"%s"+
			"\n## Severity Levels (only use when confirmed)\n"+
			"- **CRITICAL**: Bugs that will cause crashes, data corruption, security vulnerabilities, or clear logic errors.\n"+
			"- **MAJOR**: Resource leaks, clear missing error handling (visible in diff), API usage errors.\n"+
			"- **MINOR**: Edge cases that might cause problems, code that could be more robust but currently works.\n"+
			"- **NIT**: Pure style issues, only report when it truly affects readability.\n\n"+
			"## Output Format (output in Simplified Chinese, keep code snippets in original language)\n\n"+
			"### Summary\n"+
			"One sentence explaining the purpose of this change and whether it is overall correct.\n\n"+
			"### Issue List\n"+
			"Use the following format for each issue, with blank lines between issues:\n\n"+
			"**[Level] Category — Title**\n"+
			"Code: `<problematic code snippet>`\n"+
			"Issue: <issue description and impact>\n"+
			"Suggestion: <specific fix or code>\n\n"+
			"If no issues, write directly: 无问题\n\n"+
			"### Conclusion\n"+
			"No issues: LGTM, one sentence explaining it can be merged.\n"+
			"Has issues: List CRITICAL/MAJOR items that must be fixed; MINOR/NIT can be summarized in one sentence.\n\n"+
			"---\n\n"+
			"## Diff\n"+
			"```diff\n%s\n```",
		"You are a senior software engineer conducting a code review on the following git diff.\n\n",
		filePath,
		focusSection,
		langSection,
		diff,
	)
}

// buildFocusSection 构建焦点区域提示
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

// languageGuidelines 返回语言特定的检查指南
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
