package gittools

import "github.com/dswcpp/lazygit/pkg/ai/tools"

// RegisterAll registers every built-in git tool into the given registry.
// Call this once during GUI initialisation after building a Deps instance.
func RegisterAll(d *Deps, r *tools.Registry) {
	for _, t := range []tools.Tool{
		// Read-only queries
		NewGetStatusTool(d),
		NewGetStagedDiffTool(d),
		NewGetDiffTool(d),
		NewGetFileDiffTool(d),
		NewGetLogTool(d),
		NewGetBranchesTool(d),
		NewGetStashListTool(d),
		NewGetRemotesTool(d),
		NewGetTagsTool(d),
		NewGetCommitDiffTool(d),

		// Staging
		NewStageAllTool(d),
		NewStageFileTool(d),
		NewUnstageAllTool(d),
		NewUnstageFileTool(d),
		NewDiscardFileTool(d),

		// Commits
		NewCommitTool(d),
		NewAmendHeadTool(d),
		NewRevertCommitTool(d),
		NewResetSoftTool(d),
		NewResetMixedTool(d),
		NewCherryPickTool(d),

		// Branches
		NewCheckoutTool(d),
		NewCreateBranchTool(d),
		NewDeleteBranchTool(d),
		NewRenameBranchTool(d),
		NewMergeBranchTool(d),
		NewRebaseBranchTool(d),

		// Stash
		NewStashTool(d),
		NewStashPopTool(d),
		NewStashApplyTool(d),
		NewStashDropTool(d),

		// Tags
		NewCreateTagTool(d),
		NewDeleteTagTool(d),

		// Remote
		NewFetchTool(d),
		NewPushTool(d),
		NewPushForceTool(d),
	} {
		r.Register(t)
	}
}
