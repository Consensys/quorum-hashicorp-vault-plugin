package zksnarks

import (
	"context"
	"crypto/sha256"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	eddsa "github.com/consensys/gnark/crypto/signature/eddsa/bn256"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/hashicorp/vault/sdk/logical"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
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

func (uc *signPayloadUseCase) Execute(ctx context.Context, pubKeyHex, namespace, data string) (string, error) {
	logger := log.FromContext(ctx).With("namespace", namespace).With("pub_key", pubKeyHex)
	logger.Debug("signing message")

	account, err := uc.getAccountUC.Execute(ctx, pubKeyHex, namespace)
	if err != nil {
		return "", err
	}
	
	privKey := eddsa.PrivateKey{}
	privKeyB, _ := hexutil.Decode(account.PrivateKey)
	_, err = privKey.SetBytes([]byte(privKeyB))
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
