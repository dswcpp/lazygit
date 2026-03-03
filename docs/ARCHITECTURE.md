# 架构设计文档

## 📖 目录

- [系统架构](#系统架构)
- [分层设计](#分层设计)
- [核心模块](#核心模块)
- [数据流](#数据流)
- [设计模式](#设计模式)
- [扩展性设计](#扩展性设计)

---

## 系统架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────┐
│                    User Interface                        │
│              (Terminal TUI - gocui)                      │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────┴────────────────────────────────────┐
│                  Presentation Layer                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │ AI Chat  │  │ Message  │  │ Progress │             │
│  │          │  │ Box      │  │ Bar      │             │
│  └──────────┘  └──────────┘  └──────────┘             │
│  ┌──────────────────────────────────────┐              │
│  │        Controllers & Context         │              │
│  └──────────────────────────────────────┘              │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────┴────────────────────────────────────┐
│                   Business Layer                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │ AI       │  │ Commands │  │ Helpers  │             │
│  │ Providers│  │          │  │          │             │
│  └──────────┘  └──────────┘  └──────────┘             │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────┴────────────────────────────────────┐
│                     Data Layer                           │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │ Git      │  │ Config   │  │ Models   │             │
│  │ Commands │  │          │  │          │             │
│  └──────────┘  └──────────┘  └──────────┘             │
└─────────────────────────────────────────────────────────┘
```

### 技术栈

```
┌─────────────────────────────────────────┐
│ Language: Go 1.25.0                     │
├─────────────────────────────────────────┤
│ UI Framework: gocui                     │
│ Terminal: tcell/v2                      │
├─────────────────────────────────────────┤
│ Git Library: go-git/v5                  │
│ AI SDK: anthropic-sdk-go                │
├─────────────────────────────────────────┤
│ Config: YAML (gopkg.in/yaml.v3)        │
│ I18n: custom                            │
└─────────────────────────────────────────┘
```

---

## 分层设计

### 1. Presentation Layer (表示层)

**职责**: 处理用户交互和界面渲染

**组件**:
```
pkg/gui/
├── gui.go                 # GUI 主结构
├── ai_chat.go             # AI 对话界面
├── message_box.go         # 消息框
├── progress_bar.go        # 进度条
├── context/               # 上下文管理
│   ├── base_context.go
│   ├── branches_context.go
│   └── ...
└── controllers/           # 控制器
    ├── files_controller.go
    ├── branches_controller.go
    └── ...
```

**设计原则**:
- 单一职责：每个组件只负责一个功能
- 松耦合：通过接口与业务层交互
- 可测试：UI 逻辑与业务逻辑分离

### 2. Business Layer (业务层)

**职责**: 实现业务逻辑和规则

**组件**:
```
pkg/
├── ai/                    # AI 功能
│   ├── provider.go        # 接口定义
│   ├── deepseek.go        # DeepSeek 实现
│   ├── openai.go          # OpenAI 实现
│   └── anthropic.go       # Anthropic 实现
├── commands/              # Git 命令
│   ├── git_commands/      # Git 操作
│   └── oscommands/        # 系统命令
└── gui/controllers/helpers/  # 辅助函数
    ├── branches_helper.go
    ├── files_helper.go
    └── ...
```

**设计原则**:
- 业务逻辑集中
- 可复用性高
- 易于测试

### 3. Data Layer (数据层)

**职责**: 数据访问和持久化

**组件**:
```
pkg/
├── commands/
│   ├── git_commands/      # Git 数据访问
│   └── git_config/        # Git 配置
├── config/                # 应用配置
│   ├── user_config.go
│   └── app_config.go
└── commands/models/       # 数据模型
    ├── branch.go
    ├── commit.go
    └── ...
```

**设计原则**:
- 数据访问抽象
- 配置集中管理
- 模型定义清晰

---

## 核心模块

### 1. GUI 模块

#### 结构设计

```go
type Gui struct {
    *common.Common
    g          *gocui.Gui          // gocui 实例
    git        *commands.GitCommand
    os         *oscommands.OSCommand
    State      *GuiRepoState       // 仓库状态
    Config     config.AppConfigurer
    Mutexes    types.Mutexes       // 并发控制
    // ...
}
```

#### 生命周期

```
初始化 → 创建视图 → 设置键盘绑定 → 主循环 → 清理
  ↓         ↓            ↓            ↓        ↓
NewGui   setupViews  setupKeys   g.MainLoop  cleanup
```

#### 视图管理

```go
// 视图创建
func (gui *Gui) createView(name string, x0, y0, x1, y1 int) (*gocui.View, error) {
    v, err := gui.g.SetView(name, x0, y0, x1, y1, 0)
    if err != nil && err != gocui.ErrUnknownView {
        return nil, err
    }
    return v, nil
}

// 视图更新
func (gui *Gui) updateView(name string, content string) error {
    return gui.g.Update(func(*gocui.Gui) error {
        v, err := gui.g.View(name)
        if err != nil {
            return err
        }
        v.Clear()
        fmt.Fprint(v, content)
        return nil
    })
}
```

### 2. AI 模块

#### 接口设计

```go
// AIProvider 接口
type AIProvider interface {
    Complete(ctx context.Context, prompt string) (*AIResponse, error)
    CompleteStream(ctx context.Context, prompt string) (<-chan string, error)
}

// AIResponse 响应
type AIResponse struct {
    Content string
    Usage   Usage
}

// Usage 使用统计
type Usage struct {
    PromptTokens     int
    CompletionTokens int
    TotalTokens      int
}
```

#### 实现架构

```
┌─────────────────────────────────────────┐
│         AIProvider Interface            │
└────────────┬────────────────────────────┘
             │
    ┌────────┴────────┬────────┬──────────┐
    │                 │        │          │
┌───▼───┐      ┌─────▼──┐  ┌──▼────┐  ┌──▼────┐
│DeepSeek│      │OpenAI  │  │Claude │  │Ollama │
└────────┘      └────────┘  └───────┘  └───────┘
```

#### 配置管理

```go
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

### 3. Commands 模块

#### 命令封装

```go
// GitCommand 主结构
type GitCommand struct {
    Branch    *BranchCommands
    Commit    *CommitCommands
    File      *FileCommands
    Remote    *RemoteCommands
    Stash     *StashCommands
    Tag       *TagCommands
    Worktree  *WorktreeCommands
    // ...
}

// BranchCommands 分支命令
type BranchCommands struct {
    *common.Common
    gitCommon *GitCommon
}

func (c *BranchCommands) New(name string) error {
    return c.gitCommon.RunCommand("git branch %s", name)
}

func (c *BranchCommands) Delete(name string, force bool) error {
    flag := "-d"
    if force {
        flag = "-D"
    }
    return c.gitCommon.RunCommand("git branch %s %s", flag, name)
}
```

#### 命令执行流程

```
用户操作 → Controller → Helper → GitCommand → OSCommand → Git
   ↓          ↓           ↓          ↓            ↓          ↓
  按键      处理逻辑    业务逻辑   命令封装    系统调用   执行命令
```

---

## 数据流

### 1. 用户操作流

```
┌──────────┐
│  用户    │
│  按键    │
└────┬─────┘
     │
     ▼
┌──────────────┐
│ Key Binding  │
│ Handler      │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ Controller   │
│ Method       │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ Helper       │
│ Function     │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ Git Command  │
│ Execution    │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ State Update │
│ & Refresh    │
└──────────────┘
```

### 2. AI 对话流

```
┌──────────┐
│  用户    │
│  输入    │
└────┬─────┘
     │
     ▼
┌──────────────┐
│ AI Chat      │
│ Input View   │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ Build Prompt │
│ + Context    │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ AI Provider  │
│ API Call     │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ Parse        │
│ Response     │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ Update Chat  │
│ View         │
└──────────────┘
```

### 3. 状态管理流

```
┌──────────────┐
│ Git Command  │
│ Execution    │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ Load Data    │
│ from Git     │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ Update State │
│ (GuiRepoState)│
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ Refresh Views│
│ (Render)     │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ User Sees    │
│ Changes      │
└──────────────┘
```

---

## 设计模式

### 1. MVC 模式

```
Model (数据层)
  ├── commands/models/
  └── State

View (视图层)
  ├── gui/context/
  └── gui/*.go (UI 组件)

Controller (控制层)
  └── gui/controllers/
```

### 2. 策略模式 (AI Providers)

```go
// 策略接口
type AIProvider interface {
    Complete(ctx context.Context, prompt string) (*AIResponse, error)
}

// 具体策略
type DeepSeekProvider struct { /* ... */ }
type OpenAIProvider struct { /* ... */ }
type AnthropicProvider struct { /* ... */ }

// 上下文
type AIClient struct {
    provider AIProvider
}

func (c *AIClient) SetProvider(provider AIProvider) {
    c.provider = provider
}
```

### 3. 工厂模式 (UI 组件)

```go
// 工厂方法
func (gui *Gui) ShowMessageBox(config MessageBoxConfig) error {
    mb := &MessageBox{
        config: config,
        gui:    gui,
    }
    return mb.Show()
}

func (gui *Gui) ShowProgressBar(config ProgressBarConfig) *ProgressBar {
    pb := &ProgressBar{
        config: config,
        gui:    gui,
    }
    pb.Start()
    return pb
}
```

### 4. 观察者模式 (事件系统)

```go
// 事件类型
type EventType int

const (
    EventRefresh EventType = iota
    EventStateChange
    EventError
)

// 事件监听器
type EventListener interface {
    OnEvent(event Event)
}

// 事件管理器
type EventManager struct {
    listeners map[EventType][]EventListener
}

func (em *EventManager) Subscribe(eventType EventType, listener EventListener) {
    em.listeners[eventType] = append(em.listeners[eventType], listener)
}

func (em *EventManager) Publish(event Event) {
    for _, listener := range em.listeners[event.Type] {
        listener.OnEvent(event)
    }
}
```

### 5. 命令模式 (Git 操作)

```go
// 命令接口
type Command interface {
    Execute() error
    Undo() error
}

// 具体命令
type CommitCommand struct {
    message string
    files   []string
}

func (c *CommitCommand) Execute() error {
    // 执行提交
}

func (c *CommitCommand) Undo() error {
    // 撤销提交
}

// 命令调用者
type CommandInvoker struct {
    history []Command
}

func (ci *CommandInvoker) Execute(cmd Command) error {
    if err := cmd.Execute(); err != nil {
        return err
    }
    ci.history = append(ci.history, cmd)
    return nil
}
```

---

## 扩展性设计

### 1. 插件系统 (未来)

```go
// 插件接口
type Plugin interface {
    Name() string
    Version() string
    Init(gui *Gui) error
    OnEvent(event Event) error
}

// 插件管理器
type PluginManager struct {
    plugins map[string]Plugin
}

func (pm *PluginManager) Register(plugin Plugin) error {
    pm.plugins[plugin.Name()] = plugin
    return plugin.Init(gui)
}
```

### 2. 自定义命令

```yaml
customCommands:
  - key: 'c'
    command: 'git commit -m "{{.Form.Message}}"'
    context: 'files'
    prompts:
      - type: 'input'
        title: 'Commit message'
        key: 'Message'
```

### 3. 主题系统

```yaml
theme:
  activeBorderColor:
    - green
    - bold
  inactiveBorderColor:
    - white
  selectedLineBgColor:
    - blue
```

### 4. 钩子系统

```go
// 钩子类型
type HookType int

const (
    HookPreCommit HookType = iota
    HookPostCommit
    HookPrePush
    HookPostPush
)

// 钩子函数
type HookFunc func(context HookContext) error

// 钩子管理器
type HookManager struct {
    hooks map[HookType][]HookFunc
}

func (hm *HookManager) Register(hookType HookType, fn HookFunc) {
    hm.hooks[hookType] = append(hm.hooks[hookType], fn)
}

func (hm *HookManager) Execute(hookType HookType, context HookContext) error {
    for _, fn := range hm.hooks[hookType] {
        if err := fn(context); err != nil {
            return err
        }
    }
    return nil
}
```

---

## 性能优化

### 1. 懒加载

```go
// 延迟加载大型数据
type LazyLoader struct {
    data   interface{}
    loader func() (interface{}, error)
    loaded bool
}

func (ll *LazyLoader) Get() (interface{}, error) {
    if !ll.loaded {
        data, err := ll.loader()
        if err != nil {
            return nil, err
        }
        ll.data = data
        ll.loaded = true
    }
    return ll.data, nil
}
```

### 2. 缓存机制

```go
// 简单缓存
type Cache struct {
    data map[string]CacheEntry
    ttl  time.Duration
}

type CacheEntry struct {
    value      interface{}
    expiration time.Time
}

func (c *Cache) Get(key string) (interface{}, bool) {
    entry, ok := c.data[key]
    if !ok || time.Now().After(entry.expiration) {
        return nil, false
    }
    return entry.value, true
}

func (c *Cache) Set(key string, value interface{}) {
    c.data[key] = CacheEntry{
        value:      value,
        expiration: time.Now().Add(c.ttl),
    }
}
```

### 3. 并发控制

```go
// 使用 sync.Mutex 保护共享资源
type SafeState struct {
    mu    sync.RWMutex
    state map[string]interface{}
}

func (s *SafeState) Get(key string) interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.state[key]
}

func (s *SafeState) Set(key string, value interface{}) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.state[key] = value
}
```

---

## 安全性设计

### 1. 输入验证

```go
func validateInput(input string) error {
    // 检查长度
    if len(input) > MaxInputLength {
        return errors.New("input too long")
    }

    // 检查特殊字符
    if containsDangerousChars(input) {
        return errors.New("invalid characters")
    }

    return nil
}
```

### 2. 命令注入防护

```go
// 使用参数化命令
func (c *GitCommand) SafeCommand(args ...string) error {
    cmd := exec.Command("git", args...)
    return cmd.Run()
}

// 避免字符串拼接
// ❌ 不安全
cmd := fmt.Sprintf("git commit -m '%s'", message)

// ✅ 安全
cmd := exec.Command("git", "commit", "-m", message)
```

### 3. 敏感信息保护

```go
// 不在日志中输出敏感信息
func (c *AIClient) logRequest(prompt string) {
    // 移除敏感信息
    sanitized := removeSensitiveInfo(prompt)
    log.Info("AI request:", sanitized)
}

// 配置文件权限检查
func checkConfigPermissions(path string) error {
    info, err := os.Stat(path)
    if err != nil {
        return err
    }

    // 检查权限
    if info.Mode().Perm() > 0600 {
        return errors.New("config file permissions too open")
    }

    return nil
}
```

---

## 相关文档

- [项目概述](./PROJECT_OVERVIEW.md)
- [开发指南](./DEVELOPMENT_GUIDE.md)
- [API 参考](./API_REFERENCE.md)

---

**版本**: v1.0.0
**最后更新**: 2024
**状态**: ✅ 完整
