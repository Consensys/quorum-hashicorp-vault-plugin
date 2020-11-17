package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/entities"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/utils"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/consensys/quorum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
)

// createAccountUseCase is a use case to create a new Ethereum account
type createAccountUseCase struct {
	storage logical.Storage
}

// NewCreateAccountUseCase creates a new CreateAccountUseCase
func NewCreateAccountUseCase() CreateAccountUseCase {
	return &createAccountUseCase{}
}

func (uc createAccountUseCase) WithStorage(storage logical.Storage) CreateAccountUseCase {
	uc.storage = storage
	return &uc
}

// Execute creates a new Ethereum account and stores it in the Vault
func (uc *createAccountUseCase) Execute(ctx context.Context, namespace, importedPrivKey string) (*entities.ETHAccount, error) {
	logger := apputils.Logger(ctx).With("namespace", namespace)
	logger.Debug("creating new Ethereum account")

	var privKey = new(ecdsa.PrivateKey)
	var err error
	if importedPrivKey == "" {
		privKey, err = generatePrivKey(ctx)
	} else {
		privKey, err = retrievePrivKey(ctx, importedPrivKey)
	}
	if err != nil {
		return nil, err
	}

	account := &entities.ETHAccount{
		PrivateKey:          hex.EncodeToString(crypto.FromECDSA(privKey)),
		Address:             crypto.PubkeyToAddress(privKey.PublicKey).Hex(),
		PublicKey:           hexutil.Encode(crypto.FromECDSAPub(&privKey.PublicKey)),
		CompressedPublicKey: hexutil.Encode(crypto.CompressPubkey(&privKey.PublicKey)),
		Namespace:           namespace,
	}

	entry, err := logical.StorageEntryJSON(utils.ComputeKey(account.Address, account.Namespace), account)
	if err != nil {
		errMessage := "failed to create account entry"
		apputils.Logger(ctx).With("error", err).Error(errMessage)
		return nil, err
	}

	err = uc.storage.Put(ctx, entry)
	if err != nil {
		errMessage := "failed to store account in vault"
		apputils.Logger(ctx).With("error", err).Error(errMessage)
		return nil, err
	}

	logger.With("address", account.Address).Info("Ethereum account created successfully")
	return account, nil
}

func retrievePrivKey(ctx context.Context, importedPrivKey string) (*ecdsa.PrivateKey, error) {
	privKey, err := crypto.HexToECDSA(importedPrivKey)
	if err != nil {
		errMessage := "failed to import Ethereum private key, please verify that the provided private key is valid"
		apputils.Logger(ctx).With("error", err).Error(errMessage)
		return nil, err
	}

	return privKey, nil
}

func generatePrivKey(ctx context.Context) (*ecdsa.PrivateKey, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		errMessage := "failed to generate Ethereum private key"
		apputils.Logger(ctx).With("error", err).Error(errMessage)
		return nil, err
	}

	return privKey, nil
}
