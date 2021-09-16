package zksnarks

import (
	"context"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/entities"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/storage"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/logical"
)

type getAccountUseCase struct {
	storage logical.Storage
}

func NewGetAccountUseCase() usecases.GetZksAccountUseCase {
	return &getAccountUseCase{}
}

func (uc getAccountUseCase) WithStorage(storage logical.Storage) usecases.GetZksAccountUseCase {
	uc.storage = storage
	return &uc
}

func (uc *getAccountUseCase) Execute(ctx context.Context, pubKey, namespace string) (*entities.ZksAccount, error) {
	logger := log.FromContext(ctx).With("namespace", namespace).With("pub_key", pubKey)
	logger.Debug("getting zk-snarks account")

	account := &entities.ZksAccount{}
	err := storage.GetJSON(ctx, uc.storage, storage.ComputeZksStorageKey(pubKey, namespace), account)
	if err != nil {
		return nil, err
	}

	logger.Debug("zk-snarks account found successfully")
	return account, nil
}
