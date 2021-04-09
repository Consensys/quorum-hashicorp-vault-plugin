package testutils

import (
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func FakePrivateETHTransactionParams() *entities.PrivateETHTransactionParams {
	return &entities.PrivateETHTransactionParams{
		PrivateFrom:   "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
		PrivateFor:    []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="},
		PrivateTxType: entities.PrivateTxTypeRestricted,
	}
}

func FakeETHAccount() *entities.ETHAccount {
	return &entities.ETHAccount{
		Namespace:           "_",
		Address:             ethcommon.HexToAddress(RandHexString(12)).String(),
		PublicKey:           ethcommon.HexToHash(RandHexString(12)).String(),
		CompressedPublicKey: ethcommon.HexToHash(RandHexString(12)).String(),
	}
}
