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
		Description: "获取当前仓库状态（分支、工作区文件、rebase/merge 进度）",
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermReadOnly,
	}
}

func (t *GetStatusTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("当前分支: %s\n", t.d.GetCheckedOutBranch()))

	state := t.d.GetWorkingTreeState()
	if state.Any() {
		desc := workingTreeStateDesc(state)
		sb.WriteString(fmt.Sprintf("⚠ 正在进行: %s\n", desc))
	}

	files := t.d.GetFiles()
	if len(files) == 0 {
		sb.WriteString("工作区: 干净\n")
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
		sb.WriteString(fmt.Sprintf("变更文件: %d 个（已暂存 %d，未暂存 %d，未追踪 %d）\n",
			len(files), staged, unstaged, untracked))
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
		Description: "获取暂存区（staged）diff",
		Params: map[string]tools.ParamSchema{
			"max_lines": {Type: "int", Description: "最多返回行数（默认 300，0 表示不限制）"},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *GetStagedDiffTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	diff, err := t.d.Diff.GetDiff(true)
	if err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("获取暂存区 diff 失败: %v", err)}
	}
	if diff == "" {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: "暂存区为空"}
	}
	maxLines := intParam(call.Params, "max_lines", 300)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: truncateDiff(diff, maxLines)}
}

// GetDiffTool returns the diff of unstaged changes.
type GetDiffTool struct{ d *Deps }

func NewGetDiffTool(d *Deps) tools.Tool { return &GetDiffTool{d} }

func (t *GetDiffTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_diff",
		Description: "获取工作区未暂存的 diff",
		Params: map[string]tools.ParamSchema{
			"max_lines": {Type: "int", Description: "最多返回行数（默认 300，0 表示不限制）"},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *GetDiffTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	diff, err := t.d.Diff.GetDiff(false)
	if err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("获取 diff 失败: %v", err)}
	}
	if diff == "" {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: "没有未暂存的变更"}
	}
	maxLines := intParam(call.Params, "max_lines", 300)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: truncateDiff(diff, maxLines)}
}

// GetFileDiffTool returns the diff for a specific file.
type GetFileDiffTool struct{ d *Deps }

func NewGetFileDiffTool(d *Deps) tools.Tool { return &GetFileDiffTool{d} }

func (t *GetFileDiffTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_file_diff",
		Description: "获取指定文件的 diff（可选择 staged 或 unstaged）",
		Params: map[string]tools.ParamSchema{
			"path":   {Type: "string", Description: "文件路径", Required: true},
			"staged": {Type: "bool", Description: "true=暂存区 diff，false=工作区 diff（默认 false）"},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *GetFileDiffTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	path := strParam(call.Params, "path", "")
	if path == "" {
		return tools.ToolResult{CallID: call.ID, Output: "缺少 path 参数"}
	}
	staged := boolParam(call.Params, "staged", false)
	for _, f := range t.d.GetFiles() {
		if f.Path == path {
			diff := t.d.WorkingTree.WorktreeFileDiff(f, true, staged)
			if diff == "" {
				label := "未暂存"
				if staged {
					label = "暂存区"
				}
				return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("%s 没有变更: %s", label, path)}
			}
			return tools.ToolResult{CallID: call.ID, Success: true, Output: diff}
		}
	}
	return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("文件不在工作区变更列表中: %s", path)}
}

// GetLogTool returns recent commits.
type GetLogTool struct{ d *Deps }

func NewGetLogTool(d *Deps) tools.Tool { return &GetLogTool{d} }

func (t *GetLogTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_log",
		Description: "获取最近的提交记录",
		Params: map[string]tools.ParamSchema{
			"count": {Type: "int", Description: "返回条数，默认 15，最多 50"},
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
		Description: "列出本地分支（当前分支标 *）",
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
			sb.WriteString(fmt.Sprintf("... 还有 %d 个\n", len(branches)-30))
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
			// has both: commits to push and commits to pull
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
		Description: "列出所有 stash 条目",
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermReadOnly,
	}
}

func (t *GetStashListTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	stashes := t.d.GetStashEntries()
	if len(stashes) == 0 {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: "没有储藏的变更"}
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
		Description: "列出所有远程仓库",
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermReadOnly,
	}
}

func (t *GetRemotesTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	remotes := t.d.GetRemotes()
	if len(remotes) == 0 {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: "没有配置远程仓库"}
	}
	var sb strings.Builder
	for _, r := range remotes {
		sb.WriteString(fmt.Sprintf("%s  (%d 个分支)\n", r.Name, len(r.Branches)))
	}
	return tools.ToolResult{CallID: call.ID, Success: true, Output: sb.String()}
}

// GetTagsTool lists tags.
type GetTagsTool struct{ d *Deps }

func NewGetTagsTool(d *Deps) tools.Tool { return &GetTagsTool{d} }

func (t *GetTagsTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_tags",
		Description: "列出所有 tag",
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermReadOnly,
	}
}

func (t *GetTagsTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	tags := t.d.GetTags()
	if len(tags) == 0 {
		return tools.ToolResult{CallID: call.ID, Success: true, Output: "没有 tag"}
	}
	var sb strings.Builder
	for i, tag := range tags {
		if i >= 20 {
			sb.WriteString(fmt.Sprintf("... 还有 %d 个\n", len(tags)-20))
			break
		}
		sb.WriteString(tag.Name + "\n")
	}
	return tools.ToolResult{CallID: call.ID, Success: true, Output: sb.String()}
}

// GetCommitDiffTool returns the diff for a specific commit.
type GetCommitDiffTool struct{ d *Deps }

func NewGetCommitDiffTool(d *Deps) tools.Tool { return &GetCommitDiffTool{d} }

func (t *GetCommitDiffTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "get_commit_diff",
		Description: "获取指定提交的 diff（默认 HEAD）",
		Params: map[string]tools.ParamSchema{
			"hash": {Type: "string", Description: "提交 hash，留空表示 HEAD"},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *GetCommitDiffTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	hash := strParam(call.Params, "hash", "HEAD")
	diff, err := t.d.Commit.GetCommitDiff(hash)
	if err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("获取提交 diff 失败: %v", err)}
	}
	return tools.ToolResult{CallID: call.ID, Success: true, Output: diff}
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
func truncateDiff(diff string, maxLines int) string {
	if maxLines <= 0 {
		return diff
	}
	lines := strings.SplitAfter(diff, "\n")
	if len(lines) <= maxLines {
		return diff
	}
	truncated := strings.Join(lines[:maxLines], "")
	return truncated + fmt.Sprintf("\n... (已截断，共 %d 行，仅显示前 %d 行)", len(lines), maxLines)
}
