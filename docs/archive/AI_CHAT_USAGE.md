# AI 对话功能使用说明

## 功能概述

lazygit 现在支持交互式 AI 对话功能，可以与 AI 助手进行多轮对话，获取 Git 相关的帮助、建议和指导。

## 功能特性

### 1. 交互式对话界面

```
╭────────────────────────────────────────────────────────────────╮
│ AI 对话 (Esc: 关闭, Ctrl+L: 清空历史)                          │
├────────────────────────────────────────────────────────────────┤
│                                                                │
│ AI: 15:30:45                                                   │
│   你好！我是 AI 助手，可以帮你解答 Git 相关问题、生成命令、   │
│   分析代码等。有什么我可以帮助你的吗？                         │
│                                                                │
│ 你: 15:31:20                                                   │
│   如何撤销最近的一次提交？                                     │
│                                                                │
│ AI: 15:31:25                                                   │
│   你可以使用以下命令：                                         │
│   - git reset --soft HEAD~1  # 保留更改                       │
│   - git reset --hard HEAD~1  # 丢弃更改                       │
│   在 lazygit 中，你可以在 commits 视图中选择提交，按 'g'      │
│   然后选择 reset 选项。                                        │
│                                                                │
├────────────────────────────────────────────────────────────────┤
│ 输入消息 (Enter: 发送)                                         │
│ 你的问题...                                                    │
╰────────────────────────────────────────────────────────────────╯
```

### 2. 核心功能

- **多轮对话**：支持连续对话，AI 会记住上下文
- **对话历史**：显示完整的对话历史，可滚动查看
- **实时上下文**：自动包含当前仓库状态（分支、文件、提交等）
- **时间戳**：每条消息都显示发送时间
- **视觉区分**：用户消息和 AI 回复使用不同颜色区分
- **流畅交互**：异步处理，不阻塞界面

### 3. 键盘快捷键

| 快捷键 | 功能 | 说明 |
|--------|------|------|
| Enter | 发送消息 | 在输入框中按 Enter 发送当前消息 |
| Esc | 关闭对话 | 关闭 AI 对话窗口 |
| Ctrl+L | 清空历史 | 清空对话历史（需确认） |
| Tab | 切换焦点 | 在对话视图和输入框之间切换 |
| ↑/↓ | 滚动 | 在对话视图中上下滚动 |

## 使用方法

### 基本使用

```go
// 打开 AI 对话
gui.ShowAIChat()
```

### 从菜单打开

```go
menuItems := []*types.MenuItem{
    {
        Label: "AI 对话",
        OnPress: func() error {
            return gui.ShowAIChat()
        },
        Key: 'a',
    },
}
```

### 从快捷键打开

在 keybindings 配置中添加：

```go
{
    Key: 'A',
    Handler: gui.ShowAIChat,
    Description: "打开 AI 对话",
}
```

## 应用场景

### 场景 1: Git 命令咨询

**用户**: "如何查看某个文件的修改历史？"

**AI**: "你可以使用以下命令：
- `git log -- <文件路径>` - 查看文件的提交历史
- `git log -p -- <文件路径>` - 查看文件的详细修改内容
- `git blame <文件路径>` - 查看每一行的最后修改者

在 lazygit 中，你可以：
1. 在 files 视图中选择文件
2. 按 'l' 查看文件的提交历史"

### 场景 2: 问题诊断

**用户**: "为什么我的推送被拒绝了？"

**AI**: "推送被拒绝通常有以下几个原因：

1. **远程有新提交**：其他人已经推送了新的提交
   - 解决：先 pull 再 push

2. **分支保护**：远程分支设置了保护规则
   - 解决：通过 Pull Request 提交

3. **强制推送冲突**：你的历史与远程不一致
   - 解决：使用 `git push --force-with-lease`（谨慎使用）

根据你当前的仓库状态，建议先执行 `git pull --rebase` 然后再推送。"

### 场景 3: 最佳实践建议

**用户**: "我应该如何组织我的提交？"

**AI**: "良好的提交组织遵循以下原则：

1. **原子性**：每个提交只做一件事
2. **完整性**：每个提交都是可工作的状态
3. **清晰的消息**：使用 Conventional Commits 格式
   - feat: 新功能
   - fix: 修复
   - docs: 文档
   - refactor: 重构

4. **逻辑分组**：相关的修改放在一起
5. **频繁提交**：小步快跑，便于回滚

在 lazygit 中，你可以使用 stage 功能精确控制每个提交包含的内容。"

### 场景 4: 冲突解决指导

**用户**: "遇到合并冲突了，怎么办？"

**AI**: "解决合并冲突的步骤：

1. **识别冲突文件**
   - 在 lazygit 中，冲突文件会标记为红色

2. **打开冲突文件**
   - 选择文件，按 Enter 打开编辑器

3. **理解冲突标记**
   ```
   <<<<<<< HEAD
   你的修改
   =======
   他人的修改
   >>>>>>> branch-name
   ```

4. **解决冲突**
   - 决定保留哪些代码
   - 删除冲突标记
   - 确保代码逻辑正确

5. **标记为已解决**
   - 保存文件
   - 在 lazygit 中 stage 该文件

6. **完成合并**
   - 提交合并结果

