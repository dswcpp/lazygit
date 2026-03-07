package helpers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jesseduffield/gocui"
	"github.com/stretchr/testify/assert"
)

func TestAIChatSessionRender_DoesNotAutoScrollWithoutScrollToBottom(t *testing.T) {
	view := newFilledAIChatHelperView("aiChat", 32)
	view.ScrollDown(5)
	_, beforeOriginY := view.Origin()
	scrollToBottom := false

	applyAIChatAutoScroll(view, &scrollToBottom)

	_, afterOriginY := view.Origin()
	assert.Equal(t, beforeOriginY, afterOriginY)
	assert.False(t, scrollToBottom)
}

func TestApplyAIChatAutoScroll_ScrollsToBottomWhenRequested(t *testing.T) {
	view := newFilledAIChatHelperView("aiChat", 32)
	scrollToBottom := true

	applyAIChatAutoScroll(view, &scrollToBottom)

	_, originY := view.Origin()
	assert.Equal(t, maxAIChatHelperViewOriginY(view), originY)
	assert.False(t, scrollToBottom)
}

func TestResetAIChatInputView_ClearsTextAreaState(t *testing.T) {
	view := gocui.NewView("aiChatInput", 0, 0, 40, 5, gocui.OutputNormal)
	view.TextArea.TypeCharacter("first command")
	view.RenderTextArea()
	view.SetCursor(3, 0)
	view.SetOrigin(2, 1)

	ResetAIChatInputView(view)

	cursorX, cursorY := view.Cursor()
	originX, originY := view.Origin()
	assert.Equal(t, "", view.Buffer())
	assert.Equal(t, "", view.TextArea.GetContent())
	assert.Equal(t, 0, cursorX)
	assert.Equal(t, 0, cursorY)
	assert.Equal(t, 0, originX)
	assert.Equal(t, 0, originY)
}

// Note: deriveAIChatStatus tests removed as the function has been refactored
// into AIChatSession methods (deriveStatus, getStatusPresentation)

func newFilledAIChatHelperView(name string, lineCount int) *gocui.View {
	view := gocui.NewView(name, 0, 0, 40, 10, gocui.OutputNormal)
	view.Wrap = true

	lines := make([]string, 0, lineCount)
	for i := 0; i < lineCount; i++ {
		lines = append(lines, fmt.Sprintf("line %02d", i))
	}
	view.SetContent(strings.Join(lines, "\n"))

	return view
}

func maxAIChatHelperViewOriginY(view *gocui.View) int {
	maxOriginY := view.ViewLinesHeight() - view.InnerHeight()
	if maxOriginY < 0 {
		return 0
	}
	return maxOriginY
}
