package gui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/gui/style"
)

// MessageType 消息类型
type MessageType int

const (
	MessageTypeInfo MessageType = iota
	MessageTypeSuccess
	MessageTypeWarning
	MessageTypeError
	MessageTypeQuestion
)

// MessageBoxConfig 消息框配置
type MessageBoxConfig struct {
	Type    MessageType // 消息类型
	Title   string      // 标题
	Message string      // 消息内容
	Details string      // 详细信息（可选）
	Buttons []string    // 按钮列表
	Width   int         // 宽度
	Height  int         // 高度（0表示自动）
}

// MessageBox 消息框
type MessageBox struct {
	config       MessageBoxConfig
	view         *gocui.View
	gui          *Gui
	selectedBtn  int
	onClose      func(buttonIndex int)
	done         chan int
}

// NewMessageBox 创建消息框
func (gui *Gui) NewMessageBox(config MessageBoxConfig) *MessageBox {
	if config.Width == 0 {
		config.Width = 60
	}
	if len(config.Buttons) == 0 {
		config.Buttons = []string{"确定"}
	}

	return &MessageBox{
		config:      config,
		gui:         gui,
		selectedBtn: 0,
		done:        make(chan int),
	}
}

// getIcon 获取消息类型图标
func (mt MessageType) getIcon() string {
	icons := map[MessageType]string{
		MessageTypeInfo:     "ℹ",
		MessageTypeSuccess:  "✓",
		MessageTypeWarning:  "⚠",
		MessageTypeError:    "✗",
		MessageTypeQuestion: "?",
	}
	return icons[mt]
}

// getColor 获取消息类型颜色
func (mt MessageType) getColor() style.TextStyle {
	colors := map[MessageType]style.TextStyle{
		MessageTypeInfo:     style.FgCyan,
		MessageTypeSuccess:  style.FgGreen,
		MessageTypeWarning:  style.FgYellow,
		MessageTypeError:    style.FgRed,
		MessageTypeQuestion: style.FgCyan,
	}
	return colors[mt]
}

// Render 渲染消息框
func (mb *MessageBox) Render() string {
	var lines []string

	// 标题行（带图标）
	icon := mb.config.Type.getIcon()
	color := mb.config.Type.getColor()
	title := fmt.Sprintf(" %s %s", color.Sprint(icon), mb.config.Title)
	lines = append(lines, title)

	// 分隔线
	lines = append(lines, " "+strings.Repeat("─", mb.config.Width-4))

	// 空行
	lines = append(lines, "")

	// 消息内容（自动换行）
	messageLines := mb.wrapText(mb.config.Message, mb.config.Width-6)
	for _, line := range messageLines {
		lines = append(lines, "  "+line)
	}

	// 详细信息
	if mb.config.Details != "" {
		lines = append(lines, "")
		detailLines := mb.wrapText(mb.config.Details, mb.config.Width-6)
		for _, line := range detailLines {
			lines = append(lines, "  "+style.FgDefault.Sprint(line))
		}
	}

	// 空行
	lines = append(lines, "")

	// 分隔线
	lines = append(lines, " "+strings.Repeat("─", mb.config.Width-4))

	// 按钮行
	buttonLine := mb.renderButtons()
	lines = append(lines, buttonLine)

	return strings.Join(lines, "\n")
}

// wrapText 文本自动换行
func (mb *MessageBox) wrapText(text string, width int) []string {
	if width <= 0 {
		width = 50
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	currentLine := ""
	for _, word := range words {
		if len(currentLine)+len(word)+1 <= width {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// renderButtons 渲染按钮
func (mb *MessageBox) renderButtons() string {
	var buttons []string

	for i, btn := range mb.config.Buttons {
		if i == mb.selectedBtn {
			// 选中的按钮
			buttons = append(buttons, style.FgBlack.SetBg(style.NewBasicColor(color.BgCyan)).SetBold().Sprintf(" %s ", btn))
		} else {
			// 未选中的按钮
			buttons = append(buttons, fmt.Sprintf("[ %s ]", btn))
		}
	}

	// 居中显示按钮
	buttonStr := strings.Join(buttons, "  ")
	padding := (mb.config.Width - len(buttonStr) - 2) / 2
	if padding < 0 {
		padding = 0
	}

	return " " + strings.Repeat(" ", padding) + buttonStr
}

// ShowMessageBox 显示消息框
func (gui *Gui) ShowMessageBox(config MessageBoxConfig, onClose func(buttonIndex int)) *MessageBox {
	mb := gui.NewMessageBox(config)
	mb.onClose = onClose

	// 创建弹出窗口
	gui.createMessageBoxPopup(mb)

	// 设置键盘绑定
	gui.setMessageBoxKeyBindings(mb)

	return mb
}

// createMessageBoxPopup 创建消息框弹出窗口
func (gui *Gui) createMessageBoxPopup(mb *MessageBox) {
	width := mb.config.Width
	height := mb.config.Height

	// 如果高度为0，自动计算
	if height == 0 {
		// 计算内容行数
		lines := strings.Split(mb.Render(), "\n")
		height = len(lines) + 2
	}

	maxX, maxY := gui.g.Size()
	x0 := (maxX - width) / 2
	y0 := (maxY - height) / 2
	x1 := x0 + width
	y1 := y0 + height

	view, err := gui.g.SetView("messageBox", x0, y0, x1, y1, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return
	}

	view.Frame = true
	view.Title = ""
	mb.view = view

	// 渲染内容
	view.Clear()
	fmt.Fprint(view, mb.Render())

	gui.g.SetViewOnTop("messageBox")
	gui.g.SetCurrentView("messageBox")
}

// setMessageBoxKeyBindings 设置消息框键盘绑定
func (gui *Gui) setMessageBoxKeyBindings(mb *MessageBox) {
	// 左箭头 - 上一个按钮
	gui.g.SetKeybinding("messageBox", gocui.KeyArrowLeft, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		if mb.selectedBtn > 0 {
			mb.selectedBtn--
			mb.view.Clear()
			fmt.Fprint(mb.view, mb.Render())
		}
		return nil
	})

	// 右箭头 - 下一个按钮
	gui.g.SetKeybinding("messageBox", gocui.KeyArrowRight, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		if mb.selectedBtn < len(mb.config.Buttons)-1 {
			mb.selectedBtn++
			mb.view.Clear()
			fmt.Fprint(mb.view, mb.Render())
		}
		return nil
	})

	// Tab - 下一个按钮
	gui.g.SetKeybinding("messageBox", gocui.KeyTab, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		mb.selectedBtn = (mb.selectedBtn + 1) % len(mb.config.Buttons)
		mb.view.Clear()
		fmt.Fprint(mb.view, mb.Render())
		return nil
	})

	// Enter - 确认
	gui.g.SetKeybinding("messageBox", gocui.KeyEnter, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		gui.CloseMessageBox(mb, mb.selectedBtn)
		return nil
	})

	// Esc - 取消（选择最后一个按钮）
	gui.g.SetKeybinding("messageBox", gocui.KeyEsc, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		gui.CloseMessageBox(mb, len(mb.config.Buttons)-1)
		return nil
	})

	// 数字键 1-9 - 快速选择按钮
	for i := 0; i < 9 && i < len(mb.config.Buttons); i++ {
		idx := i
		key := gocui.Key('1' + i)
		gui.g.SetKeybinding("messageBox", key, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
			gui.CloseMessageBox(mb, idx)
			return nil
		})
	}
}

