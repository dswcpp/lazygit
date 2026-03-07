package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type StringOccurrence struct {
	Text     string   `json:"text"`
	Count    int      `json:"count"`
	Files    []string `json:"files"`
	Category string   `json:"category"`
}

type TranslationKey struct {
	Key          string
	EnglishText  string
	ChineseText  string
	Category     string
	Count        int
	IsFormatted  bool
	Files        []string
}

func main() {
	// 读取 JSON 数据
	data, err := os.ReadFile("chinese_strings.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var occurrences map[string]*StringOccurrence
	if err := json.Unmarshal(data, &occurrences); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// 生成翻译键
	keys := generateTranslationKeys(occurrences)

	// 按分类和频率排序
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Category != keys[j].Category {
			return keys[i].Category < keys[j].Category
		}
		return keys[i].Count > keys[j].Count
	})

	// 生成输出
	generateGoStructFields(keys)
	fmt.Println("\n" + strings.Repeat("=", 80) + "\n")
	generateEnglishTranslations(keys)
	fmt.Println("\n" + strings.Repeat("=", 80) + "\n")
	generateChineseJSON(keys)
}

func generateTranslationKeys(occurrences map[string]*StringOccurrence) []*TranslationKey {
	var keys []*TranslationKey

	for text, occ := range occurrences {
		key := generateKeyName(text, occ.Category)
		englishText := translateToEnglish(text)

		keys = append(keys, &TranslationKey{
			Key:          key,
			EnglishText:  englishText,
			ChineseText:  text,
			Category:     occ.Category,
			Count:        occ.Count,
			IsFormatted:  strings.Contains(text, "%"),
			Files:        occ.Files,
		})
	}

	return keys
}

func generateKeyName(text, category string) string {
	// 清理文本
	cleaned := strings.TrimSpace(text)
	cleaned = regexp.MustCompile(`[\\n\\t]+`).ReplaceAllString(cleaned, "")
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")

	// 特殊处理常见词汇
	commonWords := map[string]string{
		"取消":     "Cancel",
		"确认":     "Confirm",
		"确定":     "OK",
		"成功":     "Success",
		"失败":     "Failed",
		"错误":     "Error",
		"警告":     "Warning",
		"是":      "Yes",
		"否":      "No",
		"未知":     "Unknown",
		"执行中":    "Executing",
		"思考中":    "Thinking",
		"空闲":     "Idle",
		"已取消":    "Cancelled",
		"可输入下一条指令": "CanInputNext",
	}

	if key, ok := commonWords[cleaned]; ok {
		return categoryPrefix(category) + key
	}

	// 生成描述性键名
	if len(cleaned) > 50 {
		// 长文本，使用摘要
		return categoryPrefix(category) + summarizeText(cleaned)
	}

	// 短文本，直接翻译
	words := []rune(cleaned)
	if len(words) <= 10 {
		return categoryPrefix(category) + transliterateShort(cleaned)
	}

	return categoryPrefix(category) + summarizeText(cleaned)
}

func categoryPrefix(category string) string {
	switch category {
	case "AI Agent":
		return "AIAgent"
	case "AI Tools":
		return "AITool"
	case "AI Skills":
		return "AISkill"
	case "GUI":
		return "AI"
	default:
		return "AI"
	}
}

func summarizeText(text string) string {
	// 提取关键词
	keywords := extractKeywords(text)
	if len(keywords) == 0 {
		return "Text" + fmt.Sprintf("%d", hashString(text)%1000)
	}

	result := ""
	for _, kw := range keywords {
		if len(result) > 0 {
			result += capitalize(kw)
		} else {
			result += kw
		}
	}

	if len(result) > 40 {
		result = result[:40]
	}

	return result
}

