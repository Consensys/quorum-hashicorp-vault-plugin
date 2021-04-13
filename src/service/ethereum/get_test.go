package ethereum

import (
	"fmt"
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

func (s *ethereumCtrlTestSuite) TestEthereumController_Get() {
	path := s.controller.Paths()[2]
	getOperation := path.Operations[logical.ReadOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, fmt.Sprintf("ethereum/accounts/%s", framework.GenericNameRegex(formatters.IDLabel)), path.Pattern)
		assert.NotEmpty(t, getOperation)
	})

	s.T().Run("should define correct properties", func(t *testing.T) {
		properties := getOperation.Properties()

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
		account := apputils.FakeETHAccount()
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				formatters.NamespaceHeader: {account.Namespace},
			},
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel: account.Address,
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel: formatters.AddressFieldSchema,
			},
		}

		s.getAccountUC.EXPECT().Execute(gomock.Any(), account.Address, account.Namespace).Return(account, nil)

		response, err := getOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, account.Address, response.Data["address"])
		assert.Equal(t, account.PublicKey, response.Data["publicKey"])
		assert.Equal(t, account.CompressedPublicKey, response.Data["compressedPublicKey"])
		assert.Equal(t, account.Namespace, response.Data["namespace"])
	})

	s.T().Run("should map errors correctly and return the correct http status", func(t *testing.T) {
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel: "myAddress",
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel: formatters.AddressFieldSchema,
			},
		}
		expectedErr := errors.NotFoundError("error")

		s.getAccountUC.EXPECT().Execute(gomock.Any(), "myAddress", "").Return(nil, expectedErr)

		response, err := getOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Data[logical.HTTPStatusCode])
	})
}
