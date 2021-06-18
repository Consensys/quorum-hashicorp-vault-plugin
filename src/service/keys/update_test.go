package keys

import (
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func (s *keysCtrlTestSuite) TestKeysController_Update() {
	path := s.controller.Paths()[2]
	updateOperation := path.Operations[logical.UpdateOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, fmt.Sprintf("keys/%s", framework.GenericNameRegex(formatters.IDLabel)), path.Pattern)
		assert.NotEmpty(t, updateOperation)
	})

	s.T().Run("should define correct properties", func(t *testing.T) {
		properties := updateOperation.Properties()

		assert.NotEmpty(t, properties.Description)
		assert.NotEmpty(t, properties.Summary)
		assert.NotEmpty(t, properties.Examples[0].Description)
		assert.NotEmpty(t, properties.Examples[0].Response)
		assert.NotEmpty(t, properties.Examples[0].Data)
		assert.NotEmpty(t, properties.Responses[200])
		assert.NotEmpty(t, properties.Responses[400])
		assert.NotEmpty(t, properties.Responses[404])
		assert.NotEmpty(t, properties.Responses[500])
	})

	s.T().Run("handler should execute the correct use case", func(t *testing.T) {
		key := utils.FakeKey()
		tags := map[string]string{
			"tag1": "tagValue1",
			"tag2": "tagValue2",
		}
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				formatters.NamespaceHeader: {key.Namespace},
			},
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel:   key.ID,
				formatters.TagsLabel: tags,
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel:   formatters.AddressFieldSchema,
				formatters.TagsLabel: formatters.TagsFieldSchema,
			},
		}

		s.updateKeyUC.EXPECT().Execute(gomock.Any(), key.Namespace, key.ID, key.Tags).Return(key, nil)

		response, err := updateOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, key.PublicKey, response.Data[formatters.PublicKeyLabel])
		assert.Equal(t, key.Namespace, response.Data[formatters.NamespaceLabel])
		assert.Equal(t, key.Algorithm, response.Data[formatters.AlgorithmLabel])
		assert.Equal(t, key.Curve, response.Data[formatters.CurveLabel])
		assert.Equal(t, key.ID, response.Data[formatters.IDLabel])
		assert.Equal(t, key.Tags, response.Data[formatters.TagsLabel])
	})

	s.T().Run("should map errors correctly and return the correct http status", func(t *testing.T) {
		key := utils.FakeKey()
		tags := map[string]string{
			"tag1": "tagValue1",
			"tag2": "tagValue2",
		}
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel:   key.ID,
				formatters.TagsLabel: tags,
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel:   formatters.AddressFieldSchema,
				formatters.TagsLabel: formatters.TagsFieldSchema,
			},
		}
		expectedErr := errors.NotFoundError("error")

		s.updateKeyUC.EXPECT().Execute(gomock.Any(), "", key.ID, key.Tags).Return(nil, expectedErr)

		response, err := updateOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Data[logical.HTTPStatusCode])
	})
}
