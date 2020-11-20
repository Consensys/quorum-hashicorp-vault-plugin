package utils

import (
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/ethereum/go-ethereum/common"
)

func ExampleETHAccount() *entities.ETHAccount {
	return &entities.ETHAccount{
		Namespace:           "tenant0",
		Address:             common.HexToAddress("0x5962e04754cbE29a544f5104Ca100d94738Fd5d4").String(),
		PrivateKey:          "0b0232595b77568d99364bede133839ccbcb40775967a7eacd15d355c96288b5",
		PublicKey:           common.HexToHash("0x0433d7f005495fb6c0a34e22336dc3adcf4064553d5e194f77126bcac6da19491e0bab2772115cd284605d3bba94b69dc8c7a215021b58bcc87a70c9a440a3ff83").String(),
		CompressedPublicKey: common.HexToHash("0x0333d7f005495fb6c0a34e22336dc3adcf4064553d5e194f77126bcac6da19491e").String(),
	}
}

func Example500Response() framework.Response {
	return framework.Response{
		Description: "Internal server error",
		Example: &logical.Response{
			Data: map[string]interface{}{
				"error": "an unexpected error occurred. Please retry later or contact an administrator",
			},
		},
	}
}

func Example400Response() framework.Response {
	return framework.Response{
		Description: "Bad request",
		Example: &logical.Response{
			Data: map[string]interface{}{
				"error": "error message bad request",
			},
		},
	}
}

func Example404Response() framework.Response {
	return framework.Response{
		Description: "Not found",
		Example: &logical.Response{
			Data: map[string]interface{}{
				"error": "error message resource not found",
			},
		},
	}
}

func Example200Response() *framework.Response {
	return &framework.Response{
		Description: "Success",
		Example:     formatters.FormatAccountResponse(ExampleETHAccount()),
	}
}
