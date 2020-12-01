package ethereum

import (
	"context"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/testutils/mocks"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/hashicorp/go-hclog"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListNamespaces_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := apputils.WithLogger(context.Background(), hclog.New(&hclog.LoggerOptions{}))

	usecase := NewListNamespacesUseCase().WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		// We have 3 namespaces: "ns1/ns2", "_" and "tenant0"
		// We also have accounts with no namespace (starting directly with "ethereum/")
		expectedNamespaces := []string{"ns1/ns2", "_", "tenant0", ""}

		gomock.InOrder(
			mockStorage.EXPECT().List(ctx, "").Return([]string{"ns1/", "_/", "tenant0/", "ethereum/"}, nil),

			// ns1/ns2 with 1 account
			mockStorage.EXPECT().List(ctx, "ns1/").Return([]string{"ns2/"}, nil),
			mockStorage.EXPECT().List(ctx, "ns1/ns2/").Return([]string{"ethereum/"}, nil),

			// _ with 1 account
			mockStorage.EXPECT().List(ctx, "_/").Return([]string{"ethereum/"}, nil),

			// tenant0 with 2 accounts
			mockStorage.EXPECT().List(ctx, "tenant0/").Return([]string{"ethereum/", "ethereum/"}, nil),
		)

		namespaces, err := usecase.Execute(ctx)

		assert.NoError(t, err)
		assert.Contains(t, namespaces, expectedNamespaces[0])
		assert.Contains(t, namespaces, expectedNamespaces[1])
		assert.Contains(t, namespaces, expectedNamespaces[2])
		assert.Contains(t, namespaces, expectedNamespaces[3])
	})

	t.Run("should fail with same error if List fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockStorage.EXPECT().List(ctx, gomock.Any()).Return(nil, expectedErr)

		keys, err := usecase.Execute(ctx)

		assert.Nil(t, keys)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail with same error if recursive List fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		gomock.InOrder(
			mockStorage.EXPECT().List(ctx, "").Return([]string{"ns1/", "_/", "tenant0/", "ethereum/"}, nil),
			mockStorage.EXPECT().List(ctx, gomock.Any()).Return(nil, expectedErr),
		)
		keys, err := usecase.Execute(ctx)

		assert.Nil(t, keys)
		assert.Equal(t, expectedErr, err)
	})
}
