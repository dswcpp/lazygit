package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/repocontext"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

const defaultMaxPlanSteps = 15
const defaultStepTimeout = 30 * time.Second // 每个步骤的默认超时时间

// toolAliases 工具名别名映射，用于容错常见的工具名错误
var toolAliases = map[string]string{
	"add":      "stage_all",
	"git_add":  "stage_all",
	"unstage":  "unstage_all",
	"switch":   "checkout",
	"branch":   "create_branch",
}

// planningSystemPrompt 是规划阶段专用的 system prompt。
const planningSystemPrompt = `你是 lazygit 内置 AI，负责分析用户需求并制定 Git 操作计划。

## 工作流程

1. 调用只读工具（get_status、get_diff 等）收集必要信息
2. **如需生成提交信息**：
   - 先调用 get_staged_diff 获取暂存区 diff
   - 然后调用 commit_msg 工具生成提交信息（返回的内容直接用作 commit 的 message 参数）
   - **重要**：commit_msg 只能在规划阶段调用，不能放入执行计划
3. **如需生成分支名**：
   - 调用 branch_name 工具生成分支名
   - **重要**：branch_name 只能在规划阶段调用，不能放入执行计划
4. 信息收集完毕后，输出一个 ` + "```plan" + ` 块，内含完整执行计划
5. ` + "```plan" + ` 块之后附上一段简短的自然语言说明，提示用户可以输入 Y 确认、N 取消，或补充说明
6. 严禁在规划阶段调用任何写操作工具

## 重要：工具名规范

**必须使用下方工具列表中的准确工具名**，不要使用 git 命令名：
- ✅ 暂存文件：stage_all（暂存所有）或 stage_file（暂存单个文件）
- ❌ 不要使用：add、git_add
- ✅ 提交：commit（参数 message）
- ❌ 不要使用：git_commit
- ✅ 切换分支：checkout
- ❌ 不要使用：switch
- ✅ 创建分支：create_branch
- ❌ 不要使用：branch

## 特殊工具说明

**commit_msg 和 branch_name 是辅助工具，只能在规划阶段调用**：
- 在规划阶段调用 commit_msg 获取提交信息
- 将返回的提交信息作为 commit 工具的 message 参数
- **不要**把 commit_msg 放入执行计划的 steps 中

示例：
` + "```tool" + `
{"name": "commit_msg", "params": {"diff": "..."}}
` + "```" + `
返回: "feat: 添加用户登录功能"

然后在执行计划中：
` + "```plan" + `
{
  "steps": [
    {"tool": "commit", "params": {"message": "feat: 添加用户登录功能"}}
  ]
}
` + "```" + `

## 计划格式

` + "```plan" + `
{
  "summary": "整体描述（一句话）",
  "steps": [
    {
      "id": "1",
      "description": "步骤的人类可读描述",
      "tool": "工具名",
      "params": {"参数名": "具体值"},
      "critical": true
    }
  ]
}
` + "```" + `

## 注意事项

- 所有步骤的 params 必须是具体值，不能留占位符
- critical=true 表示该步骤失败则中止整个执行
- critical=false 表示失败后跳过并继续
- 只包含必要步骤`

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
//	阶段二（PhaseExecuting）：用户输入 Y 触发执行；输入 N 取消；
//	  输入其他内容则视为补充说明，重新进入规划循环调整计划。
//
// 整个交互完全在聊天窗口中完成，不弹出任何对话框。
type TwoPhaseAgent struct {
	provider     provider.Provider
	fullRegistry *tools.Registry // 完整注册表（执行阶段使用）
	readRegistry *tools.Registry // 只读注册表（规划阶段使用）
	session      *Session
	tr           *aii18n.Translator  // i18n translator
	maxPlanSteps int
	stepTimeout  time.Duration // 每个步骤的超时时间
	// planMessages 保存规划阶段的完整消息历史（含工具调用结果）。
	// 在 replan 时直接追加用户反馈后继续循环，避免重复调用只读工具。
	planMessages []provider.Message
	// toolCallHistory 记录每个工具调用的次数，防止无限循环
	toolCallHistory map[string]int
}

