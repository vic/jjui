package common

type ShowDetailsOperation struct{}

func (s *ShowDetailsOperation) Render() string {
	return ""
}

func (s *ShowDetailsOperation) RenderPosition() RenderPosition {
	return RenderPositionAfter
}
