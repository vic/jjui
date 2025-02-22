package common

type RebaseSource int

const (
	RebaseSourceRevision RebaseSource = iota
	RebaseSourceBranch
)

type RebaseTarget int

const (
	RebaseTargetDestination RebaseTarget = iota
	RebaseTargetAfter
	RebaseTargetBefore
)

type RebaseOperation struct {
	From   string
	To     string
	Source RebaseSource
	Target RebaseTarget
}

func (r RebaseOperation) RenderPosition() RenderPosition {
	if r.Target == RebaseTargetAfter {
		return RenderPositionBefore
	}
	if r.Target == RebaseTargetDestination {
		return RenderPositionBefore
	}
	return RenderPositionAfter
}

func (r RebaseOperation) Render() string {
	if r.Target == RebaseTargetDestination {
		return DropStyle.Render("<< onto >>")
	}
	if r.Target == RebaseTargetAfter {
		return DropStyle.Render("<< after >>")
	}
	if r.Target == RebaseTargetBefore {
		return DropStyle.Render("<< before >>")
	}
	return ""
}
