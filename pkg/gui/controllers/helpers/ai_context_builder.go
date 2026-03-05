package helpers

import (
	"github.com/dswcpp/lazygit/pkg/ai/repocontext"
)

// GuiContextBuilder implements repocontext.Builder using lazygit's live GUI model.
// It is the single authoritative source replacing the three duplicate buildGitContext()
// implementations that previously existed in ai_helper.go, ai_chat_helper.go, and
// ai_command_helper.go.
type GuiContextBuilder struct {
	c *HelperCommon
}

// NewGuiContextBuilder creates a ContextBuilder backed by the live GUI model.
func NewGuiContextBuilder(c *HelperCommon) repocontext.Builder {
	return &GuiContextBuilder{c: c}
}

// Build constructs a RepoContext snapshot from the current GUI model state.
func (b *GuiContextBuilder) Build() repocontext.RepoContext {
	m := b.c.Model()

	// Files
	files := make([]repocontext.FileStatus, 0, len(m.Files))
	for _, f := range m.Files {
		files = append(files, repocontext.FileStatus{
			Path:        f.Path,
			ShortStatus: f.ShortStatus,
			HasStaged:   f.HasStagedChanges,
			HasUnstaged: f.HasUnstagedChanges,
			IsTracked:   f.Tracked,
		})
	}

	// Recent commits (cap at 10)
	commitLimit := 10
	if len(m.Commits) < commitLimit {
		commitLimit = len(m.Commits)
	}
	commits := make([]repocontext.CommitSummary, 0, commitLimit)
	for _, c := range m.Commits[:commitLimit] {
		commits = append(commits, repocontext.CommitSummary{
			ShortHash: c.ShortHash(),
			Message:   c.Name,
			Author:    c.AuthorName,
		})
	}

	// Upstream info
	upstreamRemote := ""
	aheadForPull := ""
	behindForPull := ""
	if len(m.Branches) > 0 {
		cur := m.Branches[0]
		if cur.IsTrackingRemote() {
			upstreamRemote = cur.UpstreamRemote
			aheadForPull = cur.AheadForPull
			behindForPull = cur.BehindForPull
		}
	}

	// Working tree state
	workingTreeState := ""
	state := m.WorkingTreeStateAtLastCommitRefresh
	if state.Any() {
		switch {
		case state.Rebasing:
			workingTreeState = "rebase"
		case state.Merging:
			workingTreeState = "merge"
		case state.CherryPicking:
			workingTreeState = "cherry-pick"
		case state.Reverting:
			workingTreeState = "revert"
		}
	}

	// Branches (names only, cap at 30)
	branchLimit := 30
	if len(m.Branches) < branchLimit {
		branchLimit = len(m.Branches)
	}
	branchNames := make([]string, 0, branchLimit)
	for _, br := range m.Branches[:branchLimit] {
		branchNames = append(branchNames, br.Name)
	}

	// Tags (names only, cap at 20)
	tagLimit := 20
	if len(m.Tags) < tagLimit {
		tagLimit = len(m.Tags)
	}
	tagNames := make([]string, 0, tagLimit)
	for _, t := range m.Tags[:tagLimit] {
		tagNames = append(tagNames, t.Name)
	}

	return repocontext.RepoContext{
		CurrentBranch:    m.CheckedOutBranch,
		UpstreamRemote:   upstreamRemote,
		AheadForPull:     aheadForPull,
		BehindForPull:    behindForPull,
		WorkingTreeState: workingTreeState,
		Files:            files,
		RecentCommits:    commits,
		StashCount:       len(m.StashEntries),
		Branches:         branchNames,
		Tags:             tagNames,
	}
}