// NewTwoPhaseAgent 创建 TwoPhaseAgent。
//   - fullRegistry: 包含所有工具的注册表（执行阶段使用）
//   - readRegistry: 仅包含只读工具 + SkillTool 的注册表（规划阶段使用）
func NewTwoPhaseAgent(
	p provider.Provider,
	fullRegistry *tools.Registry,
	readRegistry *tools.Registry,
	session *Session,
	tr *aii18n.Translator,
) *TwoPhaseAgent {
	return &TwoPhaseAgent{
		provider:     p,
		fullRegistry: fullRegistry,
		readRegistry: readRegistry,
		session:      session,
		tr:           tr,
		maxPlanSteps: defaultMaxPlanSteps,
		stepTimeout:  defaultStepTimeout,
	}
}

// Session 返回会话，供 GUI 层读取 UIMessages 进行渲染。
func (a *TwoPhaseAgent) Session() *Session { return a.session }

// Send 是聊天终端风格的统一入口，处理用户发来的每一条消息。
// 行为由当前阶段决定：
//   - PhasePlanning / PhaseDone / PhaseCancelled → 开始新的规划
//   - PhaseWaitingConfirm → Y 执行 / N 取消 / 其他文本 → 调整计划重新规划
//   - PhaseExecuting → 忽略（执行中不接受新指令）
//
// 必须从非 UI goroutine 调用。
// onUpdate 在每次会话状态变化后被调用（与 Send 同一 goroutine），
// GUI 需通过 gocui.Update 等机制切换到 UI 线程刷新视图。
func (a *TwoPhaseAgent) Send(
	ctx context.Context,
	userMsg string,
	repoCtx repocontext.RepoContext,
	onUpdate func(),
) error {
	switch a.session.Phase {
	case PhaseWaitingConfirm:
		return a.handleConfirmation(ctx, userMsg, repoCtx, onUpdate)
	case PhaseExecuting:
		a.session.AddSystemNote("执行中，请稍候...")
		if onUpdate != nil {
			onUpdate()
		}
		return nil
	default: // PhasePlanning / PhaseDone / PhaseCancelled
		return a.startPlan(ctx, userMsg, repoCtx, onUpdate)
	}
}

// startPlan 开始一次全新的规划（清空上一轮状态）。
func (a *TwoPhaseAgent) startPlan(
	ctx context.Context,
	userMsg string,
	repoCtx repocontext.RepoContext,
	onUpdate func(),
) error {
	// 重置状态
	a.session.Reset()
	a.planMessages = nil
	a.toolCallHistory = make(map[string]int) // 重置工具调用历史
	a.session.SetPhase(PhasePlanning)

	// 构建规划阶段 system prompt（含工具列表）
	sysPrompt := planningSystemPrompt
	if toolSection := a.readRegistry.SystemPromptSection(tools.PermReadOnly); toolSection != "" {
		sysPrompt += "\n\n" + toolSection
	}

	// 初始用户消息：仓库上下文 + 用户指令
	initMsg := fmt.Sprintf("## 当前仓库状态\n\n%s\n\n## 用户指令\n\n%s",
		repoCtx.CompactString(a.tr), userMsg)

	a.session.AddUserMessage(initMsg)
	if onUpdate != nil {
		onUpdate()
	}

	// 初始化规划消息历史（独立于 session.providerMessages，带专用 system prompt）
	a.planMessages = []provider.Message{
		{Role: provider.RoleSystem, Content: sysPrompt},
		{Role: provider.RoleUser, Content: initMsg},
	}

	return a.planLoop(ctx, onUpdate)
}

