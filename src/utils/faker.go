package utils

import (
	"math/rand"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/ethereum/go-ethereum/common"
)

func FakeETHAccount() *entities.ETHAccount {
	return &entities.ETHAccount{
		Address:             common.HexToAddress(randHexString(12)).String(),
		PublicKey:           common.HexToHash(randHexString(12)).String(),
		CompressedPublicKey: common.HexToHash(randHexString(12)).String(),
		Namespace:           "_",
	}
}

func FakeZksAccount() *entities.ZksAccount {
	return &entities.ZksAccount{
		Algorithm: entities.ZksAlgorithmEDDSA,
		Curve: entities.ZksCurveBN256,
		PublicKey: "0x10d10e9f17a88d51c42380c14da49e237b4b3f03c5cdce8f470ca782506eb5f113733d92b86e28b7e6354bb88a2d6bb9b104b0de3698b993f735f31cc979f7bd",
		PrivateKey: common.HexToHash(randHexString(96)).String(),
		Namespace: "_",
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
