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
		Description: "暂存所有工作区变更",
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermWriteLocal,
	}
}

func (t *StageAllTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	if err := t.d.WorkingTree.StageAll(false); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("暂存全部失败: %v", err)}
	}
	t.d.Refresh(ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: "已暂存所有变更"}
}

// StageFileTool stages a single file.
type StageFileTool struct{ d *Deps }

func NewStageFileTool(d *Deps) tools.Tool { return &StageFileTool{d} }

func (t *StageFileTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "stage_file",
		Description: "Stage specified file",
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
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("暂存文件失败: %v", err)}
	}
	t.d.Refresh(ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已暂存: %s", path)}
}

// UnstageAllTool unstages all staged changes (reset mixed HEAD).
type UnstageAllTool struct{ d *Deps }

func NewUnstageAllTool(d *Deps) tools.Tool { return &UnstageAllTool{d} }

func (t *UnstageAllTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "unstage_all",
		Description: "取消所有暂存（git reset HEAD）",
		Params:      map[string]tools.ParamSchema{},
		Permission:  tools.PermWriteLocal,
	}
}

func (t *UnstageAllTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	if err := t.d.WorkingTree.ResetMixed("HEAD"); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("取消所有暂存失败: %v", err)}
	}
	t.d.Refresh(ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: "已取消所有暂存"}
}

// UnstageFileTool unstages a single file.
type UnstageFileTool struct{ d *Deps }

func NewUnstageFileTool(d *Deps) tools.Tool { return &UnstageFileTool{d} }

func (t *UnstageFileTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "unstage_file",
		Description: "Unstage specified file",
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
	// Determine whether the file is tracked (affects unstage command)
	tracked := true
	for _, f := range t.d.GetFiles() {
		if f.Path == path {
			tracked = f.Tracked
			break
		}
	}
	if err := t.d.WorkingTree.UnStageFile([]string{path}, tracked); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("Failed to unstage: %v", err)}
	}
	t.d.Refresh(ScopeFiles)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已取消暂存: %s", path)}
}

// DiscardFileTool discards all changes in a file (restores to HEAD).
type DiscardFileTool struct{ d *Deps }

func NewDiscardFileTool(d *Deps) tools.Tool { return &DiscardFileTool{d} }

func (t *DiscardFileTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "discard_file",
		Description: "丢弃指定文件的所有变更（恢复到 HEAD）",
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
			return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已丢弃 %s 的所有变更", path)}
		}
	}
	return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("文件不在工作区: %s", path)}
}
