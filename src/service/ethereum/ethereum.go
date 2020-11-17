package ethereum

import (
	ethereum "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/use-cases"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	namespaceLabel  = "namespace"
	privateKeyLabel = "privateKey"
)

var namespaceFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeString,
	Description: "Namespace in which to store the account",
	Required:    false,
	Default:     "",
}

type controller struct {
	useCases ethereum.UseCases
	logger   hclog.Logger
}

func NewController(useCases ethereum.UseCases, logger hclog.Logger) *controller {
	return &controller{
		useCases: useCases,
		logger:   logger,
	}
}

// Paths returns the list of paths
func (c *controller) Paths() []*framework.Path {
	return framework.PathAppend(
		[]*framework.Path{
			c.pathAccounts(),
			c.pathImportAccount(),
		},
	)
}

func (c *controller) pathAccounts() *framework.Path {
	return &framework.Path{
		Pattern:      "ethereum/accounts",
		HelpSynopsis: "Creates a new Ethereum account",
		Fields: map[string]*framework.FieldSchema{
			namespaceLabel: namespaceFieldSchema,
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: c.NewCreateOperation(),
			logical.UpdateOperation: c.NewCreateOperation(),
		},
	}
}

func (c *controller) pathImportAccount() *framework.Path {
	return &framework.Path{
		Pattern: "ethereum/accounts/import",
		Fields: map[string]*framework.FieldSchema{
			privateKeyLabel: {
				Type:        framework.TypeString,
				Description: "Private key in hexadecimal format",
				Required:    true,
			},
			namespaceLabel: namespaceFieldSchema,
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: c.NewImportOperation(),
			logical.UpdateOperation: c.NewImportOperation(),
		},
		HelpSynopsis: "Imports an Ethereum account",
	}
}