// handleConfirmation 处理 PhaseWaitingConfirm 阶段的用户输入。
func (a *TwoPhaseAgent) handleConfirmation(
	ctx context.Context,
	userMsg string,
	repoCtx repocontext.RepoContext,
	onUpdate func(),
) error {
	a.session.AddUserMessage(userMsg)
	if onUpdate != nil {
		onUpdate()
	}

	switch {
	case isConfirmMsg(userMsg):
		return a.execute(ctx, a.session.Plan, onUpdate)

	case isDenyMsg(userMsg):
		a.session.SetPhase(PhaseCancelled)
		a.session.AddSystemNote("已取消执行计划。")
		if onUpdate != nil {
			onUpdate()
		}
		return nil

	default:
		// 用户提供了补充说明，重新规划（复用已有工具调用结果）
		return a.replan(ctx, userMsg, onUpdate)
	}
}

// replan 在用户提供补充说明后，追加反馈到现有规划历史并继续规划循环。
// 不重新收集只读工具数据，直接让 LLM 基于已有信息修改计划。
func (a *TwoPhaseAgent) replan(ctx context.Context, feedback string, onUpdate func()) error {
	a.session.SetPhase(PhasePlanning)
	// 清除旧计划（session.Plan 置 nil）
	a.session.Plan = nil

	// 向规划历史追加用户反馈，让 LLM 修改计划
	feedbackMsg := fmt.Sprintf("用户对上述计划有如下反馈，请根据反馈调整计划并重新输出 ```plan 块：\n\n%s", feedback)
	a.planMessages = append(a.planMessages, provider.Message{
		Role:    provider.RoleUser,
		Content: feedbackMsg,
	})

	if onUpdate != nil {
		onUpdate()
	}
	return a.planLoop(ctx, onUpdate)
}

