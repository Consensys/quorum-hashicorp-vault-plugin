package ethereum

import (
	"context"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/testutils"
	mocks2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/testutils/mocks"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-hclog"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks2.NewMockStorage(ctrl)
	ctx := apputils.WithLogger(context.Background(), hclog.New(&hclog.LoggerOptions{}))

	usecase := NewCreateAccountUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully by generating a private key", func(t *testing.T) {
		fakeAccount := testutils.FakeETHAccount()
		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		account, err := usecase.Execute(ctx, fakeAccount.Namespace, "")

		assert.NoError(t, err)
		assert.Equal(t, fakeAccount.Namespace, account.Namespace)
		assert.True(t, common.IsHexAddress(account.Address))
	})

	t.Run("should execute use case successfully by importing a private key", func(t *testing.T) {
		privKey := "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"

		fakeAccount := testutils.FakeETHAccount()
		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(nil)

		account, err := usecase.Execute(ctx, fakeAccount.Namespace, privKey)

		assert.NoError(t, err)
		assert.Equal(t, fakeAccount.Namespace, account.Namespace)
		assert.Equal(t, "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5", account.Address)
	})

	t.Run("should fail with same error if Put fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockStorage.EXPECT().Put(ctx, gomock.Any()).Return(expectedErr)

		account, err := usecase.Execute(ctx, "namespace", "")
		assert.Nil(t, account)
		assert.Equal(t, expectedErr, err)
	})
}
