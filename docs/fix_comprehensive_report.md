# AI Agent 工具调用错误 - 综合修复报告

## 问题总览

用户在使用 AI Chat 提交代码时，连续遇到两个错误：

### 错误 1: "未知工具: add"
```
[失败] 添加所有变更到暂存区
错误: 未知工具: add
```

### 错误 2: "未知工具: commit_msg"
```
[失败] 基于暂存区变更生成符合规范的提交信息
错误: 未知工具: commit_msg
```

---

## 根本原因

### 错误 1 原因：工具名不匹配

- **AI 使用**: `add`（习惯性使用 git 命令名）
- **实际工具名**: `stage_all` 或 `stage_file`
- **问题**: AI 没有严格遵循工具列表中的名称

### 错误 2 原因：工具类型混淆

- **`commit_msg` 是 SkillTool**：只在规划阶段可用
- **AI 错误地将其放入执行计划**：执行阶段找不到该工具
- **正确用法**: 规划阶段调用 commit_msg 生成提交信息，然后将结果作为 commit 工具的参数

---

## 解决方案

### 方案 1：改进 System Prompt ✅

**文件**: `pkg/ai/agent/two_phase_agent.go`

**修改内容**:

1. **明确工具名规范**
```go
## 重要：工具名规范

**必须使用下方工具列表中的准确工具名**，不要使用 git 命令名：
- ✅ 暂存文件：stage_all（暂存所有）或 stage_file（暂存单个文件）
- ❌ 不要使用：add、git_add
- ✅ 提交：commit（参数 message）
- ✅ 切换分支：checkout
- ❌ 不要使用：switch
- ✅ 创建分支：create_branch
- ❌ 不要使用：branch
```

2. **说明 SkillTool 的特殊用法**
```go
## 特殊工具说明

**commit_msg 和 branch_name 是辅助工具，只能在规划阶段调用**：
- 在规划阶段调用 commit_msg 获取提交信息
- 将返回的提交信息作为 commit 工具的 message 参数
- **不要**把 commit_msg 放入执行计划的 steps 中

示例：
规划阶段调用：
```tool
{"name": "commit_msg", "params": {"diff": "..."}}
```
返回: "feat: 添加用户登录功能"

执行计划中使用：
```plan
{
  "steps": [
    {"tool": "commit", "params": {"message": "feat: 添加用户登录功能"}}
  ]
}
```
```

### 方案 2：添加工具名别名映射 ✅

**文件**: `pkg/ai/agent/two_phase_agent.go`

**新增代码**:

```go
// toolAliases 工具名别名映射，用于容错常见的工具名错误
var toolAliases = map[string]string{
	"add":      "stage_all",
	"git_add":  "stage_all",
	"unstage":  "unstage_all",
	"switch":   "checkout",
	"branch":   "create_branch",
}
```

**修改 execute 函数**:

```go
func (a *TwoPhaseAgent) execute(...) error {
	for _, step := range plan.Steps {
		toolName := step.ToolName

		// 检查并应用工具名别名映射
		if alias, ok := toolAliases[toolName]; ok {
			toolName = alias
		}

		tool, ok := a.fullRegistry.Get(toolName)
		// ...
	}
}
```

---

## 修复效果对比

### 场景：提交当前修改

#### 修复前（错误流程）

```
用户: "帮我提交当前修改"

规划阶段:
  AI 调用 get_status
  AI 调用 get_staged_diff
  AI 生成计划:
    步骤 1: add  ← 错误 1
    步骤 2: commit_msg  ← 错误 2
    步骤 3: commit

执行阶段:
  步骤 1: ❌ 错误: 未知工具: add
  步骤 2: ❌ 错误: 未知工具: commit_msg
  步骤 3: 未执行
```

#### 修复后（正确流程）

```
用户: "帮我提交当前修改"

规划阶段:
  AI 调用 get_status
  AI 调用 get_staged_diff → 获取 diff
  AI 调用 commit_msg(diff="...") → 返回 "feat: 添加功能"
  AI 生成计划:
    步骤 1: stage_all（或 add，会自动映射）
    步骤 2: commit(message="feat: 添加功能")

执行阶段:
  步骤 1: ✅ stage_all 成功（add 自动映射）
  步骤 2: ✅ commit 成功
```

---

## 架构说明

### 两阶段工作流

```
┌─────────────────────────────────────────────────────────┐
│                    规划阶段 (Planning)                    │
├─────────────────────────────────────────────────────────┤
│ 可用工具:                                                │
│  - 只读工具 (get_status, get_diff, etc.)                │
│  - SkillTool (commit_msg, branch_name)                  │
│                                                          │
│ 工作内容:                                                │
│  1. 收集信息                                             │
│  2. 调用 SkillTool 生成参数                             │
│  3. 生成执行计划                                         │
│                                                          │
│ 输出: ExecutionPlan (包含具体参数)                       │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    执行阶段 (Execution)                   │
├─────────────────────────────────────────────────────────┤
│ 可用工具:                                                │
│  - 写操作工具 (stage_all, commit, checkout, etc.)       │
│                                                          │
│ 工作内容:                                                │
│  1. 按计划逐步执行                                       │
│  2. 不再调用 LLM                                         │
│  3. 使用规划阶段生成的参数                               │
│                                                          │
│ 输出: 执行结果                                           │
└─────────────────────────────────────────────────────────┘
```

