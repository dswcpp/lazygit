package helpers

import (
	"github.com/dswcpp/lazygit/pkg/ai"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

type AIHelper struct {
	c *HelperCommon
}

func NewAIHelper(c *HelperCommon) *AIHelper {
	return &AIHelper{c: c}
}

// OpenAISettingsMenu opens a menu to configure AI settings in-app.
// Changes take effect immediately and are persisted to the config file.
func (self *AIHelper) OpenAISettingsMenu() error {
	cfg := self.c.UserConfig().AI

	toggleLabel := self.c.Tr.AISettingsEnable
	if cfg.Enabled {
		toggleLabel = self.c.Tr.AISettingsDisable
	}

	items := []*types.MenuItem{
		{
			Label: toggleLabel,
			OnPress: func() error {
				self.c.UserConfig().AI.Enabled = !self.c.UserConfig().AI.Enabled
				return self.saveAndReloadAI()
			},
			Key: 'e',
		},
		{
			Label: self.c.Tr.AISettingsSetAPIKey,
			OnPress: func() error {
				return self.promptAndSave(
					self.c.Tr.AISettingsAPIKeyPrompt,
					self.c.UserConfig().AI.APIKey,
					func(val string) { self.c.UserConfig().AI.APIKey = val },
				)
			},
			Key: 'k',
		},
		{
			Label: self.c.Tr.AISettingsSetProvider,
			OnPress: func() error {
				return self.openProviderMenu()
			},
			Key:       'p',
			OpensMenu: true,
		},
		{
			Label: self.c.Tr.AISettingsSetModel,
			OnPress: func() error {
				return self.promptAndSave(
					self.c.Tr.AISettingsModelPrompt,
					self.c.UserConfig().AI.Model,
					func(val string) { self.c.UserConfig().AI.Model = val },
				)
			},
			Key: 'm',
		},
		{
			Label: self.c.Tr.AISettingsSetEndpoint,
			OnPress: func() error {
				return self.promptAndSave(
					self.c.Tr.AISettingsEndpointPrompt,
					self.c.UserConfig().AI.Endpoint,
					func(val string) { self.c.UserConfig().AI.Endpoint = val },
				)
			},
			Key: 'u',
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.AISettings,
		Items: items,
	})
}

func (self *AIHelper) openProviderMenu() error {
	providers := []struct {
		name string
		key  rune
	}{
		{"deepseek", 'd'},
		{"openai", 'o'},
		{"ollama", 'l'},
		{"custom", 'c'},
	}

	items := make([]*types.MenuItem, len(providers))
	for i, p := range providers {
		provider := p
		items[i] = &types.MenuItem{
			Label: provider.name,
			OnPress: func() error {
				self.c.UserConfig().AI.Provider = provider.name
				return self.saveAndReloadAI()
			},
			Key:    provider.key,
			Widget: types.MakeMenuRadioButton(self.c.UserConfig().AI.Provider == provider.name),
		}
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.AISettingsSetProvider,
		Items: items,
	})
}

func (self *AIHelper) promptAndSave(prompt, initialValue string, apply func(string)) error {
	self.c.Prompt(types.PromptOpts{
		Title:          prompt,
		InitialContent: initialValue,
		HandleConfirm: func(val string) error {
			apply(val)
			return self.saveAndReloadAI()
		},
	})
	return nil
}

// saveAndReloadAI persists config to disk and re-initialises the AI client.
func (self *AIHelper) saveAndReloadAI() error {
	if err := self.c.GetConfig().SaveUserConfig(); err != nil {
		return err
	}

	newClient, err := ai.NewClient(self.c.UserConfig().AI)
	if err != nil {
		return err
	}
	self.c.AI = newClient
	self.c.Toast(self.c.Tr.AISettingsSaved)
	return nil
}
