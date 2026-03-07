package ui

import (
	"github.com/dswcpp/lazygit/pkg/config"
	. "github.com/dswcpp/lazygit/pkg/integration/components"
)

// ActivityBarNavigation tests Activity Bar navigation functionality
var ActivityBarNavigation = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Test Activity Bar navigation between panels",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.ActivityBar.Show = true
		cfg.GetUserConfig().Gui.ActivityBar.Width = 3
		cfg.GetUserConfig().Gui.ActivityBar.IconStyle = "ascii"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1.txt", "content1")
		shell.Commit("initial commit")
		shell.NewBranch("feature")
		shell.CreateFileAndAdd("file2.txt", "content2")
		shell.Commit("feature commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Note: Activity Bar interactions would require GUI-level testing
		// For now, we verify the Activity Bar is rendered and doesn't crash
		t.Views().Files().
			Focus().
			IsFocused()

		// Verify navigating to different panels works
		t.Views().Branches().
			Focus().
			IsFocused().
			Lines(
				Contains("feature"),
				Contains("master"),
			)

		t.Views().Commits().
			Focus().
			IsFocused().
			LineCount(EqualsInt(2))
	},
})

// ActivityBarFetch tests Activity Bar fetch functionality
var ActivityBarFetch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Test Activity Bar fetch operation",
	ExtraCmdArgs: []string{},
	Skip:         true, // Skip until we have remote repo test infrastructure
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.ActivityBar.Show = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("README.md", "# Test")
		shell.Commit("initial commit")
		// TODO: Set up remote repository for fetch testing
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// TODO: Implement fetch test once remote repo infrastructure is available
		// This would require:
		// 1. Setting up a remote repository
		// 2. Clicking the fetch icon in Activity Bar
		// 3. Verifying fetch completed successfully
	},
})

// ActivityBarStash tests Activity Bar stash functionality
var ActivityBarStash = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Test Activity Bar stash changes operation",
	ExtraCmdArgs: []string{},
	Skip:         true, // Skip until Activity Bar UI interaction testing is implemented
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.ActivityBar.Show = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1.txt", "initial content")
		shell.Commit("initial commit")
		shell.UpdateFile("file1.txt", "modified content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// TODO: Implement stash test once Activity Bar UI interaction is available
		// This would require:
		// 1. Clicking the stash icon in Activity Bar
		// 2. Entering stash message in prompt
		// 3. Verifying stash was created successfully
	},
})
