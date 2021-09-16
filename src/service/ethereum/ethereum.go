package ethereum

import (
	"fmt"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/formatters"
	usecases "github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type controller struct {
	useCases usecases.ETHUseCases
	logger   log.Logger
}

func NewController(useCases usecases.ETHUseCases, logger log.Logger) *controller {
	if logger == nil {
		logger = log.Default()
	}

	return &controller{
		useCases: useCases,
		logger:   logger.Named("ethereum"),
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
		Pattern:      "ethereum/accounts/?",
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
		Pattern:      "namespaces/ethereum/?",
		HelpSynopsis: "Lists all ethereum namespaces",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: c.NewListNamespacesOperation(),
			logical.ReadOperation: c.NewListNamespacesOperation(),
		},
	}
}

func (c *controller) pathAccount() *framework.Path {
	return &framework.Path{
		Pattern:      fmt.Sprintf("ethereum/accounts/%s", framework.GenericNameRegex(formatters.IDLabel)),
		HelpSynopsis: "Get, update or delete an Ethereum account",
		Fields: map[string]*framework.FieldSchema{
			formatters.IDLabel: formatters.AddressFieldSchema,
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
		Pattern: fmt.Sprintf("ethereum/accounts/%s/sign", framework.GenericNameRegex(formatters.IDLabel)),
		Fields: map[string]*framework.FieldSchema{
			formatters.IDLabel: formatters.AddressFieldSchema,
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
		Pattern: fmt.Sprintf("ethereum/accounts/%s/sign-transaction", framework.GenericNameRegex(formatters.IDLabel)),
		Fields: map[string]*framework.FieldSchema{
			formatters.IDLabel:       formatters.AddressFieldSchema,
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
		Pattern: fmt.Sprintf("ethereum/accounts/%s/sign-quorum-private-transaction", framework.GenericNameRegex(formatters.IDLabel)),
		Fields: map[string]*framework.FieldSchema{
			formatters.IDLabel:       formatters.AddressFieldSchema,
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
		Pattern: fmt.Sprintf("ethereum/accounts/%s/sign-eea-transaction", framework.GenericNameRegex(formatters.IDLabel)),
		Fields: map[string]*framework.FieldSchema{
			formatters.IDLabel:             formatters.AddressFieldSchema,
			formatters.NonceLabel:          formatters.NonceFieldSchema,
			formatters.ToLabel:             formatters.ToFieldSchema,
			formatters.ChainIDLabel:        formatters.ChainIDFieldSchema,
			formatters.DataLabel:           formatters.DataFieldSchema,
			formatters.PrivateFromLabel:    formatters.PrivateFromFieldSchema,
			formatters.PrivateForLabel:     formatters.PrivateForFieldSchema,
			formatters.PrivacyGroupIDLabel: formatters.PrivacyGroupIDFieldSchema,
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: c.NewSignEEATransactionOperation(),
			logical.UpdateOperation: c.NewSignEEATransactionOperation(),
		},
		HelpSynopsis: "Signs an EEA private transaction using an existing account",
	}
}
