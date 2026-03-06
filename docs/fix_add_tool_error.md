# 问题分析：AI Agent 调用了不存在的 `add` 工具

## 问题描述

用户在使用 AI Chat 提交代码时遇到错误：

```
[失败] 添加所有变更到暂存区(包括修改的文件和未追踪的新文件)
错误: 未知工具: add
```

## 根本原因

AI Agent 在生成执行计划时，使用了 `add` 作为工具名，但实际注册的工具名是：
- `stage_all` - 暂存所有变更
- `stage_file` - 暂存指定文件

**为什么会发生这个问题？**

1. **AI 的习惯性思维**：AI 熟悉 `git add` 命令，可能自然地使用 `add` 作为工具名
2. **工具名不够直观**：`stage_all` 虽然准确，但不如 `add` 直观
3. **System Prompt 可能不够强调**：虽然工具列表在 prompt 中，但 AI 可能没有严格遵循

## 解决方案

### 方案 1：改进 System Prompt（推荐）✅

在 `planningSystemPrompt` 中明确强调工具名：

```go
const planningSystemPrompt = `你是 lazygit 内置 AI，负责分析用户需求并制定 Git 操作计划。

## 工作流程

1. 调用只读工具（get_status、get_diff 等）收集必要信息
2. 如需生成提交信息，调用 commit_msg 工具；如需生成分支名，调用 branch_name 工具
3. 信息收集完毕后，输出一个 ` + "```plan" + ` 块，内含完整执行计划
4. ` + "```plan" + ` 块之后附上一段简短的自然语言说明，提示用户可以输入 Y 确认、N 取消，或补充说明
5. 严禁在规划阶段调用任何写操作工具

## 重要提示

**必须使用下方列出的准确工具名**，不要使用 git 命令名（如 add、commit）：
- 暂存文件使用 stage_all 或 stage_file，不是 add
- 提交使用 commit，参数是 message
- 切换分支使用 checkout，不是 switch

## 计划格式
...
`
```

### 方案 2：添加工具名别名映射

在执行阶段添加别名映射，自动转换常见错误：

```go
// 工具名别名映射
var toolAliases = map[string]string{
	"add":      "stage_all",
	"git_add":  "stage_all",
	"unstage":  "unstage_all",
	"switch":   "checkout",
	"branch":   "create_branch",
}

func (a *TwoPhaseAgent) execute(
	ctx context.Context,
	plan *ExecutionPlan,
	onUpdate func(),
) error {
	for _, step := range plan.Steps {
		toolName := step.ToolName

		// 检查别名
		if alias, ok := toolAliases[toolName]; ok {
			toolName = alias
		}

		tool, ok := a.fullRegistry.Get(toolName)
		// ...
	}
}
```

### 方案 3：注册 `add` 作为 `stage_all` 的别名工具

创建一个包装工具：

```go
// AddTool 是 StageAllTool 的别名，为了兼容 AI 的习惯
type AddTool struct{ d *Deps }

func NewAddTool(d *Deps) tools.Tool { return &AddTool{d} }

func (t *AddTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "add",
		Description: "暂存所有工作区变更（stage_all 的别名）",
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermWriteLocal,
	}
}

func (t *AddTool) Execute(ctx context.Context, call tools.ToolCall) tools.ToolResult {
	// 直接调用 StageAllTool
	return NewStageAllTool(t.d).Execute(ctx, call)
}
```

然后在 `register.go` 中注册：

```go
func RegisterAll(d *Deps, r *tools.Registry, p provider.Provider) {
	for _, t := range []tools.Tool{
		// ... 其他工具
		NewStageAllTool(d),
		NewAddTool(d), // 添加别名工具
		// ...
	} {
		r.Register(t)
	}
}
```

## 推荐实施顺序

1. **立即实施**：方案 1（改进 System Prompt）
   - 工作量：5 分钟
   - 效果：明确告诉 AI 使用正确的工具名
   - 风险：低

2. **短期实施**：方案 2（别名映射）
   - 工作量：30 分钟
   - 效果：自动容错，提升用户体验
   - 风险：低

3. **可选实施**：方案 3（注册别名工具）
   - 工作量：15 分钟
   - 效果：最彻底的解决方案
   - 风险：低，但会增加工具数量

## 临时解决方法（用户侧）

在 AI 生成计划后，如果看到错误的工具名，可以：

1. **输入补充说明**：
   ```
   请使用 stage_all 工具，不是 add
   ```

2. **手动执行**：
   - 按 Esc 退出 AI Chat
   - 手动暂存文件（按 `a` 键）
   - 手动提交（按 `c` 键）

## 测试验证

修复后，测试以下场景：

```bash
# 场景 1：提交所有变更
用户: "帮我提交当前修改"
预期: AI 使用 stage_all 工具

# 场景 2：提交指定文件
用户: "只提交 README.md"
预期: AI 使用 stage_file 工具

# 场景 3：创建分支并提交
用户: "创建 feature/test 分支并提交当前修改"
预期: AI 使用 create_branch 和 stage_all 工具
```

## 相关文件

- `pkg/ai/agent/two_phase_agent.go` - System Prompt 定义
- `pkg/ai/tools/git/staging.go` - 暂存工具实现
- `pkg/ai/tools/git/register.go` - 工具注册
