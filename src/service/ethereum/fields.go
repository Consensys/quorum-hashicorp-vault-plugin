package ethereum

import "github.com/hashicorp/vault/sdk/framework"

const (
	privateKeyLabel = "privateKey"
	addressLabel    = "address"
	dataLabel       = "data"

	namespaceHeader = "X-Vault-Namespace"
)

var addressFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeString,
	Description: "Address of the account",
	Required:    true,
}
