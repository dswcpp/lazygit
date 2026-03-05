package gittools

import (
	"github.com/dswcpp/lazygit/pkg/commands/git_commands"
	"github.com/dswcpp/lazygit/pkg/commands/models"
	"github.com/dswcpp/lazygit/pkg/utils"
)

// RefreshScope identifies which UI views need updating after a tool executes.
type RefreshScope int

const (
	ScopeFiles    RefreshScope = iota
	ScopeBranches
	ScopeCommits
	ScopeStash
	ScopeTags
	ScopeRemotes
)

// Deps bundles all runtime dependencies that git tools need.
// The GUI layer creates one Deps instance and shares it across all tool instances.
type Deps struct {
	// Git command groups (from pkg/commands/git_commands)
	WorkingTree *git_commands.WorkingTreeCommands
	Branch      *git_commands.BranchCommands
	Commit      *git_commands.CommitCommands
	Stash       *git_commands.StashCommands
	Tag         *git_commands.TagCommands
	Sync        *git_commands.SyncCommands
	Diff        *git_commands.DiffCommands
	Rebase      *git_commands.RebaseCommands

	// Model readers — closures that return the current GUI model state.
	// Using closures lets tools always see the latest state after each operation.
	GetFiles            func() []*models.File
	GetBranches         func() []*models.Branch
	GetCommits          func() []*models.Commit
	GetStashEntries     func() []*models.StashEntry
	GetRemotes          func() []*models.Remote
	GetTags             func() []*models.Tag
	GetCheckedOutBranch func() string
	GetWorkingTreeState func() models.WorkingTreeState
	GetHashPool         func() *utils.StringPool

	// Refresh triggers a UI view refresh after a tool modifies repository state.
	Refresh func(scopes ...RefreshScope)
}
