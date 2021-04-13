package keys

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/crypto"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	usecases "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	eddsa "github.com/consensys/gnark/crypto/signature/eddsa/bn256"
	"github.com/consensys/quorum/common/hexutil"
	crypto2 "github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/vault/sdk/logical"
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

	key := &entities.Key{
		ID:        id,
		Algorithm: algo,
		Curve:     curve,
		Namespace: namespace,
		Tags:      tags,
	}

	switch {
	case algo == entities.EDDSA && curve == entities.BN256:
		privKey, err := uc.eddsaBN256(importedPrivKey)
		if err != nil {
			errMessage := "failed to generate EDDSA/BN256 key pair"
			logger.With("error", err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		key.PrivateKey = hexutil.Encode(privKey.Bytes())
		key.PublicKey = hexutil.Encode(privKey.Public().Bytes())
	case algo == entities.ECDSA && curve == entities.Secp256k1:
		privKey, err := uc.ecdsaSecp256k1(importedPrivKey)
		if err != nil {
			errMessage := "failed to generate Secp256k1/ECDSA key pair"
			logger.With("error", err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		key.PrivateKey = hex.EncodeToString(crypto2.FromECDSA(privKey))
		key.PublicKey = hexutil.Encode(crypto2.FromECDSAPub(&privKey.PublicKey))
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

func (*createKeyUseCase) eddsaBN256(importedPrivKey string) (eddsa.PrivateKey, error) {
	if importedPrivKey == "" {
		key, err := crypto.NewBN256()
		if err != nil {
			return key, errors.CryptoOperationError(err.Error())
		}

		return key, nil
	} else {
		key, err := crypto.ImportBN256(importedPrivKey)
		if err != nil {
			return key, errors.InvalidParameterError(err.Error())
		}

		return key, nil
	}
}

func (*createKeyUseCase) ecdsaSecp256k1(importedPrivKey string) (*ecdsa.PrivateKey, error) {
	if importedPrivKey == "" {
		key, err := crypto.NewSecp256k1()
		if err != nil {
			return key, errors.CryptoOperationError(err.Error())
		}

		return key, nil

	} else {
		key, err := crypto.ImportSecp256k1(importedPrivKey)
		if err != nil {
			return key, errors.InvalidParameterError(err.Error())
		}

		return key, nil
	}
}
