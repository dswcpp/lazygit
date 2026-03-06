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
				Description: e.getCommandDescription(cmd),
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
				Description: e.getGitSubcommandDescription(subcmd),
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
				Description: e.getFlagDescription(subcommand, flag),
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
				Description: e.c.Tr.CompletionBranch,
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
				Description: e.c.Tr.CompletionRemote,
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
				Description: e.c.Tr.CompletionCommitRef,
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
				Description: e.c.Tr.CompletionTag,
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
				Description: e.getFileStatusDescription(file),
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

func (e *CompletionEngine) getCommandDescription(cmd string) string {
	descriptions := map[string]string{
		"git":  e.c.Tr.CompletionGitDesc,
		"cd":   e.c.Tr.CompletionCdDesc,
		"ls":   e.c.Tr.CompletionLsDesc,
		"pwd":  e.c.Tr.CompletionPwdDesc,
		"cat":  e.c.Tr.CompletionCatDesc,
		"grep": e.c.Tr.CompletionGrepDesc,
		"find": e.c.Tr.CompletionFindDesc,
	}
	return descriptions[cmd]
}

func (e *CompletionEngine) getGitSubcommandDescription(subcmd string) string {
	descriptions := map[string]string{
		"add":         e.c.Tr.CompletionGitAddDesc,
		"commit":      e.c.Tr.CompletionGitCommitDesc,
		"push":        e.c.Tr.CompletionGitPushDesc,
		"pull":        e.c.Tr.CompletionGitPullDesc,
		"checkout":    e.c.Tr.CompletionGitCheckoutDesc,
		"switch":      e.c.Tr.CompletionGitSwitchDesc,
		"branch":      e.c.Tr.CompletionGitBranchDesc,
		"merge":       e.c.Tr.CompletionGitMergeDesc,
		"rebase":      e.c.Tr.CompletionGitRebaseDesc,
		"reset":       e.c.Tr.CompletionGitResetDesc,
		"revert":      e.c.Tr.CompletionGitRevertDesc,
		"stash":       e.c.Tr.CompletionGitStashDesc,
		"log":         e.c.Tr.CompletionGitLogDesc,
		"diff":        e.c.Tr.CompletionGitDiffDesc,
		"status":      e.c.Tr.CompletionGitStatusDesc,
		"tag":         e.c.Tr.CompletionGitTagDesc,
		"fetch":       e.c.Tr.CompletionGitFetchDesc,
		"clone":       e.c.Tr.CompletionGitCloneDesc,
		"init":        e.c.Tr.CompletionGitInitDesc,
		"clean":       e.c.Tr.CompletionGitCleanDesc,
		"cherry-pick": e.c.Tr.CompletionGitCherryPickDesc,
		"show":        e.c.Tr.CompletionGitShowDesc,
		"rm":          e.c.Tr.CompletionGitRmDesc,
		"mv":          e.c.Tr.CompletionGitMvDesc,
		"grep":        e.c.Tr.CompletionGitGrepDesc,
		"bisect":      e.c.Tr.CompletionGitBisectDesc,
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

func (e *CompletionEngine) getFlagDescription(subcommand, flag string) string {
	descriptions := map[string]map[string]string{
		"commit": {
			"--amend":       e.c.Tr.CompletionFlagAmendDesc,
			"--no-edit":     e.c.Tr.CompletionFlagNoEditDesc,
			"-m":            e.c.Tr.CompletionFlagMDesc,
			"-a":            e.c.Tr.CompletionFlagADesc,
			"--all":         e.c.Tr.CompletionFlagAllDesc,
			"--fixup":       e.c.Tr.CompletionFlagFixupDesc,
			"--signoff":     e.c.Tr.CompletionFlagSignoffDesc,
			"-S":            e.c.Tr.CompletionFlagSDesc,
			"--no-verify":   e.c.Tr.CompletionFlagNoVerifyDesc,
			"--allow-empty": e.c.Tr.CompletionFlagAllowEmptyDesc,
		},
		"push": {
			"--force":            e.c.Tr.CompletionFlagForceDesc,
			"--force-with-lease": e.c.Tr.CompletionFlagForceWithLeaseDesc,
			"--set-upstream":     e.c.Tr.CompletionFlagSetUpstreamDesc,
			"-u":                 e.c.Tr.CompletionFlagUDesc,
			"--tags":             e.c.Tr.CompletionFlagTagsDesc,
			"--delete":           e.c.Tr.CompletionFlagDeleteDesc,
			"--dry-run":          e.c.Tr.CompletionFlagDryRunDesc,
			"--all":              e.c.Tr.CompletionFlagAllBranchesDesc,
		},
		"reset": {
			"--soft":  e.c.Tr.CompletionFlagSoftDesc,
			"--mixed": e.c.Tr.CompletionFlagMixedDesc,
			"--hard":  e.c.Tr.CompletionFlagHardDesc,
		},
	}

	if subcmdFlags, ok := descriptions[subcommand]; ok {
		if desc, ok := subcmdFlags[flag]; ok {
			return desc
		}
	}

	return ""
}

func (e *CompletionEngine) getFileStatusDescription(file *models.File) string {
	if file.HasMergeConflicts {
		return e.c.Tr.CompletionStatusConflicted
	}
	if file.HasStagedChanges && file.HasUnstagedChanges {
		return e.c.Tr.CompletionStatusPartiallyStaged
	}
	if file.HasStagedChanges {
		return e.c.Tr.CompletionStatusStaged
	}
	if file.HasUnstagedChanges {
		return e.c.Tr.CompletionStatusModified
	}
	if !file.Tracked {
		return e.c.Tr.CompletionStatusUntracked
	}
	return e.c.Tr.CompletionStatusTracked
}
