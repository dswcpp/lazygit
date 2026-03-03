package helpers

import (
	"strings"
	"testing"

	"github.com/dswcpp/lazygit/pkg/i18n"
)

// mockTr creates a mock translation set for testing
func mockTr() *i18n.TranslationSet {
	return &i18n.TranslationSet{
		AIDiffSkipped:             "[Skipped %s: %s]",
		AIDiffBinaryFile:          "binary file",
		AIDiffLockOrGeneratedFile: "lock/generated file",
		AIDiffChangeStats:         "# Change Statistics",
		AIDiffFilesCount:          "- Files: %d",
		AIDiffFileTypes:           "- File types: ",
		AIDiffChangeScale:         "- Change scale: +%d/-%d lines",
		AIDiffMajorChanges:        "- Major changes: ",
		AIDiffNewFile:             "[New]",
		AIDiffDeletedFile:         "[Deleted]",
		AIDiffRenamedFile:         "[Renamed]",
		AIDiffModifiedFile:        "[Modified]",
		AIDiffTruncated:           "[%s diff is large, %d lines total, truncated to first %d lines]",
		AIDiffSmartTruncated:      "[%s smart truncated: preserved %d/%d lines of key code (function signatures, important comments, etc.)]",
	}
}

// TestCalculateFilePriority tests the file priority calculation.
func TestCalculateFilePriority(t *testing.T) {
	tests := []struct {
		path     string
		expected int
	}{
		// Core business logic - highest priority
		{"pkg/ai/ai.go", 100},
		{"src/service/user.go", 100},
		{"lib/utils/helper.go", 100},
		{"internal/core/engine.go", 100},

		// API/controllers - high priority
		{"api/handler.go", 90},
		{"controller/user_controller.go", 90},
		{"handler/auth_handler.go", 90},
		{"service/auth_service.go", 90},

		// Models - medium-high priority
		{"model/user.go", 80},
		{"schema/user_schema.go", 80},
		{"entity/product.go", 80},

		// Config files - medium priority
		{"config.yaml", 70},
		{"settings.json", 70},
		{".env", 70},

		// Test files - low priority
		{"user_test.go", 50},
		{"handler.test.ts", 50},
		{"test/integration_test.go", 50},

		// Documentation - lowest priority
		{"README.md", 30},
		{"docs/api.md", 30},
		{"CHANGELOG.txt", 30},

		// Default priority
		{"random.file", 60},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := calculateFilePriority(tt.path)
			if result != tt.expected {
				t.Errorf("calculateFilePriority(%q) = %d, want %d", tt.path, result, tt.expected)
			}
		})
	}
}

