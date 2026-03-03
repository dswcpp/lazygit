package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/ai"
	"github.com/dswcpp/lazygit/pkg/commands/models"
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
			Label: "测试当前 Profile 连接",
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
		// Show first-time wizard instead of error
		return self.ShowFirstTimeWizard()
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
					return self.HandleAIError(err)
				}

				response := strings.TrimSpace(result.Content)
				if response == "" {
					return errors.New("AI 返回了空响应。请稍后重试或检查 AI 设置（Ctrl+A）")
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

// buildGitContext collects comprehensive repository information to include in the AI prompt.
// This enhanced version provides much richer context for better AI command generation.
func (self *AIHelper) buildGitContext() string {
	var sb strings.Builder

	// Current branch with tracking information
	var currentBranch *models.Branch
	if len(self.c.Model().Branches) > 0 {
		currentBranch = self.c.Model().Branches[0]
	}
	if currentBranch != nil {
		sb.WriteString(fmt.Sprintf("当前分支: %s\n", currentBranch.Name))

		// Branch tracking status (ahead/behind remote)
		if currentBranch.IsTrackingRemote() {
			// AheadForPull and BehindForPull are strings, check if they're not empty
			ahead := currentBranch.AheadForPull
			behind := currentBranch.BehindForPull

			if ahead != "" || behind != "" {
				if ahead != "" && behind != "" {
					sb.WriteString(fmt.Sprintf("  跟踪远程分支: %s（本地领先 %s 个提交，落后 %s 个提交）\n", currentBranch.UpstreamRemote, ahead, behind))
				} else if ahead != "" {
					sb.WriteString(fmt.Sprintf("  跟踪远程分支: %s（本地领先 %s 个提交，未推送）\n", currentBranch.UpstreamRemote, ahead))
				} else if behind != "" {
					sb.WriteString(fmt.Sprintf("  跟踪远程分支: %s（本地落后 %s 个提交，需要拉取）\n", currentBranch.UpstreamRemote, behind))
				}
			} else {
				sb.WriteString(fmt.Sprintf("  跟踪远程分支: %s（已同步）\n", currentBranch.UpstreamRemote))
			}
		} else {
			sb.WriteString("  未跟踪远程分支（本地分支）\n")
		}
	} else {
		sb.WriteString("当前分支: " + self.c.Model().CheckedOutBranch + "\n")
	}

	// Working tree state (merge/rebase/cherry-pick in progress)
	workingTreeState := self.c.Git().Status.WorkingTreeState()
	if workingTreeState.Any() {
		sb.WriteString(fmt.Sprintf("\n工作树状态: %s\n", workingTreeState.Title(self.c.Tr)))
	}

	// Staged and unstaged changes count
	stagedCount := 0
	unstagedCount := 0
	untrackedCount := 0
	files := self.c.Model().Files
	for _, f := range files {
		if f.HasStagedChanges {
			stagedCount++
		}
		if f.HasUnstagedChanges {
			unstagedCount++
		}
		if !f.Tracked {
			untrackedCount++
		}
	}
	if stagedCount > 0 || unstagedCount > 0 || untrackedCount > 0 {
		sb.WriteString(fmt.Sprintf("\n变更统计: 已暂存 %d，未暂存 %d，未跟踪 %d\n", stagedCount, unstagedCount, untrackedCount))
	}

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

	// Working tree files (detailed, limited to first 20)
	if len(files) > 0 {
		sb.WriteString("\n变更文件:\n")
		displayLimit := 20
		if len(files) < displayLimit {
			displayLimit = len(files)
		}
		for i := 0; i < displayLimit; i++ {
			f := files[i]
			sb.WriteString(fmt.Sprintf("  %s %s\n", f.ShortStatus, f.Path))
		}
		if len(files) > displayLimit {
			sb.WriteString(fmt.Sprintf("  ...（还有 %d 个文件）\n", len(files)-displayLimit))
		}
	}

	// Stash list (if any)
	stashEntries := self.c.Model().StashEntries
	if len(stashEntries) > 0 {
		sb.WriteString(fmt.Sprintf("\n贮藏列表: %d 个贮藏\n", len(stashEntries)))
		stashLimit := 5
		if len(stashEntries) < stashLimit {
			stashLimit = len(stashEntries)
		}
		for i := 0; i < stashLimit; i++ {
			sb.WriteString(fmt.Sprintf("  %s\n", stashEntries[i].Name))
		}
		if len(stashEntries) > stashLimit {
			sb.WriteString(fmt.Sprintf("  ...（还有 %d 个贮藏）\n", len(stashEntries)-stashLimit))
		}
	}

	return sb.String()
}

// ShowFirstTimeWizard guides new users through AI setup when they first try to use AI features.
// This improves the onboarding experience by making configuration more discoverable.
func (self *AIHelper) ShowFirstTimeWizard() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: "欢迎使用 AI 功能 - 首次配置向导",
		Items: []*types.MenuItem{
			{
				Label: "使用 DeepSeek（推荐）",
				OnPress: func() error {
					return self.setupProvider("deepseek", "deepseek-reasoner", "https://api.deepseek.com/v1")
				},
				Key: 'd',
			},
			{
				Label: "使用 OpenAI",
				OnPress: func() error {
					return self.setupProvider("openai", "gpt-4o-mini", "https://api.openai.com/v1")
				},
				Key: 'o',
			},
			{
				Label: "使用 Anthropic Claude",
				OnPress: func() error {
					return self.setupProvider("anthropic", "claude-sonnet-4-6", "https://api.anthropic.com/v1")
				},
				Key: 'a',
			},
			{
				Label: "使用 Ollama（本地模型）",
				OnPress: func() error {
					return self.setupProvider("ollama", "llama3", "http://localhost:11434/v1")
				},
				Key: 'l',
			},
			{
				Label: "稍后配置（进入 AI 设置）",
				OnPress: func() error {
					return self.OpenAISettingsMenu()
				},
				Key: 's',
			},
			{
				Label: "取消",
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
		Title: fmt.Sprintf("设置 %s API Key", provider),
		FindSuggestionsFunc: func(currentText string) []*types.Suggestion {
			// Suggest environment variable references
			return []*types.Suggestion{
				{Value: "${DEEPSEEK_API_KEY}", Label: "使用环境变量: DEEPSEEK_API_KEY"},
				{Value: "${OPENAI_API_KEY}", Label: "使用环境变量: OPENAI_API_KEY"},
				{Value: "${ANTHROPIC_API_KEY}", Label: "使用环境变量: ANTHROPIC_API_KEY"},
			}
		},
		HandleConfirm: func(apiKey string) error {
			apiKey = strings.TrimSpace(apiKey)
			if apiKey == "" {
				return errors.New("API Key 不能为空")
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
				Title:  "AI 配置完成",
				Prompt: fmt.Sprintf("✓ %s Profile 已创建并激活！\n\n是否测试连接？", profileName),
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
	if self.c.AI == nil {
		return errors.New("AI 未启用。请先启用 AI 并配置 Profile")
	}

	self.loadingHelper.WithCenteredLoadingStatus("正在测试 AI 连接...", func(_ gocui.Task) error {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Send a simple test prompt
		testPrompt := "请回复 'OK' 确认连接正常。"
		result, err := self.c.AI.Complete(ctx, testPrompt)

		if err != nil {
			// Use friendly error handling
			friendlyErr := self.HandleAIError(err)
			self.c.OnUIThread(func() error {
				self.c.Alert("AI 连接测试失败", friendlyErr.Error())
				return nil
			})
			return friendlyErr
		}

		// Check if we got a response
		response := strings.TrimSpace(result.Content)
		if response == "" {
			err := errors.New("AI 返回了空响应")
			self.c.OnUIThread(func() error {
				self.c.Alert("AI 连接测试失败", err.Error())
				return nil
			})
			return err
		}

		// Success
		self.c.OnUIThread(func() error {
			profile := self.c.UserConfig().AI.GetActiveProfile()
			profileInfo := "未知"
			if profile != nil {
				profileInfo = fmt.Sprintf("%s / %s", profile.Provider, profile.Model)
			}

			self.c.Toast(fmt.Sprintf("✓ AI 连接测试成功！\nProfile: %s\n响应: %s", profileInfo, response))
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
		return errors.New("AI 请求超时。请稍后重试或在 AI 设置中调整超时时间（Ctrl+A → 编辑 Profile → Timeout）")
	}

	// Authentication errors
	if strings.Contains(errStr, "401") ||
		strings.Contains(errStr, "unauthorized") ||
		strings.Contains(errStr, "API key") ||
		strings.Contains(errStr, "Invalid API key") {
		return errors.New("API 密钥无效。请检查 AI 设置中的密钥是否正确（Ctrl+A → 编辑 Profile → API Key）")
	}

	// Rate limiting
	if strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "rate limit") ||
		strings.Contains(errStr, "too many requests") {
		return errors.New("API 请求频率超限。请稍后重试或考虑更换提供商（Ctrl+A → 切换 Profile）")
	}

	// Network errors
	if strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "network") {
		return errors.New("网络连接失败。请检查网络连接或 AI 设置中的 Endpoint 配置（Ctrl+A → 编辑 Profile → Endpoint）")
	}

	// Model not found / invalid model
	if strings.Contains(errStr, "model not found") ||
		strings.Contains(errStr, "invalid model") ||
		strings.Contains(errStr, "404") {
		return errors.New("模型不可用。请在 AI 设置中选择其他模型（Ctrl+A → 编辑 Profile → Model）")
	}

	// Quota exceeded
	if strings.Contains(errStr, "quota") ||
		strings.Contains(errStr, "insufficient_quota") ||
		strings.Contains(errStr, "balance") {
		return errors.New("API 配额已用尽。请检查账户余额或更换提供商（Ctrl+A → 切换 Profile）")
	}

	// Context length exceeded
	if strings.Contains(errStr, "context length") ||
		strings.Contains(errStr, "maximum context") ||
		strings.Contains(errStr, "token limit") {
		return errors.New("输入内容过长超出模型限制。请减少暂存的文件数量或在 AI 设置中增加 MaxTokens（Ctrl+A → 编辑 Profile → MaxTokens）")
	}

	// Generic error with AI prefix for clarity
	return fmt.Errorf("AI 请求失败: %v\n\n提示：按 Ctrl+A 检查 AI 设置或切换提供商", err)
}

// SuggestBranchName uses AI to suggest a branch name based on working tree changes.
// Returns a suggested branch name in kebab-case format (e.g., "feature/add-user-auth").
func (self *AIHelper) SuggestBranchName() (string, error) {
	if self.c.AI == nil {
		return "", errors.New("AI 功能未启用。按 Ctrl+A 配置 AI 设置")
	}

	// Analyze working tree changes
	files := self.c.Model().Files
	if len(files) == 0 {
		return "", errors.New("工作区没有变更。请先进行一些修改后再使用 AI 建议分支名称")
	}

	// Build summary of changes
	var changesSummary strings.Builder
	stagedFiles := []string{}
	unstagedFiles := []string{}

	for _, f := range files {
		if f.HasStagedChanges {
			stagedFiles = append(stagedFiles, f.Path)
		} else if f.HasUnstagedChanges {
			unstagedFiles = append(unstagedFiles, f.Path)
		}
	}

	changesSummary.WriteString("变更的文件：\n")
	if len(stagedFiles) > 0 {
		changesSummary.WriteString("已暂存：\n")
		for i, file := range stagedFiles {
			if i >= 15 { // Limit to first 15 files
				changesSummary.WriteString(fmt.Sprintf("  ... 还有 %d 个文件\n", len(stagedFiles)-15))
				break
			}
			changesSummary.WriteString(fmt.Sprintf("  - %s\n", file))
		}
	}
	if len(unstagedFiles) > 0 {
		changesSummary.WriteString("未暂存：\n")
		for i, file := range unstagedFiles {
			if i >= 15 { // Limit to first 15 files
				changesSummary.WriteString(fmt.Sprintf("  ... 还有 %d 个文件\n", len(unstagedFiles)-15))
				break
			}
			changesSummary.WriteString(fmt.Sprintf("  - %s\n", file))
		}
	}

	// Get diff for more context (limit size to avoid token overflow)
	rawDiff, err := self.c.Git().Diff.GetDiff(false) // All changes, not just staged
	if err != nil {
		rawDiff = "" // Ignore diff errors, use file list only
	}

	// Truncate diff if too large
	const maxDiffChars = 8000
	diff := rawDiff
	if len(diff) > maxDiffChars {
		diff = diff[:maxDiffChars] + "\n[diff 已截断，仅显示前 8000 字符]"
	}

	// Build prompt for AI
	prompt := fmt.Sprintf(
		"你是一个 git 分支命名助手。\n\n"+
			"任务：根据以下变更内容，建议一个规范的 git 分支名称。\n\n"+
			"分支命名规则：\n"+
			"- 格式：<类型>/<简短描述>\n"+
			"- 类型：feature（新功能）、bugfix（修复 bug）、refactor（重构）、docs（文档）、test（测试）、chore（杂项）\n"+
			"- 描述：使用 kebab-case（小写字母和连字符），简洁明了，不超过 50 字符\n"+
			"- 示例：feature/user-authentication、bugfix/login-crash、refactor/api-client\n\n"+
			"变更内容：\n%s\n"+
			"diff 摘要：\n```diff\n%s\n```\n\n"+
			"要求：\n"+
			"- 只输出一个分支名称（不要包含任何解释或引号）\n"+
			"- 必须符合上述格式和规则\n"+
			"- 使用英文命名\n",
		changesSummary.String(),
		diff,
	)

	// Call AI with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := self.c.AI.Complete(ctx, prompt)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", errors.New("AI 分支命名建议已取消")
		}
		return "", self.HandleAIError(err)
	}

	// Clean up the response
	branchName := strings.TrimSpace(result.Content)
	branchName = strings.Trim(branchName, "\"'`")       // Remove quotes
	branchName = strings.ReplaceAll(branchName, " ", "-") // Replace spaces with hyphens

	// Validate format (should be <type>/<description>)
	if !strings.Contains(branchName, "/") {
		// If AI didn't follow format, prepend "feature/"
		branchName = "feature/" + branchName
	}

	// Ensure lowercase and valid characters
	branchName = strings.ToLower(branchName)
	// Remove invalid characters (git branch names can't have certain chars)
	invalidChars := []string{"~", "^", ":", "?", "*", "[", "\\", "..", "@{", "//"}
	for _, char := range invalidChars {
		branchName = strings.ReplaceAll(branchName, char, "")
	}

	return branchName, nil
}

// GeneratePRDescription uses AI to generate a pull request description based on commits and diff.
// Returns a formatted PR description suitable for GitHub/GitLab/etc.
func (self *AIHelper) GeneratePRDescription(fromBranch string, toBranch string) (string, error) {
	if self.c.AI == nil {
		return "", errors.New("AI 功能未启用。按 Ctrl+A 配置 AI 设置")
	}

	// Get commits in the current branch (commits ahead of base branch)
	currentBranch := self.c.Model().Branches[0]
	commits := self.c.Model().Commits

	if len(commits) == 0 {
		return "", errors.New("当前分支没有提交可用于生成 PR 描述")
	}

	// Build commit history summary (limit to recent commits to avoid token overflow)
	var commitsSummary strings.Builder
	commitsSummary.WriteString("提交历史：\n")
	maxCommits := 20
	for i, commit := range commits {
		if i >= maxCommits {
			commitsSummary.WriteString(fmt.Sprintf("\n... 还有 %d 个提交\n", len(commits)-maxCommits))
			break
		}
		// Format: hash - message (author)
		commitsSummary.WriteString(fmt.Sprintf("- %s - %s (%s)\n",
			commit.Hash()[:8],
			commit.Name,
			commit.AuthorName,
		))
	}

	// Get diff from base branch to current HEAD
	// Use git diff to compare branches
	baseBranchRef := toBranch
	if baseBranchRef == "" {
		baseBranchRef = "origin/main" // Default to main branch
	}

	// Get diff between base and current branch (use three-dot notation for merge base)
	rawDiff, err := self.c.Git().Diff.GetDiff(false, baseBranchRef+"...HEAD")
	if err != nil {
		// If diff fails, try to get recent commits diff instead
		rawDiff = fmt.Sprintf("[无法获取完整 diff: %v]\n提交数量: %d\n", err, len(commits))
	}

	// Truncate diff if too large
	const maxDiffChars = 15000
	diff := rawDiff
	if len(diff) > maxDiffChars {
		diff = diff[:maxDiffChars] + "\n[diff 已截断，仅显示前 15000 字符]"
	}

	// Get branch tracking info for context
	branchContext := ""
	if currentBranch.IsTrackingRemote() {
		branchContext = fmt.Sprintf("源分支: %s/%s\n目标分支: %s\n",
			currentBranch.UpstreamRemote,
			currentBranch.UpstreamBranch,
			toBranch,
		)
	} else {
		branchContext = fmt.Sprintf("源分支: %s (本地)\n目标分支: %s\n",
			fromBranch,
			toBranch,
		)
	}

	// Build prompt for AI
	prompt := fmt.Sprintf(
		"你是一个 Pull Request 描述生成助手。\n\n"+
			"任务：根据以下提交历史和代码变更，生成一个专业的 PR 描述。\n\n"+
			"%s\n"+
			"%s\n"+
			"代码变更：\n```diff\n%s\n```\n\n"+
			"PR 描述格式要求：\n"+
			"## 摘要\n"+
			"一句话概括本次 PR 的主要目的和价值。\n\n"+
			"## 变更内容\n"+
			"- 列出主要的功能变更、bug 修复或重构（3-5 个要点）\n\n"+
			"## 技术细节（可选）\n"+
			"- 如有重要的技术实现细节或架构变更，简要说明\n\n"+
			"## 测试\n"+
			"- [ ] 单元测试已通过\n"+
			"- [ ] 手动测试已完成\n"+
			"- [ ] 代码审查已准备就绪\n\n"+
			"输出要求：\n"+
			"- 使用简体中文（代码和技术术语除外）\n"+
			"- 使用 Markdown 格式\n"+
			"- 简洁专业，突出重点\n"+
			"- 不要包含标题（# Pull Request），直接从摘要开始\n",
		branchContext,
		commitsSummary.String(),
		diff,
	)

	// Call AI with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	result, err := self.c.AI.Complete(ctx, prompt)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", errors.New("AI PR 描述生成已取消")
		}
		return "", self.HandleAIError(err)
	}

	// Clean up the response
	description := strings.TrimSpace(result.Content)

	if description == "" {
		return "", errors.New("AI 返回了空响应。请稍后重试或检查 AI 设置（Ctrl+A）")
	}

	return description, nil
}
