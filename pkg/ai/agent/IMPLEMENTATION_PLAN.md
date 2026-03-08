# CodeReviewAgent重构实施计划

## 📅 总体时间线

**预计时间**: 5-7天
**当前阶段**: 准备阶段（第1天）

---

## ✅ 已完成工作

### 第1天：准备工作（已完成）

- [x] 创建详细评审报告 (`REVIEW_REPORT.md`)
- [x] 创建重构示例代码 (`code_review_agent_refactored.go`)
- [x] 创建统一Checkpointer接口 (`checkpointer_unified.go`)
- [x] 创建迁移指南 (`MIGRATION_GUIDE.md`)
- [x] 创建测试文件 (`code_review_agent_v2_test.go`)
- [x] 完善buildReviewPrompt方法

---

## 🚀 待完成工作

### 第2天：核心功能实现

#### 任务1：完善CodeReviewAgentV2实现
- [ ] 补充缺失的辅助函数（buildFocusSection, languageGuidelines）
- [ ] 添加错误处理和恢复逻辑
- [ ] 实现完整的超时控制
- [ ] 添加日志记录

**预计时间**: 3小时

**实施步骤**:
```bash
# 1. 复制辅助函数
# 从 code_review_agent.go 复制以下函数到 code_review_agent_refactored.go:
# - buildFocusSection
# - languageGuidelines
# - detectLanguage

# 2. 测试编译
go build ./pkg/ai/agent/...
```

#### 任务2：更新Checkpointer
- [ ] 将现有checkpointer.go中的实现迁移到checkpointer_unified.go
- [ ] 更新TwoPhaseAgent使用新的泛型接口
- [ ] 确保向后兼容

**预计时间**: 2小时

**实施步骤**:
```bash
# 1. 备份现有文件
cp pkg/ai/agent/checkpointer.go pkg/ai/agent/checkpointer_old.go.bak

# 2. 更新TwoPhaseAgent
# 修改 two_phase_agent.go 中的 Checkpointer 类型引用

# 3. 测试编译
go build ./pkg/ai/agent/...
```

#### 任务3：运行单元测试
- [ ] 运行code_review_agent_v2_test.go
- [ ] 修复测试失败
- [ ] 添加更多边界条件测试

**预计时间**: 2小时

**实施步骤**:
```bash
# 运行测试
go test -v ./pkg/ai/agent/ -run TestCodeReviewAgentV2

# 查看覆盖率
go test -cover ./pkg/ai/agent/ -run TestCodeReviewAgentV2
```

---

### 第3天：集成和GUI更新

#### 任务1：更新GUI层
- [ ] 修改 `ai_code_review_helper.go`
- [ ] 将 `NewCodeReviewAgent` 改为 `NewCodeReviewAgentV2`
- [ ] 测试GUI交互

**预计时间**: 2小时

**实施步骤**:
```bash
# 1. 备份GUI文件
cp pkg/gui/controllers/helpers/ai_code_review_helper.go \
   pkg/gui/controllers/helpers/ai_code_review_helper_old.go.bak

# 2. 修改代码
# 在 ai_code_review_helper.go 中:
# - 将 agent.NewCodeReviewAgent 改为 agent.NewCodeReviewAgentV2
# - 更新 checkpointer 创建代码

# 3. 编译测试
go build ./cmd/lazygit
```

#### 任务2：端到端测试
- [ ] 启动lazygit
- [ ] 测试代码评审功能
- [ ] 测试追问功能
- [ ] 测试检查点恢复

**预计时间**: 3小时

**测试清单**:
```
□ 基本评审流程
  □ 选择文件
  □ 触发评审
  □ 查看流式输出
  □ 验证评审结果

□ 追问功能
  □ 评审完成后追问
  □ 查看追问回复
  □ 多轮追问

□ 检查点恢复
  □ 评审中断
  □ 重启lazygit
  □ 验证状态恢复

□ 错误处理
  □ 网络超时
  □ 无效diff
  □ API错误
```

#### 任务3：性能测试
- [ ] 运行性能基准测试
- [ ] 对比新旧版本性能
- [ ] 优化性能瓶颈

**预计时间**: 2小时

**实施步骤**:
```bash
# 运行基准测试
go test -bench=. -benchmem ./pkg/ai/agent/ -run=^$ -bench BenchmarkCodeReviewAgentV2

# 对比结果
# 新版本应该与旧版本性能相当或更好
```

---

### 第4天：文档和清理

#### 任务1：更新文档
- [ ] 更新API文档
- [ ] 添加使用示例
- [ ] 更新CHANGELOG

**预计时间**: 2小时

**文档清单**:
```
□ API文档
  - CodeReviewAgentV2 接口说明
  - 节点函数说明
  - 状态转换图

□ 使用示例
  - 基本评审示例
  - 追问示例
  - 检查点恢复示例

□ CHANGELOG
  - 新增功能
  - 破坏性变更（如果有）
  - 迁移指南链接
```

#### 任务2：代码清理
- [ ] 删除旧的备份文件
- [ ] 统一代码风格
- [ ] 添加必要的注释

**预计时间**: 1小时

**实施步骤**:
```bash
# 1. 删除备份文件
rm pkg/ai/agent/*_old.go.bak

# 2. 格式化代码
go fmt ./pkg/ai/agent/...

# 3. 运行linter
golangci-lint run ./pkg/ai/agent/...
```

#### 任务3：准备提交
- [ ] 审查所有修改
- [ ] 编写详细的commit message
- [ ] 创建PR

**预计时间**: 1小时

