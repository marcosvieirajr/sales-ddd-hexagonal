package types

type Sex int

const (
	SexNotInformed Sex = iota
	SexMale
	SexFemale
	SexOther
)
