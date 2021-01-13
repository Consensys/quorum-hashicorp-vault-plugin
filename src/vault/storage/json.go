package storage

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
)

func StoreJSON(ctx context.Context, storage logical.Storage, key string, data interface{}) error {
	entry, err := logical.StorageEntryJSON(key, data)
	if err != nil {
		return err
	}

	err = storage.Put(ctx, entry)
	if err != nil {
		return err
	}

	return nil
}

func GetJSON(ctx context.Context, storage logical.Storage, key string, account interface{}) error {
	entry, err := storage.Get(ctx, key)
	if err != nil {
		return err
	}

	if entry == nil {
		return logical.CodedError(404, "account could not be found")
	}
	
	err = entry.DecodeJSON(&account)
	if err != nil {
		return err
	}
	
	return nil
}
