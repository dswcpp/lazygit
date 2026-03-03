package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dswcpp/lazygit/pkg/i18n"
)

// maxLinesPerFileDiff is the maximum number of diff lines to include per file
// before truncating. This prevents a single large file from drowning out other
// changes in the prompt.
const maxLinesPerFileDiff = 200

// filePriority represents a file block with its priority for sorting.
type filePriority struct {
	block    string
	path     string
	priority int
	lines    int
}

// aiDiffSkipGlobs lists exact filenames or glob patterns (matched against the
// base filename) whose diffs add no value for commit message generation:
// lock files, minified assets, generated protobuf code, etc.
var aiDiffSkipGlobs = []string{
	// Dependency lock files
	"yarn.lock",
	"package-lock.json",
	"pnpm-lock.yaml",
	"npm-shrinkwrap.json",
	"go.sum",
	"Cargo.lock",
	"Gemfile.lock",
	"poetry.lock",
	"composer.lock",
	"packages.lock.json",
	"*.lock",
	// Protobuf / gRPC generated
	"*.pb.go",
	"*.pb.ts",
	"*.pb.js",
	"*.pb.cc",
	"*.pb.h",
	// Minified assets
	"*.min.js",
	"*.min.css",
	// Other generated files
	"*.generated.go",
	"*.generated.ts",
	// Debug / vendor archives that end up in diffs
	"npm-debug.log",
	"yarn-debug.log",
	"yarn-error.log",
}

// FilterDiffForAI preprocesses a raw staged git diff before it is sent to the
// AI model. It:
//
//   - Replaces binary-file blocks with a one-line note.
//   - Replaces lock-file / generated-file blocks with a one-line note.
//   - Sorts files by priority (core business logic first, docs last).
//   - Adds semantic summaries for each file ([新增]/[修改]/[删除]).
//   - Intelligently truncates large files while preserving key signatures.
//   - Prepends an enhanced summary header with statistics.
//
// The returned string is suitable for direct inclusion in the AI prompt.
func FilterDiffForAI(rawDiff string, tr *i18n.TranslationSet) string {
	if strings.TrimSpace(rawDiff) == "" {
		return rawDiff
	}

	blocks := splitByFileHeader(rawDiff)

	// Step 1: Parse and prioritize blocks
	prioritizedBlocks := prioritizeFileBlocks(blocks)

	// Step 2: Process blocks with semantic summaries and smart truncation
	var out strings.Builder
	included, skipped := 0, 0
	totalInsertions, totalDeletions := 0, 0
	filesByExt := make(map[string]int)
	largestChanges := []fileChange{}

	for _, pb := range prioritizedBlocks {
		if strings.TrimSpace(pb.block) == "" {
			continue
		}

		// Check if should skip
		if reason := diffSkipReason(pb.block, pb.path, tr); reason != "" {
			fmt.Fprintf(&out, tr.AIDiffSkipped+"\n", pb.path, reason)
			skipped++
			continue
		}

		// Count file extension
		ext := getFileExtension(pb.path)
		filesByExt[ext]++

		// Count changes
		insertions, deletions := countDiffChanges(pb.block)
		totalInsertions += insertions
		totalDeletions += deletions

		// Track largest changes
		if insertions+deletions > 50 {
			largestChanges = append(largestChanges, fileChange{
				path:       pb.path,
				insertions: insertions,
				deletions:  deletions,
			})
		}

		// Add semantic summary
		summary := generateSemanticSummary(pb.block, pb.path, insertions, deletions, tr)
		out.WriteString(summary)
		out.WriteString("\n")

		// Smart truncate and add content
		processed := smartTruncateBlock(pb.block, pb.path, tr)
		out.WriteString(processed)
		out.WriteString("\n")

		included++
	}

	// Step 3: Generate enhanced header
	header := generateEnhancedHeader(included, skipped, totalInsertions, totalDeletions, filesByExt, largestChanges, tr)

	return header + out.String()
}

// splitByFileHeader splits a unified diff into per-file blocks.
// Each block starts with a "diff --git " line.
func splitByFileHeader(diff string) []string {
	var blocks []string
	var cur strings.Builder

	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "diff --git ") && cur.Len() > 0 {
			blocks = append(blocks, cur.String())
			cur.Reset()
		}
		cur.WriteString(line)
		cur.WriteByte('\n')
	}
	if cur.Len() > 0 {
		blocks = append(blocks, cur.String())
	}
	return blocks
}

