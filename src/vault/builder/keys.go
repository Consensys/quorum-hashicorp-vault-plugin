package builder

import (
	usecases "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/keys"
)

type keysUseCases struct {
	createKey      usecases.CreateKeyUseCase
	getKey         usecases.GetKeyUseCase
	listKeys       usecases.ListKeysUseCase
	listNamespaces usecases.ListKeysNamespacesUseCase
	sign           usecases.KeysSignUseCase
}

func NewKeysUseCases() usecases.KeysUseCases {
	getKeyUC := keys.NewGetKeyUseCase()
	return &keysUseCases{
		createKey:      keys.NewCreateKeyUseCase(),
		getKey:         getKeyUC,
		listKeys:       keys.NewListKeysUseCase(),
		listNamespaces: keys.NewListNamespacesUseCase(),
		sign:           keys.NewSignUseCase(getKeyUC),
	}
}

func (ucs *keysUseCases) CreateKey() usecases.CreateKeyUseCase {
	return ucs.createKey
}

func (ucs *keysUseCases) GetKey() usecases.GetKeyUseCase {
	return ucs.getKey
}

func (ucs *keysUseCases) ListKeys() usecases.ListKeysUseCase {
	return ucs.listKeys
}

func (ucs *keysUseCases) ListNamespaces() usecases.ListKeysNamespacesUseCase {
	return ucs.listNamespaces
}

func (ucs *keysUseCases) SignPayload() usecases.KeysSignUseCase {
	return ucs.sign
}
