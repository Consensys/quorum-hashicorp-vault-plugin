package storage

import (
	"context"
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"strings"

	"github.com/hashicorp/vault/sdk/logical"
)

const (
	EthereumSecretsPath = "ethereum"
	ZkSnarksSecretsPath = "zk-snarks"
	KeysSecretPath      = "keys"
)

func GetEthereumNamespaces(ctx context.Context, storage logical.Storage, namespace string, namespaceSet map[string]bool) error {
	return getNamespaces(ctx, storage, EthereumSecretsPath, namespace, namespaceSet)
}

func GetZkSnarksNamespaces(ctx context.Context, storage logical.Storage, namespace string, namespaceSet map[string]bool) error {
	return getNamespaces(ctx, storage, ZkSnarksSecretsPath, namespace, namespaceSet)
}

func GetKeysNamespaces(ctx context.Context, storage logical.Storage, namespace string, namespaceSet map[string]bool) error {
	return getNamespaces(ctx, storage, KeysSecretPath, namespace, namespaceSet)
}

func getNamespaces(ctx context.Context, storage logical.Storage, secretsPath, namespace string, namespaceSet map[string]bool) error {
	logger := log.FromContext(ctx).With("namespace", namespace)

	if strings.HasSuffix(namespace, secretsPath+"/") {
		namespace := strings.TrimSuffix(namespace, secretsPath+"/")
		namespaceSet[namespace] = true
		return nil
	}

	keys, err := storage.List(ctx, namespace)
	if err != nil {
		errMessage := "failed to get keys"
		logger.With("error", err).Error(errMessage)
		return errors.StorageError(errMessage)
	}

	for _, key := range keys {
		err := getNamespaces(ctx, storage, secretsPath, namespace+key, namespaceSet)
		if err != nil {
			return err
		}
	}

	return nil
}

func ComputeEthereumStorageKey(accountID, namespace string) string {
	return computeStorageKey(EthereumSecretsPath, accountID, namespace)
}

func ComputeZksStorageKey(accountID, namespace string) string {
	return computeStorageKey(ZkSnarksSecretsPath, accountID, namespace)
}

func ComputeKeysStorageKey(id, namespace string) string {
	path := fmt.Sprintf("%s/%s", KeysSecretPath, id)
	if namespace != "" {
		path = fmt.Sprintf("%s/%s", namespace, path)
	}

	return path
}

func computeStorageKey(secretsPath, accountID, namespace string) string {
	path := fmt.Sprintf("%s/accounts/%s", secretsPath, accountID)
	if namespace != "" {
		path = fmt.Sprintf("%s/%s", namespace, path)
	}

	return path
}
