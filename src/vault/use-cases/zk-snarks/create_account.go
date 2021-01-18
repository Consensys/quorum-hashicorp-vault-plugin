package zksnarks

import (
	"bytes"
	"context"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	eddsa "github.com/consensys/gnark/crypto/signature/eddsa/bn256"
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

	var seed = make([]byte, 32)
	for i, v := range utils.GenerateRandomSeed(32) {
		seed[i] = v
	}

	// Usually standards implementations of eddsa do not require the choice of a specific hash function (usually it's SHA256). 
	// Here we needed to allow the choice of the hash so we can chose a hash function that is easily programmable in a snark circuit.
	// Same hFunc should be used for sign and verify
	privKey, err := eddsa.GenerateKey(bytes.NewReader(seed))
	if err != nil {
		errMessage := "failed to generate key"
		logger.With("error", err).Error(errMessage)
		return nil, err
	}

	pubKey := privKey.Public()
	account := &entities.ZksAccount{
		Algorithm:  entities.ZksAlgorithmEDDSA,
		Curve:      entities.ZksCurveBN256,
		PrivateKey: hexutil.Encode(privKey.Bytes()),
		PublicKey:  hexutil.Encode(pubKey.Bytes()),
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
