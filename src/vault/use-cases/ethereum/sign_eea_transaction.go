package ethereum

import (
	"context"
	signing "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/crypto/ethereum"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	usecases "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/consensys/quorum/core/types"
	"github.com/consensys/quorum/crypto"
	"github.com/hashicorp/vault/sdk/logical"
)

// signEEATxUseCase is a use case to sign a Quorum private transaction using an existing account
type signEEATxUseCase struct {
	getAccountUC usecases.GetAccountUseCase
}

// NewSignEEATransactionUseCase creates a new signEEATxUseCase
func NewSignEEATransactionUseCase(getAccountUC usecases.GetAccountUseCase) usecases.SignEEATransactionUseCase {
	return &signEEATxUseCase{
		getAccountUC: getAccountUC,
	}
}

func (uc signEEATxUseCase) WithStorage(storage logical.Storage) usecases.SignEEATransactionUseCase {
	uc.getAccountUC = uc.getAccountUC.WithStorage(storage)
	return &uc
}

// Execute signs an EEA transaction
func (uc *signEEATxUseCase) Execute(
	ctx context.Context,
	address, namespace, chainID string,
	tx *types.Transaction,
	privateArgs *entities.PrivateETHTransactionParams,
) (string, error) {
	logger := log.FromContext(ctx).With("namespace", namespace).With("address", address)
	logger.Debug("signing eea private transaction")

	account, err := uc.getAccountUC.Execute(ctx, address, namespace)
	if err != nil {
		return "", err
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		errMessage := "failed to parse private key"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	signature, err := signing.SignEEATransaction(tx, privateArgs, chainID, ecdsaPrivKey)
	if err != nil {
		errMessage := "failed to sign eea transaction"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	logger.Info("eea private transaction signed successfully")
	return hexutil.Encode(signature), nil
}
