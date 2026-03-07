package gittools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dswcpp/lazygit/pkg/ai/tools"
	"github.com/dswcpp/lazygit/pkg/commands/models"
)

// GetStatusTool returns the current branch, working tree state, and changed files.
type GetStatusTool struct{ d *Deps }

func NewGetStatusTool(d *Deps) tools.Tool { return &GetStatusTool{d} }

func (t *GetStatusTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_status",
		Description: t.d.Tr.ToolGetStatusDesc(),
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermReadOnly,
	}
}

func (t *GetStatusTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Branch: %s\n", t.d.GetCheckedOutBranch()))

	state := t.d.GetWorkingTreeState()
	if state.Any() {
		sb.WriteString(t.d.Tr.ToolStatusInProgress(workingTreeStateDesc(state)) + "\n")
	}

	files := t.d.GetFiles()
	if len(files) == 0 {
		sb.WriteString(t.d.Tr.ToolStatusClean() + "\n")
	} else {
		staged, unstaged, untracked := 0, 0, 0
		for _, f := range files {
			if f.HasStagedChanges {
				staged++
			}
			if f.HasUnstagedChanges {
				unstaged++
			}
			if !f.Tracked {
				untracked++
			}
		}
		sb.WriteString(t.d.Tr.ToolStatusFiles(len(files), staged, unstaged, untracked) + "\n")
		for _, f := range files {
			sb.WriteString(fmt.Sprintf("  %s %s\n", f.ShortStatus, f.Path))
		}
	}
	return tools.ToolResult{CallID: call.ID, Success: true, Output: sb.String()}
}

// GetStagedDiffTool returns the diff of staged changes.
type GetStagedDiffTool struct{ d *Deps }

func NewGetStagedDiffTool(d *Deps) tools.Tool { return &GetStagedDiffTool{d} }

func (t *GetStagedDiffTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_staged_diff",
		Description: t.d.Tr.ToolGetStagedDiffDesc(),
		Params: map[string]tools.ParamSchema{
			"max_lines": {Type: "int", Description: t.d.Tr.ToolMaxLines()},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *GetStagedDiffTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	diff, err := t.d.Diff.GetDiff(true)
	if err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolGetStagedDiffFailed(err)}
	}
	if diff == "" {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolStagedDiffEmpty()}
	}
	maxLines := intParam(call.Params, "max_lines", 300)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: truncateDiff(diff, maxLines, t.d.Tr)}
}

// GetDiffTool returns the diff of unstaged changes.
type GetDiffTool struct{ d *Deps }

func NewGetDiffTool(d *Deps) tools.Tool { return &GetDiffTool{d} }

func (t *GetDiffTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_diff",
		Description: t.d.Tr.ToolGetDiffDesc(),
		Params: map[string]tools.ParamSchema{
			"max_lines": {Type: "int", Description: t.d.Tr.ToolMaxLines()},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *GetDiffTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	diff, err := t.d.Diff.GetDiff(false)
	if err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolGetDiffFailed(err)}
	}
	if diff == "" {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolUnstagedDiffEmpty()}
	}
	maxLines := intParam(call.Params, "max_lines", 300)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: truncateDiff(diff, maxLines, t.d.Tr)}
}

// GetFileDiffTool returns the diff for a specific file.
type GetFileDiffTool struct{ d *Deps }

func NewGetFileDiffTool(d *Deps) tools.Tool { return &GetFileDiffTool{d} }

func (t *GetFileDiffTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_file_diff",
		Description: t.d.Tr.ToolGetFileDiffDesc(),
		Params: map[string]tools.ParamSchema{
			"path":      {Type: "string", Description: t.d.Tr.ToolFilePath(), Required: true},
			"staged":    {Type: "bool", Description: t.d.Tr.ToolGetFileDiffStagedParam()},
			"max_lines": {Type: "int", Description: t.d.Tr.ToolMaxLines()},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *GetFileDiffTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	path := strParam(call.Params, "path", "")
	if path == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingPathParam()}
	}
	staged := boolParam(call.Params, "staged", false)
	maxLines := intParam(call.Params, "max_lines", 300)
	for _, f := range t.d.GetFiles() {
		if f.Path == path {
			diff := t.d.WorkingTree.WorktreeFileDiff(f, true, staged)
			if diff == "" {
				label := t.d.Tr.ToolWorkingDir()
				if staged {
					label = t.d.Tr.ToolStagingArea()
				}
				return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("%s: %s — %s", label, path, t.d.Tr.ToolNoChanges())}
			}
			return tools.ToolResult{CallID: call.ID, Success: true, Output: truncateDiff(diff, maxLines, t.d.Tr)}
		}
	}
	return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolFileNotInWorkdir(path)}
}

// GetLogTool returns recent commits.
type GetLogTool struct{ d *Deps }

