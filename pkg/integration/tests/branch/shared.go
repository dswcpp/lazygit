package branch

import (
	"github.com/dswcpp/lazygit/pkg/config"
	. "github.com/dswcpp/lazygit/pkg/integration/components"
	"github.com/samber/lo"
)

func checkRemoteBranches(t *TestDriver, keys config.KeybindingConfig, remoteName string, expectedBranches []string) {
	t.Views().Remotes().
		Focus().
		NavigateToLine(Contains(remoteName)).
		PressEnter()

	t.Views().
		RemoteBranches().
		Lines(
			lo.Map(expectedBranches, func(branch string, _ int) *TextMatcher { return Equals(branch) })...,
		).
		Press(keys.Universal.Return)

	t.Views().
		Branches().
		Focus()
}
