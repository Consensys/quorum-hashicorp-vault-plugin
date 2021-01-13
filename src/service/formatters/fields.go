package formatters

import "github.com/hashicorp/vault/sdk/framework"

const (
	PrivateKeyLabel     = "privateKey"
	AccountIDLabel      = "id"
	DataLabel           = "data"
	NonceLabel          = "nonce"
	ToLabel             = "to"
	AmountLabel         = "amount"
	GasPriceLabel       = "gasPrice"
	GasLimitLabel       = "gasLimit"
	ChainIDLabel        = "chainID"
	PrivateFromLabel    = "privateFrom"
	PrivateForLabel     = "privateFor"
	PrivacyGroupIDLabel = "privacyGroupID"

	NamespaceHeader = "X-Vault-Namespace"
)

var AddressFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeString,
	Description: "Address of the account",
	Required:    true,
}

var NonceFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeInt,
	Description: "Nonce of the transaction",
	Required:    true,
}

var ToFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeString,
	Description: "Recipient of the transaction. Empty for contract deployments",
}

var AmountFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeString,
	Description: "Amount of ETH (in wei) to transfer",
}

var GasPriceFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeString,
	Description: "The gas price for the transaction (in wei)",
	Required:    true,
}

var GasLimitFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeInt,
	Description: "The gas limit for the transaction",
	Required:    true,
}

var ChainIDFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeString,
	Description: "Network ID of the chain where the transaction will be deployed",
	Required:    true,
}

var DataFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeString,
	Description: "Data of the transaction",
}

var PrivateFromFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeString,
	Description: "EEA PrivateFrom address in base64 format",
}

var PrivateForFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeCommaStringSlice,
	Description: "EEA PrivateFor addresses in base64 format",
}

var PrivacyGroupIDFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeString,
	Description: "EEA PrivacyGroupID address in base64 format",
}
