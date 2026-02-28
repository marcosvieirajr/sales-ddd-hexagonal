package kernel_test

import (
	"errors"
	"testing"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/kernel"
	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	t.Run("should return value when no error is provided", func(t *testing.T) {
		got := kernel.Must("expected value", nil)
		assert.Equal(t, "expected value", got)
	})

	t.Run("should panic when an error is provided", func(t *testing.T) {
		err := errors.New("some error")
		assert.Panics(t, func() {
			kernel.Must("value", err)
		})
	})
}

func TestGenerateID(t *testing.T) {
	id := kernel.GenerateID()
	assert.NotEmpty(t, id)
}
