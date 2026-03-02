package gui

import (
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/gui/style"
)

// ProgressBarStyle 进度条样式
type ProgressBarStyle int

const (
	ProgressBarStyleBlock ProgressBarStyle = iota // 方块填充
	ProgressBarStyleDot                           // 点状填充
	ProgressBarStyleArrow                         // 箭头填充
	ProgressBarStyleGradient                      // 渐变填充
	ProgressBarStyleASCII                         // ASCII 兼容
)

// SpinnerStyle 旋转动画样式
type SpinnerStyle int

const (
	SpinnerStyleBraille SpinnerStyle = iota // Braille 点阵
	SpinnerStyleLine                         // 线条旋转
	SpinnerStyleArrow                        // 箭头旋转
	SpinnerStyleDot                          // 点旋转
	SpinnerStyleCircle                       // 圆圈旋转
)

// ProgressBarConfig 进度条配置
type ProgressBarConfig struct {
	Title          string           // 标题
	Message        string           // 消息
	Total          int64            // 总量
	Current        int64            // 当前值
	Width          int              // 进度条宽度
	ShowPercentage bool             // 显示百分比
	ShowStats      bool             // 显示统计信息
	Style          ProgressBarStyle // 进度条样式
	SpinnerStyle   SpinnerStyle     // 旋转动画样式
	Indeterminate  bool             // 是否为不确定进度
}

// ProgressBar 进度条
type ProgressBar struct {
	config     ProgressBarConfig
	startTime  time.Time
	lastUpdate time.Time
	spinnerIdx int
	view       *gocui.View
	done       chan bool
	gui        *Gui
}

// NewProgressBar 创建进度条
func (gui *Gui) NewProgressBar(config ProgressBarConfig) *ProgressBar {
	if config.Width == 0 {
		config.Width = 30
	}
	if !config.ShowPercentage {
		config.ShowPercentage = true
	}
	return &ProgressBar{
		config:     config,
		startTime:  time.Now(),
		lastUpdate: time.Now(),
		spinnerIdx: 0,
		done:       make(chan bool),
		gui:        gui,
	}
}

// Update 更新进度
func (pb *ProgressBar) Update(current int64, message string) {
	pb.config.Current = current
	if message != "" {
		pb.config.Message = message
	}
	pb.lastUpdate = time.Now()
}

// SetTotal 设置总量
func (pb *ProgressBar) SetTotal(total int64) {
	pb.config.Total = total
}

// Close 关闭进度条
func (pb *ProgressBar) Close() {
	close(pb.done)
}

