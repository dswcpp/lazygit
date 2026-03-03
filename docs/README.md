# Lazygit 增强版文档中心

欢迎来到 Lazygit 增强版的文档中心！这里包含了项目的完整文档。

## 📚 文档导航

### 🚀 快速开始
- **[项目概述](./PROJECT_OVERVIEW.md)** - 了解项目的核心特性和技术架构
- **[快速入门](../QUICK_START.md)** - 5 分钟快速上手指南

### ✨ 新增功能文档
- **[AI 功能指南](./AI_FEATURES.md)** - AI 对话、代码审查、智能辅助
- **[UI 组件指南](./UI_COMPONENTS.md)** - 消息框、进度条、活动栏

### 🛠️ 开发文档
- **[开发指南](./DEVELOPMENT_GUIDE.md)** - 开发环境搭建、代码规范、测试
- **[架构设计](./ARCHITECTURE.md)** - 系统架构、设计模式、扩展性
- **[API 参考](./API_REFERENCE.md)** - 完整的 API 文档

### 📖 原生功能文档
- [配置指南](./Config.md) - 完整的配置选项说明
- [自定义命令](./Custom_Command_Keybindings.md) - 自定义命令和快捷键
- [自定义分页器](./Custom_Pagers.md) - 配置自定义分页器
- [Fixup 提交](./Fixup_Commits.md) - 使用 fixup 提交
- [范围选择](./Range_Select.md) - 范围选择功能
- [搜索功能](./Searching.md) - 搜索和过滤
- [堆叠分支](./Stacked_Branches.md) - 堆叠分支工作流
- [撤销操作](./Undoing.md) - 撤销和重做

### ⌨️ 快捷键参考
- [中文快捷键](./keybindings/Keybindings_zh-CN.md)
- [English Keybindings](./keybindings/Keybindings_en.md)
- [日本語キーバインド](./keybindings/Keybindings_ja.md)
- [한국어 키바인딩](./keybindings/Keybindings_ko.md)

### 👨‍💻 开发者文档
- [代码库指南](./dev/Codebase_Guide.md)
- [集成测试](./dev/Integration_Tests.md)
- [性能分析](./dev/Profiling.md)

---

## 🎯 按场景查找

### 我想学习如何使用
→ 从 [项目概述](./PROJECT_OVERVIEW.md) 开始，然后查看 [快速入门](../QUICK_START.md)

### 我想使用 AI 功能
→ 查看 [AI 功能指南](./AI_FEATURES.md)

### 我想自定义界面
→ 查看 [UI 组件指南](./UI_COMPONENTS.md) 和 [配置指南](./Config.md)

### 我想参与开发
→ 查看 [开发指南](./DEVELOPMENT_GUIDE.md) 和 [架构设计](./ARCHITECTURE.md)

### 我想查看 API
→ 查看 [API 参考](./API_REFERENCE.md)

---

## 🆕 新增功能亮点

### AI 功能
- 🤖 **AI 对话系统** - 交互式多轮对话，智能上下文感知
- 🔍 **AI 代码审查** - 智能代码分析和改进建议
- 💡 **AI 辅助功能** - 分支命名、PR 描述生成

### UI 增强
- 📬 **消息框系统** - 5 种消息类型，自定义按钮
- 📊 **进度条系统** - 确定/不确定进度，多种样式
- 📱 **活动栏** - VSCode 风格侧边栏导航

---

## 📊 文档统计

| 类别 | 文档数量 | 状态 |
|------|---------|------|
| 快速开始 | 2 | ✅ |
| 功能文档 | 3 | ✅ |
| 开发文档 | 3 | ✅ |
| 原生功能 | 8 | ✅ |
| 快捷键 | 9 | ✅ |
| 开发者 | 3 | ✅ |

**总计**: 28 篇文档

---

## 📖 文档版本

- **当前版本**: v1.0.0
- **最后更新**: 2024-03-03
- **状态**: ✅ 完整

---

## 🔗 相关链接

- [GitHub 仓库](https://github.com/dswcpp/lazygit)
- [问题反馈](https://github.com/dswcpp/lazygit/issues)
- [原版 Lazygit](https://github.com/jesseduffield/lazygit)

---

## 💡 贡献文档

发现文档问题或想要改进？欢迎提交 PR！

1. Fork 项目
2. 创建文档分支
3. 修改或添加文档
4. 提交 PR

---

**维护者**: dswcpp
**许可证**: MIT
**语言**: 简体中文 / English

### 📂 文档目录结构

```
docs/
├── README.md                      # 文档中心首页
├── PROJECT_OVERVIEW.md            # 项目概述
├── AI_FEATURES.md                 # AI 功能指南
├── UI_COMPONENTS.md               # UI 组件指南
├── DEVELOPMENT_GUIDE.md           # 开发指南
├── ARCHITECTURE.md                # 架构设计
├── API_REFERENCE.md               # API 参考
├── Config.md                      # 配置指南
├── features/                      # 功能文档
│   ├── README.md
│   ├── AI_CHAT_USAGE_V2.md       # AI 对话使用指南
│   ├── UI_FEATURES_SUMMARY.md    # UI 功能总结
│   ├── README_UI_FEATURES.md     # UI 功能概览
│   └── QUICK_START.md            # 快速入门
├── design/                        # 设计文档
│   ├── README.md
│   ├── AI_CHAT_DESIGN_V2.md      # AI 对话设计
│   ├── AI_DIALOGUE_IMPROVEMENT_PLAN.md
│   ├── SHELL_COMMAND_REDESIGN.md
│   └── SSH_SUPPORT_EVALUATION.md
├── archive/                       # 归档文档
│   ├── README.md
│   ├── AI_CHAT_USAGE.md          # AI 对话 v1（旧版）
│   ├── P0_IMPLEMENTATION_SUMMARY.md
│   └── HIGH_PRIORITY_FIXES_COMPLETED.md
├── keybindings/                   # 快捷键文档
│   └── Keybindings_*.md
└── dev/                           # 开发者文档
    ├── README.md
    ├── Codebase_Guide.md
    ├── Busy.md
    ├── Demo_Recordings.md
    ├── Profiling.md
    └── ...
```

