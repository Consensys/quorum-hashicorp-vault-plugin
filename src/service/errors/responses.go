package errors

import (
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/hashicorp/vault/sdk/logical"
	"net/http"
)

func WriteHTTPError(req *logical.Request, err error) (*logical.Response, error) {
	switch {
	case errors.IsAlreadyExistsError(err):
		return logical.RespondWithStatusCode(logical.ErrorResponse(err.Error()), req, http.StatusConflict)
	case errors.IsNotFoundError(err):
		return logical.RespondWithStatusCode(logical.ErrorResponse(err.Error()), req, http.StatusNotFound)
	case errors.IsInvalidParameterError(err), errors.IsEncodingError(err):
		return logical.RespondWithStatusCode(logical.ErrorResponse(err.Error()), req, http.StatusUnprocessableEntity)
	case errors.IsInvalidFormatError(err):
		return logical.RespondWithStatusCode(logical.ErrorResponse(err.Error()), req, http.StatusBadRequest)
	default:
		return logical.RespondWithStatusCode(logical.ErrorResponse("internal server error. Please retry or contact an administrator"), req, http.StatusInternalServerError)
	}
}
