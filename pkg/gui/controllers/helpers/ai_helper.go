package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/ai"
	"github.com/dswcpp/lazygit/pkg/config"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

type AIHelper struct {
	c             *HelperCommon
	loadingHelper *LoadingHelper
}

func NewAIHelper(c *HelperCommon, loadingHelper *LoadingHelper) *AIHelper {
	return &AIHelper{c: c, loadingHelper: loadingHelper}
}

// OpenAISettingsMenu opens the top-level AI settings menu.
// Changes take effect immediately and are persisted to the config file.
func (self *AIHelper) OpenAISettingsMenu() error {
	cfg := self.c.UserConfig().AI

	toggleLabel := self.c.Tr.AISettingsEnable
	if cfg.Enabled {
		toggleLabel = self.c.Tr.AISettingsDisable
	}

	activeProfile := cfg.GetActiveProfile()
	activeProfileName := ""
	if activeProfile != nil {
		activeProfileName = activeProfile.Name
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
			Label:   fmt.Sprintf("%s: %s", self.c.Tr.AISettingsActiveProfile, activeProfileName),
			Tooltip: self.c.Tr.AISettingsSwitchProfile,
			OnPress: func() error {
				return self.openSwitchProfileMenu()
			},
			Key:       'p',
			OpensMenu: true,
		},
		{
			Label: self.c.Tr.AISettingsEditProfile,
			OnPress: func() error {
				return self.openEditActiveProfileMenu()
			},
			Key:       'f',
			OpensMenu: true,
		},
		{
			Label: self.c.Tr.AISettingsAddProfile,
			OnPress: func() error {
				return self.openAddProfileMenu()
			},
			Key: 'a',
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.AISettings,
		Items: items,
	})
}

// openSwitchProfileMenu shows a radio-button list of all profiles for switching.
func (self *AIHelper) openSwitchProfileMenu() error {
	cfg := self.c.UserConfig().AI
	if len(cfg.Profiles) == 0 {
		return errors.New(self.c.Tr.AISettingsNoProfiles)
	}

	items := make([]*types.MenuItem, len(cfg.Profiles))
	for i, p := range cfg.Profiles {
		profile := p
		isActive := profile.Name == cfg.ActiveProfile
		items[i] = &types.MenuItem{
			Label: fmt.Sprintf("%s  (%s / %s)", profile.Name, profile.Provider, profile.Model),
			OnPress: func() error {
				self.c.UserConfig().AI.ActiveProfile = profile.Name
				return self.saveAndReloadAI()
			},
			Widget: types.MakeMenuRadioButton(isActive),
		}
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.AISettingsSwitchProfile,
		Items: items,
	})
}

// openEditActiveProfileMenu opens the edit sub-menu for the currently active profile.
func (self *AIHelper) openEditActiveProfileMenu() error {
	cfg := self.c.UserConfig().AI
	idx := -1
	for i, p := range cfg.Profiles {
		if p.Name == cfg.ActiveProfile {
			idx = i
			break
		}
	}
	if idx == -1 && len(cfg.Profiles) > 0 {
		idx = 0
	}
	if idx == -1 {
		return errors.New(self.c.Tr.AISettingsNoProfiles)
	}
	return self.openEditProfileMenu(idx)
}

