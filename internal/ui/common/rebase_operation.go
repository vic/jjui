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

func (r RebaseOperation) RendersAfter() bool {
	return false
}
