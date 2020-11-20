package ethereum

import (
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type controller struct {
	useCases usecases.UseCases
	logger   hclog.Logger
}

func NewController(useCases usecases.UseCases, logger hclog.Logger) *controller {
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
			c.pathAccount(),
			c.pathSignPayload(),
		},
	)
}

func (c *controller) pathAccounts() *framework.Path {
	return &framework.Path{
		Pattern:      "ethereum/accounts",
		HelpSynopsis: "Creates a new Ethereum account",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: c.NewCreateOperation(),
			logical.UpdateOperation: c.NewCreateOperation(),
			logical.ListOperation:   c.NewListOperation(),
			logical.ReadOperation:   c.NewListOperation(),
		},
	}
}

func (c *controller) pathAccount() *framework.Path {
	return &framework.Path{
		Pattern:      fmt.Sprintf("ethereum/accounts/%s", framework.GenericNameRegex("address")),
		HelpSynopsis: "Get, update or delete an Ethereum account",
		Fields: map[string]*framework.FieldSchema{
			addressLabel: addressFieldSchema,
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: c.NewGetOperation(),
		},
		ExistenceCheck: c.ExistenceHandler,
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
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: c.NewImportOperation(),
			logical.UpdateOperation: c.NewImportOperation(),
		},
		HelpSynopsis: "Imports an Ethereum account",
	}
}

func (c *controller) pathSignPayload() *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("ethereum/accounts/%s/sign", framework.GenericNameRegex("address")),
		Fields: map[string]*framework.FieldSchema{
			addressLabel: addressFieldSchema,
			dataLabel: {
				Type:        framework.TypeString,
				Description: "data to sign",
				Required:    true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: c.NewSignPayloadOperation(),
			logical.UpdateOperation: c.NewSignPayloadOperation(),
		},
		HelpSynopsis: "Signs an arbitrary message using an existing Ethereum account",
	}
}

func getNamespace(req *logical.Request) string {
	namespace := ""

	if val, hasVal := req.Headers[namespaceHeader]; hasVal {
		namespace = val[0]
	}

	return namespace
}