// openEditProfileMenu opens a sub-menu to edit a single profile at the given index.
func (self *AIHelper) openEditProfileMenu(idx int) error {
	profile := &self.c.UserConfig().AI.Profiles[idx]

	items := []*types.MenuItem{
		{
			Label: fmt.Sprintf("%s: %s", self.c.Tr.AISettingsProfileName, profile.Name),
			OnPress: func() error {
				return self.promptField(self.c.Tr.AISettingsProfileNamePrompt, profile.Name,
					func(val string) { self.c.UserConfig().AI.Profiles[idx].Name = val })
			},
			Key: 'n',
		},
		{
			Label: fmt.Sprintf("%s: %s", self.c.Tr.AISettingsSetProvider, profile.Provider),
			OnPress: func() error {
				return self.openProviderMenuForProfile(idx)
			},
			Key:       'p',
			OpensMenu: true,
		},
		{
			Label: fmt.Sprintf("%s: %s", self.c.Tr.AISettingsSetAPIKey, maskKey(profile.APIKey)),
			OnPress: func() error {
				return self.promptField(self.c.Tr.AISettingsAPIKeyPrompt, profile.APIKey,
					func(val string) { self.c.UserConfig().AI.Profiles[idx].APIKey = val })
			},
			Key: 'k',
		},
		{
			Label: fmt.Sprintf("%s: %s", self.c.Tr.AISettingsSetModel, profile.Model),
			OnPress: func() error {
				return self.promptField(self.c.Tr.AISettingsModelPrompt, profile.Model,
					func(val string) { self.c.UserConfig().AI.Profiles[idx].Model = val })
			},
			Key: 'm',
		},
		{
			Label: fmt.Sprintf("%s: %s", self.c.Tr.AISettingsSetEndpoint, profile.Endpoint),
			OnPress: func() error {
				return self.promptField(self.c.Tr.AISettingsEndpointPrompt, profile.Endpoint,
					func(val string) { self.c.UserConfig().AI.Profiles[idx].Endpoint = val })
			},
			Key: 'u',
		},
		{
			Label: fmt.Sprintf("%s: %d", self.c.Tr.AISettingsMaxTokens, profile.MaxTokens),
			OnPress: func() error {
				current := ""
				if profile.MaxTokens > 0 {
					current = fmt.Sprintf("%d", profile.MaxTokens)
				}
				return self.promptField(self.c.Tr.AISettingsMaxTokensPrompt, current,
					func(val string) {
						n := 0
						fmt.Sscanf(val, "%d", &n)
						self.c.UserConfig().AI.Profiles[idx].MaxTokens = n
					})
			},
			Key: 't',
		},
		{
			Label: fmt.Sprintf("%s: %d", self.c.Tr.AISettingsTimeout, profile.Timeout),
			OnPress: func() error {
				current := ""
				if profile.Timeout > 0 {
					current = fmt.Sprintf("%d", profile.Timeout)
				}
				return self.promptField(self.c.Tr.AISettingsTimeoutPrompt, current,
					func(val string) {
						n := 0
						fmt.Sscanf(val, "%d", &n)
						self.c.UserConfig().AI.Profiles[idx].Timeout = n
					})
			},
			Key: 'o',
		},
		{
			Label: self.c.Tr.AISettingsDeleteProfile,
			OnPress: func() error {
				return self.deleteProfile(idx)
			},
			Key: 'd',
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: fmt.Sprintf("%s: %s", self.c.Tr.AISettingsEditProfile, profile.Name),
		Items: items,
	})
}

// openProviderMenuForProfile shows a radio-button provider picker for a given profile index.
func (self *AIHelper) openProviderMenuForProfile(idx int) error {
	providers := []struct {
		name string
		key  rune
	}{
		{"deepseek", 'd'},
		{"openai", 'o'},
		{"ollama", 'l'},
		{"anthropic", 'a'},
		{"custom", 'c'},
	}

	items := make([]*types.MenuItem, len(providers))
	for i, p := range providers {
		prov := p
		isSelected := self.c.UserConfig().AI.Profiles[idx].Provider == prov.name
		items[i] = &types.MenuItem{
			Label: prov.name,
			OnPress: func() error {
				self.c.UserConfig().AI.Profiles[idx].Provider = prov.name
				return self.saveAndReloadAI()
			},
			Key:    prov.key,
			Widget: types.MakeMenuRadioButton(isSelected),
		}
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.AISettingsSetProvider,
		Items: items,
	})
}

// openAddProfileMenu prompts for a name and creates a new profile with defaults.
func (self *AIHelper) openAddProfileMenu() error {
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.AISettingsNewProfileNamePrompt,
		HandleConfirm: func(name string) error {
			name = strings.TrimSpace(name)
			if name == "" {
				return nil
			}
			newProfile := config.AIProfileConfig{
				Name:      name,
				Provider:  "deepseek",
				Model:     "deepseek-chat",
				MaxTokens: 500,
				Timeout:   60,
			}
			self.c.UserConfig().AI.Profiles = append(self.c.UserConfig().AI.Profiles, newProfile)
			self.c.UserConfig().AI.ActiveProfile = name
			return self.saveAndReloadAI()
		},
	})
	return nil
}

// deleteProfile removes the profile at idx after confirmation.
func (self *AIHelper) deleteProfile(idx int) error {
	profiles := self.c.UserConfig().AI.Profiles
	if len(profiles) <= 1 {
		return errors.New(self.c.Tr.AISettingsCannotDeleteLastProfile)
	}
	profileName := profiles[idx].Name
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AISettingsDeleteProfileTitle,
		Prompt: fmt.Sprintf(self.c.Tr.AISettingsDeleteProfilePrompt, profileName),
		HandleConfirm: func() error {
			cfg := &self.c.UserConfig().AI
			cfg.Profiles = append(cfg.Profiles[:idx], cfg.Profiles[idx+1:]...)
			// If we deleted the active profile, switch to first available
			if cfg.ActiveProfile == profileName {
				cfg.ActiveProfile = cfg.Profiles[0].Name
			}
			return self.saveAndReloadAI()
		},
	})
	return nil
}

