package zksnarks

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

type zksCtrlTestSuite struct {
	suite.Suite
	createAccountUC                *mocks.MockCreateZksAccountUseCase
	getAccountUC                   *mocks.MockGetZksAccountUseCase
	listAccountsUC                 *mocks.MockListZksAccountsUseCase
	listNamespacesUC               *mocks.MockListZksNamespacesUseCase
	signPayloadUC                  *mocks.MockZksSignUseCase
	storage                        *mocks2.MockStorage
	ctx                            context.Context
	controller                     *controller
}

func (s *zksCtrlTestSuite) CreateAccount() usecases.CreateZksAccountUseCase {
	return s.createAccountUC
}

func (s *zksCtrlTestSuite) GetAccount() usecases.GetZksAccountUseCase {
	return s.getAccountUC
}

func (s *zksCtrlTestSuite) ListAccounts() usecases.ListZksAccountsUseCase {
	return s.listAccountsUC
}

func (s *zksCtrlTestSuite) ListNamespaces() usecases.ListZksNamespacesUseCase {
	return s.listNamespacesUC
}

func (s *zksCtrlTestSuite) SignPayload() usecases.ZksSignUseCase {
	return s.signPayloadUC
}

var _ usecases.ZksUseCases = &zksCtrlTestSuite{}

func TestZksController(t *testing.T) {
	s := new(zksCtrlTestSuite)
	suite.Run(t, s)
}

func (s *zksCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.createAccountUC = mocks.NewMockCreateZksAccountUseCase(ctrl)
	s.getAccountUC = mocks.NewMockGetZksAccountUseCase(ctrl)
	s.listAccountsUC = mocks.NewMockListZksAccountsUseCase(ctrl)
	s.listNamespacesUC = mocks.NewMockListZksNamespacesUseCase(ctrl)
	s.signPayloadUC = mocks.NewMockZksSignUseCase(ctrl)
	s.controller = NewController(s, hclog.Default())
	s.storage = mocks2.NewMockStorage(ctrl)
	s.ctx = context.Background()

	s.createAccountUC.EXPECT().WithStorage(s.storage).Return(s.createAccountUC).AnyTimes()
	s.getAccountUC.EXPECT().WithStorage(s.storage).Return(s.getAccountUC).AnyTimes()
	s.listAccountsUC.EXPECT().WithStorage(s.storage).Return(s.listAccountsUC).AnyTimes()
	s.listNamespacesUC.EXPECT().WithStorage(s.storage).Return(s.listNamespacesUC).AnyTimes()
	s.signPayloadUC.EXPECT().WithStorage(s.storage).Return(s.signPayloadUC).AnyTimes()
}
