package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeKey(t *testing.T) {
	t.Run("should compute key by adding the namespace", func(t *testing.T) {
		key := ComputeKey("address", "namespace")

		assert.Equal(t, "namespace/ethereum/accounts/address", key)
	})

	t.Run("should compute key without namespace", func(t *testing.T) {
		key := ComputeKey("address", "")

		assert.Equal(t, "ethereum/accounts/address", key)
	})
}
