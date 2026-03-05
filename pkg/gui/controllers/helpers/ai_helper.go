package helpers

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/ai"
	aiprovider "github.com/dswcpp/lazygit/pkg/ai/provider"
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
			Label: self.c.Tr.TestCurrentProfile,
			OnPress: func() error {
				return self.testCurrentProfile()
			},
			Key: 'x',
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

// saveAndReloadAI persists config to disk and re-initialises the AI manager.
func (self *AIHelper) saveAndReloadAI() error {
	if err := self.c.GetConfig().SaveUserConfig(); err != nil {
		return err
	}

	newManager, err := ai.NewManager(self.c.UserConfig().AI, nil)
	if err != nil {
		return err
	}
	if newManager != nil {
		newManager.SetContextBuilder(NewGuiContextBuilder(self.c))
		RegisterGitTools(self.c, newManager)
	}
	self.c.AIManager = newManager

	self.c.Toast(self.c.Tr.AISettingsSaved)
	return nil
}

// OpenAIAssistant opens an interactive prompt where the user describes a git
// task. The AI generates the shell/git commands needed, shows them for
// confirmation, then executes them via a subprocess.
func (self *AIHelper) OpenAIAssistant() error {
	if self.c.AIManager == nil {
		return self.ShowFirstTimeWizard()
	}

	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.AIAssistantPrompt,
		HandleConfirm: func(userQuery string) error {
			if strings.TrimSpace(userQuery) == "" {
				return nil
			}
			self.loadingHelper.WithCenteredLoadingStatus(self.c.Tr.AIAssistantStatus, func(_ gocui.Task) error {
				repoCtx := self.c.AIManager.RepoContext().CompactString()
				prompt := fmt.Sprintf(
					self.c.Tr.AIAssistantSystemPrompt+
						self.c.Tr.AIAssistantRules+
						self.c.Tr.AIAssistantRepoState+
						self.c.Tr.AIAssistantUserRequest,
					repoCtx,
					userQuery,
				)

				result, err := self.c.AIManager.Provider().Complete(context.Background(), []aiprovider.Message{
					{Role: aiprovider.RoleUser, Content: prompt},
				})
				if err != nil {
					return self.HandleAIError(err)
				}

				response := strings.TrimSpace(result.Content)
				if response == "" {
					return errors.New(self.c.Tr.AIEmptyResponse)
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
	normalizedCommands := normalizeAICommands(commands)
	if len(normalizedCommands) == 0 {
		return errors.New(self.c.Tr.AIAssistantNoCommands)
	}
	if err := validateAICommands(normalizedCommands); err != nil {
		return err
	}

	preview := strings.Join(normalizedCommands, "\n")
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AIAssistantTitle,
		Prompt: self.c.Tr.AIAssistantConfirmExecute + "\n\n" + preview,
		HandleConfirm: func() error {
			cmdStr := buildSequentialCommandScript(normalizedCommands, self.c.OS().Platform.OS)
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
	return ExtractCommandsFromMessage(response)
}

// ExtractCommandsFromMessage 从 AI 消息文本中提取可执行命令。
// 优先提取代码块（```...```）中的内容，其次提取以 git / $ 开头的行。
func ExtractCommandsFromMessage(message string) []string {
	var cmds []string

	// 1. 提取代码块内容
	lines := strings.Split(message, "\n")
	inBlock := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inBlock = !inBlock
			continue
		}
		if inBlock {
			cmd := normalizeAICommandLine(trimmed)
			if cmd != "" && !strings.HasPrefix(cmd, "#") {
				cmds = append(cmds, cmd)
			}
		}
	}

	// 2. 若代码块中无内容，则从正文中提取明显的命令行
	if len(cmds) == 0 {
		for _, line := range lines {
			cmd := normalizeAICommandLine(line)
			if looksLikeCommand(cmd) {
				cmds = append(cmds, cmd)
			}
		}
	}

	return cmds
}

var commandListPrefixRe = regexp.MustCompile(`^\s*(?:[-*•]+\s+|\d+[.)]\s+)`)

// knownCommands 是允许从 AI 正文（非代码块）中提取的可执行程序白名单。
// 代码块内的命令不受此限制，由调用方直接提取。
var knownCommands = map[string]bool{
	// Git 工具
	"git": true, "gh": true, "hub": true,
	// Shell 内置 / 常用工具
	"cd": true, "ls": true, "mkdir": true, "rm": true, "cp": true, "mv": true,
	"touch": true, "cat": true, "echo": true, "grep": true, "find": true,
	"chmod": true, "chown": true, "ln": true, "curl": true, "wget": true,
	"ssh": true, "scp": true, "rsync": true,
	// Shell
	"sh": true, "bash": true, "zsh": true, "fish": true,
	// Node.js
	"npm": true, "yarn": true, "pnpm": true, "npx": true, "node": true,
	// Python
	"python": true, "python3": true, "pip": true, "pip3": true,
	// Go
	"go": true,
	// 构建工具
	"make": true, "cmake": true, "cargo": true, "mvn": true, "gradle": true,
}

func normalizeAICommandLine(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}

	// 去掉列表前缀（如 "1. "、"- "、"• "）
	line = commandListPrefixRe.ReplaceAllString(line, "")
	line = strings.TrimSpace(line)

	// 去掉常见 shell 提示符前缀
	line = strings.TrimPrefix(line, "$ ")
	line = strings.TrimPrefix(line, "> ")
	line = strings.TrimSpace(line)

	// 去掉包裹式单行反引号（精确去除首尾各一个，避免误删多重反引号）
	if strings.HasPrefix(line, "`") && strings.HasSuffix(line, "`") && len(line) > 1 {
		line = strings.TrimSpace(line[1 : len(line)-1])
	}

	return line
}

func looksLikeCommand(line string) bool {
	if line == "" {
		return false
	}

	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false
	}

	// 只允许白名单中的已知可执行程序，防止说明性文字被误判为命令
	return knownCommands[fields[0]]
}

// ConfirmAndSilentExecute 展示待执行命令，用户确认后在后台静默执行（不弹出终端）并刷新界面。
func (self *AIHelper) ConfirmAndSilentExecute(commands []string) error {
	normalizedCommands := normalizeAICommands(commands)
	if len(normalizedCommands) == 0 {
		return errors.New(self.c.Tr.AIAssistantSilentNoCommands)
	}
	if err := validateAICommands(normalizedCommands); err != nil {
		return err
	}

	preview := strings.Join(normalizedCommands, "\n")
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AIAssistantTitle,
		Prompt: self.c.Tr.AIAssistantConfirmSilentExecute + "\n\n" + preview,
		HandleConfirm: func() error {
			self.c.LogAction("AI silent execute")
			return self.c.WithWaitingStatus(self.c.Tr.AIAssistantExecuting, func(_ gocui.Task) error {
				cmdStr := buildSequentialCommandScript(normalizedCommands, self.c.OS().Platform.OS)
				if err := self.c.OS().Cmd.NewShell(cmdStr, self.c.UserConfig().OS.ShellFunctionsFile).Run(); err != nil {
					return fmt.Errorf("%s: %w", self.c.Tr.AIAssistantExecuteError, err)
				}
				self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				return nil
			})
		},
	})
	return nil
}