// TestGenerateSemanticSummary tests semantic summary generation.
func TestGenerateSemanticSummary(t *testing.T) {
	tr := mockTr()
	tests := []struct {
		name       string
		block      string
		path       string
		insertions int
		deletions  int
		expected   string
	}{
		{
			name:       "new file",
			block:      "diff --git a/new.go b/new.go\nnew file mode 100644",
			path:       "new.go",
			insertions: 10,
			deletions:  0,
			expected:   "### [New] new.go",
		},
		{
			name:       "deleted file",
			block:      "diff --git a/old.go b/old.go\ndeleted file mode 100644",
			path:       "old.go",
			insertions: 0,
			deletions:  10,
			expected:   "### [Deleted] old.go",
		},
		{
			name:       "renamed file",
			block:      "diff --git a/old.go b/new.go\nrename from old.go\nrename to new.go",
			path:       "new.go",
			insertions: 0,
			deletions:  0,
			expected:   "### [Renamed] new.go <- old.go",
		},
		{
			name:       "modified file",
			block:      "diff --git a/file.go b/file.go\n+added\n-removed",
			path:       "file.go",
			insertions: 50,
			deletions:  20,
			expected:   "### [Modified] file.go (+50/-20)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateSemanticSummary(tt.block, tt.path, tt.insertions, tt.deletions, tr)
			if result != tt.expected {
				t.Errorf("generateSemanticSummary() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestSmartTruncateBlock tests smart truncation.
func TestSmartTruncateBlock(t *testing.T) {
	tr := mockTr()
	// Create a large block with function signatures
	lines := []string{
		"diff --git a/file.go b/file.go",
		"--- a/file.go",
		"+++ b/file.go",
		"@@ -1,10 +1,10 @@",
		"package main",
	}

	// Add many regular lines
	for i := 0; i < 300; i++ {
		lines = append(lines, "+regular line")
	}

	// Add important lines
	lines = append(lines, "+func ImportantFunction() {")
	lines = append(lines, "+// TODO: implement this")
	lines = append(lines, "+}")

	block := strings.Join(lines, "\n")

	result := smartTruncateBlock(block, "file.go", tr)

	// Should be truncated
	if !strings.Contains(result, "smart truncated") {
		t.Error("Expected truncation message")
	}

	// Should preserve function signature
	if !strings.Contains(result, "func ImportantFunction") {
		t.Error("Expected to preserve function signature")
	}

	// Should preserve TODO comment
	if !strings.Contains(result, "TODO") {
		t.Error("Expected to preserve TODO comment")
	}
}

// TestPrioritizeFileBlocks tests file block prioritization.
func TestPrioritizeFileBlocks(t *testing.T) {
	blocks := []string{
		"diff --git a/README.md b/README.md\n+doc change",
		"diff --git a/pkg/core/engine.go b/pkg/core/engine.go\n+core change",
		"diff --git a/test/unit_test.go b/test/unit_test.go\n+test change",
		"diff --git a/api/handler.go b/api/handler.go\n+api change",
	}

	prioritized := prioritizeFileBlocks(blocks)

	// Check order: core > api > test > docs
	if len(prioritized) != 4 {
		t.Fatalf("Expected 4 blocks, got %d", len(prioritized))
	}

	// Core should be first
	if !strings.Contains(prioritized[0].path, "engine.go") {
		t.Errorf("Expected core file first, got %s", prioritized[0].path)
	}

	// API should be second
	if !strings.Contains(prioritized[1].path, "handler.go") {
		t.Errorf("Expected API file second, got %s", prioritized[1].path)
	}

	// Test should be third
	if !strings.Contains(prioritized[2].path, "test") {
		t.Errorf("Expected test file third, got %s", prioritized[2].path)
	}

	// Docs should be last
	if !strings.Contains(prioritized[3].path, "README") {
		t.Errorf("Expected docs file last, got %s", prioritized[3].path)
	}
}

// TestCountDiffChanges tests change counting.
func TestCountDiffChanges(t *testing.T) {
	block := `diff --git a/file.go b/file.go
--- a/file.go
+++ b/file.go
@@ -1,5 +1,7 @@
 unchanged line
+added line 1
+added line 2
-removed line 1
-removed line 2
-removed line 3
 unchanged line`

	insertions, deletions := countDiffChanges(block)

	if insertions != 2 {
		t.Errorf("Expected 2 insertions, got %d", insertions)
	}

	if deletions != 3 {
		t.Errorf("Expected 3 deletions, got %d", deletions)
	}
}

// TestGetFileExtension tests file extension extraction.
func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"file.go", ".go"},
		{"path/to/file.ts", ".ts"},
		{"config.yaml", ".yaml"},
		{"README.md", ".md"},
		{".gitignore", "other"},
		{"no-extension", "other"},
		{"path/to/.hidden", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := getFileExtension(tt.path)
			if result != tt.expected {
				t.Errorf("getFileExtension(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

// TestIsNewFile tests new file detection.
func TestIsNewFile(t *testing.T) {
	tests := []struct {
		block    string
		expected bool
	}{
		{"diff --git a/new.go b/new.go\nnew file mode 100644", true},
		{"diff --git a/file.go b/file.go\n+change", false},
	}

	for _, tt := range tests {
		result := isNewFile(tt.block)
		if result != tt.expected {
			t.Errorf("isNewFile() = %v, want %v", result, tt.expected)
		}
	}
}

// TestIsDeletedFile tests deleted file detection.
func TestIsDeletedFile(t *testing.T) {
	tests := []struct {
		block    string
		expected bool
	}{
		{"diff --git a/old.go b/old.go\ndeleted file mode 100644", true},
		{"diff --git a/file.go b/file.go\n+change", false},
	}

	for _, tt := range tests {
		result := isDeletedFile(tt.block)
		if result != tt.expected {
			t.Errorf("isDeletedFile() = %v, want %v", result, tt.expected)
		}
	}
}

// TestIsRenamed tests renamed file detection.
func TestIsRenamed(t *testing.T) {
	tests := []struct {
		block    string
		expected bool
	}{
		{"diff --git a/old.go b/new.go\nrename from old.go\nrename to new.go", true},
		{"diff --git a/file.go b/file.go\n+change", false},
	}

	for _, tt := range tests {
		result := isRenamed(tt.block)
		if result != tt.expected {
			t.Errorf("isRenamed() = %v, want %v", result, tt.expected)
		}
	}
}
