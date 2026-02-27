package domain_test

import (
	"errors"
	"testing"

	"github.com/marcosvieirajr/sales/domain"
	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	t.Run("should return value when no error is provided", func(t *testing.T) {
		got := domain.Must("expected value", nil)
		assert.Equal(t, "expected value", got)
	})

	t.Run("should panic when an error is provided", func(t *testing.T) {
		err := errors.New("some error")
		assert.Panics(t, func() {
			domain.Must("value", err)
		})
	})
}

func TestGenerateID(t *testing.T) {
	id := domain.GenerateID()
	assert.NotEmpty(t, id)
}
