package tools

// PermissionLevel classifies the risk level of a tool's operations.
// Higher values represent more potentially destructive or irreversible actions.
type PermissionLevel int

const (
	// PermReadOnly — read-only queries, no side effects (get_status, get_diff, get_log…)
	PermReadOnly PermissionLevel = iota
	// PermWriteLocal — modifies local repository state (commit, stage, branch create/delete…)
	PermWriteLocal
	// PermWriteRemote — communicates with remote (fetch, push)
	PermWriteRemote
	// PermDestructive — hard to reverse (reset --hard, force push, discard changes…)
	PermDestructive
)

// RequiresConfirm reports whether this permission level should prompt the user
// for explicit confirmation before the tool is executed.
func (p PermissionLevel) RequiresConfirm() bool {
	return p >= PermWriteLocal
}

// String returns a human-readable label for the permission level.
func (p PermissionLevel) String() string {
	switch p {
	case PermReadOnly:
		return "只读"
	case PermWriteLocal:
		return "本地写入"
	case PermWriteRemote:
		return "远程写入"
	case PermDestructive:
		return "危险操作"
	default:
		return "未知"
	}
}
