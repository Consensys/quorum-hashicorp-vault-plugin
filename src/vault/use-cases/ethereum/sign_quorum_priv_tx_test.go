package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils/mocks"
	mocks2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/mocks"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSignQuorumPrivateTransaction_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetAccountUC := mocks2.NewMockGetAccountUseCase(ctrl)
	mockStorage := mocks.NewMockStorage(ctrl)
	ctx := log.Context(context.Background(), log.Default())
	address := "0xaddress"
	namespace := "namespace"
	tx := quorumtypes.NewTransaction(
		0,
		common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(10000000000),
		21000,
		big.NewInt(10000000000),
		[]byte{},
	)

	mockGetAccountUC.EXPECT().WithStorage(mockStorage).Return(mockGetAccountUC)
	usecase := NewSignQuorumPrivateTransactionUseCase(mockGetAccountUC).WithStorage(mockStorage)

	t.Run("should execute use case successfully", func(t *testing.T) {
		fakeAccount := apputils.FakeETHAccount()
		fakeAccount.PrivateKey = "5385714a2f6d69ca034f56a5268833216ffb8fba7229c39569bc4c5f42cde97c"
		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(fakeAccount, nil)

		signature, err := usecase.Execute(ctx, address, namespace, tx)

		assert.NoError(t, err)
		assert.Equal(t, "0xefa9c4498397ee12e341f6acf81072bbf0c8fb4e4e1813ac96fd3860baa28bb931aecd59811beaffc71a4ef008882d3c13537a2f733be7643fdfea4ea77f3ded00", signature)
	})

	t.Run("should fail with CryptoOperationError if creation of ECDSA private key fails", func(t *testing.T) {
		fakeAccount := apputils.FakeETHAccount()
		fakeAccount.PrivateKey = "invalidPrivKey"
		mockGetAccountUC.EXPECT().Execute(ctx, address, namespace).Return(fakeAccount, nil)

		signature, err := usecase.Execute(ctx, address, namespace, tx)

		assert.Empty(t, signature)
		assert.Error(t, err)
	})

	t.Run("should fail with same error if Get account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")

		mockGetAccountUC.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(nil, expectedErr)

		signature, err := usecase.Execute(ctx, address, namespace, tx)

		assert.Empty(t, signature)
		assert.Equal(t, expectedErr, err)
	})
}
