package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// rawPlan mirrors the JSON structure of a ```plan block.
type rawPlan struct {
	Summary string       `json:"summary"`
	Steps   []rawPlanStep `json:"steps"`
}

type rawPlanStep struct {
	ID          string         `json:"id"`
	Description string         `json:"description"`
	Tool        string         `json:"tool"`
	Params      map[string]any `json:"params"`
	Critical    bool           `json:"critical"`
}

// ParsedPlan is the parsed output of a ```plan block.
// Declared here (not in the agent package) to avoid an import cycle:
// tools → agent would be circular since agent already imports tools.
type ParsedPlan struct {
	Summary string
	Steps   []ParsedStep
}

// ParsedStep is a single step parsed from a ```plan block.
type ParsedStep struct {
	ID          string
	Description string
	ToolName    string
	Params      map[string]any
	Critical    bool
}

// ParsePlan extracts the first ```plan ... ``` block from an LLM response and
// returns the structured plan.  Returns false when no valid plan block exists.
func ParsePlan(text string) (ParsedPlan, bool) {
	lines := strings.Split(text, "\n")
	inBlock := false
	var blockLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inBlock {
			if trimmed == "```plan" {
				inBlock = true
				blockLines = blockLines[:0]
			}
			continue
		}
		if trimmed == "```" {
			inBlock = false
			if p, ok := parsePlanJSON(strings.Join(blockLines, "\n")); ok {
				return p, true
			}
			blockLines = nil
			continue
		}
		blockLines = append(blockLines, line)
	}
	return ParsedPlan{}, false
}

// StripPlanBlock removes the first ```plan ... ``` block from text,
// returning only the human-readable portions.
func StripPlanBlock(text string) string {
	var sb strings.Builder
	lines := strings.Split(text, "\n")
	inBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inBlock {
			if trimmed == "```plan" {
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

func parsePlanJSON(raw string) (ParsedPlan, bool) {
	raw = strings.TrimSpace(raw)
	var r rawPlan
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		return ParsedPlan{}, false
	}
	if len(r.Steps) == 0 {
		return ParsedPlan{}, false
	}
	steps := make([]ParsedStep, 0, len(r.Steps))
	for i, s := range r.Steps {
		if s.Tool == "" {
			continue
		}
		id := s.ID
		if id == "" {
			id = fmt.Sprintf("%d", i+1)
		}
		params := s.Params
		if params == nil {
			params = map[string]any{}
		}
		steps = append(steps, ParsedStep{
			ID:          id,
			Description: s.Description,
			ToolName:    s.Tool,
			Params:      params,
			Critical:    s.Critical,
		})
	}
	if len(steps) == 0 {
		return ParsedPlan{}, false
	}
	return ParsedPlan{Summary: r.Summary, Steps: steps}, true
}

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
