package ethereum

import (
	"encoding/base64"
	"math/big"

	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/core/types"
)

func GetEncodedPrivateFrom(privateFrom string) ([]byte, error) {
	privateFromEncoded, err := base64.StdEncoding.DecodeString(privateFrom)
	if err != nil {
		return nil, err
	}

	return privateFromEncoded, nil
}

func GetEncodedPrivateRecipient(privacyGroupID string, privateFor []string) (interface{}, error) {
	var privateRecipientEncoded interface{}
	var err error
	if privacyGroupID != "" {
		privateRecipientEncoded, err = base64.StdEncoding.DecodeString(privacyGroupID)
		if err != nil {
			return nil, err
		}
	} else {
		var privateForByteSlice [][]byte
		for _, v := range privateFor {
			b, der := base64.StdEncoding.DecodeString(v)
			if der != nil {
				return nil, err
			}
			privateForByteSlice = append(privateForByteSlice, b)
		}
		privateRecipientEncoded = privateForByteSlice
	}

	return privateRecipientEncoded, nil
}

func GetEIP155Signer(chainID string) types.Signer {
	chainIDBigInt := new(big.Int)
	chainIDBigInt, _ = chainIDBigInt.SetString(chainID, 10)
	return types.NewEIP155Signer(chainIDBigInt)
}

func GetQuorumPrivateTxSigner() quorumtypes.Signer {
	return quorumtypes.QuorumPrivateTxSigner{}
}
