package controllers

import (
	"errors"

	"github.com/jesseduffield/gocui"
	"github.com/samber/lo"
	"github.com/dswcpp/lazygit/pkg/commands/models"
	"github.com/dswcpp/lazygit/pkg/gui/context"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

type ActivityBarController struct {
	baseController
	*ListControllerTrait[*models.ActivityBarItem]
	c            *ControllerCommon
	handlePull   func() error
	handlePush   func() error
}

var _ types.IController = &ActivityBarController{}

func NewActivityBarController(
	c *ControllerCommon,
	handlePull func() error,
	handlePush func() error,
) *ActivityBarController {
	return &ActivityBarController{
		baseController: baseController{},
		c:              c,
		handlePull:     handlePull,
		handlePush:     handlePush,
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().ActivityBar,
			c.Contexts().ActivityBar.GetSelected,
			c.Contexts().ActivityBar.GetSelectedItems,
		),
	}
}

func (self *ActivityBarController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.withItem(self.handleSelect),
			Description: self.c.Tr.Execute,
		},
		{
			ViewName: "activityBar",
			Key:      gocui.MouseLeft,
			Handler:  self.withItem(self.handleClick),
		},
	}

	return bindings
}

func (self *ActivityBarController) Context() types.Context {
	return self.context()
}

func (self *ActivityBarController) context() *context.ActivityBarContext {
	return self.c.Contexts().ActivityBar
}

func (self *ActivityBarController) handleSelect(item *models.ActivityBarItem) error {
	if item.IsSeparator() {
		return nil
	}

	switch item.Type {
	case models.ActivityTypeNavigation:
		return self.handleNavigation(item)
	case models.ActivityTypeAction:
		return self.handleAction(item)
	case models.ActivityTypeTool:
		return self.handleTool(item)
	case models.ActivityTypeCustom:
		return self.handleCustomCommand(item)
	}

	return nil
}

func (self *ActivityBarController) handleClick(item *models.ActivityBarItem) error {
	return self.handleSelect(item)
}

func (self *ActivityBarController) handleNavigation(item *models.ActivityBarItem) error {
	contextMap := map[string]types.Context{
		"status":   self.c.Contexts().Status,
		"files":    self.c.Contexts().Files,
		"branches": self.c.Contexts().Branches,
		"commits":  self.c.Contexts().LocalCommits,
		"stash":    self.c.Contexts().Stash,
	}

	if ctx, ok := contextMap[item.Action]; ok {
		self.c.Context().Push(ctx, types.OnFocusOpts{})
		return nil
	}

	return nil
}

func (self *ActivityBarController) handleAction(item *models.ActivityBarItem) error {
	switch item.Action {
	case "pull":
		return self.handlePull()
	case "push":
		return self.handlePush()
	case "fetch":
		return self.handleFetch()
	case "stash":
		return self.handleStashAllChanges()
	case "merge":
		// Switch to branches panel for merge operation
		self.c.Context().Push(self.c.Contexts().Branches, types.OnFocusOpts{})
		return nil
	case "rebase":
		// Switch to branches panel for rebase operation
		self.c.Context().Push(self.c.Contexts().Branches, types.OnFocusOpts{})
		return nil
	default:
		return nil
	}
}

func (self *ActivityBarController) handleTool(item *models.ActivityBarItem) error {
	switch item.Action {
	case "settings":
		return self.handleOpenConfig()
	case "help":
		// Open menu panel which contains keybindings help
		self.c.Context().Push(self.c.Contexts().Menu, types.OnFocusOpts{})
		return nil
	default:
		return nil
	}
}

// handleFetch fetches from remote
func (self *ActivityBarController) handleFetch() error {
	return self.c.WithWaitingStatus(self.c.Tr.FetchingStatus, func(task gocui.Task) error {
		self.c.LogAction("Fetch")
		err := self.c.Git().Sync.Fetch(task)

		self.c.Refresh(types.RefreshOptions{
			Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.REMOTES, types.TAGS},
			Mode:  types.SYNC,
		})

		return err
	})
}

// handleStashAllChanges stashes all changes
func (self *ActivityBarController) handleStashAllChanges() error {
	if !self.c.Helpers().WorkingTree.IsWorkingTreeDirtyExceptSubmodules() {
		return errors.New(self.c.Tr.NoFilesToStash)
	}

	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.StashChanges,
		HandleConfirm: func(stashComment string) error {
			self.c.LogAction(self.c.Tr.Actions.Stash)
			return self.c.WithWaitingStatus(self.c.Tr.Actions.Stash, func(task gocui.Task) error {
				if err := self.c.Git().Stash.Push(stashComment); err != nil {
					return err
				}
				self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				return nil
			})
		},
	})

	return nil
}

// handleOpenConfig opens the git config file for editing
func (self *ActivityBarController) handleOpenConfig() error {
	// Get user config paths (may be multiple)
	confPaths := self.c.GetConfig().GetUserConfigPaths()

	if len(confPaths) == 0 {
		return errors.New(self.c.Tr.NoConfigFileFoundErr)
	}

	// If only one config file, edit it directly
	if len(confPaths) == 1 {
		return self.c.Helpers().Files.EditFiles([]string{confPaths[0]})
	}

	// If multiple config files, show menu to choose
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.EditConfig,
		Items: lo.Map(confPaths, func(path string, _ int) *types.MenuItem {
			return &types.MenuItem{
				Label: path,
				OnPress: func() error {
					return self.c.Helpers().Files.EditFiles([]string{path})
				},
			}
		}),
	})
}

func (self *ActivityBarController) handleCustomCommand(item *models.ActivityBarItem) error {
	if item.CustomCmd == "" {
		return nil
	}

	// Execute custom command using the shell command runner
	// This integrates with lazygit's existing custom command system which handles
	// command execution safely
	cmdObj := self.c.OS().Cmd.NewShell(item.CustomCmd, self.c.UserConfig().OS.ShellFunctionsFile)
	return self.c.RunSubprocessAndRefresh(cmdObj)
}