// parseGitDiffFilePath extracts the destination path from a "diff --git a/X b/Y" line.
func parseGitDiffFilePath(block string) string {
	first := strings.SplitN(block, "\n", 2)[0]
	// Format: "diff --git a/<path> b/<path>" — take everything after " b/"
	if i := strings.Index(first, " b/"); i >= 0 {
		return first[i+3:]
	}
	return first
}

// diffSkipReason returns a short reason string when the block should be omitted,
// or an empty string when the block should be included in the prompt.
func diffSkipReason(block, filePath string, tr *i18n.TranslationSet) string {
	if strings.Contains(block, "Binary files") || strings.Contains(block, "GIT binary patch") {
		return tr.AIDiffBinaryFile
	}
	base := baseName(filePath)
	for _, pattern := range aiDiffSkipGlobs {
		if matched, _ := matchGlob(pattern, base); matched {
			return tr.AIDiffLockOrGeneratedFile
		}
	}
	return ""
}

// limitFileBlock truncates a file block to at most maxLinesPerFileDiff lines,
// appending a note when truncation occurs.
func limitFileBlock(block, filePath string, tr *i18n.TranslationSet) string {
	lines := strings.Split(block, "\n")
	if len(lines) <= maxLinesPerFileDiff {
		return block
	}
	truncated := strings.Join(lines[:maxLinesPerFileDiff], "\n")
	truncated += "\n" + fmt.Sprintf(tr.AIDiffTruncated, filePath, len(lines), maxLinesPerFileDiff) + "\n"
	return truncated
}

// baseName returns the last path component of a forward-slash or
// backslash separated path (git paths always use forward slashes).
func baseName(path string) string {
	if i := strings.LastIndexAny(path, "/\\"); i >= 0 {
		return path[i+1:]
	}
	return path
}

// matchGlob matches a shell glob pattern against a name.
// Supports the same wildcards as filepath.Match: *, ?, [range].
func matchGlob(pattern, name string) (bool, error) {
	// Fast path: no wildcard → direct equality check.
	if !strings.ContainsAny(pattern, "*?[") {
		return pattern == name, nil
	}
	return globMatch(pattern, name)
}

// globMatch is a minimal glob matcher that handles *, ?, and [range].
// It avoids importing path/filepath to stay OS-agnostic with / separators.
func globMatch(pattern, name string) (bool, error) {
	for len(pattern) > 0 {
		var chunk string
		chunk, pattern, _ = strings.Cut(pattern, "*")
		if chunk == "" {
			if pattern == "" {
				return true, nil
			}
			// '*' consumed: try matching the rest from each position in name.
			for i := 0; i <= len(name); i++ {
				if ok, _ := globMatch(pattern, name[i:]); ok {
					return true, nil
				}
			}
			return false, nil
		}
		// Match chunk literally, character by character (handles ? and [range] naively).
		if len(name) < len(chunk) {
			return false, nil
		}
		for i, pc := range chunk {
			if i >= len(name) {
				return false, nil
			}
			nc := rune(name[i])
			if pc != '?' && pc != nc {
				return false, nil
			}
		}
		name = name[len(chunk):]
	}
	return name == "", nil
}

// getFileExtension returns the file extension including the dot (e.g., ".go", ".ts")
// or "other" if no extension is found.
func getFileExtension(path string) string {
	// Get the base name first
	base := baseName(path)

	// Check if it's a hidden file (starts with .)
	if strings.HasPrefix(base, ".") {
		return "other"
	}

	if idx := strings.LastIndex(base, "."); idx >= 0 {
		return base[idx:]
	}
	return "other"
}

// countDiffChanges counts the number of insertions and deletions in a diff block.
func countDiffChanges(block string) (insertions, deletions int) {
	for _, line := range strings.Split(block, "\n") {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			insertions++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			deletions++
		}
	}
	return
}

// fileChange represents a file with its change statistics.
type fileChange struct {
	path       string
	insertions int
	deletions  int
}

// prioritizeFileBlocks sorts file blocks by priority.
// Priority order: core business logic > API/controllers > config > tests > docs
func prioritizeFileBlocks(blocks []string) []filePriority {
	priorities := make([]filePriority, 0, len(blocks))

	for _, block := range blocks {
		if strings.TrimSpace(block) == "" {
			continue
		}

		path := parseGitDiffFilePath(block)
		priority := calculateFilePriority(path)
		lines := len(strings.Split(block, "\n"))

		priorities = append(priorities, filePriority{
			block:    block,
			path:     path,
			priority: priority,
			lines:    lines,
		})
	}

	// Sort by priority (descending)
	for i := 0; i < len(priorities); i++ {
		for j := i + 1; j < len(priorities); j++ {
			if priorities[j].priority > priorities[i].priority {
				priorities[i], priorities[j] = priorities[j], priorities[i]
			}
		}
	}

	return priorities
}

