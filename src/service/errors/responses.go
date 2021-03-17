package errors

import "github.com/hashicorp/vault/sdk/logical"

func BadRequestError(req *logical.Request, msg string) (*logical.Response, error) {
	return logical.RespondWithStatusCode(logical.ErrorResponse(msg), req, 400)
}
