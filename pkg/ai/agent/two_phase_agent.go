package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/repocontext"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

const defaultMaxPlanSteps = 15
const defaultStepTimeout = 30 * time.Second
const maxPlanMessages = 50 // 最大规划消息历史数量，防止内存无限增长

// toolAliases 工具名别名映射，用于容错常见的工具名错误
var toolAliases = map[string]string{
	"add":     "stage_all",
	"git_add": "stage_all",
	"unstage": "unstage_all",
	"switch":  "checkout",
	"branch":  "create_branch",
}

// confirmKeywords 触发执行的关键词（小写匹配）
var confirmKeywords = []string{
	"y", "yes", "是", "确认", "执行", "ok", "好", "好的", "继续", "同意", "可以",
}

// denyKeywords 触发取消的关键词（小写匹配）
var denyKeywords = []string{
	"n", "no", "否", "取消", "算了", "不", "不要", "停", "放弃",
}

func isConfirmMsg(msg string) bool {
	s := strings.ToLower(strings.TrimSpace(msg))
	for _, kw := range confirmKeywords {
		if s == kw {
			return true
		}
	}
	return false
}

func isDenyMsg(msg string) bool {
	s := strings.ToLower(strings.TrimSpace(msg))
	for _, kw := range denyKeywords {
		if s == kw {
			return true
		}
	}
	return false
}

// TwoPhaseAgent 实现聊天终端风格的两阶段工作流：
//
//	阶段一（PhasePlanning）：LLM 使用只读工具 + SkillTool 收集信息，
//	  输出执行计划后设置 PhaseWaitingConfirm 并返回，等待用户下一条消息。
//
//	阶段二（PhaseExecuting）：用户输入 Y/Yes 触发执行；输入 N/No 取消；
//	  输入其他内容则视为补充说明，重新进入规划循环调整计划。
//
// 整个交互完全在聊天窗口中完成，不弹出任何对话框。
//
// Thread Safety: Send() is NOT thread-safe. Callers must serialize calls.
type TwoPhaseAgent struct {
	mu           sync.Mutex // protects state
	provider     provider.Provider
	fullRegistry *tools.Registry // 完整注册表（执行阶段使用）
	readRegistry *tools.Registry // 只读注册表（规划阶段使用）
	session      *Session
	tr           *aii18n.Translator // i18n translator
	maxPlanSteps int
	stepTimeout  time.Duration
	state        GraphState   // single source of truth for control-flow state
	graph        *Graph       // compiled once in NewTwoPhaseAgent; reused on every run
	checkpointer Checkpointer // optional; nil means no persistence
	threadID     string       // identifies this conversation in the checkpointer
}

// NewTwoPhaseAgent 创建 TwoPhaseAgent。
func NewTwoPhaseAgent(
	p provider.Provider,
	fullRegistry *tools.Registry,
	readRegistry *tools.Registry,
	session *Session,
	tr *aii18n.Translator,
) *TwoPhaseAgent {
	a := &TwoPhaseAgent{
		provider:     p,
		fullRegistry: fullRegistry,
		readRegistry: readRegistry,
		session:      session,
		tr:           tr,
		maxPlanSteps: defaultMaxPlanSteps,
		stepTimeout:  defaultStepTimeout,
		state: GraphState{
			Phase:           PhasePlanning,
			ToolCallHistory: make(map[string]int),
		},
	}
	a.graph = a.buildGraph() // compiled once; stateless, safe to reuse across runs
	return a
}

// SetCheckpointer attaches a Checkpointer to this agent.
// threadID identifies the conversation; use a stable ID (e.g. repo path + session UUID)
// so the checkpoint survives across process restarts.
// If a saved state exists and has a pending ResumeFrom, the agent restores it
// automatically so the user can pick up the interrupted conversation.
func (a *TwoPhaseAgent) SetCheckpointer(c Checkpointer, threadID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.checkpointer = c
	a.threadID = threadID
	if saved, ok := c.Load(threadID); ok && saved.ResumeFrom != "" {
		a.state = saved
	}
}

// saveCheckpoint persists state when the graph suspends at a human interrupt.
// state is passed explicitly so the node never needs to write back to a.state
// before returning — keeping it a pure state transformer.
func (a *TwoPhaseAgent) saveCheckpoint(state GraphState) {
	if a.checkpointer != nil {
		_ = a.checkpointer.Save(a.threadID, state)
	}
}

// clearCheckpoint removes the checkpoint when the conversation ends.
func (a *TwoPhaseAgent) clearCheckpoint() {
	if a.checkpointer != nil {
		a.checkpointer.Clear(a.threadID)
	}
}

// getGraph returns the compiled graph, building it lazily if not yet done.
// NewTwoPhaseAgent pre-compiles the graph for efficiency; this fallback
// supports test agents constructed via struct literals.
func (a *TwoPhaseAgent) getGraph() *Graph {
	if a.graph == nil {
		a.graph = a.buildGraph()
	}
	return a.graph
}

// Phase returns the current agent phase (thread-safe).
func (a *TwoPhaseAgent) Phase() AgentPhase {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state.Phase
}

