# AI 模块翻译键映射表

## 通用词汇 (Common)

| 中文 | 翻译键 | 英文 |
|------|--------|------|
| 取消 | AICancel | Cancel |
| 确定 | AIOK | OK |
| 确认 | AIConfirm | Confirm |
| 是 | AIYes | Yes |
| 否 | AINo | No |
| 成功 | AISuccess | Success |
| 失败 | AIFailed | Failed |
| 错误 | AIError | Error |
| 警告 | AIWarning | Warning |
| 未知 | AIUnknown | Unknown |
| 执行中 | AIExecuting | Executing |
| 思考中 | AIThinking | Thinking |
| 空闲 | AIIdle | Idle |
| 已取消 | AICancelled | Cancelled |
| 正在思考 | AIThinkingInProgress | Thinking... |

## AI Agent

| 中文 | 翻译键 | 英文 |
|------|--------|------|
| 规划阶段不允许调用工具: %s | AIAgentToolNotAllowedInPlanning | Tool not allowed in planning phase: %s |
| 关键步骤失败: %s — %s | AIAgentCriticalStepFailed | Critical step failed: %s — %s |
| ⏱️ 步骤执行超时（%v）: %s | AIAgentStepTimeout | ⏱️ Step execution timeout (%v): %s |
| [用户拒绝] 工具 %s 未被执行，请据此调整后续操作。 | AIAgentUserRejectedTool | [User rejected] Tool %s was not executed, please adjust subsequent operations. |
| 手动解决冲突后继续 | AIAgentResolveConflictManually | Resolve conflict manually and continue |
| 设置上游分支 | AIAgentSetUpstreamBranch | Set upstream branch |
| 冲突 | AIAgentConflict | Conflict |
| 工具名 | AIAgentToolName | Tool name |
| 先暂存文件（stage_all 或 stage_file） | AIAgentStageFilesFirst | Stage files first (stage_all or stage_file) |
| \n\n💡 可能的原因： | AIAgentPossibleReasons | \n\n💡 Possible reasons: |
| feat: 添加用户登录功能 | AIAgentExampleCommitMsg | feat: add user login feature |
| 不要 | AIAgentDont | Don't |
| ## 当前仓库状态\n\n%s\n\n## 用户指令\n\n%s | AIAgentRepoStatusAndUserInstruction | ## Current Repository Status\n\n%s\n\n## User Instruction\n\n%s |

## AI Tools

| 中文 | 翻译键 | 英文 |
|------|--------|------|
| 缺少 %s 参数 | AIToolMissingParam | Missing %s parameter |
| 缺少 name 参数 | AIToolMissingNameParam | Missing name parameter |
| 缺少 path 参数 | AIToolMissingPathParam | Missing path parameter |
| 缺少 message 参数 | AIToolMissingMessageParam | Missing message parameter |
| 缺少 hash 参数 | AIToolMissingHashParam | Missing hash parameter |
| 文件路径 | AIToolFilePath | File path |
| 分支名称 | AIToolBranchName | Branch name |
| tag 名称 | AIToolTagName | Tag name |
| 标签名称 | AIToolTagName | Tag name |
| 提交信息 | AIToolCommitMessage | Commit message |
| 无变更 | AIToolNoChanges | No changes |
| 工作区 | AIToolWorkingDir | Working directory |
| 暂存区 | AIToolStagingArea | Staging area |
| 目标 ref 或 hash（优先） | AIToolTargetRefOrHash | Target ref or hash (preferred) |
| 回退步数（ref 为空时使用，默认 1） | AIToolResetSteps | Reset steps (used when ref is empty, default 1) |
| stash 索引，默认 0 | AIToolStashIndex | Stash index, default 0 |
| 最多返回行数（默认 300，0 表示不限制） | AIToolMaxLines | Maximum lines to return (default 300, 0 for unlimited) |
| 指向的 ref（默认 HEAD） | AIToolTargetRef | Target ref (default HEAD) |
| push 配置错误: %v | AIToolPushConfigError | Push configuration error: %v |
| 已将当前分支 rebase 到 %s | AIToolRebasedTo | Rebased current branch to %s |
| 重命名失败: %v | AIToolRenameFailed | Rename failed: %v |
| 丢弃变更失败: %v | AIToolDiscardChangesFailed | Discard changes failed: %v |
| 参数 | AIToolParam | Parameter |
| 值 | AIToolValue | Value |

## AI Skills

