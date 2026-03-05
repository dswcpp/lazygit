package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

// CreateUIFeaturesTestMenu 创建 UI 功能测试总菜单
func (gui *Gui) CreateUIFeaturesTestMenu() error {
	menuItems := []*types.MenuItem{
		{
			Label: "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━",
		},
		{
			Label: "💬 AI 对话功能",
			OnPress: func() error {
				return gui.createAIChatTestMenu()
			},
			Key:       'a',
			OpensMenu: true,
		},
		{
			Label: "📋 消息框功能",
			OnPress: func() error {
				return gui.createMessageBoxTestMenu()
			},
			Key:       'm',
			OpensMenu: true,
		},
		{
			Label: "📊 进度条功能",
			OnPress: func() error {
				return gui.createProgressBarTestMenu()
			},
			Key:       'p',
			OpensMenu: true,
		},
		{
			Label: "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━",
		},
		{
			Label: "🎨 综合演示",
			OnPress: func() error {
				return gui.runComprehensiveDemo()
			},
			Key: 'd',
		},
		{
			Label: "🔄 实际场景演示",
			OnPress: func() error {
				return gui.createRealWorldScenariosMenu()
			},
			Key:       'r',
			OpensMenu: true,
		},
		{
			Label: "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━",
		},
		{
			Label: "📖 快速入门指南",
			OnPress: func() error {
				return gui.showQuickStartGuide()
			},
			Key: 'q',
		},
		{
			Label: "❓ 功能帮助",
			OnPress: func() error {
				return gui.showFeaturesHelp()
			},
			Key: 'h',
		},
	}

	return gui.c.Menu(types.CreateMenuOptions{
		Title: "🎨 UI 功能测试中心",
		Items: menuItems,
	})
}

// runComprehensiveDemo 运行综合演示
func (gui *Gui) runComprehensiveDemo() error {
	// 演示流程：消息框 -> 进度条 -> AI 对话

	// 1. 显示欢迎消息
	gui.ShowInfo(
		"综合演示",
		"欢迎使用 lazygit UI 功能演示！",
		"接下来将依次展示：\n"+
			"1. 消息框功能\n"+
			"2. 进度条功能\n"+
			"3. AI 对话功能",
	)

	// 2. 演示不同类型的消息框
	go func() {
		// 等待用户关闭欢迎消息
		// 这里简化处理，实际应该等待用户操作
	}()

	return nil
}

// createRealWorldScenariosMenu 创建实际场景演示菜单
func (gui *Gui) createRealWorldScenariosMenu() error {
	menuItems := []*types.MenuItem{
		{
			Label: "场景 1: Git Push 流程",
			OnPress: func() error {
				return gui.demoGitPushScenario()
			},
			Key: '1',
		},
		{
			Label: "场景 2: 合并冲突处理",
			OnPress: func() error {
				return gui.demoMergeConflictScenario()
			},
			Key: '2',
		},
		{
			Label: "场景 3: 分支删除确认",
			OnPress: func() error {
				return gui.demoBranchDeleteScenario()
			},
			Key: '3',
		},
		{
			Label: "场景 4: AI 辅助问题解决",
			OnPress: func() error {
				return gui.demoAIAssistScenario()
			},
			Key: '4',
		},
		{
			Label: "场景 5: 大文件操作",
			OnPress: func() error {
				return gui.demoLargeFileScenario()
			},
			Key: '5',
		},
	}

	return gui.c.Menu(types.CreateMenuOptions{
		Title: "🔄 实际场景演示",
		Items: menuItems,
	})
}

// demoGitPushScenario 演示 Git Push 场景
func (gui *Gui) demoGitPushScenario() error {
	// 1. 显示确认对话框
	gui.ShowConfirm(
		"确认推送",
		"确定要推送到远程仓库 'origin/main' 吗？",
		func() {
			// 2. 显示进度条
			pb := gui.ShowProgressBar(ProgressBarConfig{
				Title:          "正在推送到远程仓库...",
				Total:          10 * 1024 * 1024, // 10 MB
				ShowPercentage: true,
				ShowStats:      true,
				Style:          ProgressBarStyleBlock,
			})

			// 3. 模拟推送过程
			go func() {
				for i := int64(0); i <= pb.config.Total; i += 256 * 1024 {
					pb.Update(i, "")
					// time.Sleep(100 * time.Millisecond)
				}
				pb.Close()

				// 4. 显示成功消息
				gui.g.Update(func(*gocui.Gui) error {
					gui.ShowSuccess(
						"推送成功",
						"已成功推送到 origin/main",
						"提交数: 3\n传输数据: 10.0 MB\n耗时: 2.5 秒",
					)
					return nil
				})
			}()
		},
	)

	return nil
}

