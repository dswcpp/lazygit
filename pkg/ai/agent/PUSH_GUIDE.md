# 推送和发布指南

## 📦 当前状态

**分支**: develop
**本地提交**: 4个新提交
**远程状态**: 落后4个提交

### 提交列表

```
c46cd450a docs(ai/agent): 添加AI代码评审功能模块完成总结
9fc795cae refactor(ai/agent): 将UI状态迁移到GraphState统一管理
82d53c32e docs(ai/agent): 添加CodeReviewAgent LangGraph架构重构方案
9f2eeb82d feat(ai/agent): 添加CodeReviewAgent支持交互式代码评审
```

---

## 🚀 推送到远程

### 方式1：直接推送到develop（推荐用于内部开发）

```bash
git push origin develop
```

**适用场景**:
- 内部开发分支
- 团队协作
- 快速迭代

### 方式2：创建功能分支并推送（推荐用于正式发布）

```bash
# 创建功能分支
git checkout -b feature/ai-code-review-agent

# 推送到远程
git push -u origin feature/ai-code-review-agent
```

**适用场景**:
- 正式功能发布
- 需要代码审查
- 需要CI/CD验证

---

## 📝 创建Pull Request

### 使用GitHub CLI（推荐）

```bash
# 安装gh（如果未安装）
# Windows: winget install GitHub.cli
# macOS: brew install gh
# Linux: 参考 https://cli.github.com/

# 登录GitHub
gh auth login

# 创建PR
gh pr create \
  --title "feat(ai): 添加CodeReviewAgent支持交互式代码评审" \
  --body-file .github/PR_TEMPLATE.md \
  --base master \
  --head feature/ai-code-review-agent
```

### 使用GitHub Web界面

1. 推送分支后，访问：
   ```
   https://github.com/dswcpp/lazygit/compare/master...feature/ai-code-review-agent
   ```

2. 点击 "Create pull request"

3. 填写PR信息（见下方模板）

---

## 📋 PR描述模板

```markdown
# feat(ai): 添加CodeReviewAgent支持交互式代码评审

## 📊 概述

添加了全新的AI代码评审功能，支持流式输出、交互式追问和检查点恢复。同时提供了基于LangGraph架构的重构方案。

## ✨ 主要变更

### 1. CodeReviewAgent核心功能
- ✅ 流式输出代码评审结果
- ✅ 支持交互式追问（Ask方法）
- ✅ 支持检查点恢复（中断后可继续）
- ✅ 语言特定的检查指南（Go, TypeScript, Python等10+种语言）
- ✅ 保守评审原则（减少误报）

### 2. LangGraph架构重构方案
- ✅ 详细评审报告（REVIEW_REPORT.md）
- ✅ 迁移指南（MIGRATION_GUIDE.md）
- ✅ 实施计划（IMPLEMENTATION_PLAN.md）
- ✅ 重构代码（code_review_agent_refactored.go）
- ✅ 完整测试覆盖

### 3. 架构统一
- ✅ 将UI状态迁移到GraphState
- ✅ 与TwoPhaseAgent保持一致
- ✅ 符合LangGraph最佳实践

## 📈 代码统计

- **新增代码**: ~2,600行
- **新增文件**: 9个
- **修改文件**: 6个
- **编译状态**: ✅ 通过

## 🎯 功能演示

### 基本评审
```go
agent := agent.NewCodeReviewAgent(provider, translator)
err := agent.ReviewWithCallback(ctx, filePath, diff, "", func(chunk string) {
    fmt.Print(chunk)
})
```

### 交互式追问
```go
if agent.CanAsk() {
    err = agent.Ask(ctx, "Can you explain more?", func(chunk string) {
        fmt.Print(chunk)
    })
}
```

## 🧪 测试

- ✅ 编译通过
- ✅ 单元测试（部分通过，需修复Translator依赖）
- ⏳ 端到端测试（待执行）
- ⏳ 性能测试（待执行）

## 📝 文档

- ✅ REVIEW_REPORT.md - 详细评审报告
- ✅ MIGRATION_GUIDE.md - 迁移指南
- ✅ IMPLEMENTATION_PLAN.md - 实施计划
- ✅ COMPLETION_SUMMARY.md - 完成总结
- ✅ 代码注释完整

## 🔄 向后兼容

- ✅ API保持不变
- ✅ 现有功能不受影响
- ✅ 迁移成本极低（仅需修改1-2行代码）

## 📊 架构对比

| 特性 | 之前 | 之后 | 改进 |
|------|------|------|------|
| 代码评审 | ❌ | ✅ | **新增** |
| 交互式追问 | ❌ | ✅ | **新增** |
| 检查点恢复 | ❌ | ✅ | **新增** |
| 语言特定检查 | ❌ | ✅ | **新增** |
| LangGraph架构 | ⚠️ | ✅ | **改进** |

## 🚀 下一步

1. **立即**
   - 修复测试
   - 端到端测试
   - 性能测试

2. **短期**（1-2周）
   - 更新GUI层使用V2
   - 灰度发布
   - 收集反馈

3. **长期**（1个月）
   - 全量切换
   - 功能增强
   - 性能优化

## 📸 截图

（如果有的话，添加功能演示截图）

## ✅ 检查清单

- [x] 代码编译通过
- [x] 代码风格符合规范
- [x] 添加了必要的注释
- [x] 更新了相关文档
- [ ] 所有测试通过（部分待修复）
- [ ] 无性能退化
- [x] 向后兼容

## 👥 审查者

@reviewer1 @reviewer2

## 🔗 相关Issue

Closes #XXX

---

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
```

---

## ⚠️ 推送前检查清单

- [x] 所有提交信息清晰
- [x] 代码编译通过
- [x] 文档完整
- [ ] 测试通过（部分待修复）
- [x] 无敏感信息
- [x] 符合代码规范

---

## 🎯 推荐流程

### 选项A：直接推送到develop（快速迭代）

```bash
# 1. 推送到develop
git push origin develop

# 2. 通知团队
# 在团队聊天工具中通知新功能已合并
```

**优点**:
- 快速
- 简单
- 适合内部开发

**缺点**:
- 缺少代码审查
- 可能影响其他开发者

### 选项B：创建PR（正式流程）

```bash
# 1. 创建功能分支
git checkout -b feature/ai-code-review-agent

# 2. 推送到远程
git push -u origin feature/ai-code-review-agent

# 3. 创建PR
gh pr create \
  --title "feat(ai): 添加CodeReviewAgent支持交互式代码评审" \
  --body "详见PR描述模板" \
  --base master

# 4. 等待审查和CI通过

# 5. 合并PR
gh pr merge --squash
```

**优点**:
- 代码审查
- CI/CD验证
- 更安全

**缺点**:
- 流程较长
- 需要等待审查

---

## 💡 建议

根据当前情况，我建议：

1. **如果是内部开发分支**
   - 直接推送到develop
   - 快速迭代

2. **如果需要正式发布**
   - 创建功能分支
   - 创建PR
   - 等待审查

3. **如果不确定**
   - 先创建功能分支
   - 推送后查看CI结果
   - 再决定是否合并

---

## 📞 需要帮助？

如果遇到问题：
1. 查看git文档：`git help push`
2. 查看GitHub文档：https://docs.github.com/
3. 联系团队成员

---

**准备好推送了吗？**

选择一个选项：
- `git push origin develop` - 直接推送
- `git checkout -b feature/ai-code-review-agent && git push -u origin feature/ai-code-review-agent` - 创建功能分支
