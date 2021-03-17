package zksnarks

import (
	"context"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/crypto"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	usecases "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type createAccountUseCase struct {
	storage logical.Storage
}

func NewCreateAccountUseCase() usecases.CreateZksAccountUseCase {
	return &createAccountUseCase{}
}

func (uc createAccountUseCase) WithStorage(storage logical.Storage) usecases.CreateZksAccountUseCase {
	uc.storage = storage
	return &uc
}

func (uc *createAccountUseCase) Execute(ctx context.Context, namespace string) (*entities.ZksAccount, error) {
	logger := log.FromContext(ctx).With("namespace", namespace)
	logger.Debug("creating new zk-snarks bn256 account")

	privKey, err := crypto.NewBN256()
	if err != nil {
		errMessage := "failed to generate key"
		logger.With("error", err).Error(errMessage)
		return nil, err
	}

	account := &entities.ZksAccount{
		Algorithm:  entities.EDDSA,
		Curve:      entities.BN256,
		PrivateKey: hexutil.Encode(privKey.Bytes()),
		PublicKey:  hexutil.Encode(privKey.Public().Bytes()),
		Namespace:  namespace,
	}

	err = storage.StoreJSON(ctx, uc.storage, storage.ComputeZksStorageKey(account.PublicKey, account.Namespace), account)
	if err != nil {
		errMessage := "failed to create account entry"
		logger.With("error", err).Error(errMessage)
		return nil, err
	}

	logger.With("pub_key", account.PublicKey).Info("zk-snarks bn256 account created successfully")
	return account, nil
}
