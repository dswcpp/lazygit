package context

import (
	"github.com/dswcpp/lazygit/pkg/common"
	"github.com/dswcpp/lazygit/pkg/gui/types"
)

type ContextCommon struct {
	*common.Common
	types.IGuiCommon
}
