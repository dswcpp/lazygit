package gittools

import (
	"context"
	"fmt"

	"github.com/dswcpp/lazygit/pkg/ai/tools"
	"github.com/dswcpp/lazygit/pkg/commands/git_commands"
)

// FetchTool fetches from all remotes in the background.
type FetchTool struct{ d *Deps }

func NewFetchTool(d *Deps) tools.Tool { return &FetchTool{d} }

func (t *FetchTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "fetch",
		Description: "从远程拉取最新引用（git fetch）",
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermWriteRemote,
	}
}

func (t *FetchTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	if err := t.d.Sync.FetchBackground(); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("fetch 失败: %v", err)}
	}
	t.d.Refresh(ScopeBranches, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: "fetch 完成"}
}

// PushTool pushes the current branch to its upstream (normal push only).
type PushTool struct{ d *Deps }

func NewPushTool(d *Deps) tools.Tool { return &PushTool{d} }

func (t *PushTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "push",
		Description: "推送当前分支到远程（git push）。如需强制推送请使用 push_force 工具",
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermWriteRemote,
	}
}

func (t *PushTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	cmdObj, err := t.d.Sync.PushCmdObj(nil, git_commands.PushOpts{})
	if err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolPushConfigError(err)}
	}
	if err := cmdObj.Run(); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("push 失败: %v（请确认远程配置和认证）", err)}
	}
	t.d.Refresh(ScopeBranches, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: "push 成功"}
}

// PushForceTool force-pushes using --force-with-lease (safer than --force).
type PushForceTool struct{ d *Deps }

func NewPushForceTool(d *Deps) tools.Tool { return &PushForceTool{d} }

func (t *PushForceTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "push_force",
		Description: "强制推送（git push --force-with-lease）：若远程有未拉取的提交则自动中止，比 --force 更安全。仍会覆盖远程历史，需谨慎",
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermDestructive,
	}
}

func (t *PushForceTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	cmdObj, err := t.d.Sync.PushCmdObj(nil, git_commands.PushOpts{ForceWithLease: true})
	if err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolPushConfigError(err)}
	}
	if err := cmdObj.Run(); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("force push 失败: %v", err)}
	}
	t.d.Refresh(ScopeBranches, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: "force push（--force-with-lease）成功"}
}
