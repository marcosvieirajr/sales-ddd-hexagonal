package types

import "github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/errs"

var ErrInvalidSex = errs.New("SEX.INVALID", "invalid sex")

// Sex represents the biological sex of a person.
type Sex struct{ value int }

var (
	SexNotInformed = Sex{0} // SexNotInformed is the zero value, used when sex is not provided.
	SexMale        = Sex{1}
	SexFemale      = Sex{2}
	SexOther       = Sex{3}
)

var sexToString = map[Sex]string{
	SexNotInformed: "not_informed",
	SexMale:        "male",
	SexFemale:      "female",
	SexOther:       "other",
}

// String returns the string representation of the Sex.
func (s Sex) String() string {
	if str, ok := sexToString[s]; ok {
		return str
	}
	return "unknown"
}

// MarshalText provides support for logging and any marshal needs.
func (s Sex) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// Equals checks if two Sex values are equal.
func (s Sex) Equals(other Sex) bool {
	return s.value == other.value
}

// ParseSex converts an int to the corresponding Sex value.
// If the input does not match any known value, it returns an error and an empty Sex value.
func ParseSex(value int) (Sex, error) {
	s := Sex{value}
	if _, ok := sexToString[s]; !ok {
		return Sex{}, ErrInvalidSex
	}
	return s, nil
}
