package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"

	cryptoutils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/crypto"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	usecases "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/consensys/quorum/common/hexutil"

	"github.com/consensys/quorum/crypto"
)

// createAccountUseCase is a use case to create a new Ethereum account
type createAccountUseCase struct {
	storage logical.Storage
}

// NewCreateAccountUseCase creates a new CreateAccountUseCase
func NewCreateAccountUseCase() usecases.CreateAccountUseCase {
	return &createAccountUseCase{}
}

func (uc createAccountUseCase) WithStorage(storage logical.Storage) usecases.CreateAccountUseCase {
	uc.storage = storage
	return &uc
}

// Execute creates a new Ethereum account and stores it in the Vault
func (uc *createAccountUseCase) Execute(ctx context.Context, namespace, importedPrivKey string) (*entities.ETHAccount, error) {
	logger := log.FromContext(ctx).With("namespace", namespace)
	logger.Debug("creating new Ethereum account")

	var privKey = new(ecdsa.PrivateKey)
	var err error
	if importedPrivKey == "" {
		privKey, err = cryptoutils.NewSecp256k1()
		if err != nil {
			errMessage := "failed to generate Ethereum private key"
			logger.With("error", err).Error(errMessage)
			return nil, errors.CryptoOperationError(errMessage)
		}
	} else {
		privKey, err = cryptoutils.ImportSecp256k1(importedPrivKey)
		if err != nil {
			errMessage := "failed to import Ethereum private key, please verify that the provided private key is valid"
			logger.With("error", err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}
	}

	account := &entities.ETHAccount{
		PrivateKey:          hex.EncodeToString(crypto.FromECDSA(privKey)),
		Address:             crypto.PubkeyToAddress(privKey.PublicKey).Hex(),
		PublicKey:           hexutil.Encode(crypto.FromECDSAPub(&privKey.PublicKey)),
		CompressedPublicKey: hexutil.Encode(crypto.CompressPubkey(&privKey.PublicKey)),
		Namespace:           namespace,
	}

	err = storage.StoreJSON(ctx, uc.storage, storage.ComputeEthereumStorageKey(account.Address, account.Namespace), account)
	if err != nil {
		return nil, err
	}

	logger.With("address", account.Address).Info("Ethereum account created successfully")
	return account, nil
}
