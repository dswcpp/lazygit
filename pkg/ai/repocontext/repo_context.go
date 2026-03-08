package repocontext

import (
	"fmt"
	"strings"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
)

// FileStatus describes a single file in the working tree.
type FileStatus struct {
	Path        string
	ShortStatus string
	HasStaged   bool
	HasUnstaged bool
	IsTracked   bool
}

// CommitSummary is a compact representation of a commit for prompt injection.
type CommitSummary struct {
	ShortHash string
	Message   string
	Author    string
}

// RepoContext is a snapshot of the current repository state.
// It is the single authoritative source of context passed to AI components.
type RepoContext struct {
	CurrentBranch    string
	UpstreamRemote   string
	AheadForPull     string
	BehindForPull    string
	// WorkingTreeState is "" when clean, or "rebase" | "merge" | "cherry-pick" | "revert"
	WorkingTreeState string
	Files            []FileStatus
	RecentCommits    []CommitSummary
	StashCount       int
	Branches         []string
	Tags             []string
}

// CompactString generates a concise text representation for inclusion in prompts.
func (r RepoContext) CompactString(tr *aii18n.Translator) string {
	var sb strings.Builder

	sb.WriteString(tr.RepoBranch(r.CurrentBranch))

	if r.WorkingTreeState != "" {
		sb.WriteString(tr.RepoInProgress(r.WorkingTreeState))
	}

	if r.UpstreamRemote != "" {
		if r.AheadForPull != "" || r.BehindForPull != "" {
			sb.WriteString(tr.RepoRemoteAheadBehind(r.UpstreamRemote, r.AheadForPull, r.BehindForPull))
		} else {
			sb.WriteString(tr.RepoRemoteSynced(r.UpstreamRemote))
		}
	}

	if len(r.Files) == 0 {
		sb.WriteString(tr.RepoWorkingDirClean())
	} else {
		staged, unstaged, untracked := 0, 0, 0
		for _, f := range r.Files {
			if f.HasStaged {
				staged++
			}
			if f.HasUnstaged {
				unstaged++
			}
			if !f.IsTracked {
				untracked++
			}
		}
		sb.WriteString(tr.RepoChanges(len(r.Files), staged, unstaged, untracked))
		limit := len(r.Files)
		if limit > 10 {
			limit = 10
		}
		for i := 0; i < limit; i++ {
			sb.WriteString(fmt.Sprintf("  %s %s\n", r.Files[i].ShortStatus, r.Files[i].Path))
		}
		if len(r.Files) > 10 {
			sb.WriteString(tr.MoreItems(len(r.Files) - 10))
		}
	}

	if len(r.RecentCommits) > 0 {
		sb.WriteString(tr.RepoRecentCommits())
		for _, c := range r.RecentCommits {
			sb.WriteString(fmt.Sprintf("  %s  %s  <%s>\n", c.ShortHash, c.Message, c.Author))
		}
	}

	if r.StashCount > 0 {
		sb.WriteString(tr.RepoStashCount(r.StashCount))
	}

	return strings.TrimSpace(sb.String())
}
