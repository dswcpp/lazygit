package tools

import (
	"fmt"
	"sort"
	"strings"
)

// Registry holds all registered tools and provides lookup by name.
type Registry struct {
	tools map[string]Tool
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{tools: make(map[string]Tool)}
}

// Register adds a tool to the registry. Panics on duplicate names.
func (r *Registry) Register(tool Tool) {
	name := tool.Schema().Name
	if _, exists := r.tools[name]; exists {
		panic(fmt.Sprintf("ai/tools: duplicate tool registration: %q", name))
	}
	r.tools[name] = tool
}

// Clear removes all registered tools. Call before re-registering tools
// (e.g. when resetHelpersAndControllers is called multiple times).
func (r *Registry) Clear() {
	r.tools = make(map[string]Tool)
}

// Get looks up a tool by name.
func (r *Registry) Get(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// All returns all registered tools in deterministic (alphabetical) order.
func (r *Registry) All() []Tool {
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]Tool, 0, len(names))
	for _, name := range names {
		out = append(out, r.tools[name])
	}
	return out
}

// ByMaxPermission returns tools whose permission level is ≤ maxPerm.
func (r *Registry) ByMaxPermission(maxPerm PermissionLevel) []Tool {
	var out []Tool
	for _, t := range r.All() {
		if t.Schema().Permission <= maxPerm {
			out = append(out, t)
		}
	}
	return out
}

// SystemPromptSection generates a formatted tool list for injection into a system prompt.
// Only tools with permission ≤ maxPerm are included.
func (r *Registry) SystemPromptSection(maxPerm PermissionLevel) string {
	tools := r.ByMaxPermission(maxPerm)
	if len(tools) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("## 可用工具\n\n")
	sb.WriteString("在回复中嵌入以下格式的代码块来调用工具：\n\n")
	sb.WriteString("```tool\n{\"name\": \"工具名\", \"params\": {\"参数\": \"值\"}}\n```\n\n")

	// Group by permission level for readability
	groups := map[PermissionLevel][]Tool{}
	for _, t := range tools {
		perm := t.Schema().Permission
		groups[perm] = append(groups[perm], t)
	}

	perms := []PermissionLevel{PermReadOnly, PermWriteLocal, PermWriteRemote, PermDestructive}
	labels := map[PermissionLevel]string{
		PermReadOnly:    "### 只读工具（直接执行）",
		PermWriteLocal:  "### 本地写入工具（需用户确认）",
		PermWriteRemote: "### 远程操作工具（需用户确认）",
		PermDestructive: "### 危险操作工具（需用户确认）",
	}

	for _, perm := range perms {
		group, ok := groups[perm]
		if !ok || len(group) == 0 {
			continue
		}
		sb.WriteString(labels[perm] + "\n")
		for _, t := range group {
			schema := t.Schema()
			paramDesc := formatParamDesc(schema.Params)
			sb.WriteString(fmt.Sprintf("- **%s** — %s%s\n", schema.Name, schema.Description, paramDesc))
		}
		sb.WriteString("\n")
	}

	return strings.TrimRight(sb.String(), "\n")
}

func formatParamDesc(params map[string]ParamSchema) string {
	if len(params) == 0 {
		return ""
	}
	var parts []string
	// Sort for determinism
	names := make([]string, 0, len(params))
	for k := range params {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, name := range names {
		p := params[name]
		req := ""
		if p.Required {
			req = "*"
		}
		parts = append(parts, fmt.Sprintf("%s%s(%s)", name, req, p.Type))
	}
	return " | 参数: " + strings.Join(parts, ", ")
}
