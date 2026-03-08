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
		Description: t.d.Tr.ToolCheckoutDesc(),
		Params: map[string]tools.ParamSchema{
			"name": {Type: "string", Description: t.d.Tr.ToolCheckoutNameParam(), Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *CheckoutTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	name := strParam(call.Params, "name", "")
	if name == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingNameParam()}
	}
	if err := t.d.Branch.Checkout(name, git_commands.CheckoutOptions{}); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolCheckoutFailed(err)}
	}
	t.d.Refresh(ScopeFiles, ScopeBranches, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolCheckoutSuccess(name)}
}

// CreateBranchTool creates a new local branch.
type CreateBranchTool struct{ d *Deps }

func NewCreateBranchTool(d *Deps) tools.Tool { return &CreateBranchTool{d} }

func (t *CreateBranchTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "create_branch",
		Description: t.d.Tr.ToolCreateBranchDesc(),
		Params: map[string]tools.ParamSchema{
			"name":     {Type: "string", Description: t.d.Tr.ToolCreateBranchNameParam(), Required: true},
			"base":     {Type: "string", Description: t.d.Tr.ToolCreateBranchBaseParam()},
			"checkout": {Type: "bool", Description: t.d.Tr.ToolCreateBranchCheckoutParam()},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *CreateBranchTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	name := strParam(call.Params, "name", "")
	if name == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingNameParam()}
	}
	base := strParam(call.Params, "base", "HEAD")
	doCheckout := boolParam(call.Params, "checkout", true)
	if doCheckout {
		if err := t.d.Branch.New(name, base); err != nil {
			return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolCreateBranchFailed(err)}
		}
		t.d.Refresh(ScopeFiles, ScopeBranches, ScopeCommits)
		return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolCreateBranchSuccess(name, base)}
	}
	if err := t.d.Branch.NewWithoutCheckout(name, base); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolCreateBranchFailed(err)}
	}
	t.d.Refresh(ScopeBranches)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolCreateBranchNoCheckoutSuccess(name, base)}
}

// DeleteBranchTool deletes a local branch.
type DeleteBranchTool struct{ d *Deps }

func NewDeleteBranchTool(d *Deps) tools.Tool { return &DeleteBranchTool{d} }

func (t *DeleteBranchTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "delete_branch",
		Description: t.d.Tr.ToolDeleteBranchDesc(),
		Params: map[string]tools.ParamSchema{
			"name":  {Type: "string", Description: t.d.Tr.ToolBranchName(), Required: true},
			"force": {Type: "bool", Description: t.d.Tr.ToolDeleteBranchForceParam()},
		},
		Permission: tools.PermDestructive,
	}
}

func (t *DeleteBranchTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	name := strParam(call.Params, "name", "")
	if name == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingNameParam()}
	}
	force := boolParam(call.Params, "force", false)
	if err := t.d.Branch.LocalDelete([]string{name}, force); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolDeleteBranchFailed(err)}
	}
	t.d.Refresh(ScopeBranches)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolDeleteBranchSuccess(name)}
}

// RenameBranchTool renames a local branch.
type RenameBranchTool struct{ d *Deps }

func NewRenameBranchTool(d *Deps) tools.Tool { return &RenameBranchTool{d} }

func (t *RenameBranchTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "rename_branch",
		Description: t.d.Tr.ToolRenameBranchDesc(),
		Params: map[string]tools.ParamSchema{
			"old":  {Type: "string", Description: t.d.Tr.ToolRenameBranchOldParam(), Required: true},
			"name": {Type: "string", Description: t.d.Tr.ToolBranchName(), Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *RenameBranchTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	old := strParam(call.Params, "old", "")
	name := strParam(call.Params, "name", "")
	if old == "" || name == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingOldOrNameParam()}
	}
	if err := t.d.Branch.Rename(old, name); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolRenameFailed(err)}
	}
	t.d.Refresh(ScopeBranches)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolRenameBranchSuccess(old, name)}
}

// MergeBranchTool merges a branch into the current branch.
type MergeBranchTool struct{ d *Deps }

func NewMergeBranchTool(d *Deps) tools.Tool { return &MergeBranchTool{d} }

func (t *MergeBranchTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "merge_branch",
		Description: t.d.Tr.ToolMergeBranchDesc(),
		Params: map[string]tools.ParamSchema{
			"name": {Type: "string", Description: t.d.Tr.ToolMergeBranchNameParam(), Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *MergeBranchTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	name := strParam(call.Params, "name", "")
	if name == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingNameParam()}
	}
	if err := t.d.Branch.Merge(name, git_commands.MERGE_VARIANT_REGULAR); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMergeBranchFailed(err)}
	}
	t.d.Refresh(ScopeCommits, ScopeBranches, ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolMergeBranchSuccess(name)}
}

// RebaseBranchTool rebases the current branch onto another branch.
type RebaseBranchTool struct{ d *Deps }

func NewRebaseBranchTool(d *Deps) tools.Tool { return &RebaseBranchTool{d} }

func (t *RebaseBranchTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "rebase_branch",
		Description: t.d.Tr.ToolRebaseBranchDesc(),
		Params: map[string]tools.ParamSchema{
			"target": {Type: "string", Description: t.d.Tr.ToolRebaseBranchTargetParam(), Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *RebaseBranchTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	target := strParam(call.Params, "target", "")
	if target == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingTargetParam()}
	}
	if err := t.d.Rebase.RebaseBranch(target); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolRebaseBranchFailed(err)}
	}
	t.d.Refresh(ScopeCommits, ScopeBranches, ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolRebasedTo(target)}
}

// ensure fmt is used
var _ = fmt.Sprintf
