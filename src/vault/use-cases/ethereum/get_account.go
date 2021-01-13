package ethereum

import (
	"context"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/logical"
)

// getAccountUseCase is a use case to get an Ethereum account
type getAccountUseCase struct {
	storage logical.Storage
}

// NewGetAccountUseCase creates a new GetAccountUseCase
func NewGetAccountUseCase() usecases.GetAccountUseCase {
	return &getAccountUseCase{}
}

func (uc getAccountUseCase) WithStorage(storage logical.Storage) usecases.GetAccountUseCase {
	uc.storage = storage
	return &uc
}

// Execute creates a new Ethereum account and stores it in the Vault
func (uc *getAccountUseCase) Execute(ctx context.Context, address, namespace string) (*entities.ETHAccount, error) {
	logger := log.FromContext(ctx).With("namespace", namespace).With("address", address)
	logger.Debug("getting Ethereum account")

	account := &entities.ETHAccount{}
	err := storage.GetJSON(ctx, uc.storage, 
		storage.ComputeEthereumStorageKey(address, namespace), account)
	if err != nil {
		logger.With("error", err).Error("failed to retrieve account from vault")
		return nil, err
	}

	logger.Debug("Ethereum account found successfully")
	return account, nil
}
