package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractCommandsFromMessage(t *testing.T) {
	message := `下面给你 4 条命令：
1. git fetch origin
2. git checkout main
3. git pull --rebase origin main
4. git status`

	cmds := ExtractCommandsFromMessage(message)

	assert.Equal(t, []string{
		"git fetch origin",
		"git checkout main",
		"git pull --rebase origin main",
		"git status",
	}, cmds)
}

func TestParseAICommands_IgnoresNarrativeText(t *testing.T) {
	response := `我建议按下面顺序执行：

1. git add .
2. git commit -m "chore: update"
3. git push origin HEAD

执行完再告诉我结果。`

	cmds := parseAICommands(response)

	assert.Equal(t, []string{
		"git add .",
		`git commit -m "chore: update"`,
		"git push origin HEAD",
	}, cmds)
}

func TestBuildSequentialCommandScript_Windows(t *testing.T) {
	script := buildSequentialCommandScript([]string{
		"cd repo",
		"git status",
	}, "windows")

	assert.Equal(t, "@echo off\nsetlocal\ncd repo\nif errorlevel 1 exit /b 1\ngit status\nif errorlevel 1 exit /b 1", script)
}

func TestBuildSequentialCommandScript_Posix(t *testing.T) {
	script := buildSequentialCommandScript([]string{
		"cd repo",
		"git status",
	}, "linux")

	assert.Equal(t, "set -e\ncd repo\ngit status", script)
}

func TestHasUnquotedGitCommitMessage(t *testing.T) {
	assert.True(t, hasUnquotedGitCommitMessage("git commit -m Auto commit before merge"))
	assert.False(t, hasUnquotedGitCommitMessage(`git commit -m "Auto commit before merge"`))
	assert.False(t, hasUnquotedGitCommitMessage("git status"))
}

func TestValidateAICommands(t *testing.T) {
	err := validateAICommands([]string{
		"git add -A",
		"git commit -m Auto commit before merge",
		"git push",
	})
	assert.Error(t, err)

	err = validateAICommands([]string{
		"git add -A",
		`git commit -m "Auto commit before merge"`,
		"git push",
	})
	assert.NoError(t, err)
}
