package types

import "github.com/dswcpp/lazygit/pkg/commands/models"

type RefRange struct {
	From models.Ref
	To   models.Ref
}
