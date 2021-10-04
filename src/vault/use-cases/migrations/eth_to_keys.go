package migrations

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
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
		status:       map[string]*entities.MigrationStatus{},
	}
}

func (uc ethToKeysUseCase) WithStorage(storage logical.Storage) usecases.EthereumToKeysUseCase {
	uc.storage = storage
	return &uc
}

func (uc ethToKeysUseCase) Status(ctx context.Context, namespace string) (*entities.MigrationStatus, error) {
	logger := log.FromContext(ctx).With("namespace", namespace)

	status := uc.status[namespace]
	if status == nil {
		errMessage := "migration could not be found"
		logger.Warn(errMessage)
		return nil, errors.NotFoundError(errMessage)
	}

	return status, nil
}

func (uc *ethToKeysUseCase) Execute(ctx context.Context, namespace string) error {
	logger := log.FromContext(ctx).With("namespace", namespace)

	addresses, err := uc.ethUseCases.ListAccounts().WithStorage(uc.storage).Execute(ctx, namespace)
	if err != nil {
		return err
	}

	status := &entities.MigrationStatus{
		Status:    "pending",
		StartTime: time.Now(),
		Total:     len(addresses),
	}
	uc.status[namespace] = status

	go func() {
		newCtx := log.Context(context.Background(), logger)

		for _, address := range addresses {
			account, der := uc.ethUseCases.GetAccount().WithStorage(uc.storage).Execute(newCtx, address, namespace)
			if der != nil {
				status.Status = "failure"
				status.Error = der
				return
			}

			// Private keys are stored in hex format without "0x" prefix, they must be transformed to base64
			privKey, der := hex.DecodeString(account.PrivateKey)
			if der != nil {
				status.Status = "failure"
				status.Error = errors.EncodingError("failed to decode private key")
				return
			}

			// The ID of the key is the address of the ETH account
			_, der = uc.keysUseCases.CreateKey().WithStorage(uc.storage).Execute(
				newCtx,
				namespace,
				address,
				entities.ECDSA,
				entities.Secp256k1,
				base64.URLEncoding.EncodeToString(privKey),
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

	logger.With("total", len(addresses)).Info("migration from ethereum to keys namespace initiated")
	return nil
}
