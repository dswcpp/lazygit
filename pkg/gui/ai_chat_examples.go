package gui

import (
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

// createAIChatTestMenu 创建 AI 对话测试菜单
func (gui *Gui) createAIChatTestMenu() error {
	menuItems := []*types.MenuItem{
		{
			Label: "打开 AI 对话",
			OnPress: func() error {
				return gui.ShowAIChat()
			},
		},
	}

	return gui.c.Menu(types.CreateMenuOptions{
		Title: "AI 功能测试菜单",
		Items: menuItems,
	})
}

// AI 对话使用示例

// Example1: 基本对话
// 用户可以直接调用 ShowAIChat() 打开对话界面
func (gui *Gui) ExampleOpenAIChat() error {
	return gui.ShowAIChat()
}

// Example2: 从快捷键打开
// 在 keybindings 中绑定快捷键
// {
//     Key: 'A',
//     Handler: gui.ShowAIChat,
//     Description: "打开 AI 对话",
// }

// Example3: 从菜单打开
// menuItems := []*types.MenuItem{
//     {
//         Label: "AI 对话",
//         OnPress: func() error {
//             return gui.ShowAIChat()
//         },
//     },
// }

// Example4: 检查 AI 是否启用
func (gui *Gui) ExampleCheckAIEnabled() error {
	if gui.c.AIManager == nil {
		gui.ShowError(gui.c.Tr.AIChatNotEnabled, "请先在设置中启用并配置 AI 功能。")
		return nil
	}
	return gui.ShowAIChat()
}

// Example5: 集成到现有工作流
// 例如：在查看 diff 时提供 AI 分析选项
func (gui *Gui) ExampleAIAnalyzeDiff() error {
	// 获取当前 diff
	// diff := gui.getCurrentDiff()

	// 打开 AI 对话并预填充问题
	// return gui.ShowAIChatWithMessage("请分析这个 diff: " + diff)

	return gui.ShowAIChat()
}

// 实际应用场景示例

// Scenario1: Git 问题咨询
// 用户: "如何撤销最近的一次提交？"
// AI: "你可以使用以下命令：
//      - git reset --soft HEAD~1  # 保留更改
//      - git reset --hard HEAD~1  # 丢弃更改
//      在 lazygit 中，你可以在 commits 视图中选择提交，按 'g' 然后选择 reset 选项。"

// Scenario2: 代码审查建议
// 用户: "这个提交有什么问题吗？"
// AI: "根据当前提交信息，我注意到：
//      1. 提交信息不够详细
//      2. 包含了多个不相关的更改
//      建议：将这个提交拆分成多个小提交，每个提交只做一件事。"

// Scenario3: 分支管理建议
// 用户: "我应该如何管理我的分支？"
// AI: "Git 分支管理最佳实践：
//      1. main/master 分支保持稳定
//      2. 为每个功能创建独立的 feature 分支
//      3. 定期合并主分支的更新
//      4. 完成后通过 PR 合并回主分支
//      5. 合并后删除已完成的分支"

// Scenario4: 冲突解决指导
// 用户: "遇到合并冲突了，怎么办？"
// AI: "解决合并冲突的步骤：
//      1. 在 lazygit 中，冲突文件会标记为红色
//      2. 选择冲突文件，按 Enter 打开编辑器
//      3. 查找 <<<<<<<, =======, >>>>>>> 标记
//      4. 决定保留哪些代码，删除冲突标记
//      5. 保存文件后，在 lazygit 中 stage 该文件
//      6. 继续合并或提交"

// Scenario5: 性能优化建议
// 用户: "仓库太大了，如何优化？"
// AI: "优化大型仓库的方法：
//      1. 使用 git gc 清理不必要的文件
//      2. 使用 git prune 删除不可达对象
//      3. 考虑使用 Git LFS 管理大文件
//      4. 使用 .gitignore 排除不必要的文件
//      5. 定期清理旧的分支和标签"
