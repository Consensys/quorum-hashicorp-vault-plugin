package formatters

import (
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/entities"
	"github.com/hashicorp/vault/sdk/logical"
)

func FormatAccountResponse(account *entities.ETHAccount) *logical.Response {
	return &logical.Response{
		Data: map[string]interface{}{
			"address":             account.Address,
			"publicKey":           account.PublicKey,
			"compressedPublicKey": account.CompressedPublicKey,
			"namespace":           account.Namespace,
		},
	}
}
