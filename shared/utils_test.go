package shared_test

import (
	"errors"
	"testing"

	"github.com/marcosvieirajr/sales-ddd-hexagonal/shared"
	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	t.Run("should return value when no error is provided", func(t *testing.T) {
		got := shared.Must("expected value", nil)
		assert.Equal(t, "expected value", got)
	})

	t.Run("should panic when an error is provided", func(t *testing.T) {
		err := errors.New("some error")
		assert.Panics(t, func() {
			shared.Must("value", err)
		})
	})
}

func TestGenerateID(t *testing.T) {
	id := shared.GenerateID()
	assert.NotEmpty(t, id)
}
