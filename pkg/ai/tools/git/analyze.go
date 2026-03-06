package gittools

import (
	"context"
	"fmt"
	"strings"

	"github.com/dswcpp/lazygit/pkg/ai/provider"
	"github.com/dswcpp/lazygit/pkg/ai/tools"
)

// AnalyzeChangesTool 智能分析文件变更，逐个文件调用AI分析后整合结果。
// 适用于大模型上下文限制场景，避免一次性发送过大的diff。
type AnalyzeChangesTool struct {
	d        *Deps
	provider provider.Provider
}

func NewAnalyzeChangesTool(d *Deps, p provider.Provider) tools.Tool {
	return &AnalyzeChangesTool{d: d, provider: p}
}

func (t *AnalyzeChangesTool) Schema() tools.ToolSchema {
	return tools.ToolSchema{
		Name:        "analyze_changes",
		Description: "智能分析当前变更：逐个文件分析diff并整合结果（适用于大量变更场景）",
		Params: map[string]tools.ParamSchema{
			"staged": {
				Type:        "bool",
				Description: "true=分析暂存区，false=分析工作区（默认 false）",
			},
			"focus": {
				Type:        "string",
				Description: "分析重点（如：安全问题、性能优化、代码质量等），留空则全面分析",
			},
		},
		Permission: tools.PermReadOnly,
	}
}

func (t *AnalyzeChangesTool) Execute(ctx context.Context, call tools.ToolCall) tools.ToolResult {
	staged := boolParam(call.Params, "staged", false)
	focus := strParam(call.Params, "focus", "")

	// 1. 获取变更文件列表
	files := t.d.GetFiles()
	if len(files) == 0 {
		return tools.ToolResult{
			CallID:  call.ID,
			Success: true,
			Output:  "工作区干净，没有变更文件",
		}
	}

	// 2. 筛选需要分析的文件
	var targetFiles []string
	for _, f := range files {
		if staged && f.HasStagedChanges {
			targetFiles = append(targetFiles, f.Path)
		} else if !staged && f.HasUnstagedChanges {
			targetFiles = append(targetFiles, f.Path)
		}
	}

	if len(targetFiles) == 0 {
		label := "工作区"
		if staged {
			label = "暂存区"
		}
		return tools.ToolResult{
			CallID:  call.ID,
			Success: true,
			Output:  fmt.Sprintf("%s没有变更", label),
		}
	}

	// 3. 逐个文件分析
	var analyses []fileAnalysis
	for _, path := range targetFiles {
		if ctx.Err() != nil {
			return tools.ToolResult{
				CallID: call.ID,
				Output: "分析被取消",
			}
		}

		analysis, err := t.analyzeFile(ctx, path, staged, focus)
		if err != nil {
			// 单个文件失败不中断整体分析
			analyses = append(analyses, fileAnalysis{
				Path:  path,
				Error: err.Error(),
			})
			continue
		}
		analyses = append(analyses, analysis)
	}

	// 4. 整合结果
	summary := t.buildSummary(analyses, focus)

	return tools.ToolResult{
		CallID:  call.ID,
		Success: true,
		Output:  summary,
	}
}

type fileAnalysis struct {
	Path     string
	Summary  string // AI生成的单文件分析摘要
	Error    string // 如果分析失败，记录错误信息
	LinesDiff int   // diff行数
}

// analyzeFile 调用AI分析单个文件的diff
func (t *AnalyzeChangesTool) analyzeFile(ctx context.Context, path string, staged bool, focus string) (fileAnalysis, error) {
	// 获取文件diff
	var diff string
	for _, f := range t.d.GetFiles() {
		if f.Path == path {
			diff = t.d.WorkingTree.WorktreeFileDiff(f, true, staged)
			break
		}
	}

	if diff == "" {
		return fileAnalysis{
			Path:    path,
			Summary: "无变更",
		}, nil
	}

	// 构建分析提示词
	prompt := t.buildAnalysisPrompt(path, diff, focus)

	// 调用AI分析
	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: "你是代码审查专家，擅长分析代码变更。请简洁、准确地分析diff内容。"},
		{Role: provider.RoleUser, Content: prompt},
	}

	result, err := t.provider.Complete(ctx, messages)
	if err != nil {
		return fileAnalysis{}, fmt.Errorf("AI分析失败: %w", err)
	}

	return fileAnalysis{
		Path:      path,
		Summary:   strings.TrimSpace(result.Content),
		LinesDiff: len(strings.Split(diff, "\n")),
	}, nil
}

// buildAnalysisPrompt 构建单文件分析提示词
func (t *AnalyzeChangesTool) buildAnalysisPrompt(path, diff, focus string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## 文件: %s\n\n", path))
	sb.WriteString("请分析以下diff，用2-3句话总结：\n")

	if focus != "" {
		sb.WriteString(fmt.Sprintf("**分析重点**: %s\n\n", focus))
	} else {
		sb.WriteString("- 主要变更内容\n")
		sb.WriteString("- 潜在问题（如有）\n")
		sb.WriteString("- 改进建议（如有）\n\n")
	}

	sb.WriteString("```diff\n")
	sb.WriteString(diff)
	sb.WriteString("\n```\n")

	return sb.String()
}

// buildSummary 整合所有文件的分析结果
func (t *AnalyzeChangesTool) buildSummary(analyses []fileAnalysis, focus string) string {
	var sb strings.Builder

	// 标题
	if focus != "" {
		sb.WriteString(fmt.Sprintf("# 变更分析报告（重点：%s）\n\n", focus))
	} else {
		sb.WriteString("# 变更分析报告\n\n")
	}

	// 统计信息
	successCount := 0
	failCount := 0
	totalLines := 0
	for _, a := range analyses {
		if a.Error == "" {
			successCount++
			totalLines += a.LinesDiff
		} else {
			failCount++
		}
	}

	sb.WriteString(fmt.Sprintf("**文件数**: %d 个（成功分析 %d，失败 %d）\n", len(analyses), successCount, failCount))
	sb.WriteString(fmt.Sprintf("**总变更行数**: 约 %d 行\n\n", totalLines))

	// 逐个文件的分析结果
	sb.WriteString("## 详细分析\n\n")
	for i, a := range analyses {
		sb.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, a.Path))

		if a.Error != "" {
			sb.WriteString(fmt.Sprintf("❌ **分析失败**: %s\n\n", a.Error))
		} else if a.Summary == "无变更" {
			sb.WriteString("ℹ️ 无变更\n\n")
		} else {
			sb.WriteString(a.Summary)
			sb.WriteString("\n\n")
		}
	}

	// 整体建议（可选）
	if successCount > 1 {
		sb.WriteString("## 整体建议\n\n")
		sb.WriteString("建议在提交前：\n")
		sb.WriteString("1. 确认所有变更符合预期\n")
		sb.WriteString("2. 运行测试确保功能正常\n")
		sb.WriteString("3. 检查是否有遗漏的文件\n")
	}

	return sb.String()
}
