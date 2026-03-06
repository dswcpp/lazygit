package tools

import (
	"context"
	"fmt"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/repocontext"
	"github.com/dswcpp/lazygit/pkg/ai/skills"
)

// SkillTool 将一个 Skill 包装为只读 Tool，供规划阶段 LLM 调用。
//
// 规划阶段的 LLM 可以通过 ```tool 块调用 SkillTool，预先计算出
// 具体的值（例如提交信息、分支名），再将这些值写入执行计划的参数中。
// 这样执行阶段只需按参数直接调用写工具，无需再次询问 LLM。
type SkillTool struct {
	skill      skills.Skill
	prov       provider.Provider
	repoCtxFn  func() repocontext.RepoContext
	schema     ToolSchema
}

// NewSkillTool 创建一个 SkillTool。
//   - skill:     要包装的 Skill 实例
//   - prov:      LLM provider，Skill 内部需要调用 LLM
//   - repoCtxFn: 返回当前仓库上下文的闭包（懒求值，确保总是最新状态）
//   - description: 展示给规划阶段 LLM 的工具描述
//   - params:    工具参数 schema（Skill 所需的 extra 参数）
func NewSkillTool(
	skill skills.Skill,
	prov provider.Provider,
	repoCtxFn func() repocontext.RepoContext,
	description string,
	params map[string]ParamSchema,
) Tool {
	return &SkillTool{
		skill:     skill,
		prov:      prov,
		repoCtxFn: repoCtxFn,
		schema: ToolSchema{
			Name:        skill.Name(),
			Description: description,
			Params:      params,
			Permission:  PermReadOnly, // 所有 SkillTool 均为只读
		},
	}
}

func (t *SkillTool) Schema() ToolSchema { return t.schema }

func (t *SkillTool) Execute(ctx context.Context, call ToolCall) ToolResult {
	// 将 call.Params 转为 skills.Input.Extra（map[string]any 类型一致，直接复用）
	extra := make(map[string]any, len(call.Params))
	for k, v := range call.Params {
		extra[k] = v
	}

	out, err := t.skill.Execute(ctx, t.prov, skills.Input{
		RepoCtx: t.repoCtxFn(),
		Extra:   extra,
	})
	if err != nil {
		return ToolResult{
			CallID:  call.ID,
			Success: false,
			Output:  fmt.Sprintf("skill %s 执行失败: %v", t.schema.Name, err),
		}
	}
	return ToolResult{
		CallID:  call.ID,
		Success: true,
		Output:  out.Content,
	}
}
