# AI 对话功能完善设计方案

## 当前状态评估

### 已实现功能

#### 1. **AI 代码审查** (Ctrl+X)
- **位置**: 文件面板、提交文件面板
- **功能**: 对选中文件的 diff 进行流式 AI 审查
- **实现**: `pkg/gui/controllers/helpers/ai_code_review_helper.go`
- **特性**:
  - 流式 SSE 响应
  - 语言感知的审查指南（Go/TypeScript/Python/Rust/Java/C/C++/Shell/SQL）
  - 保守审查原则，避免误报
  - 弹窗显示结果（可缩放、复制）

#### 2. **AI 命令助手** (Ctrl+Y)
- **位置**: 全局快捷键
- **功能**: 根据自然语言描述生成 Git/Shell 命令
- **实现**: `pkg/gui/controllers/helpers/ai_helper.go` - `OpenAIAssistant()`
- **流程**:
  1. 用户输入任务描述
  2. AI 生成命令建议
  3. 用户确认后执行

#### 3. **AI 设置管理** (Ctrl+A)
- **位置**: 全局快捷键
- **功能**: Profile 管理（添加、编辑、删除、切换）
- **支持的提供商**: DeepSeek, OpenAI, Anthropic, Ollama, Custom
- **配置项**: API Key, Model, Endpoint, MaxTokens, Timeout, CustomHeaders

#### 4. **AI 提交信息生成** (Ctrl+G) ⚠️ 配置中存在但未实现
- **状态**: 仅在默认配置中定义了快捷键，无实际功能

### 架构组件

**提供商接口** (`pkg/ai/ai.go`):
```go
type Provider interface {
    Complete(ctx context.Context, prompt string) (Result, error)
    CompleteStream(ctx context.Context, prompt string, onChunk func(string)) error
}
```

**实现的提供商**:
- `openai_provider.go`: OpenAI 兼容（DeepSeek, OpenAI, Ollama, custom）
- `anthropic_provider.go`: Anthropic Claude

**辅助工具**:
- `ai_diff_filter.go`: Diff 过滤器（跳过锁文件/二进制文件，限制文件大小到 200 行）

---

## 核心问题分析

### 问题 1：功能覆盖不完整（HIGH 严重性）

#### 1.1 缺失的关键 AI 功能
- ❌ **AI 提交信息生成**: 配置中存在但未实现
- ❌ **AI 分支命名建议**: 根据变更内容建议分支名
- ❌ **AI PR 描述生成**: 根据提交历史生成 PR 描述
- ❌ **AI 冲突解决助手**: 分析冲突并提供解决建议
- ❌ **AI 代码解释**: 解释选中的代码或 diff 的作用
- ❌ **AI 重构建议**: 分析代码并提供重构建议
- ❌ **AI 测试生成**: 根据代码生成单元测试

#### 1.2 现有功能的局限
- **代码审查**: 仅支持单文件审查，不能整体审查多个文件的关联性
- **命令助手**: 上下文有限，缺乏仓库状态感知
- **无对话历史**: 不支持追问或迭代改进

### 问题 2：用户体验不足（HIGH 严重性）

#### 2.1 可发现性差
- ❌ AI 功能未在帮助文档中列出（`docs/keybindings/Keybindings_zh-CN.md` 中无 AI 相关内容）
- ❌ 无视觉提示指示 AI 功能已启用
- ❌ 新用户不知道何时可以使用 AI 辅助

#### 2.2 反馈机制不足
- ⚠️ 流式响应已实现，但无进度指示器
- ❌ 长时间 AI 请求无法取消
- ❌ 错误信息不够用户友好（直接显示 API 错误）
- ❌ 无重试机制

#### 2.3 交互流程僵化
- ❌ 无对话历史管理
- ❌ 无法针对 AI 回复追问
- ❌ 无法编辑并重新发送 prompt
- ❌ 无法保存/加载常用 prompt 模板

### 问题 3：集成深度不够（MEDIUM 严重性）

