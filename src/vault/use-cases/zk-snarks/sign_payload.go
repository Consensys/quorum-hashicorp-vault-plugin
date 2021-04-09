package zksnarks

import (
	"context"
	"crypto/sha256"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	eddsa "github.com/consensys/gnark/crypto/signature/eddsa/bn256"
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

	account, err := uc.getAccountUC.Execute(ctx, pubKey, namespace)
	if err != nil {
		return "", err
	}

	privKey := eddsa.PrivateKey{}
	privKeyB, _ := hexutil.Decode(account.PrivateKey)
	_, err = privKey.SetBytes(privKeyB)
	if err != nil {
		errMsg := "failed to deserialize private key"
		logger.With("error", err).Error(errMsg)
		return "", errors.EncodingError(errMsg)
	}

	signatureB, err := privKey.Sign([]byte(data), sha256.New())
	if err != nil {
		errMessage := "failed to sign payload using EDDSA"
		logger.With("error", err).Error(errMessage)
		return "", err
	}

	logger.Info("payload signed successfully")
	return hexutil.Encode(signatureB), nil
}
