package ethereum

import (
	"context"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/testutils/mocks"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/testutils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/ethereum/utils"
	"github.com/hashicorp/go-hclog"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListAccounts_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := apputils.WithLogger(context.Background(), hclog.New(&hclog.LoggerOptions{}))

	usecase := NewListAccountsUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		fakeAccount := testutils.FakeETHAccount()
		expectedKeys := []string{fakeAccount.Address}
		mockStorage.EXPECT().List(ctx, utils.ComputeKey("", fakeAccount.Namespace)).Return(expectedKeys, nil)

		keys, err := usecase.Execute(ctx, fakeAccount.Namespace)

		assert.NoError(t, err)
		assert.Equal(t, expectedKeys, keys)
	})

	t.Run("should fail with same error if List fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockStorage.EXPECT().List(ctx, gomock.Any()).Return(nil, expectedErr)

		keys, err := usecase.Execute(ctx, "namespace")

		assert.Nil(t, keys)
		assert.Equal(t, expectedErr, err)
	})
}
