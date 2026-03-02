# 消息框功能使用说明

## 功能概述

lazygit 现在支持美观的消息框提示，用于显示错误、警告、信息、成功消息和确认对话框，保持终端原生 TUI 风格。

## 功能特性

### 1. 五种消息类型

#### 信息消息框 (Info)
```
╭──────────────────────────────────────────────────╮
│ ℹ 提示                                            │
│ ──────────────────────────────────────────────── │
│                                                  │
│  当前分支已经是最新的，无需拉取更新。              │
│                                                  │
│  最后更新时间: 2024-01-01 12:00:00               │
│                                                  │
│ ──────────────────────────────────────────────── │
│              [ 确定 ]                            │
╰──────────────────────────────────────────────────╯
```

#### 成功消息框 (Success)
```
╭──────────────────────────────────────────────────╮
│ ✓ 操作成功                                        │
│ ──────────────────────────────────────────────── │
│                                                  │
│  分支 'feature-x' 已成功删除。                    │
│                                                  │
│  本地分支和远程分支都已删除。                      │
│                                                  │
│ ──────────────────────────────────────────────── │
│              [ 确定 ]                            │
╰──────────────────────────────────────────────────╯
```

#### 警告消息框 (Warning)
```
╭──────────────────────────────────────────────────╮
│ ⚠ 警告                                           │
│ ──────────────────────────────────────────────── │
│                                                  │
│  你即将强制推送到远程分支 'main'，这将覆盖远程    │
│  仓库的历史记录！                                 │
│                                                  │
│  此操作不可撤销，请谨慎操作。                      │
│                                                  │
│ ──────────────────────────────────────────────── │
│              [ 确定 ]                            │
╰──────────────────────────────────────────────────╯
```

#### 错误消息框 (Error)
```
╭──────────────────────────────────────────────────╮
│ ✗ 操作失败                                        │
│ ──────────────────────────────────────────────── │
│                                                  │
│  无法连接到远程仓库，请检查网络连接。              │
│                                                  │
│  错误代码: ECONNREFUSED                          │
│  主机: github.com                                │
│  端口: 443                                       │
│                                                  │
│ ──────────────────────────────────────────────── │
│              [ 确定 ]                            │
╰──────────────────────────────────────────────────╯
```

#### 问题对话框 (Question)
```
╭──────────────────────────────────────────────────╮
│ ? 确认删除                                        │
│ ──────────────────────────────────────────────── │
│                                                  │
│  确定要删除分支 'feature-x' 吗？此操作不可撤销。   │
│                                                  │
│ ──────────────────────────────────────────────── │
│          [ 确定 ]  [ 取消 ]                      │
╰──────────────────────────────────────────────────╯
```

### 2. 消息类型图标和颜色

| 类型 | 图标 | 颜色 | 用途 |
|------|------|------|------|
| Info | ℹ | 青色 | 一般信息提示 |
| Success | ✓ | 绿色 | 操作成功提示 |
| Warning | ⚠ | 黄色 | 警告信息 |
| Error | ✗ | 红色 | 错误信息 |
| Question | ? | 青色 | 确认对话框 |

### 3. 按钮导航

- **左/右箭头键**: 切换按钮
- **Tab 键**: 循环切换按钮
- **Enter 键**: 确认当前选中的按钮
- **Esc 键**: 取消（选择最后一个按钮）
- **数字键 1-9**: 快速选择对应按钮

### 4. 特殊功能

- **自动换行**: 长文本自动换行显示
- **详细信息**: 支持显示额外的详细信息
- **自定义按钮**: 支持自定义按钮文本和数量
- **自动关闭**: 支持定时自动关闭

## 代码使用示例

### 示例 1: 显示错误消息

```go
gui.ShowError(
    "推送失败",
    "无法推送到远程仓库。",
    "错误代码: ECONNREFUSED\n主机: github.com\n端口: 443",
)
```

### 示例 2: 显示警告消息

```go
gui.ShowWarning(
    "警告",
    "你即将强制推送到远程分支 'main'，这将覆盖远程仓库的历史记录！",
    "此操作不可撤销，请谨慎操作。",
)
```

### 示例 3: 显示信息消息

```go
gui.ShowInfo(
    "提示",
    "当前分支已经是最新的，无需拉取更新。",
    "最后更新时间: 2024-01-01 12:00:00",
)
```

