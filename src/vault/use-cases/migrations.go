package usecases

import (
	"context"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/entities"
	"github.com/hashicorp/vault/sdk/logical"
)

//go:generate mockgen -source=migrations.go -destination=mocks/migrations.go -package=mocks

type MigrationsUseCases interface {
	EthereumToKeys() EthereumToKeysUseCase
}

type EthereumToKeysUseCase interface {
	Execute(ctx context.Context, namespace string) error
	WithStorage(storage logical.Storage) EthereumToKeysUseCase
	Status(ctx context.Context, namespace string) (*entities.MigrationStatus, error)
}
