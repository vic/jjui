package squash

import (
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	From string
}

func (s Operation) Render() string {
	return common.DropStyle.Render("<< into >>")
}

func (s Operation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionGlyph
}

func NewOperation(from string) Operation {
	return Operation{
		From: from,
	}
}
