package zksnarks

import (
	"context"
	"crypto/sha256"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	eddsa "github.com/consensys/gnark/crypto/signature/eddsa/bn256"
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
	
	// @TODO Generate new keys till we can deserialize stored keys correctly
	publicKey, privKey := eddsa.New(account.Seed, sha256.New())
	////
	
	logger.With("public key", account.PublicKey).Debug("signing with account")
	
	signature, err := eddsa.Sign([]byte(data), publicKey, privKey)
	if err != nil {
		errMessage := "failed to sign payload using EDDSA"
		logger.With("error", err).Error(errMessage)
		return "", err
	}

	logger.Info("payload signed successfully")
	// @TODO Integrate gnark serialization for signature
	return signature.R.X.String(), nil
}
