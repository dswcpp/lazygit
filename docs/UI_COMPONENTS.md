# UI 组件使用指南

## 📖 目录

- [概述](#概述)
- [消息框系统](#消息框系统)
- [进度条系统](#进度条系统)
- [活动栏](#活动栏)
- [使用示例](#使用示例)
- [最佳实践](#最佳实践)

---

## 概述

Lazygit 增强版提供了三个精心设计的 UI 组件，保持终端原生 TUI 风格的同时，提供现代化的交互体验。

```
┌──────────────────────────────────────────────────┐
│  📬 消息框   │  📊 进度条   │  📱 活动栏      │
│  提示确认    │  进度显示    │  快速导航       │
└──────────────────────────────────────────────────┘
```

---

## 消息框系统

### 功能特性

#### 5 种消息类型

```
ℹ️  Info     - 信息提示（蓝色）
✅ Success  - 成功反馈（绿色）
⚠️  Warning  - 警告提示（黄色）
❌ Error    - 错误提示（红色）
❓ Question - 询问确认（青色）
```

#### 核心功能

- 图标和颜色编码
- 自定义按钮配置
- 键盘导航支持
- 自动换行显示
- 详细信息展开
- 自动关闭功能

### 快捷键

| 快捷键 | 功能 |
|--------|------|
| **Enter** | 确认当前按钮 |
| **Esc** | 取消（最后一个按钮） |
| **←/→** | 切换按钮 |
| **Tab** | 循环切换按钮 |
| **1-9** | 快速选择按钮 |

### 便捷方法

#### 1. ShowError - 错误提示

```go
gui.ShowError(
    "推送失败",
    "无法连接到远程仓库",
    "错误详情: ECONNREFUSED",
)
```

**效果**:
```
┌─ ❌ 推送失败 ────────────────────────┐
│                                      │
│  无法连接到远程仓库                  │
│                                      │
│  详细信息:                           │
│  错误详情: ECONNREFUSED              │
│                                      │
│         [ 确定 ]                     │
└──────────────────────────────────────┘
```

#### 2. ShowWarning - 警告提示

```go
gui.ShowWarning(
    "未保存的更改",
    "当前有未提交的更改",
    "切换分支前请先提交或暂存",
)
```

#### 3. ShowInfo - 信息提示

```go
gui.ShowInfo(
    "操作提示",
    "按 'h' 查看帮助",
    "更多快捷键请参考文档",
)
```

#### 4. ShowSuccess - 成功反馈

```go
gui.ShowSuccess(
    "推送成功",
    "已推送 3 个提交到 origin/main",
    "",
)
```

#### 5. ShowConfirm - 确认对话框

```go
gui.ShowConfirm(
    "确认删除",
    "确定要删除分支 'feature/old' 吗？",
    func() {
        // 用户点击确认后执行
        gui.deleteBranch("feature/old")
    },
)
```

**效果**:
```
┌─ ❓ 确认删除 ────────────────────────┐
│                                      │
│  确定要删除分支 'feature/old' 吗？   │
│                                      │
│      [ 确定 ]    [ 取消 ]            │
└──────────────────────────────────────┘
```

#### 6. ShowYesNoCancel - 三选项对话框

```go
gui.ShowYesNoCancel(
    "未保存的更改",
    "如何处理未提交的更改？",
    func() {
        // 用户选择 "是"
        gui.stashChanges()
    },
    func() {
        // 用户选择 "否"
        gui.discardChanges()
    },
)
```

#### 7. ShowMessageBox - 自定义消息框

```go
gui.ShowMessageBox(MessageBoxConfig{
    Type:    MessageTypeQuestion,
    Title:   "选择操作",
    Message: "如何处理冲突？",
    Details: "当前有 3 个文件存在冲突",
    Buttons: []string{"手动解决", "使用我们的", "使用他们的", "取消"},
}, func(buttonIndex int) {
    switch buttonIndex {
    case 0:
        gui.manualResolve()
    case 1:
        gui.resolveWithOurs()
    case 2:
        gui.resolveWithTheirs()
    case 3:
        // 取消
    }
})
```

#### 8. ShowAutoCloseMessage - 自动关闭消息

```go
gui.ShowAutoCloseMessage(
    MessageTypeSuccess,
    "操作成功",
    "文件已保存",
    3 * time.Second,  // 3 秒后自动关闭
)
```

### 使用场景

#### 场景 1: Git 操作错误

```go
func (gui *Gui) handlePushError(err error) {
    gui.ShowError(
        "推送失败",
        "无法推送到远程仓库",
        err.Error(),
    )
}
```

#### 场景 2: 危险操作确认

```go
func (gui *Gui) confirmForceDelete() {
    gui.ShowWarning(
        "危险操作",
        "强制删除分支将丢失所有未合并的提交",
        "此操作不可撤销",
    )
}
```

#### 场景 3: 操作成功反馈

```go
func (gui *Gui) afterCommit() {
    gui.ShowAutoCloseMessage(
        MessageTypeSuccess,
        "提交成功",
        "已创建提交 abc1234",
        2 * time.Second,
    )
}
```

---

## 进度条系统

### 功能特性

#### 两种进度模式

**1. 确定进度** - 显示百分比和进度条
```
┌─ 正在推送... ────────────────────────┐
│                                      │
│  ████████████░░░░░░░░░░░░░░░  45%   │
│                                      │
│  已完成: 9.2 MB / 20.5 MB            │
│  速度: 1.5 MB/s                      │
│  剩余时间: 7 秒                      │
│                                      │
└──────────────────────────────────────┘
```

**2. 不确定进度** - 显示旋转动画
```
┌─ 正在克隆... ────────────────────────┐
│                                      │
│  ⠋ 正在处理...                       │
│                                      │
│  已用时间: 15 秒                     │
│                                      │
└──────────────────────────────────────┘
```

#### 5 种进度条样式

```
block:    ████████████░░░░░░░░░░░░░░░
dot:      ●●●●●●●●●●●●○○○○○○○○○○○○○○
arrow:    >>>>>>>>>>>>>>>------------
gradient: ▓▓▓▓▓▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒░░░░░░
ascii:    [============>          ]
```

#### 5 种旋转动画

```
braille:  ⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏
line:     | / - \
arrow:    ← ↖ ↑ ↗ → ↘ ↓ ↙
dot:      ⣾ ⣽ ⣻ ⢿ ⡿ ⣟ ⣯ ⣷
circle:   ◐ ◓ ◑ ◒
```

### 基本使用

#### 1. 确定进度

```go
// 创建进度条
pb := gui.ShowProgressBar(ProgressBarConfig{
    Title:          "正在推送...",
    Total:          20 * 1024 * 1024,  // 20 MB
    ShowPercentage: true,
    ShowStats:      true,
})

// 更新进度
go func() {
    for i := int64(0); i <= pb.config.Total; i += 512 * 1024 {
        pb.Update(i, "")
        time.Sleep(200 * time.Millisecond)
    }
    pb.Close()
}()
```

#### 2. 不确定进度

```go
pb := gui.ShowProgressBar(ProgressBarConfig{
    Title:         "正在克隆...",
    Indeterminate: true,
    Message:       "正在下载对象...",
})

go func() {
    // 执行长时间操作
    err := git.Clone(url)
    pb.Close()

    if err != nil {
        gui.g.Update(func(*gocui.Gui) error {
            gui.ShowError("克隆失败", err.Error())
            return nil
        })
    }
}()
```

#### 3. 动态切换模式

```go
pb := gui.ShowProgressBar(ProgressBarConfig{
    Title:         "正在处理...",
    Indeterminate: true,
})

go func() {
    // 第一阶段：不确定进度
    pb.Update(0, "正在准备...")
    time.Sleep(2 * time.Second)

    // 切换到确定进度
    pb.SwitchToDeterminate(100)

    // 第二阶段：确定进度
    for i := int64(0); i <= 100; i++ {
        pb.Update(i, fmt.Sprintf("处理中 %d/100", i))
        time.Sleep(50 * time.Millisecond)
    }

    pb.Close()
}()
```

### 配置选项

在 `~/.config/lazygit/config.yml` 中配置：

```yaml
gui:
  progressBar:
    style: "block"           # 样式: block, dot, arrow, gradient, ascii
    spinnerStyle: "braille"  # 动画: braille, line, arrow, dot, circle
    width: 30                # 宽度
    showPercentage: true     # 显示百分比
    showStats: true          # 显示统计信息
```

### 使用场景

#### 场景 1: Git Push

```go
func (gui *Gui) pushWithProgress() error {
    pb := gui.ShowProgressBar(ProgressBarConfig{
        Title:          "正在推送...",
        Total:          totalSize,
        ShowPercentage: true,
        ShowStats:      true,
    })

    go func() {
        err := gui.git.Push(func(current, total int64) {
            pb.Update(current, "")
        })
        pb.Close()

        gui.g.Update(func(*gocui.Gui) error {
            if err != nil {
                gui.ShowError("推送失败", err.Error())
            } else {
                gui.ShowSuccess("推送成功", "")
            }
            return nil
        })
    }()

    return nil
}
```

#### 场景 2: Git Clone

```go
func (gui *Gui) cloneWithProgress(url string) error {
    pb := gui.ShowProgressBar(ProgressBarConfig{
        Title:         "正在克隆...",
        Indeterminate: true,
    })

    go func() {
        err := gui.git.Clone(url, func(stage string) {
            pb.Update(0, stage)
        })
        pb.Close()

        gui.g.Update(func(*gocui.Gui) error {
            if err != nil {
                gui.ShowError("克隆失败", err.Error())
            } else {
                gui.ShowSuccess("克隆成功", "")
            }
            return nil
        })
    }()

    return nil
}
```

#### 场景 3: 大文件操作

```go
func (gui *Gui) processLargeFile(path string) error {
    fileInfo, _ := os.Stat(path)
    total := fileInfo.Size()

    pb := gui.ShowProgressBar(ProgressBarConfig{
        Title:          "正在处理文件...",
        Total:          total,
        ShowPercentage: true,
        ShowStats:      true,
    })

    go func() {
        // 处理文件
        processed := int64(0)
        for processed < total {
            // ... 处理逻辑
            processed += chunkSize
            pb.Update(processed, "")
        }
        pb.Close()
    }()

    return nil
}
```

---

## 活动栏

### 功能特性

- VSCode 风格侧边栏
- 快速导航
- 状态指示
- 可自定义布局
- 图标支持

### 使用方法

```go
// 显示活动栏
gui.ShowActivityBar()

// 隐藏活动栏
gui.HideActivityBar()

// 切换活动栏
gui.ToggleActivityBar()
```

### 配置选项

```yaml
gui:
  activityBar:
    enabled: true
    position: "left"  # left, right
    width: 5
    icons:
      files: "📁"
      branches: "🌿"
      commits: "📝"
      stash: "📦"
      tags: "🏷️"
```

---

## 使用示例

### 示例 1: 标准操作流程

```go
func (gui *Gui) standardOperation() error {
    // 1. 确认操作
    gui.ShowConfirm("确认操作", "确定要执行吗？", func() {
        // 2. 显示进度
        pb := gui.ShowProgressBar(ProgressBarConfig{
            Title: "正在处理...",
            Total: total,
        })

        // 3. 执行操作
        go func() {
            err := doSomething(func(current int64) {
                pb.Update(current, "")
            })
            pb.Close()

            // 4. 显示结果
            gui.g.Update(func(*gocui.Gui) error {
                if err != nil {
                    gui.ShowError("操作失败", err.Error())
                } else {
                    gui.ShowSuccess("操作成功", "")
                }
                return nil
            })
        }()
    })
    return nil
}
```

### 示例 2: 多选项操作

```go
func (gui *Gui) multiOptionOperation() error {
    gui.ShowMessageBox(MessageBoxConfig{
        Type:    MessageTypeQuestion,
        Title:   "选择操作",
        Message: "如何处理未提交的更改？",
        Buttons: []string{"暂存", "丢弃", "取消"},
    }, func(buttonIndex int) {
        switch buttonIndex {
        case 0:
            gui.stashChanges()
        case 1:
            gui.discardChanges()
        case 2:
            // 取消
        }
    })
    return nil
}
```

### 示例 3: 错误处理

```go
func (gui *Gui) handleError(err error) {
    if err == nil {
        return
    }

    // 根据错误类型显示不同的消息
    switch {
    case errors.Is(err, ErrNetworkTimeout):
        gui.ShowWarning(
            "网络超时",
            "连接超时，请检查网络",
            err.Error(),
        )
    case errors.Is(err, ErrPermissionDenied):
        gui.ShowError(
            "权限不足",
            "没有权限执行此操作",
            err.Error(),
        )
    default:
        gui.ShowError(
            "操作失败",
            "发生未知错误",
            err.Error(),
        )
    }
}
```

---

## 最佳实践

### 1. 消息框使用原则

✅ **推荐**:
- 重要操作前使用确认对话框
- 错误信息使用 ShowError
- 成功反馈使用 ShowSuccess
- 危险操作使用 Warning 类型
- 简短消息使用自动关闭

❌ **避免**:
- 过度使用消息框打断用户
- 不重要的信息使用 Toast 即可
- 消息内容过长（使用 Details 字段）
- 连续弹出多个消息框

### 2. 进度条使用原则

✅ **推荐**:
- 超过 2 秒的操作显示进度条
- 能计算进度的使用确定进度
- 不能计算的使用不确定进度
- 操作完成后立即关闭
- 在 goroutine 中更新进度

❌ **避免**:
- 短时间操作（< 1 秒）显示进度条
- 忘记调用 Close() 关闭进度条
- 在主线程中更新进度
- 频繁创建和销毁进度条

### 3. 活动栏使用原则

✅ **推荐**:
- 保持图标简洁明了
- 使用一致的图标风格
- 提供快捷键切换
- 根据屏幕大小调整

❌ **避免**:
- 图标过多导致混乱
- 图标含义不明确
- 占用过多屏幕空间

### 4. 通用原则

✅ **推荐**:
- 保持 UI 一致性
- 提供键盘快捷键
- 考虑用户体验
- 测试各种场景

❌ **避免**:
- 阻塞主线程
- 忽略错误处理
- 过度使用动画
- 忽视性能影响

---

## 性能优化

### 1. 消息框

```go
// 避免频繁创建
// 使用消息队列
type MessageQueue struct {
    messages chan Message
}

func (mq *MessageQueue) Add(msg Message) {
    select {
    case mq.messages <- msg:
    default:
        // 队列满，丢弃旧消息
    }
}
```

### 2. 进度条

```go
// 限制更新频率
type ThrottledProgressBar struct {
    pb           *ProgressBar
    lastUpdate   time.Time
    updateInterval time.Duration
}

func (tpb *ThrottledProgressBar) Update(current int64) {
    now := time.Now()
    if now.Sub(tpb.lastUpdate) < tpb.updateInterval {
        return
    }
    tpb.pb.Update(current, "")
    tpb.lastUpdate = now
}
```

### 3. 活动栏

```go
// 懒加载图标
func (ab *ActivityBar) loadIcon(name string) string {
    if icon, ok := ab.iconCache[name]; ok {
        return icon
    }
    icon := ab.loadIconFromConfig(name)
    ab.iconCache[name] = icon
    return icon
}
```

---

## 故障排除

### 问题 1: 消息框不显示

**原因**: 视图创建失败

**解决**:
```go
if gui.g == nil {
    return errors.New("GUI not initialized")
}
```

### 问题 2: 进度条卡住

**原因**: 忘记调用 Close()

**解决**:
```go
defer pb.Close()
```

### 问题 3: 界面闪烁

**原因**: 更新频率过高

**解决**:
```go
// 限制更新频率到 100ms
time.Sleep(100 * time.Millisecond)
```

---

## 相关文档

- [项目概述](./PROJECT_OVERVIEW.md)
- [AI 功能指南](./AI_FEATURES.md)
- [开发指南](./DEVELOPMENT_GUIDE.md)
- [消息框详细文档](../MESSAGEBOX_USAGE.md)
- [进度条详细文档](../PROGRESSBAR_USAGE.md)

---

**版本**: v1.0.0
**最后更新**: 2024
**状态**: ✅ 完整
