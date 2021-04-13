package ethereum

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

func (s *ethereumCtrlTestSuite) TestEthereumController_Create() {
	path := s.controller.Paths()[0]
	createOperation := path.Operations[logical.CreateOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, "ethereum/accounts/?", path.Pattern)
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
		account := apputils.FakeETHAccount()
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				formatters.NamespaceHeader: {account.Namespace},
			},
		}

		s.createAccountUC.EXPECT().Execute(gomock.Any(), account.Namespace, "").Return(account, nil)

		response, err := createOperation.Handler()(s.ctx, request, &framework.FieldData{})

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
		expectedErr := errors.NotFoundError("error")

		s.createAccountUC.EXPECT().Execute(gomock.Any(), "", "").Return(nil, expectedErr)

		response, err := createOperation.Handler()(s.ctx, request, &framework.FieldData{})

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Data[logical.HTTPStatusCode])
	})
}