func normalizeAICommands(commands []string) []string {
	normalizedCommands := make([]string, 0, len(commands))
	for _, cmd := range commands {
		normalized := normalizeAICommandLine(cmd)
		if looksLikeCommand(normalized) {
			normalizedCommands = append(normalizedCommands, normalized)
		}
	}
	return normalizedCommands
}

func validateAICommands(commands []string) error {
	for i, cmd := range commands {
		if hasUnquotedGitCommitMessage(cmd) {
			return fmt.Errorf("第 %d 条命令可能无效：git commit -m 的提交信息包含空格时必须加引号", i+1)
		}
	}
	return nil
}

func hasUnquotedGitCommitMessage(command string) bool {
	trimmed := strings.TrimSpace(command)
	if !strings.HasPrefix(trimmed, "git commit") {
		return false
	}

	idx := strings.Index(trimmed, " -m ")
	if idx == -1 {
		return false
	}

	msg := strings.TrimSpace(trimmed[idx+4:])
	if msg == "" {
		return false
	}

	first := msg[0]
	if first == '"' || first == '\'' {
		return false
	}

	return len(strings.Fields(msg)) > 1
}

func buildSequentialCommandScript(commands []string, osName string) string {
	if len(commands) == 0 {
		return ""
	}

	if osName == "windows" {
		lines := []string{"@echo off", "setlocal"}
		for _, cmd := range commands {
			lines = append(lines, cmd)
			lines = append(lines, "if errorlevel 1 exit /b 1")
		}
		return strings.Join(lines, "\n")
	}

	lines := []string{"set -e"}
	lines = append(lines, commands...)
	return strings.Join(lines, "\n")
}

