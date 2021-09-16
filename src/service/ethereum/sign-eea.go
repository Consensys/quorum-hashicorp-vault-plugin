package ethereum

import (
	"context"
	errors2 "github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/errors"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/formatters"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *controller) NewSignEEATransactionOperation() *framework.PathOperation {
	exampleAccount := utils.ExampleETHAccount()

	return &framework.PathOperation{
		Callback:    c.signEEATransactionHandler(),
		Summary:     "Signs an EEA private transaction using an existing account",
		Description: "Signs an EEA private transaction using ECDSA and the private key of an existing account",
		Examples: []framework.RequestExample{
			{
				Description: "Signs an EEA transaction",
				Data: map[string]interface{}{
					formatters.IDLabel:          exampleAccount.Address,
					formatters.NonceLabel:       0,
					formatters.ToLabel:          "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
					formatters.ChainIDLabel:     "1",
					formatters.DataLabel:        "0xfeee...",
					formatters.PrivateFromLabel: "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
					formatters.PrivateForLabel:  []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=", "B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="},
				},
				Response: utils.Example200ResponseSignature(),
			},
		},
		Responses: map[int][]framework.Response{
			200: {*utils.Example200ResponseSignature()},
			400: {utils.Example400Response()},
			404: {utils.Example404Response()},
			500: {utils.Example500Response()},
		},
	}
}

func (c *controller) signEEATransactionHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		address := data.Get(formatters.IDLabel).(string)
		chainID := data.Get(formatters.ChainIDLabel).(string)
		namespace := formatters.GetRequestNamespace(req)

		if chainID == "" {
			return errors.ParseHTTPError(errors2.InvalidFormatError("chainID must be provided"))
		}

		tx, privateArgs, err := formatters.FormatSignEEATransactionRequest(data)
		if err != nil {
			return errors.ParseHTTPError(err)
		}

		ctx = log.Context(ctx, c.logger)
		signature, err := c.useCases.SignEEATransaction().WithStorage(req.Storage).Execute(ctx, address, namespace, chainID, tx, privateArgs)
		if err != nil {
			return errors.ParseHTTPError(err)
		}

		return formatters.FormatSignatureResponse(signature), nil
	}
}
