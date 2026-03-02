package gui

import (
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

// createMessageBoxTestMenu 创建消息框测试菜单
func (gui *Gui) createMessageBoxTestMenu() error {
	menuItems := []*types.MenuItem{
		{
			Label: "错误消息框",
			OnPress: func() error {
				return gui.TestMessageBoxError()
			},
		},
		{
			Label: "警告消息框",
			OnPress: func() error {
				return gui.TestMessageBoxWarning()
			},
		},
		{
			Label: "信息消息框",
			OnPress: func() error {
				return gui.TestMessageBoxInfo()
			},
		},
		{
			Label: "成功消息框",
			OnPress: func() error {
				return gui.TestMessageBoxSuccess()
			},
		},
		{
			Label: "确认对话框",
			OnPress: func() error {
				return gui.TestMessageBoxConfirm()
			},
		},
		{
			Label: "是/否/取消对话框",
			OnPress: func() error {
				return gui.TestMessageBoxYesNoCancel()
			},
		},
		{
			Label: "自定义按钮",
			OnPress: func() error {
				return gui.TestMessageBoxCustomButtons()
			},
		},
		{
			Label: "自动关闭消息框",
			OnPress: func() error {
				return gui.TestMessageBoxAutoClose()
			},
		},
		{
			Label: "长文本消息框",
			OnPress: func() error {
				return gui.TestMessageBoxLongText()
			},
		},
		{
			Label: "测试所有消息类型",
			OnPress: func() error {
				return gui.TestMessageBoxAllTypes()
			},
		},
	}

	return gui.c.Menu(types.CreateMenuOptions{
		Title: "消息框测试菜单",
		Items: menuItems,
	})
}