#### 3.1 工作流割裂
- ❌ AI 功能与 lazygit 核心工作流脱节
- ❌ 提交面板无 AI 建议按钮
- ❌ 冲突解决 UI 无 AI 辅助入口
- ❌ 分支创建流程无 AI 命名建议

#### 3.2 上下文感知不足
- **命令助手**: 仅包含基础 prompt，不包含：
  - 当前分支状态
  - 未推送的提交
  - 工作区文件状态
  - 最近的操作历史
- **代码审查**: 不考虑相关文件或架构上下文

### 问题 4：配置和文档（MEDIUM 严重性）

#### 4.1 文档缺失
- ❌ 无 AI 功能使用指南
- ❌ 快捷键未在自动生成的帮助中显示
- ❌ 无最佳实践或示例

#### 4.2 配置体验
- ⚠️ API Key 明文配置（支持 `${VAR}` 但用户可能不知道）
- ❌ 无配置向导（首次使用需手动编辑 YAML）
- ❌ 无内置 provider 测试功能

---

## 改进方案

### 阶段 1：完善现有功能（1-2 天）

#### 1.1 实现 AI 提交信息生成 ✅
**文件**: `pkg/gui/controllers/helpers/ai_commit_message_helper.go` (新建)

```go
func (self *AICommitMessageHelper) GenerateCommitMessage() error {
    // 1. 获取暂存的 diff
    diff, err := self.c.Git().Diff.GetStagedDiff()
    if err != nil || diff == "" {
        return errors.New(self.c.Tr.NoStagedChanges)
    }

    // 2. 过滤 diff
    filteredDiff := helpers.FilterDiffForAI(diff)

    // 3. 构建 prompt
    prompt := self.buildCommitMessagePrompt(filteredDiff)

    // 4. 调用 AI 生成
    result, err := self.c.AI().Complete(context.Background(), prompt)
    if err != nil {
        return err
    }

    // 5. 解析结果并填充到提交消息框
    self.c.Contexts().CommitMessage.SetContent(result.Content)
    self.c.Contexts().CommitMessage.RenderCommitLength()

    return nil
}
```

**快捷键绑定**: 已配置 Ctrl+G，仅需连接到 CommitMessageController

**Prompt 模板**:
```
You are a Git commit message expert. Generate a conventional commit message for the following staged changes.

Rules:
1. Use conventional commit format: <type>(<scope>): <description>
2. Types: feat, fix, refactor, docs, test, chore, perf, ci
3. Keep the first line under 72 characters
4. Add a detailed body if changes are complex (explain WHY, not WHAT)
5. Use Chinese for description and body

Staged changes:
{filteredDiff}

Output only the commit message, no explanations.
```

#### 1.2 增强代码审查上下文
**改进**: 在审查 prompt 中包含文件路径和相关信息

```go
func (self *AICodeReviewHelper) buildReviewPrompt(filePath, diff, language string) string {
    return fmt.Sprintf(`You are a code reviewer. Review the following %s code changes in file: %s

Focus on:
- Logic errors and bugs
- Security vulnerabilities
- Performance issues
- Code style violations
- Best practices for %s

File: %s
Changes:
%s

Provide concise, actionable feedback. If the code looks good, say "LGTM" briefly.
`, language, filePath, language, filePath, diff)
}
```

#### 1.3 改进错误处理和用户反馈
**文件**: `pkg/gui/controllers/helpers/ai_helper.go`

```go
func (self *AIHelper) callAIWithUserFeedback(
    ctx context.Context,
    prompt string,
    onSuccess func(Result) error,
) error {
    // 显示加载指示器
    self.c.OnUIThread(func() error {
        self.loadingHelper.WithWaitingStatus("正在请求 AI...", func() error {
            result, err := self.c.AI().Complete(ctx, prompt)
            if err != nil {
                // 用户友好的错误处理
                return self.handleAIError(err)
            }
            return onSuccess(result)
        })
        return nil
    })
    return nil
}

func (self *AIHelper) handleAIError(err error) error {
    if strings.Contains(err.Error(), "context deadline exceeded") {
        return errors.New("AI 请求超时，请稍后重试或调整超时设置")
    }
    if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "API key") {
        return errors.New("API 密钥无效，请检查 AI 设置")
    }
    if strings.Contains(err.Error(), "429") {
        return errors.New("API 请求频率限制，请稍后重试")
    }
    return fmt.Errorf("AI 请求失败: %v", err)
}
```

