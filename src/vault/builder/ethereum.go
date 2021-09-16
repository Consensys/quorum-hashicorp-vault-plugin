package builder

import (
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases"
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/use-cases/ethereum"
)

type ethereumUseCases struct {
	createAccount       usecases.CreateAccountUseCase
	getAccount          usecases.GetAccountUseCase
	listAccounts        usecases.ListAccountsUseCase
	listNamespaces      usecases.ListNamespacesUseCase
	sign                usecases.SignUseCase
	signTx              usecases.SignTransactionUseCase
	signQuorumPrivateTx usecases.SignQuorumPrivateTransactionUseCase
	signEEATx           usecases.SignEEATransactionUseCase
}

func NewEthereumUseCases() usecases.ETHUseCases {
	getAccount := ethereum.NewGetAccountUseCase()
	return &ethereumUseCases{
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

func (ucs *ethereumUseCases) CreateAccount() usecases.CreateAccountUseCase {
	return ucs.createAccount
}

func (ucs *ethereumUseCases) GetAccount() usecases.GetAccountUseCase {
	return ucs.getAccount
}

func (ucs *ethereumUseCases) ListAccounts() usecases.ListAccountsUseCase {
	return ucs.listAccounts
}

func (ucs *ethereumUseCases) ListNamespaces() usecases.ListNamespacesUseCase {
	return ucs.listNamespaces
}

func (ucs *ethereumUseCases) SignPayload() usecases.SignUseCase {
	return ucs.sign
}

func (ucs *ethereumUseCases) SignTransaction() usecases.SignTransactionUseCase {
	return ucs.signTx
}

func (ucs *ethereumUseCases) SignQuorumPrivateTransaction() usecases.SignQuorumPrivateTransactionUseCase {
	return ucs.signQuorumPrivateTx
}

func (ucs *ethereumUseCases) SignEEATransaction() usecases.SignEEATransactionUseCase {
	return ucs.signEEATx
}
