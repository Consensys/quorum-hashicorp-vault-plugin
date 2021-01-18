package zksnarks

import (
	"context"
	"fmt"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	mocks2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignPayload_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks2.NewMockStorage(ctrl)
	mockGetAccountUC := mocks.NewMockGetZksAccountUseCase(ctrl)
	ctx := log.Context(context.Background(), log.Default())
	address := "0xaddress"
	namespace := "namespace"

	mockGetAccountUC.EXPECT().WithStorage(mockStorage).Return(mockGetAccountUC).AnyTimes()

	usecase := NewSignUseCase(mockGetAccountUC).WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		account := apputils.FakeZksAccount()
		account.PrivateKey = "0x10d10e9f17a88d51c42380c14da49e237b4b3f03c5cdce8f470ca782506eb5f113733d92b86e28b7e6354bb88a2d6bb9b104b0de3698b993f735f31cc979f7bd517ec47059e4830a2e75ca48845a9a2be571276aa654e3c3cea01abf7a092885511df98a80d9328b44d8bed35cb37e6aa1d3cda11420e77d2baf77c6172a1d98"
		
		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(account, nil)
		
		signature, err := usecase.Execute(ctx, address, namespace, "my data to sign")
		
		assert.NoError(t, err)
		assert.NotEmpty(t, signature)
	})

	t.Run("should fail with same error if Get Account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		
		mockGetAccountUC.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedErr)
		
		signature, err := usecase.Execute(ctx, address, namespace, "my data to sign")
		
		assert.Empty(t, signature)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail if creation of EDDSA private key fails", func(t *testing.T) {
		account := apputils.FakeZksAccount()
		account.PrivateKey = "account.PrivateKey"
		
		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(account, nil)
		
		signature, err := usecase.Execute(ctx, address, namespace, "my data to sign")
		
		assert.Empty(t, signature)
		assert.Error(t, err)
	})
}
