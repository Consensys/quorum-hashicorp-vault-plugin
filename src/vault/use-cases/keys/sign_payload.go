package keys

import (
	"context"
	"crypto/sha256"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/hashicorp/go-hclog"

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
		return uc.signEDDSA(logger, key.PrivateKey, data)
	case key.Algorithm == entities.ECDSA && key.Curve == entities.Secp256k1:
		return uc.signECDSA(logger, key.PrivateKey, data)
	default:
		errMessage := "invalid signing algorithm/elliptic curve combination"
		logger.Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}
}

func (uc *signPayloadUseCase) signECDSA(logger hclog.Logger, privKeyString, data string) (string, error) {
	ecdsaPrivKey, err := crypto.HexToECDSA(privKeyString)
	if err != nil {
		errMessage := "failed to parse ECDSA private key"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	signatureB, err := crypto.Sign(crypto.Keccak256([]byte(data)), ecdsaPrivKey)
	if err != nil {
		errMessage := "failed to sign payload with ECDSA"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	return hexutil.Encode(signatureB), nil
}

func (uc *signPayloadUseCase) signEDDSA(logger hclog.Logger, privKeyString, data string) (string, error) {
	privKey := eddsa.PrivateKey{}
	privKeyB, _ := hexutil.Decode(privKeyString)
	_, err := privKey.SetBytes(privKeyB)
	if err != nil {
		errMessage := "failed to parse EDDSA private key"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	signatureB, err := privKey.Sign([]byte(data), sha256.New())
	if err != nil {
		errMessage := "failed to sign payload with EDDSA"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	return hexutil.Encode(signatureB), nil
}
