package helpers

import (
	"strings"

	"github.com/dswcpp/lazygit/pkg/commands/models"
)

// Completion represents a single completion suggestion.
type Completion struct {
	Text        string // The completion text
	Description string // Description of what this completion does
	Type        string // Type: "command", "flag", "branch", "file", "tag", "remote"
	Priority    int    // Priority for sorting (higher = more important)
}

// CompletionEngine provides intelligent command completion based on context.
type CompletionEngine struct {
	c *HelperCommon

	// Cached data
	gitSubcommands []string
	gitFlags       map[string][]string
}

// NewCompletionEngine creates a new completion engine.
func NewCompletionEngine(c *HelperCommon) *CompletionEngine {
	return &CompletionEngine{
		c:              c,
		gitSubcommands: GetGitSubcommands(),
		gitFlags:       GetGitFlagCompletions(),
	}
}

// Complete provides completions for the given input.
func (e *CompletionEngine) Complete(input string) []Completion {
	parts := parseCommandInput(input)

	if len(parts) == 0 {
		return nil
	}

	// Case 1: Completing main command (e.g., "gi" -> "git")
	if len(parts) == 1 {
		return e.completeMainCommand(parts[0])
	}

	// Case 2: Completing git subcommand (e.g., "git com" -> "git commit")
	if parts[0] == "git" && len(parts) == 2 {
		return e.completeGitSubcommand(parts[1])
	}

	// Case 3: Completing arguments (flags, branches, files, etc.)
	if parts[0] == "git" && len(parts) >= 2 {
		return e.completeArguments(parts)
	}

	return nil
}

// completeMainCommand completes the main command (e.g., "git").
func (e *CompletionEngine) completeMainCommand(partial string) []Completion {
	commands := []string{"git", "cd", "ls", "pwd", "cat", "grep", "find"}

	var completions []Completion
	for _, cmd := range commands {
		if strings.HasPrefix(cmd, partial) {
			completions = append(completions, Completion{
				Text:        cmd,
				Description: getCommandDescription(cmd),
				Type:        "command",
				Priority:    10,
			})
		}
	}

	return completions
}

// completeGitSubcommand completes git subcommands.
func (e *CompletionEngine) completeGitSubcommand(partial string) []Completion {
	var completions []Completion

	for _, subcmd := range e.gitSubcommands {
		if strings.HasPrefix(subcmd, partial) {
			completions = append(completions, Completion{
				Text:        subcmd,
				Description: getGitSubcommandDescription(subcmd),
				Type:        "command",
				Priority:    getPriorityForSubcommand(subcmd),
			})
		}
	}

	return completions
}

// completeArguments completes arguments based on the git subcommand.
func (e *CompletionEngine) completeArguments(parts []string) []Completion {
	if len(parts) < 2 {
		return nil
	}

	subcommand := parts[1]
	lastPart := parts[len(parts)-1]

	var completions []Completion

	// Check if completing a flag
	if strings.HasPrefix(lastPart, "-") {
		completions = append(completions, e.completeFlags(subcommand, lastPart)...)
	}

	// Context-aware completions based on subcommand
	switch subcommand {
	case "checkout", "switch", "merge", "rebase", "branch":
		completions = append(completions, e.completeBranches(lastPart)...)

	case "push", "pull", "fetch":
		if !strings.HasPrefix(lastPart, "-") {
			completions = append(completions, e.completeRemotesAndBranches(lastPart)...)
		}

	case "reset", "revert", "show", "diff":
		if !strings.HasPrefix(lastPart, "-") {
			completions = append(completions, e.completeCommitRefs(lastPart)...)
		}

	case "tag":
		if !strings.HasPrefix(lastPart, "-") {
			completions = append(completions, e.completeTags(lastPart)...)
		}

	case "add", "rm", "restore":
		if !strings.HasPrefix(lastPart, "-") {
			completions = append(completions, e.completeFiles(lastPart)...)
		}
	}

	return completions
}

// completeFlags completes git command flags.
func (e *CompletionEngine) completeFlags(subcommand string, partial string) []Completion {
	flags, ok := e.gitFlags[subcommand]
	if !ok {
		return nil
	}

	var completions []Completion
	for _, flag := range flags {
		if strings.HasPrefix(flag, partial) {
			completions = append(completions, Completion{
				Text:        flag,
				Description: getFlagDescription(subcommand, flag),
				Type:        "flag",
				Priority:    8,
			})
		}
	}

	return completions
}

