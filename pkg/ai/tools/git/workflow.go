package gittools

import (
	"context"

	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

// AbortOperationTool aborts an in-progress rebase, merge, or cherry-pick.
type AbortOperationTool struct{ d *Deps }

func NewAbortOperationTool(d *Deps) tools.Tool { return &AbortOperationTool{d} }

func (t *AbortOperationTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "abort_operation",
		Description: t.d.Tr.ToolAbortOperationDesc(),
		Params: map[string]tools.ParamSchema{
			"type": {Type: "string", Description: t.d.Tr.ToolAbortOperationTypeParam(), Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *AbortOperationTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	opType := strParam(call.Params, "type", "")
	switch opType {
	case "rebase", "merge", "cherry-pick":
		// valid
	default:
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolUnknownOperationType(opType, "rebase | merge | cherry-pick")}
	}
	if err := t.d.Rebase.GenericMergeOrRebaseAction(opType, "abort"); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolAbortFailed(opType, err)}
	}
	t.d.Refresh(ScopeFiles, ScopeCommits, ScopeBranches)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolAbortSuccess(opType)}
}

// ContinueOperationTool continues a paused rebase or merge after resolving conflicts.
type ContinueOperationTool struct{ d *Deps }

func NewContinueOperationTool(d *Deps) tools.Tool { return &ContinueOperationTool{d} }

func (t *ContinueOperationTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "continue_operation",
		Description: t.d.Tr.ToolContinueOperationDesc(),
		Params: map[string]tools.ParamSchema{
			"type": {Type: "string", Description: t.d.Tr.ToolContinueOperationTypeParam()},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *ContinueOperationTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	opType := strParam(call.Params, "type", "rebase")
	switch opType {
	case "rebase", "merge":
		// valid
	default:
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolUnknownOperationType(opType, "rebase | merge")}
	}
	if err := t.d.Rebase.GenericMergeOrRebaseAction(opType, "continue"); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolContinueFailed(opType, err)}
	}
	t.d.Refresh(ScopeFiles, ScopeCommits, ScopeBranches)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolContinueSuccess(opType)}
}
