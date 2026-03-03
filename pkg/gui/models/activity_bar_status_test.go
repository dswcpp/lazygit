package models

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivityBarStatus_SetOperationInProgress(t *testing.T) {
	status := NewActivityBarStatus()

	// 测试设置操作为进行中
	status.SetOperationInProgress("pull", true)
	assert.True(t, status.IsOperationInProgress("pull"))

	// 测试设置操作为完成
	status.SetOperationInProgress("pull", false)
	assert.False(t, status.IsOperationInProgress("pull"))

	// 测试未设置的操作
	assert.False(t, status.IsOperationInProgress("push"))
}

func TestActivityBarStatus_MultipleOperations(t *testing.T) {
	status := NewActivityBarStatus()

	// 同时设置多个操作
	status.SetOperationInProgress("pull", true)
	status.SetOperationInProgress("push", true)
	status.SetOperationInProgress("fetch", true)

	assert.True(t, status.IsOperationInProgress("pull"))
	assert.True(t, status.IsOperationInProgress("push"))
	assert.True(t, status.IsOperationInProgress("fetch"))

	// 完成其中一个操作
	status.SetOperationInProgress("push", false)
	assert.True(t, status.IsOperationInProgress("pull"))
	assert.False(t, status.IsOperationInProgress("push"))
	assert.True(t, status.IsOperationInProgress("fetch"))
}

func TestActivityBarStatus_SpinnerAnimation(t *testing.T) {
	status := NewActivityBarStatus()

	// 初始帧应该是 0
	assert.Equal(t, 0, status.GetSpinnerFrame())

	// 测试 spinner 字符
	firstChar := status.GetSpinnerChar()
	assert.NotEmpty(t, firstChar)

	// 推进 spinner
	status.AdvanceSpinner()
	assert.Equal(t, 1, status.GetSpinnerFrame())
	secondChar := status.GetSpinnerChar()
	assert.NotEqual(t, firstChar, secondChar)

	// 推进到最后一帧后应该循环回 0
	for i := 0; i < 7; i++ {
		status.AdvanceSpinner()
	}
	assert.Equal(t, 0, status.GetSpinnerFrame())
}

func TestActivityBarStatus_SpinnerCharacters(t *testing.T) {
	status := NewActivityBarStatus()
	expectedChars := []string{
		"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧",
	}

	for i, expectedChar := range expectedChars {
		// 设置帧
		for j := 0; j < i; j++ {
			status.AdvanceSpinner()
		}

		char := status.GetSpinnerChar()
		assert.Equal(t, expectedChar, char)

		// 重置以便下次循环
		status.Reset()
	}
}

func TestActivityBarStatus_Reset(t *testing.T) {
	status := NewActivityBarStatus()

	// 设置一些状态
	status.SetOperationInProgress("pull", true)
	status.SetOperationInProgress("push", true)
	status.AdvanceSpinner()
	status.AdvanceSpinner()

	// 验证状态已设置
	assert.True(t, status.IsOperationInProgress("pull"))
	assert.True(t, status.IsOperationInProgress("push"))
	assert.Equal(t, 2, status.GetSpinnerFrame())

	// 重置
	status.Reset()

	// 验证所有状态已清除
	assert.False(t, status.IsOperationInProgress("pull"))
	assert.False(t, status.IsOperationInProgress("push"))
	assert.Equal(t, 0, status.GetSpinnerFrame())
}

func TestActivityBarStatus_ConcurrentAccess(t *testing.T) {
	status := NewActivityBarStatus()
	var wg sync.WaitGroup

	// 并发写入
	operations := []string{"pull", "push", "fetch", "merge", "rebase"}
	for _, op := range operations {
		wg.Add(1)
		go func(operation string) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				status.SetOperationInProgress(operation, true)
				status.SetOperationInProgress(operation, false)
			}
		}(op)
	}

	// 并发读取
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				for _, op := range operations {
					_ = status.IsOperationInProgress(op)
				}
			}
		}()
	}

	// 并发推进 spinner
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			status.AdvanceSpinner()
			_ = status.GetSpinnerChar()
		}
	}()

	// 等待所有 goroutine 完成（不应该有 race condition 或死锁）
	wg.Wait()

	// 测试通过就说明没有并发问题
}

func TestActivityBarStatus_EdgeCases(t *testing.T) {
	status := NewActivityBarStatus()

	// 测试空字符串操作
	status.SetOperationInProgress("", true)
	assert.True(t, status.IsOperationInProgress(""))

	// 测试特殊字符操作名
	specialOp := "pull-push-fetch"
	status.SetOperationInProgress(specialOp, true)
	assert.True(t, status.IsOperationInProgress(specialOp))

	// 多次设置同一个操作为 true
	status.SetOperationInProgress("pull", true)
	status.SetOperationInProgress("pull", true)
	status.SetOperationInProgress("pull", true)
	assert.True(t, status.IsOperationInProgress("pull"))

	// 多次设置同一个操作为 false
	status.SetOperationInProgress("pull", false)
	status.SetOperationInProgress("pull", false)
	assert.False(t, status.IsOperationInProgress("pull"))
}
