package gittools

import (
	"context"
	"fmt"

	"github.com/dswcpp/lazygit/pkg/ai/tools"
	"github.com/dswcpp/lazygit/pkg/commands/git_commands"
)

// CheckoutTool switches to an existing branch.
type CheckoutTool struct{ d *Deps }

func NewCheckoutTool(d *Deps) tools.Tool { return &CheckoutTool{d} }

func (t *CheckoutTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "checkout",
		Description: "切换到指定分支",
		Params: map[string]tools.ParamSchema{
			"name": {Type: "string", Description: "分支名", Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *CheckoutTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	name := strParam(call.Params, "name", "")
	if name == "" {
		return tools.ToolResult{CallID: call.ID, Output: "缺少 name 参数"}
	}
	if err := t.d.Branch.Checkout(name, git_commands.CheckoutOptions{}); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("切换分支失败: %v", err)}
	}
	t.d.Refresh(ScopeFiles, ScopeBranches, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已切换到: %s", name)}
}

// CreateBranchTool creates a new local branch.
type CreateBranchTool struct{ d *Deps }

func NewCreateBranchTool(d *Deps) tools.Tool { return &CreateBranchTool{d} }

func (t *CreateBranchTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "create_branch",
		Description: "创建新分支",
		Params: map[string]tools.ParamSchema{
			"name": {Type: "string", Description: "新分支名", Required: true},
			"base": {Type: "string", Description: "基础 ref（默认 HEAD）"},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *CreateBranchTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	name := strParam(call.Params, "name", "")
	if name == "" {
		return tools.ToolResult{CallID: call.ID, Output: "缺少 name 参数"}
	}
	base := strParam(call.Params, "base", "HEAD")
	if err := t.d.Branch.New(name, base); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("创建分支失败: %v", err)}
	}
	t.d.Refresh(ScopeBranches)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已创建分支 %s（基于 %s）", name, base)}
}

// DeleteBranchTool deletes a local branch.
type DeleteBranchTool struct{ d *Deps }

func NewDeleteBranchTool(d *Deps) tools.Tool { return &DeleteBranchTool{d} }

func (t *DeleteBranchTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "delete_branch",
		Description: "删除本地分支（分支须已合并）",
		Params: map[string]tools.ParamSchema{
			"name": {Type: "string", Description: "分支名", Required: true},
		},
		Permission: tools.PermDestructive,
	}
}

func (t *DeleteBranchTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	name := strParam(call.Params, "name", "")
	if name == "" {
		return tools.ToolResult{CallID: call.ID, Output: "缺少 name 参数"}
	}
	if err := t.d.Branch.LocalDelete([]string{name}, false); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("删除分支失败: %v（请确认分支已合并）", err)}
	}
	t.d.Refresh(ScopeBranches)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已删除分支: %s", name)}
}

// RenameBranchTool renames a local branch.
type RenameBranchTool struct{ d *Deps }

func NewRenameBranchTool(d *Deps) tools.Tool { return &RenameBranchTool{d} }

func (t *RenameBranchTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "rename_branch",
		Description: "重命名本地分支",
		Params: map[string]tools.ParamSchema{
			"old":  {Type: "string", Description: "当前分支名", Required: true},
			"name": {Type: "string", Description: "新分支名", Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *RenameBranchTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	old := strParam(call.Params, "old", "")
	name := strParam(call.Params, "name", "")
	if old == "" || name == "" {
		return tools.ToolResult{CallID: call.ID, Output: "缺少 old 或 name 参数"}
	}
	if err := t.d.Branch.Rename(old, name); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("重命名失败: %v", err)}
	}
	t.d.Refresh(ScopeBranches)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已将分支 %s 重命名为 %s", old, name)}
}

// MergeBranchTool merges a branch into the current branch.
type MergeBranchTool struct{ d *Deps }

func NewMergeBranchTool(d *Deps) tools.Tool { return &MergeBranchTool{d} }

func (t *MergeBranchTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "merge_branch",
		Description: "将指定分支合并到当前分支",
		Params: map[string]tools.ParamSchema{
			"name": {Type: "string", Description: "要合并的分支名", Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *MergeBranchTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	name := strParam(call.Params, "name", "")
	if name == "" {
		return tools.ToolResult{CallID: call.ID, Output: "缺少 name 参数"}
	}
	if err := t.d.Branch.Merge(name, git_commands.MERGE_VARIANT_REGULAR); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("merge 失败: %v", err)}
	}
	t.d.Refresh(ScopeCommits, ScopeBranches, ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已将 %s 合并到当前分支", name)}
}
