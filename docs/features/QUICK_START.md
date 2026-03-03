# lazygit UI 功能快速入门指南

## 🚀 5 分钟快速上手

### 第一步：了解三大功能

lazygit 现在包含三个强大的 UI 功能组件：

```
┌─────────────────────────────────────────────────────────┐
│  📬 消息框    │  📊 进度条    │  🤖 AI 对话           │
│  提示和确认   │  进度显示     │  智能助手             │
└─────────────────────────────────────────────────────────┘
```

### 第二步：测试功能

在代码中调用测试菜单：

```go
gui.CreateUIFeaturesTestMenu()
```

这将打开一个统一的测试中心，包含所有功能的演示。

### 第三步：基本使用

#### 消息框 - 3 行代码

```go
// 显示错误
gui.ShowError("操作失败", "无法连接到服务器")

// 确认操作
gui.ShowConfirm("确认删除", "确定要删除吗？", func() {
    // 执行删除
})

// 自定义按钮
gui.ShowMessageBox(MessageBoxConfig{
    Type:    MessageTypeQuestion,
    Title:   "选择操作",
    Message: "如何处理？",
    Buttons: []string{"保存", "丢弃", "取消"},
}, func(buttonIndex int) {
    // 处理选择
})
```

#### 进度条 - 5 行代码

```go
// 显示进度条
pb := gui.ShowProgressBar(ProgressBarConfig{
    Title:          "正在处理...",
    Total:          100,
    ShowPercentage: true,
})

// 更新进度
go func() {
    for i := int64(0); i <= 100; i++ {
        pb.Update(i, "")
        time.Sleep(50 * time.Millisecond)
    }
    pb.Close()
}()
```

#### AI 对话 - 1 行代码

```go
// 打开 AI 对话
gui.ShowAIChat()
```

## 📚 常见场景

### 场景 1：Git 操作确认

```go
func (gui *Gui) confirmPush() error {
    gui.ShowConfirm(
        "确认推送",
        "确定要推送到 origin/main 吗？",
        func() {
            // 显示进度
            pb := gui.ShowProgressBar(ProgressBarConfig{
                Title: "正在推送...",
                Total: fileSize,
            })

            // 执行推送
            go func() {
                // ... 推送逻辑
                pb.Close()

                // 显示结果
                gui.g.Update(func(*gocui.Gui) error {
                    gui.ShowSuccess("推送成功", "已推送 3 个提交")
                    return nil
                })
            }()
        },
    )
    return nil
}
```

### 场景 2：错误处理

```go
func (gui *Gui) handleError(err error) {
    gui.ShowError(
        "操作失败",
        "无法完成操作",
        err.Error(),
    )
}
```

### 场景 3：长时间操作

```go
func (gui *Gui) cloneRepo(url string) error {
    pb := gui.ShowProgressBar(ProgressBarConfig{
        Title:         "正在克隆仓库...",
        Indeterminate: true,
    })

    go func() {
        // 执行克隆
        err := git.Clone(url)
        pb.Close()

        if err != nil {
            gui.g.Update(func(*gocui.Gui) error {
                gui.ShowError("克隆失败", err.Error())
                return nil
            })
        } else {
            gui.g.Update(func(*gocui.Gui) error {
                gui.ShowSuccess("克隆成功", "仓库已克隆到本地")
                return nil
            })
        }
    }()

    return nil
}
```

### 场景 4：AI 辅助

```go
func (gui *Gui) getHelp() error {
    if gui.c.AI == nil {
        gui.ShowError("AI 未启用", "请先配置 AI")
        return nil
    }
    return gui.ShowAIChat()
}
```

## 🎯 最佳实践

### 1. 消息框使用原则

✅ **推荐**：
- 重要操作前使用确认对话框
- 错误信息使用 ShowError
- 成功反馈使用 ShowSuccess
- 危险操作使用 Warning 类型

❌ **避免**：
- 过度使用消息框打断用户
- 不重要的信息使用 Toast 即可
- 消息内容过长（使用 Details 字段）

### 2. 进度条使用原则

✅ **推荐**：
- 超过 2 秒的操作显示进度条
- 能计算进度的使用确定进度
- 不能计算的使用不确定进度
- 操作完成后立即关闭

❌ **避免**：
- 短时间操作（< 1 秒）显示进度条
- 忘记调用 Close() 关闭进度条
- 在主线程中更新进度（使用 goroutine）

### 3. AI 对话使用原则

✅ **推荐**：
- 复杂问题使用 AI 对话
- 利用预设问题快速提问
- 使用 Ctrl+P 浏览常见问题
- 保存有价值的对话记录

❌ **避免**：
- 简单问题过度依赖 AI
- 在对话中包含敏感信息
- 忘记检查 AI 是否启用

