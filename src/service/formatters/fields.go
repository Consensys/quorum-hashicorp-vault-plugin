package formatters

import "github.com/hashicorp/vault/sdk/framework"

const (
	PrivateKeyLabel          = "private_key"
	IDLabel                  = "id"
	DataLabel                = "data"
	NonceLabel               = "nonce"
	ToLabel                  = "to"
	AmountLabel              = "amount"
	GasPriceLabel            = "gas_price"
	GasLimitLabel            = "gas_limit"
	ChainIDLabel             = "chain_id"
	PrivateFromLabel         = "private_from"
	PrivateForLabel          = "private_for"
	PrivacyGroupIDLabel      = "privacy_group_id"
	TagsLabel                = "tags"
	CurveLabel               = "curve"
	AddressLabel             = "address"
	PublicKeyLabel           = "public_key"
	CompressedPublicKeyLabel = "compressed_public_key"
	NamespaceLabel           = "namespace"
	SignatureLabel           = "signature"
	AlgorithmLabel           = "signing_algorithm"
	VersionLabel             = "version"
	CreatedAtLabel           = "created_at"
	UpdatedAtLabel           = "updated_at"
	SourceNamespace          = "source_namespace"

	NamespaceHeader = "X-Vault-Namespace"
)

var IDFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeString,
	Description: "ID of the key pair",
	Required:    true,
}

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

var TagsFieldSchema = &framework.FieldSchema{
	Type:        framework.TypeKVPairs,
	Description: "Tags",
	Required:    true,
}