// demoMergeConflictScenario 演示合并冲突场景
func (gui *Gui) demoMergeConflictScenario() error {
	conflictFiles := []string{
		"src/main.go",
		"pkg/config/config.go",
		"README.md",
	}

	details := "冲突文件:\n"
	for _, file := range conflictFiles {
		details += "  - " + file + "\n"
	}

	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeWarning,
		Title:   "合并冲突",
		Message: "合并 'feature-branch' 到 'main' 时发现冲突。",
		Details: details,
		Buttons: []string{"解决冲突", "中止合并", "查看详情"},
	}, func(buttonIndex int) {
		switch buttonIndex {
		case 0:
			gui.c.Toast("打开冲突解决界面...")
		case 1:
			gui.ShowConfirm(
				"确认中止",
				"确定要中止合并吗？所有合并进度将丢失。",
				func() {
					gui.c.Toast("已中止合并")
				},
			)
		case 2:
			// 打开 AI 对话询问如何解决
			gui.ShowAIChat()
		}
	})

	return nil
}

// demoBranchDeleteScenario 演示分支删除场景
func (gui *Gui) demoBranchDeleteScenario() error {
	branchName := "feature-old-feature"

	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeQuestion,
		Title:   "确认删除分支",
		Message: "确定要删除分支 '" + branchName + "' 吗？",
		Details: "此分支包含:\n" +
			"  • 15 个提交\n" +
			"  • 最后更新: 30 天前\n" +
			"  • 已合并到 main\n" +
			"  • 远程分支也将被删除",
		Buttons: []string{"仅删除本地", "删除本地和远程", "取消"},
	}, func(buttonIndex int) {
		switch buttonIndex {
		case 0:
			gui.c.Toast("已删除本地分支")
			gui.ShowSuccess(
				"删除成功",
				"本地分支 '"+branchName+"' 已删除",
				"远程分支保留",
			)
		case 1:
			// 显示进度
			pb := gui.ShowProgressBar(ProgressBarConfig{
				Title:         "正在删除分支...",
				Message:       "正在删除远程分支...",
				Indeterminate: true,
				SpinnerStyle:  SpinnerStyleBraille,
			})

			go func() {
				// time.Sleep(2 * time.Second)
				pb.Close()

				gui.g.Update(func(*gocui.Gui) error {
					gui.ShowSuccess(
						"删除成功",
						"分支 '"+branchName+"' 已完全删除",
						"本地和远程分支都已删除",
					)
					return nil
				})
			}()
		}
	})

	return nil
}

// demoAIAssistScenario 演示 AI 辅助场景
func (gui *Gui) demoAIAssistScenario() error {
	if gui.c.AIManager == nil {
		gui.ShowError(
			"AI 未启用",
			"此演示需要启用 AI 功能。",
			"请在设置中配置 AI 后再试。",
		)
		return nil
	}

	gui.ShowInfo(
		"AI 辅助演示",
		"即将打开 AI 对话，你可以询问任何 Git 相关问题。",
		"示例问题:\n"+
			"  • 如何撤销提交？\n"+
			"  • 如何解决合并冲突？\n"+
			"  • 如何清理大型仓库？\n\n"+
			"提示: 按 Ctrl+P 可以查看预设问题",
	)

	// 延迟打开 AI 对话
	go func() {
		// time.Sleep(1 * time.Second)
		gui.g.Update(func(*gocui.Gui) error {
			return gui.ShowAIChat()
		})
	}()

	return nil
}

// demoLargeFileScenario 演示大文件操作场景
func (gui *Gui) demoLargeFileScenario() error {
	gui.ShowWarning(
		"大文件检测",
		"检测到大文件 'assets/video.mp4' (250 MB)",
		"建议使用 Git LFS 管理大文件。\n\n"+
			"是否继续添加到 Git？",
	)

	// 这里可以添加更多交互
	return nil
}

