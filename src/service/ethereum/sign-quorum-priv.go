package ethereum

import (
	"context"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *controller) NewSignQuorumPrivateTransactionOperation() *framework.PathOperation {
	exampleAccount := utils.ExampleETHAccount()

	return &framework.PathOperation{
		Callback:    c.signQuorumPrivateTransactionHandler(),
		Summary:     "Signs a Quorum private transaction using an existing account",
		Description: "Signs a Quorum private transaction using ECDSA and the private key of an existing account",
		Examples: []framework.RequestExample{
			{
				Description: "Signs a Quorum private transaction",
				Data: map[string]interface{}{
					formatters.AddressLabel:  exampleAccount.Address,
					formatters.NonceLabel:    0,
					formatters.ToLabel:       "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
					formatters.DataLabel:     "0xfeee...",
					formatters.AmountLabel:   "0",
					formatters.GasPriceLabel: "0",
					formatters.GasLimitLabel: 21000,
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

func (c *controller) signQuorumPrivateTransactionHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		address := data.Get(formatters.AddressLabel).(string)
		namespace := getNamespace(req)

		tx, err := formatters.FormatSignQuorumPrivateTransactionRequest(data)
		if err != nil {
			return nil, err
		}

		ctx = utils.WithLogger(ctx, c.logger)
		signature, err := c.useCases.SignQuorumPrivateTransaction().WithStorage(req.Storage).Execute(ctx, address, namespace, tx)
		if err != nil {
			return nil, err
		}

		return formatters.FormatSignatureResponse(signature), nil
	}
}
