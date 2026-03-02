# 任务完成检查清单

当完成一个开发任务后，应该执行以下步骤：

## 1. 代码格式化
```bash
make format
```
或
```bash
gofumpt -l -w .
```

确保所有代码符合 gofumpt 的格式要求。

## 2. 代码检查 (Linting)
```bash
make lint
```
或
```bash
./scripts/golangci-lint-shim.sh run
```

修复所有 linting 错误和警告。

## 3. 运行测试

### 单元测试
```bash
make unit-test
```
确保所有单元测试通过。

### 集成测试（如果相关）
```bash
make integration-test-all
```
如果修改了核心功能，运行集成测试。

## 4. 生成自动生成的文件（如果需要）
```bash
make generate
```

如果添加了新的测试或修改了需要生成文档的内容，运行此命令。

## 5. 更新文档（如果需要）
- 如果添加了新功能，更新 `README.md`
- 如果修改了配置选项，更新 `docs/Config.md`
- 如果添加了新的自定义命令，更新 `docs/Custom_Command_Keybindings.md`
- 如果修改了键绑定，更新 `docs/keybindings/`

## 6. 国际化（如果添加了新文本）
- 在 `pkg/i18n/english.go` 中添加新的翻译字段
- 翻译会通过 Crowdin 由社区完成

## 7. 提交前检查
- [ ] 代码已格式化
- [ ] 通过所有 linting 检查
- [ ] 通过所有相关测试
- [ ] 添加了必要的测试（如果是新功能或 bug 修复）
- [ ] 更新了相关文档
- [ ] 提交信息清晰且有意义
- [ ] 使用功能分支（不是 master 分支）

## 8. 提交
```bash
git add <files>
git commit -m "type: description"
```

提交信息应该：
- 解释为什么做这个变更
- 遵循良好的提交信息规范
- 如果是多个独立的变更，分成多个提交

## 9. 推送和创建 PR
```bash
git push -u origin <branch-name>
```

然后在 GitHub 上创建 Pull Request。

## 注意事项

### 重构和行为变更
- 应该在不同的提交中
- 使用 fixup commits 在 review 期间迭代

### 测试要求
- 如果添加了应该被测试的代码，必须添加测试
- 单元测试文件以 `_test.go` 结尾
- 集成测试参见 `pkg/integration/README.md`

### 性能考虑
- 启动应该快速
- 如果在启动时运行慢的操作，使其非阻塞
- 考虑命令的性能影响

### 调试
如果需要调试：
1. 终端 1：`make run-debug` 或 `go run main.go -debug`
2. 终端 2：`make print-log` 或 `lazygit --logs`
3. 在代码中使用 `gui.Log.Warn("message")` 或 `self.c.Log.Warn("message")`
4. 设置日志级别：`LOG_LEVEL=warn go run main.go -debug`
