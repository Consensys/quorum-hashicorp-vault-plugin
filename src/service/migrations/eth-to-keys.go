package migrations

import (
	"context"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/errors"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/formatters"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *controller) NewEthereumToKeysOperation() *framework.PathOperation {
	return &framework.PathOperation{
		Callback:    c.ethToKeysHandler(),
		Summary:     "Migrates Ethereum accounts to the keys namespace",
		Description: "Migrates Ethereum accounts to the keys namespace by copying the data from the ethereum/accounts namespace",
		Examples: []framework.RequestExample{
			{
				Description: "Migrates the current Ethereum accounts to the keys namespace",
				Response:    nil,
			},
		},
		Responses: map[int][]framework.Response{
			204: {},
			500: {utils.Example500Response()},
		},
	}
}

func (c *controller) NewEthereumToKeysStatusOperation() *framework.PathOperation {
	return &framework.PathOperation{
		Callback:    c.ethToKeysStatusHandler(),
		Summary:     "Checks the status of the migration",
		Description: "Checks the status of the migration for the given namespace",
		Examples: []framework.RequestExample{
			{
				Description: "Checks the status of the migration",
				Response:    nil,
			},
		},
		Responses: map[int][]framework.Response{
			204: {},
			500: {utils.Example500Response()},
		},
	}
}

func (c *controller) ethToKeysHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		namespace := formatters.GetRequestNamespace(req)

		ctx = log.Context(ctx, c.logger)
		err := c.useCases.EthereumToKeys().Execute(ctx, req.Storage, namespace)
		if err != nil {
			return errors.ParseHTTPError(err)
		}

		return &logical.Response{}, nil
	}
}

func (c *controller) ethToKeysStatusHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		namespace := formatters.GetRequestNamespace(req)

		ctx = log.Context(ctx, c.logger)
		status, err := c.useCases.EthereumToKeys().Status(ctx, namespace)
		if err != nil {
			return errors.ParseHTTPError(err)
		}

		return formatters.FormatMigrationStatusResponse(status), nil
	}
}
