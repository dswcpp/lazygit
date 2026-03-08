package gittools

import (
	"context"

	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

// StashTool saves working tree changes to the stash.
type StashTool struct{ d *Deps }

func NewStashTool(d *Deps) tools.Tool { return &StashTool{d} }

func (t *StashTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "stash",
		Description: t.d.Tr.ToolStashDesc(),
		Params: map[string]tools.ParamSchema{
			"message": {Type: "string", Description: t.d.Tr.ToolStashMsgParam()},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *StashTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	msg := strParam(call.Params, "message", "WIP")
	if err := t.d.Stash.Push(msg); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolStashFailed(err)}
	}
	t.d.Refresh(ScopeFiles, ScopeStash)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolStashSuccess(msg)}
}

// StashPopTool pops the latest stash entry.
type StashPopTool struct{ d *Deps }

func NewStashPopTool(d *Deps) tools.Tool { return &StashPopTool{d} }

func (t *StashPopTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "stash_pop",
		Description: t.d.Tr.ToolStashPopDesc(),
		Params: map[string]tools.ParamSchema{
			"index": {Type: "int", Description: t.d.Tr.ToolStashIndex()},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *StashPopTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	idx := intParam(call.Params, "index", 0)
	if err := t.d.Stash.Pop(idx); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolStashPopFailed(err)}
	}
	t.d.Refresh(ScopeFiles, ScopeStash)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolStashPopSuccess(idx)}
}

// StashApplyTool applies a stash entry without removing it.
type StashApplyTool struct{ d *Deps }

func NewStashApplyTool(d *Deps) tools.Tool { return &StashApplyTool{d} }

func (t *StashApplyTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "stash_apply",
		Description: t.d.Tr.ToolStashApplyDesc(),
		Params: map[string]tools.ParamSchema{
			"index": {Type: "int", Description: t.d.Tr.ToolStashIndex()},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *StashApplyTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	idx := intParam(call.Params, "index", 0)
	if err := t.d.Stash.Apply(idx); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolStashApplyFailed(err)}
	}
	t.d.Refresh(ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolStashApplySuccess(idx)}
}

// StashDropTool deletes a stash entry.
type StashDropTool struct{ d *Deps }

func NewStashDropTool(d *Deps) tools.Tool { return &StashDropTool{d} }

func (t *StashDropTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "stash_drop",
		Description: t.d.Tr.ToolStashDropDesc(),
		Params: map[string]tools.ParamSchema{
			"index": {Type: "int", Description: t.d.Tr.ToolStashIndex()},
		},
		Permission: tools.PermDestructive,
	}
}

func (t *StashDropTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	idx := intParam(call.Params, "index", 0)
	if err := t.d.Stash.Drop(idx); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolStashDropFailed(err)}
	}
	t.d.Refresh(ScopeStash)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolStashDropSuccess(idx)}
}
