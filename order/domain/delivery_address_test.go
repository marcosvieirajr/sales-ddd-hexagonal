package order_test

import (
	"reflect"
	"testing"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/order/domain"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDeliveryAddress(t *testing.T) {
	t.Run("should successfully create a delivery address with valid input", func(t *testing.T) {
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
		tests := []struct {
			name string
			args args
			want *order.DeliveryAddress
		}{
			{
				name: "should create a valid delivery address",
				args: args{
					cep:        "12345-678",
					street:     "Street",
					number:     "123",
					complement: "Complement",
					district:   "District",
					city:       "City",
					state:      "BA",
					country:    "Country",
				},
				want: shared.Must(order.NewDeliveryAddress(
					"12345-678",
					"Street",
					"123",
					"Complement",
					"District",
					"City",
					"BA",
					"Country",
				)),
			},
			{
				name: "should create a valid delivery address without complement",
				args: args{
					cep:        "12345-678",
					street:     "Street",
					number:     "123",
					complement: "",
					district:   "District",
					city:       "City",
					state:      "BA",
					country:    "Country",
				},
				want: shared.Must(order.NewDeliveryAddress(
					"12345-678",
					"Street",
					"123",
					"",
					"District",
					"City",
					"BA",
					"Country",
				)),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := order.NewDeliveryAddress(tt.args.cep, tt.args.street, tt.args.number, tt.args.complement, tt.args.district, tt.args.city, tt.args.state, tt.args.country)
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("should return an error when required fields are empty", func(t *testing.T) {
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
		tests := []struct {
			name    string
			args    args
			wantErr error
		}{
			{
				name: "should return an error for delivery address with whitespace fields - street",
				args: args{
					cep:        "12345-678",
					street:     "",
					number:     "123",
					complement: "Complement",
					district:   "District",
					city:       "City",
					state:      "BA",
					country:    "Country",
				},
				wantErr: order.ErrInvalidStreet,
			}, {
				name: "should return an error for delivery address with whitespace fields - number",
				args: args{
					cep:        "12345-678",
					street:     "Street",
					number:     "",
					complement: "Complement",
					district:   "District",
					city:       "City",
					state:      "BA",
					country:    "Country",
				},
				wantErr: order.ErrInvalidNumber,
			},
			{
				name: "should return an error for delivery address with whitespace fields - district",
				args: args{
					cep:        "12345-678",
					street:     "Street",
					number:     "123",
					complement: "Complement",
					district:   "",
					city:       "City",
					state:      "BA",
					country:    "Country",
				},
				wantErr: order.ErrInvalidDistrict,
			},
			{
				name: "should return an error for delivery address with whitespace fields - city",
				args: args{
					cep:        "12345-678",
					street:     "Street",
					number:     "123",
					complement: "Complement",
					district:   "District",
					city:       "",
					state:      "BA",
					country:    "Country",
				},
				wantErr: order.ErrInvalidCity,
			},
			{
				name: "should return an error for delivery address with whitespace fields - country",
				args: args{
					cep:        "12345-678",
					street:     "Street",
					number:     "123",
					complement: "Complement",
					district:   "District",
					city:       "City",
					state:      "BA",
					country:    "",
				},
				wantErr: order.ErrInvalidCountry,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := order.NewDeliveryAddress(tt.args.cep, tt.args.street, tt.args.number, tt.args.complement, tt.args.district, tt.args.city, tt.args.state, tt.args.country)
				require.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})

	t.Run("should return an error when CEP is invalid", func(t *testing.T) {
		tests := []struct {
			name    string
			cep     string
			wantErr error
		}{
			{name: "whitespace", cep: "", wantErr: order.ErrInvalidCEP},
			{name: "missing hyphen", cep: "12345678", wantErr: order.ErrInvalidCEP},
			{name: "too many digits after hyphen", cep: "123456-789", wantErr: order.ErrInvalidCEP},
			{name: "hyphen in wrong position", cep: "12-345678", wantErr: order.ErrInvalidCEP},
			{name: "non-numeric characters", cep: "ABCDE-123", wantErr: order.ErrInvalidCEP},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := order.NewDeliveryAddress(tt.cep, "Street", "123", "",
					"District", "City", "BA", "Country")
				require.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})

	t.Run("should return an error when state is invalid", func(t *testing.T) {
		tests := []struct {
			name    string
			state   string
			wantErr error
		}{
			{name: "invalid UF code", state: "AA", wantErr: order.ErrInvalidState},
			{name: "full state name instead of UF", state: "State", wantErr: order.ErrInvalidState},
			{name: "single character", state: "A", wantErr: order.ErrInvalidState},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := order.NewDeliveryAddress("12345-678", "Street", "123", "",
					"District", "City", tt.state, "Country")
				require.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}

func TestDeliveryAddress_Equals(t *testing.T) {
	baseAddr := shared.Must(order.NewDeliveryAddress(
		"12345-678", "Street", "123", "",
		"District", "City", "BA", "Country",
	))

	tests := []struct {
		name  string
		other *order.DeliveryAddress
		want  bool
	}{
		{
			name:  "should return true for equal delivery addresses",
			other: shared.Must(order.NewDeliveryAddress("12345-678", "Street", "123", "", "District", "City", "BA", "Country")),
			want:  true,
		},
		{
			name:  "should return false for different delivery addresses",
			other: shared.Must(order.NewDeliveryAddress("12345-678", "Street n2", "123", "", "District", "City", "BA", "Country")),
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

// This test ensures that all fields of the DeliveryAddress struct, as value object,
// are unexported, which is a common way to enforce immutability. By making the fields unexported,
// we prevent external code from modifying the state of the DeliveryAddress after it has been created.
func TestDeliveryAddress_MustBeImmutable(t *testing.T) {
	typ := reflect.TypeOf(order.DeliveryAddress{})
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.PkgPath == "" {
			t.Fatalf("field %q is exported", f.Name)
		}
	}
}
