package operations

type Noop struct{}

func (n *Noop) RenderPosition() RenderPosition {
	return RenderPositionNil
}

func (n *Noop) Render() string {
	return ""
}
