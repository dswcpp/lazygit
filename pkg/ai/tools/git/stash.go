package gittools

import (
	"context"
	"fmt"

	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

// StashTool saves working tree changes to the stash.
type StashTool struct{ d *Deps }

func NewStashTool(d *Deps) tools.Tool { return &StashTool{d} }

func (t *StashTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "stash",
		Description: "储藏当前工作区变更",
		Params: map[string]tools.ParamSchema{
			"message": {Type: "string", Description: "储藏描述（可选）"},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *StashTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	msg := strParam(call.Params, "message", "AI stash")
	if err := t.d.Stash.Push(msg); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("stash 失败: %v", err)}
	}
	t.d.Refresh(ScopeFiles, ScopeStash)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已储藏变更: %s", msg)}
}

// StashPopTool pops the latest stash entry.
type StashPopTool struct{ d *Deps }

func NewStashPopTool(d *Deps) tools.Tool { return &StashPopTool{d} }

func (t *StashPopTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "stash_pop",
		Description: "恢复指定 stash 并从 stash 列表中删除",
		Params: map[string]tools.ParamSchema{
			"index": {Type: "int", Description: "stash 索引，默认 0（最近的 stash）"},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *StashPopTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	idx := intParam(call.Params, "index", 0)
	if err := t.d.Stash.Pop(idx); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("stash pop 失败: %v", err)}
	}
	t.d.Refresh(ScopeFiles, ScopeStash)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已恢复 stash[%d]", idx)}
}

// StashApplyTool applies a stash entry without removing it.
type StashApplyTool struct{ d *Deps }

func NewStashApplyTool(d *Deps) tools.Tool { return &StashApplyTool{d} }

func (t *StashApplyTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "stash_apply",
		Description: "应用指定 stash（保留 stash 条目）",
		Params: map[string]tools.ParamSchema{
			"index": {Type: "int", Description: "stash 索引，默认 0"},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *StashApplyTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	idx := intParam(call.Params, "index", 0)
	if err := t.d.Stash.Apply(idx); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("stash apply 失败: %v", err)}
	}
	t.d.Refresh(ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已应用 stash[%d]（条目保留）", idx)}
}

// StashDropTool deletes a stash entry.
type StashDropTool struct{ d *Deps }

func NewStashDropTool(d *Deps) tools.Tool { return &StashDropTool{d} }

func (t *StashDropTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "stash_drop",
		Description: "删除指定 stash 条目",
		Params: map[string]tools.ParamSchema{
			"index": {Type: "int", Description: "stash 索引，默认 0"},
		},
		Permission: tools.PermDestructive,
	}
}

func (t *StashDropTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	idx := intParam(call.Params, "index", 0)
	if err := t.d.Stash.Drop(idx); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("stash drop 失败: %v", err)}
	}
	t.d.Refresh(ScopeStash)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已删除 stash[%d]", idx)}
}
