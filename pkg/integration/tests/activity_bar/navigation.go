package activity_bar

import (
	"github.com/dswcpp/lazygit/pkg/config"
	. "github.com/dswcpp/lazygit/pkg/integration/components"
)

// TODO: 完成 Activity Bar 的集成测试
// 需要验证以下功能：
// 1. Activity Bar 视图正确显示
// 2. 导航项点击可以切换面板（status, files, branches, commits, stash）
// 3. 当前面板显示蓝色圆点指示器
// 4. Git 操作（pull, push, fetch, stash, merge, rebase）正确执行
// 5. 操作进行中显示 spinner 动画
// 6. 某些操作在特定状态下显示为禁用（灰色）
// 7. 自定义配置正确加载
// 8. 图标样式（nerd/emoji/ascii）正确显示

var Navigation = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "使用 Activity Bar 导航到不同的面板",
	ExtraCmdArgs: []string{},
	Skip:         true, // TODO: 在 Activity Bar UI 完全集成后启用此测试
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ActivityBar.Show = true
		config.GetUserConfig().Gui.ActivityBar.Width = 3
		config.GetUserConfig().Gui.ActivityBar.IconStyle = "ascii"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file1.txt", "content1").
			Commit("initial commit").
			NewBranch("feature").
			Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// TODO: 实现测试逻辑
		// 示例：
		// t.Views().ActivityBar().
		// 	Focus().
		// 	SelectNextItem(). // 选择 "files"
		// 	Press(keys.Universal.Select)
		//
		// t.Views().Files().
		// 	IsFocused().
		// 	Lines(
		// 		Contains("file1.txt"),
		// 	)
	},
})

var GitOperations = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "使用 Activity Bar 执行 Git 操作",
	ExtraCmdArgs: []string{},
	Skip:         true, // TODO: 在 Activity Bar UI 完全集成后启用此测试
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ActivityBar.Show = true
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file1.txt", "content1").
			Commit("initial commit").
			AddRemote("origin", "https://github.com/user/repo.git")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// TODO: 测试 pull, push, fetch, stash 等操作
	},
})

var VisualStates = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "验证 Activity Bar 的视觉状态（当前面板指示器、spinner、禁用状态）",
	ExtraCmdArgs: []string{},
	Skip:         true, // TODO: 在 Activity Bar UI 完全集成后启用此测试
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ActivityBar.Show = true
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file1.txt", "content1").
			Commit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// TODO: 验证视觉反馈
		// - 当前面板的蓝色圆点
		// - 操作进行中的 spinner
		// - 禁用操作的灰色显示
	},
})

var CustomConfiguration = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "验证自定义 Activity Bar 配置",
	ExtraCmdArgs: []string{},
	Skip:         true, // TODO: 在 Activity Bar UI 完全集成后启用此测试
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ActivityBar.Show = true
		config.GetUserConfig().Gui.ActivityBar.IconStyle = "emoji"
		// TODO: 添加自定义项目列表配置
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1.txt", "content1").Commit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// TODO: 验证自定义配置生效
	},
})
