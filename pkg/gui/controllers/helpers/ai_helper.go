package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/ai"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

type AIHelper struct {
	c             *HelperCommon
	loadingHelper *LoadingHelper
}

func NewAIHelper(c *HelperCommon, loadingHelper *LoadingHelper) *AIHelper {
	return &AIHelper{c: c, loadingHelper: loadingHelper}
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
