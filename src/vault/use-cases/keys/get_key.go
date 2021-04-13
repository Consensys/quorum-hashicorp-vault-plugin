package keys

import (
	"context"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	usecases "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/logical"
)

type getKeyUseCase struct {
	storage logical.Storage
}

func NewGetKeyUseCase() usecases.GetKeyUseCase {
	return &getKeyUseCase{}
}

func (uc getKeyUseCase) WithStorage(storage logical.Storage) usecases.GetKeyUseCase {
	uc.storage = storage
	return &uc
}

func (uc *getKeyUseCase) Execute(ctx context.Context, id, namespace string) (*entities.Key, error) {
	logger := log.FromContext(ctx).With("namespace", namespace).With("id", id)
	logger.Debug("getting key pair")

	key := &entities.Key{}
	err := storage.GetJSON(ctx, uc.storage, storage.ComputeKeysStorageKey(id, namespace), key)
	if err != nil {
		return nil, err
	}

	logger.Debug("key pair found successfully")
	return key, nil
}
