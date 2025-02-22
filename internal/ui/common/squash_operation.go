package common

type SquashOperation struct {
	From string
}

func (s SquashOperation) Render() string {
	return DropStyle.Render("<< into >>")
}

func (s SquashOperation) RenderPosition() RenderPosition {
	return RenderPositionGlyph
}

func NewSquashOperation(from string) SquashOperation {
	return SquashOperation{
		From: from,
	}
}