### 示例 4: 显示成功消息

```go
gui.ShowSuccess(
    "操作成功",
    "分支 'feature-x' 已成功删除。",
    "本地分支和远程分支都已删除。",
)
```

### 示例 5: 确认对话框

```go
gui.ShowConfirm(
    "确认删除",
    "确定要删除分支 'feature-x' 吗？此操作不可撤销。",
    func() {
        // 用户点击"确定"后执行的操作
        gui.c.Toast("已确认删除")
    },
)
```

### 示例 6: 是/否/取消对话框

```go
gui.ShowYesNoCancel(
    "保存更改",
    "检测到未保存的更改，是否保存？",
    func() {
        // 用户点击"是"
        gui.c.Toast("已保存")
    },
    func() {
        // 用户点击"否"
        gui.c.Toast("已放弃更改")
    },
)
```

### 示例 7: 自定义按钮

```go
gui.ShowMessageBox(MessageBoxConfig{
    Type:    MessageTypeQuestion,
    Title:   "选择操作",
    Message: "检测到未提交的更改，如何处理？",
    Buttons: []string{"暂存", "丢弃", "取消"},
}, func(buttonIndex int) {
    switch buttonIndex {
    case 0:
        gui.c.Toast("已暂存更改")
    case 1:
        gui.c.Toast("已丢弃更改")
    case 2:
        gui.c.Toast("已取消操作")
    }
})
```

### 示例 8: 自动关闭消息框

```go
gui.ShowAutoCloseMessage(
    MessageTypeSuccess,
    "操作成功",
    "文件已保存，此消息将在 3 秒后自动关闭。",
    3*time.Second,
)
```

## API 参考

### MessageBoxConfig 结构

```go
type MessageBoxConfig struct {
    Type    MessageType // 消息类型
    Title   string      // 标题
    Message string      // 消息内容
    Details string      // 详细信息（可选）
    Buttons []string    // 按钮列表
    Width   int         // 宽度（0表示默认60）
    Height  int         // 高度（0表示自动）
}
```

### MessageType 常量

```go
const (
    MessageTypeInfo     MessageType = iota // 信息
    MessageTypeSuccess                     // 成功
    MessageTypeWarning                     // 警告
    MessageTypeError                       // 错误
    MessageTypeQuestion                    // 问题
)
```

### 便捷方法

```go
// 显示错误消息框
gui.ShowError(title, message string, details ...string)

// 显示警告消息框
gui.ShowWarning(title, message string, details ...string)

// 显示信息消息框
gui.ShowInfo(title, message string, details ...string)

// 显示成功消息框
gui.ShowSuccess(title, message string, details ...string)

// 显示确认对话框（确定/取消）
gui.ShowConfirm(title, message string, onConfirm func())

// 显示是/否/取消对话框
gui.ShowYesNoCancel(title, message string, onYes, onNo func())

// 显示自动关闭的消息框
gui.ShowAutoCloseMessage(msgType MessageType, title, message string, duration time.Duration)

// 显示完全自定义的消息框
gui.ShowMessageBox(config MessageBoxConfig, onClose func(buttonIndex int)) *MessageBox
```

## 实际应用场景

### Git Push 错误处理

```go
func (gui *Gui) handleGitPushError(err error) {
    gui.ShowError(
        "推送失败",
        "无法推送到远程仓库。",
        err.Error(),
    )
}
```

### 确认强制推送

```go
func (gui *Gui) confirmForcePush(branch string) {
    gui.ShowMessageBox(MessageBoxConfig{
        Type:    MessageTypeWarning,
        Title:   "确认强制推送",
        Message: "你即将强制推送到远程分支 '" + branch + "'，这将覆盖远程仓库的历史记录！此操作不可撤销。",
        Buttons: []string{"确认", "取消"},
    }, func(buttonIndex int) {
        if buttonIndex == 0 {
            // 执行强制推送
            gui.c.Toast("正在强制推送...")
        }
    })
}
```

### 确认删除分支

```go
func (gui *Gui) confirmDeleteBranch(branch string, hasRemote bool) {
    message := "确定要删除分支 '" + branch + "' 吗？"
    if hasRemote {
        message = "确定要删除本地和远程分支 '" + branch + "' 吗？"
    }

    gui.ShowConfirm(
        "确认删除",
        message,
        func() {
            // 执行删除操作
            gui.c.Toast("分支已删除")
        },
    )
}
```

