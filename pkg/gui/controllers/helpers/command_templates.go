package helpers

import (
	"sort"
	"strings"
)

// CommandTemplate represents a predefined Git command template with placeholders.
type CommandTemplate struct {
	Category     string   // Category: commit, branch, remote, stash, reset, log, clean, tag
	Command      string   // Command template with {placeholder} syntax
	Description  string   // Chinese description
	Placeholders []string // Placeholder descriptions for user guidance
	Icon         string   // Icon for visual identification
	Priority     int      // Priority for sorting (higher = more important)
}

// TemplateEngine manages Git command templates and provides search/filtering.
type TemplateEngine struct {
	templates []CommandTemplate
	categories map[string][]CommandTemplate
}

// NewTemplateEngine creates a new template engine with predefined Git commands.
func NewTemplateEngine() *TemplateEngine {
	templates := GetDefaultGitTemplates()

	engine := &TemplateEngine{
		templates:  templates,
		categories: make(map[string][]CommandTemplate),
	}

	// Index by category for faster lookups
	for _, tmpl := range templates {
		engine.categories[tmpl.Category] = append(engine.categories[tmpl.Category], tmpl)
	}

	return engine
}

// GetDefaultGitTemplates returns a comprehensive list of common Git command templates.
func GetDefaultGitTemplates() []CommandTemplate {
	return []CommandTemplate{
		// === Commit Operations ===
		{
			Category:     "commit",
			Command:      "git commit -m \"{message}\"",
			Description:  "提交暂存的更改",
			Placeholders: []string{"message: 提交消息"},
			Icon:         "📝",
			Priority:     10,
		},
		{
			Category:    "commit",
			Command:     "git commit --amend --no-edit",
			Description: "修改最后一次提交（不改消息）",
			Icon:        "✏️",
			Priority:    9,
		},
		{
			Category:     "commit",
			Command:      "git commit -a -m \"{message}\"",
			Description:  "提交所有已跟踪文件的更改",
			Placeholders: []string{"message: 提交消息"},
			Icon:         "📦",
			Priority:     8,
		},
		{
			Category:     "commit",
			Command:      "git commit --amend -m \"{message}\"",
			Description:  "修改最后一次提交消息",
			Placeholders: []string{"message: 提交消息"},
			Icon:         "✏️",
			Priority:     7,
		},
		{
			Category:    "commit",
			Command:     "git commit --amend",
			Description: "修改最后一次提交（在编辑器中修改）",
			Icon:        "✏️",
			Priority:    6,
		},

		// === Branch Operations ===
		{
			Category:     "branch",
			Command:      "git checkout -b {branch-name}",
			Description:  "创建并切换到新分支",
			Placeholders: []string{"branch-name: 分支名"},
			Icon:         "🌿",
			Priority:     10,
		},
		{
			Category:     "branch",
			Command:      "git branch {branch-name}",
			Description:  "创建新分支",
			Placeholders: []string{"branch-name: 分支名"},
			Icon:         "🌿",
			Priority:     8,
		},
		{
			Category:     "branch",
			Command:      "git branch -d {branch-name}",
			Description:  "删除已合并的分支",
			Placeholders: []string{"branch-name: 分支名"},
			Icon:         "🗑️",
			Priority:     5,
		},
		{
			Category:     "branch",
			Command:      "git branch -D {branch-name}",
			Description:  "强制删除分支（危险）",
			Placeholders: []string{"branch-name: 分支名"},
			Icon:         "⚠️",
			Priority:     3,
		},
		{
			Category:    "branch",
			Command:     "git branch -m {old-name} {new-name}",
			Description: "重命名分支",
			Placeholders: []string{"old-name: 旧分支名", "new-name: 新分支名"},
			Icon:        "📝",
			Priority:    6,
		},

		// === Remote Operations ===
		{
			Category:     "remote",
			Command:      "git push origin {branch-name}",
			Description:  "推送分支到远程",
			Placeholders: []string{"branch-name: 分支名"},
			Icon:         "⬆️",
			Priority:     10,
		},
		{
			Category:    "remote",
			Command:     "git push --force-with-lease",
			Description: "安全的强制推送",
			Icon:        "⚡",
			Priority:    5,
		},
		{
			Category:    "remote",
			Command:     "git push --force",
			Description: "强制推送（危险，可能覆盖他人的提交）",
			Icon:        "⚠️",
			Priority:    2,
		},
		{
			Category:    "remote",
			Command:     "git pull --rebase",
			Description: "使用 rebase 方式拉取",
			Icon:        "⬇️",
			Priority:    8,
		},
		{
			Category:    "remote",
			Command:     "git fetch --all --prune",
			Description: "获取所有远程分支并清理已删除的引用",
			Icon:        "🔄",
			Priority:    7,
		},

		// === Stash Operations ===
		{
			Category:     "stash",
			Command:      "git stash save \"{message}\"",
			Description:  "保存当前工作区到 stash",
			Placeholders: []string{"message: 提交消息"},
			Icon:         "💾",
			Priority:     10,
		},
		{
			Category:    "stash",
			Command:     "git stash pop",
			Description: "恢复并删除最近的 stash",
			Icon:        "📤",
			Priority:    9,
		},
		{
			Category:    "stash",
			Command:     "git stash apply",
			Description: "应用最近的 stash（不删除）",
			Icon:        "📋",
			Priority:    7,
		},
		{
			Category:     "stash",
			Command:      "git stash apply stash@{{n}}",
			Description:  "应用指定的 stash",
			Placeholders: []string{"n: stash 索引"},
			Icon:         "📋",
			Priority:     5,
		},
		{
			Category:    "stash",
			Command:     "git stash list",
			Description: "查看所有 stash 列表",
			Icon:        "📜",
			Priority:    6,
		},

		// === Reset/Revert Operations ===
		{
			Category:    "reset",
			Command:     "git reset HEAD~1",
			Description: "撤销最后一次提交（保留更改）",
			Icon:        "↩️",
			Priority:    8,
		},
		{
			Category:    "reset",
			Command:     "git reset --soft HEAD~1",
			Description: "撤销提交但保留暂存区",
			Icon:        "↩️",
			Priority:    7,
		},
		{
			Category:    "reset",
			Command:     "git reset --hard HEAD~1",
			Description: "撤销提交并丢弃所有更改（危险）",
			Icon:        "⚠️",
			Priority:    3,
		},
		{
			Category:     "reset",
			Command:      "git revert {commit-hash}",
			Description:  "反转指定提交（创建新提交）",
			Placeholders: []string{"commit-hash: 提交哈希"},
			Icon:         "🔄",
			Priority:     7,
		},
		{
			Category:    "reset",
			Command:     "git reset --hard origin/{branch-name}",
			Description: "强制重置到远程分支状态（危险）",
			Placeholders: []string{"branch-name: 分支名"},
			Icon:        "⚠️",
			Priority:    2,
		},

		// === Log/Diff Operations ===
		{
			Category:    "log",
			Command:     "git log --oneline --graph --all -n 20",
			Description: "查看图形化提交历史（最近 20 条）",
			Icon:        "📊",
			Priority:    10,
		},
		{
			Category:     "log",
			Command:      "git log --author=\"{author}\" --oneline",
			Description:  "查看指定作者的提交",
			Placeholders: []string{"author: 作者名称"},
			Icon:         "👤",
			Priority:     5,
		},
		{
			Category:    "log",
			Command:     "git log --since=\"{date}\" --oneline",
			Description: "查看指定日期后的提交",
			Placeholders: []string{"date: 日期，如 2024-01-01 或 1.week.ago"},
			Icon:        "📅",
			Priority:    4,
		},
		{
			Category:    "diff",
			Command:     "git diff HEAD~1",
			Description: "查看与上次提交的差异",
			Icon:        "🔍",
			Priority:    8,
		},
		{
			Category:     "diff",
			Command:      "git diff {branch1}..{branch2}",
			Description:  "比较两个分支的差异",
			Placeholders: []string{"branch1: 分支1", "branch2: 分支2"},
			Icon:         "🔍",
			Priority:     6,
		},

		// === Clean Operations ===
		{
			Category:    "clean",
			Command:     "git clean -fd",
			Description: "删除未跟踪的文件和目录",
			Icon:        "🧹",
			Priority:    5,
		},
		{
			Category:    "clean",
			Command:     "git clean -fdx",
			Description: "删除未跟踪和忽略的文件（危险）",
			Icon:        "⚠️",
			Priority:    2,
		},
		{
			Category:    "clean",
			Command:     "git clean -fdn",
			Description: "预览将被删除的文件（不实际删除）",
			Icon:        "👁️",
			Priority:    7,
		},

		// === Tag Operations ===
		{
			Category:     "tag",
			Command:      "git tag {tag-name}",
			Description:  "创建轻量标签",
			Placeholders: []string{"tag-name: 标签名"},
			Icon:         "🏷️",
			Priority:     7,
		},
		{
			Category:     "tag",
			Command:      "git tag -a {tag-name} -m \"{message}\"",
			Description:  "创建带注释的标签",
			Placeholders: []string{"tag-name: 标签名", "message: 提交消息"},
			Icon:         "🏷️",
			Priority:     6,
		},
		{
			Category:    "tag",
			Command:     "git push origin --tags",
			Description: "推送所有标签到远程",
			Icon:        "⬆️",
			Priority:    5,
		},
		{
			Category:     "tag",
			Command:      "git tag -d {tag-name}",
			Description:  "删除本地标签",
			Placeholders: []string{"tag-name: 标签名"},
			Icon:         "🗑️",
			Priority:     4,
		},
	}
}

