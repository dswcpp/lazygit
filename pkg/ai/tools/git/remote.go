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

// PushTool pushes the current branch to its upstream.
type PushTool struct{ d *Deps }

func NewPushTool(d *Deps) tools.Tool { return &PushTool{d} }

func (t *PushTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "push",
		Description: "推送当前分支到远程（git push）",
		Params: map[string]tools.ParamSchema{
			"force": {Type: "bool", Description: "是否强制推送（默认 false）"},
		},
		Permission: tools.PermWriteRemote,
	}
}

func (t *PushTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	force := boolParam(call.Params, "force", false)
	cmdObj, err := t.d.Sync.PushCmdObj(nil, git_commands.PushOpts{Force: force})
	if err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("push 配置错误: %v", err)}
	}
	if err := cmdObj.Run(); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("push 失败: %v（请确认远程配置和认证）", err)}
	}
	t.d.Refresh(ScopeBranches, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: "push 成功"}
}