**Commit Message模板**:
```
feat(ai/agent): 重构CodeReviewAgent为LangGraph架构

## 主要变更

1. **新增CodeReviewAgentV2**
   - 基于Graph节点的控制流
   - 纯函数式节点设计
   - 支持超时控制

2. **统一Checkpointer接口**
   - 使用泛型 Checkpointer[T any]
   - 向后兼容的类型别名

3. **改进错误处理**
   - 详细的错误分类
   - 恢复建议

## 架构改进

- 与TwoPhaseAgent架构完全一致
- 清晰的控制流图
- 更好的可测试性
- 易于扩展

## 向后兼容

- API保持不变（ReviewWithCallback, Ask）
- 现有调用代码仅需修改1行

## 测试

- 单元测试覆盖率 > 80%
- 端到端测试通过
- 性能无退化

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
```

---

### 第5-7天：灰度发布和监控

#### 任务1：灰度发布
- [ ] 配置feature flag（如果有）
- [ ] 50%用户使用新版本
- [ ] 收集用户反馈

**预计时间**: 持续监控

#### 任务2：监控和优化
- [ ] 监控错误率
- [ ] 监控性能指标
- [ ] 修复发现的问题

**预计时间**: 持续监控

#### 任务3：全量切换
- [ ] 100%用户使用新版本
- [ ] 移除旧代码（可选）
- [ ] 更新文档

**预计时间**: 1天

---

## 📊 进度跟踪

### 总体进度

```
准备阶段    ████████████████████ 100% (已完成)
核心实现    ░░░░░░░░░░░░░░░░░░░░   0% (待开始)
集成测试    ░░░░░░░░░░░░░░░░░░░░   0% (待开始)
文档清理    ░░░░░░░░░░░░░░░░░░░░   0% (待开始)
灰度发布    ░░░░░░░░░░░░░░░░░░░░   0% (待开始)
```

### 详细任务状态

| 阶段 | 任务 | 状态 | 预计时间 | 实际时间 |
|------|------|------|----------|----------|
| 准备 | 评审报告 | ✅ 完成 | 2h | 2h |
| 准备 | 重构示例 | ✅ 完成 | 3h | 3h |
| 准备 | 统一接口 | ✅ 完成 | 1h | 1h |
| 准备 | 测试文件 | ✅ 完成 | 2h | 2h |
| 核心 | 完善实现 | ⏳ 待开始 | 3h | - |
| 核心 | 更新Checkpointer | ⏳ 待开始 | 2h | - |
| 核心 | 单元测试 | ⏳ 待开始 | 2h | - |
| 集成 | 更新GUI | ⏳ 待开始 | 2h | - |
| 集成 | 端到端测试 | ⏳ 待开始 | 3h | - |
| 集成 | 性能测试 | ⏳ 待开始 | 2h | - |
| 文档 | 更新文档 | ⏳ 待开始 | 2h | - |
| 文档 | 代码清理 | ⏳ 待开始 | 1h | - |
| 文档 | 准备提交 | ⏳ 待开始 | 1h | - |

---

## 🎯 下一步行动

### 立即执行（今天）

1. **补充辅助函数**
   ```bash
   # 编辑 code_review_agent_refactored.go
   # 从 code_review_agent.go 复制:
   # - buildFocusSection
   # - languageGuidelines
   ```

2. **运行测试**
   ```bash
   go test -v ./pkg/ai/agent/ -run TestCodeReviewAgentV2
   ```

3. **修复编译错误**
   ```bash
   go build ./pkg/ai/agent/...
   ```

### 明天执行

1. **更新GUI层**
2. **端到端测试**
3. **性能测试**

---

## 🚨 风险和缓解措施

### 风险1：编译错误
**缓解**:
- 保留旧代码作为备份
- 渐进式修改
- 频繁编译测试

### 风险2：功能回归
**缓解**:
- 完整的测试覆盖
- 端到端测试
- 用户验收测试

### 风险3：性能退化
**缓解**:
- 性能基准测试
- 对比新旧版本
- 优化热点代码

### 风险4：用户体验变差
**缓解**:
- 保持API兼容
- 灰度发布
- 快速回滚机制

---

## 📞 支持和协作

### 需要帮助时

1. **查看文档**
   - REVIEW_REPORT.md
   - MIGRATION_GUIDE.md
   - code_review_agent_refactored.go

2. **运行测试**
   ```bash
   go test -v ./pkg/ai/agent/
   ```

3. **查看示例**
   - code_review_example.go
   - code_review_agent_v2_test.go

### 协作流程

1. **每日站会** - 同步进度和问题
2. **代码审查** - PR提交后及时审查
3. **问题跟踪** - 使用Issue跟踪问题

---

## ✅ 验收标准

### 功能验收
- [ ] 所有现有功能正常工作
- [ ] 新功能（超时控制）生效
- [ ] 检查点恢复正常
- [ ] 追问功能正常

### 质量验收
- [ ] 单元测试覆盖率 > 80%
- [ ] 所有测试通过
- [ ] 无编译警告
- [ ] Linter检查通过

### 性能验收
- [ ] 评审延迟 < 3秒
- [ ] 内存占用 < 10MB
- [ ] 无性能退化

### 文档验收
- [ ] API文档完整
- [ ] 使用示例清晰
- [ ] CHANGELOG更新

---

## 🎉 完成标志

当以下所有条件满足时，重构完成：

1. ✅ 所有测试通过
2. ✅ 文档更新完成
3. ✅ 代码审查通过
4. ✅ 用户验收通过
5. ✅ 性能验收通过
6. ✅ 灰度发布成功
7. ✅ 全量切换完成

**预计完成日期**: 2026-03-15

---

## 📝 备注

- 本计划可根据实际情况调整
- 优先保证质量，时间可适当延长
- 遇到问题及时沟通，不要独自解决
- 保持代码整洁，注释清晰
