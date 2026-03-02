# VSCode 风格活动栏 - 使用说明

## 🎉 实现完成

我们已经成功为 lazygit 实现了 VSCode 风格的活动栏功能！

## ✨ 功能特性

### 1. **非侵入式设计**
- ✅ 默认关闭，不影响现有用户
- ✅ 通过配置文件启用
- ✅ 原有布局和功能完全保留

### 2. **VSCode 风格界面**
```
┌──┬────────────┬──────────────────────┐
│📁│ status     │ main                 │
│🌿│ files      │                      │
│📊│ branches   │                      │
│💾│ commits    │                      │
│📦│ stash      │                      │
│──│            │                      │
│⬇️│            │                      │  ← 当前选中（蓝色加粗）
│⬆️│            │                      │
│🔄│            │                      │
│💾│            │                      │
│🔀│            │                      │
│♻️│            │                      │
│──│            │                      │
│⚙️│            │                      │
│❓│            │                      │
└──┴────────────┴──────────────────────┘
```

### 3. **三大功能区**

#### 📁 **导航区**（顶部）
- 📁 Status - 状态面板
- 🌿 Files - 文件列表
- 📊 Branches - 分支列表
- 💾 Commits - 提交历史
- 📦 Stash - 暂存区

#### ⚡ **操作区**（中部）
- ⬇️ Pull - 拉取代码
- ⬆️ Push - 推送代码
- 🔄 Fetch - 获取更新
- 💾 Stash - 暂存更改
- 🔀 Merge - 合并分支
- ♻️ Rebase - 变基操作

#### 🔧 **工具区**（底部）
- ⚙️ Settings - 设置
- ❓ Help - 帮助

## 🚀 如何启用

### 方法 1：编辑配置文件

编辑 `~/.config/lazygit/config.yml`（Linux/Mac）或 `%APPDATA%\lazygit\config.yml`（Windows）：

```yaml
gui:
  activityBar:
    show: true   # 启用活动栏
    width: 3     # 宽度（2-10 个字符）
```

### 方法 2：使用测试配置

使用项目根目录下的 `test_activity_bar_config.yml`：

```bash
cd /e/code/go/lazygit
./lazygit.exe --use-config-file test_activity_bar_config.yml
# 或使用短参数
./lazygit.exe -ucf test_activity_bar_config.yml
```

## ⌨️ 键盘操作

| 按键 | 功能 |
|------|------|
| `↑` 或 `k` | 上移选择 |
| `↓` 或 `j` | 下移选择 |
| `Enter` | 执行选中项 |

## 🖱️ 鼠标操作

| 操作 | 功能 |
|------|------|
| 左键点击 | 选择并执行 |
| 滚轮上滚 | 上移选择 |
| 滚轮下滚 | 下移选择 |

## 📝 配置选项

```yaml
gui:
  activityBar:
    show: false    # 是否显示活动栏（默认：false）
    width: 3       # 活动栏宽度（默认：3，范围：2-10）
```

## 🎨 视觉效果

- **正常状态**：默认颜色显示图标
- **选中状态**：蓝色加粗显示
- **分隔符**：`──` 分隔不同功能区

## 📂 修改的文件

### 新增文件
- `pkg/gui/activity_bar.go` - 活动栏核心逻辑

### 修改的文件
- `pkg/config/user_config.go` - 添加配置结构
- `pkg/gui/controllers/helpers/window_arrangement_helper.go` - 布局系统
- `pkg/gui/types/views.go` - 视图注册
- `pkg/gui/gui.go` - 状态管理
- `pkg/gui/views.go` - 视图初始化
- `pkg/gui/keybindings.go` - 键盘和鼠标绑定

## 🔧 技术细节

### 架构设计
- **条件性加载**：仅在配置启用时创建相关组件
- **独立模块**：所有新代码在独立文件中
- **零破坏性**：原有逻辑完全保留

### 布局系统
使用 `boxlayout` 包实现响应式布局：
- 活动栏固定宽度（默认 3 字符）
- 原有面板自适应剩余空间
- 支持动态显示/隐藏

## 🐛 已知限制

### 当前版本（MVP）
1. **操作功能**：Git 操作（Pull/Push/Fetch 等）目前只显示提示信息
2. **工具功能**：设置和帮助功能目前只显示提示信息

### 后续增强计划
- [ ] 连接实际的 Git 操作
- [ ] 添加加载状态动画
- [ ] 支持自定义图标和操作
- [ ] 添加 Tooltip 提示
- [ ] 支持快捷键（Ctrl+1~5）快速跳转

## 🎯 测试步骤

1. **编译项目**
```bash
cd /e/code/go/lazygit
go build -o lazygit.exe
```

2. **启用活动栏**
```bash
./lazygit.exe -c test_activity_bar_config.yml
```

3. **测试功能**
- 使用上下键或鼠标滚轮移动选择
- 点击或按 Enter 执行操作
- 测试导航功能（切换到不同面板）

## 💡 使用建议

### 适合启用的场景
- ✅ 宽屏显示器（>100 列）
- ✅ 喜欢鼠标操作
- ✅ 需要快速切换面板
- ✅ 新手用户（提高可发现性）

### 不建议启用的场景
- ❌ 小屏幕终端（<80 列）
- ❌ 纯键盘用户（已有快捷键）
- ❌ 追求极简界面

## 🤝 贡献

如果你想增强活动栏功能，可以修改：
- `pkg/gui/activity_bar.go` - 添加新的操作项
- `pkg/config/user_config.go` - 添加新的配置选项

## 📄 许可证

遵循 lazygit 项目的 MIT 许可证。

---

**享受你的新活动栏！** 🎉
