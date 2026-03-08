# 修复报告：AI Agent "未知工具: commit_msg" 错误

## 问题描述

**错误信息**:
```
[失败] 基于暂存区变更生成符合规范的提交信息
错误: 未知工具: commit_msg
```

**发生时机**: 用户确认执行计划后，在执行阶段

## 根本原因分析

### 问题根源

`commit_msg` 是一个 **SkillTool**，它的设计用途和限制如下：

1. **只在规划阶段可用**
   - `commit_msg` 注册在 `readRegistry`（只读工具注册表）
   - 规划阶段的 LLM 可以调用它来生成提交信息
   - 返回的提交信息应该作为参数传递给执行阶段的 `commit` 工具

2. **执行阶段不可用**
   - 执行阶段使用 `fullRegistry`（完整工具注册表）
   - `fullRegistry` 中没有 `commit_msg`，因为它是辅助工具，不是执行工具

### 错误流程

```
规划阶段:
  AI 调用 get_staged_diff → 获取 diff
  AI 调用 commit_msg → 生成提交信息
  AI 生成计划:
    步骤 1: stage_all
    步骤 2: commit_msg  ← 错误！不应该放在执行计划中

执行阶段:
  执行步骤 1: stage_all ✅
  执行步骤 2: commit_msg ❌ 错误：未知工具
```

### 正确流程

```
规划阶段:
  AI 调用 get_staged_diff → 获取 diff
  AI 调用 commit_msg → 返回 "feat: 添加用户登录"
  AI 生成计划:
    步骤 1: stage_all
    步骤 2: commit, params: {"message": "feat: 添加用户登录"}

执行阶段:
  执行步骤 1: stage_all ✅
  执行步骤 2: commit ✅
```

## 解决方案

### 修改 System Prompt

**文件**: `pkg/ai/agent/two_phase_agent.go`

**修改内容**: 在 `planningSystemPrompt` 中明确说明 `commit_msg` 和 `branch_name` 的使用方式：

```go
## 工作流程

1. 调用只读工具（get_status、get_diff 等）收集必要信息
2. **如需生成提交信息**：
   - 先调用 get_staged_diff 获取暂存区 diff
   - 然后调用 commit_msg 工具生成提交信息（返回的内容直接用作 commit 的 message 参数）
   - **重要**：commit_msg 只能在规划阶段调用，不能放入执行计划
3. **如需生成分支名**：
   - 调用 branch_name 工具生成分支名
   - **重要**：branch_name 只能在规划阶段调用，不能放入执行计划

## 特殊工具说明

**commit_msg 和 branch_name 是辅助工具，只能在规划阶段调用**：
- 在规划阶段调用 commit_msg 获取提交信息
- 将返回的提交信息作为 commit 工具的 message 参数
- **不要**把 commit_msg 放入执行计划的 steps 中

示例：
```tool
{"name": "commit_msg", "params": {"diff": "..."}}
```
返回: "feat: 添加用户登录功能"

然后在执行计划中：
```plan
{
  "steps": [
    {"tool": "commit", "params": {"message": "feat: 添加用户登录功能"}}
  ]
}
```
```

## 为什么这样设计？

### 两阶段架构的优势

1. **规划阶段**（Planning Phase）
   - 使用只读工具收集信息
   - 使用 SkillTool（如 commit_msg）生成参数
   - 生成完整的执行计划
   - **优点**：可以多次调用 LLM，反复调整计划

2. **执行阶段**（Execution Phase）
   - 按计划逐步执行写操作
   - 不再调用 LLM
   - **优点**：快速、确定性、可回滚

### SkillTool 的作用

SkillTool 是一个桥梁，让规划阶段的 LLM 可以：
- 调用需要 AI 能力的辅助功能（生成提交信息、分支名等）
- 将生成的结果嵌入到执行计划的参数中
- 避免在执行阶段再次调用 LLM（提升性能和可靠性）

## 相关概念对比

| 概念 | 用途 | 注册位置 | 可用阶段 |
|------|------|---------|---------|
| **Tool** | 执行 Git 操作 | fullRegistry | 执行阶段 |
| **SkillTool** | 生成参数（提交信息、分支名等） | readRegistry | 规划阶段 |
| **Skill** | AI 能力封装 | skillMap | 通过 SkillTool 间接调用 |

## 测试验证

### 场景 1：提交当前修改

**用户输入**: "帮我提交当前修改"

**预期流程**:

