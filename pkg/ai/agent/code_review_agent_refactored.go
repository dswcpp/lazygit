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

// ────────────────────────────────────────────────────────────────────────────
// CodeReviewAgent - LangGraph架构重构版本
// ────────────────────────────────────────────────────────────────────────────

// 代码评审节点ID（复用TwoPhaseAgent的NodeID类型）
const (
	NodeReviewInit     NodeID = "review_init"     // 初始化评审
	NodeReviewing      NodeID = "reviewing"       // 执行评审
	NodeReviewDone     NodeID = "review_done"     // 评审完成
	NodeWaitQuestion   NodeID = "wait_question"   // 等待用户追问
	NodeHandleQuestion NodeID = "handle_question" // 处理追问
	NodeReviewEnd      NodeID = "end"             // 结束
)

// CodeReviewAgentV2 代码评审Agent（LangGraph架构）
// 特性：
// 1. 基于Graph节点的控制流
// 2. 纯函数式节点（所有状态更新通过返回值）
// 3. 支持检查点恢复
// 4. 支持交互式追问
// 5. 超时控制和错误恢复
type CodeReviewAgentV2 struct {
	mu           sync.Mutex
	provider     provider.Provider
	tr           *aii18n.Translator
	state        CodeReviewState
	graph        *CodeReviewGraph
	checkpointer CodeReviewCheckpointer
	threadID     string
	timeout      time.Duration
}

// NewCodeReviewAgentV2 创建代码评审Agent（LangGraph架构）
func NewCodeReviewAgentV2(p provider.Provider, tr *aii18n.Translator) *CodeReviewAgentV2 {
	a := &CodeReviewAgentV2{
		provider: p,
		tr:       tr,
		timeout:  30 * time.Second,
	}
	a.graph = a.buildGraph()
	return a
}

// SetCheckpointer 设置检查点器（支持中断恢复）
func (a *CodeReviewAgentV2) SetCheckpointer(c CodeReviewCheckpointer, threadID string) {
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

// Review 执行代码评审（统一入口）
func (a *CodeReviewAgentV2) Review(
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

	// 执行Graph
	return a.runGraph(ctx, NodeReviewInit, onChunk)
}

// Ask 追问（交互式评审）
func (a *CodeReviewAgentV2) Ask(
	ctx context.Context,
	question string,
	onChunk func(string),
) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 检查状态
	if !a.canAsk() {
		return fmt.Errorf("cannot ask question in phase: %s", a.state.Phase)
	}

	// 更新状态
	a.state = a.state.WithUserQuestion(question).WithPhase(PhaseReviewInteractive)

	// 执行Graph（从NodeHandleQuestion开始）
	return a.runGraph(ctx, NodeHandleQuestion, onChunk)
}

// GetState 获取当前状态（线程安全）
func (a *CodeReviewAgentV2) GetState() CodeReviewState {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state
}

// Phase 返回当前阶段
func (a *CodeReviewAgentV2) Phase() ReviewPhase {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state.Phase
}

// CanAsk 是否可以追问
func (a *CodeReviewAgentV2) CanAsk() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.canAsk()
}

func (a *CodeReviewAgentV2) canAsk() bool {
	return a.state.Phase == PhaseReviewDone || a.state.Phase == PhaseReviewWaiting
}

// ────────────────────────────────────────────────────────────────────────────
// Graph执行
// ────────────────────────────────────────────────────────────────────────────

// runGraph 执行Graph（从指定节点开始）
func (a *CodeReviewAgentV2) runGraph(
	ctx context.Context,
	startNode NodeID,
	onChunk func(string),
) error {
	newState, err := a.graph.Run(ctx, startNode, a.state, onChunk)
	if err != nil {
		return err
	}
	a.state = newState
	return nil
}

// ────────────────────────────────────────────────────────────────────────────
// 节点函数（纯函数，所有状态更新通过返回值）
// ────────────────────────────────────────────────────────────────────────────

// nodeReviewInit 初始化评审节点
func (a *CodeReviewAgentV2) nodeReviewInit(
	ctx context.Context,
	state CodeReviewState,
	onChunk func(string),
) (NodeID, CodeReviewState, error) {
	// 构建系统prompt和用户prompt
	systemPrompt := a.tr.SkillCodeReviewSystemPrompt()
	userPrompt := a.buildReviewPrompt(state.FilePath, state.Language, state.Focus, state.Diff)

	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: systemPrompt},
		{Role: provider.RoleUser, Content: userPrompt},
	}

	state = state.WithMessages(messages)
	return NodeReviewing, state, nil
}

// nodeReviewing 执行评审节点
func (a *CodeReviewAgentV2) nodeReviewing(
	ctx context.Context,
	state CodeReviewState,
	onChunk func(string),
) (NodeID, CodeReviewState, error) {
	state = state.WithPhase(PhaseReviewing)

	// 带超时的context
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	// 流式调用
	var buffer strings.Builder
	err := a.provider.CompleteStream(ctx, state.Messages, func(chunk string) {
		buffer.WriteString(chunk)
		if onChunk != nil {
			onChunk(chunk)
		}
	})

	if err != nil {
		return NodeReviewEnd, state.WithError(err.Error()), err
	}

	result := buffer.String()
	state = state.WithResult(result)

	return NodeReviewDone, state, nil
}

