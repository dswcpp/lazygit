package helpers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dswcpp/lazygit/pkg/commands/git_commands"
	"github.com/dswcpp/lazygit/pkg/commands/models"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

// AIActionType 表示 AI 可执行的操作类型
type AIActionType string

const (
	// ── 查询类（只读）──────────────────────────────────────────
	ActionGetStatus     AIActionType = "get_status"
	ActionGetStagedDiff AIActionType = "get_staged_diff"
	ActionGetDiff       AIActionType = "get_diff"
	ActionGetLog        AIActionType = "get_log"
	ActionGetBranches   AIActionType = "get_branches"
	ActionGetStashList  AIActionType = "get_stash_list"
	ActionGetRemotes    AIActionType = "get_remotes"
	ActionGetTags       AIActionType = "get_tags"
	ActionGetCommitDiff AIActionType = "get_commit_diff"

	// ── 文件暂存 ──────────────────────────────────────────────
	ActionStageAll    AIActionType = "stage_all"
	ActionStageFile   AIActionType = "stage_file"
	ActionUnstageAll  AIActionType = "unstage_all"
	ActionUnstageFile AIActionType = "unstage_file"
	ActionDiscardFile AIActionType = "discard_file"

	// ── 提交 ──────────────────────────────────────────────────
	ActionCommit       AIActionType = "commit"
	ActionAmendHead    AIActionType = "amend_head"
	ActionRevertCommit AIActionType = "revert_commit"
	ActionResetSoft    AIActionType = "reset_soft"
	ActionResetMixed   AIActionType = "reset_mixed"
	ActionCherryPick   AIActionType = "cherry_pick"

	// ── 分支 ──────────────────────────────────────────────────
	ActionCheckout     AIActionType = "checkout"
	ActionNewBranch    AIActionType = "create_branch"
	ActionDeleteBranch AIActionType = "delete_branch"
	ActionRenameBranch AIActionType = "rename_branch"
	ActionMergeBranch  AIActionType = "merge_branch"

	// ── Stash ─────────────────────────────────────────────────
	ActionStash      AIActionType = "stash"
	ActionStashPop   AIActionType = "stash_pop"
	ActionStashApply AIActionType = "stash_apply"
	ActionStashDrop  AIActionType = "stash_drop"

	// ── Tag ───────────────────────────────────────────────────
	ActionCreateTag AIActionType = "create_tag"
	ActionDeleteTag AIActionType = "delete_tag"

	// ── 远程同步 ──────────────────────────────────────────────
	ActionFetch AIActionType = "fetch"
	ActionPush  AIActionType = "push"
)

// AIAction 是 AI 在回复中嵌入的操作指令
type AIAction struct {
	Type    AIActionType `json:"type"`
	// 通用参数
	Path    string `json:"path,omitempty"`
	Message string `json:"message,omitempty"`
	Name    string `json:"name,omitempty"`
	// 分支/提交操作
	Base   string `json:"base,omitempty"`
	Old    string `json:"old,omitempty"`
	Hash   string `json:"hash,omitempty"`
	Steps  int    `json:"steps,omitempty"`
	// 列表类
	Count  int    `json:"count,omitempty"`
	Index  int    `json:"index,omitempty"`
	Ref    string `json:"ref,omitempty"`
	// Push
	Force  bool   `json:"force,omitempty"`
}

// AIActionResult 是操作执行结果
type AIActionResult struct {
	Type    AIActionType
	Success bool
	Output  string
}

// parseActionsFromResponse 从 AI 回复文本中提取 <action>...</action> 块
func parseActionsFromResponse(text string) []AIAction {
	var actions []AIAction
	remaining := text
	for {
		start := strings.Index(remaining, "<action>")
		end := strings.Index(remaining, "</action>")
		if start == -1 || end == -1 || end <= start {
			break
		}
		jsonStr := strings.TrimSpace(remaining[start+8 : end])
		var action AIAction
		if err := json.Unmarshal([]byte(jsonStr), &action); err == nil {
			actions = append(actions, action)
		}
		remaining = remaining[end+9:]
	}
	return actions
}

