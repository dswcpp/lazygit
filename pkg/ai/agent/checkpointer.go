package agent

import "sync"

// Checkpointer persists and restores GraphState to/from durable storage.
//
// LangGraph analogy: this is the persistence layer (MemorySaver / SqliteSaver)
// that enables human-in-the-loop resumption across process restarts.
// When the graph suspends at nodeWaitHuman the agent saves its state;
// on the next app start the caller can restore it so the user picks up
// exactly where they left off.
type Checkpointer interface {
	// Save writes the current state for a given conversation thread.
	Save(threadID string, state GraphState) error

	// Load returns the saved state for threadID, or (zero, false) if absent.
	Load(threadID string) (GraphState, bool)

	// Clear removes the checkpoint for threadID (call after execution completes
	// or the user cancels, so stale state doesn't surface on the next launch).
	Clear(threadID string)
}

// MemoryCheckpointer stores state in-process only.
// State is lost when the process exits; useful for tests and as a
// lightweight default when cross-restart persistence isn't needed.
type MemoryCheckpointer struct {
	mu     sync.RWMutex
	states map[string]GraphState
}

// NewMemoryCheckpointer creates a ready-to-use MemoryCheckpointer.
func NewMemoryCheckpointer() *MemoryCheckpointer {
	return &MemoryCheckpointer{states: make(map[string]GraphState)}
}

func (c *MemoryCheckpointer) Save(threadID string, state GraphState) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.states[threadID] = state
	return nil
}

func (c *MemoryCheckpointer) Load(threadID string) (GraphState, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.states[threadID]
	return s, ok
}

func (c *MemoryCheckpointer) Clear(threadID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.states, threadID)
}

// ────────────────────────────────────────────────────────────────────────────
// CodeReviewCheckpointer - Specialized checkpointer for CodeReviewAgent
// ────────────────────────────────────────────────────────────────────────────

// CodeReviewCheckpointer persists and restores CodeReviewState.
type CodeReviewCheckpointer interface {
	Save(threadID string, state CodeReviewState) error
	Load(threadID string) (CodeReviewState, bool)
	Clear(threadID string)
}

// MemoryCodeReviewCheckpointer stores CodeReviewState in-process only.
type MemoryCodeReviewCheckpointer struct {
	mu     sync.RWMutex
	states map[string]CodeReviewState
}

// NewMemoryCodeReviewCheckpointer creates a ready-to-use MemoryCodeReviewCheckpointer.
func NewMemoryCodeReviewCheckpointer() *MemoryCodeReviewCheckpointer {
	return &MemoryCodeReviewCheckpointer{states: make(map[string]CodeReviewState)}
}

func (c *MemoryCodeReviewCheckpointer) Save(threadID string, state CodeReviewState) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.states[threadID] = state
	return nil
}

func (c *MemoryCodeReviewCheckpointer) Load(threadID string) (CodeReviewState, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.states[threadID]
	return s, ok
}

func (c *MemoryCodeReviewCheckpointer) Clear(threadID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.states, threadID)
}
