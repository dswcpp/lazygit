package config

import (
	"github.com/dswcpp/lazygit/pkg/config"
	. "github.com/dswcpp/lazygit/pkg/integration/components"
)

var RemoteNamedStar = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Having a config remote.*",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			SetConfig("remote.*.prune", "true").
			CreateNCommits(2)
	},
	SetupConfig: func(cfg *config.AppConfig) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// here we're just asserting that we haven't panicked upon starting lazygit
		t.Views().Commits().
			Lines(
				AnyString(),
				AnyString(),
			)
	},
})
