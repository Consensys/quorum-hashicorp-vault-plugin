package src

import (
	"context"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/ethereum"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/keys"
	zksnarks "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/zk-snarks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/builder"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// NewVaultBackend returns the Hashicorp Vault backend
func NewVaultBackend(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	vaultPlugin := &framework.Backend{
		Help: "Orchestrate Hashicorp Vault Plugin",
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"ethereum/accounts/",
				"zk-snarks/accounts/",
				"keys/",
			},
		},
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}

	ethereumController := ethereum.NewController(builder.NewEthereumUseCases(), conf.Logger)
	zkSnarksController := zksnarks.NewController(builder.NewZkSnarksUseCases(), conf.Logger)
	keysController := keys.NewController(builder.NewKeysUseCases(), conf.Logger)
	vaultPlugin.Paths = framework.PathAppend(ethereumController.Paths(), zkSnarksController.Paths(), keysController.Paths())

	if err := vaultPlugin.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return vaultPlugin, nil
}
