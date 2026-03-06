package helpers

import (
	"github.com/dswcpp/lazygit/pkg/ai"
	gittools "github.com/dswcpp/lazygit/pkg/ai/tools/git"
	"github.com/dswcpp/lazygit/pkg/commands/models"
	"github.com/dswcpp/lazygit/pkg/gui/types"
	"github.com/dswcpp/lazygit/pkg/utils"
)

// RegisterGitTools builds a gittools.Deps from the live GUI model and registers
// all built-in git tools into the Manager's Registry.
// Safe to call multiple times (e.g. when resetHelpersAndControllers re-runs);
// clears the registry first to prevent duplicate-registration panics.
func RegisterGitTools(c *HelperCommon, mgr *ai.Manager) {
	mgr.Registry().Clear()
	deps := buildGitToolDeps(c)
	gittools.RegisterAll(deps, mgr.Registry(), mgr.Provider())
}

// buildGitToolDeps constructs a Deps that reads from the live GUI model and
// dispatches git commands through the existing git_commands layer.
func buildGitToolDeps(c *HelperCommon) *gittools.Deps {
	return &gittools.Deps{
		// Git command groups
		WorkingTree: c.Git().WorkingTree,
		Branch:      c.Git().Branch,
		Commit:      c.Git().Commit,
		Stash:       c.Git().Stash,
		Tag:         c.Git().Tag,
		Sync:        c.Git().Sync,
		Diff:        c.Git().Diff,
		Rebase:      c.Git().Rebase,

		// Live model readers — closures capture c so they always see current state
		GetFiles:            func() []*models.File { return c.Model().Files },
		GetBranches:         func() []*models.Branch { return c.Model().Branches },
		GetCommits:          func() []*models.Commit { return c.Model().Commits },
		GetStashEntries:     func() []*models.StashEntry { return c.Model().StashEntries },
		GetRemotes:          func() []*models.Remote { return c.Model().Remotes },
		GetTags:             func() []*models.Tag { return c.Model().Tags },
		GetCheckedOutBranch: func() string { return c.Model().CheckedOutBranch },
		GetWorkingTreeState: func() models.WorkingTreeState {
			return c.Model().WorkingTreeStateAtLastCommitRefresh
		},
		GetHashPool: func() *utils.StringPool { return c.Model().HashPool },

		// Refresh: map RefreshScope → types.RefreshableView and trigger ASYNC refresh
		Refresh: func(scopes ...gittools.RefreshScope) {
			views := make([]types.RefreshableView, 0, len(scopes))
			for _, s := range scopes {
				switch s {
				case gittools.ScopeFiles:
					views = append(views, types.FILES)
				case gittools.ScopeBranches:
					views = append(views, types.BRANCHES)
				case gittools.ScopeCommits:
					views = append(views, types.COMMITS)
				case gittools.ScopeStash:
					views = append(views, types.STASH)
				case gittools.ScopeTags:
					views = append(views, types.TAGS)
				case gittools.ScopeRemotes:
					views = append(views, types.REMOTES)
				}
			}
			c.Refresh(types.RefreshOptions{
				Mode:  types.ASYNC,
				Scope: views,
			})
		},
	}
}
