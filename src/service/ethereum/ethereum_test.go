package ethereum

import (
	"context"
	mocks2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/testutils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/mocks"
	"github.com/hashicorp/go-hclog"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ethereumCtrlTestSuite struct {
	suite.Suite
	createAccountUC *mocks.MockCreateAccountUseCase
	getAccountUC    *mocks.MockGetAccountUseCase
	listAccountsUC  *mocks.MockListAccountsUseCase
	signPayloadUC   *mocks.MockSignUseCase
	storage         *mocks2.MockStorage
	ctx             context.Context
	controller      *controller
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

func (s *ethereumCtrlTestSuite) SignPayload() usecases.SignUseCase {
	return s.signPayloadUC
}

func (s *ethereumCtrlTestSuite) SignTransaction() usecases.SignTransactionUseCase {
	return nil
}

func (s *ethereumCtrlTestSuite) SignQuorumPrivateTransaction() usecases.SignQuorumPrivateTransactionUseCase {
	return nil
}

func (s *ethereumCtrlTestSuite) SignEEATransaction() usecases.SignEEATransactionUseCase {
	return nil
}

var _ usecases.UseCases = &ethereumCtrlTestSuite{}

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
	s.signPayloadUC = mocks.NewMockSignUseCase(ctrl)
	s.controller = NewController(s, hclog.Default())
	s.storage = mocks2.NewMockStorage(ctrl)
	s.ctx = context.Background()

	s.createAccountUC.EXPECT().WithStorage(s.storage).Return(s.createAccountUC).AnyTimes()
	s.getAccountUC.EXPECT().WithStorage(s.storage).Return(s.getAccountUC).AnyTimes()
	s.listAccountsUC.EXPECT().WithStorage(s.storage).Return(s.listAccountsUC).AnyTimes()
	s.signPayloadUC.EXPECT().WithStorage(s.storage).Return(s.signPayloadUC).AnyTimes()
}
