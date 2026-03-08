package gittools

import (
	"context"
	"fmt"

	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

// StageAllTool stages all working tree changes.
type StageAllTool struct{ d *Deps }

func NewStageAllTool(d *Deps) tools.Tool { return &StageAllTool{d} }

func (t *StageAllTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "stage_all",
		Description: t.d.Tr.ToolStageAllDesc(),
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermWriteLocal,
	}
}

func (t *StageAllTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	if err := t.d.WorkingTree.StageAll(false); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolStageAllFailed(err)}
	}
	t.d.Refresh(ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolStageAllSuccess()}
}

// StageFileTool stages a single file.
type StageFileTool struct{ d *Deps }

func NewStageFileTool(d *Deps) tools.Tool { return &StageFileTool{d} }

func (t *StageFileTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "stage_file",
		Description: t.d.Tr.ToolStageFileDesc(),
		Params: map[string]tools.ParamSchema{
			"path": {Type: "string", Description: t.d.Tr.ToolFilePath(), Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *StageFileTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	path := strParam(call.Params, "path", "")
	if path == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingPathParam()}
	}
	if err := t.d.WorkingTree.StageFile(path); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolStageFileFailed(err)}
	}
	t.d.Refresh(ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolStageFileSuccess(path)}
}

// UnstageAllTool unstages all staged changes (reset mixed HEAD).
type UnstageAllTool struct{ d *Deps }

func NewUnstageAllTool(d *Deps) tools.Tool { return &UnstageAllTool{d} }

func (t *UnstageAllTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "unstage_all",
		Description: t.d.Tr.ToolUnstageAllDesc(),
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermWriteLocal,
	}
}

func (t *UnstageAllTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	if err := t.d.WorkingTree.ResetMixed("HEAD"); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolUnstageAllFailed(err)}
	}
	t.d.Refresh(ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolUnstageAllSuccess()}
}

// UnstageFileTool unstages a single file.
type UnstageFileTool struct{ d *Deps }

func NewUnstageFileTool(d *Deps) tools.Tool { return &UnstageFileTool{d} }

func (t *UnstageFileTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "unstage_file",
		Description: t.d.Tr.ToolUnstageFileDesc(),
		Params: map[string]tools.ParamSchema{
			"path": {Type: "string", Description: t.d.Tr.ToolFilePath(), Required: true},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *UnstageFileTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	path := strParam(call.Params, "path", "")
	if path == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingPathParam()}
	}
	tracked := true
	for _, f := range t.d.GetFiles() {
		if f.Path == path {
			tracked = f.Tracked
			break
		}
	}
	if err := t.d.WorkingTree.UnStageFile([]string{path}, tracked); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolUnstageFileFailed(err)}
	}
	t.d.Refresh(ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: t.d.Tr.ToolUnstageFileSuccess(path)}
}

// DiscardFileTool discards all changes in a file (restores to HEAD).
type DiscardFileTool struct{ d *Deps }

func NewDiscardFileTool(d *Deps) tools.Tool { return &DiscardFileTool{d} }

func (t *DiscardFileTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "discard_file",
		Description: t.d.Tr.ToolDiscardFileDesc(),
		Params: map[string]tools.ParamSchema{
			"path": {Type: "string", Description: t.d.Tr.ToolFilePath(), Required: true},
		},
		Permission: tools.PermDestructive,
	}
}

func (t *DiscardFileTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	path := strParam(call.Params, "path", "")
	if path == "" {
		return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolMissingPathParam()}
	}
	for _, f := range t.d.GetFiles() {
		if f.Path == path {
			if err := t.d.WorkingTree.DiscardAllFileChanges(f); err != nil {
				return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolDiscardChangesFailed(err)}
			}
			t.d.Refresh(ScopeFiles)
			return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("Discarded all changes in: %s", path)}
		}
	}
	return tools.ToolResult{CallID: call.ID, Output: t.d.Tr.ToolFileNotInWorkdir(path)}
}