需要我详细解释某个步骤吗？"

### 场景 5: 分支管理建议

**用户**: "我有很多旧分支，如何清理？"

**AI**: "清理旧分支的安全步骤：

1. **查看已合并的分支**
   ```bash
   git branch --merged
   ```

2. **删除本地已合并分支**
   ```bash
   git branch -d <分支名>
   ```

3. **删除远程已合并分支**
   ```bash
   git push origin --delete <分支名>
   ```

4. **清理远程追踪分支**
   ```bash
   git fetch --prune
   ```

在 lazygit 中：
1. 切换到 branches 视图
2. 选择要删除的分支
3. 按 'd' 删除本地分支
4. 按 'D' 删除远程分支

建议：删除前确认分支已经合并，避免丢失工作。"

### 场景 6: 性能优化

**用户**: "仓库太大了，如何优化？"

**AI**: "优化大型仓库的方法：

1. **清理不必要的文件**
   ```bash
   git gc --aggressive --prune=now
   ```

2. **使用 Git LFS 管理大文件**
   ```bash
   git lfs install
   git lfs track "*.psd"
   ```

3. **浅克隆**（对于只需要最近历史的情况）
   ```bash
   git clone --depth 1 <仓库地址>
   ```

4. **清理旧的引用**
   ```bash
   git reflog expire --expire=now --all
   git gc --prune=now
   ```

5. **使用 .gitignore**
   - 排除构建产物
   - 排除依赖目录
   - 排除临时文件

当前你的仓库大小：[根据实际情况分析]
建议优先执行：[具体建议]"

## 技术实现

### 对话上下文

AI 对话会自动包含以下上下文信息：

1. **当前分支**：正在工作的分支名称
2. **变更文件**：未提交的文件列表
3. **最近提交**：最近的提交信息
4. **对话历史**：最近 5 条对话记录

### 消息格式

```go
type ChatMessage struct {
    Role      string    // "user" 或 "assistant"
    Content   string    // 消息内容
    Timestamp time.Time // 时间戳
}
```

### API 调用

```go
// 构建提示词
prompt := chat.buildPrompt(userMessage)

// 调用 AI
result, err := chat.gui.c.AI.Complete(chat.ctx, prompt)

// 处理回复
chat.messages = append(chat.messages, ChatMessage{
    Role:      "assistant",
    Content:   strings.TrimSpace(result.Content),
    Timestamp: time.Now(),
})
```

## 配置要求

### 前置条件

1. **启用 AI 功能**
   - 在设置中启用 AI
   - 配置 AI profile（provider, API key, model）

2. **支持的 AI 提供商**
   - DeepSeek
   - OpenAI
   - Anthropic
   - Ollama（本地）
   - Custom（自定义）

### 配置示例

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

## 注意事项

1. **网络连接**
   - AI 对话需要网络连接（除非使用本地 Ollama）
   - 请求可能需要几秒钟时间

2. **API 配额**
   - 注意 API 使用配额和费用
   - 建议设置合理的 maxTokens 限制

3. **隐私安全**
   - 对话内容会发送到 AI 服务商
   - 不要在对话中包含敏感信息（密码、密钥等）
   - 仓库上下文只包含基本信息（分支名、文件名、提交消息）

4. **对话历史**
   - 对话历史保存在内存中
   - 关闭对话窗口后历史会清空
   - 可以使用 Ctrl+L 手动清空历史

5. **错误处理**
   - 如果 AI 未启用，会显示错误提示
   - 如果 API 调用失败，会显示错误消息
   - 可以重新发送消息重试

## 文件结构

```
pkg/gui/
├── ai_chat.go              # AI 对话核心实现
├── ai_chat_examples.go     # 使用示例和测试菜单
└── controllers/helpers/
    └── ai_helper.go        # AI 助手功能（命令生成等）
```

## 与其他 AI 功能的关系

lazygit 包含多个 AI 功能：

1. **AI 对话**（本功能）
   - 交互式多轮对话
   - 问答、咨询、建议

2. **AI 助手**
   - 单次命令生成
   - 自动执行 Git 命令

3. **AI 代码审查**
   - 分析代码变更
   - 提供审查意见

这些功能可以配合使用，提供全方位的 AI 辅助。

## 未来改进

- [ ] 支持流式输出（实时显示 AI 回复）
- [ ] 保存对话历史到文件
- [ ] 支持导出对话记录
- [ ] 支持代码高亮显示
- [ ] 支持 Markdown 格式渲染
- [ ] 支持快捷回复模板
- [ ] 支持对话分支（多个独立对话）
- [ ] 集成到更多 lazygit 工作流中

## 相关文档

- [MESSAGEBOX_DESIGN.md](./MESSAGEBOX_DESIGN.md) - MessageBox 设计文档
- [MESSAGEBOX_USAGE.md](./MESSAGEBOX_USAGE.md) - MessageBox 使用文档
- [PROGRESSBAR_DESIGN.md](./PROGRESSBAR_DESIGN.md) - ProgressBar 设计文档
- [PROGRESSBAR_USAGE.md](./PROGRESSBAR_USAGE.md) - ProgressBar 使用文档

---

**实现完成**: 2024
**状态**: ✅ 已实现并可用