// Search returns templates matching the search query.
// Uses fuzzy matching for flexible search.
func (e *TemplateEngine) Search(query string) []CommandTemplate {
	if query == "" {
		return e.GetAll()
	}

	query = strings.ToLower(query)
	var results []CommandTemplate

	for _, tmpl := range e.templates {
		score := e.matchScore(tmpl, query)
		if score > 0 {
			results = append(results, tmpl)
		}
	}

	// Sort by priority (higher first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Priority > results[j].Priority
	})

	return results
}

// matchScore calculates how well a template matches the query.
// Returns 0 if no match, higher scores for better matches.
func (e *TemplateEngine) matchScore(tmpl CommandTemplate, query string) float64 {
	cmdLower := strings.ToLower(tmpl.Command)
	descLower := strings.ToLower(tmpl.Description)

	// Exact command prefix match (highest priority)
	if strings.HasPrefix(cmdLower, query) {
		return 100.0
	}

	// Command contains query
	if strings.Contains(cmdLower, query) {
		return 80.0
	}

	// Description contains query
	if strings.Contains(descLower, query) {
		return 60.0
	}

	// Category matches
	if strings.ToLower(tmpl.Category) == query {
		return 50.0
	}

	// Fuzzy match in command
	if e.fuzzyMatch(cmdLower, query) {
		return 40.0
	}

	return 0
}

