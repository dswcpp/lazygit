package gittools

import (
	"context"
	"fmt"

	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

// CreateTagTool creates a lightweight tag.
type CreateTagTool struct{ d *Deps }

func NewCreateTagTool(d *Deps) tools.Tool { return &CreateTagTool{d} }

func (t *CreateTagTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "create_tag",
		Description: "创建轻量 tag",
		Params: map[string]tools.ParamSchema{
			"name": {Type: "string", Description: t.d.Tr.ToolTagName(), Required: true},
			"ref":  {Type: "string", Description: t.d.Tr.ToolTargetRef()},
		},
		Permission: tools.PermWriteLocal,
	}
}

func (t *CreateTagTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	name := strParam(call.Params, "name", "")
	if name == "" {
		return tools.ToolResult{CallID: call.ID, Output: "缺少 name 参数"}
	}
	ref := strParam(call.Params, "ref", "HEAD")
	if err := t.d.Tag.CreateLightweightObj(name, ref, false).Run(); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("Failed to create tag: %v", err)}
	}
	t.d.Refresh(ScopeTags)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已创建 tag %s（指向 %s）", name, ref)}
}

// DeleteTagTool deletes a local tag.
type DeleteTagTool struct{ d *Deps }

func NewDeleteTagTool(d *Deps) tools.Tool { return &DeleteTagTool{d} }

func (t *DeleteTagTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "delete_tag",
		Description: "删除本地 tag",
		Params: map[string]tools.ParamSchema{
			"name": {Type: "string", Description: t.d.Tr.ToolTagName(), Required: true},
		},
		Permission: tools.PermDestructive,
	}
}

func (t *DeleteTagTool) Execute(_ context.Context, call tools.ToolCall) tools.ToolResult {
	name := strParam(call.Params, "name", "")
	if name == "" {
		return tools.ToolResult{CallID: call.ID, Output: "缺少 name 参数"}
	}
	if err := t.d.Tag.LocalDelete(name); err != nil {
		return tools.ToolResult{CallID: call.ID, Output: fmt.Sprintf("Failed to delete tag: %v", err)}
	}
	t.d.Refresh(ScopeTags)
	return tools.ToolResult{CallID: call.ID, Success: true, Output: fmt.Sprintf("已删除本地 tag: %s", name)}
}