// planLoop 是规划阶段的核心循环：LLM 调用只读工具 → 解析 plan 块 → 设置 PhaseWaitingConfirm。
func (a *TwoPhaseAgent) planLoop(ctx context.Context, onUpdate func()) error {
	for step := 0; step < a.maxPlanSteps; step++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// 检查是否应该使用流式输出（当前没有工具调用历史时，可能是纯文本回复）
		shouldStream := len(a.planMessages) > 1 && !a.hasRecentToolCalls()

		var rawContent string
		var err error

		if shouldStream {
			// 使用流式输出
			rawContent, err = a.streamPlanResponse(ctx, onUpdate)
		} else {
			// 使用非流式输出（需要解析结构化内容）
			result, completeErr := a.provider.Complete(ctx, a.planMessages)
			err = completeErr
			if err == nil {
				rawContent = result.Content
			}
		}

		if err != nil {
			return err
		}

		// 检查是否输出了 plan 块
		if parsed, ok := tools.ParsePlan(rawContent); ok {
			plan := a.buildExecutionPlan(parsed)

			// 验证计划的有效性
			if errors := a.validatePlan(plan); len(errors) > 0 {
				errMsg := "❌ 计划包含以下错误，请修正：\n\n"
				for i, err := range errors {
					errMsg += fmt.Sprintf("%d. %s\n", i+1, err)
				}
				errMsg += "\n请重新生成正确的执行计划。"

				a.session.AddSystemNote("计划验证失败")
				a.planMessages = append(a.planMessages, provider.Message{
					Role:    provider.RoleUser,
					Content: errMsg,
				})

				if onUpdate != nil {
					onUpdate()
				}
				continue // 继续规划循环，让 AI 修正计划
			}

			a.session.SetPlan(plan)

			// 展示 plan 块之外的自然语言说明（含提示用户输入 Y/N 的文字）
			displayText := tools.StripPlanBlock(rawContent)
			if displayText == "" {
				displayText = a.defaultConfirmPrompt(plan)
			}
			a.session.AddAssistantMessage(displayText)

			// 保留 LLM 这轮完整回复到规划历史（供 replan 时参考）
			a.planMessages = append(a.planMessages, provider.Message{
				Role:    provider.RoleAssistant,
				Content: rawContent,
			})

			a.session.SetPhase(PhaseWaitingConfirm)
			if onUpdate != nil {
				onUpdate()
			}
			return nil
		}

		// 解析工具调用
		toolCalls := tools.ParseToolCalls(rawContent)
		displayText := tools.StripToolBlocks(rawContent)

		// 把 LLM 的完整回复加入规划历史
		a.planMessages = append(a.planMessages, provider.Message{
			Role:    provider.RoleAssistant,
			Content: rawContent,
		})

		if displayText != "" {
			a.session.AddAssistantMessage(displayText)
		}
		if onUpdate != nil {
			onUpdate()
		}

		if len(toolCalls) == 0 {
			// LLM 既没有调用工具也没有输出计划，给一个提示继续
			hint := provider.Message{
				Role:    provider.RoleUser,
				Content: "请继续分析。收集到足够信息后，输出 ```plan 块。",
			}
			a.planMessages = append(a.planMessages, hint)
			continue
		}

		// 执行只读工具调用，结果反馈给 LLM
		for _, call := range toolCalls {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			// 检测重复调用，防止无限循环
			callKey := fmt.Sprintf("%s:%v", call.Name, call.Params)
			a.toolCallHistory[callKey]++

			if a.toolCallHistory[callKey] > 3 {
				// 同一个工具调用超过 3 次，给出警告
				warnMsg := fmt.Sprintf(
					"⚠️ 警告：工具 %s 已被调用 %d 次（参数相同）。\n"+
						"请避免重复调用相同的工具。如果已收集足够信息，请直接输出 ```plan 块。",
					call.Name, a.toolCallHistory[callKey])

				a.session.AddSystemNote(warnMsg)
				a.planMessages = append(a.planMessages, provider.Message{
					Role:    provider.RoleUser,
					Content: "[系统] " + warnMsg,
				})

				if onUpdate != nil {
					onUpdate()
				}
				continue // 跳过这次调用
			}

			tool, ok := a.readRegistry.Get(call.Name)
			if !ok {
				errMsg := a.tr.AgentToolNotAllowedInPlanning(call.Name)
				a.session.AddSystemNote(errMsg)
				a.planMessages = append(a.planMessages, provider.Message{
					Role:    provider.RoleUser,
					Content: fmt.Sprintf("[系统] %s", errMsg),
				})
				if onUpdate != nil {
					onUpdate()
				}
				continue
			}

			a.session.AddToolCall(call)
			if onUpdate != nil {
				onUpdate()
			}

			toolResult := tool.Execute(ctx, call)
			a.session.AddToolResult(toolResult, call.Name)
			a.planMessages = append(a.planMessages, provider.Message{
				Role:    provider.RoleUser,
				Content: fmt.Sprintf("[工具结果 %s]\n%s", call.Name, toolResult.Output),
			})
			if onUpdate != nil {
				onUpdate()
			}
		}
	}

	return fmt.Errorf("规划阶段超过最大步数 (%d)，未能生成执行计划", a.maxPlanSteps)
}

