package membership

type Status int

const (
	Alive Status = iota
	Suspect
	Dead
)