// showQuickStartGuide 显示快速入门指南
func (gui *Gui) showQuickStartGuide() error {
	guide := `
╔═══════════════════════════════════════════════════════════╗
║              🚀 UI 功能快速入门指南                       ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  📋 消息框 (MessageBox)                                   ║
║  ─────────────────────────────────────────────────────   ║
║    用途: 显示提示、警告、错误和确认对话框                 ║
║                                                           ║
║    基本用法:                                              ║
║      gui.ShowError(title, message, details)               ║
║      gui.ShowWarning(title, message, details)             ║
║      gui.ShowInfo(title, message, details)                ║
║      gui.ShowSuccess(title, message, details)             ║
║      gui.ShowConfirm(title, message, onConfirm)           ║
║                                                           ║
║    快捷键:                                                ║
║      Enter    - 确认当前按钮                              ║
║      Esc      - 取消（选择最后一个按钮）                  ║
║      ←/→      - 切换按钮                                  ║
║      1-9      - 快速选择按钮                              ║
║                                                           ║
║  📊 进度条 (ProgressBar)                                  ║
║  ─────────────────────────────────────────────────────   ║
║    用途: 显示长时间操作的进度                             ║
║                                                           ║
║    基本用法:                                              ║
║      pb := gui.ShowProgressBar(config)                    ║
║      pb.Update(current, message)                          ║
║      pb.Close()                                           ║
║                                                           ║
║    两种模式:                                              ║
║      • 确定进度: 显示百分比和统计信息                     ║
║      • 不确定进度: 显示旋转动画                           ║
║                                                           ║
║  💬 AI 对话 (AI Chat)                                     ║
║  ─────────────────────────────────────────────────────   ║
║    用途: 与 AI 助手对话，获取 Git 帮助                    ║
║                                                           ║
║    基本用法:                                              ║
║      gui.ShowAIChat()                                     ║
║                                                           ║
║    核心功能:                                              ║
║      • 多轮对话，保持上下文                               ║
║      • 自动包含仓库状态                                   ║
║      • 预设问题快速入口 (Ctrl+P)                          ║
║      • 输入历史导航 (↑/↓)                                 ║
║                                                           ║
║    快捷键:                                                ║
║      Enter      - 发送消息                                ║
║      Ctrl+P     - 预设问题                                ║
║      Ctrl+L     - 清空历史                                ║
║      Ctrl+C     - 复制回复                                ║
║      ?          - 显示帮助                                ║
║                                                           ║
║  💡 使用技巧                                              ║
║  ─────────────────────────────────────────────────────   ║
║    1. 消息框支持自定义按钮和回调                          ║
║    2. 进度条可以动态切换确定/不确定模式                   ║
║    3. AI 对话会自动包含当前仓库上下文                     ║
║    4. 所有功能都支持键盘操作，无需鼠标                    ║
║    5. 可以通过配置文件自定义样式和行为                    ║
║                                                           ║
║  📚 更多信息                                              ║
║  ─────────────────────────────────────────────────────   ║
║    • MESSAGEBOX_USAGE.md  - 消息框详细文档                ║
║    • PROGRESSBAR_USAGE.md - 进度条详细文档                ║
║    • AI_CHAT_USAGE_V2.md  - AI 对话详细文档               ║
║    • UI_FEATURES_SUMMARY.md - 功能总结                    ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝

按任意键关闭...
`

	gui.ShowInfo("快速入门指南", guide)
	return nil
}

// showFeaturesHelp 显示功能帮助
func (gui *Gui) showFeaturesHelp() error {
	help := `
╔═══════════════════════════════════════════════════════════╗
║                    🎨 UI 功能帮助                         ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  可用功能                                                 ║
║  ─────────────────────────────────────────────────────   ║
║    1. 消息框 (MessageBox)                                 ║
║       • 5 种消息类型                                      ║
║       • 自定义按钮                                        ║
║       • 键盘导航                                          ║
║       • 自动关闭                                          ║
║                                                           ║
║    2. 进度条 (ProgressBar)                                ║
║       • 确定/不确定进度                                   ║
║       • 5 种样式 + 5 种动画                               ║
║       • 实时统计                                          ║
║       • 配置支持                                          ║
║                                                           ║
║    3. AI 对话 (AI Chat)                                   ║
║       • 多轮对话                                          ║
║       • 智能上下文                                        ║
║       • 预设问题                                          ║
║       • 丰富快捷键                                        ║
║                                                           ║
║  测试功能                                                 ║
║  ─────────────────────────────────────────────────────   ║
║    • 各功能独立测试菜单                                   ║
║    • 综合演示                                             ║
║    • 实际场景演示                                         ║
║    • 快速入门指南                                         ║
║                                                           ║
║  文档位置                                                 ║
║  ─────────────────────────────────────────────────────   ║
║    项目根目录下的 Markdown 文件:                          ║
║    • MESSAGEBOX_DESIGN.md                                 ║
║    • MESSAGEBOX_USAGE.md                                  ║
║    • PROGRESSBAR_DESIGN.md                                ║
║    • PROGRESSBAR_USAGE.md                                 ║
║    • AI_CHAT_DESIGN_V2.md                                 ║
║    • AI_CHAT_USAGE_V2.md                                  ║
║    • UI_FEATURES_SUMMARY.md                               ║
║                                                           ║
║  获取帮助                                                 ║
║  ─────────────────────────────────────────────────────   ║
║    • 在任何功能中按 ? 查看帮助                            ║
║    • 查看相应的文档文件                                   ║
║    • 运行测试菜单体验功能                                 ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝

按任意键关闭...
`

	gui.ShowInfo("功能帮助", help)
	return nil
}
