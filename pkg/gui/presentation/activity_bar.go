package presentation

import (
	"github.com/dswcpp/lazygit/pkg/commands/models"
	"github.com/dswcpp/lazygit/pkg/config"
	"github.com/dswcpp/lazygit/pkg/gui/presentation/icons"
	"github.com/dswcpp/lazygit/pkg/gui/style"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

// GetActivityBarDisplayStrings returns display strings for activity bar items
func GetActivityBarDisplayStrings(
	items []*models.ActivityBarItem,
	activityBarConfig config.ActivityBarConfig,
	currentContext types.Context,
	userConfig *config.UserConfig,
	activityBarStatus types.IActivityBarStatus,
) [][]string {
	result := make([][]string, len(items))

	for i, item := range items {
		// 分隔符显示为空行
		if item.IsSeparator() {
			result[i] = []string{""}
			continue
		}

		// 获取图标
		icon := icons.GetActivityBarIcon(item.Name, userConfig)

		// 添加操作进行中的旋转动画
		if activityBarStatus != nil && item.IsActionItem() && activityBarStatus.IsOperationInProgress(item.Action) {
			spinner := activityBarStatus.GetSpinnerChar()
			icon = spinner + " " + icon
		}

		// 添加当前面板指示器
		prefix := " "
		if isCurrentContext(item, currentContext) {
			prefix = style.FgBlue.SetBold().Sprint("●")
		}

		displayIcon := prefix + icon
		result[i] = []string{displayIcon}
	}

	return result
}

// isCurrentContext checks if the item represents the current active context
func isCurrentContext(item *models.ActivityBarItem, currentCtx types.Context) bool {
	if !item.IsNavigationItem() {
		return false
	}

	if currentCtx == nil {
		return false
	}

	// 映射 activity bar action 到 context key
	contextKeyMap := map[string]string{
		"status":   "status",
		"files":    "files",
		"branches": "localBranches",
		"commits":  "commits",
		"stash":    "stash",
	}

	expectedKey := contextKeyMap[item.Action]
	currentKey := string(currentCtx.GetKey())

	return expectedKey == currentKey
}