// promptField opens a prompt to edit a single profile field and saves on confirm.
func (self *AIHelper) promptField(prompt, initialValue string, apply func(string)) error {
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

// maskKey returns a masked version of an API key for display (shows last 4 chars).
func maskKey(key string) string {
	if len(key) <= 4 {
		return strings.Repeat("*", len(key))
	}
	return strings.Repeat("*", len(key)-4) + key[len(key)-4:]
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

// OpenAIAssistant opens an interactive prompt where the user describes a git
// task. The AI generates the shell/git commands needed, shows them for
// confirmation, then executes them via a subprocess.
func (self *AIHelper) OpenAIAssistant() error {
	if self.c.AI == nil {
		return errors.New(self.c.Tr.AINotEnabled)
	}

	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.AIAssistantPrompt,
		HandleConfirm: func(userQuery string) error {
			if strings.TrimSpace(userQuery) == "" {
				return nil
			}
			self.loadingHelper.WithCenteredLoadingStatus(self.c.Tr.AIAssistantStatus, func(_ gocui.Task) error {
				repoCtx := self.buildGitContext()
				prompt := fmt.Sprintf(
					"你是一个 git 命令生成器。根据用户需求和仓库状态，生成需要执行的 shell/git 命令。\n\n"+
						"规则：\n"+
						"- 只输出可直接执行的命令，每行一条\n"+
						"- 不输出任何解释、注释（#开头）或 markdown\n"+
						"- 命令按执行顺序排列\n"+
						"- 如果需求无法用 git 命令安全完成，第一行输出：CANNOT_EXECUTE: <原因>\n\n"+
						"当前仓库状态：\n%s\n"+
						"用户需求：%s",
					repoCtx,
					userQuery,
				)

				result, err := self.c.AI.Complete(context.Background(), prompt)
				if err != nil {
					return err
				}

				response := strings.TrimSpace(result.Content)
				if response == "" {
					return errors.New("AI: empty response from model")
				}

				if strings.HasPrefix(response, "CANNOT_EXECUTE:") {
					reason := strings.TrimSpace(strings.TrimPrefix(response, "CANNOT_EXECUTE:"))
					self.c.OnUIThread(func() error {
						self.c.Alert(self.c.Tr.AIAssistantTitle, reason)
						return nil
					})
					return nil
				}

				commands := parseAICommands(response)
				if len(commands) == 0 {
					return errors.New(self.c.Tr.AIAssistantNoCommands)
				}

				self.c.OnUIThread(func() error {
					return self.confirmAndExecuteCommands(commands)
				})
				return nil
			})
			return nil
		},
	})
	return nil
}

// confirmAndExecuteCommands shows the AI-generated commands to the user and
// executes them via a subprocess on confirmation.
// Must be called on the UI thread.
func (self *AIHelper) confirmAndExecuteCommands(commands []string) error {
	preview := strings.Join(commands, "\n")
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AIAssistantTitle,
		Prompt: self.c.Tr.AIAssistantConfirmExecute + "\n\n" + preview,
		HandleConfirm: func() error {
			cmdStr := strings.Join(commands, " && ")
			self.c.LogAction("AI git assistant")
			return self.c.RunSubprocessAndRefresh(
				self.c.OS().Cmd.NewShell(cmdStr, self.c.UserConfig().OS.ShellFunctionsFile),
			)
		},
	})
	return nil
}

// parseAICommands splits the AI response into individual shell commands,
// discarding blank lines and comment lines.
func parseAICommands(response string) []string {
	var cmds []string
	for _, line := range strings.Split(response, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		cmds = append(cmds, line)
	}
	return cmds
}

// buildGitContext collects current repository information to include in the AI prompt.
func (self *AIHelper) buildGitContext() string {
	var sb strings.Builder

	sb.WriteString("当前分支: " + self.c.Model().CheckedOutBranch + "\n")

	// Recent commits (up to 10)
	commits := self.c.Model().Commits
	limit := 10
	if len(commits) < limit {
		limit = len(commits)
	}
	if limit > 0 {
		sb.WriteString("\n最近提交:\n")
		for _, commit := range commits[:limit] {
			sb.WriteString(fmt.Sprintf("  %s %s\n", commit.ShortHash(), commit.Name))
		}
	}

	// Working tree files
	files := self.c.Model().Files
	if len(files) > 0 {
		sb.WriteString("\n变更文件:\n")
		for _, f := range files {
			sb.WriteString(fmt.Sprintf("  %s %s\n", f.ShortStatus, f.Path))
		}
	}

	return sb.String()
}
