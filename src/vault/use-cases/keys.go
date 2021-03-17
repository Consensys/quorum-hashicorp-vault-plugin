package usecases

import (
	"context"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/hashicorp/vault/sdk/logical"
)

//go:generate mockgen -source=keys.go -destination=mocks/keys.go -package=mocks

type KeysUseCases interface {
	CreateKey() CreateKeyUseCase
	GetKey() GetKeyUseCase
	ListKeys() ListKeysUseCase
	ListNamespaces() ListKeysNamespacesUseCase
	SignPayload() KeysSignUseCase
}

type CreateKeyUseCase interface {
	Execute(ctx context.Context, namespace, id, algo, curve, importedPrivKey string, tags map[string]string) (*entities.Key, error)
	WithStorage(storage logical.Storage) CreateKeyUseCase
}

type GetKeyUseCase interface {
	Execute(ctx context.Context, id, namespace string) (*entities.Key, error)
	WithStorage(storage logical.Storage) GetKeyUseCase
}

type ListKeysUseCase interface {
	Execute(ctx context.Context, namespace string) ([]string, error)
	WithStorage(storage logical.Storage) ListKeysUseCase
}

type KeysSignUseCase interface {
	Execute(ctx context.Context, id, namespace, data string) (string, error)
	WithStorage(storage logical.Storage) KeysSignUseCase
}

type ListKeysNamespacesUseCase interface {
	Execute(ctx context.Context) ([]string, error)
	WithStorage(storage logical.Storage) ListKeysNamespacesUseCase
}