func extractKeywords(text string) []string {
	// 简单的关键词提取
	keywordMap := map[string]string{
		"缺少":   "Missing",
		"参数":   "Param",
		"失败":   "Failed",
		"成功":   "Success",
		"错误":   "Error",
		"工具":   "Tool",
		"分支":   "Branch",
		"提交":   "Commit",
		"文件":   "File",
		"路径":   "Path",
		"名称":   "Name",
		"消息":   "Message",
		"标签":   "Tag",
		"远程":   "Remote",
		"本地":   "Local",
		"暂存":   "Stage",
		"变更":   "Change",
		"冲突":   "Conflict",
		"合并":   "Merge",
		"推送":   "Push",
		"拉取":   "Pull",
		"执行":   "Execute",
		"规划":   "Planning",
		"步骤":   "Step",
		"超时":   "Timeout",
		"关键":   "Critical",
		"用户":   "User",
		"拒绝":   "Rejected",
		"AI":   "AI",
		"未启用": "NotEnabled",
		"生成":   "Generate",
		"分析":   "Analyze",
	}

	var keywords []string
	for cn, en := range keywordMap {
		if strings.Contains(text, cn) {
			keywords = append(keywords, en)
		}
	}

	return keywords
}

func transliterateShort(text string) string {
	// 对短文本进行音译或意译
	mapping := map[string]string{
		"取消":  "Cancel",
		"确认":  "Confirm",
		"确定":  "OK",
		"成功":  "Success",
		"失败":  "Failed",
		"是":   "Yes",
		"否":   "No",
		"未知":  "Unknown",
		"空闲":  "Idle",
		"执行中": "Executing",
		"思考中": "Thinking",
	}

	if en, ok := mapping[text]; ok {
		return en
	}

	return "Text"
}

func translateToEnglish(text string) string {
	// 简单的翻译映射（实际应该使用翻译 API）
	translations := map[string]string{
		"取消":           "Cancel",
		"确认":           "Confirm",
		"确定":           "OK",
		"成功":           "Success",
		"失败":           "Failed",
		"错误":           "Error",
		"警告":           "Warning",
		"是":            "Yes",
		"否":            "No",
		"未知":           "Unknown",
		"执行中":          "Executing",
		"思考中":          "Thinking",
		"空闲":           "Idle",
		"已取消":          "Cancelled",
		"可输入下一条指令":     "You can input the next command",
		"AI 未启用":       "AI not enabled",
		"缺少 name 参数":   "Missing name parameter",
		"缺少 path 参数":   "Missing path parameter",
		"缺少 message 参数": "Missing message parameter",
		"缺少 hash 参数":   "Missing hash parameter",
		"文件路径":         "File path",
		"分支名称":         "Branch name",
		"标签名称":         "Tag name",
		"提交信息":         "Commit message",
		"无变更":          "No changes",
		"工作区":          "Working directory",
		"暂存区":          "Staging area",
	}

	// 直接匹配
	if en, ok := translations[text]; ok {
		return en
	}

	// 格式化字符串处理
	for cn, en := range translations {
		if strings.Contains(text, cn) {
			result := strings.ReplaceAll(text, cn, en)
			return result
		}
	}

	// 默认返回原文（需要手动翻译）
	return "[TODO: " + truncate(text, 50) + "]"
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func hashString(s string) int {
	h := 0
	for _, c := range s {
		h = 31*h + int(c)
	}
	if h < 0 {
		h = -h
	}
	return h
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func generateGoStructFields(keys []*TranslationKey) {
	fmt.Println("// Add these fields to TranslationSet in pkg/i18n/english.go")
	fmt.Println()

	currentCategory := ""
	for _, key := range keys {
		if key.Category != currentCategory {
			currentCategory = key.Category
			fmt.Printf("\n\t// %s\n", currentCategory)
		}
		fmt.Printf("\t%-50s string\n", key.Key)
	}
}

func generateEnglishTranslations(keys []*TranslationKey) {
	fmt.Println("// Add these to EnglishTranslationSet() in pkg/i18n/english.go")
	fmt.Println()

	currentCategory := ""
	for _, key := range keys {
		if key.Category != currentCategory {
			currentCategory = key.Category
			fmt.Printf("\n\t\t// %s\n", currentCategory)
		}
		fmt.Printf("\t\t%-50s \"%s\",\n", key.Key+":", escapeString(key.EnglishText))
	}
}

func generateChineseJSON(keys []*TranslationKey) {
	fmt.Println("// Add these to pkg/i18n/translations/zh-CN.json")
	fmt.Println("{")

	for i, key := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Printf("  \"%s\": \"%s\"%s\n", key.Key, escapeJSON(key.ChineseText), comma)
	}

	fmt.Println("}")
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