#### 1.4 添加 AI 功能到帮助文档
**文件**: 自动生成系统需更新 `pkg/i18n/*/translation.go`

确保以下键已正确翻译并在帮助中可见:
- `AIAssistant`: "AI 命令助手"
- `AICodeReview`: "AI 代码审查"
- `AISettings`: "AI 设置"
- `AIGenerateCommitMessage`: "AI 生成提交信息"

---

### 阶段 2：新增核心功能（3-5 天）

#### 2.1 AI 分支命名建议
**触发位置**: 创建分支时（`n` 键）
**文件**: `pkg/gui/controllers/helpers/ai_branch_naming_helper.go` (新建)

```go
func (self *AIBranchNamingHelper) SuggestBranchName() error {
    // 1. 分析工作区变更
    files := self.c.Model().Files
    if len(files) == 0 {
        return errors.New("无工作区变更可分析")
    }

    // 2. 获取变更摘要
    diff, _ := self.c.Git().Diff.GetUnstagedDiff()

    // 3. Prompt
    prompt := fmt.Sprintf(`Based on the following code changes, suggest a git branch name.

Rules:
- Use kebab-case (lowercase with hyphens)
- Format: <type>/<short-description>
- Types: feature, bugfix, hotfix, refactor, docs, test
- Keep under 40 characters
- Be specific and descriptive

Changes:
%s

Output only the branch name, no explanations.
`, helpers.FilterDiffForAI(diff))

    // 4. 调用 AI
    result, err := self.c.AI().Complete(context.Background(), prompt)
    if err != nil {
        return err
    }

    // 5. 填充到输入框
    return self.c.Prompt(types.PromptOpts{
        Title:         "创建新分支",
        InitialContent: strings.TrimSpace(result.Content),
        HandleConfirm: func(branchName string) error {
            return self.c.Git().Branch.New(branchName, "")
        },
    })
}
```

**集成**: 在 `BranchesController` 中添加 `Ctrl+Shift+N` 快捷键调用 AI 建议

#### 2.2 AI PR 描述生成
**触发位置**: 创建 PR 时
**文件**: `pkg/gui/controllers/helpers/ai_pr_helper.go` (新建)

```go
func (self *AIPRHelper) GeneratePRDescription(baseBranch, headBranch string) error {
    // 1. 获取提交历史
    commits, err := self.c.Git().Loaders.Commits.GetCommits(
        loaders.GetCommitsOptions{
            Limit:      100,
            FilterPath: "",
            RefName:    fmt.Sprintf("%s..%s", baseBranch, headBranch),
        },
    )
    if err != nil {
        return err
    }

    // 2. 获取完整 diff
    diff, err := self.c.Git().Diff.GetDiff(baseBranch, headBranch)
    if err != nil {
        return err
    }

    // 3. 构建 prompt
    commitSummary := lo.Map(commits, func(c *models.Commit, _ int) string {
        return fmt.Sprintf("- %s: %s", c.Hash[:7], c.Name)
    })

    prompt := fmt.Sprintf(`Generate a GitHub Pull Request description for the following changes.

Branch: %s -> %s

Commits:
%s

Diff (filtered):
%s

Format:
## Summary
<1-3 bullet points summarizing the changes>

## Changes
<Detailed list of changes by category>

## Test Plan
<How to test these changes>

Use Chinese for the description.
`, headBranch, baseBranch, strings.Join(commitSummary, "\n"), helpers.FilterDiffForAI(diff))

    // 4. 流式生成并显示
    return self.showStreamingResult(prompt, "AI 生成 PR 描述")
}
```

