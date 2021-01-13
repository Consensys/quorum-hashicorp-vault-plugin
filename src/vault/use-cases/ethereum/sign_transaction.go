package ethereum

import (
	"context"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/vault/sdk/logical"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/crypto/ethereum/signing"

	"github.com/consensys/quorum/common/hexutil"
)

// signTxUseCase is a use case to sign an ethereum transaction using an existing account
type signTxUseCase struct {
	getAccountUC usecases.GetAccountUseCase
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase
func NewSignTransactionUseCase(getAccountUC usecases.GetAccountUseCase) usecases.SignTransactionUseCase {
	return &signTxUseCase{
		getAccountUC: getAccountUC,
	}
}

func (uc signTxUseCase) WithStorage(storage logical.Storage) usecases.SignTransactionUseCase {
	uc.getAccountUC = uc.getAccountUC.WithStorage(storage)
	return &uc
}

// Execute signs an ethereum transaction
func (uc *signTxUseCase) Execute(ctx context.Context, address, namespace, chainID string, tx *ethtypes.Transaction) (string, error) {
	logger := log.FromContext(ctx).With("namespace", namespace).With("address", address)
	logger.Debug("signing ethereum transaction")

	account, err := uc.getAccountUC.Execute(ctx, address, namespace)
	if err != nil {
		return "", err
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		errMessage := "failed to parse private key"
		logger.With("error", err).Error(errMessage)
		return "", err
	}

	signature, err := signing.SignTransaction(tx, ecdsaPrivKey, signing.GetEIP155Signer(chainID))
	if err != nil {
		errMessage := "failed to sign transaction using ECDSA"
		logger.With("error", err).Error(errMessage)
		return "", err
	}

	logger.Info("ethereum transaction signed successfully")
	return hexutil.Encode(signature), nil
}
