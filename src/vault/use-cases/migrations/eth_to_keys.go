package migrations

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/log"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/entities"
	"sync"
	"time"

	usecases "github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/logical"
)

type ethToKeysUseCase struct {
	ethUseCases  usecases.ETHUseCases
	keysUseCases usecases.KeysUseCases
	status       map[string]*entities.MigrationStatus
	mux          sync.RWMutex
}

type migrationAccount struct {
	address   string
	namespace string
}

func NewEthToKeysUseCase(ethUseCases usecases.ETHUseCases, keysUseCases usecases.KeysUseCases) usecases.EthereumToKeysUseCase {
	return &ethToKeysUseCase{
		ethUseCases:  ethUseCases,
		keysUseCases: keysUseCases,
		status:       map[string]*entities.MigrationStatus{},
		mux:          sync.RWMutex{},
	}
}

func (uc *ethToKeysUseCase) Status(ctx context.Context, namespace string) (*entities.MigrationStatus, error) {
	logger := log.FromContext(ctx).With("namespace", namespace)

	status := uc.status[namespace]
	if status == nil {
		errMessage := "migration could not be found"
		logger.Warn(errMessage)
		return nil, errors.NotFoundError(errMessage)
	}

	return status, nil
}

func (uc *ethToKeysUseCase) Execute(ctx context.Context, storage logical.Storage, sourceNamespace, destinationNamespace string) error {
	logger := log.FromContext(ctx).With("source_namespace", sourceNamespace, "destination_namespace", destinationNamespace)

	if uc.getStatus(destinationNamespace) != nil && uc.getStatus(destinationNamespace).Status == "pending" {
		errMessage := "migration is currently running, please check its status"
		logger.Warn(errMessage)
		return errors.AlreadyExistsError(errMessage)
	}

	namespaces := []string{sourceNamespace}
	var err error
	if sourceNamespace == "*" {
		namespaces, err = uc.ethUseCases.ListNamespaces().WithStorage(storage).Execute(ctx)
		if err != nil {
			return err
		}
	}

	var accounts []migrationAccount
	for _, namespace := range namespaces {
		currAddresses, err := uc.ethUseCases.ListAccounts().WithStorage(storage).Execute(ctx, namespace)
		if err != nil {
			return err
		}

		for _, address := range currAddresses {
			accounts = append(accounts, migrationAccount{
				address:   address,
				namespace: namespace,
			})
		}
	}

	status := &entities.MigrationStatus{
		Status:    "pending",
		StartTime: time.Now(),
		Total:     len(accounts),
	}
	uc.writeStatus(destinationNamespace, status)

	go func() {
		newCtx := log.Context(context.Background(), logger)

		for _, acc := range accounts {
			retrievedAccount, der := uc.ethUseCases.GetAccount().WithStorage(storage).Execute(newCtx, acc.address, acc.namespace)
			if der != nil {
				status.Status = "failure"
				status.Error = der
				return
			}

			// Private keys are stored in hex format without "0x" prefix, they must be transformed to base64
			privKey, der := hex.DecodeString(retrievedAccount.PrivateKey)
			if der != nil {
				errMessage := "failed to decode private key"
				logger.With("error", err).Error(errMessage)
				status.Status = "failure"
				status.Error = errors.EncodingError(errMessage)
				return
			}

			// The ID of the key is the address of the ETH account
			_, der = uc.keysUseCases.CreateKey().WithStorage(storage).Execute(
				newCtx,
				destinationNamespace,
				acc.address,
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

	logger.With("total", len(accounts)).Info("migration from ethereum to keys namespace initiated")
	return nil
}

func (uc *ethToKeysUseCase) getStatus(namespace string) *entities.MigrationStatus {
	uc.mux.RLock()
	defer uc.mux.RUnlock()

	return uc.status[namespace]
}

func (uc *ethToKeysUseCase) writeStatus(namespace string, status *entities.MigrationStatus) {
	uc.mux.Lock()
	defer uc.mux.Unlock()

	uc.status[namespace] = status
}
