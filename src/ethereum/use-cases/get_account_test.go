package ethereum

import (
	"context"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/testutils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/use-cases/utils"
	mocks2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/testutils/mocks"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAccount_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks2.NewMockStorage(ctrl)
	ctx := apputils.WithLogger(context.Background(), hclog.New(&hclog.LoggerOptions{}))

	usecase := NewGetAccountUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		fakeAccount := testutils.FakeETHAccount()
		expectedEntry, _ := logical.StorageEntryJSON(utils.ComputeKey(fakeAccount.Address, fakeAccount.Namespace), fakeAccount)
		mockStorage.EXPECT().Get(ctx, utils.ComputeKey(fakeAccount.Address, fakeAccount.Namespace)).Return(expectedEntry, nil)

		account, err := usecase.Execute(ctx, fakeAccount.Address, fakeAccount.Namespace)

		assert.NoError(t, err)
		assert.Equal(t, fakeAccount.Namespace, account.Namespace)
		assert.True(t, common.IsHexAddress(account.Address))
	})

	t.Run("should fail with same error if Get fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockStorage.EXPECT().Get(ctx, gomock.Any()).Return(nil, expectedErr)

		account, err := usecase.Execute(ctx, "0xaddress", "namespace")

		assert.Nil(t, account)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should return nil if nothing is found", func(t *testing.T) {
		mockStorage.EXPECT().Get(ctx, gomock.Any()).Return(nil, nil)

		account, err := usecase.Execute(ctx, "0xaddress", "namespace")

		assert.Nil(t, account)
		assert.Nil(t, err)
	})
}
