package controllers

import (
	"slices"
	"strings"

	"github.com/dswcpp/lazygit/pkg/gui/controllers/helpers"
	"github.com/dswcpp/lazygit/pkg/gui/types"
	"github.com/dswcpp/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type ShellCommandAction struct {
	c                *ControllerCommon
	enhancedHelper   *helpers.EnhancedShellCommandHelper
	aiMode           bool
}

func (self *ShellCommandAction) Call() error {
	// Initialize enhanced helper if needed
	if self.enhancedHelper == nil {
		self.enhancedHelper = helpers.NewEnhancedShellCommandHelper(
			self.c.HelperCommon,
			self.c.Helpers().AI,
		)
	}

	// Determine title based on mode
	title := self.c.Tr.ShellCommand
	if self.enhancedHelper.IsAIMode() {
		title = self.c.Tr.ShellCommandAIMode
	}

	self.c.Prompt(types.PromptOpts{
		Title:               title,
		FindSuggestionsFunc: self.GetEnhancedSuggestionsFunc(),
		AllowEditSuggestion: true,
		PreserveWhitespace:  true,
		HandleConfirm: func(command string) error {
			// Check for dangerous commands before execution
			return self.confirmAndExecuteCommand(command)
		},
		HandleDeleteSuggestion: func(index int) error {
			// index is the index in the _filtered_ list of suggestions, so we
			// need to map it back to the full list. There's no really good way
			// to do this, but fortunately we keep the items in the
			// ShellCommandsHistory unique, which allows us to simply search
			// for it by string.
			item := self.c.Contexts().Suggestions.GetItems()[index].Value
			fullIndex := lo.IndexOf(self.c.GetAppState().ShellCommandsHistory, item)
			if fullIndex == -1 {
				// Should never happen, but better be safe
				return nil
			}

			self.c.GetAppState().ShellCommandsHistory = slices.Delete(
				self.c.GetAppState().ShellCommandsHistory, fullIndex, fullIndex+1)
			self.c.SaveAppStateAndLogError()
			self.c.Contexts().Suggestions.RefreshSuggestions()
			return nil
		},
	})

	return nil
}

// GetEnhancedSuggestionsFunc returns enhanced suggestions combining AI, templates, completions, and history.
func (self *ShellCommandAction) GetEnhancedSuggestionsFunc() func(string) []*types.Suggestion {
	return func(input string) []*types.Suggestion {
		if self.enhancedHelper == nil {
			// Fallback to simple history suggestions
			return self.GetShellCommandsHistorySuggestionsFunc()(input)
		}

		// Get intelligent suggestions
		suggestions := self.enhancedHelper.GetSuggestions(input)

		// Convert to types.Suggestion format
		return self.enhancedHelper.ConvertToTypesSuggestions(suggestions)
	}
}

// GetShellCommandsHistorySuggestionsFunc returns simple history-based suggestions (fallback).
func (self *ShellCommandAction) GetShellCommandsHistorySuggestionsFunc() func(string) []*types.Suggestion {
	return func(input string) []*types.Suggestion {
		history := self.c.GetAppState().ShellCommandsHistory

		return helpers.FilterFunc(history, self.c.UserConfig().Gui.UseFuzzySearch())(input)
	}
}

// confirmAndExecuteCommand checks for dangerous commands and confirms before execution.
func (self *ShellCommandAction) confirmAndExecuteCommand(command string) error {
	command = strings.TrimSpace(command)
	if command == "" {
		return nil
	}

	// Save to history if appropriate
	if self.shouldSaveCommand(command) {
		self.c.GetAppState().ShellCommandsHistory = utils.Limit(
			lo.Uniq(append([]string{command}, self.c.GetAppState().ShellCommandsHistory...)),
			1000,
		)
		self.c.SaveAppStateAndLogError()
	}

	// Check for dangerous commands
	if self.enhancedHelper != nil {
		isSafe, riskReason := self.enhancedHelper.CheckCommandSafety(command)
		if !isSafe {
			// Show confirmation dialog for dangerous commands
			return self.confirmDangerousCommand(command, riskReason)
		}
	}

	// Execute command
	return self.executeCommand(command)
}

// confirmDangerousCommand shows a confirmation dialog for dangerous commands.
func (self *ShellCommandAction) confirmDangerousCommand(command string, riskReason string) error {
	// Get safer alternative if available
	alternative := ""
	if self.enhancedHelper != nil {
		alternative = self.enhancedHelper.GetSafterAlternative(command)
	}

	promptMessage := riskReason
	if alternative != "" {
		promptMessage += "\n\n" + alternative
	}
	promptMessage += "\n\n确定要执行此命令吗？"

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.ShellCommandDangerousWarning,
		Prompt: promptMessage,
		HandleConfirm: func() error {
			return self.executeCommand(command)
		},
	})

	return nil
}

// executeCommand executes the shell command.
func (self *ShellCommandAction) executeCommand(command string) error {
	self.c.LogAction(self.c.Tr.Actions.CustomCommand)
	return self.c.RunSubprocessAndRefresh(
		self.c.OS().Cmd.NewShell(command, self.c.UserConfig().OS.ShellFunctionsFile),
	)
}

// ToggleAIMode toggles AI mode for the shell command helper.
func (self *ShellCommandAction) ToggleAIMode() bool {
	if self.enhancedHelper != nil {
		return self.enhancedHelper.ToggleAIMode()
	}
	return false
}

// IsAIMode returns whether AI mode is currently active.
func (self *ShellCommandAction) IsAIMode() bool {
	if self.enhancedHelper != nil {
		return self.enhancedHelper.IsAIMode()
	}
	return false
}

// this mimics the shell functionality `ignorespace`
// which doesn't save a command to history if it starts with a space
func (self *ShellCommandAction) shouldSaveCommand(command string) bool {
	return !strings.HasPrefix(command, " ")
}
