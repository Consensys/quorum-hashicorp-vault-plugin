package keys

import (
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

func (s *keysCtrlTestSuite) TestKeysController_ListNamespaces() {
	path := s.controller.Paths()[4]
	listOperation := path.Operations[logical.ListOperation]

	s.T().Run("should define the correct path", func(t *testing.T) {
		assert.Equal(t, "namespaces/keys/?", path.Pattern)
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
		expectedList := []string{"ns1/ns2", "_"}
		request := &logical.Request{
			Storage: s.storage,
		}

		s.listNamespacesUC.EXPECT().Execute(gomock.Any()).Return(expectedList, nil)

		response, err := listOperation.Handler()(s.ctx, request, &framework.FieldData{})
		require.NoError(t, err)

		assert.Equal(t, expectedList, response.Data["keys"])
	})

	s.T().Run("should map errors correctly and return the correct http status", func(t *testing.T) {
		request := &logical.Request{
			Storage: s.storage,
		}
		expectedErr := errors.NotFoundError("error")

		s.listNamespacesUC.EXPECT().Execute(gomock.Any()).Return(nil, expectedErr)

		_, err := listOperation.Handler()(s.ctx, request, &framework.FieldData{})

		assert.Equal(t, err, logical.ErrUnsupportedPath)
	})
}
