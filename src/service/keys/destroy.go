package keys

import (
	"context"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/errors"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/formatters"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *controller) NewDestroyOperation() *framework.PathOperation {
	return &framework.PathOperation{
		Callback:    c.destroyHandler(),
		Summary:     "Destroys an existing key pair",
		Description: "Destroys an existing key pair. The key will not be recoverable after this operation is performed",
		Examples: []framework.RequestExample{
			{
				Description: "Destroys an existing key pair",
				Data: map[string]interface{}{
					formatters.IDLabel: "my-key",
				},
			},
		},
		Responses: map[int][]framework.Response{
			204: {},
			404: {utils.Example404Response()},
			500: {utils.Example500Response()},
		},
	}
}

func (c *controller) destroyHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		id := data.Get(formatters.IDLabel).(string)
		namespace := formatters.GetRequestNamespace(req)

		ctx = log.Context(ctx, c.logger)
		err := c.useCases.DestroyKey().WithStorage(req.Storage).Execute(ctx, namespace, id)
		if err != nil {
			return errors.ParseHTTPError(err)
		}

		return &logical.Response{}, nil
	}
}
