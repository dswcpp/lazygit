# Lazygit 代码风格和约定

## Go 语言风格

### 接收者命名
- **允许使用 `self`** 作为结构体方法的接收者名称（与标准 Go 惯例不同）
- 示例：`func (self *MyStruct) MyMethod() {}`

### 接口命名
- 允许使用 'I' 前缀命名接口（如 `IMyInterface`），而不是标准的 'er' 后缀
- 当接口有多个方法时，这种方式更清晰

### 接口实现声明
- 显式声明结构体实现接口：
```go
var _ MyInterface = &MyStruct{}
```
- 这使意图更清晰，如果未满足接口会在编译时报错

## 代码格式化

### 工具
- **gofumpt** - 比 gofmt 更严格的格式化工具
- 安装：`go install mvdan.cc/gofumpt@latest`
- 运行：`gofumpt -l -w .`

### VSCode 配置
在 `.vscode/settings.json` 中设置：
```json
{
  "gopls": {
    "formatting.gofumpt": true
  }
}
```

## 代码检查 (Linting)

### 工具
- **golangci-lint** - 集成多个 linter
- 配置文件：`.golangci.yml`

### 启用的 Linters
- copyloopvar - 检查循环变量复制
- errorlint - 错误处理检查
- exhaustive - 穷尽性检查
- intrange - 整数范围检查
- makezero - make 零值检查
- nakedret - 裸返回检查（禁止所有裸返回）
- prealloc - 预分配检查
- revive - 代码质量检查
- thelper - 测试辅助函数检查
- unconvert - 不必要的类型转换
- unparam - 未使用的参数
- wastedassign - 浪费的赋值

### 特殊规则
- 允许使用 `self` 或 `this` 作为接收者名称
- 错误字符串可以大写（与标准 Go 不同）
- 允许使用 `!(a && b)` 而不强制改为 `!a || !b`

## 文本约定

### 大小写
- **使用 Sentence case**（句子大小写）
- **不使用 Title Case**（标题大小写）
- **不使用 all-lowercase**（全小写）
- 示例：
  - ✅ "Stage individual lines"
  - ❌ "Stage Individual Lines"
  - ❌ "stage individual lines"

## 国际化 (i18n)

### 添加新文本
1. 在 `pkg/i18n/english.go` 的 `TranslationSet` 结构体中添加新字段
2. 在 `EnglishTranslationSet()` 方法中添加实际内容
3. 通过 `gui.Tr.YourNewText` 或 `self.c.Tr.YourNewText` 访问

### 翻译管理
- 翻译通过 [Crowdin](https://crowdin.com/project/lazygit/) 管理
- 不要直接编辑 `pkg/i18n/translations/` 中的翻译文件

## 字体要求
- 开发环境需要使用 [Nerd Fonts](https://www.nerdfonts.com)
- 代码中（特别是测试）可能包含 Nerd Font 图标字符

## 提交历史

### 原则
- 保持干净和有用的提交历史
- 不会在合并时 squash 提交
- 重构和行为变更应该在不同的提交中
- 追求最小化提交：每个独立的变更应该在自己的提交中
- 使用 fixup commits 在 review 期间迭代

### 提交信息
- 遵循 [好的提交信息规范](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)
- 解释为什么做这个变更，而不仅仅是做了什么

## 分支策略
- 从 `master` 分支创建功能分支
- **不要从 fork 的 master 分支提交 PR**
- 使用功能分支，以便维护者可以推送变更
