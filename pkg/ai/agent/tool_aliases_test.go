package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToolAliases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		found    bool
	}{
		{"add", "stage_all", true},
		{"git_add", "stage_all", true},
		{"unstage", "unstage_all", true},
		{"switch", "checkout", true},
		{"branch", "create_branch", true},
		{"stage_all", "", false}, // 不在别名表中，应该直接使用
		{"commit", "", false},    // 不在别名表中，应该直接使用
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			alias, found := toolAliases[tt.input]
			assert.Equal(t, tt.found, found)
			if found {
				assert.Equal(t, tt.expected, alias)
			}
		})
	}
}
