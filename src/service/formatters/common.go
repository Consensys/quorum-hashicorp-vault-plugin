package formatters

import "github.com/hashicorp/vault/sdk/logical"

func FormatSignatureResponse(signature string) *logical.Response {
	return &logical.Response{
		Data: map[string]interface{}{
			SignatureLabel: signature,
		},
	}
}
