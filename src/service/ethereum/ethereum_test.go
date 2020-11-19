package ethereum

import (
	"context"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/testutils"
	ethereum "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/use-cases"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/use-cases/mocks"
	mocks2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/testutils/mocks"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ethereumCtrlTestSuite struct {
	suite.Suite
	createAccountUC *mocks.MockCreateAccountUseCase
	getAccountUC    *mocks.MockGetAccountUseCase
	listAccountsUC  *mocks.MockListAccountsUseCase
	storage         *mocks2.MockStorage
	ctx             context.Context
	controller      *controller
}

func (s *ethereumCtrlTestSuite) CreateAccount() ethereum.CreateAccountUseCase {
	return s.createAccountUC
}

func (s *ethereumCtrlTestSuite) GetAccount() ethereum.GetAccountUseCase {
	return s.getAccountUC
}

func (s *ethereumCtrlTestSuite) ListAccounts() ethereum.ListAccountsUseCase {
	return s.listAccountsUC
}

func (s *ethereumCtrlTestSuite) SignPayload() ethereum.SignUseCase {
	return nil
}

func (s *ethereumCtrlTestSuite) SignTransaction() ethereum.SignTransactionUseCase {
	return nil
}

func (s *ethereumCtrlTestSuite) SignQuorumPrivateTransaction() ethereum.SignQuorumPrivateTransactionUseCase {
	return nil
}

func (s *ethereumCtrlTestSuite) SignEEATransaction() ethereum.SignEEATransactionUseCase {
	return nil
}

var _ ethereum.UseCases = &ethereumCtrlTestSuite{}

func TestEthereumController(t *testing.T) {
	s := new(ethereumCtrlTestSuite)
	suite.Run(t, s)
}

func (s *ethereumCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.createAccountUC = mocks.NewMockCreateAccountUseCase(ctrl)
	s.getAccountUC = mocks.NewMockGetAccountUseCase(ctrl)
	s.listAccountsUC = mocks.NewMockListAccountsUseCase(ctrl)
	s.controller = NewController(s, hclog.Default())
	s.storage = mocks2.NewMockStorage(ctrl)
	s.ctx = context.Background()

	s.createAccountUC.EXPECT().WithStorage(s.storage).Return(s.createAccountUC).AnyTimes()
	s.getAccountUC.EXPECT().WithStorage(s.storage).Return(s.getAccountUC).AnyTimes()
	s.listAccountsUC.EXPECT().WithStorage(s.storage).Return(s.listAccountsUC).AnyTimes()
}

func (s *ethereumCtrlTestSuite) TestEthereumController_Create() {
	path := s.controller.Paths()[0]
	createOperation := path.Operations[logical.CreateOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, "ethereum/accounts", path.Pattern)
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
		account := testutils.FakeETHAccount()
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				namespaceHeader: {account.Namespace},
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

	s.T().Run("should return same error if use case fails", func(t *testing.T) {
		request := &logical.Request{
			Storage: s.storage,
		}
		expectedErr := fmt.Errorf("error")

		s.createAccountUC.EXPECT().Execute(gomock.Any(), "", "").Return(nil, expectedErr)

		response, err := createOperation.Handler()(s.ctx, request, &framework.FieldData{})

		assert.Empty(t, response)
		assert.Equal(t, expectedErr, err)
	})
}

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

func (s *ethereumCtrlTestSuite) TestEthereumController_Get() {
	path := s.controller.Paths()[2]
	getOperation := path.Operations[logical.ReadOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, fmt.Sprintf("ethereum/accounts/%s", framework.GenericNameRegex("address")), path.Pattern)
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
		account := testutils.FakeETHAccount()
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				namespaceHeader: {account.Namespace},
			},
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				addressLabel: account.Address,
			},
			Schema: map[string]*framework.FieldSchema{
				addressLabel: addressFieldSchema,
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

	s.T().Run("should return same error if use case fails", func(t *testing.T) {
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				addressLabel: "myAddress",
			},
			Schema: map[string]*framework.FieldSchema{
				addressLabel: addressFieldSchema,
			},
		}
		expectedErr := fmt.Errorf("error")

		s.getAccountUC.EXPECT().Execute(gomock.Any(), "myAddress", "").Return(nil, expectedErr)

		response, err := getOperation.Handler()(s.ctx, request, data)

		assert.Empty(t, response)
		assert.Equal(t, expectedErr, err)
	})
}

func (s *ethereumCtrlTestSuite) TestEthereumController_List() {
	path := s.controller.Paths()[0]
	listOperation := path.Operations[logical.ListOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, "ethereum/accounts", path.Pattern)
		assert.NotEmpty(t, listOperation)
	})

	s.T().Run("should define correct properties", func(t *testing.T) {
		properties := listOperation.Properties()

		assert.NotEmpty(t, properties.Description)
		assert.NotEmpty(t, properties.Summary)
		assert.NotEmpty(t, properties.Examples[0].Description)
		assert.NotEmpty(t, properties.Examples[0].Response)
		assert.NotEmpty(t, properties.Responses[200])
		assert.NotEmpty(t, properties.Responses[500])
	})

	s.T().Run("handler should execute the correct use case", func(t *testing.T) {
		account := testutils.FakeETHAccount()
		expectedList := []string{account.Address}
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				namespaceHeader: {account.Namespace},
			},
		}
		data := &framework.FieldData{}

		s.listAccountsUC.EXPECT().Execute(gomock.Any(), account.Namespace).Return(expectedList, nil)

		response, err := listOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, expectedList, response.Data["keys"])
	})

	s.T().Run("should return same error if use case fails", func(t *testing.T) {
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{}
		expectedErr := fmt.Errorf("error")

		s.listAccountsUC.EXPECT().Execute(gomock.Any(), "").Return(nil, expectedErr)

		response, err := listOperation.Handler()(s.ctx, request, data)

		assert.Empty(t, response)
		assert.Equal(t, expectedErr, err)
	})
}
