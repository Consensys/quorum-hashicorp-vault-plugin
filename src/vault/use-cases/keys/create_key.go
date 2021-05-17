package keys

import (
	"context"
	"crypto/ecdsa"
	crypto2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/crypto"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/encoding"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	usecases "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/vault/sdk/logical"
	"time"
)

type createKeyUseCase struct {
	storage logical.Storage
}

func NewCreateKeyUseCase() usecases.CreateKeyUseCase {
	return &createKeyUseCase{}
}

func (uc createKeyUseCase) WithStorage(storage logical.Storage) usecases.CreateKeyUseCase {
	uc.storage = storage
	return &uc
}

func (uc *createKeyUseCase) Execute(ctx context.Context, namespace, id, algo, curve, importedPrivKey string, tags map[string]string) (*entities.Key, error) {
	logger := log.FromContext(ctx).
		With("namespace", namespace).
		With("algorithm", algo).
		With("curve", curve).
		With("id", id)
	logger.Debug("creating new key")

	timeNow := time.Now()
	key := &entities.Key{
		ID:        id,
		Algorithm: algo,
		Curve:     curve,
		Namespace: namespace,
		Tags:      tags,
		Version:   1,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	switch {
	case algo == entities.EDDSA && curve == entities.BN254:
		privKey, err := uc.eddsaBN254(importedPrivKey)
		if err != nil {
			errMessage := "failed to generate EDDSA/BN254 key pair"
			logger.With("error", err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		key.PrivateKey = encoding.EncodeToBase64(privKey.Bytes())
		key.PublicKey = encoding.EncodeToBase64(privKey.Public().Bytes())
	case algo == entities.ECDSA && curve == entities.Secp256k1:
		privKey, err := uc.ecdsaSecp256k1(importedPrivKey)
		if err != nil {
			errMessage := "failed to generate Secp256k1/ECDSA key pair"
			logger.With("error", err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		key.PrivateKey = encoding.EncodeToBase64(crypto.FromECDSA(privKey))
		key.PublicKey = encoding.EncodeToBase64(crypto.FromECDSAPub(&privKey.PublicKey))
	default:
		errMessage := "invalid signing algorithm/elliptic curve combination"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	err := storage.StoreJSON(ctx, uc.storage, storage.ComputeKeysStorageKey(id, key.Namespace), key)
	if err != nil {
		return nil, err
	}

	logger.With("pub_key", key.PublicKey).Info("key pair created successfully")
	return key, nil
}

func (*createKeyUseCase) eddsaBN254(importedPrivKey string) (eddsa.PrivateKey, error) {
	if importedPrivKey == "" {
		key, err := crypto2.NewBN254()
		if err != nil {
			return key, err
		}

		return key, nil
	}

	key := eddsa.PrivateKey{}
	privKeyBytes, err := encoding.DecodeBase64(importedPrivKey)
	if err != nil {
		return key, err
	}

	_, err = key.SetBytes(privKeyBytes)
	if err != nil {
		return key, err
	}

	return key, nil
}

func (*createKeyUseCase) ecdsaSecp256k1(importedPrivKey string) (*ecdsa.PrivateKey, error) {
	if importedPrivKey == "" {
		key, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}

		return key, nil
	}

	privKeyBytes, err := encoding.DecodeBase64(importedPrivKey)
	if err != nil {
		return nil, err
	}

	key, err := crypto.ToECDSA(privKeyBytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