// calculateFilePriority assigns a priority score to a file path.
// Higher score = higher priority (processed first).
func calculateFilePriority(path string) int {
	lowerPath := strings.ToLower(path)

	// Highest priority: core business logic
	// Match /pkg/, /src/, /lib/, /internal/, /core/ anywhere in path
	if strings.Contains(lowerPath, "/pkg/") || strings.HasPrefix(lowerPath, "pkg/") ||
		strings.Contains(lowerPath, "/src/") || strings.HasPrefix(lowerPath, "src/") ||
		strings.Contains(lowerPath, "/lib/") || strings.HasPrefix(lowerPath, "lib/") ||
		strings.Contains(lowerPath, "/internal/") || strings.HasPrefix(lowerPath, "internal/") ||
		strings.Contains(lowerPath, "/core/") || strings.HasPrefix(lowerPath, "core/") {
		return 100
	}

	// High priority: API/controllers/handlers/services
	if strings.Contains(lowerPath, "/api/") || strings.HasPrefix(lowerPath, "api/") ||
		strings.Contains(lowerPath, "/controller/") || strings.HasPrefix(lowerPath, "controller/") ||
		strings.Contains(lowerPath, "/handler/") || strings.HasPrefix(lowerPath, "handler/") ||
		strings.Contains(lowerPath, "/service/") || strings.HasPrefix(lowerPath, "service/") ||
		strings.Contains(lowerPath, "/router/") || strings.HasPrefix(lowerPath, "router/") {
		return 90
	}

	// Medium-high priority: models/schemas
	if strings.Contains(lowerPath, "/model/") || strings.HasPrefix(lowerPath, "model/") ||
		strings.Contains(lowerPath, "/schema/") || strings.HasPrefix(lowerPath, "schema/") ||
		strings.Contains(lowerPath, "/entity/") || strings.HasPrefix(lowerPath, "entity/") ||
		strings.Contains(lowerPath, "/dto/") || strings.HasPrefix(lowerPath, "dto/") {
		return 80
	}

	// Medium priority: config files
	if strings.HasSuffix(lowerPath, ".yaml") || strings.HasSuffix(lowerPath, ".yml") ||
		strings.HasSuffix(lowerPath, ".json") || strings.HasSuffix(lowerPath, ".toml") ||
		strings.HasSuffix(lowerPath, ".env") {
		return 70
	}

	// Low priority: test files
	if strings.Contains(lowerPath, "_test.") || strings.Contains(lowerPath, ".test.") ||
		strings.Contains(lowerPath, "/test/") || strings.Contains(lowerPath, "/tests/") ||
		strings.HasPrefix(lowerPath, "test/") || strings.HasPrefix(lowerPath, "tests/") ||
		strings.HasSuffix(lowerPath, "_spec.") {
		return 50
	}

	// Lowest priority: documentation
	if strings.HasSuffix(lowerPath, ".md") || strings.HasSuffix(lowerPath, ".txt") ||
		strings.HasSuffix(lowerPath, ".rst") || strings.Contains(lowerPath, "/docs/") ||
		strings.Contains(lowerPath, "/doc/") || strings.HasPrefix(lowerPath, "docs/") ||
		strings.HasPrefix(lowerPath, "doc/") {
		return 30
	}

	// Default priority
	return 60
}

// generateSemanticSummary creates a semantic summary for a file change.
func generateSemanticSummary(block, path string, insertions, deletions int, tr *i18n.TranslationSet) string {
	if isNewFile(block) {
		return fmt.Sprintf("### %s %s", tr.AIDiffNewFile, path)
	}
	if isDeletedFile(block) {
		return fmt.Sprintf("### %s %s", tr.AIDiffDeletedFile, path)
	}
	if isRenamed(block) {
		oldPath := extractRenamedFrom(block)
		if oldPath != "" {
			return fmt.Sprintf("### %s %s <- %s", tr.AIDiffRenamedFile, path, oldPath)
		}
		return fmt.Sprintf("### %s %s", tr.AIDiffRenamedFile, path)
	}

	return fmt.Sprintf("### %s %s (+%d/-%d)", tr.AIDiffModifiedFile, path, insertions, deletions)
}

