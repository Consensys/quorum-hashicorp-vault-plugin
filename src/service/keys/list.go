package keys

import (
	"context"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/errors"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *controller) NewListOperation() *framework.PathOperation {
	return &framework.PathOperation{
		Callback:    c.listHandler(),
		Summary:     "Gets a list of all key pair ids",
		Description: "Gets a list of all key pair ids optionally filtered by namespace",
		Examples: []framework.RequestExample{
			{
				Description: "Gets all key pair ids",
				Response: &framework.Response{
					Description: "Success",
					Example:     logical.ListResponse([]string{utils.ExampleKey().ID}),
				},
			},
		},
		Responses: map[int][]framework.Response{
			200: {framework.Response{
				Description: "Success",
				Example:     logical.ListResponse([]string{utils.ExampleKey().ID}),
			}},
			500: {utils.Example500Response()},
		},
	}
}

func (c *controller) listHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		namespace := formatters.GetRequestNamespace(req)

		ctx = log.Context(ctx, c.logger)
		keys, err := c.useCases.ListKeys().WithStorage(req.Storage).Execute(ctx, namespace)
		if err != nil {
			return errors.WriteHTTPError(req, err)
		}

		return logical.ListResponse(keys), nil
	}
}