### 显示合并冲突

```go
func (gui *Gui) showMergeConflict(conflictFiles []string) {
    details := "冲突文件:\n"
    for _, file := range conflictFiles {
        details += "  - " + file + "\n"
    }

    gui.ShowMessageBox(MessageBoxConfig{
        Type:    MessageTypeWarning,
        Title:   "合并冲突",
        Message: "合并过程中发现冲突，请解决冲突后再提交。",
        Details: details,
        Buttons: []string{"解决冲突", "中止合并"},
    }, func(buttonIndex int) {
        if buttonIndex == 0 {
            // 打开冲突解决界面
            gui.c.Toast("打开冲突解决界面")
        } else {
            // 中止合并
            gui.c.Toast("已中止合并")
        }
    })
}
```

### 显示暂存选项

```go
func (gui *Gui) showStashOptions() {
    gui.ShowMessageBox(MessageBoxConfig{
        Type:    MessageTypeQuestion,
        Title:   "未提交的更改",
        Message: "检测到未提交的更改，如何处理？",
        Buttons: []string{"暂存", "丢弃", "取消"},
    }, func(buttonIndex int) {
        switch buttonIndex {
        case 0:
            gui.c.Toast("已暂存更改")
        case 1:
            gui.ShowConfirm(
                "确认丢弃",
                "确定要丢弃所有未提交的更改吗？此操作不可撤销。",
                func() {
                    gui.c.Toast("已丢弃更改")
                },
            )
        }
    })
}
```

### 显示操作成功

```go
func (gui *Gui) showOperationSuccess(operation, details string) {
    gui.ShowAutoCloseMessage(
        MessageTypeSuccess,
        operation+"成功",
        details,
        2*time.Second,
    )
}
```

## 测试功能

项目中包含了完整的测试示例，可以通过以下方式测试：

1. 在代码中调用测试菜单：
```go
gui.createMessageBoxTestMenu()
```

2. 测试菜单包含以下选项：
   - 错误消息框
   - 警告消息框
   - 信息消息框
   - 成功消息框
   - 确认对话框
   - 是/否/取消对话框
   - 自定义按钮
   - 自动关闭消息框
   - 长文本消息框
   - 测试所有消息类型

## 文件结构

```
pkg/gui/
├── message_box.go              # 消息框核心实现
├── message_box_examples.go     # 使用示例
└── message_box_test_menu.go    # 测试菜单
```

## 注意事项

1. **线程安全**：消息框的显示会自动在 GUI 线程中执行，可以安全地从后台协程调用。

2. **资源清理**：消息框关闭时会自动清理键盘绑定和视图，无需手动清理。

3. **按钮数量**：
   - 默认只有一个"确定"按钮
   - 最多支持 9 个按钮（可通过数字键 1-9 快速选择）
   - 建议不超过 3-4 个按钮以保持界面简洁

4. **文本长度**：
   - 消息内容会自动换行
   - 默认宽度为 60 字符
   - 可通过 `Width` 参数自定义宽度

5. **详细信息**：
   - 详细信息以较淡的颜色显示
   - 适合显示错误堆栈、日志等技术信息

6. **回调函数**：
   - `onClose` 回调接收按钮索引（0 表示第一个按钮）
   - 按 Esc 键会选择最后一个按钮
   - 回调函数可以为 `nil`（仅显示信息，不需要处理）

## 设计原则

1. **终端原生风格**：使用 Unicode 字符绘制边框，保持 TUI 风格
2. **清晰的视觉层次**：通过图标、颜色、分隔线区分不同部分
3. **直观的交互**：支持键盘导航，快捷键符合直觉
4. **灵活的配置**：支持自定义按钮、宽度、高度等
5. **一致的 API**：提供便捷方法和完全自定义两种使用方式

## 相关文档

- [MESSAGEBOX_DESIGN.md](./MESSAGEBOX_DESIGN.md) - 详细设计文档
- [PROGRESSBAR_DESIGN.md](./PROGRESSBAR_DESIGN.md) - ProgressBar 设计文档
- [PROGRESSBAR_USAGE.md](./PROGRESSBAR_USAGE.md) - ProgressBar 使用文档

---

**实现完成**: 2024
**状态**: ✅ 已实现并可用