// Plan returns the current execution plan, or nil if not yet produced (thread-safe).
func (a *TwoPhaseAgent) Plan() *ExecutionPlan {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state.Plan
}

// Session 返回会话，供 GUI 层读取 UIMessages 进行渲染。
// DEPRECATED: UIMessages 已迁移到 GraphState，此方法同步 state.UIMessages 到 session 以保持向后兼容。
func (a *TwoPhaseAgent) Session() *Session {
	a.mu.Lock()
	defer a.mu.Unlock()
	// 同步 GraphState 的 UIMessages 到 Session（向后兼容）
	a.session.UIMessages = a.state.UIMessages
	return a.session
}

// Send is the unified chat entry point.
//
// Routing logic (LangGraph analogy: "which node to start from"):
//   - PhaseExecuting         → ignore (execution in progress)
//   - state.ResumeFrom != "" → resume graph from the saved checkpoint
//   - otherwise              → start a fresh planning run
//
// Thread Safety: NOT thread-safe. Callers must serialize calls to Send().
// onUpdate is called after each state change on the same goroutine;
// the GUI must switch to the UI thread (e.g. gocui.Update) before rendering.
func (a *TwoPhaseAgent) Send(
	ctx context.Context,
	userMsg string,
	repoCtx repocontext.RepoContext,
	onUpdate func(),
) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Ignore input while the execution phase is running.
	if a.state.Phase == PhaseExecuting {
		a.state = a.state.WithSystemNote(a.tr.TwoPhaseAgentExecuting())
		if onUpdate != nil {
			onUpdate()
		}
		return nil
	}

	// Resume from a pending interrupt (e.g. PhaseWaitingConfirm).
	if a.state.ResumeFrom != "" {
		a.state = a.state.WithUserMessage(userMsg)
		if onUpdate != nil {
			onUpdate()
		}
		return a.resume(ctx, userMsg, onUpdate)
	}

	// Fresh plan (PhasePlanning / PhaseDone / PhaseCancelled).
	return a.startPlan(ctx, userMsg, repoCtx, onUpdate)
}

// startPlan 开始一次全新的规划（清空上一轮状态）。
func (a *TwoPhaseAgent) startPlan(
	ctx context.Context,
	userMsg string,
	repoCtx repocontext.RepoContext,
	onUpdate func(),
) error {
	// 重置状态，清除上一轮检查点
	a.clearCheckpoint()
	a.state.Reset() // Phase=PhasePlanning, Plan=nil, PlanMessages=nil, ToolCallHistory cleared

	// 构建规划阶段 system prompt（含工具列表）
	sysPrompt := a.tr.BuildPlanningSystemPrompt()
	if toolSection := a.readRegistry.SystemPromptSection(tools.PermReadOnly); toolSection != "" {
		sysPrompt += "\n\n" + toolSection
	}

	// 初始用户消息：仓库上下文 + 用户指令
	initMsg := fmt.Sprintf("%s%s\n\n%s%s",
		a.tr.TwoPhaseAgentRepoStatusTitle(), repoCtx.CompactString(a.tr),
		a.tr.TwoPhaseAgentUserInstructionTitle(), userMsg)

	a.state = a.state.WithUserMessage(initMsg)
	if onUpdate != nil {
		onUpdate()
	}

	// 初始化规划消息历史（独立于 ProviderMessages，带专用 system prompt）
	a.state = a.state.WithPlanMessages([]provider.Message{
		{Role: provider.RoleSystem, Content: sysPrompt},
		{Role: provider.RoleUser, Content: initMsg},
	})

	return a.planLoop(ctx, onUpdate)
}

// resume continues a graph run that was suspended at a human interrupt point.
// It writes the user's reply into GraphState.HumanInput, then resumes the
// compiled graph from the saved checkpoint (typically NodeHandleConfirmation).
func (a *TwoPhaseAgent) resume(ctx context.Context, userMsg string, onUpdate func()) error {
	resumeFrom := a.state.ResumeFrom
	// Use temporary state to ensure atomicity: only update a.state on success
	tempState := a.state.WithResumeFrom("").WithHumanInput(userMsg)
	newState, err := a.getGraph().Run(ctx, resumeFrom, tempState, onUpdate)
	if err != nil {
		return err // a.state unchanged on error
	}
	a.state = newState
	return nil
}

// planLoop drives the planning phase through the compiled agent graph,
// starting at NodePlan. The graph threads state through each node; on
// return the agent's state is updated to reflect all node transitions.
func (a *TwoPhaseAgent) planLoop(ctx context.Context, onUpdate func()) error {
	newState, err := a.getGraph().Run(ctx, NodePlan, a.state, onUpdate)
	if err != nil {
		return err
	}
	a.state = newState
	return nil
}

// getAIResponse 获取 AI 的响应（流式或非流式）。
// messages 是调用方从 state 中传入的规划历史，不直接访问 a.state。
// 纯函数：通过返回 state 更新状态。
func (a *TwoPhaseAgent) getAIResponse(ctx context.Context, state GraphState, messages []provider.Message, onUpdate func()) (string, GraphState, error) {
	shouldStream := len(messages) > 1 && !a.hasRecentToolCalls(messages)

	if shouldStream {
		return a.streamPlanResponse(ctx, state, messages, onUpdate)
	}

	result, err := a.provider.Complete(ctx, messages)
	if err != nil {
		return "", state, err
	}
	return result.Content, state, nil
}