// fuzzyMatch performs simple fuzzy matching.
func (e *TemplateEngine) fuzzyMatch(text, pattern string) bool {
	if len(pattern) == 0 {
		return true
	}

	patternIdx := 0
	for _, char := range text {
		if patternIdx < len(pattern) && char == rune(pattern[patternIdx]) {
			patternIdx++
		}
	}

	return patternIdx == len(pattern)
}

// GetByCategory returns all templates in a specific category.
func (e *TemplateEngine) GetByCategory(category string) []CommandTemplate {
	return e.categories[category]
}

// GetAll returns all templates sorted by priority.
func (e *TemplateEngine) GetAll() []CommandTemplate {
	results := make([]CommandTemplate, len(e.templates))
	copy(results, e.templates)

	sort.Slice(results, func(i, j int) bool {
		return results[i].Priority > results[j].Priority
	})

	return results
}

// GetCategories returns all available categories.
func (e *TemplateEngine) GetCategories() []string {
	categories := make([]string, 0, len(e.categories))
	for cat := range e.categories {
		categories = append(categories, cat)
	}
	return categories
}

// FillPlaceholders replaces placeholders in a command template with actual values.
// Returns the filled command and the position of the first unfilled placeholder.
func FillPlaceholders(command string, values map[string]string) (filled string, nextPlaceholder int) {
	filled = command
	nextPlaceholder = -1

	// Find all placeholders {name}
	start := 0
	for {
		openIdx := strings.Index(filled[start:], "{")
		if openIdx == -1 {
			break
		}
		openIdx += start

		closeIdx := strings.Index(filled[openIdx:], "}")
		if closeIdx == -1 {
			break
		}
		closeIdx += openIdx

		placeholder := filled[openIdx+1 : closeIdx]
		if value, ok := values[placeholder]; ok {
			filled = filled[:openIdx] + value + filled[closeIdx+1:]
			start = openIdx + len(value)
		} else {
			// Unfilled placeholder found
			if nextPlaceholder == -1 {
				nextPlaceholder = openIdx
			}
			start = closeIdx + 1
		}
	}

	return filled, nextPlaceholder
}

// ExtractPlaceholders extracts all placeholder names from a command template.
func ExtractPlaceholders(command string) []string {
	var placeholders []string
	start := 0

	for {
		openIdx := strings.Index(command[start:], "{")
		if openIdx == -1 {
			break
		}
		openIdx += start

		closeIdx := strings.Index(command[openIdx:], "}")
		if closeIdx == -1 {
			break
		}
		closeIdx += openIdx

		placeholder := command[openIdx+1 : closeIdx]
		placeholders = append(placeholders, placeholder)
		start = closeIdx + 1
	}

	return placeholders
}
