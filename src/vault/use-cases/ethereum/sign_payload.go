package ethereum

import (
	"context"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/consensys/quorum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
)

// signPayloadUseCase is a use case to sign an arbitrary payload usign an existing Ethereum account
type signPayloadUseCase struct {
	getAccountUC usecases.GetAccountUseCase
}

// NewSignUseCase creates a new SignUseCase
func NewSignUseCase(getAccountUC usecases.GetAccountUseCase) usecases.SignUseCase {
	return &signPayloadUseCase{
		getAccountUC: getAccountUC,
	}
}

func (uc signPayloadUseCase) WithStorage(storage logical.Storage) usecases.SignUseCase {
	uc.getAccountUC = uc.getAccountUC.WithStorage(storage)
	return &uc
}

// Execute signs an arbitrary payload using an existing Ethereum account
func (uc *signPayloadUseCase) Execute(ctx context.Context, address, namespace, data string) (string, error) {
	logger := apputils.Logger(ctx).With("namespace", namespace).With("address", address)
	logger.Debug("signing message")

	account, err := uc.getAccountUC.Execute(ctx, address, namespace)
	if err != nil {
		return "", err
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		errMessage := "failed to parse secp256k1 private key"
		logger.With("error", err).Error(errMessage)
		return "", err
	}

	signature, err := crypto.Sign(crypto.Keccak256([]byte(data)), ecdsaPrivKey)
	if err != nil {
		errMessage := "failed to sign payload using ECDSA"
		logger.With("error", err).Error(errMessage)
		return "", err
	}

	logger.Info("payload signed successfully")
	return hexutil.Encode(signature), nil
}
