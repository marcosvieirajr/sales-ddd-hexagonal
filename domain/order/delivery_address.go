package order

import (
	"errors"
	"regexp"
	"strings"

	"github.com/marcosvieirajr/sales/domain"
	"github.com/marcosvieirajr/sales/domain/errs"
)

var (
	ErrInvalidCEP      = errs.New("DELIVERY_ADDRESS.INVALID_CEP_FORMAT", "invalid CEP: must be in the format 12345-678")
	ErrInvalidStreet   = errs.New("DELIVERY_ADDRESS.INVALID_STREET", "street cannot be null or whitespace")
	ErrInvalidNumber   = errs.New("DELIVERY_ADDRESS.INVALID_NUMBER", "number cannot be null or whitespace")
	ErrInvalidDistrict = errs.New("DELIVERY_ADDRESS.INVALID_DISTRICT", "district cannot be null or whitespace")
	ErrInvalidCity     = errs.New("DELIVERY_ADDRESS.INVALID_CITY", "city cannot be null or whitespace")
	ErrInvalidState    = errs.New("DELIVERY_ADDRESS.INVALID_STATE", "invalid state: must be a valid Brazilian state (UF)")
	ErrInvalidCountry  = errs.New("DELIVERY_ADDRESS.INVALID_COUNTRY", "country cannot be null or whitespace")
)

// DeliveryAddress is an immutable value object representing a Brazilian postal address.
// All fields are unexported to enforce construction through [NewDeliveryAddress] and
// to prevent external mutation. Two DeliveryAddress values are equal when every field
// is equal (see [DeliveryAddress.Equals]).
type DeliveryAddress struct {
	cep        string
	street     string
	number     string
	complement string
	district   string
	city       string
	state      string
	country    string
}

// NewDeliveryAddress constructs and validates a [DeliveryAddress] value object.
// All fields except complement are required (non-empty, non-whitespace).
// cep must follow the Brazilian postal format "12345-678" and state must be a valid
// two-letter UF code (e.g. "SP", "RJ"). complement may be an empty string.
//
// If multiple fields are invalid, all violations are collected and returned as a
// single joined error, allowing callers to inspect every failure via [errors.Is].
func NewDeliveryAddress(cep, street, number, complement, district, city, state, country string) (*DeliveryAddress, error) {
	if err := errors.Join(
		domain.CheckNotNullOrWhiteSpace(street, ErrInvalidStreet),
		domain.CheckNotNullOrWhiteSpace(number, ErrInvalidNumber),
		domain.CheckNotNullOrWhiteSpace(district, ErrInvalidDistrict),
		domain.CheckNotNullOrWhiteSpace(city, ErrInvalidCity),
		domain.CheckNotNullOrWhiteSpace(country, ErrInvalidCountry),
		domain.CheckMatchRegex(cep, cepRegex, ErrInvalidCEP),
		checkValidState(state),
	); err != nil {
		return nil, err
	}

	da := DeliveryAddress{
		cep:        cep,
		street:     street,
		number:     number,
		complement: complement,
		district:   district,
		city:       city,
		state:      state,
		country:    country,
	}

	return &da, nil
}

// Equals reports whether da and other represent the same postal address by
// comparing every field for equality. It returns false if other is nil.
func (da *DeliveryAddress) Equals(other *DeliveryAddress) bool {
	if other == nil {
		return false
	}
	return *da == *other
}

func checkValidState(state string) error {
	state = strings.ToUpper(state)
	if _, ok := validStates[state]; !ok {
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
