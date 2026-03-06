# AI 智能变更分析工具

## 功能说明

`analyze_changes` 工具提供智能的文件变更分析功能，特别适用于大量文件变更的场景。它会：

1. **逐个文件分析**：避免一次性发送过大的 diff 内容给 AI
2. **智能整合**：将多个文件的分析结果整合成结构化报告
3. **灵活配置**：支持自定义分析重点（安全、性能、代码质量等）

## 使用方式

### 在 AI Chat 中使用

直接在 AI 聊天窗口中输入：

```
分析当前的变更
```

或者指定分析重点：

```
分析暂存区的变更，重点关注安全问题
```

### 工具参数

- `staged` (bool):
  - `true` - 分析暂存区变更
  - `false` - 分析工作区变更（默认）

- `focus` (string):
  - 分析重点，如："安全问题"、"性能优化"、"代码质量"
  - 留空则进行全面分析

## 输出示例

```markdown
# 变更分析报告（重点：安全问题）

**文件数**: 3 个（成功分析 3，失败 0）
**总变更行数**: 约 150 行

## 详细分析

### 1. pkg/gui/ai_chat.go

主要变更：添加了 AI 聊天输入处理逻辑
- 新增 aiChatInputEditor 函数处理键盘输入
- 实现了 Enter 发送、Esc 关闭等交互
- 建议：考虑添加输入内容的长度限制

### 2. pkg/gui/controllers/helpers/ai_chat_helper.go

主要变更：实现了聊天会话管理
- 新增 AIChatSession 结构体
- 实现消息渲染和滚动控制
- 潜在问题：未发现明显安全问题

### 3. pkg/gui/controllers/helpers/ai_chat_helper_test.go

主要变更：添加了单元测试
- 测试了滚动行为
- 测试了状态推导逻辑
- 覆盖率良好

## 整体建议

建议在提交前：
1. 确认所有变更符合预期
2. 运行测试确保功能正常
3. 检查是否有遗漏的文件
```

## 技术实现

### 架构设计

```
用户请求
    ↓
获取变更文件列表
    ↓
逐个文件分析 (并发控制)
    ↓
    ├─ 获取文件 diff
    ├─ 构建分析提示词
    ├─ 调用 AI 分析
    └─ 收集结果
    ↓
整合所有分析结果
    ↓
生成结构化报告
```

### 上下文优化策略

1. **分片处理**：每个文件独立分析，避免超出上下文限制
2. **简洁提示词**：针对单文件的分析提示词保持简洁
3. **结果缓存**：分析结果在内存中缓存，支持快速重新整合

### 错误处理

- 单个文件分析失败不会中断整体流程
- 失败的文件会在报告中标注错误信息
- 支持 context 取消，可随时中断分析

## 扩展性

### 添加新的分析维度

可以通过修改 `buildAnalysisPrompt` 函数来添加新的分析维度：

```go
func (t *AnalyzeChangesTool) buildAnalysisPrompt(path, diff, focus string) string {
    // 根据 focus 参数定制提示词
    switch focus {
    case "性能":
        return buildPerformancePrompt(path, diff)
    case "安全":
        return buildSecurityPrompt(path, diff)
    default:
        return buildGeneralPrompt(path, diff)
    }
}
```

### 支持更多文件类型

可以根据文件扩展名使用不同的分析策略：

```go
ext := filepath.Ext(path)
switch ext {
case ".go":
    prompt = buildGoAnalysisPrompt(diff)
case ".js", ".ts":
    prompt = buildJSAnalysisPrompt(diff)
default:
    prompt = buildGeneralPrompt(diff)
}
```

## 性能考虑

- **并发控制**：当前实现是串行分析，可以考虑添加并发控制
- **diff 截断**：对于超大文件，可以考虑只分析关键部分
- **缓存机制**：相同 diff 的分析结果可以缓存

## 测试

运行测试：

```bash
go test ./pkg/ai/tools/git/... -v -run TestAnalyzeChanges
```

## 相关文件

- `pkg/ai/tools/git/analyze.go` - 主要实现
- `pkg/ai/tools/git/analyze_test.go` - 单元测试
- `pkg/ai/tools/git/register.go` - 工具注册