// CloseMessageBox 关闭消息框
func (gui *Gui) CloseMessageBox(mb *MessageBox, buttonIndex int) {
	// 删除键盘绑定
	gui.g.DeleteViewKeybindings("messageBox")

	// 删除视图
	gui.g.DeleteView("messageBox")

	// 调用回调
	if mb.onClose != nil {
		mb.onClose(buttonIndex)
	}

	// 发送完成信号
	select {
	case mb.done <- buttonIndex:
	default:
	}
}

// 便捷方法

// ShowError 显示错误消息框
func (gui *Gui) ShowError(title, message string, details ...string) {
	detailStr := ""
	if len(details) > 0 {
		detailStr = details[0]
	}

	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeError,
		Title:   title,
		Message: message,
		Details: detailStr,
		Buttons: []string{"确定"},
	}, nil)
}

// ShowWarning 显示警告消息框
func (gui *Gui) ShowWarning(title, message string, details ...string) {
	detailStr := ""
	if len(details) > 0 {
		detailStr = details[0]
	}

	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeWarning,
		Title:   title,
		Message: message,
		Details: detailStr,
		Buttons: []string{"确定"},
	}, nil)
}

// ShowInfo 显示信息消息框
func (gui *Gui) ShowInfo(title, message string, details ...string) {
	detailStr := ""
	if len(details) > 0 {
		detailStr = details[0]
	}

	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeInfo,
		Title:   title,
		Message: message,
		Details: detailStr,
		Buttons: []string{"确定"},
	}, nil)
}

// ShowSuccess 显示成功消息框
func (gui *Gui) ShowSuccess(title, message string, details ...string) {
	detailStr := ""
	if len(details) > 0 {
		detailStr = details[0]
	}

	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeSuccess,
		Title:   title,
		Message: message,
		Details: detailStr,
		Buttons: []string{"确定"},
	}, nil)
}

// ShowConfirm 显示确认对话框
func (gui *Gui) ShowConfirm(title, message string, onConfirm func()) {
	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeQuestion,
		Title:   title,
		Message: message,
		Buttons: []string{"确定", "取消"},
	}, func(buttonIndex int) {
		if buttonIndex == 0 && onConfirm != nil {
			onConfirm()
		}
	})
}

// ShowYesNoCancel 显示是/否/取消对话框
func (gui *Gui) ShowYesNoCancel(title, message string, onYes, onNo func()) {
	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeQuestion,
		Title:   title,
		Message: message,
		Buttons: []string{"是", "否", "取消"},
	}, func(buttonIndex int) {
		switch buttonIndex {
		case 0:
			if onYes != nil {
				onYes()
			}
		case 1:
			if onNo != nil {
				onNo()
			}
		}
	})
}

// ShowAutoCloseMessage 显示自动关闭的消息框
func (gui *Gui) ShowAutoCloseMessage(msgType MessageType, title, message string, duration time.Duration) {
	mb := gui.ShowMessageBox(MessageBoxConfig{
		Type:    msgType,
		Title:   title,
		Message: message,
		Buttons: []string{"确定"},
	}, nil)

	// 自动关闭
	go func() {
		time.Sleep(duration)
		gui.g.Update(func(*gocui.Gui) error {
			gui.CloseMessageBox(mb, 0)
			return nil
		})
	}()
}
