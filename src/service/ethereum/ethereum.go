package ethereum

import (
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
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
			c.pathSignTransaction(),
			c.pathSignQuorumPrivate(),
			c.pathSignEEA(),
			c.pathNamespaces(),
		},
	)
}

func (c *controller) pathAccounts() *framework.Path {
	return &framework.Path{
		Pattern:      "ethereum/accounts",
		HelpSynopsis: "Creates a new Ethereum account or list them",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: c.NewCreateOperation(),
			logical.UpdateOperation: c.NewCreateOperation(),
			logical.ListOperation:   c.NewListOperation(),
			logical.ReadOperation:   c.NewListOperation(),
		},
	}
}

func (c *controller) pathNamespaces() *framework.Path {
	return &framework.Path{
		Pattern:      "ethereum/namespaces",
		HelpSynopsis: "Lists all ethereum namespaces",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: c.NewListNamespacesOperation(),
			logical.ReadOperation: c.NewListNamespacesOperation(),
		},
	}
}

func (c *controller) pathAccount() *framework.Path {
	return &framework.Path{
		Pattern:      fmt.Sprintf("ethereum/accounts/%s", framework.GenericNameRegex("address")),
		HelpSynopsis: "Get, update or delete an Ethereum account",
		Fields: map[string]*framework.FieldSchema{
			formatters.AddressLabel: formatters.AddressFieldSchema,
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: c.NewGetOperation(),
		},
	}
}

func (c *controller) pathImportAccount() *framework.Path {
	return &framework.Path{
		Pattern: "ethereum/accounts/import",
		Fields: map[string]*framework.FieldSchema{
			formatters.PrivateKeyLabel: {
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
			formatters.AddressLabel: formatters.AddressFieldSchema,
			formatters.DataLabel: {
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

func (c *controller) pathSignTransaction() *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("ethereum/accounts/%s/sign-transaction", framework.GenericNameRegex("address")),
		Fields: map[string]*framework.FieldSchema{
			formatters.AddressLabel:  formatters.AddressFieldSchema,
			formatters.NonceLabel:    formatters.NonceFieldSchema,
			formatters.ToLabel:       formatters.ToFieldSchema,
			formatters.AmountLabel:   formatters.AmountFieldSchema,
			formatters.GasPriceLabel: formatters.GasPriceFieldSchema,
			formatters.GasLimitLabel: formatters.GasLimitFieldSchema,
			formatters.ChainIDLabel:  formatters.ChainIDFieldSchema,
			formatters.DataLabel:     formatters.DataFieldSchema,
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: c.NewSignTransactionOperation(),
			logical.UpdateOperation: c.NewSignTransactionOperation(),
		},
		HelpSynopsis: "Signs an Ethereum transaction using an existing account",
	}
}

func (c *controller) pathSignQuorumPrivate() *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("ethereum/accounts/%s/sign-quorum-private-transaction", framework.GenericNameRegex("address")),
		Fields: map[string]*framework.FieldSchema{
			formatters.AddressLabel:  formatters.AddressFieldSchema,
			formatters.NonceLabel:    formatters.NonceFieldSchema,
			formatters.ToLabel:       formatters.ToFieldSchema,
			formatters.AmountLabel:   formatters.AmountFieldSchema,
			formatters.GasPriceLabel: formatters.GasPriceFieldSchema,
			formatters.GasLimitLabel: formatters.GasLimitFieldSchema,
			formatters.DataLabel:     formatters.DataFieldSchema,
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: c.NewSignQuorumPrivateTransactionOperation(),
			logical.UpdateOperation: c.NewSignQuorumPrivateTransactionOperation(),
		},
		HelpSynopsis: "Signs a Quorum private transaction using an existing account",
	}
}

func (c *controller) pathSignEEA() *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("ethereum/accounts/%s/sign-eea-transaction", framework.GenericNameRegex("address")),
		Fields: map[string]*framework.FieldSchema{
			formatters.AddressLabel:        formatters.AddressFieldSchema,
			formatters.NonceLabel:          formatters.NonceFieldSchema,
			formatters.ToLabel:             formatters.ToFieldSchema,
			formatters.ChainIDLabel:        formatters.ChainIDFieldSchema,
			formatters.DataLabel:           formatters.DataFieldSchema,
			formatters.PrivateFromLabel:    formatters.PrivateFromFielSchema,
			formatters.PrivateForLabel:     formatters.PrivateForFielSchema,
			formatters.PrivacyGroupIDLabel: formatters.PrivacyGroupIDFielSchema,
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: c.NewSignEEATransactionOperation(),
			logical.UpdateOperation: c.NewSignEEATransactionOperation(),
		},
		HelpSynopsis: "Signs an EEA private transaction using an existing account",
	}
}

func getNamespace(req *logical.Request) string {
	namespace := ""

	if val, hasVal := req.Headers[formatters.NamespaceHeader]; hasVal {
		namespace = val[0]
	}

	return namespace
}
