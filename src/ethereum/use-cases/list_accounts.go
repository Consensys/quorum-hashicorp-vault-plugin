package ethereum

import (
	"context"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/use-cases/utils"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/logical"
)

// listAccountsUseCase is a use case to get a list of Ethereum accounts
type listAccountsUseCase struct {
	storage logical.Storage
}

// NewListAccountUseCase creates a new ListAccountsUseCase
func NewListAccountsUseCase() ListAccountsUseCase {
	return &listAccountsUseCase{}
}

func (uc listAccountsUseCase) WithStorage(storage logical.Storage) ListAccountsUseCase {
	uc.storage = storage
	return &uc
}

// Execute creates a new Ethereum account and stores it in the Vault
func (uc *listAccountsUseCase) Execute(ctx context.Context, namespace string) ([]string, error) {
	logger := apputils.Logger(ctx).With("namespace", namespace)
	logger.Debug("listing Ethereum accounts")

	return uc.storage.List(ctx, utils.ComputeKey("", namespace))
}
