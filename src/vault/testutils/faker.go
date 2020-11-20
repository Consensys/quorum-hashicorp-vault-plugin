package testutils

import (
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/ethereum/go-ethereum/common"
	"math/rand"
)

func FakeETHAccount() *entities.ETHAccount {
	return &entities.ETHAccount{
		Address:             common.HexToAddress(randHexString(12)).String(),
		PublicKey:           common.HexToHash(randHexString(12)).String(),
		CompressedPublicKey: common.HexToHash(randHexString(12)).String(),
		Namespace:           "_",
	}
}

func randHexString(n int) string {
	var letterRunes = []rune("abcdef0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
