package ai

import (
	"context"
	"fmt"

	"github.com/dswcpp/lazygit/pkg/ai/agent"
	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/repocontext"
	"github.com/dswcpp/lazygit/pkg/ai/skills"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
	"github.com/dswcpp/lazygit/pkg/config"
)

// Manager is the top-level facade for all AI features.
// GUI components depend on *Manager instead of *Client.
//
// Migration note: Manager is introduced alongside the existing *Client.
// Components migrate to Manager incrementally; *Client remains available
// during the transition period and is accessible via Manager.LegacyClient().
type Manager struct {
	prov       provider.Provider
	registry   *tools.Registry
	ctxBuilder repocontext.Builder
	skillMap   map[string]skills.Skill
}

// agentSystemPrompt is the default system prompt for the chat agent.
const agentSystemPrompt = `你是 lazygit 的内置 AI Agent，可以直接操控 Git 仓库。
用户对你说话就像对同事下指令：你要主动思考、自主执行，而不是给出建议让用户手动操作。
在回复中嵌入工具调用代码块来执行操作，格式如下：

` + "```tool" + `
{"name": "工具名", "params": {"参数": "值"}}
` + "```" + `

每次只调用一个工具，等待结果后再决定下一步。
如果任务已完成，用自然语言告知用户结果，不要再调用工具。

## 提交规范

需要执行提交（commit）时，**绝对不要询问用户提交信息**，按以下步骤自动完成：
1. 调用 get_staged_diff 获取暂存区变更内容
2. 根据 diff 内容自行生成符合 Conventional Commits 规范的提交信息（如 feat: ...、fix: ...、refactor: ... 等）
3. 直接调用 commit 工具完成提交

## 分支命名规范

需要新建分支时，**绝对不要询问用户分支名**，根据用户描述的功能/目的自动生成：
- 格式：类型/简短描述，用 kebab-case，如 feature/user-login、fix/null-pointer、chore/update-deps
- 类型参考：feature（新功能）、fix（修复）、refactor（重构）、chore（杂项）、docs（文档）、test（测试）
- 直接调用 create_branch 工具，默认同时切换到新分支`

// NewManager creates a Manager from the user's active AI profile.
// Returns nil, nil when AI is disabled or no active profile is configured.
func NewManager(cfg config.AIConfig, ctxBuilder repocontext.Builder) (*Manager, error) {
	prov, err := provider.NewFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	if prov == nil {
		return nil, nil
	}
	m := &Manager{
		prov:       prov,
		registry:   tools.NewRegistry(),
		ctxBuilder: ctxBuilder,
		skillMap:   make(map[string]skills.Skill),
	}
	// Register built-in skills
	for _, sk := range []skills.Skill{
		skills.NewCommitMsgSkill(),
		skills.NewBranchNameSkill(),
		skills.NewPRDescSkill(),
		skills.NewCodeReviewSkill(),
		skills.NewShellCmdSkill(),
	} {
		m.skillMap[sk.Name()] = sk
	}
	return m, nil
}

// Provider returns the underlying provider for direct use by skills and the agent.
func (m *Manager) Provider() provider.Provider { return m.prov }

// Registry returns the tool registry.
// The GUI layer registers git tools here during initialisation.
func (m *Manager) Registry() *tools.Registry { return m.registry }

// SetContextBuilder injects the repository context builder.
// Called by the GUI layer after the Manager is created, once the GUI model is available.
func (m *Manager) SetContextBuilder(b repocontext.Builder) {
	m.ctxBuilder = b
}

// RepoContext builds a current snapshot of the repository state.
func (m *Manager) RepoContext() repocontext.RepoContext {
	if m.ctxBuilder == nil {
		return repocontext.RepoContext{}
	}
	return m.ctxBuilder.Build()
}

// RunSkill executes a named skill and returns the output.
// Returns an error if the skill is not registered or the provider call fails.
func (m *Manager) RunSkill(ctx context.Context, name string, extra map[string]any) (skills.Output, error) {
	sk, ok := m.skillMap[name]
	if !ok {
		return skills.Output{}, fmt.Errorf("unknown skill: %q", name)
	}
	return sk.Execute(ctx, m.prov, skills.Input{
		RepoCtx: m.RepoContext(),
		Extra:   extra,
	})
}

// NewAgent creates a new Agent backed by this Manager's provider and tool registry.
// Each call creates an independent session.
// The systemPrompt parameter may be empty to use the default agent prompt.
func (m *Manager) NewAgent(systemPrompt string, confirmFn agent.ConfirmFunc) *agent.Agent {
	if systemPrompt == "" {
		systemPrompt = agentSystemPrompt
	}
	// Inject tool list into system prompt
	toolSection := m.registry.SystemPromptSection(tools.PermDestructive)
	if toolSection != "" {
		systemPrompt += "\n\n" + toolSection
	}
	session := agent.NewSession(systemPrompt)
	return agent.NewAgent(m.prov, m.registry, session, confirmFn)
}

// NewTwoPhaseAgent 创建两阶段 Agent（规划 → 聊天确认 → 执行）。
//
// 规划阶段只使用只读工具（PermReadOnly），写操作工具在执行阶段才会被访问。
// skillTools 为可选的 SkillTool 列表（通过 tools.NewSkillTool 构建），
// 会被注入到规划阶段的只读注册表中，供 LLM 预计算提交信息、分支名等。
//
// 用户通过聊天输入 Y/N 或补充说明来完成确认，无需弹窗。
func (m *Manager) NewTwoPhaseAgent(skillTools []tools.Tool) *agent.TwoPhaseAgent {
	// 构建只读注册表：从完整注册表中筛选只读工具，再加入 SkillTool
	readReg := tools.NewRegistry()
	for _, t := range m.registry.ByMaxPermission(tools.PermReadOnly) {
		readReg.Register(t)
	}
	for _, st := range skillTools {
		readReg.Register(st)
	}

	session := agent.NewSession("") // TwoPhaseAgent 使用自己的 system prompt，此处留空
	return agent.NewTwoPhaseAgent(m.prov, m.registry, readReg, session)
}

// DefaultSkillTools 从已注册的 Skill 构建默认 SkillTool 列表，
// 用于注入规划阶段的只读注册表，让规划 LLM 能预计算提交信息、分支名等。
func (m *Manager) DefaultSkillTools() []tools.Tool {
	repoCtxFn := func() repocontext.RepoContext { return m.RepoContext() }
	out := make([]tools.Tool, 0, 2)

	if sk, ok := m.skillMap["commit_msg"]; ok {
		out = append(out, tools.NewSkillTool(
			sk, m.prov, repoCtxFn,
			"根据暂存区的 diff 生成符合 Conventional Commits 规范的提交信息",
			map[string]tools.ParamSchema{
				"diff": {Type: "string", Required: true, Description: "git diff --staged 的输出"},
			},
		))
	}
	if sk, ok := m.skillMap["branch_name"]; ok {
		out = append(out, tools.NewSkillTool(
			sk, m.prov, repoCtxFn,
			"根据功能描述生成合适的 Git 分支名（kebab-case，带类型前缀）",
			map[string]tools.ParamSchema{
				"description": {Type: "string", Required: true, Description: "分支要实现的功能或目的"},
			},
		))
	}
	return out
}

// LegacyClient returns a *Client that wraps this Manager's provider,
// allowing existing code that depends on *Client to work unchanged
// during the incremental migration.
func (m *Manager) LegacyClient() *Client {
	return &Client{provider: &legacyProviderAdapter{p: m.prov}}
}
