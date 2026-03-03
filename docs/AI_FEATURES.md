# AI 功能完整指南

## 📖 目录

- [概述](#概述)
- [AI 对话系统](#ai-对话系统)
- [AI 代码审查](#ai-代码审查)
- [AI 辅助功能](#ai-辅助功能)
- [配置指南](#配置指南)
- [使用技巧](#使用技巧)
- [常见问题](#常见问题)

---

## 概述

Lazygit 增强版集成了强大的 AI 功能，提供智能化的 Git 操作建议、代码审查和问题诊断。

### 核心功能

```
┌─────────────────────────────────────────────────┐
│  🤖 AI 对话    │  🔍 代码审查  │  💡 智能辅助  │
│  交互式问答    │  质量分析     │  分支命名     │
│  上下文感知    │  改进建议     │  PR 描述      │
│  预设问题      │  问题检测     │  提交优化     │
└─────────────────────────────────────────────────┘
```

---

## AI 对话系统

### 功能特性

#### 1. 交互式多轮对话

AI 会记住对话历史，支持连续提问：

```
你: 如何合并分支？
AI: [解释合并方法]

你: 如果遇到冲突怎么办？
AI: [基于上下文，解释冲突解决]

你: 能给个具体例子吗？
AI: [提供详细示例]
```

#### 2. 智能上下文

AI 自动获取当前仓库状态：

```
📍 分支: feature/new-ui
📝 变更文件: 5 个
   M  src/main.go
   A  src/ui.go
   D  old/legacy.go
📌 最近提交: abc1234 - feat: add new UI
```

#### 3. 预设问题库

按 `Ctrl+P` 打开预设问题菜单：

##### 📚 基础操作
- 如何撤销最近的提交？
- 如何查看文件的修改历史？
- 如何暂存部分修改？
- 如何修改最后一次提交？

##### 🌿 分支管理
- 如何创建和切换分支？
- 如何合并分支？
- 如何删除本地和远程分支？
- 如何查看所有分支？

##### 🔄 远程操作
- 如何推送到远程仓库？
- 如何拉取远程更新？
- 如何解决推送冲突？
- 如何查看远程仓库信息？

##### 🔧 问题解决
- 遇到合并冲突怎么办？
- 如何恢复误删的文件？
- 如何清理大型仓库？
- 如何找回丢失的提交？

##### 💡 最佳实践
- 如何写好提交信息？
- 如何组织提交历史？
- 如何使用 Git Flow？
- 如何进行代码审查？

### 快捷键

#### 输入框快捷键

| 快捷键 | 功能 | 说明 |
|--------|------|------|
| **Enter** | 发送消息 | 发送当前输入 |
| **Alt+Enter** | 换行 | 多行输入 |
| **↑** | 上一条历史 | 浏览输入历史 |
| **↓** | 下一条历史 | 浏览输入历史 |
| **Ctrl+P** | 预设问题 | 打开预设问题菜单 |
| **Ctrl+K** | 停止生成 | 停止 AI 生成 |
| **Ctrl+L** | 清空历史 | 清空对话记录 |
| **Ctrl+S** | 保存对话 | 保存到文件 |
| **Tab** | 切换焦点 | 切换到对话视图 |
| **Esc** | 关闭 | 关闭对话窗口 |

#### 对话视图快捷键

| 快捷键 | 功能 | 说明 |
|--------|------|------|
| **↑/↓** | 滚动 | 上下滚动 |
| **PgUp/PgDn** | 翻页 | 快速翻页 |
| **Home** | 跳到顶部 | 滚动到开头 |
| **End** | 跳到底部 | 滚动到结尾 |
| **Ctrl+C** | 复制回复 | 复制最后一条 AI 回复 |
| **Ctrl+L** | 清空历史 | 清空对话记录 |
| **Tab** | 切换焦点 | 切换到输入框 |
| **?** | 帮助 | 显示帮助信息 |
| **Esc** | 关闭 | 关闭对话窗口 |

### 使用示例

#### 示例 1: 基础咨询

```
你: 如何撤销最近的提交？

AI: 撤销提交有几种方式：

1️⃣ 保留更改（推荐）
   git reset --soft HEAD~1
   ↳ 提交被撤销，更改保留在暂存区

2️⃣ 保留更改但不暂存
   git reset --mixed HEAD~1
   ↳ 提交被撤销，更改保留在工作区

3️⃣ 完全丢弃（危险）
   git reset --hard HEAD~1
   ↳ 提交和更改都被删除

💡 在 lazygit 中：
   • 在 commits 视图选择提交
   • 按 'g' 打开 reset 菜单
   • 选择 soft/mixed/hard reset
```

#### 示例 2: 问题诊断

```
你: 推送失败，提示 rejected

AI: 推送被拒绝通常有以下原因：

🔍 诊断步骤：
1. 检查远程是否有新提交
   git fetch origin
   git log HEAD..origin/main

2. 如果有新提交，需要先拉取
   git pull --rebase origin main

3. 解决可能的冲突后再推送
   git push origin main

💡 在 lazygit 中：
   • 按 'f' 执行 fetch
   • 按 'p' 执行 pull
   • 解决冲突后按 'P' 推送
```

#### 示例 3: 最佳实践

```
你: 如何写好提交信息？

AI: 优秀的提交信息应该：

📝 格式规范：
<type>(<scope>): <subject>

<body>

<footer>

🎯 类型（type）：
• feat: 新功能
• fix: Bug 修复
• docs: 文档更新
• style: 代码格式
• refactor: 重构
• test: 测试
• chore: 构建/工具

💡 示例：
feat(auth): 添加用户登录功能

实现了基于 JWT 的用户认证系统，包括：
- 登录接口
- Token 验证中间件
- 刷新 Token 机制

Closes #123
```

---

## AI 代码审查

### 功能特性

- 🔍 智能代码分析
- 💡 改进建议
- 🐛 潜在问题检测
- 📝 代码质量评估
- 🎯 最佳实践建议

### 使用方法

#### 方法 1: 从文件视图

```
1. 在文件视图选择文件
2. 按快捷键（需要配置）
3. AI 分析代码并给出建议
```

#### 方法 2: 从提交视图

```
1. 在提交视图选择提交
2. 按快捷键（需要配置）
3. AI 审查提交的所有更改
```

### 审查报告示例

```
🔍 代码审查报告

📊 总体评分: 8.5/10

✅ 优点:
• 代码结构清晰
• 命名规范
• 注释完整

⚠️ 改进建议:
1. 错误处理可以更细致
   位置: main.go:45
   建议: 区分不同类型的错误

2. 性能优化
   位置: utils.go:78
   建议: 使用缓存减少重复计算

3. 安全性
   位置: auth.go:23
   建议: 添加输入验证

💡 最佳实践:
• 考虑添加单元测试
• 可以提取重复代码
• 建议使用依赖注入
```

---

## AI 辅助功能

### 1. 智能分支命名

根据当前工作内容，AI 建议合适的分支名：

```
当前更改:
- 添加用户登录功能
- 实现 JWT 认证

AI 建议:
1. feature/user-authentication
2. feat/jwt-login
3. feature/auth-system
```

### 2. PR 描述生成

基于提交历史，AI 生成 PR 描述：

```
## 功能描述
实现了用户认证系统，包括登录、注册和 Token 管理。

## 主要更改
- 添加登录接口 (commit abc1234)
- 实现 JWT 中间件 (commit def5678)
- 添加用户注册功能 (commit ghi9012)

## 测试
- ✅ 单元测试覆盖率 85%
- ✅ 集成测试通过
- ✅ 手动测试完成

## 相关 Issue
Closes #123
```

### 3. 提交信息优化

AI 帮助改进提交信息：

```
原始提交信息:
"fix bug"

AI 优化后:
"fix(auth): 修复登录失败时的错误处理

- 添加详细的错误信息
- 改进错误日志记录
- 修复空指针异常

Fixes #456"
```

---

## 配置指南

### 基础配置

在 `~/.config/lazygit/config.yml` 中配置：

```yaml
ai:
  enabled: true
  activeProfile: "default"
  profiles:
    - name: "default"
      provider: "deepseek"
      apiKey: "your-api-key"
      model: "deepseek-chat"
      maxTokens: 2000
      timeout: 60
```

### 支持的 AI 提供商

#### 1. DeepSeek (推荐)

```yaml
ai:
  profiles:
    - name: "deepseek"
      provider: "deepseek"
      apiKey: "sk-xxx"
      model: "deepseek-chat"
      baseURL: "https://api.deepseek.com"
      maxTokens: 2000
      timeout: 60
```

**优点**:
- 高性价比
- 响应速度快
- 中文支持好

#### 2. OpenAI

```yaml
ai:
  profiles:
    - name: "openai"
      provider: "openai"
      apiKey: "sk-xxx"
      model: "gpt-4"
      baseURL: "https://api.openai.com/v1"
      maxTokens: 2000
      timeout: 60
```

**优点**:
- 功能强大
- 生态完善
- 文档丰富

#### 3. Anthropic (Claude)

```yaml
ai:
  profiles:
    - name: "claude"
      provider: "anthropic"
      apiKey: "sk-ant-xxx"
      model: "claude-3-opus-20240229"
      baseURL: "https://api.anthropic.com"
      maxTokens: 2000
      timeout: 60
```

**优点**:
- 推理能力强
- 安全性高
- 长文本支持好

#### 4. Ollama (本地)

```yaml
ai:
  profiles:
    - name: "ollama"
      provider: "ollama"
      model: "llama2"
      baseURL: "http://localhost:11434"
      maxTokens: 2000
      timeout: 120
```

**优点**:
- 完全本地运行
- 无需 API Key
- 隐私保护

#### 5. 自定义 API

```yaml
ai:
  profiles:
    - name: "custom"
      provider: "custom"
      apiKey: "your-key"
      model: "your-model"
      baseURL: "https://your-api.com"
      maxTokens: 2000
      timeout: 60
```

### 多 Profile 配置

```yaml
ai:
  enabled: true
  activeProfile: "deepseek"  # 默认使用的 profile
  profiles:
    - name: "deepseek"
      provider: "deepseek"
      apiKey: "sk-xxx"
      model: "deepseek-chat"

    - name: "gpt4"
      provider: "openai"
      apiKey: "sk-xxx"
      model: "gpt-4"

    - name: "local"
      provider: "ollama"
      model: "llama2"
```

### 高级配置

```yaml
ai:
  enabled: true
  activeProfile: "default"

  # 全局设置
  defaultMaxTokens: 2000
  defaultTimeout: 60
  retryCount: 3
  retryDelay: 1000

  # 代理设置
  proxy:
    enabled: true
    url: "http://localhost:7890"

  # 缓存设置
  cache:
    enabled: true
    ttl: 3600
    maxSize: 100

  profiles:
    - name: "default"
      provider: "deepseek"
      apiKey: "sk-xxx"
      model: "deepseek-chat"

      # Profile 特定设置
      temperature: 0.7
      topP: 0.9
      presencePenalty: 0.0
      frequencyPenalty: 0.0
```

---

## 使用技巧

### 技巧 1: 有效提问

✅ **好的问题**:
```
"如何在不丢失更改的情况下切换分支？"
"遇到 'detached HEAD' 状态怎么办？"
"如何查看某个文件在特定提交时的内容？"
```

❌ **不好的问题**:
```
"git"
"分支"
"帮我"
```

### 技巧 2: 利用上下文

AI 会自动获取仓库状态，可以直接问：

```
"当前分支有什么问题？"
"这些更改应该如何提交？"
"我应该先合并哪个分支？"
```

### 技巧 3: 使用预设问题

不确定怎么问时，按 `Ctrl+P` 浏览预设问题。

### 技巧 4: 保存有价值的对话

遇到有用的回答，按 `Ctrl+S` 保存对话记录。

### 技巧 5: 复制命令

AI 给出命令后，按 `Ctrl+C` 复制，直接在终端执行。

---

## 常见问题

### Q1: AI 功能无法使用？

**检查清单**:
1. 确认 AI 已启用: `ai.enabled: true`
2. 检查 API Key 是否正确
3. 检查网络连接
4. 查看日志: `lazygit --debug`

### Q2: AI 响应太慢？

**解决方案**:
1. 切换到更快的模型
2. 减少 maxTokens
3. 使用本地模型 (Ollama)
4. 检查网络延迟

### Q3: AI 回答不准确？

**改进方法**:
1. 提供更详细的问题描述
2. 使用更强大的模型 (GPT-4, Claude)
3. 增加上下文信息
4. 多轮对话澄清问题

### Q4: 如何切换 AI Provider？

```yaml
# 修改 activeProfile
ai:
  activeProfile: "gpt4"  # 切换到 GPT-4
```

### Q5: 担心隐私安全？

**建议**:
1. 使用本地模型 (Ollama)
2. 不要在对话中包含敏感信息
3. 定期清理对话历史
4. 使用自托管的 AI 服务

### Q6: API 配额用完了？

**解决方案**:
1. 切换到其他 Profile
2. 使用免费的本地模型
3. 减少 AI 使用频率
4. 升级 API 套餐

### Q7: 如何查看 AI 使用统计？

```bash
# 查看日志
lazygit --logs

# 查看配置
cat ~/.config/lazygit/config.yml
```

---

## 最佳实践

### 1. 合理使用 AI

✅ **适合使用 AI 的场景**:
- 学习新的 Git 概念
- 解决复杂的 Git 问题
- 获取最佳实践建议
- 代码审查和改进

❌ **不适合使用 AI 的场景**:
- 简单的 Git 操作
- 已经熟悉的命令
- 需要精确控制的操作

### 2. 保护隐私

- 不要在对话中包含密码、密钥
- 不要分享敏感的代码片段
- 定期清理对话历史
- 考虑使用本地模型

### 3. 验证建议

- AI 的建议仅供参考
- 重要操作前先测试
- 理解命令的含义再执行
- 保持备份习惯

### 4. 持续学习

- 通过 AI 学习 Git 知识
- 记录有用的回答
- 分享经验给团队
- 不断改进提问技巧

---

## 相关文档

- [项目概述](./PROJECT_OVERVIEW.md)
- [UI 组件指南](./UI_COMPONENTS.md)
- [配置指南](./Config.md)
- [开发指南](./DEVELOPMENT_GUIDE.md)

---

**版本**: v1.0.0
**最后更新**: 2024
**状态**: ✅ 完整
