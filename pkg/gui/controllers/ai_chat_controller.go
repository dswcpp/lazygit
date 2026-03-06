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

func (self *AIChatController) close() error {
	self.c.Context().Pop()
	return nil
}

func (self *AIChatController) Context() types.Context {
	return self.c.Contexts().AIChat
}