#### 2.3 AI 冲突解决助手
**触发位置**: 冲突解决面板
**文件**: `pkg/gui/controllers/helpers/ai_conflict_helper.go` (新建)

```go
func (self *AIConflictHelper) AnalyzeConflict(filePath string) error {
    // 1. 读取冲突文件
    content, err := self.c.Git().File.GetFileContent(filePath)
    if err != nil {
        return err
    }

    // 2. 提取冲突标记
    conflicts := self.extractConflictMarkers(content)

    // 3. Prompt
    prompt := fmt.Sprintf(`Analyze the following merge conflict and suggest a resolution strategy.

File: %s

Conflict:
%s

Explain:
1. What caused this conflict
2. Which version (ours/theirs) should be kept and why
3. Or if a manual merge is needed, provide suggestions

Use Chinese for the explanation.
`, filePath, conflicts)

    // 4. 流式显示分析
    return self.showStreamingAnalysis(prompt)
}
```

#### 2.4 AI 代码解释
**触发位置**: 文件面板、提交面板
**快捷键**: `Ctrl+E`
**文件**: `pkg/gui/controllers/helpers/ai_explain_helper.go` (新建)

```go
func (self *AIExplainHelper) ExplainCode(filePath, code string) error {
    prompt := fmt.Sprintf(`Explain what this code does in simple terms.

File: %s

Code:
%s

Provide:
1. A brief summary (1-2 sentences)
2. Key functionality
3. Any potential issues or improvements

Use Chinese for the explanation.
`, filePath, code)

    return self.showStreamingResult(prompt, "AI 代码解释")
}
```

---

### 阶段 3：交互体验增强（2-3 天）

#### 3.1 对话历史管理
**文件**: `pkg/gui/controllers/helpers/ai_conversation_helper.go` (新建)

**数据结构**:
```go
type AIConversation struct {
    ID        string
    Title     string
    Messages  []AIMessage
    CreatedAt time.Time
}

type AIMessage struct {
    Role      string // "user" | "assistant"
    Content   string
    Timestamp time.Time
}
```

**功能**:
- 保存对话历史到 `~/.config/lazygit/ai_conversations/`
- 在 AI 弹窗中添加 "继续对话" 选项
- `Ctrl+H` 快捷键查看历史对话

#### 3.2 可取消的 AI 请求
**实现**: 使用 `context.WithCancel`

```go
func (self *AIHelper) callAICancellable(prompt string) error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 在弹窗中添加 "按 Esc 取消" 提示
    self.c.OnUIThread(func() error {
        self.c.Contexts().AICodeReview.SetContent("正在请求 AI...\n\n按 Esc 取消")

        // 注册 Esc 处理
        self.registerCancelHandler(cancel)

        go func() {
            result, err := self.c.AI().CompleteStream(ctx, prompt, func(chunk string) {
                self.c.OnUIThread(func() error {
                    self.c.Contexts().AICodeReview.AppendContent(chunk)
                    return nil
                })
            })

            if err != nil {
                if errors.Is(err, context.Canceled) {
                    self.c.Toast("AI 请求已取消")
                } else {
                    self.handleAIError(err)
                }
            }
        }()

        return nil
    })

    return nil
}
```

#### 3.3 Prompt 模板系统
**配置文件**: `~/.config/lazygit/ai_prompts.yml`

```yaml
prompts:
  - name: "简洁代码审查"
    template: |
      Review the following code changes briefly.
      Focus only on critical issues (bugs, security).
      {diff}

  - name: "详细代码审查"
    template: |
      Thoroughly review the following code changes.
      Check for: logic errors, security, performance, style, best practices.
      {diff}

  - name: "重构建议"
    template: |
      Suggest refactoring improvements for:
      {code}
```

**UI**: 在 AI 功能菜单中添加 "选择 Prompt 模板" 选项

