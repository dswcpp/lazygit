package gui

import (
	"time"
)

// 这个文件包含进度条的使用示例

// TestProgressBarDeterminate 测试确定进度条
func (gui *Gui) TestProgressBarDeterminate() error {
	// 创建进度条
	pb := gui.ShowProgressBar(ProgressBarConfig{
		Title:          gui.c.Tr.AIChatPushingToRemote,
		Total:          20 * 1024 * 1024, // 20 MB
		Width:          30,
		ShowPercentage: true,
		ShowStats:      true,
		Style:          ProgressBarStyleBlock,
		Indeterminate:  false,
	})

	// 模拟推送过程
	go func() {
		for i := int64(0); i <= pb.config.Total; i += 512 * 1024 {
			pb.Update(i, "")
			time.Sleep(200 * time.Millisecond)
		}
		pb.Close()
	}()

	return nil
}

// TestProgressBarIndeterminate 测试不确定进度条
func (gui *Gui) TestProgressBarIndeterminate() error {
	// 创建不确定进度条
	pb := gui.ShowProgressBar(ProgressBarConfig{
		Title:         "正在克隆仓库...",
		Message:       "正在连接服务器...",
		Indeterminate: true,
		SpinnerStyle:  SpinnerStyleBraille,
	})

	// 模拟克隆过程
	go func() {
		time.Sleep(2 * time.Second)
		pb.Update(0, "正在接收对象...")
		time.Sleep(2 * time.Second)
		pb.Update(0, "正在解析增量...")
		time.Sleep(2 * time.Second)
		pb.Close()
	}()

	return nil
}

// TestProgressBarStyles 测试不同样式的进度条
func (gui *Gui) TestProgressBarStyles() error {
	styles := []struct {
		name  string
		style ProgressBarStyle
	}{
		{"方块样式", ProgressBarStyleBlock},
		{"点状样式", ProgressBarStyleDot},
		{"箭头样式", ProgressBarStyleArrow},
		{"渐变样式", ProgressBarStyleGradient},
		{"ASCII样式", ProgressBarStyleASCII},
	}

	for _, s := range styles {
		pb := gui.ShowProgressBar(ProgressBarConfig{
			Title:          s.name,
			Total:          100,
			Width:          30,
			ShowPercentage: true,
			Style:          s.style,
		})

		// 模拟进度
		go func(pb *ProgressBar) {
			for i := int64(0); i <= 100; i += 10 {
				pb.Update(i, "")
				time.Sleep(300 * time.Millisecond)
			}
			pb.Close()
		}(pb)

		// 等待当前进度条完成
		time.Sleep(4 * time.Second)
	}

	return nil
}

// TestProgressBarSpinners 测试不同的旋转动画
func (gui *Gui) TestProgressBarSpinners() error {
	spinners := []struct {
		name    string
		spinner SpinnerStyle
	}{
		{"Braille 点阵", SpinnerStyleBraille},
		{"线条旋转", SpinnerStyleLine},
		{"箭头旋转", SpinnerStyleArrow},
		{"点旋转", SpinnerStyleDot},
		{"圆圈旋转", SpinnerStyleCircle},
	}

	for _, s := range spinners {
		pb := gui.ShowProgressBar(ProgressBarConfig{
			Title:         s.name,
			Message:       "正在处理...",
			Indeterminate: true,
			SpinnerStyle:  s.spinner,
		})

		// 显示3秒
		go func(pb *ProgressBar) {
			time.Sleep(3 * time.Second)
			pb.Close()
		}(pb)

		time.Sleep(4 * time.Second)
	}

	return nil
}

// 实际使用示例：在 Git Push 中使用进度条
func (gui *Gui) pushWithProgress() error {
	// 显示进度条
	pb := gui.ShowProgressBar(ProgressBarConfig{
		Title:          gui.c.Tr.AIChatPushingToRemote,
		Message:        "正在连接...",
		Indeterminate:  true,
		SpinnerStyle:   SpinnerStyleBraille,
	})

	// 执行 git push（这里需要实际的 git 命令集成）
	go func() {
		// 模拟推送过程
		time.Sleep(1 * time.Second)

		// 切换到确定进度
		pb.config.Indeterminate = false
		pb.config.Total = 10 * 1024 * 1024 // 10 MB
		pb.config.ShowStats = true
		pb.Update(0, "正在推送...")

		// 模拟推送进度
		for i := int64(0); i <= pb.config.Total; i += 256 * 1024 {
			pb.Update(i, "")
			time.Sleep(100 * time.Millisecond)
		}

		pb.Close()

		// 显示成功消息
		gui.c.Toast("推送成功！")
	}()

	return nil
}

// 实际使用示例：在 Git Clone 中使用进度条
func (gui *Gui) cloneWithProgress(url string) error {
	pb := gui.ShowProgressBar(ProgressBarConfig{
		Title:         "正在克隆仓库...",
		Message:       "正在连接服务器...",
		Indeterminate: true,
		SpinnerStyle:  SpinnerStyleBraille,
	})

	go func() {
		// 模拟克隆过程
		time.Sleep(1 * time.Second)
		pb.Update(0, "正在接收对象...")

		time.Sleep(2 * time.Second)
		pb.Update(0, "正在解析增量...")

		time.Sleep(2 * time.Second)
		pb.Update(0, "正在检出文件...")

		time.Sleep(1 * time.Second)
		pb.Close()

		gui.c.Toast("克隆成功！")
	}()

	return nil
}

// 实际使用示例：在 Git Fetch 中使用进度条
func (gui *Gui) fetchWithProgress() error {
	pb := gui.ShowProgressBar(ProgressBarConfig{
		Title:          "正在获取更新...",
		Total:          5 * 1024 * 1024, // 5 MB
		Width:          30,
		ShowPercentage: true,
		ShowStats:      true,
		Style:          ProgressBarStyleGradient,
	})

	go func() {
		// 模拟获取过程
		for i := int64(0); i <= pb.config.Total; i += 128 * 1024 {
			pb.Update(i, "")
			time.Sleep(50 * time.Millisecond)
		}
		pb.Close()

		gui.c.Toast("获取成功！")
	}()

	return nil
}
