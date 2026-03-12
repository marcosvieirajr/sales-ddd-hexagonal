package customer

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/errs"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel/guard"
)

var (
	ErrInvalidCEP      = errs.New("ADDRESS.INVALID_CEP_FORMAT", "invalid CEP: must be in the format 12345-678")
	ErrInvalidStreet   = errs.New("ADDRESS.INVALID_STREET", "street cannot be null or whitespace")
	ErrInvalidNumber   = errs.New("ADDRESS.INVALID_NUMBER", "number cannot be null or whitespace")
	ErrInvalidDistrict = errs.New("ADDRESS.INVALID_DISTRICT", "district cannot be null or whitespace")
	ErrInvalidCity     = errs.New("ADDRESS.INVALID_CITY", "city cannot be null or whitespace")
	ErrInvalidState    = errs.New("ADDRESS.INVALID_STATE", "invalid state: must be a valid Brazilian state (UF)")
	ErrInvalidCountry  = errs.New("ADDRESS.INVALID_COUNTRY", "country cannot be null or whitespace")
)

// Address is an entity of the Customer aggregate that representing a Brazilian postal address.
// All fields are unexported to enforce construction through [NewAddress] and
// to prevent external mutation.
type Address struct {
	ID         string
	cep        string
	street     string
	number     string
	complement string // optional
	district   string
	city       string
	state      string
	country    string
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

// NewAddress constructs and validates a [Address] Entity.
// All fields except complement are required (non-empty, non-whitespace).
// cep must follow the Brazilian postal format "12345-678" and state must be a valid
// two-letter UF code (e.g. "SP", "RJ"). complement may be an empty string.
//
// If multiple fields are invalid, all violations are collected and returned as a
// single joined error, allowing callers to inspect every failure via [errors.Is].
func NewAddress(cep, street, number, complement, district, city, state, country string) (*Address, error) {
	if err := errors.Join(
		guard.CheckNotNullOrWhiteSpace(street, ErrInvalidStreet),
		guard.CheckNotNullOrWhiteSpace(number, ErrInvalidNumber),
		guard.CheckNotNullOrWhiteSpace(district, ErrInvalidDistrict),
		guard.CheckNotNullOrWhiteSpace(city, ErrInvalidCity),
		guard.CheckNotNullOrWhiteSpace(country, ErrInvalidCountry),
		guard.CheckMatchRegex(cep, cepRegex, ErrInvalidCEP),
		checkValidState(state),
	); err != nil {
		return nil, err
	}

	return &Address{
		ID:         kernel.NewID().String(),
		cep:        cep,
		street:     street,
		number:     number,
		complement: complement,
		district:   district,
		city:       city,
		state:      state,
		country:    country,
		CreatedAt:  time.Now().UTC(),
	}, nil
}

// Update validates the given fields using the same rules as [NewAddress] and,
// if all validations pass, replaces the receiver's fields in-place.
// Returns a joined error for all violations, leaving the receiver unchanged on failure.
func (a *Address) Update(cep, street, number, complement, district, city, state, country string) error {
	updated, err := NewAddress(cep, street, number, complement, district, city, state, country)
	if err != nil {
		return err
	}

	updated.ID = a.ID
	updated.CreatedAt = a.CreatedAt

	*a = *updated
	a.updateTimestamp()
	return nil
}

// Equals reports whether a and other are the same Address entity by comparing IDs.
// It returns false if other is nil.
func (a *Address) Equals(other *Address) bool {
	if other == nil {
		return false
	}
	return a.ID == other.ID
}

// IsZero reports whether the Address is uninitialized (nil pointer or zero-value struct).
func (a *Address) IsZero() bool {
	return a == nil || *a == Address{}
}

func (a *Address) updateTimestamp() {
	a.UpdatedAt = new(time.Now().UTC())
}

func checkValidState(state string) error {
	if _, ok := validStates[strings.ToUpper(state)]; !ok {
		return ErrInvalidState
	}
	return nil
}

// Regular expression for validating Brazilian CEP format (12345-6785: 5 digits, a hyphen, and 3 digits)
// Note: The regex is a package-level precompiled variable to avoid recompiling it on every validation of a DeliveryAddress.
var cepRegex = regexp.MustCompile(`^\d{5}-\d{3}$`)

// List of valid Brazilian states (UF) for validation. Using a map for O(1) lookups.
// Note: This is a package-level variable to avoid recreating the map on every validation of a DeliveryAddress.
var validStates = map[string]struct{}{
	"AC": {}, "AL": {}, "AP": {}, "AM": {}, "BA": {}, "CE": {}, "DF": {}, "ES": {},
	"GO": {}, "MA": {}, "MT": {}, "MS": {}, "MG": {}, "PA": {}, "PB": {}, "PR": {},
	"PE": {}, "PI": {}, "RJ": {}, "RN": {}, "RS": {}, "RO": {}, "RR": {}, "SC": {},
	"SP": {}, "SE": {}, "TO": {},
}