// handlePlanBlock 处理计划块，验证计划并更新 state。
// 返回 (handled, updatedState, error)：handled=false 表示验证失败、需重规划。
// Phase 转换（→ PhaseWaitingConfirm）由 nodeWaitHuman 负责，不在此处设置。
// 纯函数：所有状态更新通过返回值传递。
func (a *TwoPhaseAgent) handlePlanBlock(
	state GraphState,
	parsed tools.ParsedPlan,
	rawContent string,
	onUpdate func(),
) (bool, GraphState, error) {
	plan := a.buildExecutionPlan(parsed)

	if errors := a.validatePlan(plan); len(errors) > 0 {
		state = a.handlePlanValidationErrors(state, errors, onUpdate)
		return false, state, nil
	}

	state = state.WithPlan(plan)
	state = state.WithPlanUIMessage(plan)

	displayText := tools.StripPlanBlock(rawContent)
	if displayText == "" {
		displayText = a.defaultConfirmPrompt(plan)
	}
	state = state.WithAssistantMessage(displayText)

	state = state.AppendPlanMessage(provider.Message{
		Role:    provider.RoleAssistant,
		Content: rawContent,
	}, maxPlanMessages)

	if onUpdate != nil {
		onUpdate()
	}
	return true, state, nil
}

// handlePlanValidationErrors 处理计划验证错误，返回更新后的 state。
// 纯函数：通过返回值更新状态。
func (a *TwoPhaseAgent) handlePlanValidationErrors(state GraphState, errors []string, onUpdate func()) GraphState {
	var sb strings.Builder
	sb.WriteString(a.tr.TwoPhaseAgentPlanErrorsIntro())
	for i, err := range errors {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, err))
	}
	sb.WriteString(a.tr.TwoPhaseAgentPlanRegeneratePrompt())

	state = state.WithSystemNote(a.tr.TwoPhaseAgentPlanValidationFailed())
	state = state.AppendPlanMessage(provider.Message{
		Role:    provider.RoleUser,
		Content: sb.String(),
	}, maxPlanMessages)

	if onUpdate != nil {
		onUpdate()
	}
	return state
}

// handleEmptyResponse 处理空响应（无工具调用也无计划），返回更新后的 state。
func (a *TwoPhaseAgent) handleEmptyResponse(state GraphState, onUpdate func()) (GraphState, error) {
	state = state.WithEmptyResponseCount(state.EmptyResponseCount + 1)
	if state.EmptyResponseCount >= 3 {
		return state, fmt.Errorf("%s", a.tr.TwoPhaseAgentEmptyResponseError())
	}
	state = state.AppendPlanMessage(provider.Message{
		Role:    provider.RoleUser,
		Content: a.tr.TwoPhaseAgentContinueAnalysis(),
	}, maxPlanMessages)
	return state, nil
}

// executeToolCalls 执行工具调用列表，返回更新后的 state。
// 纯函数：所有状态更新通过返回值传递。
func (a *TwoPhaseAgent) executeToolCalls(
	ctx context.Context,
	state GraphState,
	toolCalls []tools.ToolCall,
	onUpdate func(),
) (GraphState, error) {
	for _, call := range toolCalls {
		if ctx.Err() != nil {
			return state, ctx.Err()
		}

		var skip bool
		skip, state = a.checkToolCallLimit(state, call, onUpdate)
		if skip {
			continue
		}

		tool, ok := a.readRegistry.Get(call.Name)
		if !ok {
			state = a.handleToolNotFound(state, call.Name, onUpdate)
			continue
		}

		state = state.WithToolCall(call)
		if onUpdate != nil {
			onUpdate()
		}

		toolResult := tool.Execute(ctx, call)
		state = state.WithToolResult(toolResult, call.Name)
		state = state.AppendPlanMessage(provider.Message{
			Role:    provider.RoleUser,
			Content: a.tr.TwoPhaseAgentToolResultPrefix(call.Name, toolResult.Output),
		}, maxPlanMessages)
		if onUpdate != nil {
			onUpdate()
		}
	}
	return state, nil
}

// toolCallKey returns a deterministic cache key for a tool call.
// json.Marshal sorts map keys alphabetically, making the key stable across
// identical calls regardless of insertion order.
func toolCallKey(call tools.ToolCall) string {
	b, _ := json.Marshal(call.Params)
	return call.Name + ":" + string(b)
}

// checkToolCallLimit 检查工具调用次数限制。
// 返回 (skip, updatedState)：skip=true 表示调用方应跳过本次调用。
// 纯函数：通过返回值更新状态（使用不可变 map 更新）。
func (a *TwoPhaseAgent) checkToolCallLimit(state GraphState, call tools.ToolCall, onUpdate func()) (skip bool, _ GraphState) {
	callKey := toolCallKey(call)
	state = state.WithToolCallIncrement(callKey)

	if state.ToolCallHistory[callKey] > 3 {
		warnMsg := a.tr.TwoPhaseAgentToolCallWarning(call.Name, state.ToolCallHistory[callKey])
		state = state.WithSystemNote(warnMsg)
		state = state.AppendPlanMessage(provider.Message{
			Role:    provider.RoleUser,
			Content: a.tr.TwoPhaseAgentSystemPrefix() + warnMsg,
		}, maxPlanMessages)
		if onUpdate != nil {
			onUpdate()
		}
		return true, state
	}
	return false, state
}

