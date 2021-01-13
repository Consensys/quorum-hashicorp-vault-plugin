package formatters

import (
	"github.com/hashicorp/vault/sdk/logical"
)

func GetRequestNamespace(req *logical.Request) string {
	namespace := ""

	if val, hasVal := req.Headers[NamespaceHeader]; hasVal {
		namespace = val[0]
	}

	return namespace
}