#### 3.4 进度指示器
**实现**: 在流式响应时显示动画

```go
func (self *AICodeReviewHelper) ReviewDiffWithProgress(filePath, diff string) error {
    self.c.OnUIThread(func() error {
        view := self.c.Views().AICodeReview
        view.Clear()
        view.Title = "AI 代码审查 - 正在分析..."

        spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
        frame := 0

        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()

        go func() {
            for range ticker.C {
                self.c.OnUIThread(func() error {
                    view.Title = fmt.Sprintf("AI 代码审查 %s", spinner[frame%len(spinner)])
                    frame++
                    return nil
                })
            }
        }()

        // ... 调用 AI CompleteStream

        return nil
    })

    return nil
}
```

---

### 阶段 4：深度集成（3-4 天）

#### 4.1 提交面板集成
**位置**: `pkg/gui/context/commit_message_context.go`

添加 AI 建议按钮到提交消息输入框:

```go
func (c *CommitMessageContext) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
    return []*types.Binding{
        // ... 现有绑定
        {
            Key:         opts.GetKey(opts.Config.CommitMessage.AIGenerateCommitMessage),
            Handler:     c.aiGenerateCommitMessage,
            Description: c.c.Tr.AIGenerateCommitMessage,
            Tooltip:     "使用 AI 根据暂存的变更生成提交信息",
        },
    }
}
```

#### 4.2 冲突解决 UI 集成
**位置**: `pkg/gui/controllers/merge_conflicts_controller.go`

在冲突解决面板添加 AI 分析选项:

```go
{
    Key:         'a', // AI analyze
    Handler:     self.aiAnalyzeConflict,
    Description: "AI 分析冲突",
},
```

#### 4.3 增强命令助手上下文
**改进**: 包含更多仓库状态信息

```go
func (self *AIHelper) buildAssistantPrompt(userInput string) string {
    // 收集上下文
    currentBranch := self.c.Helpers().Refs.GetCheckedOutRef()
    unpushedCommits := self.c.Model().Commits[:5] // 最近 5 个提交
    workingTreeState := self.c.Git().Status.WorkingTreeState()

    context := fmt.Sprintf(`Current repository state:
- Branch: %s
- Unpushed commits: %d
- Working tree state: %s
- Uncommitted changes: %d files

User request: %s

Generate appropriate git or shell commands to accomplish this task.
Explain what each command does.
`, currentBranch.Name, len(unpushedCommits), workingTreeState, len(self.c.Model().Files), userInput)

    return context
}
```

---

### 阶段 5：配置和文档（1-2 天）

#### 5.1 首次使用向导
**触发**: AI 功能未启用时按任意 AI 快捷键

```go
func (self *AIHelper) showFirstTimeWizard() error {
    return self.c.Menu(types.CreateMenuOptions{
        Title: "AI 功能配置向导",
        Items: []*types.MenuItem{
            {
                Label: "使用 DeepSeek (推荐)",
                OnPress: func() error {
                    return self.setupProvider("deepseek")
                },
            },
            {
                Label: "使用 OpenAI",
                OnPress: func() error {
                    return self.setupProvider("openai")
                },
            },
            {
                Label: "使用 Anthropic Claude",
                OnPress: func() error {
                    return self.setupProvider("anthropic")
                },
            },
            {
                Label: "稍后配置",
                OnPress: func() error { return nil },
            },
        },
    })
}
```

#### 5.2 Provider 测试功能
**位置**: AI 设置菜单

```go
{
    Label: "测试当前 Profile",
    OnPress: func() error {
        return self.testCurrentProfile()
    },
}

func (self *AIHelper) testCurrentProfile() error {
    self.loadingHelper.WithWaitingStatus("测试 AI 连接...", func() error {
        result, err := self.c.AI().Complete(context.Background(), "Say 'OK' if you can hear me.")
        if err != nil {
            return fmt.Errorf("连接失败: %v", err)
        }
        self.c.Toast(fmt.Sprintf("✓ 连接成功: %s", result.Content))
        return nil
    })
    return nil
}
```

