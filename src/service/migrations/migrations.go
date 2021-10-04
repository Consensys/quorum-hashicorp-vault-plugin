package migrations

import (
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	usecases "github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type controller struct {
	useCases usecases.MigrationsUseCases
	logger   log.Logger
}

func NewController(useCases usecases.MigrationsUseCases, logger log.Logger) *controller {
	if logger == nil {
		logger = log.Default()
	}

	return &controller{
		useCases: useCases,
		logger:   logger.Named("migrations"),
	}
}

// Paths returns the list of paths
func (c *controller) Paths() []*framework.Path {
	return framework.PathAppend(
		[]*framework.Path{
			c.pathEthereumToKeys(),
			c.pathEthereumToKeysStatus(),
		},
	)
}

func (c *controller) pathEthereumToKeys() *framework.Path {
	return &framework.Path{
		Pattern:      "migrations/ethereum-to-keys/migrate/?",
		HelpSynopsis: "Migrates the current Ethereum accounts to the keys namespace",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: c.NewEthereumToKeysOperation(),
			logical.UpdateOperation: c.NewEthereumToKeysOperation(),
		},
	}
}

func (c *controller) pathEthereumToKeysStatus() *framework.Path {
	return &framework.Path{
		Pattern:      "migrations/ethereum-to-keys/status/?",
		HelpSynopsis: "Checks the status of the migration",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: c.NewEthereumToKeysStatusOperation(),
		},
	}
}
