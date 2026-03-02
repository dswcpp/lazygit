package gui

import (
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

// createProgressBarTestMenu 创建进度条测试菜单
func (gui *Gui) createProgressBarTestMenu() error {
	menuItems := []*types.MenuItem{
		{
			Label: "确定进度条（推送示例）",
			OnPress: func() error {
				return gui.TestProgressBarDeterminate()
			},
		},
		{
			Label: "不确定进度条（克隆示例）",
			OnPress: func() error {
				return gui.TestProgressBarIndeterminate()
			},
		},
		{
			Label: "测试所有进度条样式",
			OnPress: func() error {
				return gui.TestProgressBarStyles()
			},
		},
		{
			Label: "测试所有旋转动画",
			OnPress: func() error {
				return gui.TestProgressBarSpinners()
			},
		},
		{
			Label: "模拟 Git Push",
			OnPress: func() error {
				return gui.pushWithProgress()
			},
		},
		{
			Label: "模拟 Git Clone",
			OnPress: func() error {
				return gui.cloneWithProgress("https://github.com/example/repo.git")
			},
		},
		{
			Label: "模拟 Git Fetch",
			OnPress: func() error {
				return gui.fetchWithProgress()
			},
		},
	}

	return gui.c.Menu(types.CreateMenuOptions{
		Title: "进度条测试菜单",
		Items: menuItems,
	})
}
