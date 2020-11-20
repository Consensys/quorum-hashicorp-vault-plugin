package src

import (
	"context"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/ethereum"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
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
			},
		},
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}

	ethereumController := ethereum.NewController(builder.NewEthereumUseCases(), vaultPlugin.Logger())
	vaultPlugin.Paths = ethereumController.Paths()

	ctx = utils.WithLogger(ctx, vaultPlugin.Logger())
	if err := vaultPlugin.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return vaultPlugin, nil
}
