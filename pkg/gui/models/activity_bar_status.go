package models

import (
	"sync"
)

// ActivityBarStatus tracks ongoing operations in the activity bar
type ActivityBarStatus struct {
	mu                sync.RWMutex
	ongoingOperations map[string]bool // action name -> in progress
	spinnerFrame      int             // current spinner animation frame
}

// NewActivityBarStatus creates a new activity bar status tracker
func NewActivityBarStatus() *ActivityBarStatus {
	return &ActivityBarStatus{
		ongoingOperations: make(map[string]bool),
		spinnerFrame:      0,
	}
}

// SetOperationInProgress marks an operation as in progress or completed
func (self *ActivityBarStatus) SetOperationInProgress(action string, inProgress bool) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if inProgress {
		self.ongoingOperations[action] = true
	} else {
		delete(self.ongoingOperations, action)
	}
}

// IsOperationInProgress checks if an operation is currently in progress
func (self *ActivityBarStatus) IsOperationInProgress(action string) bool {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return self.ongoingOperations[action]
}

// GetSpinnerFrame returns the current spinner animation frame
func (self *ActivityBarStatus) GetSpinnerFrame() int {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return self.spinnerFrame
}

// AdvanceSpinner advances the spinner animation to the next frame
func (self *ActivityBarStatus) AdvanceSpinner() {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.spinnerFrame = (self.spinnerFrame + 1) % 8
}

// GetSpinnerChar returns the spinner character for the current frame
func (self *ActivityBarStatus) GetSpinnerChar() string {
	// Braille spinner frames
	spinnerChars := []string{
		"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧",
	}

	frame := self.GetSpinnerFrame()
	return spinnerChars[frame]
}

// Reset clears all ongoing operations
func (self *ActivityBarStatus) Reset() {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.ongoingOperations = make(map[string]bool)
	self.spinnerFrame = 0
}
