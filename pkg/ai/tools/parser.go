package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseToolCalls extracts all tool call blocks from an LLM response.
// It looks for fenced code blocks with the "tool" language tag:
//
//	```tool
//	{"name": "stage_file", "params": {"path": "main.go"}}
//	```
func ParseToolCalls(text string) []ToolCall {
	var calls []ToolCall
	counter := 0

	lines := strings.Split(text, "\n")
	inBlock := false
	var blockLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inBlock {
			if trimmed == "```tool" {
				inBlock = true
				blockLines = blockLines[:0]
			}
			continue
		}
		if trimmed == "```" {
			inBlock = false
			if call, ok := parseToolCallJSON(strings.Join(blockLines, "\n"), counter); ok {
				calls = append(calls, call)
				counter++
			}
			blockLines = nil
			continue
		}
		blockLines = append(blockLines, line)
	}

	return calls
}

// StripToolBlocks removes all ```tool ... ``` blocks from text, returning
// only the human-readable portions of the response.
func StripToolBlocks(text string) string {
	var sb strings.Builder
	lines := strings.Split(text, "\n")
	inBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inBlock {
			if trimmed == "```tool" {
				inBlock = true
				continue
			}
			sb.WriteString(line)
			sb.WriteString("\n")
			continue
		}
		if trimmed == "```" {
			inBlock = false
		}
	}

	return strings.TrimSpace(sb.String())
}

type rawToolCall struct {
	Name   string         `json:"name"`
	Params map[string]any `json:"params"`
}

func parseToolCallJSON(raw string, idx int) (ToolCall, bool) {
	raw = strings.TrimSpace(raw)
	var r rawToolCall
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		return ToolCall{}, false
	}
	if r.Name == "" {
		return ToolCall{}, false
	}
	if r.Params == nil {
		r.Params = map[string]any{}
	}
	return ToolCall{
		ID:     fmt.Sprintf("call_%d", idx),
		Name:   r.Name,
		Params: r.Params,
	}, true
}
