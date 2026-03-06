package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/repocontext"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

const defaultMaxPlanSteps = 15

// planningSystemPrompt 是规划阶段专用的 system prompt。
const planningSystemPrompt = `你是 lazygit 内置 AI，负责分析用户需求并制定 Git 操作计划。

## 工作流程

1. 调用只读工具（get_status、get_diff 等）收集必要信息
2. 如需生成提交信息，调用 commit_msg 工具；如需生成分支名，调用 branch_name 工具
3. 信息收集完毕后，输出一个 ` + "```plan" + ` 块，内含完整执行计划
4. ` + "```plan" + ` 块之后附上一段简短的自然语言说明，提示用户可以输入 Y 确认、N 取消，或补充说明
5. 严禁在规划阶段调用任何写操作工具

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
	maxPlanSteps int
	// planMessages 保存规划阶段的完整消息历史（含工具调用结果）。
	// 在 replan 时直接追加用户反馈后继续循环，避免重复调用只读工具。
	planMessages []provider.Message
}

// NewTwoPhaseAgent 创建 TwoPhaseAgent。
//   - fullRegistry: 包含所有工具的注册表（执行阶段使用）
//   - readRegistry: 仅包含只读工具 + SkillTool 的注册表（规划阶段使用）
func NewTwoPhaseAgent(
	p provider.Provider,
	fullRegistry *tools.Registry,
	readRegistry *tools.Registry,
	session *Session,
) *TwoPhaseAgent {
	return &TwoPhaseAgent{
		provider:     p,
		fullRegistry: fullRegistry,
		readRegistry: readRegistry,
		session:      session,
		maxPlanSteps: defaultMaxPlanSteps,
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
	a.session.SetPhase(PhasePlanning)

	// 构建规划阶段 system prompt（含工具列表）
	sysPrompt := planningSystemPrompt
	if toolSection := a.readRegistry.SystemPromptSection(tools.PermReadOnly); toolSection != "" {
		sysPrompt += "\n\n" + toolSection
	}

	// 初始用户消息：仓库上下文 + 用户指令
	initMsg := fmt.Sprintf("## 当前仓库状态\n\n%s\n\n## 用户指令\n\n%s",
		repoCtx.CompactString(), userMsg)

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

		result, err := a.provider.Complete(ctx, a.planMessages)
		if err != nil {
			return err
		}

		rawContent := result.Content

		// 检查是否输出了 plan 块
		if parsed, ok := tools.ParsePlan(rawContent); ok {
			plan := a.buildExecutionPlan(parsed)
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

			tool, ok := a.readRegistry.Get(call.Name)
			if !ok {
				errMsg := fmt.Sprintf("规划阶段不允许调用工具: %s", call.Name)
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

		tool, ok := a.fullRegistry.Get(step.ToolName)
		if !ok {
			errMsg := fmt.Sprintf("未知工具: %s", step.ToolName)
			a.session.UpdateStepStatus(step.ID, StepFailed, "", errMsg)
			if onUpdate != nil {
				onUpdate()
			}
			if step.Critical {
				return fmt.Errorf("关键步骤失败: %s — %s", step.Description, errMsg)
			}
			continue
		}

		call := tools.ToolCall{
			ID:     fmt.Sprintf("exec_%s", step.ID),
			Name:   step.ToolName,
			Params: step.Params,
		}
		result := tool.Execute(ctx, call)

		if result.Success {
			a.session.UpdateStepStatus(step.ID, StepDone, result.Output, "")
		} else {
			a.session.UpdateStepStatus(step.ID, StepFailed, "", result.Output)
			if step.Critical {
				if onUpdate != nil {
					onUpdate()
				}
				return fmt.Errorf("关键步骤失败: %s — %s", step.Description, result.Output)
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

// buildExecutionPlan 将 tools.ParsedPlan 转换为 ExecutionPlan，
// 同时从工具注册表中查询每个步骤的权限级别。
func (a *TwoPhaseAgent) buildExecutionPlan(parsed tools.ParsedPlan) *ExecutionPlan {
	steps := make([]*PlanStep, 0, len(parsed.Steps))
	for _, s := range parsed.Steps {
		perm := tools.PermReadOnly
		if tool, ok := a.fullRegistry.Get(s.ToolName); ok {
			perm = tool.Schema().Permission
		}
		steps = append(steps, &PlanStep{
			ID:          s.ID,
			Description: s.Description,
			ToolName:    s.ToolName,
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