// smartTruncateBlock intelligently truncates a file block while preserving key information.
func smartTruncateBlock(block, path string, tr *i18n.TranslationSet) string {
	lines := strings.Split(block, "\n")
	if len(lines) <= maxLinesPerFileDiff {
		return block
	}

	// Preserve:
	// 1. File header (first 5 lines) - contains diff metadata
	// 2. Function/class signatures
	// 3. Important comments (TODO, FIXME, NOTE, etc.)

	preserved := []string{}
	headerLines := 5
	if len(lines) < headerLines {
		headerLines = len(lines)
	}
	preserved = append(preserved, lines[:headerLines]...)

	// Patterns for important lines
	signaturePattern := regexp.MustCompile(`^[+-]?\s*(func|class|def|interface|type|struct|const|var|export|public|private|protected)\s+`)
	importantCommentPattern := regexp.MustCompile(`^[+-]?\s*(//|#|/\*)\s*(TODO|FIXME|NOTE|IMPORTANT|WARNING|BUG|HACK)`)

	// Scan remaining lines for important content
	for i := headerLines; i < len(lines) && len(preserved) < maxLinesPerFileDiff; i++ {
		line := lines[i]
		if signaturePattern.MatchString(line) || importantCommentPattern.MatchString(line) {
			preserved = append(preserved, line)
		} else if len(preserved) < maxLinesPerFileDiff/2 {
			// Include some context lines in the first half
			preserved = append(preserved, line)
		}
	}

	result := strings.Join(preserved, "\n")
	if len(preserved) < len(lines) {
		result += "\n" + fmt.Sprintf(tr.AIDiffSmartTruncated, path, len(preserved), len(lines)) + "\n"
	}

	return result
}

// generateEnhancedHeader generates an enhanced summary header with statistics.
func generateEnhancedHeader(included, skipped, totalInsertions, totalDeletions int,
	filesByExt map[string]int, largestChanges []fileChange, tr *i18n.TranslationSet) string {

	var header strings.Builder

	header.WriteString(tr.AIDiffChangeStats + "\n")
	header.WriteString(fmt.Sprintf(tr.AIDiffFilesCount, included))
	if skipped > 0 {
		header.WriteString(fmt.Sprintf(tr.AIDiffSkippedFilesNote, skipped))
	}
	header.WriteString("\n")

	// File types
	if len(filesByExt) > 0 {
		header.WriteString(tr.AIDiffFileTypes)
		first := true
		for ext, count := range filesByExt {
			if !first {
				header.WriteString(", ")
			}
			fmt.Fprintf(&header, "%s: %d", ext, count)
			first = false
		}
		header.WriteString("\n")
	}

	// Change scale
	header.WriteString(fmt.Sprintf(tr.AIDiffChangeScale, totalInsertions, totalDeletions) + "\n")

	// Largest changes (top 3)
	if len(largestChanges) > 0 {
		// Sort by total changes
		for i := 0; i < len(largestChanges); i++ {
			for j := i + 1; j < len(largestChanges); j++ {
				if largestChanges[j].insertions+largestChanges[j].deletions >
					largestChanges[i].insertions+largestChanges[i].deletions {
					largestChanges[i], largestChanges[j] = largestChanges[j], largestChanges[i]
				}
			}
		}

		// Take top 3
		limit := 3
		if len(largestChanges) < limit {
			limit = len(largestChanges)
		}

		header.WriteString(tr.AIDiffMajorChanges)
		for i := 0; i < limit; i++ {
			if i > 0 {
				header.WriteString(", ")
			}
			fc := largestChanges[i]
			fmt.Fprintf(&header, "%s (+%d/-%d)", fc.path, fc.insertions, fc.deletions)
		}
		header.WriteString("\n")
	}

	header.WriteString("\n")
	return header.String()
}

// isNewFile checks if the block represents a new file.
func isNewFile(block string) bool {
	return strings.Contains(block, "new file mode")
}

// isDeletedFile checks if the block represents a deleted file.
func isDeletedFile(block string) bool {
	return strings.Contains(block, "deleted file mode")
}

// isRenamed checks if the block represents a renamed file.
func isRenamed(block string) bool {
	return strings.Contains(block, "rename from") || strings.Contains(block, "rename to")
}

// extractRenamedFrom extracts the old path from a rename diff.
func extractRenamedFrom(block string) string {
	for _, line := range strings.Split(block, "\n") {
		if strings.HasPrefix(line, "rename from ") {
			return strings.TrimPrefix(line, "rename from ")
		}
	}
	return ""
}
