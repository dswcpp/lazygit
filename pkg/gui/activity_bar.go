package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/gui/style"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

// ActivityBarItem 表示活动栏中的一个项目
type ActivityBarItem struct {
	Icon    string
	Name    string
	Type    string // "navigation", "action", "tool", "separator"
	Tooltip string
}

// 获取活动栏项目列表
func (gui *Gui) getActivityBarItems() []*ActivityBarItem {
	return []*ActivityBarItem{
		// 导航区
		{Icon: "📁", Name: "status", Type: "navigation", Tooltip: "Status"},
		{Icon: "🌿", Name: "files", Type: "navigation", Tooltip: "Files"},
		{Icon: "📊", Name: "branches", Type: "navigation", Tooltip: "Branches"},
		{Icon: "💾", Name: "commits", Type: "navigation", Tooltip: "Commits"},
		{Icon: "📦", Name: "stash", Type: "navigation", Tooltip: "Stash"},

		// 分隔符
		{Icon: "──", Name: "separator1", Type: "separator"},

		// 操作区
		{Icon: "⬇️", Name: "pull", Type: "action", Tooltip: "Pull"},
		{Icon: "⬆️", Name: "push", Type: "action", Tooltip: "Push"},
		{Icon: "🔄", Name: "fetch", Type: "action", Tooltip: "Fetch"},
		{Icon: "💾", Name: "stash-action", Type: "action", Tooltip: "Stash changes"},
		{Icon: "🔀", Name: "merge", Type: "action", Tooltip: "Merge"},
		{Icon: "♻️", Name: "rebase", Type: "action", Tooltip: "Rebase"},

		// 分隔符
		{Icon: "──", Name: "separator2", Type: "separator"},

		// 工具区
		{Icon: "⚙️", Name: "settings", Type: "tool", Tooltip: "Settings"},
		{Icon: "❓", Name: "help", Type: "tool", Tooltip: "Help"},
	}
}

// 初始化活动栏（仅在启用时调用）
func (gui *Gui) initActivityBar() error {
	if !gui.c.UserConfig().Gui.ActivityBar.Show {
		return nil // 未启用，直接返回
	}

	return gui.renderActivityBar()
}

// 渲染活动栏内容
func (gui *Gui) renderActivityBar() error {
	if !gui.c.UserConfig().Gui.ActivityBar.Show {
		return nil
	}

	view, err := gui.g.View("activityBar")
	if err != nil {
		return nil // 视图可能还未创建
	}

	view.Clear()
	view.Frame = false
	view.FgColor = gocui.ColorDefault

	items := gui.getActivityBarItems()
	selectedIdx := gui.State.Panels.ActivityBar.SelectedLine

	for i, item := range items {
		icon := item.Icon

		// 当前选中项添加背景色
		if i == selectedIdx {
			icon = style.FgBlue.SetBold().Sprint(icon)
		}

		fmt.Fprintln(view, icon)
	}

	return nil
}

// 处理活动栏项目选择
func (gui *Gui) handleActivityBarSelect() error {
	if !gui.c.UserConfig().Gui.ActivityBar.Show {
		return nil
	}

	items := gui.getActivityBarItems()
	selectedIdx := gui.State.Panels.ActivityBar.SelectedLine

	if selectedIdx < 0 || selectedIdx >= len(items) {
		return nil
	}

	item := items[selectedIdx]

	switch item.Type {
	case "navigation":
		return gui.handleActivityBarNavigation(item)
	case "action":
		return gui.handleActivityBarAction(item)
	case "tool":
		return gui.handleActivityBarTool(item)
	case "separator":
		return nil
	}

	return nil
}

// 处理导航操作
func (gui *Gui) handleActivityBarNavigation(item *ActivityBarItem) error {
	switch item.Name {
	case "status":
		gui.c.Context().Push(gui.State.Contexts.Status, types.OnFocusOpts{})
	case "files":
		gui.c.Context().Push(gui.State.Contexts.Files, types.OnFocusOpts{})
	case "branches":
		gui.c.Context().Push(gui.State.Contexts.Branches, types.OnFocusOpts{})
	case "commits":
		gui.c.Context().Push(gui.State.Contexts.LocalCommits, types.OnFocusOpts{})
	case "stash":
		gui.c.Context().Push(gui.State.Contexts.Stash, types.OnFocusOpts{})
	}
	return nil
}

// 处理 Git 操作
func (gui *Gui) handleActivityBarAction(item *ActivityBarItem) error {
	// 暂时返回提示信息，后续可以连接到实际的 Git 操作
	switch item.Name {
	case "pull":
		gui.c.Toast("Pulling...")
		return nil
	case "push":
		gui.c.Toast("Pushing...")
		return nil
	case "fetch":
		gui.c.Toast("Fetching...")
		return nil
	case "stash-action":
		gui.c.Toast("Stashing changes...")
		return nil
	case "merge":
		gui.c.Toast("Opening merge menu...")
		return nil
	case "rebase":
		gui.c.Toast("Opening rebase menu...")
		return nil
	}
	return nil
}

// 处理工具操作
func (gui *Gui) handleActivityBarTool(item *ActivityBarItem) error {
	switch item.Name {
	case "settings":
		gui.c.Toast("Opening settings...")
		return nil
	case "help":
		gui.c.Toast("Opening help...")
		return nil
	}
	return nil
}

// 活动栏上移
func (gui *Gui) handleActivityBarPrevLine() error {
	if !gui.c.UserConfig().Gui.ActivityBar.Show {
		return nil
	}

	items := gui.getActivityBarItems()
	if gui.State.Panels.ActivityBar.SelectedLine > 0 {
		gui.State.Panels.ActivityBar.SelectedLine--
		// 跳过分隔符
		for gui.State.Panels.ActivityBar.SelectedLine > 0 &&
			items[gui.State.Panels.ActivityBar.SelectedLine].Type == "separator" {
			gui.State.Panels.ActivityBar.SelectedLine--
		}
	}

	return gui.renderActivityBar()
}

// 活动栏下移
func (gui *Gui) handleActivityBarNextLine() error {
	if !gui.c.UserConfig().Gui.ActivityBar.Show {
		return nil
	}

	items := gui.getActivityBarItems()
	if gui.State.Panels.ActivityBar.SelectedLine < len(items)-1 {
		gui.State.Panels.ActivityBar.SelectedLine++
		// 跳过分隔符
		for gui.State.Panels.ActivityBar.SelectedLine < len(items)-1 &&
			items[gui.State.Panels.ActivityBar.SelectedLine].Type == "separator" {
			gui.State.Panels.ActivityBar.SelectedLine++
		}
	}

	return gui.renderActivityBar()
}

// 处理活动栏鼠标点击
func (gui *Gui) handleActivityBarClick(opts gocui.ViewMouseBindingOpts) error {
	if !gui.c.UserConfig().Gui.ActivityBar.Show {
		return nil
	}

	// 设置选中行
	gui.State.Panels.ActivityBar.SelectedLine = opts.Y

	// 重新渲染
	if err := gui.renderActivityBar(); err != nil {
		return err
	}

	// 执行选中项
	return gui.handleActivityBarSelect()
}
