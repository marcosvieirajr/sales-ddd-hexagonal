package order_test

import (
	"reflect"
	"testing"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== Tests   ==================== //

func TestNewDeliveryAddress(t *testing.T) {
	type args struct {
		cep        string
		street     string
		number     string
		complement string
		district   string
		city       string
		state      string
		country    string
	}

	// ==================== Success cases ==================== //
	successTests := []struct {
		name string
		args args
		want *order.DeliveryAddress
	}{
		{
			name: "should create a valid address",
			args: args{
				cep: "12345-678", street: "Street", number: "123",
				complement: "Complement", district: "District", city: "City",
				state: "BA", country: "Country",
			},
			want: kernel.Must(order.NewDeliveryAddress(
				"12345-678", "Street", "123", "Complement", "District", "City", "BA", "Country",
			)),
		},
		{
			name: "should create a valid address without complement",
			args: args{
				cep: "12345-678", street: "Street", number: "123",
				complement: "", district: "District", city: "City",
				state: "BA", country: "Country",
			},
			want: kernel.Must(order.NewDeliveryAddress(
				"12345-678", "Street", "123", "", "District", "City", "BA", "Country",
			)),
		},
	}
	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := order.NewDeliveryAddress(
				tt.args.cep, tt.args.street, tt.args.number, tt.args.complement,
				tt.args.district, tt.args.city, tt.args.state, tt.args.country,
			)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

	// ==================== Failure cases ==================== //
	failureTests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "should return an error when street is empty",
			args:    args{cep: "12345-678", street: "", number: "123", complement: "Complement", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: order.ErrInvalidStreet,
		},
		{
			name:    "should return an error when number is empty",
			args:    args{cep: "12345-678", street: "Street", number: "", complement: "Complement", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: order.ErrInvalidNumber,
		},
		{
			name:    "should return an error when district is empty",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "Complement", district: "", city: "City", state: "BA", country: "Country"},
			wantErr: order.ErrInvalidDistrict,
		},
		{
			name:    "should return an error when city is empty",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "Complement", district: "District", city: "", state: "BA", country: "Country"},
			wantErr: order.ErrInvalidCity,
		},
		{
			name:    "should return an error when country is empty",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "Complement", district: "District", city: "City", state: "BA", country: ""},
			wantErr: order.ErrInvalidCountry,
		},
		{
			name:    "should return an error when CEP is empty",
			args:    args{cep: "", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: order.ErrInvalidCEP,
		},
		{
			name:    "should return an error when CEP is missing hyphen",
			args:    args{cep: "12345678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: order.ErrInvalidCEP,
		},
		{
			name:    "should return an error when CEP has too many digits after hyphen",
			args:    args{cep: "12345-7890", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: order.ErrInvalidCEP,
		},
		{
			name:    "should return an error when CEP has hyphen in wrong position",
			args:    args{cep: "12-345678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: order.ErrInvalidCEP,
		},
		{
			name:    "should return an error when CEP has non-numeric characters",
			args:    args{cep: "ABCDE-123", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: order.ErrInvalidCEP,
		},
		{
			name:    "should return an error when state is an invalid UF code",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "AA", country: "Country"},
			wantErr: order.ErrInvalidState,
		},
		{
			name:    "should return an error when state is a full state name instead of UF",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "State", country: "Country"},
			wantErr: order.ErrInvalidState,
		},
		{
			name:    "should return an error when state is a single character",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "A", country: "Country"},
			wantErr: order.ErrInvalidState,
		},
	}
	for _, tt := range failureTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := order.NewDeliveryAddress(
				tt.args.cep, tt.args.street, tt.args.number, tt.args.complement,
				tt.args.district, tt.args.city, tt.args.state, tt.args.country,
			)

			require.Error(t, err)
			assert.Nil(t, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeliveryAddress_Equals(t *testing.T) {
	baseAddr := kernel.Must(order.NewDeliveryAddress(
		"12345-678", "Street", "123", "",
		"District", "City", "BA", "Country",
	))

	tests := []struct {
		name  string
		other *order.DeliveryAddress
		want  bool
	}{
		// ==================== Success cases ==================== //
		{
			name:  "should return true for equal delivery addresses",
			other: kernel.Must(order.NewDeliveryAddress("12345-678", "Street", "123", "", "District", "City", "BA", "Country")),
			want:  true,
		},
		// ==================== Failure cases ==================== //
		{
			name:  "should return false for different delivery addresses",
			other: kernel.Must(order.NewDeliveryAddress("12345-678", "Street n2", "123", "", "District", "City", "BA", "Country")),
			want:  false,
		},
		{
			name:  "should return false for nil delivery address",
			other: nil,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := baseAddr.Equals(tt.other)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeliveryAddress_IsZero(t *testing.T) {
	tests := []struct {
		name string
		addr *order.DeliveryAddress
		want bool
	}{
		{
			name: "should return true for nil pointer",
			addr: nil,
			want: true,
		},
		{
			name: "should return true for zero-value struct",
			addr: &order.DeliveryAddress{},
			want: true,
		},
		{
			name: "should return false for a valid address",
			addr: kernel.Must(order.NewDeliveryAddress("12345-678", "Street", "123", "", "District", "City", "BA", "Country")),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.addr.IsZero()

			assert.Equal(t, tt.want, got)
		})
	}
}

// This test ensures that all fields of the DeliveryAddress struct, as value object,
// are unexported, preventing external mutation after construction.
func TestDeliveryAddress_MustBeImmutable(t *testing.T) {
	typ := reflect.TypeOf(order.DeliveryAddress{})
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.PkgPath == "" {
			t.Fatalf("field %q is exported", f.Name)
		}
	}
}
