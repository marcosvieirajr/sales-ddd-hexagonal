package errs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/shared/errs"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := errs.New("TEST.CODE", "test message")

	assert.NotNil(t, err)
	assert.Equal(t, errs.ErrorCode("TEST.CODE"), err.Code)
	assert.Equal(t, "test message", err.Message)
	assert.Nil(t, err.Err)
}

func TestWrap(t *testing.T) {
	underlying := fmt.Errorf("underlying cause")

	err := errs.Wrap("TEST.CODE", "test message", underlying)

	assert.NotNil(t, err)
	assert.Equal(t, errs.ErrorCode("TEST.CODE"), err.Code)
	assert.Equal(t, "test message", err.Message)
	assert.Equal(t, underlying, err.Err)
}

func TestDomainError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *errs.DomainError
		want string
	}{
		{
			name: "should format error without underlying error",
			err:  errs.New("TEST.CODE", "test message"),
			want: "[TEST.CODE] test message",
		},
		{
			name: "should format error with underlying error",
			err:  errs.Wrap("TEST.CODE", "test message", fmt.Errorf("underlying cause")),
			want: "[TEST.CODE] test message: underlying cause",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Error())
		})
	}
}

func TestDomainError_Unwrap(t *testing.T) {
	underlying := fmt.Errorf("underlying cause")

	err := errs.Wrap("TEST.CODE", "test message", underlying)

	assert.Equal(t, underlying, errors.Unwrap(err))
}

func TestDomainError_Is(t *testing.T) {
	sentinel := errs.New("TEST.CODE", "test message")

	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{
			name:   "should match errors with the same code",
			err:    sentinel,
			target: errs.New("TEST.CODE", "different message"),
			want:   true,
		},
		{
			name:   "should match wrapped copy against sentinel",
			err:    sentinel.Wrap(fmt.Errorf("some cause")),
			target: sentinel,
			want:   true,
		},
		{
			name:   "should not match errors with different codes",
			err:    sentinel,
			target: errs.New("OTHER.CODE", "test message"),
			want:   false,
		},
		{
			name:   "should not match non-DomainError target",
			err:    sentinel,
			target: fmt.Errorf("plain error"),
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, errors.Is(tt.err, tt.target))
		})
	}
}

func TestDomainError_Wrap(t *testing.T) {
	sentinel := errs.New("TEST.CODE", "test message")
	underlying := fmt.Errorf("underlying cause")

	wrapped := sentinel.Wrap(underlying)

	assert.Equal(t, sentinel.Code, wrapped.Code)
	assert.Equal(t, sentinel.Message, wrapped.Message)
	assert.Equal(t, underlying, wrapped.Err)
	assert.Nil(t, sentinel.Err, "original sentinel should not be modified")
}
