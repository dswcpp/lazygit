package helpers

import (
	"errors"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/dswcpp/lazygit/pkg/gui/presentation"
)

// LoadingHelper shows a centered, framed loading overlay with a spinning
// animation while a long-running task executes. It is designed to be
// reusable across different parts of the codebase.
type LoadingHelper struct {
	c *HelperCommon
}

func NewLoadingHelper(c *HelperCommon) *LoadingHelper {
	return &LoadingHelper{c: c}
}

// WithCenteredLoadingStatus displays a centered loading popup with the given
// message while f executes on a worker goroutine. The overlay is hidden
// automatically when f returns.
func (self *LoadingHelper) WithCenteredLoadingStatus(message string, f func(gocui.Task) error) {
	self.c.OnWorker(func(task gocui.Task) error {
		self.c.OnUIThread(func() error {
			return self.showLoadingView(message)
		})

		stop := make(chan struct{})
		go self.animateSpinner(message, stop)

		err := f(task)

		close(stop)
		self.c.OnUIThread(func() error {
			self.c.Views().Loading.Visible = false
			return nil
		})

		return err
	})
}

// showLoadingView positions and reveals the loading view centered on screen.
// Must be called on the UI thread.
func (self *LoadingHelper) showLoadingView(message string) error {
	view := self.c.Views().Loading
	width, height := self.c.GocuiGui().Size()

	// Content: " <spinner> <message>"  →  2 + len(message) chars
	contentWidth := len([]rune(message)) + 4
	if contentWidth < 16 {
		contentWidth = 16
	}

	// In gocui: content_width = x1 - x0 - 1, so x1 = x0 + contentWidth + 1
	totalWidth := contentWidth + 2
	x0 := width/2 - totalWidth/2
	x1 := x0 + contentWidth + 1
	y0 := height/2 - 2
	y1 := y0 + 2

	_, err := self.c.GocuiGui().SetView(view.Name(), x0, y0, x1, y1, 0)
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		return err
	}

	view.Visible = true
	self.renderContent(view, message)
	return nil
}

// animateSpinner runs in a goroutine, updating the spinner frame on the UI
// thread at the configured spinner rate until stop is closed.
func (self *LoadingHelper) animateSpinner(message string, stop chan struct{}) {
	rate := self.c.UserConfig().Gui.Spinner.Rate
	ticker := time.NewTicker(time.Millisecond * time.Duration(rate))
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			self.c.OnUIThread(func() error {
				view := self.c.Views().Loading
				if !view.Visible {
					return nil
				}
				self.renderContent(view, message)
				return nil
			})
		}
	}
}

func (self *LoadingHelper) renderContent(view *gocui.View, message string) {
	spinner := presentation.Loader(time.Now(), self.c.UserConfig().Gui.Spinner)
	self.c.SetViewContent(view, " "+spinner+" "+message)
}
