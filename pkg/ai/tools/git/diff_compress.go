package gittools

import (
	"fmt"
	"strings"
)

// compressDiff applies AI-agent-style context compression to a unified diff.
//
// Instead of truncating at an arbitrary line boundary (which hides changes),
// it progressively reduces detail while preserving structural completeness:
// at every compression level the LLM still sees every changed file and every
// changed hunk — only the amount of surrounding context shrinks.
//
// Compression levels (tried in order until the result fits maxLines):
//
//	Level 0  – verbatim (diff already fits)
//	Level 1  – trim context to 3 lines per hunk
//	Level 2  – trim context to 1 line per hunk
//	Level 3  – strip all context lines (changed lines + @@ headers only)
//	Level 4  – replace each hunk body with "+N / -M" statistics
//	Level 5  – replace each file with a single-line summary (file stats only)
//
// maxLines <= 0 means no limit (returns diff verbatim).
func compressDiff(diff string, maxLines int) string {
	if maxLines <= 0 || countDiffLines(diff) <= maxLines {
		return diff
	}

	files := parseDiff(diff)

	// Levels 1–3: reduce context lines progressively
	for _, ctx := range []int{3, 1, 0} {
		result := renderDiff(files, renderModeContext, ctx)
		if countDiffLines(result) <= maxLines {
			return result + fmt.Sprintf(
				"\n\n[diff compressed – showing %d context lines per hunk; use get_file_diff for full details]",
				ctx,
			)
		}
	}

	// Level 4: hunk statistics only
	result := renderDiff(files, renderModeHunkStats, 0)
	if countDiffLines(result) <= maxLines {
		return result + "\n\n[diff compressed – hunk statistics only; use get_file_diff for full hunks]"
	}

	// Level 5: file statistics only
	return renderDiff(files, renderModeFileStats, 0) +
		"\n\n[diff compressed – file statistics only; use get_file_diff for per-file diff]"
}

// ── internal types ────────────────────────────────────────────────────────────

type parsedFile struct {
	header string // "diff --git …", index, ---, +++ lines
	hunks  []parsedHunk
}

type parsedHunk struct {
	header string      // "@@ -x,y +x,y @@ …" line
	lines  []parsedLine
}

type parsedLine struct {
	text string // full line content (includes leading ' ', '+', or '-')
	kind byte   // ' ' = context, '+' = added, '-' = removed
}

type renderMode int

const (
	renderModeContext   renderMode = iota // full hunk body, context trimmed
	renderModeHunkStats                   // @@ header + "+N -M" per hunk
	renderModeFileStats                   // one summary line per file
)

// ── parser ────────────────────────────────────────────────────────────────────

// parseDiff splits a unified diff into per-file structures.
func parseDiff(diff string) []parsedFile {
	var files []parsedFile
	var cur *parsedFile
	var curHunk *parsedHunk

	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "diff --git "):
			if cur != nil {
				if curHunk != nil {
					cur.hunks = append(cur.hunks, *curHunk)
					curHunk = nil
				}
				files = append(files, *cur)
			}
			cur = &parsedFile{header: line}
			curHunk = nil

		case cur != nil && strings.HasPrefix(line, "@@"):
			if curHunk != nil {
				cur.hunks = append(cur.hunks, *curHunk)
			}
			curHunk = &parsedHunk{header: line}

		case cur != nil && curHunk == nil:
			// Between "diff --git" and the first "@@" (index, ---, +++ lines)
			cur.header += "\n" + line

		case cur != nil && curHunk != nil:
			kind := byte(' ')
			if len(line) > 0 && (line[0] == '+' || line[0] == '-') {
				kind = line[0]
			}
			curHunk.lines = append(curHunk.lines, parsedLine{text: line, kind: kind})
		}
	}

	// Flush last file
	if cur != nil {
		if curHunk != nil {
			cur.hunks = append(cur.hunks, *curHunk)
		}
		files = append(files, *cur)
	}

	return files
}

// ── renderer ──────────────────────────────────────────────────────────────────

func renderDiff(files []parsedFile, mode renderMode, contextLines int) string {
	var sb strings.Builder
	for i, f := range files {
		if i > 0 {
			sb.WriteByte('\n')
		}
		switch mode {
		case renderModeFileStats:
			sb.WriteString(fileStatSummary(f))
		case renderModeHunkStats:
			sb.WriteString(f.header)
			for _, h := range f.hunks {
				sb.WriteByte('\n')
				sb.WriteString(hunkStatSummary(h))
			}
		default: // renderModeContext
			sb.WriteString(f.header)
			for _, h := range f.hunks {
				sb.WriteByte('\n')
				sb.WriteString(h.header)
				sb.WriteByte('\n')
				sb.WriteString(renderHunkContext(h, contextLines))
			}
		}
	}
	return sb.String()
}

// renderHunkContext returns the hunk body with context lines trimmed to at
// most maxCtx lines before and after each changed block.
func renderHunkContext(h parsedHunk, maxCtx int) string {
	// Mark each line as changed (+/-) or context
	changed := make([]bool, len(h.lines))
	for i, l := range h.lines {
		if l.kind == '+' || l.kind == '-' {
			changed[i] = true
		}
	}

	// Include lines within maxCtx distance of any changed line
	include := make([]bool, len(h.lines))
	for i := range h.lines {
		if !changed[i] {
			continue
		}
		include[i] = true
		for d := 1; d <= maxCtx; d++ {
			if i-d >= 0 {
				include[i-d] = true
			}
			if i+d < len(h.lines) {
				include[i+d] = true
			}
		}
	}

	var sb strings.Builder
	ellipsis := false
	for i, l := range h.lines {
		if include[i] {
			if ellipsis {
				sb.WriteString(" ...\n")
				ellipsis = false
			}
			sb.WriteString(l.text)
			sb.WriteByte('\n')
		} else {
			ellipsis = true
		}
	}
	return sb.String()
}

// hunkStatSummary returns a compact one-line representation of a hunk.
func hunkStatSummary(h parsedHunk) string {
	added, removed := 0, 0
	for _, l := range h.lines {
		switch l.kind {
		case '+':
			added++
		case '-':
			removed++
		}
	}
	return fmt.Sprintf("%s  [+%d / -%d lines]", h.header, added, removed)
}

// fileStatSummary returns a compact one-line representation of a file.
func fileStatSummary(f parsedFile) string {
	added, removed := 0, 0
	for _, h := range f.hunks {
		for _, l := range h.lines {
			switch l.kind {
			case '+':
				added++
			case '-':
				removed++
			}
		}
	}
	// Extract target path from "diff --git a/foo b/foo"
	firstLine := strings.SplitN(f.header, "\n", 2)[0]
	parts := strings.Fields(firstLine)
	path := ""
	if len(parts) >= 4 {
		path = strings.TrimPrefix(parts[3], "b/")
	}
	return fmt.Sprintf("%-60s  +%d / -%d", path, added, removed)
}

// ── util ──────────────────────────────────────────────────────────────────────

func countDiffLines(s string) int {
	return strings.Count(s, "\n") + 1
}
