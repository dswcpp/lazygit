package controllers

import (
	"github.com/dswcpp/lazygit/pkg/gui/types"
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
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.newMessage,
			Description: "发送新消息",
		},
		{
			Key:         'n',
			Handler:     self.newMessage,
			Description: "发送新消息",
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
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.close,
			Description: self.c.Tr.Close,
		},
		{
			Key:         'q',
			Handler:     self.close,
			Description: self.c.Tr.Close,
		},
	}
}

// newMessage 通过标准 Prompt 弹出输入框，获取用户消息后发送给 AI
func (self *AIChatController) newMessage() error {
	self.c.Prompt(types.PromptOpts{
		Title: "Ask AI",
		HandleConfirm: func(content string) error {
			return self.c.Helpers().AIChat.SendMessage(content)
		},
	})
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