func NewGetLogTool(d *Deps) tools.Tool { return &GetLogTool{d} }

func (t *GetLogTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_log",
		Description: t.d.Tr.ToolGetLogDesc(),
		Params: map[string]tools.ParamSchema{
			"count": {Type: "int", Description: t.d.Tr.ToolGetLogCountParam()},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *GetLogTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	count := intParam(call.Params, "count", 15)
	if count <= 0 || count > 50 {
		count = 15
	}
	commits := t.d.GetCommits()
	limit := count
	if len(commits) < limit {
		limit = len(commits)
	}
	var sb strings.Builder
	for i := 0; i < limit; i++ {
		cm := commits[i]
		date := ""
		if cm.UnixTimestamp > 0 {
			date = time.Unix(cm.UnixTimestamp, 0).Format("2006-01-02") + "  "
		}
		sb.WriteString(fmt.Sprintf("%s  %s%s  <%s>\n", cm.ShortHash(), date, cm.Name, cm.AuthorName))
	}
	return tools.ToolResult{CallID: call.ID, Success: true, Output: sb.String()}
}

// GetBranchesTool lists local branches.
type GetBranchesTool struct{ d *Deps }

func NewGetBranchesTool(d *Deps) tools.Tool { return &GetBranchesTool{d} }

func (t *GetBranchesTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_branches",
		Description: t.d.Tr.ToolGetBranchesDesc(),
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermReadOnly,
	}
}

func (t *GetBranchesTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	branches := t.d.GetBranches()
	current := t.d.GetCheckedOutBranch()
	var sb strings.Builder
	for i, b := range branches {
		if i >= 30 {
			sb.WriteString(fmt.Sprintf("... +%d\n", len(branches)-30))
			break
		}
		marker := "  "
		if b.Name == current {
			marker = "* "
		}
		tracking := ""
		switch {
		case b.UpstreamGone:
			tracking = " [upstream: gone]"
		case b.AheadForPush != "" && b.AheadForPush != "0" && b.BehindForPull != "" && b.BehindForPull != "0":
			tracking = fmt.Sprintf(" [↑%s ↓%s]", b.AheadForPush, b.BehindForPull)
		case b.AheadForPush != "" && b.AheadForPush != "0":
			tracking = fmt.Sprintf(" [↑%s]", b.AheadForPush)
		case b.BehindForPull != "" && b.BehindForPull != "0":
			tracking = fmt.Sprintf(" [↓%s]", b.BehindForPull)
		}
		sb.WriteString(fmt.Sprintf("%s%s%s\n", marker, b.Name, tracking))
	}
	return tools.ToolResult{CallID: call.ID, Success: true, Output: sb.String()}
}

// GetStashListTool lists stash entries.
type GetStashListTool struct{ d *Deps }

func NewGetStashListTool(d *Deps) tools.Tool { return &GetStashListTool{d} }

func (t *GetStashListTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_stash_list",
		Description: t.d.Tr.ToolGetStashListDesc(),
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermReadOnly,
	}
}

func (t *GetStashListTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	stashes := t.d.GetStashEntries()
	if len(stashes) == 0 {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolNoStashEntries()}
	}
	var sb strings.Builder
	for _, s := range stashes {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", s.Index, s.Name))
	}
	return tools.ToolResult{CallID: call.ID, Success: true, Output: sb.String()}
}

// GetRemotesTool lists configured remotes.
type GetRemotesTool struct{ d *Deps }

func NewGetRemotesTool(d *Deps) tools.Tool { return &GetRemotesTool{d} }

func (t *GetRemotesTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_remotes",
		Description: t.d.Tr.ToolGetRemotesDesc(),
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermReadOnly,
	}
}

func (t *GetRemotesTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	remotes := t.d.GetRemotes()
	if len(remotes) == 0 {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolNoRemotes()}
	}
	var sb strings.Builder
	for _, r := range remotes {
		sb.WriteString(fmt.Sprintf("%s  (%d branches)\n", r.Name, len(r.Branches)))
	}
	return tools.ToolResult{CallID: call.ID, Success: true, Output: sb.String()}
}

// GetTagsTool lists tags.
type GetTagsTool struct{ d *Deps }

func NewGetTagsTool(d *Deps) tools.Tool { return &GetTagsTool{d} }

func (t *GetTagsTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_tags",
		Description: t.d.Tr.ToolGetTagsDesc(),
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermReadOnly,
	}
}

func (t *GetTagsTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	tags := t.d.GetTags()
	if len(tags) == 0 {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolNoTags()}
	}
	var sb strings.Builder
	for i, tag := range tags {
		if i >= 20 {
			sb.WriteString(fmt.Sprintf("... +%d\n", len(tags)-20))
			break
		}
		sb.WriteString(tag.Name + "\n")
	}
	return tools.ToolResult{CallID: call.ID, Success: true, Output: sb.String()}
}

