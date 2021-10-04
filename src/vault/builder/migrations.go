package builder

import (
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases/migrations"
)

type migrationsUseCases struct {
	ethToKeys usecases.EthereumToKeysUseCase
}

func NewMigrationsUseCases(ethUseCases usecases.ETHUseCases, keysUseCases usecases.KeysUseCases) usecases.MigrationsUseCases {
	return &migrationsUseCases{
		ethToKeys: migrations.NewEthToKeysUseCase(ethUseCases, keysUseCases),
	}
}

func (ucs *migrationsUseCases) EthereumToKeys() usecases.EthereumToKeysUseCase {
	return ucs.ethToKeys
}
