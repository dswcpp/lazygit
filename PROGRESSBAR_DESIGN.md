# 缓冲条提示框设计与实现

## 设计目标

为 lazygit 添加一个美观的进度条提示框，用于显示长时间操作的进度，保持终端原生 TUI 风格。

## 视觉设计

### 1. 确定进度条（百分比）

```
╭─────────────────────────────────────────╮
│ ⏳ 正在推送到远程仓库...                 │
├─────────────────────────────────────────┤
│                                         │
│  [████████████░░░░░░░░] 60%            │
│                                         │
│  已推送: 12.5 MB / 20.8 MB              │
│  速度: 1.2 MB/s                         │
│  剩余时间: 约 7 秒                      │
│                                         │
╰─────────────────────────────────────────╯
```

### 2. 不确定进度条（旋转动画）

```
╭─────────────────────────────────────────╮
│ ⏳ 正在克隆仓库...                       │
├─────────────────────────────────────────┤
│                                         │
│  ⠋ 正在连接服务器...                    │
│                                         │
│  已接收: 156 个对象                     │
│  已解析: 89 个增量                      │
│                                         │
╰─────────────────────────────────────────╯
```

### 3. 多阶段进度条

```
╭─────────────────────────────────────────╮
│ ⏳ 正在构建项目...                       │
├─────────────────────────────────────────┤
│                                         │
│  ✓ 1. 安装依赖                          │
│  ✓ 2. 编译代码                          │
│  ⠋ 3. 运行测试                          │
│  ░ 4. 打包发布                          │
│                                         │
│  [████████░░░░░░░░░░] 40%              │
│                                         │
╰─────────────────────────────────────────╯
```

### 4. 简洁进度条

```
╭─────────────────────────────────────────╮
│ ⏳ 正在拉取更新...                       │
│                                         │
│  [████████████████████] 100%           │
╰─────────────────────────────────────────╯
```

## 进度条样式

### 样式 1: 方块填充
```
[████████████░░░░░░░░] 60%
```

### 样式 2: 点状填充
```
[●●●●●●●●●●●●○○○○○○○○] 60%
```

### 样式 3: 箭头填充
```
[►►►►►►►►►►►░░░░░░░░░] 60%
```

### 样式 4: 渐变填充
```
[█▓▒░░░░░░░░░░░░░░░░] 20%
[████████▓▒░░░░░░░░░] 50%
[████████████████▓▒░] 90%
```

### 样式 5: ASCII 兼容
```
[============>       ] 60%
```

## 旋转动画字符

### 样式 1: Braille 点阵（推荐）
```
⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏
```

### 样式 2: 线条旋转
```
| / - \
```

### 样式 3: 箭头旋转
```
← ↖ ↑ ↗ → ↘ ↓ ↙
```

### 样式 4: 点旋转
```
⣾ ⣽ ⣻ ⢿ ⡿ ⣟ ⣯ ⣷
```

### 样式 5: 圆圈旋转
```
◐ ◓ ◑ ◒
```

## 代码实现

### 1. 进度条结构定义

```go
package gui

import (
    "fmt"
    "strings"
    "time"
)

// ProgressBarStyle 进度条样式
type ProgressBarStyle int

const (
    ProgressBarStyleBlock ProgressBarStyle = iota  // 方块填充
    ProgressBarStyleDot                            // 点状填充
    ProgressBarStyleArrow                          // 箭头填充
    ProgressBarStyleGradient                       // 渐变填充
    ProgressBarStyleASCII                          // ASCII 兼容
)

// SpinnerStyle 旋转动画样式
type SpinnerStyle int

const (
    SpinnerStyleBraille SpinnerStyle = iota  // Braille 点阵
    SpinnerStyleLine                          // 线条旋转
    SpinnerStyleArrow                         // 箭头旋转
    SpinnerStyleDot                           // 点旋转
    SpinnerStyleCircle                        // 圆圈旋转
)

// ProgressBarConfig 进度条配置
type ProgressBarConfig struct {
    Title           string            // 标题
    Message         string            // 消息
    Total           int64             // 总量
    Current         int64             // 当前值
    Width           int               // 进度条宽度
    ShowPercentage  bool              // 显示百分比
    ShowStats       bool              // 显示统计信息
    Style           ProgressBarStyle  // 进度条样式
    SpinnerStyle    SpinnerStyle      // 旋转动画样式
    Indeterminate   bool              // 是否为不确定进度
}

// ProgressBar 进度条
type ProgressBar struct {
    config      ProgressBarConfig
    startTime   time.Time
    lastUpdate  time.Time
    spinnerIdx  int
    view        *gocui.View
}

// NewProgressBar 创建进度条
func NewProgressBar(config ProgressBarConfig) *ProgressBar {
    if config.Width == 0 {
        config.Width = 30
    }
    return &ProgressBar{
        config:     config,
        startTime:  time.Now(),
        lastUpdate: time.Now(),
        spinnerIdx: 0,
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

// Render 渲染进度条
func (pb *ProgressBar) Render() string {
    var lines []string

    // 标题行
    icon := "⏳"
    title := fmt.Sprintf("%s %s", icon, pb.config.Title)
    lines = append(lines, title)

    // 空行
    lines = append(lines, "")

    // 进度条
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
    if pb.config.ShowStats && !pb.config.Indeterminate {
        lines = append(lines, "")
        lines = append(lines, pb.renderStats())
    }

    // 空行
    lines = append(lines, "")

    return strings.Join(lines, "\n")
}

// renderProgressBar 渲染进度条
func (pb *ProgressBar) renderProgressBar() string {
    percentage := float64(pb.config.Current) / float64(pb.config.Total) * 100
    filled := int(float64(pb.config.Width) * float64(pb.config.Current) / float64(pb.config.Total))

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
    if filled < pb.config.Width {
        // 添加渐变字符
        bar += "▓"
        empty--
        if empty > 0 {
            bar += strings.Repeat("░", empty)
        }
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
    speed := float64(pb.config.Current) / elapsed.Seconds()

    // 计算剩余时间
    remaining := time.Duration(0)
    if speed > 0 {
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

    return spinner
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
```

