package types

type MaritalStatus int

const (
	MaritalStatusNotInformed MaritalStatus = iota
	MaritalStatusSingle
	MaritalStatusMarried
	MaritalStatusDivorced
	MaritalStatusWidowed
	MaritalStatusStableUnion
)
