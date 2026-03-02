# 进度条功能使用说明

## 功能概述

lazygit 现在支持美观的进度条提示框，用于显示长时间操作的进度，保持终端原生 TUI 风格。

## 功能特性

### 1. 两种进度模式

#### 确定进度（百分比显示）
```
╭─────────────────────────────────────────╮
│ ⏳ 正在推送到远程仓库...                 │
├─────────────────────────────────────────┤
│                                         │
│  [████████████░░░░░░░░] 60%            │
│                                         │
│  已完成: 12.5 MB / 20.8 MB              │
│  速度: 1.2 MB/s                         │
│  剩余时间: 约 7 秒                      │
│                                         │
╰─────────────────────────────────────────╯
```

#### 不确定进度（旋转动画）
```
╭─────────────────────────────────────────╮
│ ⏳ 正在克隆仓库...                       │
├─────────────────────────────────────────┤
│                                         │
│  ⠋ 正在连接服务器...                    │
│                                         │
╰─────────────────────────────────────────╯
```

### 2. 五种进度条样式

- **方块样式** (block): `[████████████░░░░░░░░] 60%`
- **点状样式** (dot): `[●●●●●●●●●●●●○○○○○○○○] 60%`
- **箭头样式** (arrow): `[►►►►►►►►►►►░░░░░░░░░] 60%`
- **渐变样式** (gradient): `[████████████▓▒░░░░░] 60%`
- **ASCII样式** (ascii): `[============>       ] 60%`

### 3. 五种旋转动画

- **Braille 点阵** (braille): `⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏`
- **线条旋转** (line): `| / - \`
- **箭头旋转** (arrow): `← ↖ ↑ ↗ → ↘ ↓ ↙`
- **点旋转** (dot): `⣾ ⣽ ⣻ ⢿ ⡿ ⣟ ⣯ ⣷`
- **圆圈旋转** (circle): `◐ ◓ ◑ ◒`

## 配置选项

在 `~/.config/lazygit/config.yml` 中添加配置：

```yaml
gui:
  progressBar:
    # 进度条样式: "block", "dot", "arrow", "gradient", "ascii"
    style: "block"

    # 旋转动画样式: "braille", "line", "arrow", "dot", "circle"
    spinnerStyle: "braille"

    # 进度条宽度（字符数）
    width: 30

    # 是否显示百分比
    showPercentage: true

    # 是否显示统计信息（速度、剩余时间等）
    showStats: true
```

## 代码使用示例

### 示例 1: 确定进度条

```go
// 显示确定进度的进度条
pb := gui.ShowProgressBar(ProgressBarConfig{
    Title:          "正在推送到远程仓库...",
    Total:          20 * 1024 * 1024, // 20 MB
    Width:          30,
    ShowPercentage: true,
    ShowStats:      true,
    Style:          ProgressBarStyleBlock,
    Indeterminate:  false,
})

// 在后台协程中更新进度
go func() {
    for i := int64(0); i <= pb.config.Total; i += 512 * 1024 {
        pb.Update(i, "")
        time.Sleep(200 * time.Millisecond)
    }
    pb.Close() // 完成后关闭
}()
```

### 示例 2: 不确定进度条

```go
// 显示不确定进度的进度条（旋转动画）
pb := gui.ShowProgressBar(ProgressBarConfig{
    Title:         "正在克隆仓库...",
    Message:       "正在连接服务器...",
    Indeterminate: true,
    SpinnerStyle:  SpinnerStyleBraille,
})

// 在后台协程中更新消息
go func() {
    time.Sleep(2 * time.Second)
    pb.Update(0, "正在接收对象...")

    time.Sleep(2 * time.Second)
    pb.Update(0, "正在解析增量...")

    time.Sleep(2 * time.Second)
    pb.Close()
}()
```

### 示例 3: 动态切换模式

```go
// 开始时显示不确定进度
pb := gui.ShowProgressBar(ProgressBarConfig{
    Title:         "正在推送...",
    Message:       "正在连接...",
    Indeterminate: true,
    SpinnerStyle:  SpinnerStyleBraille,
})