// execute 执行阶段二：按计划逐步执行写操作。
func (a *TwoPhaseAgent) execute(
	ctx context.Context,
	plan *ExecutionPlan,
	onUpdate func(),
) error {
	a.session.SetPhase(PhaseExecuting)
	if onUpdate != nil {
		onUpdate()
	}

	for _, step := range plan.Steps {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		a.session.UpdateStepStatus(step.ID, StepRunning, "", "")
		if onUpdate != nil {
			onUpdate()
		}

		toolName := step.ToolName

		// 检查并应用工具名别名映射
		if alias, ok := toolAliases[toolName]; ok {
			toolName = alias
		}

		tool, ok := a.fullRegistry.Get(toolName)
		if !ok {
			errMsg := a.formatToolNotFoundError(step.ToolName, toolName)
			a.session.UpdateStepStatus(step.ID, StepFailed, "", errMsg)
			if onUpdate != nil {
				onUpdate()
			}
			if step.Critical {
				return fmt.Errorf(a.tr.AgentCriticalStepFailed(step.Description, errMsg))
			}
			continue
		}

		// 为每个步骤创建带超时的 context
		stepCtx, cancel := context.WithTimeout(ctx, a.stepTimeout)

		call := tools.ToolCall{
			ID:     fmt.Sprintf("exec_%s", step.ID),
			Name:   toolName, // 使用映射后的工具名
			Params: step.Params,
		}

		// 在 goroutine 中执行工具，支持超时
		resultChan := make(chan tools.ToolResult, 1)
		go func() {
			resultChan <- tool.Execute(stepCtx, call)
		}()

		var result tools.ToolResult
		select {
		case result = <-resultChan:
			cancel() // 正常完成，取消 context
		case <-stepCtx.Done():
			cancel()
			if stepCtx.Err() == context.DeadlineExceeded {
				errMsg := a.formatTimeoutError(step, a.stepTimeout)
				a.session.UpdateStepStatus(step.ID, StepFailed, "", errMsg)
				if onUpdate != nil {
					onUpdate()
				}
				if step.Critical {
					return fmt.Errorf("关键步骤超时: %s", step.Description)
				}
				continue
			}
			return stepCtx.Err()
		}

		if result.Success {
			a.session.UpdateStepStatus(step.ID, StepDone, result.Output, "")
		} else {
			errMsg := a.formatExecutionError(step, result.Output)
			a.session.UpdateStepStatus(step.ID, StepFailed, "", errMsg)
			if step.Critical {
				if onUpdate != nil {
					onUpdate()
				}
				return fmt.Errorf(a.tr.AgentCriticalStepFailed(step.Description, errMsg))
			}
		}
		if onUpdate != nil {
			onUpdate()
		}
	}

	a.session.SetPhase(PhaseDone)
	done := plan.DoneCount()
	total := len(plan.Steps)
	failed := len(plan.FailedSteps())
	summary := fmt.Sprintf("执行完成：%d/%d 步成功", done, total)
	if failed > 0 {
		summary += fmt.Sprintf("，%d 步失败", failed)
	}
	a.session.AddSystemNote(summary)
	if onUpdate != nil {
		onUpdate()
	}
	return nil
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
	sb.WriteString("\n输入 **Y** 确认执行，**N** 取消，或直接输入补充说明来调整计划。")
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

// hasRecentToolCalls 检查最近是否有工具调用（用于判断是否应该流式输出）。
// 如果最近的消息包含工具调用，说明可能需要解析结构化内容，不适合流式。
func (a *TwoPhaseAgent) hasRecentToolCalls() bool {
	if len(a.planMessages) < 2 {
		return false
	}
	// 检查最后一条消息是否包含工具调用标记
	lastMsg := a.planMessages[len(a.planMessages)-1]
	return strings.Contains(lastMsg.Content, "[工具结果")
}

// streamPlanResponse 使用流式输出获取 LLM 响应。
// 适用于纯文本回复场景（无需解析 plan/tool 块）。
func (a *TwoPhaseAgent) streamPlanResponse(ctx context.Context, onUpdate func()) (string, error) {
	a.session.StartStreamingMessage()

	// 节流控制：避免过于频繁的 UI 刷新
	lastUpdateTime := time.Now()
	updateInterval := 50 * time.Millisecond // 每 50ms 最多刷新一次

	var fullContent strings.Builder

	err := a.provider.CompleteStream(ctx, a.planMessages, func(chunk string) {
		fullContent.WriteString(chunk)
		a.session.AppendToStreamingMessage(chunk)

		// 节流刷新
		now := time.Now()
		if now.Sub(lastUpdateTime) >= updateInterval {
			if onUpdate != nil {
				onUpdate()
			}
			lastUpdateTime = now
		}
	})

	// 完成流式输出
	a.session.FinishStreamingMessage()

	// 最后一次刷新，确保显示完整内容
	if onUpdate != nil {
		onUpdate()
	}

	if err != nil {
		return "", err
	}

	return fullContent.String(), nil
}

