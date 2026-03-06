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
		Description: "提交已暂存的变更。调用前应先用 get_staged_diff 查看变更内容，由 AI 自行生成提交信息，不要询问用户",
		Params: map[string]tools.ParamSchema{
			"message": {Type: "string", Description: "Conventional Commits 格式的提交信息（由 AI 根据 diff 生成，如 feat: add login page）", Required: true},
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
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("Failed to commit: %v", err)}
	}
	t.d.Refresh(ScopeFiles, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("提交成功: \"%s\"", msg)}
}

// AmendHeadTool rewrites the most recent commit message.
type AmendHeadTool struct{ d *Deps }

func NewAmendHeadTool(d *Deps) tools.Tool { return &AmendHeadTool{d} }

func (t *AmendHeadTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "amend_head",
		Description: "修改最近一次提交的信息（git commit --amend）",
		Params: map[string]tools.ParamSchema{
			"message": {Type: "string", Description: "New commit message", Required: true},
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
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("修改提交信息失败: %v", err)}
	}
	t.d.Refresh(ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已修改最新提交信息为: \"%s\"", msg)}
}

// RevertCommitTool reverts a commit by hash.
type RevertCommitTool struct{ d *Deps }

func NewRevertCommitTool(d *Deps) tools.Tool { return &RevertCommitTool{d} }

func (t *RevertCommitTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "revert_commit",
		Description: "撤销指定提交（创建一个新的反向提交）",
		Params: map[string]tools.ParamSchema{
			"hash": {Type: "string", Description: "要撤销的提交 hash", Required: true},
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
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("revert 失败: %v", err)}
	}
	t.d.Refresh(ScopeCommits, ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已 revert 提交: %s", hash)}
}

// ResetSoftTool performs a soft reset.
type ResetSoftTool struct{ d *Deps }

func NewResetSoftTool(d *Deps) tools.Tool { return &ResetSoftTool{d} }

func (t *ResetSoftTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "reset_soft",
		Description: "git reset --soft（保留变更到暂存区）",
		Params: map[string]tools.ParamSchema{
			"ref":   {Type: "string", Description: t.d.Tr.ToolTargetRefOrHash()},
			"steps": {Type: "int", Description: t.d.Tr.ToolResetSteps()},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *ResetSoftTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	ref := strParam(call.Params, "ref", "")
	if ref == "" {
		steps := intParam(call.Params, "steps", 1)
		if steps <= 0 {
			steps = 1
		}
		ref = fmt.Sprintf("HEAD~%d", steps)
	}
	if err := t.d.WorkingTree.ResetSoft(ref); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("reset --soft 失败: %v", err)}
	}
	t.d.Refresh(ScopeFiles, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("reset --soft 到 %s，变更已保留在暂存区", ref)}
}

// ResetMixedTool performs a mixed reset.
type ResetMixedTool struct{ d *Deps }

func NewResetMixedTool(d *Deps) tools.Tool { return &ResetMixedTool{d} }

func (t *ResetMixedTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "reset_mixed",
		Description: "git reset --mixed（保留变更到工作区，不暂存）",
		Params: map[string]tools.ParamSchema{
			"ref":   {Type: "string", Description: t.d.Tr.ToolTargetRefOrHash()},
			"steps": {Type: "int", Description: t.d.Tr.ToolResetSteps()},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *ResetMixedTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	ref := strParam(call.Params, "ref", "")
	if ref == "" {
		steps := intParam(call.Params, "steps", 1)
		if steps <= 0 {
			steps = 1
		}
		ref = fmt.Sprintf("HEAD~%d", steps)
	}
	if err := t.d.WorkingTree.ResetMixed(ref); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("reset --mixed 失败: %v", err)}
	}
	t.d.Refresh(ScopeFiles, ScopeCommits)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("reset --mixed 到 %s，变更已保留在工作区", ref)}
}

// CherryPickTool cherry-picks a commit by hash.
type CherryPickTool struct{ d *Deps }

func NewCherryPickTool(d *Deps) tools.Tool { return &CherryPickTool{d} }

func (t *CherryPickTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "cherry_pick",
		Description: "将指定提交 cherry-pick 到当前分支",
		Params: map[string]tools.ParamSchema{
			"hash": {Type: "string", Description: "要 cherry-pick 的提交 hash", Required: true},
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
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("cherry-pick 失败: %v", err)}
	}
	t.d.Refresh(ScopeCommits, ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已 cherry-pick: %s", hash)}
}
