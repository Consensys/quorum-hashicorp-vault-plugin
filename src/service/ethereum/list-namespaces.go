package ethereum

import (
	"context"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/errors"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *controller) NewListNamespacesOperation() *framework.PathOperation {
	return &framework.PathOperation{
		Callback:    c.listNamespacesHandler(),
		Summary:     "Gets a list of all Ethereum namespaces",
		Description: "Gets a list of all Ethereum namespaces",
		Examples: []framework.RequestExample{
			{
				Description: "Gets all Ethereum namespaces",
				Response: &framework.Response{
					Description: "Success",
					Example:     logical.ListResponse([]string{"ns1", "ns2"}),
				},
			},
		},
		Responses: map[int][]framework.Response{
			200: {framework.Response{
				Description: "Success",
				Example:     logical.ListResponse([]string{"ns1", "ns2"}),
			}},
			500: {utils.Example500Response()},
		},
	}
}

func (c *controller) listNamespacesHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		ctx = log.Context(ctx, c.logger)
		namespaces, err := c.useCases.ListNamespaces().WithStorage(req.Storage).Execute(ctx)
		if err != nil {
			return errors.WriteHTTPError(req, err)
		}

		return logical.ListResponse(namespaces), nil
	}
}
