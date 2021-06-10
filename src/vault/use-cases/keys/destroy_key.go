package keys

import (
	"context"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	usecases "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/logical"
)

type destroyKeyUseCase struct {
	storage logical.Storage
}

func NewDestroyKeyUseCase() usecases.DestroyKeyUseCase {
	return &destroyKeyUseCase{}
}

func (uc destroyKeyUseCase) WithStorage(storage logical.Storage) usecases.DestroyKeyUseCase {
	uc.storage = storage
	return &uc
}

func (uc *destroyKeyUseCase) Execute(ctx context.Context, namespace, id string) error {
	logger := log.FromContext(ctx).With("namespace", namespace).With("id", id)
	logger.Debug("permanently deleting key")

	err := storage.DestroyJSON(ctx, uc.storage, storage.ComputeKeysStorageKey(id, namespace))
	if err != nil {
		return err
	}

	logger.Info("key pair permanently deleted")
	return nil
}