// handleToolNotFound 处理工具未找到的情况，返回更新后的 state。
// 纯函数：通过返回值更新状态。
func (a *TwoPhaseAgent) handleToolNotFound(state GraphState, toolName string, onUpdate func()) GraphState {
	errMsg := a.tr.AgentToolNotAllowedInPlanning(toolName)
	state = state.WithSystemNote(errMsg)
	state = state.AppendPlanMessage(provider.Message{
		Role:    provider.RoleUser,
		Content: a.tr.TwoPhaseAgentSystemPrefix() + errMsg,
	}, maxPlanMessages)
	if onUpdate != nil {
		onUpdate()
	}
	return state
}

// execute drives the execution phase through the compiled agent graph.
// It uses a.state.Plan set during planning; plan steps are visited one at a
// time via nodeExecuteStep until all are done or a critical failure occurs.
func (a *TwoPhaseAgent) execute(ctx context.Context, onUpdate func()) error {
	// Update phase immediately so GUI can show execution state
	a.state = a.state.WithPhase(PhaseExecuting).WithExecStepIndex(0)
	if onUpdate != nil {
		onUpdate()
	}
	newState, err := a.getGraph().Run(ctx, NodeExecuteStep, a.state, onUpdate)
	if err != nil {
		return err // Error: a.state keeps PhaseExecuting (user can see error in that phase)
	}
	a.state = newState // Success: update to final state (typically PhaseDone)
	return nil
}

// executeStep 执行单个步骤。
// 纯函数：通过返回 state 更新状态。
func (a *TwoPhaseAgent) executeStep(
	ctx context.Context,
	state GraphState,
	step *PlanStep,
	onUpdate func(),
) (GraphState, error) {
	step.Status = StepRunning // agent owns step field mutation; Session only observes
	state = state.WithStepUpdate(step, StepRunning, "", "")
	if onUpdate != nil {
		onUpdate()
	}

	// 应用工具名别名映射
	toolName := step.ToolName
	if alias, ok := toolAliases[toolName]; ok {
		toolName = alias
	}

	// 获取工具
	tool, ok := a.fullRegistry.Get(toolName)
	if !ok {
		return a.handleStepToolNotFound(state, step, toolName, onUpdate)
	}

	// 执行工具（带超时）
	return a.executeStepWithTimeout(ctx, state, step, tool, toolName, onUpdate)
}

// handleStepToolNotFound 处理步骤中工具未找到的情况。
// 纯函数：通过修改 step 字段和返回 state 更新状态。
func (a *TwoPhaseAgent) handleStepToolNotFound(
	state GraphState,
	step *PlanStep,
	toolName string,
	onUpdate func(),
) (GraphState, error) {
	errMsg := a.formatToolNotFoundError(step.ToolName, toolName)
	step.Status = StepFailed
	step.Error = errMsg
	state = state.WithStepUpdate(step, StepFailed, "", errMsg)
	if onUpdate != nil {
		onUpdate()
	}
	if step.Critical {
		return state, fmt.Errorf("%s", a.tr.AgentCriticalStepFailed(step.Description, errMsg))
	}
	return state, nil
}

// executeStepWithTimeout 执行步骤（带超时控制）。
// 纯函数：通过返回 state 更新状态。
func (a *TwoPhaseAgent) executeStepWithTimeout(
	ctx context.Context,
	state GraphState,
	step *PlanStep,
	tool tools.Tool,
	toolName string,
	onUpdate func(),
) (GraphState, error) {
	// 为每个步骤创建带超时的 context
	stepCtx, cancel := context.WithTimeout(ctx, a.stepTimeout)
	defer cancel()

	call := tools.ToolCall{
		ID:     fmt.Sprintf("exec_%s", step.ID),
		Name:   toolName,
		Params: step.Params,
	}

	// 在 goroutine 中执行工具，支持超时
	resultChan := make(chan tools.ToolResult, 1)
	go func() {
		resultChan <- tool.Execute(stepCtx, call)
	}()

	// 等待结果或超时
	select {
	case result := <-resultChan:
		return a.handleStepResult(state, step, result, onUpdate)
	case <-stepCtx.Done():
		if stepCtx.Err() == context.DeadlineExceeded {
			return a.handleStepTimeout(state, step, onUpdate)
		}
		return state, stepCtx.Err()
	}
}

