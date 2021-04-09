package keys

import (
	"context"
	"crypto/sha256"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	usecases "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	eddsa "github.com/consensys/gnark/crypto/signature/eddsa/bn256"
	"github.com/consensys/quorum/common/hexutil"
	"github.com/consensys/quorum/crypto"
	"github.com/hashicorp/vault/sdk/logical"
)

type signPayloadUseCase struct {
	getKeyUC usecases.GetKeyUseCase
}

func NewSignUseCase(getKeyUC usecases.GetKeyUseCase) usecases.KeysSignUseCase {
	return &signPayloadUseCase{
		getKeyUC: getKeyUC,
	}
}

func (uc signPayloadUseCase) WithStorage(storage logical.Storage) usecases.KeysSignUseCase {
	uc.getKeyUC = uc.getKeyUC.WithStorage(storage)
	return &uc
}

func (uc *signPayloadUseCase) Execute(ctx context.Context, id, namespace, data string) (string, error) {
	logger := log.FromContext(ctx).With("namespace", namespace).With("id", id)
	logger.Debug("signing message")

	key, err := uc.getKeyUC.Execute(ctx, id, namespace)
	if err != nil {
		return "", err
	}

	switch {
	case key.Algorithm == entities.EDDSA && key.Curve == entities.BN256:
		signature, err := uc.signEDDSA(key.PrivateKey, data)
		if err != nil {
			errMessage := "failed to sign payload using EDDSA"
			logger.With("error", err).Error(errMessage)
			return "", errors.EncodingError(errMessage)
		}

		return signature, nil
	case key.Algorithm == entities.ECDSA && key.Curve == entities.Secp256k1:
		signature, err := uc.signECDSA(key.PrivateKey, data)
		if err != nil {
			errMessage := "failed to sign payload using ECDSA"
			logger.With("error", err).Error(errMessage)
			return "", errors.EncodingError(errMessage)
		}

		return signature, nil
	default:
		errMessage := "invalid signing algorithm/elliptic curve combination"
		logger.Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}
}

func (uc *signPayloadUseCase) signECDSA(privKeyString, data string) (string, error) {
	ecdsaPrivKey, err := crypto.HexToECDSA(privKeyString)
	if err != nil {
		return "", err
	}

	signatureB, err := crypto.Sign(crypto.Keccak256([]byte(data)), ecdsaPrivKey)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(signatureB), nil
}

func (uc *signPayloadUseCase) signEDDSA(privKeyString, data string) (string, error) {
	privKey := eddsa.PrivateKey{}
	privKeyB, _ := hexutil.Decode(privKeyString)
	_, err := privKey.SetBytes(privKeyB)
	if err != nil {
		return "", err
	}

	signatureB, err := privKey.Sign([]byte(data), sha256.New())
	if err != nil {
		return "", err
	}

	return hexutil.Encode(signatureB), nil
}
