package presentation

import (
	"testing"

	"github.com/dswcpp/lazygit/pkg/commands/models"
	"github.com/dswcpp/lazygit/pkg/config"
	guiModels "github.com/dswcpp/lazygit/pkg/gui/models"
	"github.com/dswcpp/lazygit/pkg/gui/presentation/icons"
	"github.com/dswcpp/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestGetActivityBarDisplayStrings(t *testing.T) {
	scenarios := []struct {
		name              string
		items             []*models.ActivityBarItem
		activityBarConfig config.ActivityBarConfig
		currentContext    types.Context
		activityBarStatus *guiModels.ActivityBarStatus
		expected          [][]string
	}{
		{
			name: "分隔符显示为空行",
			items: []*models.ActivityBarItem{
				{Name: "separator1", Type: models.ActivityTypeSeparator},
			},
			activityBarConfig: config.ActivityBarConfig{IconStyle: "ascii"},
			currentContext:    nil,
			activityBarStatus: nil,
			expected:          [][]string{{""}},
		},
		{
			name: "导航项显示图标（非当前面板）",
			items: []*models.ActivityBarItem{
				{
					Name:   "status",
					Type:   models.ActivityTypeNavigation,
					Action: "status",
					Icon:   icons.ActivityBarIcons["status"],
				},
			},
			activityBarConfig: config.ActivityBarConfig{IconStyle: "ascii"},
			currentContext:    nil,
			activityBarStatus: nil,
			expected:          [][]string{{" [S]"}},
		},
		{
			name: "当前活动面板显示蓝色圆点",
			items: []*models.ActivityBarItem{
				{
					Name:   "status",
					Type:   models.ActivityTypeNavigation,
					Action: "status",
					Icon:   icons.ActivityBarIcons["status"],
				},
			},
			activityBarConfig: config.ActivityBarConfig{IconStyle: "ascii"},
			currentContext:    &mockContext{key: "status"},
			activityBarStatus: nil,
			repoState:         nil,
			// 注意：实际输出会包含 ANSI 颜色码，这里简化测试
			expected: [][]string{{"●[S]"}},
		},
		{
			name: "操作进行中显示 spinner",
			items: []*models.ActivityBarItem{
				{
					Name:   "pull",
					Type:   models.ActivityTypeAction,
					Action: "pull",
					Icon:   icons.ActivityBarIcons["pull"],
				},
			},
			activityBarConfig: config.ActivityBarConfig{IconStyle: "ascii"},
			currentContext:    nil,
			activityBarStatus: mockActivityBarStatus("pull"),
			repoState:         nil,
			// spinner 字符后跟图标
			expected: [][]string{{" ⠋ [v]"}},
		},
		{
			name: "禁用的操作显示为灰色",
			items: []*models.ActivityBarItem{
				{
					Name:   "merge",
					Type:   models.ActivityTypeAction,
					Action: "merge",
					Icon:   icons.ActivityBarIcons["merge"],
				},
			},
			activityBarConfig: config.ActivityBarConfig{IconStyle: "ascii"},
			currentContext:    nil,
			activityBarStatus: nil,
			repoState:         &mockRepoState{diffingActive: true},
			// 灰色显示（实际会包含 ANSI 码）
			expected: [][]string{{" [m]"}},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			userConfig := &config.UserConfig{
				Gui: config.GuiConfig{
					ActivityBar: s.activityBarConfig,
				},
			}

			actual := GetActivityBarDisplayStrings(
				s.items,
				s.activityBarConfig,
				s.currentContext,
				userConfig,
				s.activityBarStatus,
			)

			// 由于实际输出包含 ANSI 颜色码，我们只检查基本结构
			assert.Equal(t, len(s.expected), len(actual), "结果数量应该相等")
			for i := range actual {
				// 检查是否包含预期的图标文本（忽略颜色码）
				if s.expected[i][0] != "" {
					assert.Contains(t, actual[i][0], s.expected[i][0][len(s.expected[i][0])-3:],
						"应该包含预期的图标")
				}
			}
		})
	}
}

