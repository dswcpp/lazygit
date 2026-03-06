package helpers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dswcpp/lazygit/pkg/ai/agent"
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

func TestDeriveAIChatStatus_WaitingConfirmShowsExpectedPrompt(t *testing.T) {
	sess := agent.NewSession("")
	sess.SetPhase(agent.PhaseWaitingConfirm)

	status, detail := deriveAIChatStatus(sess, false, "", "")

	assert.Equal(t, "等待确认", status)
	assert.Equal(t, "输入 Y 执行，N 取消，或输入补充说明", detail)
}

func TestDeriveAIChatStatus_ExecutingUsesLatestStepUpdate(t *testing.T) {
	sess := agent.NewSession("")
	sess.SetPhase(agent.PhaseExecuting)
	sess.UIMessages = append(sess.UIMessages, agent.UIMessage{
		Kind:    agent.KindStepUpdate,
		Content: "[执行中] 正在提交当前修改\n更多细节",
	})

	status, detail := deriveAIChatStatus(sess, false, "", "")

	assert.Equal(t, "执行中", status)
	assert.Equal(t, "[执行中] 正在提交当前修改", detail)
}

func TestDeriveAIChatStatus_UsesStoredTerminalStatusWhenAgentIsGone(t *testing.T) {
	status, detail := deriveAIChatStatus(nil, false, "已完成", "可输入下一条指令")

	assert.Equal(t, "已完成", status)
	assert.Equal(t, "可输入下一条指令", detail)
}

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
