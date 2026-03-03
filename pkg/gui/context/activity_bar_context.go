package context

import (
	"github.com/dswcpp/lazygit/pkg/commands/models"
	guiModels "github.com/dswcpp/lazygit/pkg/gui/models"
	"github.com/dswcpp/lazygit/pkg/gui/presentation"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

type ActivityBarContext struct {
	*FilteredListViewModel[*models.ActivityBarItem]
	*ListContextTrait
}

var _ types.IListContext = (*ActivityBarContext)(nil)

func NewActivityBarContext(c *ContextCommon) *ActivityBarContext {
	viewModel := NewFilteredListViewModel(
		func() []*models.ActivityBarItem { return c.Model().ActivityBarItems },
		func(item *models.ActivityBarItem) []string {
			return []string{item.Name, item.Tooltip}
		},
	)

	getDisplayStrings := func(startIdx int, endIdx int) [][]string {
		// 获取 activityBarStatus（需要类型断言）
		var activityBarStatus *guiModels.ActivityBarStatus
		if statusInterface := c.GetActivityBarStatus(); statusInterface != nil {
			activityBarStatus, _ = statusInterface.(*guiModels.ActivityBarStatus)
		}

		return presentation.GetActivityBarDisplayStrings(
			viewModel.GetItems(),
			c.UserConfig().Gui.ActivityBar,
			c.Context().Current(),
			c.UserConfig(),
			activityBarStatus,
		)
	}

	self := &ActivityBarContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				View:       c.Views().ActivityBar,
				WindowName: "activityBar",
				Key:        ACTIVITY_BAR_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT, // Side context like branches
				Focusable:  true,
			})),
			ListRenderer: ListRenderer{
				list:              viewModel,
				getDisplayStrings: getDisplayStrings,
			},
			c: c,
		},
	}

	return self
}
