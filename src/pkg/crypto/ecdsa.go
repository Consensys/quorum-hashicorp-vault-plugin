package crypto

import (
	"crypto/ecdsa"

	"github.com/consensys/quorum/crypto"
)

func ImportSecp256k1(importedPrivKey string) (*ecdsa.PrivateKey, error) {
	privKey, err := crypto.HexToECDSA(importedPrivKey)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func NewSecp256k1() (*ecdsa.PrivateKey, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return privKey, nil
}