// handleStepResult 处理步骤执行结果。
// 纯函数：通过修改 step 字段和返回 state 更新状态。
func (a *TwoPhaseAgent) handleStepResult(
	state GraphState,
	step *PlanStep,
	result tools.ToolResult,
	onUpdate func(),
) (GraphState, error) {
	if result.Success {
		step.Status = StepDone
		step.Result = result.Output
		state = state.WithStepUpdate(step, StepDone, result.Output, "")
	} else {
		errMsg := a.formatExecutionError(step, result.Output)
		step.Status = StepFailed
		step.Error = errMsg
		state = state.WithStepUpdate(step, StepFailed, "", errMsg)
		if step.Critical {
			if onUpdate != nil {
				onUpdate()
			}
			return state, fmt.Errorf("%s", a.tr.AgentCriticalStepFailed(step.Description, errMsg))
		}
	}
	if onUpdate != nil {
		onUpdate()
	}
	return state, nil
}

// handleStepTimeout 处理步骤超时。
// 纯函数：通过修改 step 字段和返回 state 更新状态。
func (a *TwoPhaseAgent) handleStepTimeout(state GraphState, step *PlanStep, onUpdate func()) (GraphState, error) {
	errMsg := a.formatTimeoutError(step, a.stepTimeout)
	step.Status = StepFailed
	step.Error = errMsg
	state = state.WithStepUpdate(step, StepFailed, "", errMsg)
	if onUpdate != nil {
		onUpdate()
	}
	if step.Critical {
		return state, fmt.Errorf("关键步骤超时: %s", step.Description)
	}
	return state, nil
}

// finishExecution 完成执行阶段，生成总结，返回更新后的 state。
// 纯函数：通过返回值更新状态。
func (a *TwoPhaseAgent) finishExecution(state GraphState, onUpdate func()) GraphState {
	state = state.WithPhase(PhaseDone)
	plan := state.Plan
	done := plan.DoneCount()
	total := len(plan.Steps)
	failed := len(plan.FailedSteps())
	summary := fmt.Sprintf("执行完成：%d/%d 步成功", done, total)
	if failed > 0 {
		summary += fmt.Sprintf("，%d 步失败", failed)
	}
	state = state.WithSystemNote(summary)
	if onUpdate != nil {
		onUpdate()
	}
	return state
}

// validatePlan 验证执行计划的有效性，返回错误列表。
// 检查项：
//   - 工具名是否存在（考虑别名映射）
//   - 必需参数是否提供
//   - 参数类型是否正确
func (a *TwoPhaseAgent) validatePlan(plan *ExecutionPlan) []string {
	var errors []string

	for _, step := range plan.Steps {
		toolName := step.ToolName

		// 检查别名映射
		if alias, ok := toolAliases[toolName]; ok {
			toolName = alias
		}

		// 检查工具是否存在
		tool, ok := a.fullRegistry.Get(toolName)
		if !ok {
			errors = append(errors, fmt.Sprintf(
				"步骤 %s: 未知工具 '%s'（请使用工具列表中的准确名称）",
				step.ID, step.ToolName))
			continue
		}

		// 验证参数
		if err := a.validateStepParams(step, tool); err != nil {
			errors = append(errors, fmt.Sprintf(
				"步骤 %s (%s): %s",
				step.ID, step.ToolName, err.Error()))
		}
	}

	return errors
}

// validateStepParams 验证步骤参数是否符合工具的 schema。
func (a *TwoPhaseAgent) validateStepParams(step *PlanStep, tool tools.Tool) error {
	schema := tool.Schema()

	// 检查必需参数
	for paramName, paramSchema := range schema.Params {
		if !paramSchema.Required {
			continue
		}

		value, ok := step.Params[paramName]
		if !ok || value == nil {
			return fmt.Errorf("缺少必需参数: %s", paramName)
		}

		// 检查空字符串
		if paramSchema.Type == "string" {
			if str, ok := value.(string); ok && strings.TrimSpace(str) == "" {
				return fmt.Errorf("参数 %s 不能为空", paramName)
			}
		}

		// 验证类型
		if err := validateParamType(paramName, value, paramSchema.Type); err != nil {
			return err
		}
	}

	return nil
}

// validateParamType 验证参数类型是否匹配。
func validateParamType(paramName string, value any, expectedType string) error {
	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("参数 %s 类型错误：期望 string，实际 %T", paramName, value)
		}
	case "int":
		switch value.(type) {
		case int, int64, float64:
			// JSON 解析可能产生 float64，需要兼容
		default:
			return fmt.Errorf("参数 %s 类型错误：期望 int，实际 %T", paramName, value)
		}
	case "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("参数 %s 类型错误：期望 bool，实际 %T", paramName, value)
		}
	}
	return nil
}

// buildExecutionPlan 将 tools.ParsedPlan 转换为 ExecutionPlan，
// 同时从工具注册表中查询每个步骤的权限级别。
func (a *TwoPhaseAgent) buildExecutionPlan(parsed tools.ParsedPlan) *ExecutionPlan {
	steps := make([]*PlanStep, 0, len(parsed.Steps))
	for _, s := range parsed.Steps {
		perm := tools.PermReadOnly
		toolName := s.ToolName

		// 应用别名映射
		if alias, ok := toolAliases[toolName]; ok {
			toolName = alias
		}

		if tool, ok := a.fullRegistry.Get(toolName); ok {
			perm = tool.Schema().Permission
		}

		steps = append(steps, &PlanStep{
			ID:          s.ID,
			Description: s.Description,
			ToolName:    s.ToolName, // 保留原始名称，执行时再映射
			Params:      s.Params,
			Permission:  perm,
			Critical:    s.Critical,
			Status:      StepPending,
		})
	}
	return &ExecutionPlan{
		Summary: parsed.Summary,
		Steps:   steps,
	}
}

