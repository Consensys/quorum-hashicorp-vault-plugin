package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	mocks2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignTransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetAccountUC := mocks.NewMockGetAccountUseCase(ctrl)
	mockStorage := mocks2.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())
	address := "0xaddress"
	namespace := "namespace"
	chainID := "1"
	tx := types.NewTransaction(
		0,
		common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(10000000000),
		21000,
		big.NewInt(10000000000),
		[]byte{},
	)

	mockGetAccountUC.EXPECT().WithStorage(mockStorage).Return(mockGetAccountUC)
	usecase := NewSignTransactionUseCase(mockGetAccountUC).WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		fakeAccount := apputils.FakeETHAccount()
		fakeAccount.PrivateKey = "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"
		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(fakeAccount, nil)

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx)

		assert.NoError(t, err)
		assert.Equal(t, "0xd35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e01", signature)
	})

	t.Run("should fail with CryptoOperationError if creation of ECDSA private key fails", func(t *testing.T) {
		fakeAccount := apputils.FakeETHAccount()
		fakeAccount.PrivateKey = "invalidPrivKey"
		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(fakeAccount, nil)

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx)

		assert.Empty(t, signature)
		assert.Error(t, err)
	})

	t.Run("should fail with same error if Get account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockGetAccountUC.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, chainID, tx)

		assert.Empty(t, signature)
		assert.Equal(t, expectedErr, err)
	})
}
