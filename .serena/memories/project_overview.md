# Lazygit 项目概览

## 项目目的
Lazygit 是一个简单的 Git 命令终端 UI 工具，旨在成为最令人愉悦的 Git UI。它通过直观的终端界面简化了 Git 操作，让用户无需记忆复杂的 Git 命令。

## 核心功能
- 交互式暂存（stage individual lines）
- 交互式 rebase
- Cherry-pick
- Bisect
- Worktrees 管理
- 自定义命令系统
- Undo/Redo 功能
- 提交图可视化

## 技术栈
- **语言**: Go 1.25.0
- **主要依赖**:
  - `github.com/jesseduffield/gocui` - UI 渲染和用户输入处理
  - `github.com/jesseduffield/go-git/v5` - Git 操作
  - `github.com/gdamore/tcell/v2` - 终端控制
  - `github.com/jesseduffield/lazycore` - lazygit 和 lazydocker 共享功能
  - `github.com/sirupsen/logrus` - 日志
  - `github.com/spf13/afero` - 文件系统抽象

## 项目结构
```
lazygit/
├── cmd/                    # 命令行入口
│   ├── i18n/              # 国际化工具
│   └── integration_test/  # 集成测试入口
├── pkg/                    # 主要代码包
│   ├── app/               # 应用入口和核心逻辑
│   ├── commands/          # Git 命令封装
│   ├── gui/               # GUI 相关代码
│   ├── i18n/              # 国际化
│   ├── integration/       # 集成测试
│   └── utils/             # 工具函数
├── docs/                   # 文档
├── test/                   # 测试文件
├── scripts/                # 脚本工具
└── vendor/                 # 依赖包
```

## 设计原则
1. **可发现性** (Discoverability) - 让新用户容易上手
2. **简单性** (Simplicity) - 常见用例应该简单
3. **安全性** (Safety) - 保护用户不犯错
4. **强大性** (Power) - 支持复杂用例
5. **速度** (Speed) - 快速响应和启动
6. **符合 Git** (Conformity with git) - 遵循 Git 的行为
7. **考虑代码库** (Think of the codebase) - 保持代码库可维护性
