# 开发指南

## 📖 目录

- [开发环境搭建](#开发环境搭建)
- [项目结构](#项目结构)
- [开发流程](#开发流程)
- [代码规范](#代码规范)
- [测试指南](#测试指南)
- [调试技巧](#调试技巧)
- [贡献指南](#贡献指南)

---

## 开发环境搭建

### 系统要求

```
操作系统: Windows 10+, macOS 10.12+, Linux
Go 版本: 1.25.0+
Git 版本: 2.20.0+
终端: 支持 256 色
```

### 安装 Go

```bash
# macOS
brew install go

# Linux
sudo apt install golang-go

# Windows
# 从 https://golang.org/dl/ 下载安装
```

### 克隆项目

```bash
git clone https://github.com/dswcpp/lazygit.git
cd lazygit
```

### 安装依赖

```bash
go mod download
go mod vendor
```

### 编译项目

```bash
# 开发编译
go build

# 生产编译
go build -ldflags="-s -w"

# 交叉编译
GOOS=linux GOARCH=amd64 go build
GOOS=darwin GOARCH=amd64 go build
GOOS=windows GOARCH=amd64 go build
```

### 运行项目

```bash
# 直接运行
go run main.go

# 运行编译后的二进制
./lazygit

# 调试模式
./lazygit --debug

# 查看日志
./lazygit --logs
```

---

## 项目结构

### 目录说明

```
lazygit/
├── cmd/                       # 命令行入口
├── pkg/                       # 核心代码
│   ├── ai/                    # AI 功能模块
│   │   ├── provider.go        # AI 提供商接口
│   │   ├── deepseek.go        # DeepSeek 实现
│   │   ├── openai.go          # OpenAI 实现
│   │   └── anthropic.go       # Anthropic 实现
│   ├── app/                   # 应用层
│   │   ├── entry_point.go     # 应用入口
│   │   └── daemon/            # 后台服务
│   ├── commands/              # Git 命令封装
│   │   ├── git_commands/      # Git 操作
│   │   ├── git_config/        # Git 配置
│   │   └── oscommands/        # 系统命令
│   ├── config/                # 配置管理
│   │   ├── user_config.go     # 用户配置
│   │   └── app_config.go      # 应用配置
│   ├── gui/                   # UI 层
│   │   ├── gui.go             # GUI 主结构
│   │   ├── ai_chat.go         # AI 对话 ⭐
│   │   ├── message_box.go     # 消息框 ⭐
│   │   ├── progress_bar.go    # 进度条 ⭐
│   │   ├── context/           # 上下文管理
│   │   ├── controllers/       # 控制器
│   │   │   ├── helpers/       # 辅助函数
│   │   │   └── *.go           # 各种控制器
│   │   ├── presentation/      # 展示层
│   │   └── types/             # 类型定义
│   ├── i18n/                  # 国际化
│   │   └── translations/      # 翻译文件
│   ├── integration/           # 集成测试
│   └── utils/                 # 工具函数
├── docs/                      # 文档
├── test/                      # 测试
├── vendor/                    # 依赖包
├── go.mod                     # Go 模块定义
├── go.sum                     # 依赖校验
├── main.go                    # 主入口
└── Makefile                   # 构建脚本
```

### 核心模块

#### 1. AI 模块 (`pkg/ai/`)

```go
// provider.go - AI 提供商接口
type AIProvider interface {
    Complete(ctx context.Context, prompt string) (*AIResponse, error)
    CompleteStream(ctx context.Context, prompt string) (<-chan string, error)
}

// deepseek.go - DeepSeek 实现
type DeepSeekProvider struct {
    apiKey  string
    model   string
    baseURL string
}

func (p *DeepSeekProvider) Complete(ctx context.Context, prompt string) (*AIResponse, error) {
    // 实现
}
```

#### 2. GUI 模块 (`pkg/gui/`)

```go
// gui.go - GUI 主结构
type Gui struct {
    *common.Common
    g          *gocui.Gui
    git        *commands.GitCommand
    os         *oscommands.OSCommand
    State      *GuiRepoState
    Config     config.AppConfigurer
    // ...
}

// ai_chat.go - AI 对话
type AIChat struct {
    gui          *Gui
    chatView     *gocui.View
    inputView    *gocui.View
    messages     []ChatMessage
    // ...
}
```

#### 3. Commands 模块 (`pkg/commands/`)

```go
// git_commands/branch.go - 分支操作
type BranchCommands struct {
    *common.Common
    gitCommon *GitCommon
}

func (c *BranchCommands) New(name string) error {
    return c.gitCommon.RunCommand("git branch %s", name)
}
```

---

## 开发流程

### 1. 创建功能分支

```bash
git checkout -b feature/your-feature
```

### 2. 开发功能

#### 添加新的 AI 功能

```go
// 1. 在 pkg/ai/ 中添加新的提供商
type NewProvider struct {
    apiKey string
    model  string
}

func (p *NewProvider) Complete(ctx context.Context, prompt string) (*AIResponse, error) {
    // 实现
}

// 2. 在 config 中添加配置支持
type AIConfig struct {
    Profiles []AIProfile `yaml:"profiles"`
}

// 3. 在 GUI 中集成
func (gui *Gui) useNewProvider() {
    // 使用新提供商
}
```

#### 添加新的 UI 组件

```go
// 1. 在 pkg/gui/ 中创建新文件
// new_component.go

type NewComponent struct {
    gui  *Gui
    view *gocui.View
    // ...
}

func (gui *Gui) ShowNewComponent() error {
    // 实现
}

// 2. 添加键盘绑定
func (gui *Gui) setNewComponentKeyBindings() {
    gui.g.SetKeybinding("newComponent", gocui.KeyEnter, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
        // 处理
        return nil
    })
}
```

### 3. 编写测试

```go
// new_component_test.go
func TestNewComponent(t *testing.T) {
    gui := setupTestGui()

    err := gui.ShowNewComponent()
    assert.NoError(t, err)

    // 验证
}
```

### 4. 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/gui/

# 运行特定测试
go test -run TestNewComponent

# 查看覆盖率
go test -cover ./...
```

### 5. 提交代码

```bash
git add .
git commit -m "feat(gui): 添加新组件"
git push origin feature/your-feature
```

### 6. 创建 Pull Request

在 GitHub 上创建 PR，描述你的更改。

---

## 代码规范

### Go 代码规范

#### 1. 命名规范

```go
// 包名：小写，简短
package gui

// 类型名：PascalCase
type MessageBox struct {}

// 函数名：camelCase (导出) 或 camelCase (私有)
func ShowMessageBox() {}
func createView() {}

// 常量：PascalCase 或 UPPER_SNAKE_CASE
const MaxRetries = 3
const DEFAULT_TIMEOUT = 60

// 变量：camelCase
var messageCount int
```

#### 2. 注释规范

```go
// Package gui provides the terminal user interface for lazygit.
package gui

// MessageBox represents a modal dialog box.
// It supports multiple message types and custom buttons.
type MessageBox struct {
    config MessageBoxConfig
    view   *gocui.View
}

// ShowMessageBox displays a message box with the given configuration.
// It returns an error if the view cannot be created.
func (gui *Gui) ShowMessageBox(config MessageBoxConfig) error {
    // Implementation
}
```

#### 3. 错误处理

```go
// ✅ 好的错误处理
func (gui *Gui) doSomething() error {
    if err := gui.validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    result, err := gui.process()
    if err != nil {
        return fmt.Errorf("process failed: %w", err)
    }

    return nil
}

// ❌ 不好的错误处理
func (gui *Gui) doSomething() error {
    gui.validate()  // 忽略错误
    result, _ := gui.process()  // 忽略错误
    return nil
}
```

#### 4. 代码组织

```go
// 结构体定义
type MessageBox struct {
    // 导出字段在前
    Config MessageBoxConfig

    // 私有字段在后
    gui         *Gui
    view        *gocui.View
    selectedBtn int
}

// 构造函数
func NewMessageBox(gui *Gui, config MessageBoxConfig) *MessageBox {
    return &MessageBox{
        Config: config,
        gui:    gui,
    }
}

// 公共方法
func (mb *MessageBox) Show() error {
    // Implementation
}

// 私有方法
func (mb *MessageBox) render() {
    // Implementation
}
```

### 提交信息规范

使用 Conventional Commits 格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型 (type)**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具

**示例**:
```
feat(ai): 添加 DeepSeek 支持

- 实现 DeepSeek API 集成
- 添加配置选项
- 更新文档

Closes #123
```

---

## 测试指南

### 单元测试

```go
// message_box_test.go
package gui

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestMessageBox_Show(t *testing.T) {
    // 准备
    gui := setupTestGui()
    config := MessageBoxConfig{
        Type:    MessageTypeInfo,
        Title:   "Test",
        Message: "Test message",
    }

    // 执行
    err := gui.ShowMessageBox(config)

    // 验证
    assert.NoError(t, err)
    assert.NotNil(t, gui.g.View("messageBox"))
}

func TestMessageBox_ButtonSelection(t *testing.T) {
    gui := setupTestGui()
    mb := NewMessageBox(gui, MessageBoxConfig{
        Buttons: []string{"OK", "Cancel"},
    })

    // 测试按钮选择
    mb.selectButton(0)
    assert.Equal(t, 0, mb.selectedBtn)

    mb.selectButton(1)
    assert.Equal(t, 1, mb.selectedBtn)
}
```

### 集成测试

```go
// integration_test.go
func TestAIChatIntegration(t *testing.T) {
    // 设置测试环境
    gui := setupIntegrationTest()

    // 打开 AI 对话
    err := gui.ShowAIChat()
    assert.NoError(t, err)

    // 发送消息
    chat := gui.aiChat
    chat.inputView.SetContent("test message")
    err = chat.sendMessage()
    assert.NoError(t, err)

    // 等待响应
    time.Sleep(2 * time.Second)

    // 验证响应
    assert.Greater(t, len(chat.messages), 1)
}
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包
go test ./pkg/gui/

# 查看覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 调试技巧

### 1. 使用调试模式

```bash
# 启动调试模式
./lazygit --debug

# 查看日志
./lazygit --logs

# 同时查看程序和日志
# 终端 1
./lazygit --debug

# 终端 2
./lazygit --logs
```

### 2. 添加日志

```go
// 使用 gui.c.Log
gui.c.Log.Info("Message box shown")
gui.c.Log.Warn("Warning message")
gui.c.Log.Error("Error occurred")

// 使用 fmt.Fprintf 输出到日志文件
fmt.Fprintf(gui.c.Log.GetLogFile(), "Debug: %v\n", data)
```

### 3. 使用 Delve 调试器

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug

# 设置断点
(dlv) break main.main
(dlv) break pkg/gui/ai_chat.go:100

# 运行
(dlv) continue

# 查看变量
(dlv) print variableName

# 单步执行
(dlv) next
(dlv) step
```

### 4. 性能分析

```bash
# CPU 分析
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# 内存分析
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# 运行时分析
go build -o lazygit
./lazygit --cpuprofile=cpu.prof
go tool pprof cpu.prof
```

---

## 贡献指南

### 1. Fork 项目

在 GitHub 上 Fork 项目到你的账号。

### 2. 克隆 Fork

```bash
git clone https://github.com/your-username/lazygit.git
cd lazygit
```

### 3. 添加上游仓库

```bash
git remote add upstream https://github.com/dswcpp/lazygit.git
```

### 4. 创建功能分支

```bash
git checkout -b feature/your-feature
```

### 5. 开发和测试

按照上面的开发流程进行开发和测试。

### 6. 提交更改

```bash
git add .
git commit -m "feat: your feature description"
```

### 7. 推送到 Fork

```bash
git push origin feature/your-feature
```

### 8. 创建 Pull Request

在 GitHub 上创建 PR，描述你的更改。

### 9. 代码审查

等待维护者审查你的代码，根据反馈进行修改。

### 10. 合并

PR 被批准后，维护者会合并你的代码。

---

## 常见问题

### Q1: 编译失败？

**检查**:
- Go 版本是否正确
- 依赖是否完整
- 环境变量是否设置

```bash
go version
go mod download
go mod vendor
```

### Q2: 测试失败？

**检查**:
- 测试环境是否正确
- 依赖是否安装
- 配置是否正确

```bash
go test -v ./...
```

### Q3: 如何添加新的 AI 提供商？

1. 在 `pkg/ai/` 中实现 `AIProvider` 接口
2. 在 `config` 中添加配置支持
3. 在 GUI 中集成
4. 添加测试
5. 更新文档

### Q4: 如何调试 UI 问题？

1. 使用 `--debug` 模式
2. 查看日志输出
3. 使用 `fmt.Fprintf` 添加调试信息
4. 检查视图创建和更新

---

## 相关资源

### 官方文档
- [Go 官方文档](https://golang.org/doc/)
- [gocui 文档](https://github.com/jroimartin/gocui)
- [原版 lazygit](https://github.com/jesseduffield/lazygit)

### 开发工具
- [VS Code](https://code.visualstudio.com/)
- [GoLand](https://www.jetbrains.com/go/)
- [Delve](https://github.com/go-delve/delve)

### 相关文档
- [项目概述](./PROJECT_OVERVIEW.md)
- [架构设计](./ARCHITECTURE.md)
- [API 参考](./API_REFERENCE.md)

---

**版本**: v1.0.0
**最后更新**: 2024
**状态**: ✅ 完整
