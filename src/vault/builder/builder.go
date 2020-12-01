package builder

import (
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/use-cases/ethereum"
)

type useCases struct {
	createAccount       usecases.CreateAccountUseCase
	getAccount          usecases.GetAccountUseCase
	listAccounts        usecases.ListAccountsUseCase
	listNamespaces      usecases.ListNamespacesUseCase
	sign                usecases.SignUseCase
	signTx              usecases.SignTransactionUseCase
	signQuorumPrivateTx usecases.SignQuorumPrivateTransactionUseCase
	signEEATx           usecases.SignEEATransactionUseCase
}

func NewEthereumUseCases() usecases.UseCases {
	getAccount := ethereum.NewGetAccountUseCase()
	return &useCases{
		createAccount:       ethereum.NewCreateAccountUseCase(),
		getAccount:          getAccount,
		listAccounts:        ethereum.NewListAccountsUseCase(),
		listNamespaces:      ethereum.NewListNamespacesUseCase(),
		sign:                ethereum.NewSignUseCase(getAccount),
		signTx:              ethereum.NewSignTransactionUseCase(getAccount),
		signQuorumPrivateTx: ethereum.NewSignQuorumPrivateTransactionUseCase(getAccount),
		signEEATx:           ethereum.NewSignEEATransactionUseCase(getAccount),
	}
}

func (ucs *useCases) CreateAccount() usecases.CreateAccountUseCase {
	return ucs.createAccount
}

func (ucs *useCases) GetAccount() usecases.GetAccountUseCase {
	return ucs.getAccount
}

func (ucs *useCases) ListAccounts() usecases.ListAccountsUseCase {
	return ucs.listAccounts
}

func (ucs *useCases) ListNamespaces() usecases.ListNamespacesUseCase {
	return ucs.listNamespaces
}

func (ucs *useCases) SignPayload() usecases.SignUseCase {
	return ucs.sign
}

func (ucs *useCases) SignTransaction() usecases.SignTransactionUseCase {
	return ucs.signTx
}

func (ucs *useCases) SignQuorumPrivateTransaction() usecases.SignQuorumPrivateTransactionUseCase {
	return ucs.signQuorumPrivateTx
}

func (ucs *useCases) SignEEATransaction() usecases.SignEEATransactionUseCase {
	return ucs.signEEATx
}
