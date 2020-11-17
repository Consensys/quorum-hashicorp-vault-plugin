package ethereum

import (
	ethereum "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/ethereum/use-cases"
)

type useCases struct {
	createAccount       ethereum.CreateAccountUseCase
	sign                ethereum.SignUseCase
	signTx              ethereum.SignTransactionUseCase
	signQuorumPrivateTx ethereum.SignQuorumPrivateTransactionUseCase
	signEEATx           ethereum.SignEEATransactionUseCase
}

func NewEthereumUseCases() ethereum.UseCases {
	return &useCases{
		createAccount: ethereum.NewCreateAccountUseCase(),
	}
}

func (ucs *useCases) CreateAccount() ethereum.CreateAccountUseCase {
	return ucs.createAccount
}

func (ucs *useCases) SignPayload() ethereum.SignUseCase {
	return ucs.sign
}

func (ucs *useCases) SignTransaction() ethereum.SignTransactionUseCase {
	return ucs.signTx
}

func (ucs *useCases) SignQuorumPrivateTransaction() ethereum.SignQuorumPrivateTransactionUseCase {
	return ucs.signQuorumPrivateTx
}

func (ucs *useCases) SignEEATransaction() ethereum.SignEEATransactionUseCase {
	return ucs.signEEATx
}