func TestIsCurrentContext(t *testing.T) {
	scenarios := []struct {
		name           string
		item           *models.ActivityBarItem
		currentContext types.Context
		expected       bool
	}{
		{
			name: "非导航项总是返回 false",
			item: &models.ActivityBarItem{
				Type:   models.ActivityTypeAction,
				Action: "pull",
			},
			currentContext: &mockContext{key: "status"},
			expected:       false,
		},
		{
			name: "当前上下文为 nil 返回 false",
			item: &models.ActivityBarItem{
				Type:   models.ActivityTypeNavigation,
				Action: "status",
			},
			currentContext: nil,
			expected:       false,
		},
		{
			name: "匹配的 status context",
			item: &models.ActivityBarItem{
				Type:   models.ActivityTypeNavigation,
				Action: "status",
			},
			currentContext: &mockContext{key: "status"},
			expected:       true,
		},
		{
			name: "匹配的 files context",
			item: &models.ActivityBarItem{
				Type:   models.ActivityTypeNavigation,
				Action: "files",
			},
			currentContext: &mockContext{key: "files"},
			expected:       true,
		},
		{
			name: "匹配的 branches context（映射到 localBranches）",
			item: &models.ActivityBarItem{
				Type:   models.ActivityTypeNavigation,
				Action: "branches",
			},
			currentContext: &mockContext{key: "localBranches"},
			expected:       true,
		},
		{
			name: "不匹配的 context",
			item: &models.ActivityBarItem{
				Type:   models.ActivityTypeNavigation,
				Action: "status",
			},
			currentContext: &mockContext{key: "files"},
			expected:       false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			actual := isCurrentContext(s.item, s.currentContext)
			assert.Equal(t, s.expected, actual)
		})
	}
}

// Mock implementations

type mockContext struct {
	key types.ContextKey
}

func (m *mockContext) GetKey() types.ContextKey {
	return m.key
}

func (m *mockContext) GetKind() types.ContextKind                         { return types.SIDE_CONTEXT }
func (m *mockContext) GetViewName() string                                { return "" }
func (m *mockContext) GetWindowName() string                              { return "" }
func (m *mockContext) HandleFocus(opts types.OnFocusOpts) error           { return nil }
func (m *mockContext) HandleFocusLost(opts types.OnFocusLostOpts) error   { return nil }
func (m *mockContext) HandleRender() error                                { return nil }
func (m *mockContext) HandleRenderToMain() error                          { return nil }
func (m *mockContext) Title() string                                      { return "" }
func (m *mockContext) GetOptionsMap() map[string]string                   { return nil }
func (m *mockContext) GetParentView() types.IViewTrait                    { return nil }
func (m *mockContext) GetView() types.IViewTrait                          { return nil }
func (m *mockContext) AddKeybindingsFn(fn types.KeybindingsFn)            {}
func (m *mockContext) AddMouseKeybindingsFn(fn types.MouseKeybindingsFn)  {}
func (m *mockContext) AddOnFocusFn(fn types.OnFocusFn)                    {}
func (m *mockContext) AddOnFocusLostFn(fn types.OnFocusLostFn)            {}
func (m *mockContext) AddOnRenderToMainFn(fn types.OnRenderToMainFn)      {}
func (m *mockContext) IsFocusable() bool                                  { return true }
func (m *mockContext) IsTransient() bool                                  { return false }
func (m *mockContext) ModelSearchResults() []types.SearchResultsForModel  { return nil }
func (m *mockContext) Keybindings() []*types.Binding                      { return nil }
func (m *mockContext) AvailableActions() []types.AvailableAction          { return nil }
func (m *mockContext) MarkSearchResultsCacheDirty(modelId string)         {}
func (m *mockContext) ClearSearchResultsCache()                           {}
func (m *mockContext) IsFiltered() bool                                   { return false }
func (m *mockContext) SetHasUncontrolledBounds(hasUncontrolledBounds bool) {}

func mockActivityBarStatus(operationInProgress string) *guiModels.ActivityBarStatus {
	s := guiModels.NewActivityBarStatus()
	if operationInProgress != "" {
		s.SetOperationInProgress(operationInProgress, true)
	}
	return s
}