// completeBranches completes branch names.
func (e *CompletionEngine) completeBranches(partial string) []Completion {
	branches := e.c.Model().Branches
	var completions []Completion

	for _, branch := range branches {
		if strings.HasPrefix(branch.Name, partial) || partial == "" {
			completions = append(completions, Completion{
				Text:        branch.Name,
				Description: "分支",
				Type:        "branch",
				Priority:    9,
			})
		}
	}

	return completions
}

// completeRemotesAndBranches completes remote names and branch names.
func (e *CompletionEngine) completeRemotesAndBranches(partial string) []Completion {
	var completions []Completion

	// Add remote names
	remotes := e.c.Model().Remotes
	for _, remote := range remotes {
		if strings.HasPrefix(remote.Name, partial) || partial == "" {
			completions = append(completions, Completion{
				Text:        remote.Name,
				Description: "远程仓库",
				Type:        "remote",
				Priority:    9,
			})
		}
	}

	// Add branch names
	completions = append(completions, e.completeBranches(partial)...)

	return completions
}

// completeCommitRefs completes commit references (hashes, HEAD~N, etc.).
func (e *CompletionEngine) completeCommitRefs(partial string) []Completion {
	var completions []Completion

	// Add special refs
	specialRefs := []string{"HEAD", "HEAD~1", "HEAD~2", "HEAD~3"}
	for _, ref := range specialRefs {
		if strings.HasPrefix(ref, partial) || partial == "" {
			completions = append(completions, Completion{
				Text:        ref,
				Description: "提交引用",
				Type:        "ref",
				Priority:    10,
			})
		}
	}

	// Add recent commit hashes
	commits := e.c.Model().Commits
	for i, commit := range commits {
		if i >= 10 { // Limit to recent 10 commits
			break
		}
		hash := commit.Hash()
		shortHash := hash[:7]
		if strings.HasPrefix(shortHash, partial) || strings.HasPrefix(hash, partial) || partial == "" {
			completions = append(completions, Completion{
				Text:        shortHash,
				Description: commit.Name,
				Type:        "commit",
				Priority:    8,
			})
		}
	}

	return completions
}

// completeTags completes tag names.
func (e *CompletionEngine) completeTags(partial string) []Completion {
	tags := e.c.Model().Tags
	var completions []Completion

	for _, tag := range tags {
		if strings.HasPrefix(tag.Name, partial) || partial == "" {
			completions = append(completions, Completion{
				Text:        tag.Name,
				Description: "标签",
				Type:        "tag",
				Priority:    8,
			})
		}
	}

	return completions
}

// completeFiles completes file paths from the current repository.
func (e *CompletionEngine) completeFiles(partial string) []Completion {
	files := e.c.Model().Files
	var completions []Completion

	for _, file := range files {
		if strings.HasPrefix(file.Path, partial) || partial == "" {
			completions = append(completions, Completion{
				Text:        file.Path,
				Description: getFileStatusDescription(file),
				Type:        "file",
				Priority:    7,
			})
		}
	}

	return completions
}

// parseCommandInput splits the input into command parts.
func parseCommandInput(input string) []string {
	if input == "" {
		return nil
	}

	// Simple split by spaces (doesn't handle quotes properly, but good enough for now)
	parts := strings.Fields(input)
	return parts
}

// GetGitSubcommands returns a list of common git subcommands.
func GetGitSubcommands() []string {
	return []string{
		"add", "bisect", "branch", "checkout", "cherry-pick", "clean",
		"clone", "commit", "diff", "fetch", "grep", "init", "log",
		"merge", "mv", "pull", "push", "rebase", "reset", "revert",
		"rm", "show", "stash", "status", "switch", "tag",
	}
}

// GetGitFlagCompletions returns common flags for each git subcommand.
func GetGitFlagCompletions() map[string][]string {
	return map[string][]string{
		"commit": {
			"--amend", "--no-edit", "-m", "-a", "--all", "--fixup",
			"--signoff", "-S", "--gpg-sign", "--no-verify", "--allow-empty",
		},
		"push": {
			"--force", "--force-with-lease", "--set-upstream", "-u",
			"--tags", "--delete", "--dry-run", "--all",
		},
		"pull": {
			"--rebase", "--no-rebase", "--ff-only", "--no-ff",
			"--tags", "--all", "--prune",
		},
		"checkout": {
			"-b", "-B", "--track", "--no-track", "-f", "--force",
			"--orphan", "--detach",
		},
		"branch": {
			"-d", "-D", "-m", "-M", "-a", "--all", "-r", "--remote",
			"--list", "--merged", "--no-merged",
		},
		"rebase": {
			"-i", "--interactive", "--continue", "--abort", "--skip",
			"--onto", "--autosquash", "--autostash",
		},
		"reset": {
			"--soft", "--mixed", "--hard", "--merge", "--keep",
		},
		"stash": {
			"save", "pop", "apply", "list", "show", "drop", "clear",
			"-u", "--include-untracked",
		},
		"log": {
			"--oneline", "--graph", "--all", "--author", "--since",
			"--until", "-n", "--follow", "--stat",
		},
		"diff": {
			"--cached", "--staged", "--stat", "--name-only",
			"--name-status", "--word-diff",
		},
		"tag": {
			"-a", "--annotate", "-m", "-d", "--delete", "-l", "--list",
		},
		"fetch": {
			"--all", "--prune", "--tags", "--dry-run",
		},
		"merge": {
			"--no-ff", "--ff-only", "--squash", "--abort", "--continue",
		},
	}
}

