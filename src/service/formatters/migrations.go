package formatters

import (
	"github.com/consensys/quorum-hashicorp-vault-plugin/src/vault/entities"
	"github.com/hashicorp/vault/sdk/logical"
)

func FormatMigrationStatusResponse(status *entities.MigrationStatus) *logical.Response {
	return &logical.Response{
		Data: map[string]interface{}{
			"status":    status.Status,
			"error":     status.Error,
			"startTime": status.StartTime,
			"endtime":   status.EndTime,
			"n":         status.N,
		},
	}
}
