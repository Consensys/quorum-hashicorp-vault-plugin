package ethereum

import (
	"context"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/vault/sdk/logical"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/crypto/ethereum/signing"

	"github.com/consensys/quorum/common/hexutil"
	quorumtypes "github.com/consensys/quorum/core/types"
)

// signQuorumPrivateTxUseCase is a use case to sign a Quorum private transaction using an existing account
type signQuorumPrivateTxUseCase struct {
	getAccountUC usecases.GetAccountUseCase
}

// NewSignQuorumPrivateTransactionUseCase creates a new signQuorumPrivateTxUseCase
func NewSignQuorumPrivateTransactionUseCase(getAccountUC usecases.GetAccountUseCase) usecases.SignQuorumPrivateTransactionUseCase {
	return &signQuorumPrivateTxUseCase{
		getAccountUC: getAccountUC,
	}
}

func (uc signQuorumPrivateTxUseCase) WithStorage(storage logical.Storage) usecases.SignQuorumPrivateTransactionUseCase {
	uc.getAccountUC = uc.getAccountUC.WithStorage(storage)
	return &uc
}

// Execute signs a Quorum private transaction
func (uc *signQuorumPrivateTxUseCase) Execute(ctx context.Context, address, namespace string, tx *quorumtypes.Transaction) (string, error) {
	logger := apputils.Logger(ctx).With("namespace", namespace).With("address", address)
	logger.Debug("signing quorum private transaction")

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

	signature, err := signing.SignQuorumPrivateTransaction(tx, ecdsaPrivKey, signing.GetQuorumPrivateTxSigner())
	if err != nil {
		errMessage := "failed to sign quorum private transaction"
		logger.With("error", err).Error(errMessage)
		return "", err
	}

	logger.Info("quorum private transaction signed successfully")
	return hexutil.Encode(signature), nil
}
