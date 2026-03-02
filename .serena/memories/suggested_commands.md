# Lazygit 开发常用命令

## 构建和运行

### 构建
```bash
make build
# 或直接使用 go
go build -gcflags='all=-N -l'
```

### 安装到系统
```bash
make install
# 或
go install
```

### 运行
```bash
make run
# 或直接运行
./lazygit
# 或使用 go run
go run main.go
```

### 调试模式运行
```bash
# 终端 1：运行 lazygit（调试模式）
make run-debug
# 或
go run main.go -debug

# 终端 2：查看日志
make print-log
# 或
go run main.go --logs
# 或（如果已安装）
lazygit --logs
```

### 设置日志级别
```bash
LOG_LEVEL=warn go run main.go -debug
```

## 测试

### 单元测试
```bash
make unit-test
# 或
go test ./... -short
```

### 集成测试
```bash
# 运行所有集成测试
make integration-test-all
# 或
go test pkg/integration/clients/*.go

# TUI 集成测试
make integration-test-tui [test-name]
# 或
go run cmd/integration_test/main.go tui [test-name]

# CLI 集成测试
make integration-test-cli [test-name]
# 或
go run cmd/integration_test/main.go cli [test-name]
```

### 运行所有测试
```bash
make test
```

## 代码质量

### 格式化代码
```bash
make format
# 或
gofumpt -l -w .
```

### 代码检查 (Linting)
```bash
make lint
# 或
./scripts/golangci-lint-shim.sh run
```

### 生成自动生成的文件
```bash
make generate
# 或
go generate ./...
```
这会生成：
- 测试列表
- Cheatsheets
- 其他自动生成的文件

## 依赖管理

### 更新 vendor 目录
```bash
make vendor
# 或
go mod vendor && go mod tidy
```

### 更新 gocui
```bash
make bump-gocui
# 或
./scripts/bump_gocui.sh
```

### 更新 lazycore
```bash
make bump-lazycore
# 或
./scripts/bump_lazycore.sh
```

## 演示录制
```bash
make record-demo [demo-name]
# 或
demo/record_demo.sh [demo-name]
```

## Windows 系统命令

### 基本命令
- `dir` - 列出目录内容（相当于 Linux 的 `ls`）
- `cd` - 切换目录
- `type` - 查看文件内容（相当于 Linux 的 `cat`）
- `findstr` - 搜索文本（相当于 Linux 的 `grep`）
- `where` - 查找可执行文件（相当于 Linux 的 `which`）

### Git 命令
- 所有标准 Git 命令在 Windows 上都可用（如果安装了 Git for Windows）
- 推荐使用 Git Bash 或 PowerShell

### Go 命令
- `go version` - 查看 Go 版本
- `go env` - 查看 Go 环境变量
- `go mod tidy` - 整理依赖
- `go mod vendor` - 创建 vendor 目录

## VSCode 任务
如果使用 VSCode，可以通过 `Cmd+Shift+P`（Mac）或 `Ctrl+Shift+P`（Windows）输入 "Run task" 来运行预定义的任务，例如：
- Bump lazycore
- 其他开发任务

## 开发环境

### 使用 Dev Container
- 需要 Docker 和 VSCode Dev Containers 扩展
- 在 VSCode 中打开项目，选择 "Reopen in Container"

### 使用 GitHub Codespace
- Fork 仓库后，点击创建 Codespace
- 使用 `go run main.go` 运行 lazygit

### 使用 Nix
```bash
# 进入开发环境
nix develop

# 构建
nix build

# 运行
nix run
```
