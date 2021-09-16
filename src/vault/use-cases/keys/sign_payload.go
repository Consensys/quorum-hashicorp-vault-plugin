package keys

import (
	"context"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/encoding"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/hashicorp/go-hclog"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/entities"
	usecases "github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases"
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

	dataBytes, err := encoding.DecodeBase64(data)
	if err != nil {
		errMessage := "data must be a base64 string"
		logger.With("error", err).Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}

	key, err := uc.getKeyUC.Execute(ctx, id, namespace)
	if err != nil {
		return "", err
	}
	privKeyB, _ := encoding.DecodeBase64(key.PrivateKey)

	switch {
	case key.Algorithm == entities.EDDSA && key.Curve == entities.BN254:
		return uc.signEDDSA(logger, privKeyB, dataBytes)
	case key.Algorithm == entities.ECDSA && key.Curve == entities.Secp256k1:
		return uc.signECDSA(logger, privKeyB, dataBytes)
	default:
		errMessage := "invalid signing algorithm/elliptic curve combination"
		logger.Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}
}

func (uc *signPayloadUseCase) signECDSA(logger hclog.Logger, privKeyB, data []byte) (string, error) {
	if len(data) != crypto.DigestLength {
		return "", errors.InvalidParameterError("data is required to be exactly %d bytes (%d)", crypto.DigestLength, len(data))
	}

	ecdsaPrivKey, err := crypto.ToECDSA(privKeyB)
	if err != nil {
		errMessage := "failed to parse ECDSA private key"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	signatureB, err := crypto.Sign(data, ecdsaPrivKey)
	if err != nil {
		errMessage := "failed to sign payload with ECDSA"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	// We remove the recID from the signature (last byte).
	return encoding.EncodeToBase64(signatureB[:len(signatureB)-1]), nil
}

func (uc *signPayloadUseCase) signEDDSA(logger hclog.Logger, privKeyB, data []byte) (string, error) {
	privKey := eddsa.PrivateKey{}
	_, err := privKey.SetBytes(privKeyB)
	if err != nil {
		errMessage := "failed to parse EDDSA private key"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	signatureB, err := privKey.Sign(data, hash.MIMC_BN254.New("seed"))
	if err != nil {
		errMessage := "failed to sign payload with EDDSA"
		logger.With("error", err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage)
	}

	return encoding.EncodeToBase64(signatureB), nil
}
