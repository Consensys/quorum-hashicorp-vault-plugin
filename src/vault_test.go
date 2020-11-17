package src

import (
	"context"
	"github.com/hashicorp/vault/sdk/logical"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVaultBackend(t *testing.T) {
	backend, err := NewVaultBackend(context.Background(), &logical.BackendConfig{})

	assert.NoError(t, err)
	assert.NotEmpty(t, backend)
}
