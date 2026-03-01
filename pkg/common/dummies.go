package common

import (
	"github.com/dswcpp/lazygit/pkg/config"
	"github.com/dswcpp/lazygit/pkg/i18n"
	"github.com/dswcpp/lazygit/pkg/utils"
	"github.com/spf13/afero"
)

func NewDummyCommon() *Common {
	tr := i18n.EnglishTranslationSet()
	cmn := &Common{
		Log: utils.NewDummyLog(),
		Tr:  tr,
		Fs:  afero.NewOsFs(),
	}
	cmn.SetUserConfig(config.GetDefaultConfig())
	return cmn
}

func NewDummyCommonWithUserConfigAndAppState(userConfig *config.UserConfig, appState *config.AppState) *Common {
	tr := i18n.EnglishTranslationSet()
	cmn := &Common{
		Log:      utils.NewDummyLog(),
		Tr:       tr,
		AppState: appState,
		// TODO: remove dependency on actual filesystem in tests and switch to using
		// in-memory for everything
		Fs: afero.NewOsFs(),
	}
	cmn.SetUserConfig(userConfig)
	return cmn
}