// Helper functions for descriptions

func getCommandDescription(cmd string) string {
	descriptions := map[string]string{
		"git":  "版本控制系统",
		"cd":   "切换目录",
		"ls":   "列出文件",
		"pwd":  "显示当前目录",
		"cat":  "显示文件内容",
		"grep": "搜索文本",
		"find": "查找文件",
	}
	return descriptions[cmd]
}

func getGitSubcommandDescription(subcmd string) string {
	descriptions := map[string]string{
		"add":         "添加文件到暂存区",
		"commit":      "提交更改",
		"push":        "推送到远程",
		"pull":        "从远程拉取",
		"checkout":    "切换分支",
		"switch":      "切换分支（新）",
		"branch":      "管理分支",
		"merge":       "合并分支",
		"rebase":      "变基",
		"reset":       "重置提交",
		"revert":      "反转提交",
		"stash":       "暂存工作区",
		"log":         "查看提交历史",
		"diff":        "查看差异",
		"status":      "查看状态",
		"tag":         "管理标签",
		"fetch":       "获取远程更新",
		"clone":       "克隆仓库",
		"init":        "初始化仓库",
		"clean":       "清理未跟踪文件",
		"cherry-pick": "挑选提交",
		"show":        "显示提交详情",
		"rm":          "删除文件",
		"mv":          "移动文件",
		"grep":        "搜索内容",
		"bisect":      "二分查找问题提交",
	}
	return descriptions[subcmd]
}

func getPriorityForSubcommand(subcmd string) int {
	// Most commonly used commands get higher priority
	highPriority := []string{"commit", "push", "pull", "checkout", "branch", "merge"}
	for _, cmd := range highPriority {
		if subcmd == cmd {
			return 10
		}
	}

	mediumPriority := []string{"add", "status", "log", "diff", "stash", "rebase"}
	for _, cmd := range mediumPriority {
		if subcmd == cmd {
			return 8
		}
	}

	return 6
}

func getFlagDescription(subcommand, flag string) string {
	descriptions := map[string]map[string]string{
		"commit": {
			"--amend":       "修改上次提交",
			"--no-edit":     "不修改提交消息",
			"-m":            "指定提交消息",
			"-a":            "提交所有已跟踪文件",
			"--all":         "提交所有已跟踪文件",
			"--fixup":       "创建 fixup 提交",
			"--signoff":     "添加 Signed-off-by 行",
			"-S":            "GPG 签名",
			"--no-verify":   "跳过 pre-commit 钩子",
			"--allow-empty": "允许空提交",
		},
		"push": {
			"--force":            "强制推送",
			"--force-with-lease": "安全强制推送",
			"--set-upstream":     "设置上游分支",
			"-u":                 "设置上游分支",
			"--tags":             "推送标签",
			"--delete":           "删除远程分支",
			"--dry-run":          "预览推送",
			"--all":              "推送所有分支",
		},
		"reset": {
			"--soft":  "保留暂存区和工作区",
			"--mixed": "保留工作区",
			"--hard":  "丢弃所有更改",
		},
	}

	if subcmdFlags, ok := descriptions[subcommand]; ok {
		if desc, ok := subcmdFlags[flag]; ok {
			return desc
		}
	}

	return ""
}

func getFileStatusDescription(file *models.File) string {
	if file.HasMergeConflicts {
		return "冲突"
	}
	if file.HasStagedChanges && file.HasUnstagedChanges {
		return "部分暂存"
	}
	if file.HasStagedChanges {
		return "已暂存"
	}
	if file.HasUnstagedChanges {
		return "已修改"
	}
	if !file.Tracked {
		return "未跟踪"
	}
	return "已跟踪"
}
