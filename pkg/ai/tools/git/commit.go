package gittools

import (
	"context"
	"fmt"

	"github.com/dswcpp/lazygit/pkg/ai/tools"
	"github.com/dswcpp/lazygit/pkg/commands/models"
)

// CommitTool creates a new commit.
type CommitTool struct{ d *Deps }

func NewCommitTool(d *Deps) tools.Tool { return &CommitTool{d} }

func (t *CommitTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "commit",
		Description: t.d.Tr.ToolCommitDesc(),
		Params: map[string]tools.ParamSchema{
			"message": {Type: "string", Description: t.d.Tr.ToolCommitMsgParam(), Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *CommitTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	msg := strParam(call.Params, "message", "")
	if msg == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingMessageParam()}
	}
	if err := t.d.Commit.CommitCmdObj(msg, "", false).Run(); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolCommitFailed(err)}
	}
	t.d.Refresh(ScopeFiles, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolCommitSuccess(msg)}
}

// AmendHeadTool rewrites the most recent commit message.
type AmendHeadTool struct{ d *Deps }

func NewAmendHeadTool(d *Deps) tools.Tool { return &AmendHeadTool{d} }

func (t *AmendHeadTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "amend_head",
		Description: t.d.Tr.ToolAmendHeadDesc(),
		Params: map[string]tools.ParamSchema{
			"message": {Type: "string", Description: t.d.Tr.ToolAmendMsgParam(), Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *AmendHeadTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	msg := strParam(call.Params, "message", "")
	if msg == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingMessageParam()}
	}
	if err := t.d.Commit.RewordLastCommit(msg, "").Run(); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolAmendFailed(err)}
	}
	t.d.Refresh(ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolAmendSuccess(msg)}
}

// RevertCommitTool reverts a commit by hash.
type RevertCommitTool struct{ d *Deps }

func NewRevertCommitTool(d *Deps) tools.Tool { return &RevertCommitTool{d} }

func (t *RevertCommitTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "revert_commit",
		Description: t.d.Tr.ToolRevertCommitDesc(),
		Params: map[string]tools.ParamSchema{
			"hash": {Type: "string", Description: t.d.Tr.ToolRevertHashParam(), Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *RevertCommitTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	hash := strParam(call.Params, "hash", "")
	if hash == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingHashParam()}
	}
	if err := t.d.Commit.Revert([]string{hash}, false); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolRevertFailed(err)}
	}
	t.d.Refresh(ScopeCommits, ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolRevertSuccess(hash)}
}

// ResetSoftTool performs a soft reset.
type ResetSoftTool struct{ d *Deps }

func NewResetSoftTool(d *Deps) tools.Tool { return &ResetSoftTool{d} }

func (t *ResetSoftTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "reset_soft",
		Description: t.d.Tr.ToolResetSoftDesc(),
		Params: map[string]tools.ParamSchema{
			"ref":   {Type: "string", Description: t.d.Tr.ToolTargetRefOrHash()},
			"steps": {Type: "int", Description: t.d.Tr.ToolResetSteps()},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *ResetSoftTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	ref := resolveRef(call.Params)
	if err := t.d.WorkingTree.ResetSoft(ref); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolResetSoftFailed(err)}
	}
	t.d.Refresh(ScopeFiles, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolResetSoftSuccess(ref)}
}

// ResetMixedTool performs a mixed reset.
type ResetMixedTool struct{ d *Deps }

func NewResetMixedTool(d *Deps) tools.Tool { return &ResetMixedTool{d} }

func (t *ResetMixedTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "reset_mixed",
		Description: t.d.Tr.ToolResetMixedDesc(),
		Params: map[string]tools.ParamSchema{
			"ref":   {Type: "string", Description: t.d.Tr.ToolTargetRefOrHash()},
			"steps": {Type: "int", Description: t.d.Tr.ToolResetSteps()},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *ResetMixedTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	ref := resolveRef(call.Params)
	if err := t.d.WorkingTree.ResetMixed(ref); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolResetMixedFailed(err)}
	}
	t.d.Refresh(ScopeFiles, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolResetMixedSuccess(ref)}
}

// ResetHardTool performs a hard reset — all uncommitted changes are discarded.
type ResetHardTool struct{ d *Deps }

func NewResetHardTool(d *Deps) tools.Tool { return &ResetHardTool{d} }

func (t *ResetHardTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "reset_hard",
		Description: t.d.Tr.ToolResetHardDesc(),
		Params: map[string]tools.ParamSchema{
			"ref":   {Type: "string", Description: t.d.Tr.ToolTargetRefOrHash()},
			"steps": {Type: "int", Description: t.d.Tr.ToolResetSteps()},
		},
		Permission: tools.PermDestructive,
	}
}

func (t *ResetHardTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	ref := resolveRef(call.Params)
	if err := t.d.WorkingTree.ResetHard(ref); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolResetHardFailed(err)}
	}
	t.d.Refresh(ScopeFiles, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolResetHardSuccess(ref)}
}

// CherryPickTool cherry-picks a commit by hash.
type CherryPickTool struct{ d *Deps }

func NewCherryPickTool(d *Deps) tools.Tool { return &CherryPickTool{d} }

func (t *CherryPickTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "cherry_pick",
		Description: t.d.Tr.ToolCherryPickDesc(),
		Params: map[string]tools.ParamSchema{
			"hash": {Type: "string", Description: t.d.Tr.ToolCherryPickHashParam(), Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *CherryPickTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	hash := strParam(call.Params, "hash", "")
	if hash == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingHashParam()}
	}
	commit := models.NewCommit(t.d.GetHashPool(), models.NewCommitOpts{Hash: hash})
	if err := t.d.Rebase.CherryPickCommits([]*models.Commit{commit}); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolCherryPickFailed(err)}
	}
	t.d.Refresh(ScopeCommits, ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolCherryPickSuccess(hash)}
}

// resolveRef builds a git ref from params: uses "ref" if provided, else "HEAD~{steps}".
func resolveRef(params map[string]any) string {
	if ref := strParam(params, "ref", ""); ref != "" {
		return ref
	}
	steps := intParam(params, "steps", 1)
	if steps <= 0 {
		steps = 1
	}
	return fmt.Sprintf("HEAD~%d", steps)
}