## 🔧 配置

### 进度条配置

在 `~/.config/lazygit/config.yml` 中：

```yaml
gui:
  progressBar:
    style: "block"           # 样式: block, dot, arrow, gradient, ascii
    spinnerStyle: "braille"  # 动画: braille, line, arrow, dot, circle
    width: 30                # 宽度
    showPercentage: true     # 显示百分比
    showStats: true          # 显示统计信息
```

### AI 配置

```yaml
ai:
  enabled: true
  activeProfile: "default"
  profiles:
    - name: "default"
      provider: "deepseek"
      apiKey: "your-api-key"
      model: "deepseek-chat"
      maxTokens: 2000
      timeout: 60
```

## 📖 快捷键速查

### 消息框

| 快捷键 | 功能 |
|--------|------|
| Enter | 确认当前按钮 |
| Esc | 取消（最后一个按钮） |
| ←/→ | 切换按钮 |
| Tab | 循环切换按钮 |
| 1-9 | 快速选择按钮 |

### AI 对话

| 快捷键 | 功能 |
|--------|------|
| Enter | 发送消息 |
| Alt+Enter | 换行 |
| ↑/↓ | 浏览历史 |
| Ctrl+P | 预设问题 |
| Ctrl+L | 清空历史 |
| Ctrl+C | 复制回复 |
| Ctrl+K | 停止生成 |
| Tab | 切换焦点 |
| Esc | 关闭 |
| ? | 帮助 |

## 🎬 演示视频

运行完整演示：

```go
gui.RunFullDemo()
```

运行实际场景演示：

```go
gui.RunRealWorldDemo()
```

## 📝 代码模板

### 模板 1：标准操作流程

```go
func (gui *Gui) standardOperation() error {
    // 1. 确认
    gui.ShowConfirm("确认操作", "确定要执行吗？", func() {
        // 2. 显示进度
        pb := gui.ShowProgressBar(ProgressBarConfig{
            Title: "正在处理...",
            Total: total,
        })

        // 3. 执行操作
        go func() {
            err := doSomething()
            pb.Close()

            // 4. 显示结果
            gui.g.Update(func(*gocui.Gui) error {
                if err != nil {
                    gui.ShowError("操作失败", err.Error())
                } else {
                    gui.ShowSuccess("操作成功", "已完成")
                }
                return nil
            })
        }()
    })
    return nil
}
```

### 模板 2：多选项操作

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
            // 取消，不做任何操作
        }
    })
    return nil
}
```

### 模板 3：AI 辅助操作

```go
func (gui *Gui) aiAssistedOperation() error {
    if gui.c.AI == nil {
        gui.ShowError("AI 未启用", "请先配置 AI")
        return nil
    }

    gui.ShowInfo(
        "AI 助手",
        "有任何问题都可以问我",
        "提示：使用 Ctrl+P 查看预设问题",
    )

    return gui.ShowAIChat()
}
```

## 🐛 故障排除

### 问题 1：消息框不显示

**原因**：可能是视图创建失败

**解决**：检查 gocui 是否正确初始化

```go
if gui.g == nil {
    return errors.New("GUI not initialized")
}
```

### 问题 2：进度条卡住

**原因**：忘记调用 Close()

**解决**：确保在 goroutine 结束时调用

```go
defer pb.Close()
```

### 问题 3：AI 对话无响应

**原因**：AI 未启用或配置错误

**解决**：检查配置

```go
if gui.c.AI == nil {
    gui.ShowError("AI 未启用", "请检查配置")
    return nil
}
```

## 📚 进阶阅读

- [消息框详细文档](./MESSAGEBOX_USAGE.md)
- [进度条详细文档](./PROGRESSBAR_USAGE.md)
- [AI 对话详细文档](./AI_CHAT_USAGE_V2.md)
- [设计文档](./UI_FEATURES_SUMMARY.md)

## 💡 提示

1. **先测试再使用**：使用 `CreateUIFeaturesTestMenu()` 测试所有功能
2. **查看示例**：每个功能都有 `*_examples.go` 文件
3. **阅读文档**：每个功能都有详细的使用文档
4. **参考演示**：运行 `RunFullDemo()` 查看完整演示

## 🎉 开始使用

现在你已经掌握了基础知识，可以开始在项目中使用这些功能了！

```go
// 打开测试菜单
gui.CreateUIFeaturesTestMenu()

// 或者直接使用
gui.ShowError("错误", "操作失败")
gui.ShowProgressBar(config)
gui.ShowAIChat()
```

祝你使用愉快！🚀

---

**版本**: v1.0
**最后更新**: 2024
**状态**: ✅ 可用
