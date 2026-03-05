package repocontext

import (
	"fmt"
	"strings"
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
func (r RepoContext) CompactString() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("分支: %s\n", r.CurrentBranch))

	if r.WorkingTreeState != "" {
		sb.WriteString(fmt.Sprintf("⚠ 正在进行: %s\n", r.WorkingTreeState))
	}

	if r.UpstreamRemote != "" {
		if r.AheadForPull != "" || r.BehindForPull != "" {
			sb.WriteString(fmt.Sprintf("远程: %s [↑%s ↓%s]\n", r.UpstreamRemote, r.AheadForPull, r.BehindForPull))
		} else {
			sb.WriteString(fmt.Sprintf("远程: %s [已同步]\n", r.UpstreamRemote))
		}
	}

	if len(r.Files) == 0 {
		sb.WriteString("工作区: 干净\n")
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
		sb.WriteString(fmt.Sprintf("变更: %d 个（暂存 %d，未暂存 %d，未追踪 %d）\n",
			len(r.Files), staged, unstaged, untracked))
		limit := len(r.Files)
		if limit > 10 {
			limit = 10
		}
		for i := 0; i < limit; i++ {
			sb.WriteString(fmt.Sprintf("  %s %s\n", r.Files[i].ShortStatus, r.Files[i].Path))
		}
		if len(r.Files) > 10 {
			sb.WriteString(fmt.Sprintf("  ... 还有 %d 个\n", len(r.Files)-10))
		}
	}

	if len(r.RecentCommits) > 0 {
		sb.WriteString("最近提交:\n")
		for _, c := range r.RecentCommits {
			sb.WriteString(fmt.Sprintf("  %s  %s  <%s>\n", c.ShortHash, c.Message, c.Author))
		}
	}

	if r.StashCount > 0 {
		sb.WriteString(fmt.Sprintf("Stash: %d 条\n", r.StashCount))
	}

	return strings.TrimSpace(sb.String())
}
