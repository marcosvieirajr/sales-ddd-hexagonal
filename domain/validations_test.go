package domain_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/marcosvieirajr/sales/domain"
	"github.com/stretchr/testify/assert"
)

var sentinelErr = fmt.Errorf("sentinel error")

func TestCheckMatchRegex(t *testing.T) {
	digitRegex := regexp.MustCompile(`^\d+$`)

	tests := []struct {
		name        string
		value       string
		wantErr error
	}{
		{
			name:        "should return nil when value matches regex",
			value:       "12345",
			wantErr: nil,
		},
		{
			name:        "should return error when value does not match regex",
			value:       "abc",
			wantErr: sentinelErr,
		},
		{
			name:        "should return error when value is empty",
			value:       "",
			wantErr: sentinelErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.CheckMatchRegex(tt.value, digitRegex, sentinelErr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCheckNotNullOrWhiteSpace(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		wantErr error
	}{
		{
			name:        "should return nil when value is non-empty",
			value:       "valid string",
			wantErr: nil,
		},
		{
			name:        "should return error when value is empty",
			value:       "",
			wantErr: sentinelErr,
		},
		{
			name:        "should return error when value contains only spaces",
			value:       "   ",
			wantErr: sentinelErr,
		},
		{
			name:        "should return error when value contains only a tab",
			value:       "\t",
			wantErr: sentinelErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.CheckNotNullOrWhiteSpace(tt.value, sentinelErr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCheckNotZeroOrNegative(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		wantErr error
	}{
		{
			name:        "should return nil when value is positive",
			value:       1.0,
			wantErr: nil,
		},
		{
			name:        "should return nil when value is a very small positive number",
			value:       0.001,
			wantErr: nil,
		},
		{
			name:        "should return error when value is zero",
			value:       0.0,
			wantErr: sentinelErr,
		},
		{
			name:        "should return error when value is negative",
			value:       -1.0,
			wantErr: sentinelErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.CheckNotZeroOrNegative(tt.value, sentinelErr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCheckNotNil(t *testing.T) {
	nonNilValue := "value"
	var typedNilPtr *string

	tests := []struct {
		name        string
		value       any
		wantErr error
	}{
		{
			name:        "should return nil when value is non-nil",
			value:       &nonNilValue,
			wantErr: nil,
		},
		{
			name:        "should return error when value is untyped nil",
			value:       nil,
			wantErr: sentinelErr,
		},
		{
			name:        "should return error when value is a typed nil pointer",
			value:       typedNilPtr,
			wantErr: sentinelErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.CheckNotNil(tt.value, sentinelErr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCheckNil(t *testing.T) {
	nonNilValue := "value"
	var typedNilPtr *string

	tests := []struct {
		name        string
		value       any
		wantErr error
	}{
		{
			name:        "should return nil when value is untyped nil",
			value:       nil,
			wantErr: nil,
		},
		{
			name:        "should return nil when value is a typed nil pointer",
			value:       typedNilPtr,
			wantErr: nil,
		},
		{
			name:        "should return error when value is non-nil",
			value:       &nonNilValue,
			wantErr: sentinelErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.CheckNil(tt.value, sentinelErr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
