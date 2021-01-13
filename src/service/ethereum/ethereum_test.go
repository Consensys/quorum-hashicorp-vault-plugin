package ethereum

import (
	"context"
	mocks2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/mocks"
	"github.com/hashicorp/go-hclog"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ethereumCtrlTestSuite struct {
	suite.Suite
	createAccountUC                *mocks.MockCreateAccountUseCase
	getAccountUC                   *mocks.MockGetAccountUseCase
	listAccountsUC                 *mocks.MockListAccountsUseCase
	listNamespacesUC               *mocks.MockListNamespacesUseCase
	signPayloadUC                  *mocks.MockSignUseCase
	signTransactionUC              *mocks.MockSignTransactionUseCase
	signQuorumPrivateTransactionUC *mocks.MockSignQuorumPrivateTransactionUseCase
	signEEATransactionUC           *mocks.MockSignEEATransactionUseCase
	storage                        *mocks2.MockStorage
	ctx                            context.Context
	controller                     *controller
}

func (s *ethereumCtrlTestSuite) CreateAccount() usecases.CreateAccountUseCase {
	return s.createAccountUC
}

func (s *ethereumCtrlTestSuite) GetAccount() usecases.GetAccountUseCase {
	return s.getAccountUC
}

func (s *ethereumCtrlTestSuite) ListAccounts() usecases.ListAccountsUseCase {
	return s.listAccountsUC
}

func (s *ethereumCtrlTestSuite) ListNamespaces() usecases.ListNamespacesUseCase {
	return s.listNamespacesUC
}

func (s *ethereumCtrlTestSuite) SignPayload() usecases.SignUseCase {
	return s.signPayloadUC
}

func (s *ethereumCtrlTestSuite) SignTransaction() usecases.SignTransactionUseCase {
	return s.signTransactionUC
}

func (s *ethereumCtrlTestSuite) SignQuorumPrivateTransaction() usecases.SignQuorumPrivateTransactionUseCase {
	return s.signQuorumPrivateTransactionUC
}

func (s *ethereumCtrlTestSuite) SignEEATransaction() usecases.SignEEATransactionUseCase {
	return s.signEEATransactionUC
}

var _ usecases.ETHUseCases = &ethereumCtrlTestSuite{}

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
	s.listNamespacesUC = mocks.NewMockListNamespacesUseCase(ctrl)
	s.signPayloadUC = mocks.NewMockSignUseCase(ctrl)
	s.signTransactionUC = mocks.NewMockSignTransactionUseCase(ctrl)
	s.signQuorumPrivateTransactionUC = mocks.NewMockSignQuorumPrivateTransactionUseCase(ctrl)
	s.signEEATransactionUC = mocks.NewMockSignEEATransactionUseCase(ctrl)
	s.controller = NewController(s, hclog.Default())
	s.storage = mocks2.NewMockStorage(ctrl)
	s.ctx = context.Background()

	s.createAccountUC.EXPECT().WithStorage(s.storage).Return(s.createAccountUC).AnyTimes()
	s.getAccountUC.EXPECT().WithStorage(s.storage).Return(s.getAccountUC).AnyTimes()
	s.listAccountsUC.EXPECT().WithStorage(s.storage).Return(s.listAccountsUC).AnyTimes()
	s.listNamespacesUC.EXPECT().WithStorage(s.storage).Return(s.listNamespacesUC).AnyTimes()
	s.signPayloadUC.EXPECT().WithStorage(s.storage).Return(s.signPayloadUC).AnyTimes()
	s.signTransactionUC.EXPECT().WithStorage(s.storage).Return(s.signTransactionUC).AnyTimes()
	s.signQuorumPrivateTransactionUC.EXPECT().WithStorage(s.storage).Return(s.signQuorumPrivateTransactionUC).AnyTimes()
	s.signEEATransactionUC.EXPECT().WithStorage(s.storage).Return(s.signEEATransactionUC).AnyTimes()
}
