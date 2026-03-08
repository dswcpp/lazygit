package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type StringOccurrence struct {
	Text     string   `json:"text"`
	Count    int      `json:"count"`
	Files    []string `json:"files"`
	Category string   `json:"category"`
}

func main() {
	// 匹配中文字符串的正则表达式
	stringPattern := regexp.MustCompile(`"([^"]*[\p{Han}]+[^"]*)"`)
	backtickPattern := regexp.MustCompile("` + \"`\" + `([^` + \"`\" + `]*[\\p{Han}]+[^` + \"`\" + `]*)` + \"`\" + `")

	occurrences := make(map[string]*StringOccurrence)

	// 遍历 pkg 目录下的所有 .go 文件
	err := filepath.Walk("pkg", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过测试文件和翻译文件
		if !strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, "_test.go") ||
			strings.Contains(path, "i18n/translations") ||
			strings.Contains(path, "i18n\\translations") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			// 跳过注释行
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") {
				continue
			}

			// 提取双引号字符串
			matches := stringPattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) > 1 {
					text := match[1]
					recordOccurrence(occurrences, text, path)
				}
			}

			// 提取反引号字符串
			matches = backtickPattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) > 1 {
					text := match[1]
					recordOccurrence(occurrences, text, path)
				}
			}
		}

		return scanner.Err()
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// 分类和排序
	categorizeStrings(occurrences)

	// 生成报告
	generateReport(occurrences)
}

func recordOccurrence(occurrences map[string]*StringOccurrence, text, file string) {
	// 清理文本
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}

	if occ, exists := occurrences[text]; exists {
		occ.Count++
		if !contains(occ.Files, file) {
			occ.Files = append(occ.Files, file)
		}
	} else {
		occurrences[text] = &StringOccurrence{
			Text:  text,
			Count: 1,
			Files: []string{file},
		}
	}
}

func categorizeStrings(occurrences map[string]*StringOccurrence) {
	for _, occ := range occurrences {
		// 根据文件路径和内容特征分类
		if len(occ.Files) > 0 {
			firstFile := occ.Files[0]
			switch {
			case strings.Contains(firstFile, "pkg/ai/agent") || strings.Contains(firstFile, "pkg\\ai\\agent"):
				occ.Category = "AI Agent"
			case strings.Contains(firstFile, "pkg/ai/tools") || strings.Contains(firstFile, "pkg\\ai\\tools"):
				occ.Category = "AI Tools"
			case strings.Contains(firstFile, "pkg/ai/skills") || strings.Contains(firstFile, "pkg\\ai\\skills"):
				occ.Category = "AI Skills"
			case strings.Contains(firstFile, "pkg/gui") || strings.Contains(firstFile, "pkg\\gui"):
				occ.Category = "GUI"
			case strings.Contains(firstFile, "pkg/commands") || strings.Contains(firstFile, "pkg\\commands"):
				occ.Category = "Commands"
			default:
				occ.Category = "Other"
			}
		}
	}
}

func generateReport(occurrences map[string]*StringOccurrence) {
	// 按分类统计
	categoryStats := make(map[string]int)
	categoryStrings := make(map[string][]*StringOccurrence)

	for _, occ := range occurrences {
		categoryStats[occ.Category]++
		categoryStrings[occ.Category] = append(categoryStrings[occ.Category], occ)
	}

	// 排序分类
	categories := make([]string, 0, len(categoryStats))
	for cat := range categoryStats {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	// 生成 Markdown 报告
	fmt.Println("# 中文字符串统计报告")
	fmt.Println()
	fmt.Printf("**总计**: %d 个唯一中文字符串\n\n", len(occurrences))

	fmt.Println("## 按分类统计")
	fmt.Println()
	for _, cat := range categories {
		fmt.Printf("- **%s**: %d 个字符串\n", cat, categoryStats[cat])
	}

	fmt.Println()
	fmt.Println("## 详细列表")
	fmt.Println()

	for _, cat := range categories {
		fmt.Printf("### %s (%d 个)\n\n", cat, categoryStats[cat])

		// 排序该分类下的字符串（按出现次数降序）
		strs := categoryStrings[cat]
		sort.Slice(strs, func(i, j int) bool {
			return strs[i].Count > strs[j].Count
		})

		// 只显示前 20 个最常见的
		limit := 20
		if len(strs) < limit {
			limit = len(strs)
		}

		for i := 0; i < limit; i++ {
			occ := strs[i]
			fmt.Printf("- `%s` (出现 %d 次)\n",
				truncate(occ.Text, 80), occ.Count)
			if len(occ.Files) <= 3 {
				for _, file := range occ.Files {
					fmt.Printf("  - %s\n", file)
				}
			} else {
				fmt.Printf("  - %s (及其他 %d 个文件)\n",
					occ.Files[0], len(occ.Files)-1)
			}
		}

		if len(strs) > limit {
			fmt.Printf("\n... 还有 %d 个字符串\n", len(strs)-limit)
		}
		fmt.Println()
	}

	// 生成 JSON 文件供后续处理
	jsonData, _ := json.MarshalIndent(occurrences, "", "  ")
	os.WriteFile("chinese_strings.json", jsonData, 0644)
	fmt.Println()
	fmt.Println("完整数据已保存到 `chinese_strings.json`")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
