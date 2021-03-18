package keys

import (
	"context"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	usecases "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/logical"
)

type listKeysUseCase struct {
	storage logical.Storage
}

func NewListKeysUseCase() usecases.ListKeysUseCase {
	return &listKeysUseCase{}
}

func (uc listKeysUseCase) WithStorage(storage logical.Storage) usecases.ListKeysUseCase {
	uc.storage = storage
	return &uc
}

func (uc *listKeysUseCase) Execute(ctx context.Context, namespace string) ([]string, error) {
	logger := log.FromContext(ctx).With("namespace", namespace)
	logger.Debug("listing key pairs")

	return uc.storage.List(ctx, storage.ComputeKeysStorageKey("", namespace))
}