// defaultConfirmPrompt 当 LLM 未在 plan 块外附上提示时，生成默认确认提示。
func (a *TwoPhaseAgent) defaultConfirmPrompt(plan *ExecutionPlan) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**计划：%s**\n\n", plan.Summary))
	for _, s := range plan.Steps {
		sb.WriteString(fmt.Sprintf("  %s. %s\n", s.ID, s.Description))
	}
	sb.WriteString("\n输入 **Y / Yes** 确认执行，**N / No** 取消，或直接输入补充说明来调整计划。")
	return sb.String()
}

// formatToolNotFoundError 格式化工具未找到错误，提供友好的错误信息和建议。
func (a *TwoPhaseAgent) formatToolNotFoundError(originalName, mappedName string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("❌ 未知工具: %s", originalName))

	if originalName != mappedName {
		sb.WriteString(fmt.Sprintf("（已尝试映射为 %s）", mappedName))
	}

	sb.WriteString("\n\n💡 可能的原因：")
	sb.WriteString("\n  • 工具名拼写错误")
	sb.WriteString("\n  • 该工具不存在于当前注册表中")

	// 提供相似工具建议
	if suggestions := a.findSimilarTools(originalName); len(suggestions) > 0 {
		sb.WriteString("\n\n📝 您是否想使用以下工具？")
		for _, s := range suggestions {
			sb.WriteString(fmt.Sprintf("\n  • %s", s))
		}
	}

	return sb.String()
}

// formatTimeoutError 格式化超时错误，提供友好的错误信息和建议。
func (a *TwoPhaseAgent) formatTimeoutError(step *PlanStep, timeout time.Duration) string {
	var sb strings.Builder
	sb.WriteString(a.tr.AgentStepTimeout(timeout.String(), step.Description))
	sb.WriteString(a.tr.AgentPossibleReasons())
	sb.WriteString("\n  • 操作耗时过长（如大文件处理）")
	sb.WriteString("\n  • 网络请求超时")
	sb.WriteString("\n  • 工具内部阻塞")
	sb.WriteString("\n\n🔧 建议：")
	sb.WriteString("\n  • 检查网络连接")
	sb.WriteString("\n  • 减小操作范围")
	sb.WriteString("\n  • 重试该操作")
	return sb.String()
}

// formatExecutionError 格式化执行错误，提供友好的错误信息和恢复建议。
func (a *TwoPhaseAgent) formatExecutionError(step *PlanStep, rawError string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("❌ 执行失败: %s", step.Description))
	sb.WriteString(fmt.Sprintf("\n\n错误详情：\n%s", rawError))

	// 根据工具类型提供特定建议
	suggestions := a.getRecoverySuggestions(step.ToolName, rawError)
	if len(suggestions) > 0 {
		sb.WriteString("\n\n🔧 恢复建议：")
		for _, s := range suggestions {
			sb.WriteString(fmt.Sprintf("\n  • %s", s))
		}
	}

	return sb.String()
}

// findSimilarTools 查找与给定名称相似的工具（简单的前缀匹配）。
func (a *TwoPhaseAgent) findSimilarTools(name string) []string {
	var similar []string
	lowerName := strings.ToLower(name)

	// 遍历所有工具，查找相似的
	allTools := a.fullRegistry.All()
	for _, tool := range allTools {
		toolName := tool.Schema().Name
		lowerToolName := strings.ToLower(toolName)

		// 前缀匹配或包含关系
		if strings.HasPrefix(lowerToolName, lowerName) ||
			strings.HasPrefix(lowerName, lowerToolName) ||
			strings.Contains(lowerToolName, lowerName) {
			similar = append(similar, toolName)
		}
	}

	// 最多返回 3 个建议
	if len(similar) > 3 {
		similar = similar[:3]
	}

	return similar
}

