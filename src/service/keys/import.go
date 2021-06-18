package keys

import (
	"context"
	errors2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/errors"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *controller) NewImportOperation() *framework.PathOperation {
	exampleKey := utils.ExampleKey()
	successExample := utils.Example200KeyResponse()

	return &framework.PathOperation{
		Callback:    c.importHandler(),
		Summary:     "Imports a key pair",
		Description: "Imports a key pair given a private key, storing it in the Vault and computing its public key and address",
		Examples: []framework.RequestExample{
			{
				Description: "Imports a key pair on the tenant0 namespace",
				Data: map[string]interface{}{
					formatters.PrivateKeyLabel: exampleKey.PrivateKey,
					formatters.CurveLabel:      exampleKey.Curve,
					formatters.AlgorithmLabel:  exampleKey.Algorithm,
					formatters.IDLabel:         exampleKey.ID,
					formatters.TagsLabel:       exampleKey.Tags,
				},
				Response: successExample,
			},
		},
		Responses: map[int][]framework.Response{
			200: {*successExample},
			400: {utils.Example400Response()},
			500: {utils.Example500Response()},
		},
	}
}

func (c *controller) importHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		namespace := formatters.GetRequestNamespace(req)
		id := data.Get(formatters.IDLabel).(string)
		curve := data.Get(formatters.CurveLabel).(string)
		algo := data.Get(formatters.AlgorithmLabel).(string)
		tags := data.Get(formatters.TagsLabel).(map[string]string)
		privateKeyString := data.Get(formatters.PrivateKeyLabel).(string)

		if id == "" {
			return errors.WriteHTTPError(req, errors2.InvalidFormatError("id must be provided"))
		}
		if curve == "" {
			return errors.WriteHTTPError(req, errors2.InvalidFormatError("curve must be provided"))
		}
		if algo == "" {
			return errors.WriteHTTPError(req, errors2.InvalidFormatError("algorithm must be provided"))
		}
		if privateKeyString == "" {
			return errors.WriteHTTPError(req, errors2.InvalidFormatError("privateKey must be provided"))
		}

		ctx = log.Context(ctx, c.logger)
		key, err := c.useCases.CreateKey().WithStorage(req.Storage).Execute(ctx, namespace, id, algo, curve, privateKeyString, tags)
		if err != nil {
			return errors.WriteHTTPError(req, err)
		}

		return formatters.FormatKeyResponse(key), nil
	}
}
