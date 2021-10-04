package migrations

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/entities"
	"time"

	usecases "github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/logical"
)

type ethToKeysUseCase struct {
	storage      logical.Storage
	ethUseCases  usecases.ETHUseCases
	keysUseCases usecases.KeysUseCases
	status       map[string]*entities.MigrationStatus
}

func NewEthToKeysUseCase(ethUseCases usecases.ETHUseCases, keysUseCases usecases.KeysUseCases) usecases.EthereumToKeysUseCase {
	return &ethToKeysUseCase{
		ethUseCases:  ethUseCases,
		keysUseCases: keysUseCases,
	}
}

func (uc ethToKeysUseCase) WithStorage(storage logical.Storage) usecases.EthereumToKeysUseCase {
	uc.storage = storage
	return &uc
}

func (uc ethToKeysUseCase) Status(namespace string) *entities.MigrationStatus {
	return uc.status[namespace]
}

func (uc *ethToKeysUseCase) Execute(ctx context.Context, namespace string) error {
	logger := log.FromContext(ctx).With("namespace", namespace)

	status := &entities.MigrationStatus{
		Status:    "pending",
		StartTime: time.Now(),
	}
	uc.status[namespace] = status

	go func() {
		addresses, err := uc.ethUseCases.ListAccounts().Execute(ctx, namespace)
		if err != nil {
			status.Error = err
			return
		}

		for _, address := range addresses {
			account, der := uc.ethUseCases.GetAccount().WithStorage(uc.storage).Execute(ctx, address, namespace)
			if der != nil {
				status.Status = "failure"
				status.Error = der
				return
			}

			privKey := fmt.Sprintf("0x%s", account.PrivateKey)

			// The ID of the key is the address of the ETH account
			_, der = uc.keysUseCases.CreateKey().WithStorage(uc.storage).Execute(
				ctx,
				namespace,
				address,
				entities.ECDSA,
				entities.Secp256k1,
				privKey,
				map[string]string{},
			)
			if der != nil {
				status.Status = "failure"
				status.Error = der
				return
			}

			status.N += 1
		}

		status.Status = "success"
		status.EndTime = time.Now()
	}()

	logger.Info("migration from ethereum to keys namespace initiated")
	return nil
}