// getRecoverySuggestions 根据工具类型和错误信息提供恢复建议。
func (a *TwoPhaseAgent) getRecoverySuggestions(toolName, errorMsg string) []string {
	var suggestions []string
	lowerError := strings.ToLower(errorMsg)

	// 通用建议
	if strings.Contains(lowerError, "permission") || strings.Contains(lowerError, "权限") {
		suggestions = append(suggestions, "检查文件或目录权限")
	}

	if strings.Contains(lowerError, "not found") || strings.Contains(lowerError, "找不到") {
		suggestions = append(suggestions, "确认文件或分支存在")
	}

	if strings.Contains(lowerError, "conflict") || strings.Contains(lowerError, "冲突") {
		suggestions = append(suggestions, "解决冲突后重试")
	}

	// 工具特定建议
	switch toolName {
	case "commit", "git_commit":
		if strings.Contains(lowerError, "nothing to commit") {
			suggestions = append(suggestions, "先暂存文件（stage_all 或 stage_file）")
		}
		if strings.Contains(lowerError, "message") {
			suggestions = append(suggestions, "提供有效的提交信息")
		}

	case "checkout", "switch":
		if strings.Contains(lowerError, "uncommitted changes") ||
			strings.Contains(lowerError, "would be overwritten") ||
			strings.Contains(lowerError, "local changes") {
			suggestions = append(suggestions, "先提交或暂存当前修改")
		}

	case "push":
		if strings.Contains(lowerError, "rejected") ||
			strings.Contains(lowerError, "failed to push") {
			suggestions = append(suggestions, "先拉取远程更新（pull）")
		}
		if strings.Contains(lowerError, "no upstream") {
			suggestions = append(suggestions, "设置上游分支")
		}

	case "merge":
		if strings.Contains(lowerError, "conflict") {
			suggestions = append(suggestions, "手动解决冲突后继续")
		}
	}

	// 如果没有特定建议，提供通用建议
	if len(suggestions) == 0 {
		suggestions = append(suggestions, "检查操作参数是否正确")
		suggestions = append(suggestions, "查看完整错误信息")
	}

	return suggestions
}

// hasRecentToolCalls 检查消息历史中最近是否有工具调用。
func (a *TwoPhaseAgent) hasRecentToolCalls(messages []provider.Message) bool {
	if len(messages) < 2 {
		return false
	}
	lastMsg := messages[len(messages)-1]
	return strings.Contains(lastMsg.Content, "[工具结果")
}

// streamPlanResponse 使用流式输出获取 LLM 响应。
// 适用于纯文本回复场景（无需解析 plan/tool 块）。
// 纯函数：通过返回 state 更新状态。
func (a *TwoPhaseAgent) streamPlanResponse(ctx context.Context, state GraphState, messages []provider.Message, onUpdate func()) (string, GraphState, error) {
	state = state.StartStreaming()

	lastUpdateTime := time.Now()
	updateInterval := 50 * time.Millisecond

	var fullContent strings.Builder

	err := a.provider.CompleteStream(ctx, messages, func(chunk string) {
		if ctx.Err() != nil {
			return // context cancelled — stop appending but let provider drain
		}
		fullContent.WriteString(chunk)
		state = state.AppendStreamingChunk(chunk)

		now := time.Now()
		if now.Sub(lastUpdateTime) >= updateInterval {
			if onUpdate != nil {
				onUpdate()
			}
			lastUpdateTime = now
		}
	})

	state = state.FinishStreaming()
	if onUpdate != nil {
		onUpdate()
	}

	if err != nil {
		return "", state, err
	}
	return fullContent.String(), state, nil
}

// ── Graph nodes ──────────────────────────────────────────────────────────────
//
// Each node performs one focused unit of work and returns the next NodeID.
// They are wired together by buildGraph() below.

// nodePlan calls the LLM once, interprets the response, and routes to the
// appropriate next node:
//   - NodeWaitHuman  if a valid plan block was produced
//   - NodeCallTools  if the LLM emitted tool calls
//   - NodePlan       to retry on empty or validation-failure responses
// 纯函数：所有状态更新通过返回值传递。
func (a *TwoPhaseAgent) nodePlan(ctx context.Context, state GraphState, onUpdate func()) (NodeID, GraphState, error) {
	if state.PlanStepCount >= a.maxPlanSteps {
		return NodeEnd, state, fmt.Errorf("%s", a.tr.TwoPhaseAgentMaxStepsExceeded(a.maxPlanSteps))
	}
	state = state.WithPlanStepCount(state.PlanStepCount + 1)

	rawContent, newState, err := a.getAIResponse(ctx, state, state.PlanMessages, onUpdate)
	state = newState
	if err != nil {
		return NodeEnd, state, err
	}

	// Plan block produced → validate, update state, suspend for human confirmation.
	if parsed, ok := tools.ParsePlan(rawContent); ok {
		state = state.WithEmptyResponseCount(0)
		handled, newState, err := a.handlePlanBlock(state, parsed, rawContent, onUpdate)
		if err != nil {
			return NodeEnd, newState, err
		}
		if handled {
			return NodeWaitHuman, newState, nil
		}
		return NodePlan, newState, nil // validation failed — LLM will try again
	}

	// Record the assistant turn regardless of what follows.
	toolCalls := tools.ParseToolCalls(rawContent)
	displayText := tools.StripToolBlocks(rawContent)
	state = state.AppendPlanMessage(provider.Message{Role: provider.RoleAssistant, Content: rawContent}, maxPlanMessages)
	if displayText != "" {
		state = state.WithAssistantMessage(displayText)
	}
	if onUpdate != nil {
		onUpdate()
	}

	// No tool calls — empty response, inject a nudge and retry.
	if len(toolCalls) == 0 {
		state, err = a.handleEmptyResponse(state, onUpdate)
		if err != nil {
			return NodeEnd, state, err
		}
		return NodePlan, state, nil
	}

	// Queue tool calls for nodeCallTools.
	state = state.WithEmptyResponseCount(0)
	state = state.WithPendingToolCalls(toolCalls)
	return NodeCallTools, state, nil
}