// GetStashDiffTool returns the diff of a specific stash entry.
type GetStashDiffTool struct{ d *Deps }

func NewGetStashDiffTool(d *Deps) tools.Tool { return &GetStashDiffTool{d} }

func (t *GetStashDiffTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_stash_diff",
		Description: t.d.Tr.ToolGetStashDiffDesc(),
		Params: map[string]tools.ParamSchema{
			"index":     {Type: "int", Description: t.d.Tr.ToolGetStashDiffIndexParam()},
			"max_lines": {Type: "int", Description: t.d.Tr.ToolMaxLines()},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *GetStashDiffTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	idx := intParam(call.Params, "index", 0)
	out, err := t.d.Stash.ShowStashEntryCmdObj(idx).RunWithOutput()
	if err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolGetStashDiffFailed(idx, err)}
	}
	if out == "" {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolStashEntryEmpty(idx)}
	}
	maxLines := intParam(call.Params, "max_lines", 300)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: truncateDiff(out, maxLines, t.d.Tr)}
}

// GetCommitDiffTool returns the diff for a specific commit.
type GetCommitDiffTool struct{ d *Deps }

func NewGetCommitDiffTool(d *Deps) tools.Tool { return &GetCommitDiffTool{d} }

func (t *GetCommitDiffTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_commit_diff",
		Description: t.d.Tr.ToolGetCommitDiffDesc(),
		Params: map[string]tools.ParamSchema{
			"hash":      {Type: "string", Description: t.d.Tr.ToolGetCommitDiffHashParam()},
			"max_lines": {Type: "int", Description: t.d.Tr.ToolMaxLines()},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *GetCommitDiffTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	hash := strParam(call.Params, "hash", "HEAD")
	maxLines := intParam(call.Params, "max_lines", 300)
	diff, err := t.d.Commit.GetCommitDiff(hash)
	if err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolGetCommitDiffFailed(err)}
	}
	return tools.ToolResult{CallID: call.ID, Success: true, Output: truncateDiff(diff, maxLines, t.d.Tr)}
}

// GetBranchDiffTool returns the diff between two branches or refs.
type GetBranchDiffTool struct{ d *Deps }

func NewGetBranchDiffTool(d *Deps) tools.Tool { return &GetBranchDiffTool{d} }

func (t *GetBranchDiffTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_branch_diff",
		Description: t.d.Tr.ToolGetBranchDiffDesc(),
		Params: map[string]tools.ParamSchema{
			"base":      {Type: "string", Description: t.d.Tr.ToolGetBranchDiffBaseParam(), Required: true},
			"target":    {Type: "string", Description: t.d.Tr.ToolGetBranchDiffTargetParam()},
			"max_lines": {Type: "int", Description: t.d.Tr.ToolMaxLines()},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *GetBranchDiffTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	base := strParam(call.Params, "base", "")
	if base == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingParam("base")}
	}
	target := strParam(call.Params, "target", "HEAD")
	maxLines := intParam(call.Params, "max_lines", 300)

	// three-dot syntax: changes on target since it diverged from base
	refRange := base + "..." + target
	diff, err := t.d.Diff.GetDiff(false, refRange, "--")
	if err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolGetBranchDiffFailed(base, target, err)}
	}
	if diff == "" {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolGetBranchDiffEmpty(base, target)}
	}
	return tools.ToolResult{CallID: call.ID, Success: true, Output: truncateDiff(diff, maxLines, t.d.Tr)}
}

// ── helpers ────────────────────────────────────────────────────────────────

func workingTreeStateDesc(state models.WorkingTreeState) string {
	switch {
	case state.Rebasing:
		return "rebase"
	case state.Merging:
		return "merge"
	case state.CherryPicking:
		return "cherry-pick"
	case state.Reverting:
		return "revert"
	default:
		return "unknown"
	}
}

func strParam(params map[string]any, key, defaultVal string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

func intParam(params map[string]any, key string, defaultVal int) int {
	if v, ok := params[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		}
	}
	return defaultVal
}

func boolParam(params map[string]any, key string, defaultVal bool) bool {
	if v, ok := params[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultVal
}

// truncateDiff limits diff output to maxLines lines.
// maxLines <= 0 means no limit.
func truncateDiff(diff string, maxLines int, tr interface {
	ToolTruncated(total, shown int) string
}) string {
	if maxLines <= 0 {
		return diff
	}
	lines := strings.SplitAfter(diff, "\n")
	if len(lines) <= maxLines {
		return diff
	}
	return strings.Join(lines[:maxLines], "") + tr.ToolTruncated(len(lines), maxLines)
}
