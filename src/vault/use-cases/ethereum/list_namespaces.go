package ethereum

import (
	"context"
	"strings"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/storage"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/hashicorp/vault/sdk/logical"
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
	logger := log.FromContext(ctx)
	logger.Debug("listing ethereum namespaces")

	namespaceSet := make(map[string]bool)
	err := storage.GetEthereumNamespaces(ctx, uc.storage, "", namespaceSet)
	if err != nil {
		errMessage := "failed to get namespace"
		logger.With("error", err).Error(errMessage)
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
