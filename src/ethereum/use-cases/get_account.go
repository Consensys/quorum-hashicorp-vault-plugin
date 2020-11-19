package ethereum

import (
	"context"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/entities"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/use-cases/utils"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/logical"
)

// getAccountUseCase is a use case to get an Ethereum account
type getAccountUseCase struct {
	storage logical.Storage
}

// NewGetAccountUseCase creates a new GetAccountUseCase
func NewGetAccountUseCase() GetAccountUseCase {
	return &getAccountUseCase{}
}

func (uc getAccountUseCase) WithStorage(storage logical.Storage) GetAccountUseCase {
	uc.storage = storage
	return &uc
}

// Execute creates a new Ethereum account and stores it in the Vault
func (uc *getAccountUseCase) Execute(ctx context.Context, address, namespace string) (*entities.ETHAccount, error) {
	logger := apputils.Logger(ctx).With("namespace", namespace).With("address", address)
	logger.Debug("getting Ethereum account")

	entry, err := uc.storage.Get(ctx, utils.ComputeKey(address, namespace))
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	account := &entities.ETHAccount{}
	err = entry.DecodeJSON(&account)
	if err != nil {
		return nil, err
	}

	logger.Info("Ethereum account found successfully")
	return account, nil
}
