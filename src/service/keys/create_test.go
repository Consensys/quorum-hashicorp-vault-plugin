package keys

import (
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"net/http"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

func (s *keysCtrlTestSuite) TestKeysController_Create() {
	path := s.controller.Paths()[0]
	createOperation := path.Operations[logical.CreateOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, "keys/?", path.Pattern)
		assert.NotEmpty(t, createOperation)
	})

	s.T().Run("should define correct properties", func(t *testing.T) {
		properties := createOperation.Properties()

		assert.NotEmpty(t, properties.Description)
		assert.NotEmpty(t, properties.Summary)
		assert.NotEmpty(t, properties.Examples[0].Description)
		assert.Empty(t, properties.Examples[0].Data)
		assert.NotEmpty(t, properties.Examples[0].Response)
		assert.NotEmpty(t, properties.Responses[200])
		assert.NotEmpty(t, properties.Responses[500])
	})

	s.T().Run("handler should execute the correct use case", func(t *testing.T) {
		key := apputils.FakeKey()
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				formatters.NamespaceHeader: {key.Namespace},
			},
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.CurveLabel: key.Curve,
				formatters.AlgoLabel:  key.Algorithm,
				formatters.IDLabel:    key.ID,
				formatters.TagsLabel:  key.Tags,
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel: formatters.IDFieldSchema,
				formatters.CurveLabel: {
					Type:        framework.TypeString,
					Description: "Elliptic curve",
					Required:    true,
				},
				formatters.AlgoLabel: {
					Type:        framework.TypeString,
					Description: "Signing algorithm",
					Required:    true,
				},
				formatters.TagsLabel: {
					Type:        framework.TypeKVPairs,
					Description: "Tags",
					Required:    true,
				},
			},
		}

		s.createKeyUC.EXPECT().Execute(gomock.Any(), key.Namespace, key.ID, key.Algorithm, key.Curve, "", key.Tags).Return(key, nil)

		response, err := createOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, key.PublicKey, response.Data["publicKey"])
		assert.Equal(t, key.Namespace, response.Data["namespace"])
		assert.Equal(t, key.Algorithm, response.Data["algorithm"])
		assert.Equal(t, key.Curve, response.Data["curve"])
		assert.Equal(t, key.ID, response.Data["id"])
		assert.Equal(t, key.Tags, response.Data["tags"])
	})

	s.T().Run("should map errors correctly and return the correct http status", func(t *testing.T) {
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.CurveLabel: "curve",
				formatters.AlgoLabel:  "algo",
				formatters.IDLabel:    "id",
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel: formatters.IDFieldSchema,
				formatters.CurveLabel: {
					Type:        framework.TypeString,
					Description: "Elliptic curve",
					Required:    true,
				},
				formatters.AlgoLabel: {
					Type:        framework.TypeString,
					Description: "Signing algorithm",
					Required:    true,
				},
				formatters.TagsLabel: {
					Type:        framework.TypeKVPairs,
					Description: "Tags",
					Required:    true,
				},
			},
		}
		expectedErr := errors.NotFoundError("error")

		s.createKeyUC.EXPECT().Execute(gomock.Any(), "", "id", "algo", "curve", "", map[string]string{}).Return(nil, expectedErr)

		response, err := createOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Data[logical.HTTPStatusCode])
	})
}
