package keys

import (
	"context"
	errors2 "github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/errors"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/formatters"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *controller) NewUpdateOperation() *framework.PathOperation {
	exampleKey := utils.ExampleKey()
	successExample := utils.Example200KeysResponse()

	return &framework.PathOperation{
		Callback:    c.updateHandler(),
		Summary:     "Updates the tags of a key pair",
		Description: "Updates the tags of a key pair",
		Examples: []framework.RequestExample{
			{
				Description: "Updates the tags of a key pair",
				Data: map[string]interface{}{
					formatters.IDLabel:   exampleKey.ID,
					formatters.TagsLabel: exampleKey.Tags,
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

func (c *controller) updateHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		id := data.Get(formatters.IDLabel).(string)
		tags := data.Get(formatters.TagsLabel).(map[string]string)
		namespace := formatters.GetRequestNamespace(req)

		if tags == nil {
			return errors.ParseHTTPError(errors2.InvalidFormatError("tags must be provided"))
		}

		ctx = log.Context(ctx, c.logger)
		key, err := c.useCases.UpdateKey().WithStorage(req.Storage).Execute(ctx, namespace, id, tags)
		if err != nil {
			return errors.ParseHTTPError(err)
		}

		return formatters.FormatKeyResponse(key), nil
	}
}
