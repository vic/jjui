package common

type None struct{}

func (n *None) RenderPosition() RenderPosition {
	return RenderPositionNil
}

func (n *None) Render() string {
	return ""
}
