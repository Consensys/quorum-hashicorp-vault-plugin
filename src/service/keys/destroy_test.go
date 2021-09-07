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
	"github.com/stretchr/testify/require"
	"testing"
)

func (s *keysCtrlTestSuite) TestKeysController_Destroy() {
	path := s.controller.Paths()[5]
	deleteOperation := path.Operations[logical.DeleteOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, fmt.Sprintf("keys/%s/destroy", framework.GenericNameRegex(formatters.IDLabel)), path.Pattern)
		assert.NotEmpty(t, deleteOperation)
	})

	s.T().Run("should define correct properties", func(t *testing.T) {
		properties := deleteOperation.Properties()

		assert.NotEmpty(t, properties.Description)
		assert.NotEmpty(t, properties.Summary)
		assert.NotEmpty(t, properties.Examples[0].Description)
		assert.NotEmpty(t, properties.Examples[0].Data)
		assert.Empty(t, properties.Responses[200])
		assert.NotEmpty(t, properties.Responses[404])
		assert.NotEmpty(t, properties.Responses[500])
	})

	s.T().Run("handler should execute the correct use case", func(t *testing.T) {
		account := utils.FakeKey()
		tags := map[string]string{
			"tag1": "tagValue1",
			"tag2": "tagValue2",
		}
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				formatters.NamespaceHeader: {account.Namespace},
			},
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel:   account.ID,
				formatters.TagsLabel: tags,
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel:   formatters.AddressFieldSchema,
				formatters.TagsLabel: formatters.TagsFieldSchema,
			},
		}

		s.destroyKeyUC.EXPECT().Execute(gomock.Any(), account.Namespace, account.ID).Return(nil)

		response, err := deleteOperation.Handler()(s.ctx, request, data)
		require.NoError(t, err)

		assert.Nil(t, response.Data)
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

		s.destroyKeyUC.EXPECT().Execute(gomock.Any(), "", key.ID).Return(expectedErr)

		_, err := deleteOperation.Handler()(s.ctx, request, data)

		assert.Equal(t, err, logical.ErrUnsupportedPath)
	})
}
