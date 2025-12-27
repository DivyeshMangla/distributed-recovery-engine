package node

type Status int

const (
	Alive Status = iota
	Suspect
	Dead
)
