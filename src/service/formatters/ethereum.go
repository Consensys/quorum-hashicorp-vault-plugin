package formatters

import (
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"math/big"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/vault/entities"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func FormatAccountResponse(account *entities.ETHAccount) *logical.Response {
	return &logical.Response{
		Data: map[string]interface{}{
			"address":             account.Address,
			"publicKey":           account.PublicKey,
			"compressedPublicKey": account.CompressedPublicKey,
			"namespace":           account.Namespace,
		},
	}
}

func FormatSignatureResponse(signature string) *logical.Response {
	return &logical.Response{
		Data: map[string]interface{}{
			"signature": signature,
		},
	}
}

func FormatSignETHTransactionRequest(requestData *framework.FieldData) (*types.Transaction, error) {
	amount, ok := new(big.Int).SetString(requestData.Get(AmountLabel).(string), 10)
	if !ok {
		return nil, errors.InvalidFormatError("invalid amount")
	}

	gasPrice, ok := new(big.Int).SetString(requestData.Get(GasPriceLabel).(string), 10)
	if !ok {
		return nil, errors.InvalidFormatError("invalid gas price")
	}

	data, err := hexutil.Decode(requestData.Get(DataLabel).(string))
	if err != nil {
		return nil, errors.InvalidFormatError("invalid data")
	}

	nonce := requestData.Get(NonceLabel).(int)
	gasLimit := requestData.Get(GasLimitLabel).(int)
	to := requestData.Get(ToLabel).(string)
	if to == "" {
		return types.NewContractCreation(uint64(nonce), amount, uint64(gasLimit), gasPrice, data), nil
	}

	return types.NewTransaction(uint64(nonce), common.HexToAddress(to), amount, uint64(gasLimit), gasPrice, data), nil
}

func FormatSignQuorumPrivateTransactionRequest(requestData *framework.FieldData) (*quorumtypes.Transaction, error) {
	amount, ok := new(big.Int).SetString(requestData.Get(AmountLabel).(string), 10)
	if !ok {
		return nil, errors.InvalidFormatError("invalid amount")
	}

	gasPrice, ok := new(big.Int).SetString(requestData.Get(GasPriceLabel).(string), 10)
	if !ok {
		return nil, errors.InvalidFormatError("invalid gas price")
	}

	data, err := hexutil.Decode(requestData.Get(DataLabel).(string))
	if err != nil {
		return nil, errors.InvalidFormatError("invalid data")
	}

	nonce := requestData.Get(NonceLabel).(int)
	gasLimit := requestData.Get(GasLimitLabel).(int)
	to := requestData.Get(ToLabel).(string)
	if to == "" {
		return quorumtypes.NewContractCreation(uint64(nonce), amount, uint64(gasLimit), gasPrice, data), nil
	}

	return quorumtypes.NewTransaction(uint64(nonce), common.HexToAddress(to), amount, uint64(gasLimit), gasPrice, data), nil
}

func FormatSignEEATransactionRequest(requestData *framework.FieldData) (tx *types.Transaction, privateArgs *entities.PrivateETHTransactionParams, err error) {
	data, err := hexutil.Decode(requestData.Get(DataLabel).(string))
	if err != nil {
		return nil, nil, errors.InvalidFormatError("invalid data")
	}

	amount := big.NewInt(0)
	gasPrice := big.NewInt(0)
	gas := uint64(0)
	to := requestData.Get(ToLabel).(string)
	nonce := requestData.Get(NonceLabel).(int)

	privateArgs = &entities.PrivateETHTransactionParams{
		PrivateFrom:    requestData.Get(PrivateFromLabel).(string),
		PrivateFor:     requestData.Get(PrivateForLabel).([]string),
		PrivacyGroupID: requestData.Get(PrivacyGroupIDLabel).(string),
		PrivateTxType:  "restricted",
	}
	if to == "" {
		return types.NewContractCreation(uint64(nonce), amount, gas, gasPrice, data), privateArgs, nil
	}

	return types.NewTransaction(uint64(nonce), common.HexToAddress(to), amount, gas, gasPrice, data), privateArgs, nil
}
