package icons

import (
	"github.com/dswcpp/lazygit/pkg/commands/models"
	"github.com/dswcpp/lazygit/pkg/config"
)

// ActivityBarIcons defines icons for each activity bar item (Nerd Fonts v3)
var ActivityBarIcons = map[string]models.IconConfig{
	// Navigation区
	"status": {
		NerdFont: "\uf0c9", //  (三横线/列表)
		Emoji:    "📋",
		ASCII:    "[S]",
		Color:    "#89ddff",
	},
	"files": {
		NerdFont: "\uf15b", //  (文件)
		Emoji:    "📁",
		ASCII:    "[F]",
		Color:    "#ffcb6b",
	},
	"branches": {
		NerdFont: BRANCH_ICON, // 󰘬 (分支) - 复用现有定义
		Emoji:    "🌿",
		ASCII:    "[B]",
		Color:    "#c3e88d",
	},
	"commits": {
		NerdFont: COMMIT_ICON, // 󰜘 (提交) - 复用现有定义
		Emoji:    "💾",
		ASCII:    "[C]",
		Color:    "#82aaff",
	},
	"stash": {
		NerdFont: STASH_ICON, //  (收藏箱) - 复用现有定义
		Emoji:    "📦",
		ASCII:    "[T]",
		Color:    "#c792ea",
	},

	// 操作区
	"pull": {
		NerdFont: "\uf01a", //  (下载箭头)
		Emoji:    "⬇️",
		ASCII:    "[v]",
		Color:    "#89ddff",
	},
	"push": {
		NerdFont: "\uf01b", //  (上传箭头)
		Emoji:    "⬆️",
		ASCII:    "[^]",
		Color:    "#89ddff",
	},
	"fetch": {
		NerdFont: "\uf021", //  (刷新/同步)
		Emoji:    "🔄",
		ASCII:    "[r]",
		Color:    "#ffcb6b",
	},
	"stash-action": {
		NerdFont: "\uf0c7", //  (软盘/保存)
		Emoji:    "💾",
		ASCII:    "[s]",
		Color:    "#c792ea",
	},
	"merge": {
		NerdFont: MERGE_COMMIT_ICON, // 󰘭 (合并) - 复用现有定义
		Emoji:    "🔀",
		ASCII:    "[m]",
		Color:    "#c3e88d",
	},
	"rebase": {
		NerdFont: "\uf126", //  (代码分支)
		Emoji:    "♻️",
		ASCII:    "[R]",
		Color:    "#ffcb6b",
	},

	// 工具区
	"settings": {
		NerdFont: "\uf013", //  (齿轮)
		Emoji:    "⚙️",
		ASCII:    "[*]",
		Color:    "#f78c6c",
	},
	"help": {
		NerdFont: "\uf059", //  (问号圆圈)
		Emoji:    "❓",
		ASCII:    "[?]",
		Color:    "#89ddff",
	},
}

// PatchActivityBarIconsForNerdFontsV2 patches icons for Nerd Fonts v2 compatibility
func PatchActivityBarIconsForNerdFontsV2() {
	// 分支图标 v2 补丁 (已在 git_icons.go 中处理)
	branchIcon := ActivityBarIcons["branches"]
	branchIcon.NerdFont = "\ue725" // v2 分支图标
	ActivityBarIcons["branches"] = branchIcon

	commitIcon := ActivityBarIcons["commits"]
	commitIcon.NerdFont = "\ufc16" // ﰖ v2 提交图标
	ActivityBarIcons["commits"] = commitIcon

	mergeIcon := ActivityBarIcons["merge"]
	mergeIcon.NerdFont = "\ufb2c" // שּׁ v2 合并图标
	ActivityBarIcons["merge"] = mergeIcon
}

// GetActivityBarIcon returns the appropriate icon based on configuration
func GetActivityBarIcon(name string, userConfig *config.UserConfig) string {
	iconCfg, ok := ActivityBarIcons[name]
	if !ok {
		return "?"
	}

	// 根据配置选择图标类型
	iconStyle := userConfig.Gui.ActivityBar.IconStyle
	if iconStyle == "" {
		iconStyle = "auto"
	}

	switch iconStyle {
	case "nerd":
		return iconCfg.NerdFont
	case "emoji":
		return iconCfg.Emoji
	case "ascii":
		return iconCfg.ASCII
	case "auto":
		// 自动检测：如果启用了 Nerd Fonts，使用 Nerd Font 图标
		if userConfig.Gui.NerdFontsVersion != "" {
			return iconCfg.NerdFont
		}
		// 否则使用 emoji
		return iconCfg.Emoji
	default:
		// 默认使用 Nerd Fonts 如果可用，否则 emoji
		if userConfig.Gui.NerdFontsVersion != "" {
			return iconCfg.NerdFont
		}
		return iconCfg.Emoji
	}
}

// GetActivityBarIconColor returns the icon color (for future color support)
func GetActivityBarIconColor(name string) string {
	if iconCfg, ok := ActivityBarIcons[name]; ok {
		return iconCfg.Color
	}
	return ""
}
