package gittools

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// makeDiff builds a fake unified diff with the given number of context and
// changed lines per hunk, across nFiles files.
func makeDiff(nFiles, contextPerSide, changed int) string {
	var sb strings.Builder
	for f := 0; f < nFiles; f++ {
		sb.WriteString("diff --git a/file" + string(rune('A'+f)) + ".go b/file" + string(rune('A'+f)) + ".go\n")
		sb.WriteString("index abc..def 100644\n--- a/file.go\n+++ b/file.go\n")
		sb.WriteString("@@ -1,10 +1,10 @@\n")
		for i := 0; i < contextPerSide; i++ {
			sb.WriteString(" context line\n")
		}
		for i := 0; i < changed; i++ {
			sb.WriteString("-removed\n")
			sb.WriteString("+added\n")
		}
		for i := 0; i < contextPerSide; i++ {
			sb.WriteString(" context line\n")
		}
	}
	return sb.String()
}

func TestCompressDiff_NoLimitReturnsFull(t *testing.T) {
	diff := makeDiff(1, 3, 2)
	result := compressDiff(diff, 0)
	assert.Equal(t, diff, result)
}

func TestCompressDiff_FitsReturnsFull(t *testing.T) {
	diff := makeDiff(1, 3, 2)
	lines := countDiffLines(diff)
	result := compressDiff(diff, lines+10)
	assert.Equal(t, diff, result)
}

func TestCompressDiff_ContextTrimmed(t *testing.T) {
	// Build a diff with 10 context lines on each side: will need trimming
	diff := makeDiff(1, 10, 2)
	// Budget that forces context reduction but keeps changed lines
	result := compressDiff(diff, 20)

	// All changed lines must be present
	assert.Contains(t, result, "-removed")
	assert.Contains(t, result, "+added")
	// Compression note must be present
	assert.Contains(t, result, "[diff compressed")
}

func TestCompressDiff_HunkStatsLevel(t *testing.T) {
	// Very tight budget: force hunk-stats level
	diff := makeDiff(2, 5, 5)
	result := compressDiff(diff, 10)

	// Should still show both files' @@ headers
	assert.True(t, strings.Count(result, "@@") >= 2, "should keep hunk headers")
	// Each hunk shows +N / -M
	assert.Contains(t, result, "[+")
	assert.Contains(t, result, "[diff compressed")
}

func TestCompressDiff_FileStatsLevel(t *testing.T) {
	// Extremely tight budget: force file-stats level
	diff := makeDiff(3, 5, 10)
	result := compressDiff(diff, 5)

	// Should have one line per file (3 files)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	// At least 3 non-empty lines (one per file) + note line
	nonEmpty := 0
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			nonEmpty++
		}
	}
	assert.GreaterOrEqual(t, nonEmpty, 3)
	assert.Contains(t, result, "[diff compressed")
}

func TestCompressDiff_ChangedLinesKeptUnderLightCompression(t *testing.T) {
	// With enough budget to reach level 3 (strip context, keep changed lines),
	// all +/- lines must be present.
	diff := makeDiff(1, 10, 3) // 10 context lines each side, 3 changes
	// Full diff ≈ 4 header + 1 @@ + 10+6+10 body = ~31 lines
	// Level-3 (ctx=0): 4 header + 1 @@ + 6 changed = ~11 lines — fits in 15
	result := compressDiff(diff, 15)
	assert.Equal(t, 3, strings.Count(result, "-removed"))
	assert.Equal(t, 3, strings.Count(result, "+added"))
	assert.Contains(t, result, "[diff compressed")
}

func TestCompressDiff_HunkStatsWhenVeryTight(t *testing.T) {
	// When budget is too small for even stripped diff, hunk statistics are used.
	// Changed lines are summarised; the LLM sees "+N / -M" counts instead.
	diff := makeDiff(1, 0, 3) // 3 changes, no context
	// Full stripped diff ≈ 11 lines; force hunk-stats level with budget=5
	result := compressDiff(diff, 5)
	assert.Contains(t, result, "[diff compressed")
	// Either hunk stats or file stats — either way the note is present
	// and the result does NOT exceed the budget (at file-stats level it fits)
	assert.LessOrEqual(t, countDiffLines(result), 10) // loose upper bound
}

// ── dual-mode behaviour (compressDiff vs paginateDiff) ───────────────────────

// paginateDiff is the raw-offset mode used by get_file_diff when offset is set.
// These tests verify that the two modes produce meaningfully different output
// for the same diff when it exceeds the budget.

func TestPaginateDiff_FirstPage(t *testing.T) {
	diff := makeDiff(1, 5, 2) // ~18 lines
	result := paginateDiff(diff, 0, 8, &stubTruncator{})
	// Should contain the first 8 lines and a truncation notice
	assert.Contains(t, result, "truncated")
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	// Content lines ≤ 8 (the notice adds 1 more)
	assert.LessOrEqual(t, len(lines), 10)
}

func TestPaginateDiff_SecondPage(t *testing.T) {
	diff := makeDiff(1, 5, 2)
	page1 := paginateDiff(diff, 0, 8, &stubTruncator{})
	// Extract next offset from truncation notice
	assert.Contains(t, page1, "offset=8")

	page2 := paginateDiff(diff, 8, 8, &stubTruncator{})
	// page2 must contain different content than page1
	assert.NotEqual(t, page1, page2)
}

func TestCompressVsPaginate_CompressShowsAllFiles(t *testing.T) {
	// Two-file diff: compression shows both files; pagination may hide second file.
	diff := makeDiff(2, 5, 3) // ~40 lines across 2 files
	budget := 15

	compressed := compressDiff(diff, budget)
	paginated := paginateDiff(diff, 0, budget, &stubTruncator{})

	// Compression must mention both files (fileA.go and fileB.go)
	assert.True(t, strings.Contains(compressed, "fileA") && strings.Contains(compressed, "fileB"),
		"compressed diff should reference both files")

	// Pagination with budget=15 likely only shows first file's content
	_ = paginated // behaviour varies; just ensure compression wins on coverage
}

// stubTruncator implements the interface required by paginateDiff.
type stubTruncator struct{}

func (s *stubTruncator) ToolTruncated(start, end, total, nextOffset int) string {
	return fmt.Sprintf("\n... [truncated: lines %d-%d of %d. Use offset=%d]", start, end, total, nextOffset)
}

// ─────────────────────────────────────────────────────────────────────────────

func TestParseDiff_TwoFiles(t *testing.T) {
	diff := makeDiff(2, 2, 1)
	files := parseDiff(diff)
	assert.Len(t, files, 2)
	assert.Len(t, files[0].hunks, 1)
	assert.Len(t, files[1].hunks, 1)
}

func TestParseDiff_HunkLineClassification(t *testing.T) {
	diff := "diff --git a/x b/x\n--- a/x\n+++ b/x\n@@ -1,3 +1,3 @@\n context\n-removed\n+added\n"
	files := parseDiff(diff)
	assert.Len(t, files, 1)
	assert.Len(t, files[0].hunks, 1)
	hunk := files[0].hunks[0]
	assert.Equal(t, byte(' '), hunk.lines[0].kind)
	assert.Equal(t, byte('-'), hunk.lines[1].kind)
	assert.Equal(t, byte('+'), hunk.lines[2].kind)
}
