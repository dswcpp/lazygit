package gui

import (
	"time"
)

// 这个文件包含消息框的使用示例

// TestMessageBoxError 测试错误消息框
func (gui *Gui) TestMessageBoxError() error {
	gui.ShowError(
		"操作失败",
		"无法连接到远程仓库，请检查网络连接。",
		"错误代码: ECONNREFUSED\n主机: github.com\n端口: 443",
	)
	return nil
}

// TestMessageBoxWarning 测试警告消息框
func (gui *Gui) TestMessageBoxWarning() error {
	gui.ShowWarning(
		"警告",
		"你即将强制推送到远程分支 'main'，这将覆盖远程仓库的历史记录！",
		"此操作不可撤销，请谨慎操作。",
	)
	return nil
}

// TestMessageBoxInfo 测试信息消息框
func (gui *Gui) TestMessageBoxInfo() error {
	gui.ShowInfo(
		"提示",
		"当前分支已经是最新的，无需拉取更新。",
		"最后更新时间: 2024-01-01 12:00:00",
	)
	return nil
}

// TestMessageBoxSuccess 测试成功消息框
func (gui *Gui) TestMessageBoxSuccess() error {
	gui.ShowSuccess(
		"操作成功",
		"分支 'feature-x' 已成功删除。",
		"本地分支和远程分支都已删除。",
	)
	return nil
}

// TestMessageBoxConfirm 测试确认对话框
func (gui *Gui) TestMessageBoxConfirm() error {
	gui.ShowConfirm(
		"确认删除",
		"确定要删除分支 'feature-x' 吗？此操作不可撤销。",
		func() {
			gui.c.Toast("已确认删除")
		},
	)
	return nil
}

// TestMessageBoxYesNoCancel 测试是/否/取消对话框
func (gui *Gui) TestMessageBoxYesNoCancel() error {
	gui.ShowYesNoCancel(
		"保存更改",
		"检测到未保存的更改，是否保存？",
		func() {
			gui.c.Toast("已保存")
		},
		func() {
			gui.c.Toast("已放弃更改")
		},
	)
	return nil
}

// TestMessageBoxCustomButtons 测试自定义按钮
func (gui *Gui) TestMessageBoxCustomButtons() error {
	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeQuestion,
		Title:   "选择操作",
		Message: "检测到未提交的更改，如何处理？",
		Buttons: []string{"暂存", "丢弃", "取消"},
	}, func(buttonIndex int) {
		switch buttonIndex {
		case 0:
			gui.c.Toast("已暂存更改")
		case 1:
			gui.c.Toast("已丢弃更改")
		case 2:
			gui.c.Toast("已取消操作")
		}
	})
	return nil
}

// TestMessageBoxAutoClose 测试自动关闭消息框
func (gui *Gui) TestMessageBoxAutoClose() error {
	gui.ShowAutoCloseMessage(
		MessageTypeSuccess,
		"操作成功",
		"文件已保存，此消息将在 3 秒后自动关闭。",
		3*time.Second,
	)
	return nil
}

// TestMessageBoxLongText 测试长文本消息框
func (gui *Gui) TestMessageBoxLongText() error {
	gui.ShowError(
		"Git 操作失败",
		"无法推送到远程仓库。远程仓库包含你本地没有的提交，这通常是因为另一个仓库已向该引用进行了推送。你可能需要先整合远程变更（例如 'git pull'）再推送。",
		"详细错误信息:\n"+
			"To https://github.com/user/repo.git\n"+
			" ! [rejected]        main -> main (fetch first)\n"+
			"error: failed to push some refs to 'https://github.com/user/repo.git'\n"+
			"hint: Updates were rejected because the remote contains work that you do\n"+
			"hint: not have locally. This is usually caused by another repository pushing\n"+
			"hint: to the same ref. You may want to first integrate the remote changes\n"+
			"hint: (e.g., 'git pull ...') before pushing again.",
	)
	return nil
}

// TestMessageBoxAllTypes 测试所有消息类型
func (gui *Gui) TestMessageBoxAllTypes() error {
	types := []struct {
		msgType MessageType
		title   string
		message string
	}{
		{MessageTypeInfo, "信息", "这是一条信息消息"},
		{MessageTypeSuccess, "成功", "操作已成功完成"},
		{MessageTypeWarning, "警告", "请注意这个警告信息"},
		{MessageTypeError, "错误", "发生了一个错误"},
		{MessageTypeQuestion, "问题", "你确定要继续吗？"},
	}

	for i, t := range types {
		time.Sleep(time.Duration(i) * 2 * time.Second)
		gui.ShowMessageBox(MessageBoxConfig{
			Type:    t.msgType,
			Title:   t.title,
			Message: t.message,
			Buttons: []string{"确定"},
		}, nil)
		time.Sleep(2 * time.Second)
	}

	return nil
}

// 实际使用示例

// handleGitPushError 处理 Git Push 错误
func (gui *Gui) handleGitPushError(err error) {
	gui.ShowError(
		"推送失败",
		"无法推送到远程仓库。",
		err.Error(),
	)
}

// confirmForcePush 确认强制推送
func (gui *Gui) confirmForcePush(branch string) {
	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeWarning,
		Title:   "确认强制推送",
		Message: "你即将强制推送到远程分支 '" + branch + "'，这将覆盖远程仓库的历史记录！此操作不可撤销。",
		Buttons: []string{"确认", "取消"},
	}, func(buttonIndex int) {
		if buttonIndex == 0 {
			// 执行强制推送
			gui.c.Toast("正在强制推送...")
		}
	})
}

// confirmDeleteBranch 确认删除分支
func (gui *Gui) confirmDeleteBranch(branch string, hasRemote bool) {
	message := "确定要删除分支 '" + branch + "' 吗？"
	if hasRemote {
		message = "确定要删除本地和远程分支 '" + branch + "' 吗？"
	}

	gui.ShowConfirm(
		"确认删除",
		message,
		func() {
			// 执行删除操作
			gui.c.Toast("分支已删除")
		},
	)
}

// showMergeConflict 显示合并冲突
func (gui *Gui) showMergeConflict(conflictFiles []string) {
	details := "冲突文件:\n"
	for _, file := range conflictFiles {
		details += "  - " + file + "\n"
	}

	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeWarning,
		Title:   "合并冲突",
		Message: "合并过程中发现冲突，请解决冲突后再提交。",
		Details: details,
		Buttons: []string{"解决冲突", "中止合并"},
	}, func(buttonIndex int) {
		if buttonIndex == 0 {
			// 打开冲突解决界面
			gui.c.Toast("打开冲突解决界面")
		} else {
			// 中止合并
			gui.c.Toast("已中止合并")
		}
	})
}

// showStashOptions 显示暂存选项
func (gui *Gui) showStashOptions() {
	gui.ShowMessageBox(MessageBoxConfig{
		Type:    MessageTypeQuestion,
		Title:   "未提交的更改",
		Message: "检测到未提交的更改，如何处理？",
		Buttons: []string{"暂存", "丢弃", "取消"},
	}, func(buttonIndex int) {
		switch buttonIndex {
		case 0:
			gui.c.Toast("已暂存更改")
		case 1:
			gui.ShowConfirm(
				"确认丢弃",
				"确定要丢弃所有未提交的更改吗？此操作不可撤销。",
				func() {
					gui.c.Toast("已丢弃更改")
				},
			)
		}
	})
}

// showOperationSuccess 显示操作成功
func (gui *Gui) showOperationSuccess(operation, details string) {
	gui.ShowAutoCloseMessage(
		MessageTypeSuccess,
		operation+"成功",
		details,
		2*time.Second,
	)
}
