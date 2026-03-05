package agent

import "github.com/dswcpp/lazygit/pkg/ai/tools"

// ConfirmFunc is called by the agent before executing a tool that requires
// user confirmation (permission level ≥ PermWriteLocal).
//
// The implementation provided by the GUI layer must block until the user
// responds, typically by using a channel to bridge the goroutine and the
// gocui UI thread:
//
//	func makeConfirmFn(c *HelperCommon) agent.ConfirmFunc {
//	    return func(toolName string, perm tools.PermissionLevel, preview string) (bool, error) {
//	        ch := make(chan bool, 1)
//	        c.OnUIThread(func() error {
//	            c.Confirm(types.ConfirmOpts{
//	                Title:  "AI 请求执行: " + toolName,
//	                Prompt: preview,
//	                HandleConfirm: func() error { ch <- true; return nil },
//	                HandleClose:   func() error { ch <- false; return nil },
//	            })
//	            return nil
//	        })
//	        return <-ch, nil
//	    }
//	}
type ConfirmFunc func(toolName string, perm tools.PermissionLevel, preview string) (approved bool, err error)

// AutoApproveAll returns a ConfirmFunc that approves every tool call without
// prompting. Intended for automated testing or trusted non-interactive runs.
func AutoApproveAll() ConfirmFunc {
	return func(_ string, _ tools.PermissionLevel, _ string) (bool, error) {
		return true, nil
	}
}

// AutoDenyWrite returns a ConfirmFunc that approves read-only tools but
// denies all write operations.
func AutoDenyWrite() ConfirmFunc {
	return func(_ string, perm tools.PermissionLevel, _ string) (bool, error) {
		return perm == tools.PermReadOnly, nil
	}
}
