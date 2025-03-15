package jj

type GraphRow struct {
	Connections [][]ConnectionType
	Commit      *Commit
	IsSelected  bool
}
