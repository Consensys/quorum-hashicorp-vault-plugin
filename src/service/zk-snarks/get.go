package zksnarks

import (
	"context"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *controller) NewGetOperation() *framework.PathOperation {
	exampleAccount := utils.ExampleZksAccount()
	successExample := utils.Example200Response()

	return &framework.PathOperation{
		Callback:    c.getHandler(),
		Summary:     "Gets an zk-snarks account",
		Description: "Gets an zk-snarks account stored in the vault at the given address and namespace",
		Examples: []framework.RequestExample{
			{
				Description: "Gets an account on the tenant0 namespace",
				Data: map[string]interface{}{
					formatters.IDLabel: exampleAccount.PublicKey,
				},
				Response: successExample,
			},
		},
		Responses: map[int][]framework.Response{
			200: {*successExample},
			400: {utils.Example400Response()},
			404: {utils.Example404Response()},
			500: {utils.Example500Response()},
		},
	}
}

func (c *controller) getHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		accountID := data.Get(formatters.IDLabel).(string)
		namespace := formatters.GetRequestNamespace(req)

		ctx = log.Context(ctx, c.logger)
		account, err := c.useCases.GetAccount().WithStorage(req.Storage).Execute(ctx, accountID, namespace)
		if err != nil {
			return nil, err
		}

		return formatters.FormatZksAccountResponse(account), nil
	}
}
