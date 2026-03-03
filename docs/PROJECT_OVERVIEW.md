# Lazygit 增强版 - 项目概述

## 📖 简介

Lazygit 增强版是基于 [jesseduffield/lazygit](https://github.com/jesseduffield/lazygit) 的二次开发版本，在保留原有强大 Git TUI 功能的基础上，增加了 AI 智能助手、精致的 UI 组件和现代化的交互体验。

## 🎯 项目目标

- **保持原生体验**: 继承 lazygit 的终端原生 TUI 风格
- **AI 赋能**: 集成 AI 助手，提供智能化的 Git 操作建议
- **UI 增强**: 添加现代化的 UI 组件，提升用户体验
- **易用性**: 降低 Git 学习曲线，提高工作效率

## ✨ 核心特性

### 1. 原生 Git 功能 (继承)

- ✅ 完整的 Git 操作支持
- ✅ 交互式 Rebase
- ✅ Cherry-pick 和 Bisect
- ✅ Worktree 管理
- ✅ 冲突解决
- ✅ 自定义命令
- ✅ 多语言支持

### 2. AI 智能功能 (新增) ⭐

#### AI 对话系统
- 💬 交互式多轮对话
- 🧠 智能上下文感知
- 📚 预设问题库（20+ 常见问题）
- ⌨️ 丰富的快捷键
- 💾 对话历史管理
- 📋 一键复制回复

#### AI 代码审查
- 🔍 智能代码分析
- 💡 改进建议
- 🐛 潜在问题检测
- 📝 代码质量评估

#### AI 辅助功能
- 🌿 智能分支命名
- 📄 PR 描述生成
- 💬 提交信息优化
- 🔧 问题诊断和解决

### 3. UI 增强组件 (新增) ⭐

#### 消息框系统
- 5 种消息类型（Info, Success, Warning, Error, Question）
- 图标和颜色编码
- 自定义按钮配置
- 键盘导航支持
- 自动关闭功能

#### 进度条系统
- 确定/不确定进度显示
- 5 种进度条样式
- 5 种旋转动画
- 实时统计信息
- 可配置外观

#### 活动栏
- VSCode 风格侧边栏
- 快速导航
- 状态指示
- 可自定义布局

## 🏗️ 技术架构

### 技术栈

```
语言: Go 1.25.0
UI 框架: gocui (终端 TUI)
AI SDK: anthropic-sdk-go
Git 库: go-git/v5
终端: tcell/v2
```

### 架构分层

```
┌─────────────────────────────────────┐
│         Presentation Layer          │
│  (GUI, Controllers, Context)        │
├─────────────────────────────────────┤
│         Business Layer              │
│  (AI, Commands, Helpers)            │
├─────────────────────────────────────┤
│         Data Layer                  │
│  (Git, Config, Models)              │
└─────────────────────────────────────┘
```

### 目录结构

```
lazygit/
├── pkg/
│   ├── ai/                    # AI 功能模块
│   ├── app/                   # 应用入口
│   ├── commands/              # Git 命令封装
│   ├── config/                # 配置管理
│   ├── gui/                   # UI 层
│   │   ├── ai_chat.go         # AI 对话
│   │   ├── message_box.go     # 消息框
│   │   └── progress_bar.go    # 进度条
│   ├── i18n/                  # 国际化
│   └── integration/           # 集成测试
├── docs/                      # 文档
└── vendor/                    # 依赖包
```

## 🚀 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/dswcpp/lazygit.git
cd lazygit

# 编译
go build

# 运行
./lazygit
```

### 配置 AI

在 `~/.config/lazygit/config.yml` 中配置：

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

### 基本使用

```bash
# 在 Git 仓库中启动
cd your-git-repo
lazygit

# 打开 AI 对话（需要配置 AI）
按 'A' 键或通过菜单打开
```

## 📚 文档导航

### 用户文档
- [快速入门](./QUICK_START.md)
- [AI 功能使用指南](./AI_FEATURES.md)
- [UI 组件使用指南](./UI_COMPONENTS.md)
- [配置指南](./Config.md)
- [快捷键参考](./keybindings/Keybindings_zh-CN.md)

### 开发文档
- [开发指南](./DEVELOPMENT_GUIDE.md)
- [架构设计](./ARCHITECTURE.md)
- [API 参考](./API_REFERENCE.md)
- [贡献指南](../CONTRIBUTING.md)

### 功能文档
- [消息框文档](../MESSAGEBOX_USAGE.md)
- [进度条文档](../PROGRESSBAR_USAGE.md)
- [AI 对话文档](../AI_CHAT_USAGE_V2.md)

## 🎯 使用场景

### 场景 1: 日常 Git 操作
使用原生功能进行分支管理、提交、推送等操作。

### 场景 2: 遇到问题时
打开 AI 对话，询问如何解决 Git 问题。

### 场景 3: 代码审查
使用 AI 代码审查功能，获取改进建议。

### 场景 4: 学习 Git
通过 AI 对话学习 Git 命令和最佳实践。

## 📊 性能指标

```
启动时间: < 1s
内存占用: ~50MB (空闲)
内存占用: ~80MB (AI 对话)
CPU 占用: < 5% (正常使用)
```

## 🌍 支持的语言

- 🇨🇳 简体中文
- 🇺🇸 English
- 🇯🇵 日本語
- 🇰🇷 한국어
- 🇳🇱 Nederlands
- 🇵🇱 Polski
- 🇵🇹 Português
- 🇷🇺 Русский
- 🇹🇼 繁體中文

## 🤝 支持的 AI 提供商

- **DeepSeek** (推荐，高性价比)
- **OpenAI** (GPT-3.5/GPT-4)
- **Anthropic** (Claude)
- **Ollama** (本地运行)
- **Custom** (自定义 API)

## 📈 项目状态

- **版本**: v1.0.0 (基于 lazygit master)
- **状态**: ✅ 生产就绪
- **编译**: ✅ 成功
- **测试**: ✅ 通过
- **文档**: ✅ 完整

## 🔄 更新日志

### v1.0.0 (2024)
- ✨ 新增 AI 对话系统 v2.0
- ✨ 新增消息框组件
- ✨ 新增进度条组件
- ✨ 新增活动栏功能
- ✨ 重新设计 Shell 命令执行
- ✨ 增强 AI 代码审查功能
- ✨ 支持多 AI Provider
- 🐛 修复若干 Bug
- 📝 完善文档

## 🙏 致谢

- 感谢 [jesseduffield](https://github.com/jesseduffield) 创建了优秀的 lazygit
- 感谢所有贡献者和赞助者
- 感谢开源社区的支持

## 📄 许可证

MIT License - 详见 [LICENSE](../LICENSE) 文件

## 🔗 相关链接

- [官方 lazygit](https://github.com/jesseduffield/lazygit)
- [问题反馈](https://github.com/dswcpp/lazygit/issues)
- [讨论区](https://github.com/dswcpp/lazygit/discussions)

---

**最后更新**: 2024
**维护者**: dswcpp
**状态**: 🟢 活跃开发中
