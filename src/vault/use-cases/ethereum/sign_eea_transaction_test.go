package ethereum

import (
	"context"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/testutils/mocks"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/testutils"
	mocks2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/mocks"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/go-hclog"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignEEATransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetAccountUC := mocks2.NewMockGetAccountUseCase(ctrl)
	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := apputils.WithLogger(context.Background(), hclog.New(&hclog.LoggerOptions{}))
	address := "0xaddress"
	namespace := "namespace"
	chainID := "1"
	tx := ethtypes.NewTransaction(
		0,
		common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(10000000000),
		21000,
		big.NewInt(10000000000),
		[]byte{},
	)

	mockGetAccountUC.EXPECT().WithStorage(mockStorage).Return(mockGetAccountUC)
	usecase := NewSignEEATransactionUseCase(mockGetAccountUC).WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		fakeAccount := testutils.FakeETHAccount()
		fakeAccount.PrivateKey = "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"
		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(fakeAccount, nil)
		privateArgs := testutils2.FakePrivateETHTransactionParams()

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx, privateArgs)

		assert.NoError(t, err)
		assert.Equal(t, "0x2424ed4546e2039c9f132222eb361286a485a5e9eade6fc5ee1c9548d5391e146d7e794a3a5aa7b553d3905f65824633aab893985112f1c48e6f51e0a8ceb02001", signature)
	})

	t.Run("should fail with same error if Get account fails", func(t *testing.T) {
		privateArgs := testutils2.FakePrivateETHTransactionParams()
		expectedErr := fmt.Errorf("error")

		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(nil, expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx, privateArgs)

		assert.Empty(t, signature)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should fail with CryptoOperationError if creation of ECDSA private key fails", func(t *testing.T) {
		privateArgs := testutils2.FakePrivateETHTransactionParams()
		fakeAccount := testutils.FakeETHAccount()
		fakeAccount.PrivateKey = "invalidPrivKey"
		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(fakeAccount, nil)

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx, privateArgs)

		assert.Empty(t, signature)
		assert.Error(t, err)
	})
}
