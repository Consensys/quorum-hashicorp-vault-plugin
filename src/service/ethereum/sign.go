package ethereum

import (
	"context"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *controller) NewSignPayloadOperation() *framework.PathOperation {
	exampleAccount := utils.ExampleETHAccount()
	exampleSignature := "0x8b9679a75861e72fa6968dd5add3bf96e2747f0f124a2e728980f91e1958367e19c2486a40fdc65861824f247603bc18255fa497ca0b8b0a394aa7a6740fdc4601"

	return &framework.PathOperation{
		Callback:    c.signPayloadHandler(),
		Summary:     "Signs an arbitrary message using an existing Ethereum account",
		Description: "Signs an arbitrary message using ECDSA and the private key of an existing Ethereum account",
		Examples: []framework.RequestExample{
			{
				Description: "Signs a message",
				Data: map[string]interface{}{
					addressLabel: exampleAccount.Address,
					dataLabel:    "my data to sign",
				},
				Response: &framework.Response{
					Description: "Success",
					Example:     formatters.FormatSignatureResponse(exampleSignature),
				},
			},
		},
		Responses: map[int][]framework.Response{
			200: {framework.Response{
				Description: "Success",
				Example:     formatters.FormatSignatureResponse(exampleSignature),
			}},
			400: {utils.Example400Response()},
			500: {utils.Example500Response()},
		},
	}
}

func (c *controller) signPayloadHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		address := data.Get("address").(string)
		payload := data.Get("data").(string)
		namespace := getNamespace(req)

		if payload == "" {
			return logical.ErrorResponse("data must be provided"), nil
		}

		ctx = utils.WithLogger(ctx, c.logger)
		signature, err := c.useCases.SignPayload().WithStorage(req.Storage).Execute(ctx, address, namespace, payload)
		if err != nil {
			return nil, err
		}

		return formatters.FormatSignatureResponse(signature), nil
	}
}
