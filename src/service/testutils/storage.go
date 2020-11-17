package testutils

import (
	"context"
	"github.com/hashicorp/vault/sdk/logical"
)

//go:generate mockgen -source=storage.go -destination=mocks/storage.go -package=mocks

// Storage is the way that logical backends are able read/write data.
type Storage interface {
	List(context.Context, string) ([]string, error)
	Get(context.Context, string) (*logical.StorageEntry, error)
	Put(context.Context, *logical.StorageEntry) error
	Delete(context.Context, string) error
}
