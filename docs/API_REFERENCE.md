# API 参考文档

## 📖 目录

- [GUI API](#gui-api)
- [AI API](#ai-api)
- [Commands API](#commands-api)
- [Config API](#config-api)
- [Utils API](#utils-api)

---

## GUI API

### MessageBox API

#### ShowError

显示错误消息框。

```go
func (gui *Gui) ShowError(title, message string, details ...string) error
```

**参数**:
- `title`: 标题
- `message`: 消息内容
- `details`: 详细信息（可选）

**示例**:
```go
gui.ShowError("推送失败", "无法连接到远程仓库", err.Error())
```

#### ShowWarning

显示警告消息框。

```go
func (gui *Gui) ShowWarning(title, message string, details ...string) error
```

#### ShowInfo

显示信息消息框。

```go
func (gui *Gui) ShowInfo(title, message string, details ...string) error
```

#### ShowSuccess

显示成功消息框。

```go
func (gui *Gui) ShowSuccess(title, message string, details ...string) error
```

#### ShowConfirm

显示确认对话框。

```go
func (gui *Gui) ShowConfirm(title, message string, onConfirm func()) error
```

**参数**:
- `title`: 标题
- `message`: 消息内容
- `onConfirm`: 确认回调函数

**示例**:
```go
gui.ShowConfirm("确认删除", "确定要删除分支吗？", func() {
    gui.deleteBranch(branchName)
})
```

#### ShowYesNoCancel

显示三选项对话框。

```go
func (gui *Gui) ShowYesNoCancel(title, message string, onYes, onNo func()) error
```

#### ShowMessageBox

显示自定义消息框。

```go
func (gui *Gui) ShowMessageBox(config MessageBoxConfig, onClose func(int)) error
```

**MessageBoxConfig**:
```go
type MessageBoxConfig struct {
    Type    MessageType  // 消息类型
    Title   string       // 标题
    Message string       // 消息内容
    Details string       // 详细信息
    Buttons []string     // 按钮列表
}
```

**MessageType**:
```go
const (
    MessageTypeInfo     MessageType = iota
    MessageTypeSuccess
    MessageTypeWarning
    MessageTypeError
    MessageTypeQuestion
)
```

**示例**:
```go
gui.ShowMessageBox(MessageBoxConfig{
    Type:    MessageTypeQuestion,
    Title:   "选择操作",
    Message: "如何处理冲突？",
    Buttons: []string{"手动解决", "使用我们的", "使用他们的", "取消"},
}, func(buttonIndex int) {
    switch buttonIndex {
    case 0:
        gui.manualResolve()
    case 1:
        gui.resolveWithOurs()
    case 2:
        gui.resolveWithTheirs()
    }
})
```

---

### ProgressBar API

#### ShowProgressBar

显示进度条。

```go
func (gui *Gui) ShowProgressBar(config ProgressBarConfig) *ProgressBar
```

**ProgressBarConfig**:
```go
type ProgressBarConfig struct {
    Title          string  // 标题
    Total          int64   // 总量（确定进度）
    Indeterminate  bool    // 不确定进度
    Message        string  // 消息
    ShowPercentage bool    // 显示百分比
    ShowStats      bool    // 显示统计信息
}
```

**返回**: `*ProgressBar` 实例

**示例**:
```go
pb := gui.ShowProgressBar(ProgressBarConfig{
    Title:          "正在推送...",
    Total:          totalSize,
    ShowPercentage: true,
    ShowStats:      true,
})
```

#### ProgressBar.Update

更新进度。

```go
func (pb *ProgressBar) Update(current int64, message string)
```

**参数**:
- `current`: 当前进度
- `message`: 状态消息

**示例**:
```go
for i := int64(0); i <= total; i += step {
    pb.Update(i, fmt.Sprintf("处理中 %d/%d", i, total))
    time.Sleep(100 * time.Millisecond)
}
```

#### ProgressBar.Close

关闭进度条。

```go
func (pb *ProgressBar) Close()
```

**示例**:
```go
defer pb.Close()
```

#### ProgressBar.SwitchToDeterminate

切换到确定进度模式。

```go
func (pb *ProgressBar) SwitchToDeterminate(total int64)
```

---

### AI Chat API

#### ShowAIChat

显示 AI 对话窗口。

```go
func (gui *Gui) ShowAIChat() error
```

**示例**:
```go
if err := gui.ShowAIChat(); err != nil {
    return err
}
```

#### CloseAIChat

关闭 AI 对话窗口。

```go
func (gui *Gui) CloseAIChat(chat *AIChat) error
```

---

## AI API

### AIProvider Interface

AI 提供商接口。

```go
type AIProvider interface {
    Complete(ctx context.Context, prompt string) (*AIResponse, error)
    CompleteStream(ctx context.Context, prompt string) (<-chan string, error)
}
```

#### Complete

同步完成请求。

```go
func (p *AIProvider) Complete(ctx context.Context, prompt string) (*AIResponse, error)
```

**参数**:
- `ctx`: 上下文
- `prompt`: 提示词

**返回**: `*AIResponse` 和错误

**AIResponse**:
```go
type AIResponse struct {
    Content string  // 响应内容
    Usage   Usage   // 使用统计
}

type Usage struct {
    PromptTokens     int
    CompletionTokens int
    TotalTokens      int
}
```

**示例**:
```go
response, err := provider.Complete(ctx, "如何撤销提交？")
if err != nil {
    return err
}
fmt.Println(response.Content)
```

#### CompleteStream

流式完成请求。

```go
func (p *AIProvider) CompleteStream(ctx context.Context, prompt string) (<-chan string, error)
```

**返回**: 字符串通道和错误

**示例**:
```go
stream, err := provider.CompleteStream(ctx, prompt)
if err != nil {
    return err
}

for chunk := range stream {
    fmt.Print(chunk)
}
```

---

### DeepSeek Provider

#### NewDeepSeekProvider

创建 DeepSeek 提供商。

```go
func NewDeepSeekProvider(config AIProfile) *DeepSeekProvider
```

**示例**:
```go
provider := NewDeepSeekProvider(AIProfile{
    APIKey:  "sk-xxx",
    Model:   "deepseek-chat",
    BaseURL: "https://api.deepseek.com",
})
```

---

### OpenAI Provider

#### NewOpenAIProvider

创建 OpenAI 提供商。

```go
func NewOpenAIProvider(config AIProfile) *OpenAIProvider
```

---

### Anthropic Provider

#### NewAnthropicProvider

创建 Anthropic 提供商。

```go
func NewAnthropicProvider(config AIProfile) *AnthropicProvider
```

---

## Commands API

### Git Commands

#### Branch Commands

```go
type BranchCommands struct {
    *common.Common
    gitCommon *GitCommon
}
```

##### New

创建新分支。

```go
func (c *BranchCommands) New(name string) error
```

##### Delete

删除分支。

```go
func (c *BranchCommands) Delete(name string, force bool) error
```

##### Checkout

切换分支。

```go
func (c *BranchCommands) Checkout(name string) error
```

##### Rename

重命名分支。

```go
func (c *BranchCommands) Rename(oldName, newName string) error
```

##### GetBranches

获取分支列表。

```go
func (c *BranchCommands) GetBranches() ([]*models.Branch, error)
```

**示例**:
```go
// 创建分支
if err := gui.git.Branch.New("feature/new"); err != nil {
    return err
}

// 切换分支
if err := gui.git.Branch.Checkout("feature/new"); err != nil {
    return err
}

// 获取分支列表
branches, err := gui.git.Branch.GetBranches()
if err != nil {
    return err
}
```

---

#### Commit Commands

```go
type CommitCommands struct {
    *common.Common
    gitCommon *GitCommon
}
```

##### Commit

创建提交。

```go
func (c *CommitCommands) Commit(message string) error
```

##### Amend

修改最后一次提交。

```go
func (c *CommitCommands) Amend(message string) error
```

##### Revert

撤销提交。

```go
func (c *CommitCommands) Revert(sha string) error
```

##### GetCommits

获取提交列表。

```go
func (c *CommitCommands) GetCommits(limit int) ([]*models.Commit, error)
```

**示例**:
```go
// 创建提交
if err := gui.git.Commit.Commit("feat: add new feature"); err != nil {
    return err
}

// 修改提交
if err := gui.git.Commit.Amend("feat: update feature"); err != nil {
    return err
}
```

---

#### File Commands

```go
type FileCommands struct {
    *common.Common
    gitCommon *GitCommon
}
```

##### Stage

暂存文件。

```go
func (c *FileCommands) Stage(path string) error
```

##### Unstage

取消暂存。

```go
func (c *FileCommands) Unstage(path string) error
```

##### Discard

丢弃更改。

```go
func (c *FileCommands) Discard(path string) error
```

##### GetFiles

获取文件列表。

```go
func (c *FileCommands) GetFiles() ([]*models.File, error)
```

---

#### Remote Commands

```go
type RemoteCommands struct {
    *common.Common
    gitCommon *GitCommon
}
```

##### Push

推送到远程。

```go
func (c *RemoteCommands) Push(remoteName, branchName string, force bool) error
```

##### Pull

从远程拉取。

```go
func (c *RemoteCommands) Pull(remoteName, branchName string) error
```

##### Fetch

获取远程更新。

```go
func (c *RemoteCommands) Fetch(remoteName string) error
```

---

## Config API

### User Config

#### LoadConfig

加载配置。

```go
func LoadConfig() (*UserConfig, error)
```

#### SaveConfig

保存配置。

```go
func (c *UserConfig) SaveConfig() error
```

#### UserConfig Structure

```go
type UserConfig struct {
    GUI GuiConfig `yaml:"gui"`
    AI  AIConfig  `yaml:"ai"`
    Git GitConfig `yaml:"git"`
}

type GuiConfig struct {
    Theme        ThemeConfig        `yaml:"theme"`
    ProgressBar  ProgressBarConfig  `yaml:"progressBar"`
    ActivityBar  ActivityBarConfig  `yaml:"activityBar"`
}

type AIConfig struct {
    Enabled       bool        `yaml:"enabled"`
    ActiveProfile string      `yaml:"activeProfile"`
    Profiles      []AIProfile `yaml:"profiles"`
}

type AIProfile struct {
    Name      string `yaml:"name"`
    Provider  string `yaml:"provider"`
    APIKey    string `yaml:"apiKey"`
    Model     string `yaml:"model"`
    BaseURL   string `yaml:"baseURL"`
    MaxTokens int    `yaml:"maxTokens"`
    Timeout   int    `yaml:"timeout"`
}
```

**示例**:
```go
// 加载配置
config, err := LoadConfig()
if err != nil {
    return err
}

// 修改配置
config.AI.Enabled = true
config.AI.ActiveProfile = "deepseek"

// 保存配置
if err := config.SaveConfig(); err != nil {
    return err
}
```

---

## Utils API

### String Utils

#### Truncate

截断字符串。

```go
func Truncate(s string, maxLen int) string
```

#### Pad

填充字符串。

```go
func Pad(s string, length int, padChar rune) string
```

---

### File Utils

#### FileExists

检查文件是否存在。

```go
func FileExists(path string) bool
```

#### ReadFile

读取文件。

```go
func ReadFile(path string) (string, error)
```

#### WriteFile

写入文件。

```go
func WriteFile(path, content string) error
```

---

### Time Utils

#### FormatDuration

格式化时长。

```go
func FormatDuration(d time.Duration) string
```

**示例**:
```go
duration := 125 * time.Second
formatted := FormatDuration(duration)  // "2m 5s"
```

---

## 错误处理

### 错误类型

```go
var (
    ErrNotFound         = errors.New("not found")
    ErrInvalidInput     = errors.New("invalid input")
    ErrNetworkTimeout   = errors.New("network timeout")
    ErrPermissionDenied = errors.New("permission denied")
    ErrAINotEnabled     = errors.New("AI not enabled")
)
```

### 错误检查

```go
if errors.Is(err, ErrNotFound) {
    // 处理未找到错误
}

if errors.Is(err, ErrNetworkTimeout) {
    // 处理超时错误
}
```

---

## 事件系统

### 事件类型

```go
const (
    EventRefresh EventType = iota
    EventStateChange
    EventError
    EventAIResponse
)
```

### 订阅事件

```go
gui.EventManager.Subscribe(EventRefresh, func(event Event) {
    // 处理刷新事件
})
```

### 发布事件

```go
gui.EventManager.Publish(Event{
    Type: EventRefresh,
    Data: refreshData,
})
```

---

## 并发控制

### Mutex

```go
gui.Mutexes.RefreshingFilesMutex.Lock()
defer gui.Mutexes.RefreshingFilesMutex.Unlock()

// 执行需要保护的操作
```

### Context

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := provider.Complete(ctx, prompt)
```

---

## 最佳实践

### 1. 错误处理

```go
// ✅ 好的错误处理
if err := gui.ShowMessageBox(config); err != nil {
    gui.c.Log.Error("Failed to show message box:", err)
    return fmt.Errorf("show message box: %w", err)
}

// ❌ 不好的错误处理
gui.ShowMessageBox(config)  // 忽略错误
```

### 2. 资源清理

```go
// ✅ 使用 defer 确保清理
pb := gui.ShowProgressBar(config)
defer pb.Close()

// 执行操作
doSomething()
```

### 3. 并发安全

```go
// ✅ 使用 gui.g.Update 更新 UI
go func() {
    result := doWork()
    gui.g.Update(func(*gocui.Gui) error {
        gui.updateView(result)
        return nil
    })
}()
```

### 4. 上下文传递

```go
// ✅ 传递上下文
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

response, err := provider.Complete(ctx, prompt)
```

---

## 相关文档

- [项目概述](./PROJECT_OVERVIEW.md)
- [开发指南](./DEVELOPMENT_GUIDE.md)
- [架构设计](./ARCHITECTURE.md)

---

**版本**: v1.0.0
**最后更新**: 2024
**状态**: ✅ 完整
