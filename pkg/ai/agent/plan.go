package agent

import "github.com/dswcpp/lazygit/pkg/ai/tools"

// AgentPhase 标识 TwoPhaseAgent 当前所处的阶段。
type AgentPhase int

const (
	PhasePlanning       AgentPhase = iota // 阶段一：分析与规划中
	PhaseWaitingConfirm                   // 等待用户在聊天中输入 Y/N 或补充说明
	PhaseExecuting                        // 阶段二：执行写操作中
	PhaseDone                             // 正常完成
	PhaseCancelled                        // 用户取消
)

// StepStatus 描述单个执行步骤的状态。
type StepStatus int

const (
	StepPending StepStatus = iota // 待执行
	StepRunning                   // 执行中
	StepDone                      // 成功完成
	StepFailed                    // 执行失败
	StepSkipped                   // 跳过（non-critical 步骤失败后继续）
)

// String 返回步骤状态的可读标签。
func (s StepStatus) String() string {
	switch s {
	case StepPending:
		return "待执行"
	case StepRunning:
		return "执行中"
	case StepDone:
		return "完成"
	case StepFailed:
		return "失败"
	case StepSkipped:
		return "跳过"
	default:
		return "未知"
	}
}

// PlanStep 是执行计划中的一个原子操作。
// 规划阶段由 LLM 填充 ID/Description/ToolName/Params/Critical；
// 执行阶段填充 Status/Result/Error。
type PlanStep struct {
	ID          string         // 步骤唯一标识，如 "1"、"2"
	Description string         // 展示给用户的人类可读描述
	ToolName    string         // 要调用的工具名
	Params      map[string]any // 预计算好的参数（规划阶段确定）
	Permission  tools.PermissionLevel
	// Critical=true：该步骤失败则中止整个执行。
	// Critical=false：失败后标记为 StepSkipped，继续后续步骤。
	Critical bool

	// 以下字段在执行阶段填充。
	Status StepStatus
	Result string // 工具返回的输出
	Error  string // 失败时的错误信息
}

// ExecutionPlan 是阶段一的输出，阶段二的输入。
type ExecutionPlan struct {
	Summary string // 整体描述，供用户确认时阅读
	Steps   []*PlanStep
}

// HasWriteOps 报告计划是否包含任何写操作（决定是否需要用户确认）。
func (p *ExecutionPlan) HasWriteOps() bool {
	for _, s := range p.Steps {
		if s.Permission >= tools.PermWriteLocal {
			return true
		}
	}
	return false
}

// DoneCount 返回已成功完成的步骤数。
func (p *ExecutionPlan) DoneCount() int {
	n := 0
	for _, s := range p.Steps {
		if s.Status == StepDone {
			n++
		}
	}
	return n
}

// FailedSteps 返回执行失败的步骤列表。
func (p *ExecutionPlan) FailedSteps() []*PlanStep {
	var out []*PlanStep
	for _, s := range p.Steps {
		if s.Status == StepFailed {
			out = append(out, s)
		}
	}
	return out
}
