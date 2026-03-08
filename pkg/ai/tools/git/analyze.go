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
		Description: t.d.Tr.AnalyzeToolDescription(),
		Params: map[string]tools.ParamSchema{
			"staged": {
				Type:        "bool",
				Description: t.d.Tr.AnalyzeToolStagedParam(),
			},
			"focus": {
				Type:        "string",
				Description: t.d.Tr.AnalyzeToolFocusParam(),
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
			Output:  t.d.Tr.AnalyzeWorkingDirClean(),
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
		label := t.d.Tr.ToolWorkingDir()
		if staged {
			label = t.d.Tr.ToolStagingArea()
		}
		return tools.ToolResult{
			CallID:  call.ID,
			Success: true,
			Output:  t.d.Tr.AnalyzeNoChanges(label),
		}
	}

	// 3. 逐个文件分析
	var analyses []fileAnalysis
	for _, path := range targetFiles {
		if ctx.Err() != nil {
			return tools.ToolResult{
				CallID: call.ID,
				Output: t.d.Tr.AnalyzeCancelled(),
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
			Summary: t.d.Tr.ToolNoChanges(),
		}, nil
	}

	// 构建分析提示词
	prompt := t.buildAnalysisPrompt(path, diff, focus)

	// 调用AI分析
	messages := []provider.Message{
		{Role: provider.RoleSystem, Content: t.d.Tr.AnalyzeCodeReviewExpert()},
		{Role: provider.RoleUser, Content: prompt},
	}

	result, err := t.provider.Complete(ctx, messages)
	if err != nil {
		return fileAnalysis{}, fmt.Errorf("%s", t.d.Tr.AnalyzeFailed(err))
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

	sb.WriteString(t.d.Tr.AnalyzeFileLabel(path))
	sb.WriteString(t.d.Tr.AnalyzePromptIntro())

	if focus != "" {
		sb.WriteString(t.d.Tr.AnalyzeFocusLabel(focus))
	} else {
		sb.WriteString(t.d.Tr.AnalyzeMainChanges())
		sb.WriteString(t.d.Tr.AnalyzePotentialIssues())
		sb.WriteString(t.d.Tr.AnalyzeImprovementSuggestions())
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
		sb.WriteString(t.d.Tr.AnalyzeReportTitleWithFocus(focus))
	} else {
		sb.WriteString(t.d.Tr.AnalyzeReportTitle())
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

	sb.WriteString(t.d.Tr.AnalyzeFileCount(len(analyses), successCount, failCount))
	sb.WriteString(t.d.Tr.AnalyzeTotalLines(totalLines))

	// 逐个文件的分析结果
	sb.WriteString(t.d.Tr.AnalyzeDetailedAnalysis())
	for i, a := range analyses {
		sb.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, a.Path))

		if a.Error != "" {
			sb.WriteString(t.d.Tr.AnalyzeAnalysisFailed(a.Error))
		} else if a.Summary == t.d.Tr.ToolNoChanges() {
			sb.WriteString(t.d.Tr.AnalyzeNoChangesInfo())
		} else {
			sb.WriteString(a.Summary)
			sb.WriteString("\n\n")
		}
	}

	// 整体建议（可选）
	if successCount > 1 {
		sb.WriteString(t.d.Tr.AnalyzeOverallSuggestions())
		sb.WriteString(t.d.Tr.AnalyzeSuggestion1())
		sb.WriteString(t.d.Tr.AnalyzeSuggestion2())
		sb.WriteString(t.d.Tr.AnalyzeSuggestion3())
	}

	return sb.String()
}