// nodeReviewDone 评审完成节点
func (a *CodeReviewAgentV2) nodeReviewDone(
	ctx context.Context,
	state CodeReviewState,
	onChunk func(string),
) (NodeID, CodeReviewState, error) {
	state = state.WithPhase(PhaseReviewWaiting)

	// 保存检查点
	a.saveCheckpoint(state)

	return NodeWaitQuestion, state, nil
}

// nodeWaitQuestion 等待用户追问节点
func (a *CodeReviewAgentV2) nodeWaitQuestion(
	ctx context.Context,
	state CodeReviewState,
	onChunk func(string),
) (NodeID, CodeReviewState, error) {
	// 设置恢复点
	state = state.WithResumeFrom(NodeHandleQuestion)
	a.saveCheckpoint(state)

	// 挂起，等待用户输入
	return NodeReviewEnd, state, nil
}

// nodeHandleQuestion 处理追问节点
func (a *CodeReviewAgentV2) nodeHandleQuestion(
	ctx context.Context,
	state CodeReviewState,
	onChunk func(string),
) (NodeID, CodeReviewState, error) {
	question := state.UserQuestion
	if question == "" {
		return NodeReviewEnd, state, fmt.Errorf("no question provided")
	}

	// 追加用户问题到消息历史
	state = state.AppendMessage(provider.Message{
		Role:    provider.RoleUser,
		Content: question,
	})

	// 带超时的context
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	// 流式调用
	var buffer strings.Builder
	err := a.provider.CompleteStream(ctx, state.Messages, func(chunk string) {
		buffer.WriteString(chunk)
		if onChunk != nil {
			onChunk(chunk)
		}
	})

	if err != nil {
		return NodeReviewEnd, state.WithError(err.Error()), err
	}

	// 追加AI回复
	answer := buffer.String()
	state = state.AppendMessage(provider.Message{
		Role:    provider.RoleAssistant,
		Content: answer,
	})

	// 清除用户问题，回到等待状态
	state = state.WithUserQuestion("").WithPhase(PhaseReviewWaiting)

	// 保存检查点
	a.saveCheckpoint(state)

	return NodeWaitQuestion, state, nil
}

// ────────────────────────────────────────────────────────────────────────────
// 辅助方法
// ────────────────────────────────────────────────────────────────────────────

// validateDiff 验证diff大小
func (a *CodeReviewAgentV2) validateDiff(diff string) error {
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
func (a *CodeReviewAgentV2) saveCheckpoint(state CodeReviewState) {
	if a.checkpointer != nil {
		_ = a.checkpointer.Save(a.threadID, state)
	}
}

// clearCheckpoint 清除检查点
func (a *CodeReviewAgentV2) clearCheckpoint() {
	if a.checkpointer != nil {
		a.checkpointer.Clear(a.threadID)
	}
}

// buildReviewPrompt 构建评审prompt
func (a *CodeReviewAgentV2) buildReviewPrompt(filePath, lang, focus, diff string) string {
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

// ────────────────────────────────────────────────────────────────────────────
// Graph构建
// ────────────────────────────────────────────────────────────────────────────

// CodeReviewGraph 代码评审Graph
type CodeReviewGraph struct {
	nodes map[NodeID]CodeReviewNodeFunc
}

// CodeReviewNodeFunc 节点函数类型
type CodeReviewNodeFunc func(
	ctx context.Context,
	state CodeReviewState,
	onChunk func(string),
) (NodeID, CodeReviewState, error)

// NewCodeReviewGraph 创建Graph
func NewCodeReviewGraph() *CodeReviewGraph {
	return &CodeReviewGraph{
		nodes: make(map[NodeID]CodeReviewNodeFunc),
	}
}

// AddNode 添加节点
func (g *CodeReviewGraph) AddNode(id NodeID, fn CodeReviewNodeFunc) {
	g.nodes[id] = fn
}

// Run 执行Graph
func (g *CodeReviewGraph) Run(
	ctx context.Context,
	startNode NodeID,
	initialState CodeReviewState,
	onChunk func(string),
) (CodeReviewState, error) {
	state := initialState
	currentNode := startNode

	for currentNode != NodeReviewEnd {
		if ctx.Err() != nil {
			return state, ctx.Err()
		}

		nodeFn, ok := g.nodes[currentNode]
		if !ok {
			return state, fmt.Errorf("unknown node: %s", currentNode)
		}

		nextNode, newState, err := nodeFn(ctx, state, onChunk)
		if err != nil {
			return newState, err
		}

		state = newState
		currentNode = nextNode
	}

	return state, nil
}

// buildGraph 构建Graph
func (a *CodeReviewAgentV2) buildGraph() *CodeReviewGraph {
	g := NewCodeReviewGraph()
	g.AddNode(NodeReviewInit, a.nodeReviewInit)
	g.AddNode(NodeReviewing, a.nodeReviewing)
	g.AddNode(NodeReviewDone, a.nodeReviewDone)
	g.AddNode(NodeWaitQuestion, a.nodeWaitQuestion)
	g.AddNode(NodeHandleQuestion, a.nodeHandleQuestion)
	return g
}

// 注意：buildFocusSection, languageGuidelines, detectLanguage
// 这些辅助函数已在 code_review_agent.go 和 code_review_state.go 中定义
// 可以直接复用，无需重复声明
