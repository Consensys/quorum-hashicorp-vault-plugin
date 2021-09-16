package keys

import (
	"context"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/errors"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/formatters"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *controller) NewGetOperation() *framework.PathOperation {
	exampleKey := utils.ExampleKey()
	successExample := utils.Example200KeyResponse()

	return &framework.PathOperation{
		Callback:    c.getHandler(),
		Summary:     "Gets a key pair",
		Description: "Gets a key pair stored in the vault at the given id and namespace",
		Examples: []framework.RequestExample{
			{
				Description: "Gets a key pair on the tenant0 namespace",
				Data: map[string]interface{}{
					formatters.IDLabel: exampleKey.ID,
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
		id := data.Get(formatters.IDLabel).(string)
		namespace := formatters.GetRequestNamespace(req)

		ctx = log.Context(ctx, c.logger)
		key, err := c.useCases.GetKey().WithStorage(req.Storage).Execute(ctx, id, namespace)
		if err != nil {
			return errors.ParseHTTPError(err)
		}

		return formatters.FormatKeyResponse(key), nil
	}
}
