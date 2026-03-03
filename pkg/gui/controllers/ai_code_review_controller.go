package controllers

import (
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

type AICodeReviewController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &AICodeReviewController{}

func NewAICodeReviewController(c *ControllerCommon) *AICodeReviewController {
	return &AICodeReviewController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *AICodeReviewController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.close,
			Description: self.c.Tr.Close,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.cancel,
			Description: "取消 AI 请求",
		},
		{
			Key:         'c',
			Handler:     self.copyToClipboard,
			Description: self.c.Tr.CopyToClipboardMenu,
		},
		{
			Key:         'z',
			Handler:     self.toggleZoom,
			Description: self.c.Tr.AICodeReviewToggleZoom,
		},
	}
}

func (self *AICodeReviewController) close() error {
	self.c.Context().Pop()
	return nil
}

func (self *AICodeReviewController) copyToClipboard() error {
	content := self.c.Views().AICodeReview.Buffer()
	if err := self.c.OS().CopyToClipboard(content); err != nil {
		return err
	}
	self.c.Toast(self.c.Tr.AICodeReviewCopiedToClipboard)
	return nil
}

func (self *AICodeReviewController) toggleZoom() error {
	self.c.Contexts().AICodeReview.Zoomed = !self.c.Contexts().AICodeReview.Zoomed
	self.c.Helpers().Confirmation.ResizeCurrentPopupPanels()
	return nil
}

func (self *AICodeReviewController) Context() types.Context {
	return self.c.Contexts().AICodeReview
}

func (self *AICodeReviewController) cancel() error {
	ctx := self.c.Contexts().AICodeReview
	if ctx.CancelFunc != nil {
		ctx.CancelFunc()
		self.c.Toast("正在取消 AI 请求...")
	}
	return nil
}
