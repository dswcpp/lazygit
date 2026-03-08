# 中文硬编码问题分析与解决方案

## 问题概述

项目中存在 **705 个唯一的硬编码中文字符串**，分布在多个模块中，严重违反了国际化（i18n）原则。

### 分布统计

| 模块 | 字符串数量 | 占比 | 严重程度 |
|------|-----------|------|---------|
| **GUI** | 366 | 51.9% | 🔴 高 |
| **AI Tools** | 181 | 25.7% | 🔴 高 |
| **AI Agent** | 96 | 13.6% | 🟡 中 |
| **AI Skills** | 43 | 6.1% | 🟡 中 |
| **Other** | 19 | 2.7% | 🟢 低 |

## 核心问题

### 1. **破坏项目国际化底线**
- lazygit 是国际化项目，已有完善的 i18n 机制（支持 9+ 语言）
- AI 模块的硬编码中文使非中文用户无法使用这些功能
- 违反了项目的编码规范和架构原则

### 2. **维护成本高**
- 同一文本在多处重复（如 "取消" 出现 7 次）
- 修改文案需要搜索所有文件
- 容易出现不一致的翻译

### 3. **代码可读性差**
- 中文字符串混杂在代码逻辑中
- 难以区分业务逻辑和展示文本
- 代码审查困难

## 解决方案

### 阶段 1: 扩展 i18n 系统（1-2 天）

#### 1.1 扩展 TranslationSet 结构

在 `pkg/i18n/english.go` 中添加 AI 相关字段：

```go
type TranslationSet struct {
    // ... 现有字段 ...

    // AI Common
    AICancel                     string
    AIConfirm                    string
    AISuccess                    string
    AIFailed                     string
    AIUnknown                    string
    AIYes                        string
    AINo                         string
    AIThinking                   string
    AIIdle                       string
    AIExecuting                  string
    AICancelled                  string

    // AI Agent
    AIAgentPlanningPhase         string
    AIAgentExecutionPhase        string
    AIAgentStepTimeout           string
    AIAgentCriticalStepFailed    string
    AIAgentUserRejected          string
    AIAgentConflictDetected      string
    AIAgentResolveConflict       string

    // AI Tools
    AIToolMissingParam           string  // "缺少 %s 参数"
    AIToolInvalidParam           string
    AIToolExecutionFailed        string
    AIToolNoChanges              string
    AIToolFilePath               string
    AIToolBranchName             string
    AIToolTagName                string
    AIToolCommitMessage          string

    // AI Skills
    AISkillCommitMsgPrompt       string
    AISkillBranchNamePrompt      string
    AISkillPRDescPrompt          string
    AISkillShellCmdPrompt        string

    // AI Chat
    AIChatNotEnabled             string
    AIChatInputPlaceholder       string
    AIChatGeneratingPlan         string
    AIChatCanInputNext           string

    // ... 更多字段 ...
}
```

#### 1.2 创建英文翻译

在 `pkg/i18n/english.go` 的 `EnglishTranslationSet()` 函数中添加：

```go
func EnglishTranslationSet() *TranslationSet {
    return &TranslationSet{
        // ... 现有翻译 ...

        // AI Common
        AICancel:                  "Cancel",
        AIConfirm:                 "Confirm",
        AISuccess:                 "Success",
        AIFailed:                  "Failed",
        AIUnknown:                 "Unknown",
        AIYes:                     "Yes",
        AINo:                      "No",
        AIThinking:                "Thinking",
        AIIdle:                    "Idle",
        AIExecuting:               "Executing",
        AICancelled:               "Cancelled",

        // AI Agent
        AIAgentPlanningPhase:      "Planning phase",
        AIAgentExecutionPhase:     "Execution phase",
        AIAgentStepTimeout:        "⏱️ Step execution timeout (%v): %s",
        AIAgentCriticalStepFailed: "Critical step failed: %s — %s",
        AIAgentUserRejected:       "[User rejected] Tool %s was not executed, please adjust subsequent operations accordingly.",
        AIAgentConflictDetected:   "Conflict detected",
        AIAgentResolveConflict:    "Resolve conflict manually and continue",

        // AI Tools
        AIToolMissingParam:        "Missing %s parameter",
        AIToolInvalidParam:        "Invalid %s parameter",
        AIToolExecutionFailed:     "Tool execution failed: %v",
        AIToolNoChanges:           "No changes",
        AIToolFilePath:            "File path",
        AIToolBranchName:          "Branch name",
        AIToolTagName:             "Tag name",
        AIToolCommitMessage:       "Commit message",

        // ... 更多翻译 ...
    }
}
```

#### 1.3 更新中文翻译文件

在 `pkg/i18n/translations/zh-CN.json` 中添加对应的中文翻译。

### 阶段 2: 创建 AI 模块的 i18n 辅助层（半天）

创建 `pkg/ai/i18n.go`：

```go
package ai

import "github.com/dswcpp/lazygit/pkg/i18n"

// Translator 为 AI 模块提供翻译功能
type Translator struct {
    tr *i18n.TranslationSet
}

// NewTranslator 创建翻译器
func NewTranslator(tr *i18n.TranslationSet) *Translator {
    return &Translator{tr: tr}
}

// Common translations
func (t *Translator) Cancel() string { return t.tr.AICancel }
func (t *Translator) Confirm() string { return t.tr.AIConfirm }
func (t *Translator) Success() string { return t.tr.AISuccess }
func (t *Translator) Failed() string { return t.tr.AIFailed }
func (t *Translator) Yes() string { return t.tr.AIYes }
func (t *Translator) No() string { return t.tr.AINo }

// Agent translations
func (t *Translator) AgentStepTimeout(duration, step string) string {
    return fmt.Sprintf(t.tr.AIAgentStepTimeout, duration, step)
}

func (t *Translator) AgentCriticalStepFailed(step, reason string) string {
    return fmt.Sprintf(t.tr.AIAgentCriticalStepFailed, step, reason)
}

// Tool translations
func (t *Translator) ToolMissingParam(param string) string {
    return fmt.Sprintf(t.tr.AIToolMissingParam, param)
}

// ... 更多辅助方法 ...
```