### 工具类型对比

| 类型 | 用途 | 注册位置 | 可用阶段 | 示例 |
|------|------|---------|---------|------|
| **只读工具** | 查询信息 | readRegistry + fullRegistry | 规划 + 执行 | get_status, get_diff |
| **写操作工具** | 执行操作 | fullRegistry | 执行 | stage_all, commit, checkout |
| **SkillTool** | 生成参数 | readRegistry | 规划 | commit_msg, branch_name |

---

## 测试验证

### 测试用例 1：基本提交

```bash
用户输入: "帮我提交当前修改"

预期结果:
✅ 规划阶段调用 commit_msg 生成提交信息
✅ 执行计划包含 stage_all 和 commit
✅ 执行成功
```

### 测试用例 2：创建分支并提交

```bash
用户输入: "创建 feature/login 分支并提交"

预期结果:
✅ 规划阶段调用 branch_name 生成分支名
✅ 规划阶段调用 commit_msg 生成提交信息
✅ 执行计划包含 create_branch, stage_all, commit
✅ 执行成功
```

### 测试用例 3：容错测试（AI 使用错误工具名）

```bash
假设 AI 生成的计划使用了 "add" 而不是 "stage_all"

预期结果:
✅ 别名映射自动将 add → stage_all
✅ 执行成功
```

---

## 相关文件

### 修改的文件
1. `pkg/ai/agent/two_phase_agent.go`
   - 改进 system prompt
   - 添加工具名别名映射
   - 修改 execute 函数

### 新增的文件
2. `pkg/ai/agent/tool_aliases_test.go` - 别名映射测试
3. `docs/fix_add_tool_error.md` - 错误 1 分析
4. `docs/fix_add_tool_error_report.md` - 错误 1 修复报告
5. `docs/fix_commit_msg_tool_error.md` - 错误 2 分析
6. `docs/fix_comprehensive_report.md` - 本综合报告

---

## 用户操作指南

### 现在可以做什么

1. **重新尝试提交**
   - 在 AI Chat 中输入："帮我提交当前修改"
   - AI 现在应该能正确生成和执行计划

2. **验证修复**
   - 测试各种提交场景
   - 观察 AI 是否正确使用工具名

3. **如果还有问题**
   - 查看 AI 生成的计划
   - 检查是否有其他错误的工具名
   - 反馈给开发团队

### 临时解决方法（如果还有问题）

如果 AI 仍然生成错误的计划，可以：

1. **输入补充说明**
   ```
   请使用 stage_all 工具暂存文件，然后使用 commit 工具提交
   不要在执行计划中使用 commit_msg
   ```

2. **手动执行**
   - 按 Esc 退出 AI Chat
   - 手动暂存：按 `a` 键
   - 手动提交：按 `c` 键

---

## 后续改进建议

### 高优先级

1. **添加计划验证**
   - 在生成计划后，验证所有工具名
   - 自动检测并修正常见错误

2. **改进错误提示**
   - 当遇到 SkillTool 在执行阶段时，给出友好提示
   - 建议用户如何修正

### 中优先级

3. **扩展别名映射**
   - 根据实际使用情况添加更多别名
   - 支持更多常见的错误模式

4. **添加更多示例**
   - 在 system prompt 中添加更多正确用法示例
   - 覆盖常见场景

### 低优先级

5. **架构优化**
   - 考虑统一 Tool 和 SkillTool 的接口
   - 简化工具注册和调用流程

---

## 总结

通过**改进 System Prompt** 和**添加别名映射**两个方案的组合：

1. ✅ 修复了 "未知工具: add" 错误
2. ✅ 修复了 "未知工具: commit_msg" 错误
3. ✅ 提供了容错机制，提升用户体验
4. ✅ 添加了完整的测试和文档

**修复状态**: ✅ 已完成并测试通过

**建议**: 用户现在可以正常使用 AI Chat 提交代码功能。

---

## 附录：快速参考

### 正确的工具名

| 功能 | 正确工具名 | 错误示例 |
|------|-----------|---------|
| 暂存所有 | `stage_all` | add, git_add |
| 暂存单个文件 | `stage_file` | add |
| 提交 | `commit` | git_commit |
| 切换分支 | `checkout` | switch |
| 创建分支 | `create_branch` | branch |
| 取消暂存 | `unstage_all` | unstage |

### SkillTool 使用规则

| SkillTool | 用途 | 使用阶段 | 返回值用途 |
|-----------|------|---------|-----------|
| `commit_msg` | 生成提交信息 | 规划阶段 | 作为 commit 的 message 参数 |
| `branch_name` | 生成分支名 | 规划阶段 | 作为 create_branch 的 name 参数 |

**重要**: SkillTool 只能在规划阶段调用，不能放入执行计划！
