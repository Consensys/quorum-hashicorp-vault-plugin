package ethereum

import (
	"fmt"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"net/http"
	"testing"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/service/formatters"
	apputils "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

func (s *ethereumCtrlTestSuite) TestEthereumController_SignQuorumPrivateTransaction() {
	path := s.controller.Paths()[5]
	signOperation := path.Operations[logical.CreateOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, fmt.Sprintf("ethereum/accounts/%s/sign-quorum-private-transaction", framework.GenericNameRegex(formatters.IDLabel)), path.Pattern)
		assert.NotEmpty(t, signOperation)
	})

	s.T().Run("should define correct properties", func(t *testing.T) {
		properties := signOperation.Properties()

		assert.NotEmpty(t, properties.Description)
		assert.NotEmpty(t, properties.Summary)
		assert.NotEmpty(t, properties.Examples[0].Description)
		assert.NotEmpty(t, properties.Examples[0].Response)
		assert.NotEmpty(t, properties.Examples[0].Data)
		assert.NotEmpty(t, properties.Responses[200])
		assert.NotEmpty(t, properties.Responses[400])
		assert.NotEmpty(t, properties.Responses[404])
		assert.NotEmpty(t, properties.Responses[500])
	})

	s.T().Run("handler should execute the correct use case", func(t *testing.T) {
		account := apputils.FakeETHAccount()
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				formatters.NamespaceHeader: {account.Namespace},
			},
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel:       account.Address,
				formatters.NonceLabel:    0,
				formatters.ToLabel:       "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
				formatters.AmountLabel:   "0",
				formatters.GasPriceLabel: "0",
				formatters.GasLimitLabel: 21000,
				formatters.DataLabel:     "0xfeee",
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel:       formatters.AddressFieldSchema,
				formatters.NonceLabel:    formatters.NonceFieldSchema,
				formatters.ToLabel:       formatters.ToFieldSchema,
				formatters.AmountLabel:   formatters.AmountFieldSchema,
				formatters.GasPriceLabel: formatters.GasPriceFieldSchema,
				formatters.GasLimitLabel: formatters.GasLimitFieldSchema,
				formatters.DataLabel:     formatters.DataFieldSchema,
			},
		}
		expectedSignature := "0x8b9679a75861e72fa6968dd5add3bf96e2747f0f124a2e728980f91e1958367e19c2486a40fdc65861824f247603bc18255fa497ca0b8b0a394aa7a6740fdc4601"
		expectedTx, _ := formatters.FormatSignQuorumPrivateTransactionRequest(data)

		s.signQuorumPrivateTransactionUC.EXPECT().Execute(gomock.Any(), account.Address, account.Namespace, expectedTx).Return(expectedSignature, nil)

		response, err := signOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, expectedSignature, response.Data[formatters.SignatureLabel])
	})

	s.T().Run("should fail with 400 if validation fails", func(t *testing.T) {
		account := apputils.FakeETHAccount()
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel: account.Address,
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel:       formatters.AddressFieldSchema,
				formatters.NonceLabel:    formatters.NonceFieldSchema,
				formatters.ToLabel:       formatters.ToFieldSchema,
				formatters.AmountLabel:   formatters.AmountFieldSchema,
				formatters.GasPriceLabel: formatters.GasPriceFieldSchema,
				formatters.GasLimitLabel: formatters.GasLimitFieldSchema,
				formatters.DataLabel:     formatters.DataFieldSchema,
			},
		}

		response, err := signOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Data[logical.HTTPStatusCode])
	})

	s.T().Run("should fail with 400 if data is invalid", func(t *testing.T) {
		account := apputils.FakeETHAccount()
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel:       account.Address,
				formatters.NonceLabel:    0,
				formatters.ToLabel:       "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
				formatters.AmountLabel:   "0",
				formatters.GasPriceLabel: "0",
				formatters.GasLimitLabel: 21000,
				formatters.DataLabel:     "",
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel:       formatters.AddressFieldSchema,
				formatters.NonceLabel:    formatters.NonceFieldSchema,
				formatters.ToLabel:       formatters.ToFieldSchema,
				formatters.AmountLabel:   formatters.AmountFieldSchema,
				formatters.GasPriceLabel: formatters.GasPriceFieldSchema,
				formatters.GasLimitLabel: formatters.GasLimitFieldSchema,
				formatters.DataLabel:     formatters.DataFieldSchema,
			},
		}

		response, err := signOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Data[logical.HTTPStatusCode])
	})

	s.T().Run("should fail with 404 if NotFoundError is returned by use case", func(t *testing.T) {
		account := apputils.FakeETHAccount()
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel:       account.Address,
				formatters.NonceLabel:    0,
				formatters.ToLabel:       "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
				formatters.AmountLabel:   "0",
				formatters.GasPriceLabel: "0",
				formatters.GasLimitLabel: 21000,
				formatters.DataLabel:     "0xfeee",
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel:       formatters.AddressFieldSchema,
				formatters.NonceLabel:    formatters.NonceFieldSchema,
				formatters.ToLabel:       formatters.ToFieldSchema,
				formatters.AmountLabel:   formatters.AmountFieldSchema,
				formatters.GasPriceLabel: formatters.GasPriceFieldSchema,
				formatters.GasLimitLabel: formatters.GasLimitFieldSchema,
				formatters.DataLabel:     formatters.DataFieldSchema,
			},
		}
		expectedErr := errors.NotFoundError("not found")

		s.signQuorumPrivateTransactionUC.EXPECT().Execute(gomock.Any(), account.Address, "", gomock.Any()).Return("", expectedErr)

		response, err := signOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.Data[logical.HTTPStatusCode])
	})

	s.T().Run("should fail with 500 if use case fails with any non mapped error", func(t *testing.T) {
		account := apputils.FakeETHAccount()
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				formatters.IDLabel:       account.Address,
				formatters.NonceLabel:    0,
				formatters.ToLabel:       "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
				formatters.AmountLabel:   "0",
				formatters.GasPriceLabel: "0",
				formatters.GasLimitLabel: 21000,
				formatters.DataLabel:     "0xfeee",
			},
			Schema: map[string]*framework.FieldSchema{
				formatters.IDLabel:       formatters.AddressFieldSchema,
				formatters.NonceLabel:    formatters.NonceFieldSchema,
				formatters.ToLabel:       formatters.ToFieldSchema,
				formatters.AmountLabel:   formatters.AmountFieldSchema,
				formatters.GasPriceLabel: formatters.GasPriceFieldSchema,
				formatters.GasLimitLabel: formatters.GasLimitFieldSchema,
				formatters.DataLabel:     formatters.DataFieldSchema,
			},
		}
		expectedErr := fmt.Errorf("error")

		s.signQuorumPrivateTransactionUC.EXPECT().Execute(gomock.Any(), account.Address, "", gomock.Any()).Return("", expectedErr)

		response, err := signOperation.Handler()(s.ctx, request, data)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Data[logical.HTTPStatusCode])
	})
}
