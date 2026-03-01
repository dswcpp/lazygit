package helpers

import (
	"fmt"
	"strings"
)

// maxLinesPerFileDiff is the maximum number of diff lines to include per file
// before truncating. This prevents a single large file from drowning out other
// changes in the prompt.
const maxLinesPerFileDiff = 200

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
//   - Truncates any file block that exceeds maxLinesPerFileDiff lines.
//   - Prepends a short summary header so the model knows what was skipped.
//
// The returned string is suitable for direct inclusion in the AI prompt.
func FilterDiffForAI(rawDiff string) string {
	if strings.TrimSpace(rawDiff) == "" {
		return rawDiff
	}

	blocks := splitByFileHeader(rawDiff)

	var out strings.Builder
	included, skipped := 0, 0

	for _, block := range blocks {
		if strings.TrimSpace(block) == "" {
			continue
		}

		filePath := parseGitDiffFilePath(block)

		if reason := diffSkipReason(block, filePath); reason != "" {
			fmt.Fprintf(&out, "[跳过 %s: %s]\n", filePath, reason)
			skipped++
			continue
		}

		out.WriteString(limitFileBlock(block, filePath))
		included++
	}

	// Prepend a one-line summary so the model has context.
	var header strings.Builder
	fmt.Fprintf(&header, "# 变更文件数: %d", included)
	if skipped > 0 {
		fmt.Fprintf(&header, "（另有 %d 个锁文件/二进制/生成文件已略去）", skipped)
	}
	header.WriteString("\n\n")

	return header.String() + out.String()
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
func diffSkipReason(block, filePath string) string {
	if strings.Contains(block, "Binary files") || strings.Contains(block, "GIT binary patch") {
		return "二进制文件"
	}
	base := baseName(filePath)
	for _, pattern := range aiDiffSkipGlobs {
		if matched, _ := matchGlob(pattern, base); matched {
			return "锁文件/生成文件"
		}
	}
	return ""
}

// limitFileBlock truncates a file block to at most maxLinesPerFileDiff lines,
// appending a note when truncation occurs.
func limitFileBlock(block, filePath string) string {
	lines := strings.Split(block, "\n")
	if len(lines) <= maxLinesPerFileDiff {
		return block
	}
	truncated := strings.Join(lines[:maxLinesPerFileDiff], "\n")
	truncated += fmt.Sprintf(
		"\n[%s 差异较大，共 %d 行，已截断至前 %d 行]\n",
		filePath, len(lines), maxLinesPerFileDiff,
	)
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
