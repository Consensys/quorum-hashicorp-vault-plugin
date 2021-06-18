package keys

import (
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"net/http"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

func (s *keysCtrlTestSuite) TestKeysController_Sign() {
	path := s.controller.Paths()[3]
	signOperation := path.Operations[logical.CreateOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, fmt.Sprintf("keys/%s/sign", framework.GenericNameRegex(formatters.IDLabel)), path.Pattern)
		assert.NotEmpty(t, signOperation)
	})

	s.T().Run("should define correct properties", func(t *testing.T) {
		properties := signOperation.Properties()

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
		account := utils.FakeKey()
		payload := "my data to sign"
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				formatters.NamespaceHeader: {account.Namespace},
			},
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel:   account.ID,
				formatters.DataLabel: payload,
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel: formatters.AddressFieldSchema,
				formatters.DataLabel: {
					Type:        framework.TypeString,
					Description: "data to sign",
					Required:    true,
				},
			},
		}
		expectedSignature := "0x8b9679a75861e72fa6968dd5add3bf96e2747f0f124a2e728980f91e1958367e19c2486a40fdc65861824f247603bc18255fa497ca0b8b0a394aa7a6740fdc4601"

		s.signPayloadUC.EXPECT().Execute(gomock.Any(), account.ID, account.Namespace, payload).Return(expectedSignature, nil)

		response, err := signOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, expectedSignature, response.Data[formatters.SignatureLabel])
	})

	s.T().Run("should map errors correctly and return the correct http status", func(t *testing.T) {
		key := utils.FakeKey()
		payload := "my data to sign"
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel:   key.ID,
				formatters.DataLabel: payload,
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel: formatters.AddressFieldSchema,
				formatters.DataLabel: {
					Type:        framework.TypeString,
					Description: "data to sign",
					Required:    true,
				},
			},
		}
		expectedErr := errors.NotFoundError("error")

		s.signPayloadUC.EXPECT().Execute(gomock.Any(), key.ID, "", payload).Return("", expectedErr)

		response, err := signOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Data[logical.HTTPStatusCode])
	})
}