| 中文 | 翻译键 | 英文 |
|------|--------|------|
| 当前分支: %s\n | AISkillCurrentBranch | Current branch: %s\n |
| - 只输出分支名，不要任何解释\n | AISkillBranchNameOnly | - Output branch name only, no explanation\n |
| - description: 小写 kebab-case，2-5 个单词\n | AISkillBranchNameFormat | - description: lowercase kebab-case, 2-5 words\n |
| Windows + Git Bash，用 && 连接命令 | AISkillWindowsGitBash | Windows + Git Bash, use && to connect commands |
| 运行环境: %s\n\n | AISkillRuntime | Runtime environment: %s\n\n |
| 输出 JSON 数组，每个元素包含:\n | AISkillOutputJSONArray | Output JSON array, each element contains:\n |
| - explanation: 中文解释（1-2 句）\n | AISkillExplanation | - explanation: Chinese explanation (1-2 sentences)\n |
| - subject: 中文，动词开头，祈使句，不超过 72 字符\n | AISkillCommitSubject | - subject: Chinese, verb-first, imperative, max 72 chars\n |
| 场景提示: 这是测试相关变更，使用 test 类型。\n | AISkillTestScenario | Scenario hint: This is test-related change, use test type.\n |
| \n请直接输出提交信息： | AISkillOutputCommitMsg | \nPlease output commit message directly: |
| 场景提示: 这是重构，优先使用 refactor 类型。\n | AISkillRefactorScenario | Scenario hint: This is refactoring, prefer refactor type.\n |
| ## 请生成包含以下部分的 PR 描述\n | AISkillGeneratePRDesc | ## Please generate PR description with the following sections\n |
| ### Summary\n一句话说明 PR 的目的。\n\n | AISkillPRSummary | ### Summary\nOne sentence describing the purpose of this PR.\n\n |
| ### Testing\n- 说明如何验证这些变更\n | AISkillPRTesting | ### Testing\n- Explain how to verify these changes\n |
| ## 代码变更\n```diff\n | AISkillCodeChanges | ## Code Changes\n```diff\n |
| \nDiff 摘要:\n```diff\n | AISkillDiffSummary | \nDiff summary:\n```diff\n |
| ## 仓库背景\n | AISkillRepoContext | ## Repository Context\n |
| ## 代码变更\n | AISkillCodeChangesTitle | ## Code Changes\n |
| ## 分支信息\n从 `%s` 合并到 `%s`\n\n | AISkillBranchInfo | ## Branch Info\nMerging from `%s` to `%s`\n\n |
| ## 提交历史\n | AISkillCommitHistory | ## Commit History\n |

## AI Chat (GUI)

| 中文 | 翻译键 | 英文 |
|------|--------|------|
| AI 未启用 | AIChatNotEnabled | AI not enabled |
| 可输入下一条指令 | AIChatCanInputNext | You can input the next command |
| 正在分析并生成执行计划 | AIChatGeneratingPlan | Analyzing and generating execution plan |
| branch-name: 分支名称 | AIChatTemplateBranchName | branch-name: branch name |
| tag-name: 标签名称 | AIChatTemplateTagName | tag-name: tag name |
| message: 提交信息 | AIChatTemplateMessage | message: commit message |
| 正在推送到远程仓库... | AIChatPushingToRemote | Pushing to remote repository... |
| 中止合并 | AIChatAbortMerge | Abort merge |
| 解决冲突 | AIChatResolveConflict | Resolve conflict |
| 冲突文件:\n | AIChatConflictFiles | Conflict files:\n |
| 合并冲突 | AIChatMergeConflict | Merge conflict |
| 检测到未提交的更改，如何处理？ | AIChatUncommittedChanges | Uncommitted changes detected, how to handle? |
| 删除成功 | AIChatDeleteSuccess | Delete successful |
| ' 吗？ | AIChatConfirmSuffix | ? |

## Other

| 中文 | 翻译键 | 英文 |
|------|--------|------|
| ... 还有 %d 个\n | AIMoreItems | ... %d more\n |
| 工作区: 干净\n | AIRepoWorkingDirClean | Working directory: clean\n |
| ⚠ 正在进行: %s\n | AIRepoInProgress | ⚠ In progress: %s\n |
| 根据功能描述生成合适的 Git 分支名（kebab-case，带类型前缀） | AIManagerGenerateBranchName | Generate appropriate Git branch name based on feature description (kebab-case with type prefix) |
| 参数 | AIManagerParam | Parameter |
| 远程: %s [已同步]\n | AIRepoRemoteSynced | Remote: %s [synced]\n |
| 变更: %d 个（暂存 %d，未暂存 %d，未追踪 %d）\n | AIRepoChanges | Changes: %d (staged %d, unstaged %d, untracked %d)\n |
| 值 | AIManagerValue | Value |
| git diff --staged 的输出 | AIManagerStagedDiff | Output of git diff --staged |
| 远程: %s [↑%s ↓%s]\n | AIRepoRemoteAheadBehind | Remote: %s [↑%s ↓%s]\n |
| 分支要实现的功能或目的 | AIManagerFeatureDesc | Feature or purpose of the branch |
| 分支: %s\n | AIRepoBranch | Branch: %s\n |
| 最近提交:\n | AIRepoRecentCommits | Recent commits:\n |
| 根据暂存区的 diff 生成符合 Conventional Commits 规范的提交信息 | AIManagerGenerateCommitMsg | Generate commit message following Conventional Commits specification based on staged diff |
| Stash: %d 条\n | AIRepoStashCount | Stash: %d entries\n |

## AI Prompts (保持原样或改为英文)

这些是给 AI 模型看的 prompt，建议：
1. 改为英文（更通用）
2. 或保持中文（如果主要用户是中文）
3. 或根据用户语言设置动态生成

暂时保持原样，后续可以优化。
