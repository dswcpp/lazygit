package controllers

import (
	"github.com/jesseduffield/gocui"
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
		// Fetch 在后台自动运行，这里只显示提示
		self.c.Toast("Fetching from remote...")
		return nil
	case "stash":
		// TODO: 需要集成 stash 创建功能
		self.c.Toast("Stash 功能即将推出")
		return nil
	case "merge":
		// TODO: 需要集成 merge 菜单
		self.c.Toast("Merge 功能即将推出")
		return nil
	case "rebase":
		// TODO: 需要集成 rebase 菜单
		self.c.Toast("Rebase 功能即将推出")
		return nil
	default:
		self.c.Toast("未知操作: " + item.Action)
		return nil
	}
}

func (self *ActivityBarController) handleTool(item *models.ActivityBarItem) error {
	switch item.Action {
	case "settings":
		// TODO: Open settings panel
		self.c.Toast("Settings 功能即将推出")
		return nil
	case "help":
		// 显示帮助提示
		self.c.Toast("按 '?' 键查看所有快捷键")
		return nil
	default:
		self.c.Toast("未知工具: " + item.Action)
		return nil
	}
}

func (self *ActivityBarController) handleCustomCommand(item *models.ActivityBarItem) error {
	if item.CustomCmd == "" {
		return nil
	}

	// TODO: Execute custom command
	// This would integrate with the existing custom commands system
	self.c.Toast("执行自定义命令: " + item.CustomCmd)
	return nil
}
