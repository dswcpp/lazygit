package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCheckpointer_SaveLoadClear(t *testing.T) {
	c := NewMemoryCheckpointer()
	threadID := "test-thread-1"

	// 初始状态：Load 返回 false
	_, ok := c.Load(threadID)
	assert.False(t, ok, "should return false before any save")

	// 保存一个有 ResumeFrom 的状态
	state := GraphState{
		Phase:      PhaseWaitingConfirm,
		ResumeFrom: NodeHandleConfirmation,
		HumanInput: "",
	}
	err := c.Save(threadID, state)
	assert.NoError(t, err)

	// Load 应该返回之前保存的状态
	loaded, ok := c.Load(threadID)
	assert.True(t, ok)
	assert.Equal(t, PhaseWaitingConfirm, loaded.Phase)
	assert.Equal(t, NodeHandleConfirmation, loaded.ResumeFrom)

	// Clear 后 Load 返回 false
	c.Clear(threadID)
	_, ok = c.Load(threadID)
	assert.False(t, ok, "should return false after clear")
}

func TestMemoryCheckpointer_MultipleThreads(t *testing.T) {
	c := NewMemoryCheckpointer()

	s1 := GraphState{Phase: PhaseWaitingConfirm, ResumeFrom: NodeHandleConfirmation}
	s2 := GraphState{Phase: PhasePlanning, ResumeFrom: ""}

	assert.NoError(t, c.Save("thread-1", s1))
	assert.NoError(t, c.Save("thread-2", s2))

	l1, ok1 := c.Load("thread-1")
	l2, ok2 := c.Load("thread-2")

	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.Equal(t, NodeHandleConfirmation, l1.ResumeFrom)
	assert.Equal(t, NodeID(""), l2.ResumeFrom)

	// 清除 thread-1 不影响 thread-2
	c.Clear("thread-1")
	_, ok1 = c.Load("thread-1")
	_, ok2 = c.Load("thread-2")
	assert.False(t, ok1)
	assert.True(t, ok2)
}

func TestSetCheckpointer_RestoresInterruptedState(t *testing.T) {
	c := NewMemoryCheckpointer()
	threadID := "session-abc"

	// 模拟：上一次进程在 nodeWaitHuman 处保存了检查点
	interrupted := GraphState{
		Phase:           PhaseWaitingConfirm,
		ResumeFrom:      NodeHandleConfirmation,
		ToolCallHistory: map[string]int{"stage_all:map[]": 1},
	}
	assert.NoError(t, c.Save(threadID, interrupted))

	// 新进程启动，创建 agent 并挂载 checkpointer
	agent := &TwoPhaseAgent{
		state: GraphState{
			Phase:           PhasePlanning,
			ToolCallHistory: make(map[string]int),
		},
	}
	agent.SetCheckpointer(c, threadID)

	// agent 应该自动恢复中断状态
	assert.Equal(t, NodeHandleConfirmation, agent.state.ResumeFrom)
	assert.Equal(t, PhaseWaitingConfirm, agent.state.Phase)
	assert.Equal(t, 1, agent.state.ToolCallHistory["stage_all:map[]"])
}

func TestSetCheckpointer_NoCheckpointDoesNothing(t *testing.T) {
	c := NewMemoryCheckpointer()
	// 没有保存过任何状态

	agent := &TwoPhaseAgent{
		state: GraphState{
			Phase:           PhasePlanning,
			ToolCallHistory: make(map[string]int),
		},
	}
	agent.SetCheckpointer(c, "new-thread")

	// 没有检查点时，state 保持原样
	assert.Equal(t, PhasePlanning, agent.state.Phase)
	assert.Empty(t, agent.state.ResumeFrom)
}

func TestSetCheckpointer_IgnoresCompletedCheckpoint(t *testing.T) {
	c := NewMemoryCheckpointer()
	threadID := "done-thread"

	// 保存一个已完成（ResumeFrom 为空）的状态——不应恢复
	completed := GraphState{
		Phase:           PhaseDone,
		ResumeFrom:      "", // 没有挂起点
		ToolCallHistory: make(map[string]int),
	}
	assert.NoError(t, c.Save(threadID, completed))

	agent := &TwoPhaseAgent{
		state: GraphState{
			Phase:           PhasePlanning,
			ToolCallHistory: make(map[string]int),
		},
	}
	agent.SetCheckpointer(c, threadID)

	// ResumeFrom 为空的检查点不应被恢复
	assert.Equal(t, PhasePlanning, agent.state.Phase)
	assert.Empty(t, agent.state.ResumeFrom)
}
