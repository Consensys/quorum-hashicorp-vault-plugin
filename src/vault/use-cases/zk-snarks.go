package usecases

import (
	"context"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/hashicorp/vault/sdk/logical"
)

//go:generate mockgen -source=zk-snarks.go -destination=mocks/zk-snarks.go -package=mocks

type ZksUseCases interface {
	CreateAccount() CreateZksAccountUseCase
	GetAccount() GetZksAccountUseCase
	ListAccounts() ListZksAccountsUseCase
	ListNamespaces() ListZksNamespacesUseCase
	SignPayload() ZksSignUseCase
}

type CreateZksAccountUseCase interface {
	Execute(ctx context.Context, namespace string) (*entities.ZksAccount, error)
	WithStorage(storage logical.Storage) CreateZksAccountUseCase
}

type GetZksAccountUseCase interface {
	Execute(ctx context.Context, pubKey, namespace string) (*entities.ZksAccount, error)
	WithStorage(storage logical.Storage) GetZksAccountUseCase
}

type ListZksAccountsUseCase interface {
	Execute(ctx context.Context, namespace string) ([]string, error)
	WithStorage(storage logical.Storage) ListZksAccountsUseCase
}

type ZksSignUseCase interface {
	Execute(ctx context.Context, pubKey, namespace, data string) (string, error)
	WithStorage(storage logical.Storage) ZksSignUseCase
}

type ListZksNamespacesUseCase interface {
	Execute(ctx context.Context) ([]string, error)
	WithStorage(storage logical.Storage) ListZksNamespacesUseCase
}
