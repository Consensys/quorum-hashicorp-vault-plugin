package src

import (
	"context"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/migrations"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/ethereum"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/keys"
	zksnarks "github.com/consensys/quorum-hashicorp-vault-plugin/src/service/zk-snarks"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/builder"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// NewVaultBackend returns the Hashicorp Vault backend
func NewVaultBackend(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	vaultPlugin := &framework.Backend{
		Help: "Quorum Hashicorp Vault Plugin",
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"ethereum/accounts/",
				"zk-snarks/accounts/",
				"keys/",
				"migrations/",
			},
		},
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}

	ethUseCases := builder.NewEthereumUseCases()
	keysUsecases := builder.NewKeysUseCases()
	ethereumController := ethereum.NewController(ethUseCases, conf.Logger)
	zkSnarksController := zksnarks.NewController(builder.NewZkSnarksUseCases(), conf.Logger)
	keysController := keys.NewController(keysUsecases, conf.Logger)
	migrationsController := migrations.NewController(builder.NewMigrationsUseCases(ethUseCases, keysUsecases), conf.Logger)
	vaultPlugin.Paths = framework.PathAppend(
		ethereumController.Paths(),
		zkSnarksController.Paths(),
		keysController.Paths(),
		migrationsController.Paths(),
	)

	if err := vaultPlugin.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return vaultPlugin, nil
}
