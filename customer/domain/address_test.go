package customer_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/customer/domain"
	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== Tests   ==================== //

func TestNewAddress(t *testing.T) {
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
		want *customer.Address
	}{
		{
			name: "should create a valid address",
			args: args{
				cep: "12345-678", street: "Street", number: "123",
				complement: "Complement", district: "District", city: "City",
				state: "BA", country: "Country",
			},
			want: kernel.Must(customer.NewAddress("12345-678", "Street", "123", "Complement", "District", "City", "BA", "Country")),
		},
		{
			name: "should create a valid address without complement",
			args: args{
				cep: "12345-678", street: "Street", number: "123",
				complement: "", district: "District", city: "City",
				state: "BA", country: "Country",
			},
			want: kernel.Must(customer.NewAddress("12345-678", "Street", "123", "", "District", "City", "BA", "Country")),
		},
	}
	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := customer.NewAddress(
				tt.args.cep, tt.args.street, tt.args.number, tt.args.complement,
				tt.args.district, tt.args.city, tt.args.state, tt.args.country,
			)

			require.NoError(t, err)
			assert.NotNil(t, got)
			assert.NotEmpty(t, got.ID)
			assert.False(t, got.CreatedAt.IsZero())
			assert.Nil(t, got.UpdatedAt)
			ignoreFields := cmpopts.IgnoreFields(customer.Address{}, "ID", "CreatedAt", "UpdatedAt")
			assert.True(t,
				cmp.Equal(got, tt.want, cmp.AllowUnexported(customer.Address{}), ignoreFields),
				cmp.Diff(got, tt.want, cmp.AllowUnexported(customer.Address{}), ignoreFields),
			)
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
			wantErr: customer.ErrInvalidStreet,
		},
		{
			name:    "should return an error when number is empty",
			args:    args{cep: "12345-678", street: "Street", number: "", complement: "Complement", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidNumber,
		},
		{
			name:    "should return an error when district is empty",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "Complement", district: "", city: "City", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidDistrict,
		},
		{
			name:    "should return an error when city is empty",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "Complement", district: "District", city: "", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidCity,
		},
		{
			name:    "should return an error when country is empty",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "Complement", district: "District", city: "City", state: "BA", country: ""},
			wantErr: customer.ErrInvalidCountry,
		},
		{
			name:    "should return an error when CEP is empty",
			args:    args{cep: "", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidCEP,
		},
		{
			name:    "should return an error when CEP is missing hyphen",
			args:    args{cep: "12345678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidCEP,
		},
		{
			name:    "should return an error when CEP has too many digits after hyphen",
			args:    args{cep: "12345-7890", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidCEP,
		},
		{
			name:    "should return an error when CEP has hyphen in wrong position",
			args:    args{cep: "12-345678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidCEP,
		},
		{
			name:    "should return an error when CEP has non-numeric characters",
			args:    args{cep: "ABCDE-123", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidCEP,
		},
		{
			name:    "should return an error when state is an invalid UF code",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "AA", country: "Country"},
			wantErr: customer.ErrInvalidState,
		},
		{
			name:    "should return an error when state is a full state name instead of UF",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "State", country: "Country"},
			wantErr: customer.ErrInvalidState,
		},
		{
			name:    "should return an error when state is a single character",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "A", country: "Country"},
			wantErr: customer.ErrInvalidState,
		},
	}
	for _, tt := range failureTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := customer.NewAddress(
				tt.args.cep, tt.args.street, tt.args.number, tt.args.complement,
				tt.args.district, tt.args.city, tt.args.state, tt.args.country,
			)

			require.Error(t, err)
			assert.Nil(t, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestAddress_Equals(t *testing.T) {
	baseAddr := kernel.Must(customer.NewAddress(
		"12345-678", "Street", "123", "",
		"District", "City", "BA", "Country",
	))

	tests := []struct {
		name  string
		other *customer.Address
		want  bool
	}{
		// ==================== Success cases ==================== //
		{
			name:  "should return true for equal addresses",
			other: baseAddr,
			want:  true,
		},
		// ==================== Failure cases ==================== //
		{
			name:  "should return false for different addresses",
			other: kernel.Must(customer.NewAddress("12345-678", "Street n2", "123", "", "District", "City", "BA", "Country")),
			want:  false,
		},
		{
			name:  "should return false for nil address",
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

func TestAddress_IsZero(t *testing.T) {
	tests := []struct {
		name string
		addr *customer.Address
		want bool
	}{
		{
			name: "should return true for nil pointer",
			addr: nil,
			want: true,
		},
		{
			name: "should return true for zero-value struct",
			addr: &customer.Address{},
			want: true,
		},
		{
			name: "should return false for a valid address",
			addr: kernel.Must(customer.NewAddress("12345-678", "Street", "123", "", "District", "City", "BA", "Country")),
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

func TestAddress_Update(t *testing.T) {
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
		want *customer.Address
	}{
		{
			name: "should update all fields",
			args: args{
				cep: "98765-432", street: "New Street", number: "456",
				complement: "New Complement", district: "New District", city: "New City",
				state: "SP", country: "New Country",
			},
			want: kernel.Must(customer.NewAddress("98765-432", "New Street", "456", "New Complement", "New District", "New City", "SP", "New Country")),
		},
		{
			name: "should update without complement",
			args: args{
				cep: "98765-432", street: "New Street", number: "456",
				complement: "", district: "New District", city: "New City",
				state: "SP", country: "New Country",
			},
			want: kernel.Must(customer.NewAddress("98765-432", "New Street", "456", "", "New District", "New City", "SP", "New Country")),
		},
	}
	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			addr := kernel.Must(customer.NewAddress(
				"12345-678", "Street", "123", "Complement", "District", "City", "BA", "Country",
			))
			beforeID := addr.ID
			beforeCreatedAt := addr.CreatedAt

			err := addr.Update(
				tt.args.cep, tt.args.street, tt.args.number, tt.args.complement,
				tt.args.district, tt.args.city, tt.args.state, tt.args.country,
			)

			require.NoError(t, err)
			assert.Equal(t, beforeID, addr.ID, "ID must not change on update")
			assert.Equal(t, beforeCreatedAt, addr.CreatedAt, "CreatedAt must not change on update")
			assert.NotNil(t, addr.UpdatedAt, "UpdatedAt must be set after update")
			ignoreFields := cmpopts.IgnoreFields(customer.Address{}, "ID", "CreatedAt", "UpdatedAt")
			assert.True(t,
				cmp.Equal(addr, tt.want, cmp.AllowUnexported(customer.Address{}), ignoreFields),
				cmp.Diff(addr, tt.want, cmp.AllowUnexported(customer.Address{}), ignoreFields),
			)
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
			wantErr: customer.ErrInvalidStreet,
		},
		{
			name:    "should return an error when number is empty",
			args:    args{cep: "12345-678", street: "Street", number: "", complement: "Complement", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidNumber,
		},
		{
			name:    "should return an error when district is empty",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "Complement", district: "", city: "City", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidDistrict,
		},
		{
			name:    "should return an error when city is empty",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "Complement", district: "District", city: "", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidCity,
		},
		{
			name:    "should return an error when country is empty",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "Complement", district: "District", city: "City", state: "BA", country: ""},
			wantErr: customer.ErrInvalidCountry,
		},
		{
			name:    "should return an error when CEP is empty",
			args:    args{cep: "", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidCEP,
		},
		{
			name:    "should return an error when CEP is missing hyphen",
			args:    args{cep: "12345678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "BA", country: "Country"},
			wantErr: customer.ErrInvalidCEP,
		},
		{
			name:    "should return an error when state is invalid",
			args:    args{cep: "12345-678", street: "Street", number: "123", complement: "", district: "District", city: "City", state: "AA", country: "Country"},
			wantErr: customer.ErrInvalidState,
		},
	}
	for _, tt := range failureTests {
		t.Run(tt.name, func(t *testing.T) {
			addr := kernel.Must(customer.NewAddress(
				"12345-678", "Street", "123", "Complement", "District", "City", "BA", "Country",
			))
			beforeID := addr.ID

			err := addr.Update(
				tt.args.cep, tt.args.street, tt.args.number, tt.args.complement,
				tt.args.district, tt.args.city, tt.args.state, tt.args.country,
			)

			require.Error(t, err)
			assert.Equal(t, beforeID, addr.ID, "ID must not change on failed update")
			assert.Nil(t, addr.UpdatedAt, "UpdatedAt must not be set on failed update")
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