// ShowFirstTimeWizard guides new users through AI setup when they first try to use AI features.
// This improves the onboarding experience by making configuration more discoverable.
func (self *AIHelper) ShowFirstTimeWizard() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.AIWelcomeWizardTitle,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.UseDeepSeekRecommended,
				OnPress: func() error {
					return self.setupProvider("deepseek", "deepseek-reasoner", "https://api.deepseek.com/v1")
				},
				Key: 'd',
			},
			{
				Label: self.c.Tr.UseOpenAI,
				OnPress: func() error {
					return self.setupProvider("openai", "gpt-4o-mini", "https://api.openai.com/v1")
				},
				Key: 'o',
			},
			{
				Label: self.c.Tr.UseAnthropicClaude,
				OnPress: func() error {
					return self.setupProvider("anthropic", "claude-sonnet-4-6", "https://api.anthropic.com/v1")
				},
				Key: 'a',
			},
			{
				Label: self.c.Tr.UseOllamaLocal,
				OnPress: func() error {
					return self.setupProvider("ollama", "llama3", "http://localhost:11434/v1")
				},
				Key: 'l',
			},
			{
				Label: self.c.Tr.ConfigureLater,
				OnPress: func() error {
					return self.OpenAISettingsMenu()
				},
				Key: 's',
			},
			{
				Label: self.c.Tr.Cancel,
				OnPress: func() error {
					return nil
				},
				Key: 'c',
			},
		},
	})
}

// setupProvider creates a new AI profile with the specified provider and prompts for API key.
func (self *AIHelper) setupProvider(provider, defaultModel, defaultEndpoint string) error {
	// Prompt for API key
	self.c.Prompt(types.PromptOpts{
		Title: fmt.Sprintf(self.c.Tr.SetupProviderAPIKey, provider),
		FindSuggestionsFunc: func(currentText string) []*types.Suggestion {
			// Suggest environment variable references
			return []*types.Suggestion{
				{Value: "${DEEPSEEK_API_KEY}", Label: "Use environment variable: DEEPSEEK_API_KEY"},
				{Value: "${OPENAI_API_KEY}", Label: "Use environment variable: OPENAI_API_KEY"},
				{Value: "${ANTHROPIC_API_KEY}", Label: "Use environment variable: ANTHROPIC_API_KEY"},
			}
		},
		HandleConfirm: func(apiKey string) error {
			apiKey = strings.TrimSpace(apiKey)
			if apiKey == "" {
				return errors.New(self.c.Tr.APIKeyCannotBeEmpty)
			}

			// Create new profile
			profileName := provider + "-default"
			newProfile := config.AIProfileConfig{
				Name:           profileName,
				Provider:       provider,
				APIKey:         apiKey,
				Model:          defaultModel,
				Endpoint:       defaultEndpoint,
				MaxTokens:      8000,
				Timeout:        300,
				EnableThinking: provider == "deepseek",
			}

			// Add to config
			cfg := &self.c.UserConfig().AI
			cfg.Enabled = true
			cfg.Profiles = append(cfg.Profiles, newProfile)
			cfg.ActiveProfile = profileName

			// Save and reload
			if err := self.saveAndReloadAI(); err != nil {
				return err
			}

			// Offer to test the connection
			self.c.Confirm(types.ConfirmOpts{
				Title:  self.c.Tr.AIConfigComplete,
				Prompt: fmt.Sprintf(self.c.Tr.AIConfigCompletePrompt, profileName),
				HandleConfirm: func() error {
					return self.testCurrentProfile()
				},
			})

			return nil
		},
	})
	return nil
}

// testCurrentProfile tests the current AI profile by sending a simple completion request.
// This helps users verify their API configuration is correct before using AI features.
func (self *AIHelper) testCurrentProfile() error {
	if self.c.AIManager == nil {
		return errors.New(self.c.Tr.AINotEnabledPleaseConfig)
	}

	self.loadingHelper.WithCenteredLoadingStatus(self.c.Tr.AITestingConnection, func(_ gocui.Task) error {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		testPrompt := "Please reply 'OK' to confirm the connection is working."
		result, err := self.c.AIManager.Provider().Complete(ctx, []aiprovider.Message{
			{Role: aiprovider.RoleUser, Content: testPrompt},
		})

		if err != nil {
			// Use friendly error handling
			friendlyErr := self.HandleAIError(err)
			self.c.OnUIThread(func() error {
				self.c.Alert(self.c.Tr.AIConnectionTestFailed, friendlyErr.Error())
				return nil
			})
			return friendlyErr
		}

		// Check if we got a response
		response := strings.TrimSpace(result.Content)
		if response == "" {
			err := errors.New(self.c.Tr.AIEmptyResponse)
			self.c.OnUIThread(func() error {
				self.c.Alert(self.c.Tr.AIConnectionTestFailed, err.Error())
				return nil
			})
			return err
		}

		// Success
		self.c.OnUIThread(func() error {
			profile := self.c.UserConfig().AI.GetActiveProfile()
			profileInfo := "Unknown"
			if profile != nil {
				profileInfo = fmt.Sprintf("%s / %s", profile.Provider, profile.Model)
			}

			self.c.Toast(fmt.Sprintf(self.c.Tr.AIConnectionTestSuccessDetail, profileInfo, response))
			return nil
		})
		return nil
	})

	return nil
}