#### 5.3 文档更新
**文件**: 新建 `docs/AI_Features.md`

```markdown
# AI 功能使用指南

## 功能概览

Lazygit 集成了多种 AI 辅助功能，帮助您更高效地完成 Git 操作。

### 1. AI 代码审查 (Ctrl+X)
在文件面板或提交文件面板中，选中一个文件后按 Ctrl+X，AI 将分析代码变更并提供审查意见。

### 2. AI 命令助手 (Ctrl+Y)
按 Ctrl+Y 打开 AI 助手，用自然语言描述您想执行的 Git 操作，AI 将生成相应的命令。

### 3. AI 生成提交信息 (Ctrl+G)
在提交面板中按 Ctrl+G，AI 将根据暂存的变更自动生成规范的提交信息。

### 4. AI 分支命名 (Ctrl+Shift+N)
创建新分支时，AI 可以根据当前的变更内容建议合适的分支名。

### 5. AI PR 描述生成
创建 Pull Request 时，AI 可以根据提交历史和变更生成详细的 PR 描述。

## 配置

### 首次设置

1. 按 Ctrl+A 打开 AI 设置
2. 选择 "添加 Profile"
3. 输入 API 密钥（支持环境变量，如 ${DEEPSEEK_API_KEY}）
4. 选择模型和提供商

### 支持的提供商

- **DeepSeek**: 性价比高，推荐使用 deepseek-reasoner
- **OpenAI**: GPT-4o, GPT-4o-mini
- **Anthropic**: Claude Sonnet 4.5
- **Ollama**: 本地模型
- **Custom**: 自定义 OpenAI 兼容 API

## 最佳实践

1. **代码审查**: 在提交前审查关键文件变更
2. **提交信息**: 让 AI 生成初稿，然后根据需要修改
3. **命令助手**: 用于复杂的 Git 操作（如 cherry-pick、rebase）
4. **隐私**: 敏感代码可以使用本地 Ollama 模型

## 故障排查

- **API 密钥错误**: 检查配置文件或环境变量
- **请求超时**: 调整 `timeout` 设置
- **频率限制**: 更换提供商或等待冷却
```

---

## 实施优先级

### P0 (必须完成)
1. ✅ 实现 AI 提交信息生成 (Ctrl+G)
2. ✅ 改进错误处理和用户反馈
3. ✅ 添加 AI 功能到帮助文档
4. ✅ 实现可取消的 AI 请求

### P1 (高优先级)
1. AI 分支命名建议
2. AI PR 描述生成
3. 进度指示器和流式反馈
4. 增强命令助手上下文
5. 首次使用向导

### P2 (中优先级)
1. AI 代码解释功能
2. Prompt 模板系统
3. 对话历史管理
4. Provider 测试功能
5. 提交面板深度集成

### P3 (可选增强)
1. AI 冲突解决助手
2. AI 重构建议
3. AI 测试生成
4. 自定义审查规则

---

## 技术实施路线图

### Week 1: 核心功能完善
- Day 1-2: 实现 AI 提交信息生成
- Day 3-4: 改进错误处理、可取消请求、进度指示器
- Day 5: 更新文档和帮助系统

### Week 2: 新功能开发
- Day 1-2: AI 分支命名 + AI PR 描述
- Day 3-4: 增强上下文、对话历史
- Day 5: Prompt 模板系统

### Week 3: 深度集成
- Day 1-2: 提交面板集成
- Day 3: 冲突解决集成
- Day 4-5: 首次使用向导、Provider 测试、文档完善

---

## 测试计划

### 单元测试
```go
// pkg/gui/controllers/helpers/ai_commit_message_helper_test.go
func TestGenerateCommitMessage(t *testing.T) {
    // Mock AI provider
    // Test diff filtering
    // Test prompt construction
    // Test commit message parsing
}
```

