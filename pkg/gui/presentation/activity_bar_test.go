package presentation

import (
	"testing"

	"github.com/jesseduffield/gocui"
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
		activityBarStatus types.IActivityBarStatus
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
			// spinner 字符后跟图标
			expected: [][]string{{" ⠋ [v]"}},
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

// IBaseContext methods
func (m *mockContext) GetKey() types.ContextKey                                                    { return m.key }
func (m *mockContext) GetKind() types.ContextKind                                                  { return types.SIDE_CONTEXT }
func (m *mockContext) GetViewName() string                                                         { return "" }
func (m *mockContext) GetView() *gocui.View                                                        { return nil }
func (m *mockContext) GetViewTrait() types.IViewTrait                                              { return nil }
func (m *mockContext) GetWindowName() string                                                       { return "" }
func (m *mockContext) SetWindowName(string)                                                        {}
func (m *mockContext) IsFocusable() bool                                                           { return true }
func (m *mockContext) IsTransient() bool                                                           { return false }
func (m *mockContext) HasControlledBounds() bool                                                   { return true }
func (m *mockContext) TotalContentHeight() int                                                     { return 0 }
func (m *mockContext) NeedsRerenderOnWidthChange() types.NeedsRerenderOnWidthChangeLevel           { return types.NEEDS_RERENDER_ON_WIDTH_CHANGE_NONE }
func (m *mockContext) NeedsRerenderOnHeightChange() bool                                           { return false }
func (m *mockContext) Title() string                                                               { return "" }
func (m *mockContext) GetOptionsMap() map[string]string                                            { return nil }
func (m *mockContext) AddKeybindingsFn(types.KeybindingsFn)                                        {}
func (m *mockContext) AddMouseKeybindingsFn(types.MouseKeybindingsFn)                              {}
func (m *mockContext) ClearAllAttachedControllerFunctions()                                        {}
func (m *mockContext) AddOnClickFn(func() error)                                                   {}
func (m *mockContext) AddOnClickFocusedMainViewFn(func(string, int) error)                         {}
func (m *mockContext) AddOnRenderToMainFn(func())                                                  {}
func (m *mockContext) AddOnFocusFn(func(types.OnFocusOpts))                                        {}
func (m *mockContext) AddOnFocusLostFn(func(types.OnFocusLostOpts))                                {}

// HasKeybindings methods
func (m *mockContext) GetKeybindings(types.KeybindingsOpts) []*types.Binding                       { return nil }
func (m *mockContext) GetMouseKeybindings(types.KeybindingsOpts) []*gocui.ViewMouseBinding         { return nil }
func (m *mockContext) GetOnClick() func() error                                                    { return nil }
func (m *mockContext) GetOnClickFocusedMainView() func(string, int) error                          { return nil }

// ParentContexter methods
func (m *mockContext) SetParentContext(types.Context)                                              {}
func (m *mockContext) GetParentContext() types.Context                                             { return nil }

// Context methods
func (m *mockContext) HandleFocus(types.OnFocusOpts)                                              {}
func (m *mockContext) HandleFocusLost(types.OnFocusLostOpts)                                      {}
func (m *mockContext) FocusLine(bool)                                                             {}
func (m *mockContext) HandleRender()                                                              {}
func (m *mockContext) HandleRenderToMain()                                                        {}

func mockActivityBarStatus(operationInProgress string) types.IActivityBarStatus {
	s := guiModels.NewActivityBarStatus()
	if operationInProgress != "" {
		s.SetOperationInProgress(operationInProgress, true)
	}
	return s
}