// HandleAIError converts raw AI errors into user-friendly Chinese messages.
// This provides better UX by translating common API errors (auth, rate limit, timeout)
// into actionable guidance for the user.
func (self *AIHelper) HandleAIError(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Timeout errors
	if strings.Contains(errStr, "context deadline exceeded") ||
		strings.Contains(errStr, "timeout") {
		return errors.New(self.c.Tr.AIRequestTimeout)
	}

	// Authentication errors
	if strings.Contains(errStr, "401") ||
		strings.Contains(errStr, "unauthorized") ||
		strings.Contains(errStr, "API key") ||
		strings.Contains(errStr, "Invalid API key") {
		return errors.New(self.c.Tr.APIKeyInvalid)
	}

	// Rate limiting
	if strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "rate limit") ||
		strings.Contains(errStr, "too many requests") {
		return errors.New(self.c.Tr.APIRateLimitExceeded)
	}

	// Network errors
	if strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "network") {
		return errors.New(self.c.Tr.NetworkConnectionFailed)
	}

	// Model not found / invalid model
	if strings.Contains(errStr, "model not found") ||
		strings.Contains(errStr, "invalid model") ||
		strings.Contains(errStr, "404") {
		return errors.New(self.c.Tr.ModelNotAvailable)
	}

	// Quota exceeded
	if strings.Contains(errStr, "quota") ||
		strings.Contains(errStr, "insufficient_quota") ||
		strings.Contains(errStr, "balance") {
		return errors.New(self.c.Tr.APIQuotaExhausted)
	}

	// Context length exceeded
	if strings.Contains(errStr, "context length") ||
		strings.Contains(errStr, "maximum context") ||
		strings.Contains(errStr, "token limit") {
		return errors.New(self.c.Tr.InputTooLong)
	}

	// Generic error with AI prefix for clarity
	return fmt.Errorf(self.c.Tr.AIGenericError, err)
}

// SuggestBranchName uses AI to suggest a branch name based on working tree changes.
// Returns a suggested branch name in kebab-case format (e.g., "feature/add-user-auth").
func (self *AIHelper) SuggestBranchName() (string, error) {
	if self.c.AIManager == nil {
		return "", errors.New(self.c.Tr.AINotEnabledConfigFirst)
	}

	if len(self.c.Model().Files) == 0 {
		return "", errors.New(self.c.Tr.NoChangesForBranchName)
	}

	rawDiff, err := self.c.Git().Diff.GetDiff(false)
	if err != nil {
		rawDiff = ""
	}
	const maxDiffChars = 8000
	if len(rawDiff) > maxDiffChars {
		rawDiff = rawDiff[:maxDiffChars]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err := self.c.AIManager.RunSkill(ctx, "branch_name", map[string]any{
		"diff": rawDiff,
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", errors.New(self.c.Tr.AIBranchNameCancelled)
		}
		return "", self.HandleAIError(err)
	}

	branchName := strings.TrimSpace(output.Content)
	branchName = strings.Trim(branchName, "\"'`")
	branchName = strings.ReplaceAll(branchName, " ", "-")
	if !strings.Contains(branchName, "/") {
		branchName = "feature/" + branchName
	}
	branchName = strings.ToLower(branchName)
	for _, char := range []string{"~", "^", ":", "?", "*", "[", "\\", "..", "@{", "//"} {
		branchName = strings.ReplaceAll(branchName, char, "")
	}
	return branchName, nil
}

// GeneratePRDescription uses AI to generate a pull request description based on commits and diff.
// Returns a formatted PR description suitable for GitHub/GitLab/etc.
func (self *AIHelper) GeneratePRDescription(fromBranch string, toBranch string) (string, error) {
	if self.c.AIManager == nil {
		return "", errors.New(self.c.Tr.AINotEnabledConfigFirst)
	}

	if len(self.c.Model().Commits) == 0 {
		return "", errors.New(self.c.Tr.NoCommitsForPRDescription)
	}

	baseBranchRef := toBranch
	if baseBranchRef == "" {
		baseBranchRef = "origin/main"
	}

	rawDiff, err := self.c.Git().Diff.GetDiff(false, baseBranchRef+"...HEAD")
	if err != nil {
		rawDiff = ""
	}
	const maxDiffChars = 15000
	if len(rawDiff) > maxDiffChars {
		rawDiff = rawDiff[:maxDiffChars]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	output, err := self.c.AIManager.RunSkill(ctx, "pr_desc", map[string]any{
		"diff":        rawDiff,
		"from_branch": fromBranch,
		"to_branch":   toBranch,
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", errors.New(self.c.Tr.AIPRDescriptionCancelled)
		}
		return "", self.HandleAIError(err)
	}

	description := strings.TrimSpace(output.Content)
	if description == "" {
		return "", errors.New(self.c.Tr.AIEmptyResponse)
	}
	return description, nil
}
