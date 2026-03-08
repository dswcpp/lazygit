package agent

import (
	"context"
	"fmt"

	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

// ExampleCodeReviewWithCheckpoint 演示如何使用 CodeReviewAgent 的检查点功能
func ExampleCodeReviewWithCheckpoint(prov provider.Provider, tr *aii18n.Translator) {
	// 1. 创建 agent
	agent := NewCodeReviewAgent(prov, tr)

	// 2. 设置检查点器（支持会话恢复）
	checkpointer := NewMemoryCodeReviewCheckpointer()
	threadID := "review-session-123"
	agent.SetCheckpointer(checkpointer, threadID)

	// 3. 执行初始评审
	ctx := context.Background()
	err := agent.ReviewWithCallback(ctx, "main.go", "diff content", "", func(chunk string) {
		fmt.Print(chunk)
	})
	if err != nil {
		fmt.Printf("Review error: %v\n", err)
		return
	}

	// 4. 检查是否可以追问
	if agent.CanAsk() {
		// 5. 追问问题
		err = agent.Ask(ctx, "Can you explain the security concern in detail?", func(chunk string) {
			fmt.Print(chunk)
		})
		if err != nil {
			fmt.Printf("Ask error: %v\n", err)
		}
	}

	// 6. 完成后清除检查点
	checkpointer.Clear(threadID)
}

// ExampleCodeReviewBatch 演示如何使用 ConversationID 进行批量评审
func ExampleCodeReviewBatch(prov provider.Provider, tr *aii18n.Translator) {
	files := []struct {
		path string
		diff string
	}{
		{"file1.go", "diff1"},
		{"file2.go", "diff2"},
		{"file3.go", "diff3"},
	}

	// 使用相同的 ConversationID 进行批量评审
	conversationID := "batch-review-456"

	for _, file := range files {
		agent := NewCodeReviewAgent(prov, tr)

		// 设置 ConversationID 以保持上下文连续性
		agent.state = agent.state.WithConversationID(conversationID)

		ctx := context.Background()
		err := agent.ReviewWithCallback(ctx, file.path, file.diff, "", func(chunk string) {
			fmt.Printf("[%s] %s", file.path, chunk)
		})
		if err != nil {
			fmt.Printf("Review error for %s: %v\n", file.path, err)
		}
	}
}
