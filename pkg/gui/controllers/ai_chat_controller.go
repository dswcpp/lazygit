package controllers

import (
	"github.com/dswcpp/lazygit/pkg/gui/types"
	"github.com/jesseduffield/gocui"
)

// AIChatController 处理 AI 聊天弹窗的键盘绑定。
// 遵循 lazygit Context 系统：键位由框架管理，resetKeybindings 后自动恢复。
type AIChatController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &AIChatController{}

func NewAIChatController(c *ControllerCommon) *AIChatController {
	return &AIChatController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *AIChatController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	historyViewName := self.Context().GetViewName()

	return []*types.Binding{
		// 'i' / Tab：把焦点切回输入条（从历史滚动视图返回输入）
		{
			Key:         'i',
			Handler:     self.focusInput,
			Description: "聚焦输入框",
		},
		{
			Key:         gocui.KeyTab,
			Handler:     self.focusInput,
			Description: "聚焦输入框",
		},
		{
			Key:         'c',
			Handler:     self.copyToClipboard,
			Description: "复制最后一条 AI 回复",
		},
		{
			Key:         'x',
			Handler:     self.executeCommands,
			Description: "执行 AI 回复中的命令",
		},
		{
			Key:         'z',
			Handler:     self.toggleZoom,
			Description: "切换全屏模式",
		},
		{
			ViewName:    historyViewName,
			Key:         gocui.KeyArrowUp,
			Handler:     self.handleScrollUp,
			Description: "向上滚动聊天历史",
		},
		{
			ViewName:    historyViewName,
			Key:         gocui.KeyArrowDown,
			Handler:     self.handleScrollDown,
			Description: "向下滚动聊天历史",
		},
		{
			ViewName:    historyViewName,
			Key:         gocui.KeyPgup,
			Handler:     self.handlePageUp,
			Description: "向上翻页",
		},
		{
			ViewName:    historyViewName,
			Key:         gocui.KeyPgdn,
			Handler:     self.handlePageDown,
			Description: "向下翻页",
		},
		{
			ViewName:    historyViewName,
			Key:         gocui.KeyHome,
			Handler:     self.handleGoTop,
			Description: "跳到顶部",
		},
		{
			ViewName:    historyViewName,
			Key:         gocui.KeyEnd,
			Handler:     self.handleGoBottom,
			Description: "跳到底部",
		},
		{
			Key:         'q',
			Handler:     self.close,
			Description: self.c.Tr.Close,
		},
	}
}

// focusInput 将 gocui 焦点切到输入条视图（光标出现，可以开始打字）
func (self *AIChatController) focusInput() error {
	_, _ = self.c.GocuiGui().SetCurrentView(self.c.Views().AIChatInput.Name())
	return nil
}

func (self *AIChatController) GetOnFocusLost() func(types.OnFocusLostOpts) {
	return func(types.OnFocusLostOpts) {
		hideAIChatPopupViews(self.c.Views())
	}
}

func (self *AIChatController) copyToClipboard() error {
	return self.c.Helpers().AIChat.CopyLastResponse()
}

func (self *AIChatController) executeCommands() error {
	return self.c.Helpers().AIChat.ExecuteLastCommands()
}

func (self *AIChatController) toggleZoom() error {
	self.c.Contexts().AIChat.Zoomed = !self.c.Contexts().AIChat.Zoomed
	self.c.Helpers().Confirmation.ResizeCurrentPopupPanels()
	return nil
}

func (self *AIChatController) handleScrollUp() error {
	scrollAIChatViewUp(self.c.Views().AIChat, resolveAIChatScrollHeight(self.c.UserConfig().Gui.ScrollHeight))
	return nil
}

func (self *AIChatController) handleScrollDown() error {
	scrollAIChatViewDown(self.c.Views().AIChat, resolveAIChatScrollHeight(self.c.UserConfig().Gui.ScrollHeight))
	return nil
}

func (self *AIChatController) handlePageUp() error {
	pageAIChatViewUp(self.c.Views().AIChat)
	return nil
}

func (self *AIChatController) handlePageDown() error {
	pageAIChatViewDown(self.c.Views().AIChat)
	return nil
}

func (self *AIChatController) handleGoTop() error {
	scrollAIChatViewToTop(self.c.Views().AIChat)
	return nil
}

func (self *AIChatController) handleGoBottom() error {
	scrollAIChatViewToBottom(self.c.Views().AIChat)
	return nil
}

func (self *AIChatController) close() error {
	self.c.Context().Pop()
	return nil
}

func (self *AIChatController) Context() types.Context {
	return self.c.Contexts().AIChat
}

func hideAIChatPopupViews(views types.Views) {
	if views.AIChat != nil {
		views.AIChat.Visible = false
	}
	if views.AIChatInput != nil {
		views.AIChatInput.Visible = false
	}
}

func resolveAIChatScrollHeight(scrollHeight int) int {
	if scrollHeight < 1 {
		return 1
	}
	return scrollHeight
}

func scrollAIChatViewUp(view *gocui.View, scrollHeight int) {
	if view == nil {
		return
	}
	view.ScrollUp(resolveAIChatScrollHeight(scrollHeight))
}

func scrollAIChatViewDown(view *gocui.View, scrollHeight int) {
	if view == nil {
		return
	}
	view.ScrollDown(resolveAIChatScrollHeight(scrollHeight))
}

func pageAIChatViewUp(view *gocui.View) {
	if view == nil {
		return
	}
	view.ScrollUp(aiChatPageDelta(view))
}

func pageAIChatViewDown(view *gocui.View) {
	if view == nil {
		return
	}
	view.ScrollDown(aiChatPageDelta(view))
}

func scrollAIChatViewToTop(view *gocui.View) {
	if view == nil {
		return
	}
	view.ScrollUp(view.ViewLinesHeight())
}

func scrollAIChatViewToBottom(view *gocui.View) {
	if view == nil {
		return
	}
	view.ScrollDown(view.ViewLinesHeight())
}

func aiChatPageDelta(view *gocui.View) int {
	if view == nil {
		return 1
	}

	delta := view.InnerHeight() - 1
	if delta < 1 {
		return 1
	}

	return delta
}
