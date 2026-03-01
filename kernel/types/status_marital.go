package types

import "github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/errs"

var ErrInvalidMaritalStatus = errs.New("MARITAL_STATUS.INVALID", "invalid marital status")

// MaritalStatus represents the marital status of a person.
type MaritalStatus struct{ value int }

var (
	MaritalStatusNotInformed = MaritalStatus{0} // MaritalStatusNotInformed is the zero value, used when marital status is not provided.
	MaritalStatusSingle      = MaritalStatus{1}
	MaritalStatusMarried     = MaritalStatus{2}
	MaritalStatusDivorced    = MaritalStatus{3}
	MaritalStatusWidowed     = MaritalStatus{4}
	MaritalStatusStableUnion = MaritalStatus{5}
)

var maritalStatusToString = map[MaritalStatus]string{
	MaritalStatusNotInformed: "not_informed",
	MaritalStatusSingle:      "single",
	MaritalStatusMarried:     "married",
	MaritalStatusDivorced:    "divorced",
	MaritalStatusWidowed:     "widowed",
	MaritalStatusStableUnion: "stable_union",
}

// String returns the string representation of the MaritalStatus.
func (m MaritalStatus) String() string {
	if str, ok := maritalStatusToString[m]; ok {
		return str
	}
	return "unknown"
}

// MarshalText provides support for logging and any marshal needs.
func (m MaritalStatus) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// Equals checks if two MaritalStatus values are equal.
func (m MaritalStatus) Equals(other MaritalStatus) bool {
	return m.value == other.value
}

// ParseMaritalStatus converts an int to the corresponding MaritalStatus value.
// If the input does not match any known value, it returns an error and an empty MaritalStatus value.
func ParseMaritalStatus(value int) (MaritalStatus, error) {
	ms := MaritalStatus{value}
	if _, ok := maritalStatusToString[ms]; !ok {
		return MaritalStatus{}, ErrInvalidMaritalStatus
	}
	return ms, nil
}
