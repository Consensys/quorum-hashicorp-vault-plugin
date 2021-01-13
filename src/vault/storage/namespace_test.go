package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeKey(t *testing.T) {
	t.Run("should compute key by adding the namespace", func(t *testing.T) {
		key := computeStorageKey("ethereum", "address", "namespace")

		assert.Equal(t, "namespace/ethereum/accounts/address", key)
	})

	t.Run("should compute key without namespace", func(t *testing.T) {
		key := computeStorageKey("ethereum", "address", "")

		assert.Equal(t, "ethereum/accounts/address", key)
	})
}
