package controllers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dswcpp/lazygit/pkg/common"
	"github.com/dswcpp/lazygit/pkg/config"
	guicontext "github.com/dswcpp/lazygit/pkg/gui/context"
	guihelpers "github.com/dswcpp/lazygit/pkg/gui/controllers/helpers"
	"github.com/dswcpp/lazygit/pkg/gui/types"
	"github.com/dswcpp/lazygit/pkg/i18n"
	"github.com/jesseduffield/gocui"
	"github.com/stretchr/testify/assert"
)

type testAIChatContexts struct {
	tree *guicontext.ContextTree
}

func (t testAIChatContexts) Contexts() *guicontext.ContextTree {
	return t.tree
}

func TestAIChatController_OnFocusLost_HidesInputView(t *testing.T) {
	views := types.Views{
		AIChat:      newTestChatView("aiChat"),
		AIChatInput: newTestChatView("aiChatInput"),
	}
	views.AIChat.Visible = true
	views.AIChatInput.Visible = true

	hideAIChatPopupViews(views)

	assert.False(t, views.AIChat.Visible)
	assert.False(t, views.AIChatInput.Visible)
}

func TestAIChatController_ScrollHandlers_ScrollByConfiguredHeight(t *testing.T) {
	view := newFilledChatView("aiChat", 30)

	scrollAIChatViewDown(view, 2)
	_, originY := view.Origin()
	assert.Equal(t, 2, originY)

	scrollAIChatViewDown(view, 0)
	_, originY = view.Origin()
	assert.Equal(t, 3, originY)

	scrollAIChatViewUp(view, 2)
	_, originY = view.Origin()
	assert.Equal(t, 1, originY)
}

func TestAIChatController_PageAndJumpHandlers_MoveToExpectedPositions(t *testing.T) {
	view := newFilledChatView("aiChat", 40)

	pageAIChatViewDown(view)
	_, originY := view.Origin()
	assert.Equal(t, aiChatPageDelta(view), originY)

	pageAIChatViewUp(view)
	_, originY = view.Origin()
	assert.Equal(t, 0, originY)

	scrollAIChatViewDown(view, 7)
	scrollAIChatViewToTop(view)
	_, originY = view.Origin()
	assert.Equal(t, 0, originY)

	scrollAIChatViewToBottom(view)
	_, originY = view.Origin()
	assert.Equal(t, maxChatViewOriginY(view), originY)
}

func TestAIChatController_Keybindings_KeepExistingActionsAndAddScrollBindings(t *testing.T) {
	controller := newAIChatControllerForBindingsTest()

	bindings := controller.GetKeybindings(types.KeybindingsOpts{})

	assert.True(t, hasBinding(bindings, "", 'i'))
	assert.True(t, hasBinding(bindings, "", gocui.KeyTab))
	assert.True(t, hasBinding(bindings, "", 'c'))
	assert.True(t, hasBinding(bindings, "", 'x'))
	assert.True(t, hasBinding(bindings, "", 'z'))
	assert.True(t, hasBinding(bindings, "", 'q'))

	assert.True(t, hasBinding(bindings, "aiChat", gocui.KeyArrowUp))
	assert.True(t, hasBinding(bindings, "aiChat", gocui.KeyArrowDown))
	assert.True(t, hasBinding(bindings, "aiChat", gocui.KeyPgup))
	assert.True(t, hasBinding(bindings, "aiChat", gocui.KeyPgdn))
	assert.True(t, hasBinding(bindings, "aiChat", gocui.KeyHome))
	assert.True(t, hasBinding(bindings, "aiChat", gocui.KeyEnd))
}

func newAIChatControllerForBindingsTest() *AIChatController {
	chatView := newTestChatView("aiChat")
	chatContext := &guicontext.AIChatContext{
		SimpleContext: guicontext.NewSimpleContext(guicontext.NewBaseContext(guicontext.NewBaseContextOpts{
			View:                  chatView,
			WindowName:            "aiChat",
			Key:                   guicontext.AI_CHAT_CONTEXT_KEY,
			Kind:                  types.TEMPORARY_POPUP,
			Focusable:             true,
			HasUncontrolledBounds: true,
		})),
	}
	contextTree := &guicontext.ContextTree{
		AIChat: chatContext,
	}

	appCommon := &common.Common{
		Tr: &i18n.TranslationSet{
			Close: "关闭",
		},
	}
	appCommon.SetUserConfig(&config.UserConfig{
		Gui: config.GuiConfig{
			ScrollHeight: 2,
		},
	})

	return NewAIChatController(&ControllerCommon{
		HelperCommon: &guihelpers.HelperCommon{
			Common:       appCommon,
			IGetContexts: testAIChatContexts{tree: contextTree},
		},
	})
}

func newTestChatView(name string) *gocui.View {
	view := gocui.NewView(name, 0, 0, 40, 10, gocui.OutputNormal)
	view.Wrap = true
	return view
}

func newFilledChatView(name string, lineCount int) *gocui.View {
	view := newTestChatView(name)
	lines := make([]string, 0, lineCount)
	for i := 0; i < lineCount; i++ {
		lines = append(lines, fmt.Sprintf("line %02d", i))
	}
	view.SetContent(strings.Join(lines, "\n"))
	return view
}

func maxChatViewOriginY(view *gocui.View) int {
	maxOriginY := view.ViewLinesHeight() - view.InnerHeight()
	if maxOriginY < 0 {
		return 0
	}
	return maxOriginY
}

func hasBinding(bindings []*types.Binding, viewName string, key types.Key) bool {
	for _, binding := range bindings {
		if binding.ViewName == viewName && binding.Key == key {
			return true
		}
	}
	return false
}
