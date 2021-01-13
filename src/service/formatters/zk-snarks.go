package formatters

import (
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/hashicorp/vault/sdk/logical"
)

func FormatZksAccountResponse(account *entities.ZksAccount) *logical.Response {
	return &logical.Response{
		Data: map[string]interface{}{
			"curve":            account.Curve,
			"signingAlgorithm": account.Algorithm,
			"publicKey":        account.PublicKey,
			"namespace":        account.Namespace,
		},
	}
}
