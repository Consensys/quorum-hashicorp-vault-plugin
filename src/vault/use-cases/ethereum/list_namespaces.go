package ethereum

import (
	"context"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/logical"
	"strings"
)

// listNamespacesUseCase is a use case to get a list of Ethereum accounts
type listNamespacesUseCase struct {
	storage logical.Storage
}

// NewListAccountUseCase creates a new ListAccountsUseCase
func NewListNamespacesUseCase() usecases.ListNamespacesUseCase {
	return &listNamespacesUseCase{}
}

func (uc listNamespacesUseCase) WithStorage(storage logical.Storage) usecases.ListNamespacesUseCase {
	uc.storage = storage
	return &uc
}

// Execute get a list of all available namespaces
func (uc *listNamespacesUseCase) Execute(ctx context.Context) ([]string, error) {
	logger := apputils.Logger(ctx)
	logger.Debug("listing ethereum namespaces")

	namespaceSet := make(map[string]bool)
	err := uc.getNamespaces(ctx, "", namespaceSet)
	if err != nil {
		return nil, err
	}

	namespaces := make([]string, 0, len(namespaceSet))
	for namespace := range namespaceSet {
		if namespace != "" {
			namespace = strings.TrimSuffix(namespace, "/")
		}
		namespaces = append(namespaces, namespace)
	}

	logger.Debug("ethereum namespaces found successfully")
	return namespaces, nil

}

func (uc *listNamespacesUseCase) getNamespaces(ctx context.Context, prefix string, namespaceSet map[string]bool) error {
	if strings.HasSuffix(prefix, "ethereum/") {
		namespace := strings.TrimSuffix(prefix, "ethereum/")
		namespaceSet[namespace] = true
		return nil
	}

	keys, err := uc.storage.List(ctx, prefix)
	if err != nil {
		errMessage := "failed to get namespace"
		apputils.Logger(ctx).With("prefix", prefix).With("error", err).Error(errMessage)
		return err
	}

	for _, key := range keys {
		err := uc.getNamespaces(ctx, prefix+key, namespaceSet)
		if err != nil {
			errMessage := "failed to get namespace"
			apputils.Logger(ctx).With("prefix", prefix).With("error", err).Error(errMessage)
			return err
		}
	}

	return nil
}