### 阶段 3: 批量替换硬编码字符串（3-5 天）

#### 3.1 优先级排序

1. **高优先级**（用户直接可见）
   - GUI 模块（366 个字符串）
   - AI Chat 相关提示

2. **中优先级**（错误消息）
   - AI Tools 错误提示（181 个）
   - AI Agent 状态消息（96 个）

3. **低优先级**（内部提示词）
   - AI Skills 的 prompt（43 个）
   - 这些主要是给 AI 模型看的，可以保持中文或改为英文

#### 3.2 替换策略

**示例 1: 简单字符串替换**

```go
// 修改前
return fmt.Errorf("缺少 name 参数")

// 修改后
return fmt.Errorf(t.tr.AIToolMissingParam, "name")
```

**示例 2: 复杂字符串替换**

```go
// 修改前
const planningSystemPrompt = `你是 lazygit 内置 AI，负责分析用户需求并制定 Git 操作计划。

## 工作流程

1. 调用只读工具（get_status、get_diff 等）收集必要信息
...`

// 修改后
func (t *Translator) PlanningSystemPrompt() string {
    return t.tr.AIAgentPlanningSystemPrompt
}
```

#### 3.3 自动化脚本

创建 `scripts/migrate_i18n.go` 来辅助批量替换：

```go
// 1. 读取 chinese_strings.json
// 2. 为每个字符串生成翻译键名
// 3. 生成 TranslationSet 字段定义
// 4. 生成英文翻译
// 5. 生成替换建议
```

### 阶段 4: 测试与验证（1-2 天）

1. **功能测试**
   - 切换到英文语言，验证所有 AI 功能正常
   - 切换到中文语言，验证翻译正确

2. **回归测试**
   - 运行现有测试套件
   - 手动测试关键 AI 功能

3. **翻译质量检查**
   - 确保翻译准确、一致
   - 检查格式化字符串的参数顺序

## 实施计划

### Week 1: 基础设施
- [ ] Day 1-2: 扩展 TranslationSet，添加所有 AI 相关字段
- [ ] Day 3: 创建 AI 模块的 i18n 辅助层
- [ ] Day 4: 编写自动化迁移脚本
- [ ] Day 5: 迁移 GUI 模块（高优先级）

### Week 2: 核心迁移
- [ ] Day 1-2: 迁移 AI Tools 模块
- [ ] Day 3: 迁移 AI Agent 模块
- [ ] Day 4: 迁移 AI Skills 模块
- [ ] Day 5: 测试与修复

### Week 3: 完善与发布
- [ ] Day 1-2: 完善翻译质量
- [ ] Day 3: 添加其他语言翻译（日语、韩语等）
- [ ] Day 4-5: 全面测试与文档更新

## 预期收益

1. **国际化支持**
   - 非中文用户可以使用 AI 功能
   - 支持 9+ 语言

2. **代码质量提升**
   - 代码逻辑与展示文本分离
   - 更易维护和审查

3. **用户体验改善**
   - 一致的翻译体验
   - 更专业的产品形象

4. **长期维护成本降低**
   - 集中管理翻译
   - 易于添加新语言

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| 翻译不准确 | 中 | 请母语者审核 |
| 遗漏字符串 | 低 | 使用自动化脚本检查 |
| 破坏现有功能 | 高 | 充分测试，分阶段发布 |
| 工作量大 | 中 | 使用自动化工具，分优先级实施 |

## 下一步行动

1. **立即行动**：审查并确认 i18n 架构设计
2. **本周内**：完成 TranslationSet 扩展和辅助层
3. **两周内**：完成高优先级模块迁移
4. **一个月内**：完成所有模块迁移并发布

## 附录

### A. 高频字符串列表

需要优先处理的高频字符串：

| 字符串 | 出现次数 | 建议翻译键 |
|--------|---------|-----------|
| 确定 | 8 | AIConfirm |
| 取消 | 7 | AICancel |
| 可输入下一条指令 | 7 | AIChatCanInputNext |
| 缺少 name 参数 | 6 | AIToolMissingParam |
| 缺少 path 参数 | 4 | AIToolMissingParam |
| 失败 | 4 | AIFailed |
| 成功 | 3 | AISuccess |

### B. 特殊处理项

**AI Prompt 字符串**：
- AI Skills 中的 prompt 模板（43 个）
- 这些是给 AI 模型看的，建议：
  1. 保持中文（如果主要用户是中文）
  2. 或改为英文（更通用）
  3. 或根据用户语言动态生成

**测试文件**：
- 测试文件中的中文字符串可以保留
- 或使用 i18n 以确保测试覆盖翻译功能

### C. 工具支持

建议开发的辅助工具：

1. **i18n 键名生成器**：根据中文自动生成合适的英文键名
2. **翻译覆盖率检查器**：确保所有硬编码字符串都已迁移
3. **翻译一致性检查器**：检查相同含义的字符串是否使用相同的翻译键
