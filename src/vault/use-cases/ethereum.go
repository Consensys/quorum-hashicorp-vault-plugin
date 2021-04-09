package usecases

import (
	"context"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/hashicorp/vault/sdk/logical"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
)

//go:generate mockgen -source=ethereum.go -destination=mocks/ethereum.go -package=mocks

type ETHUseCases interface {
	CreateAccount() CreateAccountUseCase
	GetAccount() GetAccountUseCase
	ListAccounts() ListAccountsUseCase
	ListNamespaces() ListNamespacesUseCase
	SignPayload() SignUseCase
	SignTransaction() SignTransactionUseCase
	SignQuorumPrivateTransaction() SignQuorumPrivateTransactionUseCase
	SignEEATransaction() SignEEATransactionUseCase
}

type CreateAccountUseCase interface {
	Execute(ctx context.Context, namespace, importedPrivKey string) (*entities.ETHAccount, error)
	WithStorage(storage logical.Storage) CreateAccountUseCase
}

type GetAccountUseCase interface {
	Execute(ctx context.Context, address, namespace string) (*entities.ETHAccount, error)
	WithStorage(storage logical.Storage) GetAccountUseCase
}

type ListAccountsUseCase interface {
	Execute(ctx context.Context, namespace string) ([]string, error)
	WithStorage(storage logical.Storage) ListAccountsUseCase
}

type SignUseCase interface {
	Execute(ctx context.Context, address, namespace, data string) (string, error)
	WithStorage(storage logical.Storage) SignUseCase
}

type SignTransactionUseCase interface {
	Execute(ctx context.Context, address, namespace, chainID string, tx *ethtypes.Transaction) (string, error)
	WithStorage(storage logical.Storage) SignTransactionUseCase
}

type SignQuorumPrivateTransactionUseCase interface {
	Execute(ctx context.Context, address, namespace string, tx *quorumtypes.Transaction) (string, error)
	WithStorage(storage logical.Storage) SignQuorumPrivateTransactionUseCase
}

type SignEEATransactionUseCase interface {
	Execute(
		ctx context.Context,
		address, namespace string, chainID string,
		tx *ethtypes.Transaction,
		privateArgs *entities.PrivateETHTransactionParams,
	) (string, error)
	WithStorage(storage logical.Storage) SignEEATransactionUseCase
}

type ListNamespacesUseCase interface {
	Execute(ctx context.Context) ([]string, error)
	WithStorage(storage logical.Storage) ListNamespacesUseCase
}