go func() {
    time.Sleep(1 * time.Second)

    // 连接成功后切换到确定进度
    pb.config.Indeterminate = false
    pb.config.Total = 10 * 1024 * 1024
    pb.config.ShowStats = true

    // 更新进度
    for i := int64(0); i <= pb.config.Total; i += 256 * 1024 {
        pb.Update(i, "")
        time.Sleep(100 * time.Millisecond)
    }

    pb.Close()
}()
```

## API 参考

### ProgressBarConfig 结构

```go
type ProgressBarConfig struct {
    Title          string           // 标题
    Message        string           // 消息（用于不确定进度）
    Total          int64            // 总量（字节数）
    Current        int64            // 当前值
    Width          int              // 进度条宽度
    ShowPercentage bool             // 显示百分比
    ShowStats      bool             // 显示统计信息
    Style          ProgressBarStyle // 进度条样式
    SpinnerStyle   SpinnerStyle     // 旋转动画样式
    Indeterminate  bool             // 是否为不确定进度
}
```

### ProgressBar 方法

```go
// 更新进度
pb.Update(current int64, message string)

// 设置总量
pb.SetTotal(total int64)

// 关闭进度条
pb.Close()
```

## 测试功能

项目中包含了完整的测试示例，可以通过以下方式测试：

1. 在代码中调用测试菜单：
```go
gui.createProgressBarTestMenu()
```

2. 测试菜单包含以下选项：
   - 确定进度条（推送示例）
   - 不确定进度条（克隆示例）
   - 测试所有进度条样式
   - 测试所有旋转动画
   - 模拟 Git Push
   - 模拟 Git Clone
   - 模拟 Git Fetch

## 实际应用场景

### Git Push
```go
func (gui *Gui) pushWithProgress() error {
    pb := gui.ShowProgressBar(ProgressBarConfig{
        Title:         "正在推送到远程仓库...",
        Indeterminate: true,
        SpinnerStyle:  SpinnerStyleBraille,
    })

    // 执行实际的 git push 操作
    // 根据 git 输出更新进度

    return nil
}
```

### Git Clone
```go
func (gui *Gui) cloneWithProgress(url string) error {
    pb := gui.ShowProgressBar(ProgressBarConfig{
        Title:         "正在克隆仓库...",
        Message:       "正在连接服务器...",
        Indeterminate: true,
        SpinnerStyle:  SpinnerStyleBraille,
    })

    // 执行实际的 git clone 操作
    // 根据 git 输出更新消息

    return nil
}
```

### Git Fetch
```go
func (gui *Gui) fetchWithProgress() error {
    pb := gui.ShowProgressBar(ProgressBarConfig{
        Title:          "正在获取更新...",
        Total:          5 * 1024 * 1024,
        ShowPercentage: true,
        ShowStats:      true,
        Style:          ProgressBarStyleGradient,
    })

    // 执行实际的 git fetch 操作
    // 根据 git 输出更新进度

    return nil
}
```

## 文件结构

```
pkg/gui/
├── progress_bar.go              # 进度条核心实现
├── progress_bar_examples.go     # 使用示例
└── progress_bar_test_menu.go    # 测试菜单
```

## 注意事项

1. **线程安全**：进度条的更新会自动在 GUI 线程中执行，可以安全地从后台协程调用 `Update()` 方法。

2. **资源清理**：使用完进度条后务必调用 `Close()` 方法，否则进度条会一直显示。

3. **性能考虑**：进度条每 100ms 更新一次，不要过于频繁地调用 `Update()` 方法。

4. **样式选择**：
   - 如果终端不支持 Unicode，使用 `ascii` 样式
   - 如果需要更好的视觉效果，使用 `gradient` 样式
   - 默认的 `block` 样式适合大多数场景

5. **统计信息**：只有在确定进度模式下才会显示统计信息（速度、剩余时间等）。

## 未来改进

- [ ] 支持多阶段进度条
- [ ] 支持可取消操作
- [ ] 支持进度条颜色渐变
- [ ] 集成到实际的 Git 操作中
- [ ] 添加进度条动画效果

## 相关文档

- [PROGRESSBAR_DESIGN.md](./PROGRESSBAR_DESIGN.md) - 详细设计文档
- [MESSAGEBOX_DESIGN.md](./MESSAGEBOX_DESIGN.md) - MessageBox 设计文档

---

**实现完成**: 2024
**状态**: ✅ 已实现并可用
