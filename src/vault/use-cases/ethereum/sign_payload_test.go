package ethereum

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
	mockGetAccountUC := mocks.NewMockGetAccountUseCase(ctrl)
	ctx := log.Context(context.Background(), log.Default())
	address := "0xaddress"
	namespace := "namespace"

	mockGetAccountUC.EXPECT().WithStorage(mockStorage).Return(mockGetAccountUC).AnyTimes()

	usecase := NewSignUseCase(mockGetAccountUC).WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		account := apputils.FakeETHAccount()
		account.PrivateKey = "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"

		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(account, nil)

		signature, err := usecase.Execute(ctx, address, namespace, "my data to sign")

		assert.NoError(t, err)
		assert.Equal(t, signature, "0x7107193a8683e258ada2dfa76b5e6fc145ebd98f0e6eee77cb91381201fe7bca5445beccebe164e23abe0639f089e17b24ce867be9fece8b4872cfe13d91464601")
	})

	t.Run("should fail with same error if Get Account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockGetAccountUC.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, "my data to sign")

		assert.Empty(t, signature)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail if creation of ECDSA private key fails", func(t *testing.T) {
		account := apputils.FakeETHAccount()
		account.PrivateKey = "invalidPrivKey"

		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(account, nil)

		signature, err := usecase.Execute(ctx, address, namespace, "my data to sign")

		assert.Empty(t, signature)
		assert.Error(t, err)
	})
}
