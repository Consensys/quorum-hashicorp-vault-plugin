package zksnarks

import (
	"context"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/consensys/gnark-crypto/crypto/hash"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type signPayloadUseCase struct {
	getAccountUC usecases.GetZksAccountUseCase
}

func NewSignUseCase(getAccountUC usecases.GetZksAccountUseCase) usecases.ZksSignUseCase {
	return &signPayloadUseCase{
		getAccountUC: getAccountUC,
	}
}

func (uc signPayloadUseCase) WithStorage(storage logical.Storage) usecases.ZksSignUseCase {
	uc.getAccountUC = uc.getAccountUC.WithStorage(storage)
	return &uc
}

func (uc *signPayloadUseCase) Execute(ctx context.Context, pubKey, namespace, data string) (string, error) {
	logger := log.FromContext(ctx).With("namespace", namespace).With("pub_key", pubKey)
	logger.Debug("signing message")

	dataBytes, err := hexutil.Decode(data)
	if err != nil {
		errMessage := "data must be a hex string"
		logger.With("error", err).Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}

	account, err := uc.getAccountUC.Execute(ctx, pubKey, namespace)
	if err != nil {
		return "", err
	}

	privKey := eddsa.PrivateKey{}
	privKeyB, _ := hexutil.Decode(account.PrivateKey)
	_, err = privKey.SetBytes(privKeyB)
	if err != nil {
		errMessage := "failed to parse private key"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	signatureB, err := privKey.Sign(dataBytes, hash.MIMC_BN254.New("seed"))
	if err != nil {
		errMessage := "failed to sign payload"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	logger.Info("payload signed successfully")
	return hexutil.Encode(signatureB), nil
}
