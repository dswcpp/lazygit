package agent

import (
	"context"
	"fmt"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/repocontext"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

const defaultMaxSteps = 10

// Agent runs a ReAct (Reason + Act) loop: it calls the LLM, parses any tool
// calls from the response, executes them with permission checks, feeds results
// back into the conversation, and repeats until the LLM stops calling tools or
// maxSteps is reached.
type Agent struct {
	provider  provider.Provider
	registry  *tools.Registry
	session   *Session
	confirmFn ConfirmFunc
	maxSteps  int
}

// NewAgent creates an Agent. If confirmFn is nil, all write operations are denied.
func NewAgent(p provider.Provider, r *tools.Registry, session *Session, confirmFn ConfirmFunc) *Agent {
	if confirmFn == nil {
		confirmFn = AutoDenyWrite()
	}
	return &Agent{
		provider:  p,
		registry:  r,
		session:   session,
		confirmFn: confirmFn,
		maxSteps:  defaultMaxSteps,
	}
}

// Session returns the session so the GUI layer can read UIMessages for rendering.
func (a *Agent) Session() *Session { return a.session }

// Run processes a user message through the ReAct loop.
// Must be called from a goroutine (not the UI thread).
//
// onUpdate is called after each session state change so the GUI can re-render.
// It is called on the same goroutine as Run — the GUI must use an appropriate
// mechanism (e.g. gocui.Update) to safely update the view.
func (a *Agent) Run(ctx context.Context, userMsg string, repoCtx repocontext.RepoContext, onUpdate func()) error {
	a.session.AddUserMessage(userMsg)
	if onUpdate != nil {
		onUpdate()
	}

	for step := 0; step < a.maxSteps; step++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// ── Step 1: LLM call ────────────────────────────────────────────────
		result, err := a.provider.Complete(ctx, a.session.ProviderMessages())
		if err != nil {
			return err
		}

		rawContent := result.Content
		toolCalls := tools.ParseToolCalls(rawContent)
		displayText := tools.StripToolBlocks(rawContent)

		a.session.AddAssistantMessage(displayText)
		if onUpdate != nil {
			onUpdate()
		}

		// ── Step 2: If no tool calls, we're done ────────────────────────────
		if len(toolCalls) == 0 {
			return nil
		}

		// ── Step 3: Execute each tool call ──────────────────────────────────
		for _, call := range toolCalls {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			tool, ok := a.registry.Get(call.Name)
			if !ok {
				a.session.AddToolResult(tools.ToolResult{
					CallID: call.ID,
					Output: fmt.Sprintf("未知工具: %s", call.Name),
				}, call.Name)
				if onUpdate != nil {
					onUpdate()
				}
				continue
			}

			schema := tool.Schema()
			a.session.AddToolCall(call)
			if onUpdate != nil {
				onUpdate()
			}

			// Permission check
			if schema.Permission.RequiresConfirm() {
				preview := buildConfirmPreview(call, schema)
				approved, err := a.confirmFn(call.Name, schema.Permission, preview)
				if err != nil {
					return err
				}
				if !approved {
					a.session.AddToolResult(tools.ToolResult{
						CallID: call.ID,
						Output: fmt.Sprintf("用户拒绝执行: %s", call.Name),
					}, call.Name)
					if onUpdate != nil {
						onUpdate()
					}
					// Tell the LLM the user declined so it can adapt
					a.session.providerMessages = append(a.session.providerMessages, provider.Message{
						Role:    provider.RoleUser,
						Content: fmt.Sprintf("[用户拒绝] 工具 %s 未被执行，请据此调整后续操作。", call.Name),
					})
					continue
				}
			}

			// Execute
			toolResult := tool.Execute(ctx, call)
			a.session.AddToolResult(toolResult, call.Name)
			if onUpdate != nil {
				onUpdate()
			}
		}

		// ── Step 4: Inject refresh prompt and loop ──────────────────────────
		// The next LLM call will see all tool results and decide what to do next.
	}

	// Max steps reached — add a note and stop
	a.session.AddSystemNote(fmt.Sprintf("已达到最大步数 (%d)，停止执行。", a.maxSteps))
	if onUpdate != nil {
		onUpdate()
	}
	return nil
}

// buildConfirmPreview builds a human-readable preview of what a tool will do.
func buildConfirmPreview(call tools.ToolCall, schema tools.ToolSchema) string {
	if len(call.Params) == 0 {
		return fmt.Sprintf("工具: %s\n描述: %s\n权限: %s",
			call.Name, schema.Description, schema.Permission)
	}
	params := ""
	for k, v := range call.Params {
		params += fmt.Sprintf("  %s: %v\n", k, v)
	}
	return fmt.Sprintf("工具: %s\n描述: %s\n权限: %s\n参数:\n%s",
		call.Name, schema.Description, schema.Permission, params)
}
