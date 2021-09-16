package formatters

import (
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/entities"
	"github.com/hashicorp/vault/sdk/logical"
)

func FormatZksAccountResponse(account *entities.ZksAccount) *logical.Response {
	return &logical.Response{
		Data: map[string]interface{}{
			CurveLabel:     account.Curve,
			AlgorithmLabel: account.Algorithm,
			PublicKeyLabel: account.PublicKey,
			NamespaceLabel: account.Namespace,
		},
	}
}
