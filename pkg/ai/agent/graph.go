package agent

import (
	"context"
	"fmt"
)

// NodeID identifies a node in the agent graph.
type NodeID string

const (
	// NodePlan calls the LLM and routes based on the response:
	// → NodeCallTools if the LLM emitted tool calls,
	// → NodeWaitHuman if a valid plan block was produced,
	// → NodePlan again on empty/retry responses.
	NodePlan NodeID = "plan"

	// NodeCallTools executes the read-only tools queued by NodePlan,
	// then routes back to NodePlan.
	NodeCallTools NodeID = "call_tools"

	// NodeWaitHuman records the resume checkpoint and terminates the current
	// graph run, suspending execution until the user responds via Send().
	NodeWaitHuman NodeID = "wait_human"

	// NodeHandleConfirmation is the resume entry point after NodeWaitHuman.
	// It reads GraphState.HumanInput and routes to:
	//   → NodeExecuteStep  on confirm (Y / yes / 确认 / …)
	//   → NodeEnd          on deny    (N / no  / 取消 / …)
	//   → NodePlan         on any other text (user provided feedback → replan)
	NodeHandleConfirmation NodeID = "handle_confirmation"

	// NodeExecuteStep executes a single plan step and loops back to itself
	// until all steps are done, then routes to NodeDone.
	NodeExecuteStep NodeID = "execute_step"

	// NodeDone records execution summary and terminates the graph run.
	NodeDone NodeID = "done"

	// NodeEnd is the terminal sentinel — reaching it stops the Run loop.
	NodeEnd NodeID = "__end__"
)

// NodeFunc is a single processing unit in the graph.
//
// LangGraph analogy: a node receives the current GraphState, performs work,
// and returns (next node, updated state). The graph runtime applies the
// returned state before advancing — state always flows *through* nodes,
// never around them via shared mutable references.
type NodeFunc func(ctx context.Context, state GraphState, onUpdate func()) (NodeID, GraphState, error)

// Graph is a directed graph of NodeFuncs with optional breakpoints.
//
// LangGraph analogy: this is the compiled StateGraph that drives the agent.
// Nodes replace the ad-hoc for-loops in planLoop/execute; edges (the NodeID
// return values) replace the scattered switch/if branching.
type Graph struct {
	nodes       map[NodeID]NodeFunc
	breakpoints map[NodeID]bool
}

// NewGraph creates an empty Graph ready for node registration.
func NewGraph() *Graph {
	return &Graph{
		nodes:       make(map[NodeID]NodeFunc),
		breakpoints: make(map[NodeID]bool),
	}
}

// AddNode registers fn under id. Panics on duplicate registration.
func (g *Graph) AddNode(id NodeID, fn NodeFunc) {
	if _, exists := g.nodes[id]; exists {
		panic(fmt.Sprintf("agent graph: duplicate node %q", id))
	}
	g.nodes[id] = fn
}

// SetBreakpoint marks id as an interrupt point. When Run reaches this node it
// returns nil without executing the node's function, suspending the graph.
// The caller may resume by calling Run again with a later start node.
func (g *Graph) SetBreakpoint(id NodeID) {
	g.breakpoints[id] = true
}

// Run traverses the graph from startNode, threading state through each node,
// until NodeEnd is reached, a breakpoint is hit, the context is cancelled,
// or a node returns an error. Returns the final state after all node updates.
func (g *Graph) Run(ctx context.Context, startNode NodeID, state GraphState, onUpdate func()) (GraphState, error) {
	current := startNode
	for current != NodeEnd {
		if ctx.Err() != nil {
			return state, ctx.Err()
		}
		if g.breakpoints[current] {
			// Suspended — caller handles the interrupt (e.g. waiting for user input).
			return state, nil
		}
		fn, ok := g.nodes[current]
		if !ok {
			return state, fmt.Errorf("agent graph: unknown node %q", current)
		}
		var next NodeID
		var err error
		next, state, err = fn(ctx, state, onUpdate)
		if err != nil {
			return state, err
		}
		current = next
	}
	return state, nil
}
