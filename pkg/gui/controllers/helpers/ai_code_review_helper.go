package helpers

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/ai/agent"
	aii18n "github.com/dswcpp/lazygit/pkg/ai/i18n"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

type AICodeReviewHelper struct {
	c             *HelperCommon
	loadingHelper *LoadingHelper
	aiHelper      *AIHelper
}

func NewAICodeReviewHelper(c *HelperCommon, loadingHelper *LoadingHelper, aiHelper *AIHelper) *AICodeReviewHelper {
	return &AICodeReviewHelper{c: c, loadingHelper: loadingHelper, aiHelper: aiHelper}
}

// ReviewDiff asks the user to confirm, then streams an AI code review for the
// given diff into the command log (Extras) panel.
//
// Flow:
//  1. Confirmation dialog: "Review file X?"
//  2. User confirms → centered loading overlay: "AI reviewing, please wait..."
//  3. First SSE chunk arrives → overlay closes; Extras panel header + content stream in.
//  4. Error before first chunk → overlay closes; error toast is shown.
func (self *AICodeReviewHelper) ReviewDiff(filePath string, diff string) error {
	if self.c.AIManager == nil {
		return self.aiHelper.ShowFirstTimeWizard()
	}

	if diff == "" {
		return errors.New(self.c.Tr.AICodeReviewNoDiff)
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AICodeReviewConfirmTitle,
		Prompt: fmt.Sprintf(self.c.Tr.AICodeReviewConfirmPrompt, filePath),
		HandleConfirm: func() error {
			return self.startReview(filePath, diff)
		},
	})
	return nil
}

// startReview shows the AI code review popup and launches the streaming review.
// Must be called from the UI thread (inside a Confirm HandleConfirm callback).
func (self *AICodeReviewHelper) startReview(filePath, diff string) error {
	// 创建CodeReviewAgent
	reviewAgent := agent.NewCodeReviewAgent(
		self.c.AIManager.Provider(),
		aii18n.NewTranslator(self.c.Tr),
	)

	// Prepare the popup view before pushing the context.
	aiView := self.c.Views().AICodeReview
	aiView.Clear()
	aiView.Autoscroll = true

	// Spinner frames for progress indicator
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerFrame := 0
	aiView.Title = fmt.Sprintf(" %s %s: %s ", spinner[0], self.c.Tr.AICodeReviewTitle, filePath)

	// Create cancellable context for the AI request
	ctx, cancel := context.WithCancel(context.Background())

	// Store cancel function in the context so it can be called by Esc key handler
	self.c.Contexts().AICodeReview.CancelFunc = cancel

	// Push the AI code review context to show the floating popup.
	self.c.Context().Push(self.c.Contexts().AICodeReview, types.OnFocusOpts{})

	// WithCenteredLoadingStatus runs the callback on a worker goroutine and
	// hides the overlay when the callback returns.
	self.loadingHelper.WithCenteredLoadingStatus(self.c.Tr.AICodeReviewStatus, func(_ gocui.Task) error {
		// firstChunk is closed exactly once: when the first response chunk
		// arrives (or when the request errors out). Closing it causes the
		// loading overlay to disappear.
		firstChunk := make(chan struct{})
		var once sync.Once
		signalFirst := func() { once.Do(func() { close(firstChunk) }) }

		// Start spinner animation in background
		spinnerDone := make(chan struct{})
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-spinnerDone:
					return
				case <-ticker.C:
					spinnerFrame = (spinnerFrame + 1) % len(spinner)
					self.c.OnUIThreadSync(func() error {
						aiView.Title = fmt.Sprintf(" %s %s: %s ", spinner[spinnerFrame], self.c.Tr.AICodeReviewTitle, filePath)
						return nil
					})
				}
			}
		}()

		// Streaming goroutine: runs independently after the overlay closes.
		// All UI writes use OnUIThreadSync (gocui.UpdateAsync) so that events
		// are enqueued directly from this single goroutine in order, avoiding
		// the race condition caused by OnUIThread spawning a new goroutine per
		// chunk which can arrive at the UI event queue out of order.
		go func() {
			defer func() {
				// Stop spinner
				close(spinnerDone)
				// Clear cancel function when stream completes
				self.c.Contexts().AICodeReview.CancelFunc = nil
				// Update title to show completion
				self.c.OnUIThreadSync(func() error {
					aiView.Title = fmt.Sprintf(" %s: %s ", self.c.Tr.AICodeReviewTitle, filePath)
					return nil
				})
			}()

			// 使用CodeReviewAgent执行评审
			err := reviewAgent.ReviewWithCallback(ctx, filePath, diff, "", func(chunk string) {
				signalFirst()
				self.c.OnUIThreadSync(func() error {
					fmt.Fprint(self.c.Views().AICodeReview, chunk)
					return nil
				})
			})

			if err != nil {
				signalFirst()
				self.c.OnUIThread(func() error {
					// Check if the error is due to cancellation
					if errors.Is(err, context.Canceled) {
						self.c.Toast(self.c.Tr.AICodeReviewCancelled)
						return nil
					}
					// Use friendly error handling from AIHelper
					friendlyErr := self.aiHelper.HandleAIError(err)
					self.c.Toast(friendlyErr.Error())
					return nil
				})
			}
		}()

		// Block here until the first chunk arrives → overlay hides.
		<-firstChunk
		return nil
	})

	return nil
}
