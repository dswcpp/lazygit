package gui

import (
	"github.com/dswcpp/lazygit/pkg/commands/models"
	"github.com/dswcpp/lazygit/pkg/config"
	"github.com/dswcpp/lazygit/pkg/gui/presentation/icons"
	"github.com/dswcpp/lazygit/pkg/i18n"
)

// loadActivityBarItems loads activity bar items based on user configuration
func (gui *Gui) loadActivityBarItems() {
	config := gui.UserConfig().Gui.ActivityBar

	// 如果用户定义了自定义列表，使用自定义列表
	if len(config.Items) > 0 {
		gui.State.Model.ActivityBarItems = gui.buildCustomActivityBarItems(config.Items)
		return
	}

	// 否则使用默认列表
	gui.State.Model.ActivityBarItems = gui.buildDefaultActivityBarItems()
}

// buildDefaultActivityBarItems creates the default activity bar items list
func (gui *Gui) buildDefaultActivityBarItems() []*models.ActivityBarItem {
	tr := gui.c.Tr

	return []*models.ActivityBarItem{
		// 导航区
		{
			Name:    "status",
			Type:    models.ActivityTypeNavigation,
			Action:  "status",
			Icon:    icons.ActivityBarIcons["status"],
			Tooltip: getTooltip(tr, "Status"),
		},
		{
			Name:    "files",
			Type:    models.ActivityTypeNavigation,
			Action:  "files",
			Icon:    icons.ActivityBarIcons["files"],
			Tooltip: getTooltip(tr, "Files"),
		},
		{
			Name:    "branches",
			Type:    models.ActivityTypeNavigation,
			Action:  "branches",
			Icon:    icons.ActivityBarIcons["branches"],
			Tooltip: getTooltip(tr, "Branches"),
		},
		{
			Name:    "commits",
			Type:    models.ActivityTypeNavigation,
			Action:  "commits",
			Icon:    icons.ActivityBarIcons["commits"],
			Tooltip: getTooltip(tr, "Commits"),
		},
		{
			Name:    "stash",
			Type:    models.ActivityTypeNavigation,
			Action:  "stash",
			Icon:    icons.ActivityBarIcons["stash"],
			Tooltip: getTooltip(tr, "Stash"),
		},

		// 分隔符
		{Name: "separator1", Type: models.ActivityTypeSeparator},

		// 操作区
		{
			Name:    "pull",
			Type:    models.ActivityTypeAction,
			Action:  "pull",
			Icon:    icons.ActivityBarIcons["pull"],
			Tooltip: getTooltip(tr, "Pull"),
		},
		{
			Name:    "push",
			Type:    models.ActivityTypeAction,
			Action:  "push",
			Icon:    icons.ActivityBarIcons["push"],
			Tooltip: getTooltip(tr, "Push"),
		},
		{
			Name:    "fetch",
			Type:    models.ActivityTypeAction,
			Action:  "fetch",
			Icon:    icons.ActivityBarIcons["fetch"],
			Tooltip: getTooltip(tr, "Fetch"),
		},
		{
			Name:    "stash-action",
			Type:    models.ActivityTypeAction,
			Action:  "stash",
			Icon:    icons.ActivityBarIcons["stash-action"],
			Tooltip: getTooltip(tr, "Stash changes"),
		},
		{
			Name:    "merge",
			Type:    models.ActivityTypeAction,
			Action:  "merge",
			Icon:    icons.ActivityBarIcons["merge"],
			Tooltip: getTooltip(tr, "Merge"),
		},
		{
			Name:    "rebase",
			Type:    models.ActivityTypeAction,
			Action:  "rebase",
			Icon:    icons.ActivityBarIcons["rebase"],
			Tooltip: getTooltip(tr, "Rebase"),
		},

		// 分隔符
		{Name: "separator2", Type: models.ActivityTypeSeparator},

		// 工具区
		{
			Name:    "settings",
			Type:    models.ActivityTypeTool,
			Action:  "settings",
			Icon:    icons.ActivityBarIcons["settings"],
			Tooltip: getTooltip(tr, "Settings"),
		},
		{
			Name:    "help",
			Type:    models.ActivityTypeTool,
			Action:  "help",
			Icon:    icons.ActivityBarIcons["help"],
			Tooltip: getTooltip(tr, "Help"),
		},
	}
}

// buildCustomActivityBarItems creates activity bar items from user configuration
func (gui *Gui) buildCustomActivityBarItems(configItems []config.ActivityBarItemConfig) []*models.ActivityBarItem {
	items := make([]*models.ActivityBarItem, len(configItems))

	for i, configItem := range configItems {
		item := &models.ActivityBarItem{
			Name:      configItem.Name,
			Type:      parseActivityItemType(configItem.Type),
			Action:    configItem.Action,
			CustomCmd: configItem.CustomCmd,
			Tooltip:   configItem.Tooltip,
		}

		// 如果用户提供了自定义图标，使用自定义图标
		if configItem.Icon != "" {
			item.Icon = models.IconConfig{
				NerdFont: configItem.Icon,
				Emoji:    configItem.Icon,
				ASCII:    configItem.Icon,
			}
		} else {
			// 否则使用默认图标（如果存在）
			if defaultIcon, ok := icons.ActivityBarIcons[configItem.Name]; ok {
				item.Icon = defaultIcon
			} else {
				// 如果没有默认图标，使用占位符
				item.Icon = models.IconConfig{
					NerdFont: "•",
					Emoji:    "•",
					ASCII:    "*",
				}
			}
		}

		items[i] = item
	}

	return items
}

// parseActivityItemType converts string type to ActivityItemType
func parseActivityItemType(typeStr string) models.ActivityItemType {
	switch typeStr {
	case "navigation":
		return models.ActivityTypeNavigation
	case "action":
		return models.ActivityTypeAction
	case "tool":
		return models.ActivityTypeTool
	case "separator":
		return models.ActivityTypeSeparator
	case "custom":
		return models.ActivityTypeCustom
	default:
		return models.ActivityTypeNavigation
	}
}

// getTooltip retrieves translation or returns the default value
func getTooltip(tr *i18n.TranslationSet, key string) string {
	// 尝试从翻译中获取，如果不存在则返回 key 本身
	// 这里简化处理，实际应该有更复杂的翻译逻辑
	return key
}