// nodeCallTools executes the tool calls queued by nodePlan, then returns
// to nodePlan so the LLM can reason about the results.
func (a *TwoPhaseAgent) nodeCallTools(ctx context.Context, state GraphState, onUpdate func()) (NodeID, GraphState, error) {
	newState, err := a.executeToolCalls(ctx, state, state.PendingToolCalls, onUpdate)
	if err != nil {
		return NodeEnd, newState, err
	}
	newState.PendingToolCalls = nil
	return NodePlan, newState, nil
}

// nodeExecuteStep executes the next plan step and loops back to itself until
// all steps are exhausted, then routes to NodeDone.
// 纯函数：所有状态更新通过返回值传递。
func (a *TwoPhaseAgent) nodeExecuteStep(ctx context.Context, state GraphState, onUpdate func()) (NodeID, GraphState, error) {
	if state.Plan == nil {
		// Should never happen: execution phase requires a validated plan.
		return NodeEnd, state, fmt.Errorf("internal error: nodeExecuteStep reached with nil plan")
	}
	if state.ExecStepIndex >= len(state.Plan.Steps) {
		return NodeDone, state, nil
	}
	step := state.Plan.Steps[state.ExecStepIndex]
	state = state.WithExecStepIndex(state.ExecStepIndex + 1)
	newState, err := a.executeStep(ctx, state, step, onUpdate)
	if err != nil {
		return NodeEnd, newState, err
	}
	return NodeExecuteStep, newState, nil
}

// nodeDone records the execution summary, clears the checkpoint, and terminates
// the graph run.
func (a *TwoPhaseAgent) nodeDone(_ context.Context, state GraphState, onUpdate func()) (NodeID, GraphState, error) {
	state = a.finishExecution(state, onUpdate)
	a.clearCheckpoint()
	return NodeEnd, state, nil
}

// nodeWaitHuman transitions to PhaseWaitingConfirm, records the resume
// checkpoint, persists state, and terminates the current graph run.
//
// This is the correct node for the phase transition: handlePlanBlock (a helper
// called from nodePlan) validates and builds the plan, but setting the Phase
// is this node's responsibility — it owns the "suspend" semantic.
//
// LangGraph analogy: interrupt point enabling human-in-the-loop confirmation
// without re-running the planning phase on the next Send().
func (a *TwoPhaseAgent) nodeWaitHuman(_ context.Context, state GraphState, _ func()) (NodeID, GraphState, error) {
	state = state.WithPhase(PhaseWaitingConfirm)
	state = state.WithResumeFrom(NodeHandleConfirmation)
	a.saveCheckpoint(state) // pass local state directly; never touch a.state inside a node
	return NodeEnd, state, nil
}

// nodeHandleConfirmation is the resume entry point after nodeWaitHuman.
// It reads GraphState.HumanInput and routes to:
//   - NodeExecuteStep  on confirm (y/Y/yes/Yes/YES / 确认 / …)
//   - NodeEnd          on deny    (n/N/no/No/NO   / 取消 / …)
//   - NodePlan         on any other text (user provided feedback → replan)
// 纯函数：所有状态更新通过返回值传递。
func (a *TwoPhaseAgent) nodeHandleConfirmation(_ context.Context, state GraphState, onUpdate func()) (NodeID, GraphState, error) {
	input := state.HumanInput
	state = state.WithHumanInput("")

	switch {
	case isConfirmMsg(input):
		state = state.WithPhase(PhaseExecuting)
		state = state.WithExecStepIndex(0)
		if onUpdate != nil {
			onUpdate()
		}
		return NodeExecuteStep, state, nil

	case isDenyMsg(input):
		state = state.WithPhase(PhaseCancelled)
		state = state.WithSystemNote(a.tr.TwoPhaseAgentExecutionCancelled())
		a.clearCheckpoint()
		if onUpdate != nil {
			onUpdate()
		}
		return NodeEnd, state, nil

	default:
		// User provided feedback — reset planning cursors and replan.
		state = state.WithPhase(PhasePlanning)
		state = state.WithPlan(nil)
		state = state.WithPlanStepCount(0)
		state = state.WithEmptyResponseCount(0)
		state = state.WithPendingToolCalls(nil)
		feedbackMsg := a.tr.TwoPhaseAgentUserFeedbackPrompt(input)
		state = state.AppendPlanMessage(
			provider.Message{Role: provider.RoleUser, Content: feedbackMsg},
			maxPlanMessages,
		)
		if onUpdate != nil {
			onUpdate()
		}
		return NodePlan, state, nil
	}
}

// buildGraph wires all nodes together into a compiled, reusable Graph.
// Called once from NewTwoPhaseAgent; the result is stored in a.graph.
func (a *TwoPhaseAgent) buildGraph() *Graph {
	g := NewGraph()
	g.AddNode(NodePlan, a.nodePlan)
	g.AddNode(NodeCallTools, a.nodeCallTools)
	g.AddNode(NodeWaitHuman, a.nodeWaitHuman)
	g.AddNode(NodeHandleConfirmation, a.nodeHandleConfirmation)
	g.AddNode(NodeExecuteStep, a.nodeExecuteStep)
	g.AddNode(NodeDone, a.nodeDone)
	return g
}