```
规划阶段:
  1. 调用 get_status → 查看状态
  2. 调用 get_staged_diff → 获取 diff
  3. 调用 commit_msg(diff="...") → 返回 "feat: 添加登录功能"
  4. 生成计划:
     ```plan
     {
       "steps": [
         {"id": "1", "tool": "stage_all", "params": {}},
         {"id": "2", "tool": "commit", "params": {"message": "feat: 添加登录功能"}}
       ]
     }
     ```

执行阶段:
  1. 执行 stage_all ✅
  2. 执行 commit(message="feat: 添加登录功能") ✅
```

### 场景 2：创建分支并提交

**用户输入**: "创建 feature 分支并提交当前修改"

**预期流程**:

```
规划阶段:
  1. 调用 branch_name(description="feature") → 返回 "feature/user-login"
  2. 调用 get_staged_diff → 获取 diff
  3. 调用 commit_msg(diff="...") → 返回 "feat: 添加用户登录"
  4. 生成计划:
     ```plan
     {
       "steps": [
         {"id": "1", "tool": "create_branch", "params": {"name": "feature/user-login", "checkout": true}},
         {"id": "2", "tool": "stage_all", "params": {}},
         {"id": "3", "tool": "commit", "params": {"message": "feat: 添加用户登录"}}
       ]
     }
     ```

执行阶段:
  1. 执行 create_branch ✅
  2. 执行 stage_all ✅
  3. 执行 commit ✅
```

## 常见错误模式

### ❌ 错误 1：把 commit_msg 放入执行计划

```json
{
  "steps": [
    {"tool": "stage_all"},
    {"tool": "commit_msg", "params": {"diff": "..."}},  // 错误！
    {"tool": "commit", "params": {"message": "..."}}
  ]
}
```

**问题**: 执行阶段找不到 commit_msg 工具

**修复**: 在规划阶段调用 commit_msg，将结果放入 commit 的参数

### ❌ 错误 2：commit 缺少 message 参数

```json
{
  "steps": [
    {"tool": "stage_all"},
    {"tool": "commit", "params": {}}  // 错误！缺少 message
  ]
}
```

**问题**: commit 工具需要 message 参数

**修复**: 先调用 commit_msg 生成提交信息，然后传递给 commit

### ✅ 正确模式

```
规划阶段:
  调用 commit_msg → 获取提交信息

执行计划:
  {
    "steps": [
      {"tool": "stage_all"},
      {"tool": "commit", "params": {"message": "从 commit_msg 获取的信息"}}
    ]
  }
```

## 修复效果

### 修复前

```
用户: "帮我提交当前修改"
AI 生成计划:
  步骤 1: stage_all
  步骤 2: commit_msg  ← 错误
执行: ❌ 错误: 未知工具: commit_msg
```

### 修复后

```
用户: "帮我提交当前修改"
AI 规划阶段:
  调用 commit_msg → "feat: 添加功能"
AI 生成计划:
  步骤 1: stage_all
  步骤 2: commit(message="feat: 添加功能")
执行: ✅ 成功
```

## 相关文件

### 修改的文件
- `pkg/ai/agent/two_phase_agent.go` - 改进 system prompt

### 相关文件
- `pkg/ai/tools/skill_tool.go` - SkillTool 实现
- `pkg/ai/manager.go` - DefaultSkillTools 定义
- `pkg/ai/skills/commit_msg.go` - CommitMsgSkill 实现

## 后续建议

### 短期改进

1. **添加更多示例**
   - 在 system prompt 中添加更多正确用法的示例
   - 强调 commit_msg 和 branch_name 的特殊性

2. **改进错误提示**
   - 当执行阶段遇到 commit_msg 时，给出更友好的错误提示
   - 提示用户这是规划阶段工具，不应该在执行计划中

### 中期改进

3. **计划验证**
   - 在生成计划后，验证所有工具名是否在 fullRegistry 中
   - 提前发现并修正错误

4. **自动修正**
   - 检测到 commit_msg 在执行计划中时，自动提示 AI 重新规划

### 长期改进

5. **架构优化**
   - 考虑统一 Tool 和 SkillTool 的接口
   - 或者在执行阶段也支持 SkillTool（但需要考虑性能）

## 总结

通过**改进 System Prompt**，明确说明 `commit_msg` 和 `branch_name` 的使用方式：

1. ✅ 明确这些是规划阶段工具，不能放入执行计划
2. ✅ 提供了清晰的使用示例
3. ✅ 强调了正确的工作流程

**修复状态**: ✅ 已完成

**建议**: 用户现在可以重新尝试提交操作，AI 应该能正确生成计划。
