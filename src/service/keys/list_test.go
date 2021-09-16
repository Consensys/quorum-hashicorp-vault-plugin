package keys

import (
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/consensys/quorum-hashicorp-vault-plugin/src/service/formatters"
	apputils "github.com/consensys/quorum-hashicorp-vault-plugin/src/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

func (s *keysCtrlTestSuite) TestKeysController_List() {
	path := s.controller.Paths()[0]
	listOperation := path.Operations[logical.ListOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, "keys/?", path.Pattern)
		assert.NotEmpty(t, listOperation)
	})

	s.T().Run("should define correct properties", func(t *testing.T) {
		properties := listOperation.Properties()

		assert.NotEmpty(t, properties.Description)
		assert.NotEmpty(t, properties.Summary)
		assert.NotEmpty(t, properties.Examples[0].Description)
		assert.NotEmpty(t, properties.Examples[0].Response)
		assert.NotEmpty(t, properties.Responses[200])
		assert.NotEmpty(t, properties.Responses[500])
	})

	s.T().Run("handler should execute the correct use case", func(t *testing.T) {
		key := apputils.FakeKey()
		expectedList := []string{key.ID}
		request := &logical.Request{
			Storage: s.storage,
			Headers: map[string][]string{
				formatters.NamespaceHeader: {key.Namespace},
			},
		}
		data := &framework.FieldData{}

		s.listKeysUC.EXPECT().Execute(gomock.Any(), key.Namespace).Return(expectedList, nil)

		response, err := listOperation.Handler()(s.ctx, request, data)
		require.NoError(t, err)

		assert.Equal(t, expectedList, response.Data["keys"])
	})

	s.T().Run("should map errors correctly and return the correct http status", func(t *testing.T) {
		request := &logical.Request{
			Storage: s.storage,
		}
		data := &framework.FieldData{}
		expectedErr := errors.NotFoundError("error")

		s.listKeysUC.EXPECT().Execute(gomock.Any(), "").Return(nil, expectedErr)

		_, err := listOperation.Handler()(s.ctx, request, data)

		assert.Equal(t, err, logical.ErrUnsupportedPath)
	})
}
