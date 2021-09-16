package keys

import (
	"context"
	"testing"

	mocks2 "github.com/consensys/quorum-hashicorp-vault-plugin/src/utils/mocks"
	usecases "github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases/mocks"
	"github.com/hashicorp/go-hclog"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type keysCtrlTestSuite struct {
	suite.Suite
	createKeyUC      *mocks.MockCreateKeyUseCase
	getKeyUC         *mocks.MockGetKeyUseCase
	listKeysUC       *mocks.MockListKeysUseCase
	listNamespacesUC *mocks.MockListKeysNamespacesUseCase
	signPayloadUC    *mocks.MockKeysSignUseCase
	destroyKeyUC     *mocks.MockDestroyKeyUseCase
	updateKeyUC      *mocks.MockUpdateKeyUseCase
	storage          *mocks2.MockStorage
	ctx              context.Context
	controller       *controller
}

func (s *keysCtrlTestSuite) CreateKey() usecases.CreateKeyUseCase {
	return s.createKeyUC
}

func (s *keysCtrlTestSuite) GetKey() usecases.GetKeyUseCase {
	return s.getKeyUC
}

func (s *keysCtrlTestSuite) ListKeys() usecases.ListKeysUseCase {
	return s.listKeysUC
}

func (s *keysCtrlTestSuite) ListNamespaces() usecases.ListKeysNamespacesUseCase {
	return s.listNamespacesUC
}

func (s *keysCtrlTestSuite) SignPayload() usecases.KeysSignUseCase {
	return s.signPayloadUC
}

func (s *keysCtrlTestSuite) DestroyKey() usecases.DestroyKeyUseCase {
	return s.destroyKeyUC
}

func (s *keysCtrlTestSuite) UpdateKey() usecases.UpdateKeyUseCase {
	return s.updateKeyUC
}

var _ usecases.KeysUseCases = &keysCtrlTestSuite{}

func TestKeysController(t *testing.T) {
	s := new(keysCtrlTestSuite)
	suite.Run(t, s)
}

func (s *keysCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.createKeyUC = mocks.NewMockCreateKeyUseCase(ctrl)
	s.getKeyUC = mocks.NewMockGetKeyUseCase(ctrl)
	s.listKeysUC = mocks.NewMockListKeysUseCase(ctrl)
	s.listNamespacesUC = mocks.NewMockListKeysNamespacesUseCase(ctrl)
	s.signPayloadUC = mocks.NewMockKeysSignUseCase(ctrl)
	s.destroyKeyUC = mocks.NewMockDestroyKeyUseCase(ctrl)
	s.updateKeyUC = mocks.NewMockUpdateKeyUseCase(ctrl)
	s.controller = NewController(s, hclog.Default())
	s.storage = mocks2.NewMockStorage(ctrl)
	s.ctx = context.Background()

	s.createKeyUC.EXPECT().WithStorage(s.storage).Return(s.createKeyUC).AnyTimes()
	s.getKeyUC.EXPECT().WithStorage(s.storage).Return(s.getKeyUC).AnyTimes()
	s.listKeysUC.EXPECT().WithStorage(s.storage).Return(s.listKeysUC).AnyTimes()
	s.listNamespacesUC.EXPECT().WithStorage(s.storage).Return(s.listNamespacesUC).AnyTimes()
	s.signPayloadUC.EXPECT().WithStorage(s.storage).Return(s.signPayloadUC).AnyTimes()
	s.destroyKeyUC.EXPECT().WithStorage(s.storage).Return(s.destroyKeyUC).AnyTimes()
	s.updateKeyUC.EXPECT().WithStorage(s.storage).Return(s.updateKeyUC).AnyTimes()
}
