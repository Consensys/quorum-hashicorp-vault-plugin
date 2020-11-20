package ethereum

import (
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/testutils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"testing"
)

func (s *ethereumCtrlTestSuite) TestEthereumController_Import() {
	path := s.controller.Paths()[1]
	importOperation := path.Operations[logical.CreateOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, "ethereum/accounts/import", path.Pattern)
		assert.NotEmpty(t, importOperation)
	})

	s.T().Run("should define correct properties", func(t *testing.T) {
		properties := importOperation.Properties()

		assert.NotEmpty(t, properties.Description)
		assert.NotEmpty(t, properties.Summary)
		assert.NotEmpty(t, properties.Examples[0].Description)
		assert.NotEmpty(t, properties.Examples[0].Data)
		assert.NotEmpty(t, properties.Examples[0].Response)
		assert.NotEmpty(t, properties.Responses[200])
		assert.NotEmpty(t, properties.Responses[400])
		assert.NotEmpty(t, properties.Responses[500])
	})

	s.T().Run("handler should execute the correct use case", func(t *testing.T) {
		account := testutils.FakeETHAccount()
		privKey := "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				namespaceHeader: {account.Namespace},
			},
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				privateKeyLabel: privKey,
			},
			Schema: map[string]*framework.FieldSchema{
				privateKeyLabel: {
					Type:        framework.TypeString,
					Description: "Private key in hexadecimal format",
					Required:    true,
				},
			},
		}

		s.createAccountUC.EXPECT().Execute(gomock.Any(), account.Namespace, privKey).Return(account, nil)

		response, err := importOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, account.Address, response.Data["address"])
		assert.Equal(t, account.PublicKey, response.Data["publicKey"])
		assert.Equal(t, account.CompressedPublicKey, response.Data["compressedPublicKey"])
		assert.Equal(t, account.Namespace, response.Data["namespace"])
	})

	s.T().Run("should return same error if use case fails", func(t *testing.T) {
		privKey := "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				privateKeyLabel: privKey,
			},
			Schema: map[string]*framework.FieldSchema{
				privateKeyLabel: {
					Type:        framework.TypeString,
					Description: "Private key in hexadecimal format",
					Required:    true,
				},
			},
		}
		expectedErr := fmt.Errorf("error")

		s.createAccountUC.EXPECT().Execute(gomock.Any(), "", privKey).Return(nil, expectedErr)

		response, err := importOperation.Handler()(s.ctx, request, data)

		assert.Empty(t, response)
		assert.Equal(t, expectedErr, err)
	})
}
