package gittools

import (
	"context"

	"github.com/dswcpp/lazygit/pkg/ai/tools"
	"github.com/dswcpp/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/gocui"
)

// PullTool pulls from remote (fetch + merge into current branch).
type PullTool struct{ d *Deps }

func NewPullTool(d *Deps) tools.Tool { return &PullTool{d} }

func (t *PullTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "pull",
		Description: t.d.Tr.ToolPullDesc(),
		Params: map[string]tools.ParamSchema{
			"remote": {Type: "string", Description: t.d.Tr.ToolPullRemoteParam()},
			"branch": {Type: "string", Description: t.d.Tr.ToolPullBranchParam()},
		},
		Permission: tools.PermWriteRemote,
	}
}

func (t *PullTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	remote := strParam(call.Params, "remote", "")
	branch := strParam(call.Params, "branch", "")
	// FakeTask: AI tool context has no gocui event loop; remote must support
	// keyring/SSH auth (credential prompt cannot be shown).
	task := gocui.NewFakeTask()
	defer task.Done()
	if err := t.d.Sync.Pull(task, git_commands.PullOptions{
		RemoteName: remote,
		BranchName: branch,
	}); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolPullFailed(err)}
	}
	t.d.Refresh(ScopeBranches, ScopeCommits, ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolPullSuccess()}
}

// FetchTool fetches from all remotes in the background.
type FetchTool struct{ d *Deps }

func NewFetchTool(d *Deps) tools.Tool { return &FetchTool{d} }

func (t *FetchTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "fetch",
		Description: t.d.Tr.ToolFetchDesc(),
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermWriteRemote,
	}
}

func (t *FetchTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	if err := t.d.Sync.FetchBackground(); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolFetchFailed(err)}
	}
	t.d.Refresh(ScopeBranches, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolFetchSuccess()}
}

// PushTool pushes the current branch to its upstream (normal push only).
type PushTool struct{ d *Deps }

func NewPushTool(d *Deps) tools.Tool { return &PushTool{d} }

func (t *PushTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "push",
		Description: t.d.Tr.ToolPushDesc(),
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
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolPushFailed(err)}
	}
	t.d.Refresh(ScopeBranches, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolPushSuccess()}
}

// PushForceTool force-pushes using --force-with-lease (safer than --force).
type PushForceTool struct{ d *Deps }

func NewPushForceTool(d *Deps) tools.Tool { return &PushForceTool{d} }

func (t *PushForceTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "push_force",
		Description: t.d.Tr.ToolPushForceDesc(),
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
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolPushForceFailed(err)}
	}
	t.d.Refresh(ScopeBranches, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolPushForceSuccess()}
}