### 集成测试
```go
// pkg/integration/tests/ai/commit_message.go
var AICommitMessage = NewIntegrationTest(NewIntegrationTestArgs{
    Description: "AI generates commit message from staged changes",
    Run: func(t *TestDriver, keys config.KeybindingConfig) {
        // Stage files
        // Trigger AI commit message (Ctrl+G)
        // Verify commit message populated
    },
})
```

### 手动测试清单
- [ ] AI 提交信息生成正常工作
- [ ] AI 代码审查流式输出正常
- [ ] AI 命令助手生成正确命令
- [ ] 错误处理用户友好
- [ ] 可以取消长时间请求
- [ ] 进度指示器正常显示
- [ ] 首次使用向导流程正常
- [ ] 所有 AI 快捷键在帮助中可见
- [ ] 多 Profile 切换正常
- [ ] Provider 测试功能正常

---

## 关键文件清单

### 新建文件
1. `pkg/gui/controllers/helpers/ai_commit_message_helper.go` - 提交信息生成
2. `pkg/gui/controllers/helpers/ai_branch_naming_helper.go` - 分支命名建议
3. `pkg/gui/controllers/helpers/ai_pr_helper.go` - PR 描述生成
4. `pkg/gui/controllers/helpers/ai_conflict_helper.go` - 冲突分析
5. `pkg/gui/controllers/helpers/ai_explain_helper.go` - 代码解释
6. `pkg/gui/controllers/helpers/ai_conversation_helper.go` - 对话历史
7. `docs/AI_Features.md` - AI 功能文档

### 修改文件
1. `pkg/gui/controllers/helpers/ai_helper.go` - 错误处理、可取消请求
2. `pkg/gui/controllers/helpers/ai_code_review_helper.go` - 进度指示器
3. `pkg/gui/controllers/commit_message_controller.go` - 集成 AI 生成
4. `pkg/gui/controllers/branches_controller.go` - AI 分支命名
5. `pkg/gui/controllers/global_controller.go` - 新快捷键绑定
6. `pkg/i18n/*/translation.go` - 翻译键补充
7. `docs/keybindings/Keybindings_zh-CN.md` - 自动生成系统更新

---

## 成功标准

1. ✅ 所有 P0 功能已实现并通过测试
2. ✅ AI 功能在帮助文档中可见
3. ✅ 用户可以流畅使用 AI 辅助完成 Git 操作
4. ✅ 错误处理用户友好，无崩溃
5. ✅ 首次使用体验良好（向导 + 文档）
6. ✅ 单元测试覆盖率 > 70%
7. ✅ 集成测试覆盖核心场景
8. ✅ 性能可接受（AI 请求延迟 < 10s）

---

## 风险与缓解

### 风险 1: AI 响应质量不稳定
**缓解**:
- 提供多个 Prompt 模板供用户选择
- 允许用户编辑 AI 生成的结果
- 添加 "重新生成" 功能

### 风险 2: API 费用
**缓解**:
- 默认使用性价比高的模型（如 DeepSeek）
- 提供本地 Ollama 选项
- Diff 过滤减少 token 消耗

### 风险 3: 用户隐私担忧
**缓解**:
- 文档明确说明数据流向
- 推荐本地模型用于敏感代码
- 提供禁用 AI 功能的开关

### 风险 4: 集成复杂度
**缓解**:
- 分阶段实施，每阶段独立验证
- 遵循现有 lazygit 架构模式
- 充分的单元测试和集成测试

---

## 后续演进

### 未来可能的功能
1. **AI 代码搜索**: 自然语言搜索代码
2. **AI 变更摘要**: 自动总结最近的提交
3. **AI 重构工具**: 自动化重构建议
4. **团队学习**: 从团队的 Git 历史中学习模式
5. **多模态支持**: 支持图片/图表生成

### 性能优化方向
1. 缓存 AI 响应（相同 diff 不重复请求）
2. 批量处理多个文件审查
3. 预加载常用 Prompt
4. 使用 Streaming 优化感知延迟