### 2. GUI 集成

```go
// ShowProgressBar 显示进度条
func (gui *Gui) ShowProgressBar(config ProgressBarConfig) *ProgressBar {
    pb := NewProgressBar(config)

    // 创建弹出窗口
    gui.createProgressBarPopup(pb)

    // 启动更新协程
    go gui.updateProgressBar(pb)

    return pb
}

// createProgressBarPopup 创建进度条弹出窗口
func (gui *Gui) createProgressBarPopup(pb *ProgressBar) {
    width := 50
    height := 10

    x0 := (gui.g.MaxX() - width) / 2
    y0 := (gui.g.MaxY() - height) / 2
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

    for range ticker.C {
        if pb.view == nil {
            return
        }

        pb.view.Clear()
        content := pb.Render()
        fmt.Fprint(pb.view, content)

        gui.g.Update(func(*gocui.Gui) error {
            return nil
        })

        // 检查是否完成
        if !pb.config.Indeterminate && pb.config.Current >= pb.config.Total {
            time.Sleep(500 * time.Millisecond)
            gui.CloseProgressBar()
            return
        }
    }
}

// CloseProgressBar 关闭进度条
func (gui *Gui) CloseProgressBar() {
    gui.g.DeleteView("progressBar")
}
```

### 3. 使用示例

```go
// 示例 1: 确定进度
func (gui *Gui) pushWithProgress() error {
    pb := gui.ShowProgressBar(ProgressBarConfig{
        Title:          "正在推送到远程仓库...",
        Total:          20 * 1024 * 1024, // 20 MB
        Width:          30,
        ShowPercentage: true,
        ShowStats:      true,
        Style:          ProgressBarStyleBlock,
    })

    // 模拟推送过程
    for i := int64(0); i <= pb.config.Total; i += 1024 * 1024 {
        pb.Update(i, "")
        time.Sleep(500 * time.Millisecond)
    }

    return nil
}

// 示例 2: 不确定进度
func (gui *Gui) cloneWithProgress() error {
    pb := gui.ShowProgressBar(ProgressBarConfig{
        Title:         "正在克隆仓库...",
        Message:       "正在连接服务器...",
        Indeterminate: true,
        SpinnerStyle:  SpinnerStyleBraille,
    })

    // 模拟克隆过程
    time.Sleep(5 * time.Second)
    pb.Update(0, "正在接收对象...")
    time.Sleep(3 * time.Second)

    gui.CloseProgressBar()
    return nil
}

// 示例 3: 多阶段进度
func (gui *Gui) buildWithProgress() error {
    stages := []string{
        "安装依赖",
        "编译代码",
        "运行测试",
        "打包发布",
    }

    pb := gui.ShowProgressBar(ProgressBarConfig{
        Title:          "正在构建项目...",
        Total:          int64(len(stages)),
        Width:          30,
        ShowPercentage: true,
        Style:          ProgressBarStyleBlock,
    })

    for i, stage := range stages {
        pb.Update(int64(i+1), fmt.Sprintf("✓ %d. %s", i+1, stage))
        time.Sleep(2 * time.Second)
    }

    return nil
}
```

## 配置选项

在 `config.yml` 中添加：

```yaml
gui:
  progressBar:
    # 进度条样式: "block", "dot", "arrow", "gradient", "ascii"
    style: "block"

    # 旋转动画样式: "braille", "line", "arrow", "dot", "circle"
    spinnerStyle: "braille"

    # 进度条宽度
    width: 30

    # 是否显示百分比
    showPercentage: true

    # 是否显示统计信息
    showStats: true
```

## 实现优先级

### Phase 1: 基础功能
- ✅ 基本进度条渲染
- ✅ 确定进度显示
- ✅ 不确定进度（旋转动画）

### Phase 2: 样式优化
- ✅ 多种进度条样式
- ✅ 多种旋转动画
- ✅ 统计信息显示

### Phase 3: 高级功能
- ⬜ 多阶段进度
- ⬜ 可取消操作
- ⬜ 进度条颜色渐变

## 效果预览

### 推送操作
```
╭─────────────────────────────────────────╮
│ ⏳ 正在推送到远程仓库...                 │
├─────────────────────────────────────────┤
│                                         │
│  [████████████████░░░░] 80%            │
│                                         │
│  已完成: 16.6 MB / 20.8 MB              │
│  速度: 2.1 MB/s                         │
│  剩余时间: 约 2 秒                      │
│                                         │
╰─────────────────────────────────────────╯
```

### 克隆操作
```
╭─────────────────────────────────────────╮
│ ⏳ 正在克隆仓库...                       │
├─────────────────────────────────────────┤
│                                         │
│  ⠹ 正在接收对象...                      │
│                                         │
│  已接收: 1,234 个对象                   │
│  已解析: 567 个增量                     │
│                                         │
╰─────────────────────────────────────────╯
```

---

**设计完成**: 2024
**状态**: 待实现
