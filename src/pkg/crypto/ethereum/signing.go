package ethereum

import (
	"crypto/ecdsa"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/entities"
	"math/big"

	"github.com/consensys/quorum/core/types"
	"github.com/consensys/quorum/crypto"
)

func SignTransaction(tx *types.Transaction, privKey *ecdsa.PrivateKey, signer types.Signer) ([]byte, error) {
	h := signer.Hash(tx)
	decodedSignature, err := crypto.Sign(h[:], privKey)
	if err != nil {
		return nil, err
	}

	return decodedSignature, nil
}

func SignQuorumPrivateTransaction(tx *types.Transaction, privKey *ecdsa.PrivateKey, signer types.Signer) ([]byte, error) {
	h := signer.Hash(tx)
	decodedSignature, err := crypto.Sign(h[:], privKey)
	if err != nil {
		return nil, err
	}

	return decodedSignature, nil
}

func SignEEATransaction(tx *types.Transaction, privateArgs *entities.PrivateETHTransactionParams, chainID string, privKey *ecdsa.PrivateKey) ([]byte, error) {
	chainIDBigInt, _ := new(big.Int).SetString(chainID, 10)
	privateFromEncoded, err := GetEncodedPrivateFrom(privateArgs.PrivateFrom)
	if err != nil {
		return nil, err
	}

	privateRecipientEncoded, err := GetEncodedPrivateRecipient(privateArgs.PrivacyGroupID, privateArgs.PrivateFor)
	if err != nil {
		return nil, err
	}

	hash, err := Hash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		chainIDBigInt,
		uint(0),
		uint(0),
		privateFromEncoded,
		privateRecipientEncoded,
		privateArgs.PrivateTxType,
	})
	if err != nil {
		return nil, err
	}

	signature, err := crypto.Sign(hash[:], privKey)
	if err != nil {
		return nil, err
	}

	return signature, err
}