// Render 渲染进度条
func (pb *ProgressBar) Render() string {
	var lines []string

	// 标题行
	icon := "⏳"
	title := fmt.Sprintf(" %s %s", icon, pb.config.Title)
	lines = append(lines, title)

	// 空行
	lines = append(lines, "")

	// 进度条或旋转动画
	if pb.config.Indeterminate {
		// 不确定进度 - 显示旋转动画
		spinner := pb.getSpinner()
		lines = append(lines, fmt.Sprintf("  %s %s", spinner, pb.config.Message))
	} else {
		// 确定进度 - 显示进度条
		progressBar := pb.renderProgressBar()
		lines = append(lines, "  "+progressBar)
	}

	// 统计信息
	if pb.config.ShowStats && !pb.config.Indeterminate && pb.config.Total > 0 {
		lines = append(lines, "")
		lines = append(lines, pb.renderStats())
	}

	// 空行
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

// renderProgressBar 渲染进度条
func (pb *ProgressBar) renderProgressBar() string {
	if pb.config.Total == 0 {
		return "[" + strings.Repeat("░", pb.config.Width) + "] 0%"
	}

	percentage := float64(pb.config.Current) / float64(pb.config.Total) * 100
	if percentage > 100 {
		percentage = 100
	}

	filled := int(float64(pb.config.Width) * float64(pb.config.Current) / float64(pb.config.Total))
	if filled > pb.config.Width {
		filled = pb.config.Width
	}

	var bar string
	switch pb.config.Style {
	case ProgressBarStyleBlock:
		bar = pb.renderBlockBar(filled)
	case ProgressBarStyleDot:
		bar = pb.renderDotBar(filled)
	case ProgressBarStyleArrow:
		bar = pb.renderArrowBar(filled)
	case ProgressBarStyleGradient:
		bar = pb.renderGradientBar(filled)
	case ProgressBarStyleASCII:
		bar = pb.renderASCIIBar(filled)
	default:
		bar = pb.renderBlockBar(filled)
	}

	if pb.config.ShowPercentage {
		return fmt.Sprintf("[%s] %3.0f%%", bar, percentage)
	}
	return fmt.Sprintf("[%s]", bar)
}

// renderBlockBar 方块填充样式
func (pb *ProgressBar) renderBlockBar(filled int) string {
	empty := pb.config.Width - filled
	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}

// renderDotBar 点状填充样式
func (pb *ProgressBar) renderDotBar(filled int) string {
	empty := pb.config.Width - filled
	return strings.Repeat("●", filled) + strings.Repeat("○", empty)
}

// renderArrowBar 箭头填充样式
func (pb *ProgressBar) renderArrowBar(filled int) string {
	empty := pb.config.Width - filled
	return strings.Repeat("►", filled) + strings.Repeat("░", empty)
}

// renderGradientBar 渐变填充样式
func (pb *ProgressBar) renderGradientBar(filled int) string {
	empty := pb.config.Width - filled
	bar := strings.Repeat("█", filled)
	if filled < pb.config.Width && filled > 0 {
		// 添加渐变字符
		bar += "▓"
		empty--
		if empty > 0 {
			bar += strings.Repeat("░", empty)
		}
	} else if filled == 0 {
		bar = strings.Repeat("░", pb.config.Width)
	}
	return bar
}

// renderASCIIBar ASCII 兼容样式
func (pb *ProgressBar) renderASCIIBar(filled int) string {
	empty := pb.config.Width - filled
	bar := strings.Repeat("=", filled)
	if filled < pb.config.Width && filled > 0 {
		bar += ">"
		empty--
	}
	bar += strings.Repeat(" ", empty)
	return bar
}

// renderStats 渲染统计信息
func (pb *ProgressBar) renderStats() string {
	elapsed := time.Since(pb.startTime)

	// 计算速度
	speed := float64(0)
	if elapsed.Seconds() > 0 {
		speed = float64(pb.config.Current) / elapsed.Seconds()
	}

	// 计算剩余时间
	remaining := time.Duration(0)
	if speed > 0 && pb.config.Current < pb.config.Total {
		remainingBytes := pb.config.Total - pb.config.Current
		remaining = time.Duration(float64(remainingBytes)/speed) * time.Second
	}

	var lines []string

	// 已完成 / 总量
	lines = append(lines, fmt.Sprintf("  已完成: %s / %s",
		formatBytes(pb.config.Current),
		formatBytes(pb.config.Total)))

	// 速度
	if speed > 0 {
		lines = append(lines, fmt.Sprintf("  速度: %s/s", formatBytes(int64(speed))))
	}

	// 剩余时间
	if remaining > 0 {
		lines = append(lines, fmt.Sprintf("  剩余时间: 约 %s", formatDuration(remaining)))
	}

	return strings.Join(lines, "\n")
}

// getSpinner 获取旋转动画字符
func (pb *ProgressBar) getSpinner() string {
	var frames []string

	switch pb.config.SpinnerStyle {
	case SpinnerStyleBraille:
		frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	case SpinnerStyleLine:
		frames = []string{"|", "/", "-", "\\"}
	case SpinnerStyleArrow:
		frames = []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}
	case SpinnerStyleDot:
		frames = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	case SpinnerStyleCircle:
		frames = []string{"◐", "◓", "◑", "◒"}
	default:
		frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	}

	spinner := frames[pb.spinnerIdx%len(frames)]
	pb.spinnerIdx++

	return style.FgCyan.Sprint(spinner)
}

// formatBytes 格式化字节数
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatDuration 格式化时间
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%d 秒", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%d 分钟", int(d.Minutes()))
	}
	return fmt.Sprintf("%d 小时", int(d.Hours()))
}

// ShowProgressBar 显示进度条
func (gui *Gui) ShowProgressBar(config ProgressBarConfig) *ProgressBar {
	pb := gui.NewProgressBar(config)

	// 创建弹出窗口
	gui.createProgressBarPopup(pb)

	// 启动更新协程
	go gui.updateProgressBar(pb)

	return pb
}

// createProgressBarPopup 创建进度条弹出窗口
func (gui *Gui) createProgressBarPopup(pb *ProgressBar) {
	width := 50
	height := 12

	maxX, maxY := gui.g.Size()
	x0 := (maxX - width) / 2
	y0 := (maxY - height) / 2
	x1 := x0 + width
	y1 := y0 + height

	view, err := gui.g.SetView("progressBar", x0, y0, x1, y1, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return
	}

	view.Frame = true
	view.Title = ""
	pb.view = view

	gui.g.SetViewOnTop("progressBar")
}

// updateProgressBar 更新进度条
func (gui *Gui) updateProgressBar(pb *ProgressBar) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-pb.done:
			// 显示完成状态
			time.Sleep(500 * time.Millisecond)
			gui.CloseProgressBar()
			return
		case <-ticker.C:
			if pb.view == nil {
				return
			}

			gui.g.Update(func(*gocui.Gui) error {
				pb.view.Clear()
				content := pb.Render()
				fmt.Fprint(pb.view, content)
				return nil
			})

			// 检查是否完成
			if !pb.config.Indeterminate && pb.config.Total > 0 && pb.config.Current >= pb.config.Total {
				time.Sleep(500 * time.Millisecond)
				gui.CloseProgressBar()
				return
			}
		}
	}
}

// CloseProgressBar 关闭进度条
func (gui *Gui) CloseProgressBar() {
	gui.g.Update(func(*gocui.Gui) error {
		gui.g.DeleteView("progressBar")
		return nil
	})
}
