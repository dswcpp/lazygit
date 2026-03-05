package skills

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
)

// BranchNameSkill suggests a branch name based on staged/unstaged changes.
// Extra keys:
//   - "diff" (string, optional) diff of all changes (truncated if too large)
type BranchNameSkill struct{}

func NewBranchNameSkill() Skill { return &BranchNameSkill{} }
func (s *BranchNameSkill) Name() string { return "branch_name" }

func (s *BranchNameSkill) Execute(ctx context.Context, p provider.Provider, input Input) (Output, error) {
	diff := extraStr(input.Extra, "diff")

	var sb strings.Builder
	sb.WriteString("根据以下工作区变更，为新分支推荐一个名称。\n\n")

	// Files summary from context
	staged := []string{}
	unstaged := []string{}
	for _, f := range input.RepoCtx.Files {
		if f.HasStaged {
			staged = append(staged, f.Path)
		} else if f.HasUnstaged {
			unstaged = append(unstaged, f.Path)
		}
	}

	if len(staged) > 0 {
		sb.WriteString("已暂存文件:\n")
		for i, f := range staged {
			if i >= 15 {
				sb.WriteString(fmt.Sprintf("  ... 还有 %d 个\n", len(staged)-15))
				break
			}
			sb.WriteString(fmt.Sprintf("  - %s\n", f))
		}
	}
	if len(unstaged) > 0 {
		sb.WriteString("未暂存文件:\n")
		for i, f := range unstaged {
			if i >= 15 {
				sb.WriteString(fmt.Sprintf("  ... 还有 %d 个\n", len(unstaged)-15))
				break
			}
			sb.WriteString(fmt.Sprintf("  - %s\n", f))
		}
	}

	if diff != "" {
		sb.WriteString("\nDiff 摘要:\n```diff\n")
		sb.WriteString(diff)
		sb.WriteString("\n```\n")
	}

	sb.WriteString("\n命名规则:\n")
	sb.WriteString("- 格式: <type>/<description>（如 feature/add-user-auth）\n")
	sb.WriteString("- type: feature | fix | refactor | docs | test | chore\n")
	sb.WriteString("- description: 小写 kebab-case，2-5 个单词\n")
	sb.WriteString("- 只输出分支名，不要任何解释\n")

	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: "你是一个 Git 分支命名专家。根据变更内容推荐简洁、描述性强的分支名。只输出分支名本身。"},
		{Role: provider.RoleUser, Content: sb.String()},
	}

	result, err := p.Complete(ctx, messages)
	if err != nil {
		return Output{}, err
	}

	name := cleanBranchName(strings.TrimSpace(result.Content))
	if name == "" {
		return Output{}, errors.New("AI returned empty branch name")
	}
	if !strings.Contains(name, "/") {
		name = "feature/" + name
	}
	return Output{Content: name}, nil
}

func cleanBranchName(raw string) string {
	// Strip surrounding quotes/backticks
	raw = strings.Trim(raw, "\"'`")
	raw = strings.ReplaceAll(raw, " ", "-")
	raw = strings.ToLower(raw)
	for _, ch := range []string{"~", "^", ":", "?", "*", "[", "\\", "..", "@{", "//"} {
		raw = strings.ReplaceAll(raw, ch, "")
	}
	return raw
}