// stripActionBlocks 移除文本中的 <action>...</action> 块，返回纯显示文本
func stripActionBlocks(text string) string {
	var sb strings.Builder
	remaining := text
	for {
		start := strings.Index(remaining, "<action>")
		if start == -1 {
			sb.WriteString(remaining)
			break
		}
		end := strings.Index(remaining, "</action>")
		if end == -1 {
			sb.WriteString(remaining)
			break
		}
		before := strings.TrimRight(remaining[:start], " \t\n")
		if before != "" {
			sb.WriteString(before)
			sb.WriteString("\n")
		}
		remaining = strings.TrimLeft(remaining[end+9:], "\n")
	}
	return strings.TrimSpace(sb.String())
}

// executeAction 执行单个 AI 操作（可从 goroutine 安全调用，git 命令均为进程调用）
func executeAction(c *HelperCommon, action AIAction) AIActionResult {
	res := AIActionResult{Type: action.Type}

	switch action.Type {

	// ─── 查询类 ────────────────────────────────────────────────────────────

	case ActionGetStatus:
		branch := c.Model().CheckedOutBranch
		files := c.Model().Files
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("当前分支: %s\n", branch))

		// 是否在 rebase/merge 中
		state := c.Model().WorkingTreeStateAtLastCommitRefresh
		if state.Any() {
			stateDesc := ""
			if state.Rebasing {
				stateDesc = "rebase"
			} else if state.Merging {
				stateDesc = "merge"
			} else if state.CherryPicking {
				stateDesc = "cherry-pick"
			} else if state.Reverting {
				stateDesc = "revert"
			}
			sb.WriteString(fmt.Sprintf("⚠ 正在进行: %s\n", stateDesc))
		}

		if len(files) == 0 {
			sb.WriteString("工作区干净\n")
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
		res.Success = true
		res.Output = sb.String()

	case ActionGetStagedDiff:
		diff, err := c.Git().Diff.GetDiff(true)
		if err != nil {
			res.Output = fmt.Sprintf("获取暂存区 diff 失败: %v", err)
			return res
		}
		if diff == "" {
			res.Output = "暂存区为空"
		} else {
			res.Output = diff
		}
		res.Success = true

	case ActionGetDiff:
		diff, err := c.Git().Diff.GetDiff(false)
		if err != nil {
			res.Output = fmt.Sprintf("获取 diff 失败: %v", err)
			return res
		}
		if diff == "" {
			res.Output = "没有未暂存的变更"
		} else {
			res.Output = diff
		}
		res.Success = true

	case ActionGetLog:
		count := action.Count
		if count <= 0 || count > 50 {
			count = 15
		}
		commits := c.Model().Commits
		var sb strings.Builder
		limit := count
		if len(commits) < limit {
			limit = len(commits)
		}
		for i := 0; i < limit; i++ {
			cm := commits[i]
			sb.WriteString(fmt.Sprintf("%s  %s  <%s>\n", cm.ShortHash(), cm.Name, cm.AuthorName))
		}
		res.Success = true
		res.Output = sb.String()

	case ActionGetBranches:
		branches := c.Model().Branches
		current := c.Model().CheckedOutBranch
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
			if b.AheadForPull != "" && b.AheadForPull != "0" {
				tracking = fmt.Sprintf(" [↑%s ↓%s]", b.AheadForPull, b.BehindForPull)
			}
			sb.WriteString(fmt.Sprintf("%s%s%s\n", marker, b.Name, tracking))
		}
		res.Success = true
		res.Output = sb.String()

	case ActionGetStashList:
		stashes := c.Model().StashEntries
		if len(stashes) == 0 {
			res.Output = "没有储藏的变更"
		} else {
			var sb strings.Builder
			for _, s := range stashes {
				sb.WriteString(fmt.Sprintf("[%d] %s\n", s.Index, s.Name))
			}
			res.Output = sb.String()
		}
		res.Success = true

	case ActionGetRemotes:
		remotes := c.Model().Remotes
		if len(remotes) == 0 {
			res.Output = "没有配置远程仓库"
		} else {
			var sb strings.Builder
			for _, r := range remotes {
				sb.WriteString(fmt.Sprintf("%s  (%d 个分支)\n", r.Name, len(r.Branches)))
			}
			res.Output = sb.String()
		}
		res.Success = true

	case ActionGetTags:
		tags := c.Model().Tags
		if len(tags) == 0 {
			res.Output = "没有 tag"
		} else {
			var sb strings.Builder
			for i, t := range tags {
				if i >= 20 {
					sb.WriteString(fmt.Sprintf("... 还有 %d 个\n", len(tags)-20))
					break
				}
				sb.WriteString(fmt.Sprintf("%s\n", t.Name))
			}
			res.Output = sb.String()
		}
		res.Success = true

	case ActionGetCommitDiff:
		hash := action.Hash
		if hash == "" {
			hash = "HEAD"
		}
		diff, err := c.Git().Commit.GetCommitDiff(hash)
		if err != nil {
			res.Output = fmt.Sprintf("获取提交 diff 失败: %v", err)
			return res
		}
		res.Success = true
		res.Output = diff

	// ─── 文件暂存 ──────────────────────────────────────────────────────────

	case ActionStageAll:
		if err := c.Git().WorkingTree.StageAll(false); err != nil {
			res.Output = fmt.Sprintf("暂存全部失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
		res.Success = true
		res.Output = "已暂存所有变更"

	case ActionStageFile:
		if action.Path == "" {
			res.Output = "缺少 path 参数"
			return res
		}
		if err := c.Git().WorkingTree.StageFile(action.Path); err != nil {
			res.Output = fmt.Sprintf("暂存文件失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
		res.Success = true
		res.Output = fmt.Sprintf("已暂存: %s", action.Path)

	case ActionUnstageAll:
		if err := c.Git().WorkingTree.ResetMixed("HEAD"); err != nil {
			res.Output = fmt.Sprintf("取消所有暂存失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
		res.Success = true
		res.Output = "已取消所有暂存"

	case ActionUnstageFile:
		if action.Path == "" {
			res.Output = "缺少 path 参数"
			return res
		}
		tracked := true
		for _, f := range c.Model().Files {
			if f.Path == action.Path {
				tracked = f.Tracked
				break
			}
		}
		if err := c.Git().WorkingTree.UnStageFile([]string{action.Path}, tracked); err != nil {
			res.Output = fmt.Sprintf("取消暂存失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
		res.Success = true
		res.Output = fmt.Sprintf("已取消暂存: %s", action.Path)

	case ActionDiscardFile:
		if action.Path == "" {
			res.Output = "缺少 path 参数"
			return res
		}
		for _, f := range c.Model().Files {
			if f.Path == action.Path {
				if err := c.Git().WorkingTree.DiscardAllFileChanges(f); err != nil {
					res.Output = fmt.Sprintf("丢弃变更失败: %v", err)
					return res
				}
				break
			}
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
		res.Success = true
		res.Output = fmt.Sprintf("已丢弃 %s 的所有变更", action.Path)

	// ─── 提交操作 ──────────────────────────────────────────────────────────

	case ActionCommit:
		if action.Message == "" {
			res.Output = "缺少 message 参数"
			return res
		}
		if err := c.Git().Commit.CommitCmdObj(action.Message, "", false).Run(); err != nil {
			res.Output = fmt.Sprintf("提交失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES, types.COMMITS}})
		res.Success = true
		res.Output = fmt.Sprintf("提交成功: \"%s\"", action.Message)

	case ActionAmendHead:
		if action.Message == "" {
			res.Output = "缺少 message 参数（新的提交信息）"
			return res
		}
		if err := c.Git().Commit.RewordLastCommit(action.Message, "").Run(); err != nil {
			res.Output = fmt.Sprintf("修改提交信息失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS}})
		res.Success = true
		res.Output = fmt.Sprintf("已修改最新提交信息为: \"%s\"", action.Message)

	case ActionRevertCommit:
		if action.Hash == "" {
			res.Output = "缺少 hash 参数"
			return res
		}
		if err := c.Git().Commit.Revert([]string{action.Hash}, false); err != nil {
			res.Output = fmt.Sprintf("revert 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS, types.FILES}})
		res.Success = true
		res.Output = fmt.Sprintf("已 revert 提交: %s", action.Hash)

	case ActionResetSoft:
		steps := action.Steps
		if steps <= 0 {
			steps = 1
		}
		ref := fmt.Sprintf("HEAD~%d", steps)
		if action.Hash != "" {
			ref = action.Hash
		}
		if err := c.Git().WorkingTree.ResetSoft(ref); err != nil {
			res.Output = fmt.Sprintf("reset --soft 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES, types.COMMITS}})
		res.Success = true
		res.Output = fmt.Sprintf("reset --soft 到 %s，变更已保留在暂存区", ref)

	case ActionResetMixed:
		steps := action.Steps
		if steps <= 0 {
			steps = 1
		}
		ref := fmt.Sprintf("HEAD~%d", steps)
		if action.Hash != "" {
			ref = action.Hash
		}
		if err := c.Git().WorkingTree.ResetMixed(ref); err != nil {
			res.Output = fmt.Sprintf("reset --mixed 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES, types.COMMITS}})
		res.Success = true
		res.Output = fmt.Sprintf("reset --mixed 到 %s，变更已保留在工作区（未暂存）", ref)

	case ActionCherryPick:
		if action.Hash == "" {
			res.Output = "缺少 hash 参数"
			return res
		}
		// 构造临时 commit 对象用于 cherry-pick
		commit := models.NewCommit(c.Model().HashPool, models.NewCommitOpts{Hash: action.Hash})
		if err := c.Git().Rebase.CherryPickCommits([]*models.Commit{commit}); err != nil {
			res.Output = fmt.Sprintf("cherry-pick 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS, types.FILES}})
		res.Success = true
		res.Output = fmt.Sprintf("已 cherry-pick: %s", action.Hash)

	// ─── 分支操作 ──────────────────────────────────────────────────────────

	case ActionCheckout:
		if action.Name == "" {
			res.Output = "缺少 name 参数（分支名）"
			return res
		}
		if err := c.Git().Branch.Checkout(action.Name, git_commands.CheckoutOptions{}); err != nil {
			res.Output = fmt.Sprintf("切换分支失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES, types.BRANCHES, types.COMMITS}})
		res.Success = true
		res.Output = fmt.Sprintf("已切换到: %s", action.Name)

	case ActionNewBranch:
		if action.Name == "" {
			res.Output = "缺少 name 参数"
			return res
		}
		base := action.Base
		if base == "" {
			base = "HEAD"
		}
		if err := c.Git().Branch.New(action.Name, base); err != nil {
			res.Output = fmt.Sprintf("创建分支失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES}})
		res.Success = true
		res.Output = fmt.Sprintf("已创建分支 %s（基于 %s）", action.Name, base)

	case ActionDeleteBranch:
		if action.Name == "" {
			res.Output = "缺少 name 参数"
			return res
		}
		if err := c.Git().Branch.LocalDelete([]string{action.Name}, false); err != nil {
			res.Output = fmt.Sprintf("删除分支失败: %v（如需强制删除请确保分支已合并）", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES}})
		res.Success = true
		res.Output = fmt.Sprintf("已删除分支: %s", action.Name)

	case ActionRenameBranch:
		if action.Old == "" || action.Name == "" {
			res.Output = "缺少 old（旧名）或 name（新名）参数"
			return res
		}
		if err := c.Git().Branch.Rename(action.Old, action.Name); err != nil {
			res.Output = fmt.Sprintf("重命名失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES}})
		res.Success = true
		res.Output = fmt.Sprintf("已将分支 %s 重命名为 %s", action.Old, action.Name)

	case ActionMergeBranch:
		if action.Name == "" {
			res.Output = "缺少 name 参数（要合并的分支名）"
			return res
		}
		if err := c.Git().Branch.Merge(action.Name, git_commands.MERGE_VARIANT_REGULAR); err != nil {
			res.Output = fmt.Sprintf("merge 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS, types.BRANCHES, types.FILES}})
		res.Success = true
		res.Output = fmt.Sprintf("已将 %s 合并到当前分支", action.Name)

	// ─── Stash ─────────────────────────────────────────────────────────────

	case ActionStash:
		msg := action.Message
		if msg == "" {
			msg = "AI stash"
		}
		if err := c.Git().Stash.Push(msg); err != nil {
			res.Output = fmt.Sprintf("stash 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES, types.STASH}})
		res.Success = true
		res.Output = fmt.Sprintf("已储藏变更: %s", msg)

	case ActionStashPop:
		if err := c.Git().Stash.Pop(0); err != nil {
			res.Output = fmt.Sprintf("stash pop 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES, types.STASH}})
		res.Success = true
		res.Output = "已恢复最近的 stash"

	case ActionStashApply:
		idx := action.Index
		if err := c.Git().Stash.Apply(idx); err != nil {
			res.Output = fmt.Sprintf("stash apply 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
		res.Success = true
		res.Output = fmt.Sprintf("已应用 stash[%d]（stash 条目保留）", idx)

	case ActionStashDrop:
		idx := action.Index
		if err := c.Git().Stash.Drop(idx); err != nil {
			res.Output = fmt.Sprintf("stash drop 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH}})
		res.Success = true
		res.Output = fmt.Sprintf("已删除 stash[%d]", idx)

	// ─── Tag ───────────────────────────────────────────────────────────────

	case ActionCreateTag:
		if action.Name == "" {
			res.Output = "缺少 name 参数（tag 名）"
			return res
		}
		ref := action.Ref
		if ref == "" {
			ref = "HEAD"
		}
		if err := c.Git().Tag.CreateLightweightObj(action.Name, ref, false).Run(); err != nil {
			res.Output = fmt.Sprintf("创建 tag 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.TAGS}})
		res.Success = true
		res.Output = fmt.Sprintf("已创建 tag %s（指向 %s）", action.Name, ref)

	case ActionDeleteTag:
		if action.Name == "" {
			res.Output = "缺少 name 参数（tag 名）"
			return res
		}
		if err := c.Git().Tag.LocalDelete(action.Name); err != nil {
			res.Output = fmt.Sprintf("删除 tag 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.TAGS}})
		res.Success = true
		res.Output = fmt.Sprintf("已删除本地 tag: %s", action.Name)

	// ─── 远程同步 ──────────────────────────────────────────────────────────

	case ActionFetch:
		if err := c.Git().Sync.FetchBackground(); err != nil {
			res.Output = fmt.Sprintf("fetch 失败: %v", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS}})
		res.Success = true
		res.Output = "fetch 完成"

	case ActionPush:
		cmdObj, err := c.Git().Sync.PushCmdObj(nil, git_commands.PushOpts{
			Force: action.Force,
		})
		if err != nil {
			res.Output = fmt.Sprintf("push 配置错误: %v", err)
			return res
		}
		if err := cmdObj.Run(); err != nil {
			res.Output = fmt.Sprintf("push 失败: %v（请确认远程配置和认证）", err)
			return res
		}
		c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS}})
		res.Success = true
		res.Output = "push 成功"

	default:
		res.Output = fmt.Sprintf("未知操作: %s", action.Type)
	}

	return res
}

// aiAgentSystemPrompt 返回完整的 AI Agent 系统提示
func aiAgentSystemPrompt() string {
	return `你是 lazygit 的内置 AI Agent，可以直接操控 Git 仓库。
用户对你说话就像对同事下指令：你要主动思考、自主执行，而不是给出建议让用户手动操作。

━━━ 工具调用格式 ━━━
在回复中嵌入 <action>JSON</action>，系统会自动执行并将结果反馈给你。
每次回复可以调用多个工具。工具执行是有序的，前一个结果可以影响后续决策。

━━━ 完整工具列表 ━━━

【查询工具（优先使用，了解当前状态）】
  <action>{"type":"get_status"}</action>
    → 工作区状态：变更文件、暂存情况、是否在 rebase/merge 中

  <action>{"type":"get_staged_diff"}</action>
    → 暂存区的完整 diff

  <action>{"type":"get_diff"}</action>
    → 未暂存变更的 diff

  <action>{"type":"get_log","count":15}</action>
    → 最近提交历史（hash、信息、作者）

  <action>{"type":"get_branches"}</action>
    → 所有本地分支及 ahead/behind 状态

  <action>{"type":"get_stash_list"}</action>
    → stash 列表

  <action>{"type":"get_remotes"}</action>
    → 远程仓库配置

  <action>{"type":"get_tags"}</action>
    → 所有本地 tag

  <action>{"type":"get_commit_diff","hash":"abc1234"}</action>
    → 指定提交的 diff（hash 为空则查 HEAD）

【暂存操作】
  <action>{"type":"stage_all"}</action>
  <action>{"type":"stage_file","path":"src/main.go"}</action>
  <action>{"type":"unstage_all"}</action>
  <action>{"type":"unstage_file","path":"src/main.go"}</action>
  <action>{"type":"discard_file","path":"src/main.go"}</action>   ⚠ 不可恢复

【提交操作】
  <action>{"type":"commit","message":"feat: 描述"}</action>
  <action>{"type":"amend_head","message":"fix: 修正提交信息"}</action>
  <action>{"type":"revert_commit","hash":"abc1234"}</action>
  <action>{"type":"reset_soft","steps":1}</action>                 撤销提交，变更回到暂存区
  <action>{"type":"reset_mixed","steps":1}</action>               撤销提交，变更回到工作区
  <action>{"type":"cherry_pick","hash":"abc1234"}</action>

【分支操作】
  <action>{"type":"checkout","name":"develop"}</action>
  <action>{"type":"create_branch","name":"feat/x","base":"main"}</action>
  <action>{"type":"delete_branch","name":"old-branch"}</action>
  <action>{"type":"rename_branch","old":"old-name","name":"new-name"}</action>
  <action>{"type":"merge_branch","name":"feature/x"}</action>

【Stash 操作】
  <action>{"type":"stash","message":"wip: 临时保存"}</action>
  <action>{"type":"stash_pop"}</action>
  <action>{"type":"stash_apply","index":0}</action>               应用但不删除
  <action>{"type":"stash_drop","index":0}</action>

【Tag 操作】
  <action>{"type":"create_tag","name":"v1.0.0","ref":"HEAD"}</action>
  <action>{"type":"delete_tag","name":"v1.0.0"}</action>

【远程同步】
  <action>{"type":"fetch"}</action>
  <action>{"type":"push"}</action>
  <action>{"type":"push","force":true}</action>                    ⚠ 强制推送

━━━ 行为准则 ━━━

1. 【先查询再操作】
   不确定当前状态时，先调用 get_status / get_log 等查询工具。
   例：用户说"提交代码" → 先 get_status 确认有哪些变更 → 决定是 stage_all 还是按文件暂存 → 再 commit。

2. 【自主完成多步任务】
   用户说"帮我把当前修改 stash，切到 main 分支" → 你要：
   stash → checkout main → 告诉用户已完成，无需用户手动操作任何步骤。

3. 【commit message 规范】
   遵循 Conventional Commits：feat/fix/refactor/docs/test/chore/perf/ci: 描述
   根据 staged diff 内容自动判断类型，不要让用户手动起名。

4. 【危险操作要说明】
   discard_file / reset_hard / push --force / delete_branch 执行前先说明原因和影响。
   如果是用户明确要求，直接执行；如果是你主动判断需要，先告知再执行。

5. 【不要说"请你执行..."】
   你有完整的操作能力，任何你描述的步骤，都应该由你通过工具执行，不要让用户手动运行命令。

6. 【操作完成后总结】
   执行完所有操作后，用简洁的语言告知用户完成了什么，当前状态如何。

━━━ 典型工作流示例 ━━━

示例1：「帮我提交代码」
  → get_status（查看哪些文件有变更）
  → get_staged_diff（查看暂存内容，生成提交信息）
  → 如果有文件未暂存：stage_all
  → commit（自动生成规范 commit message）
  → 告知提交结果

示例2：「我要修一个 bug，要创建新分支」
  → get_status（确认工作区干净）
  → 如果有未提交变更：stash
  → get_branches（了解当前分支结构）
  → create_branch（从合适的 base 创建）
  → checkout（切换到新分支）
  → 告知用户已就绪

示例3：「刚才的提交信息写错了」
  → get_log（确认最新提交）
  → amend_head（修改提交信息）
  → 告知完成`
}
