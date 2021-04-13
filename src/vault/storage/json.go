package storage

import (
	"context"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/errors"
	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/pkg/log"

	"github.com/hashicorp/vault/sdk/logical"
)

func StoreJSON(ctx context.Context, storage logical.Storage, key string, data interface{}) error {
	logger := log.FromContext(ctx).With("key", key)

	entry, err := logical.StorageEntryJSON(key, data)
	if err != nil {
		errMessage := "failed to create JSON entry"
		logger.With("error", err).Error(errMessage)
		return errors.StorageError(errMessage)
	}

	err = storage.Put(ctx, entry)
	if err != nil {
		errMessage := "failed to store entry"
		logger.With("error", err).Error(errMessage)
		return errors.StorageError(errMessage)
	}

	return nil
}

func GetJSON(ctx context.Context, storage logical.Storage, key string, entity interface{}) error {
	logger := log.FromContext(ctx).With("key", key)

	entry, err := storage.Get(ctx, key)
	if err != nil {
		errMessage := "failed to get resource"
		logger.With("error", err).Error(errMessage)
		return errors.StorageError(errMessage)
	}

	if entry == nil {
		errMessage := "resource could not be found"
		logger.With("error", err).Error(errMessage)
		return errors.NotFoundError(errMessage)
	}

	err = entry.DecodeJSON(&entity)
	if err != nil {
		errMessage := "could not decode entity"
		logger.With("error", err).Error(errMessage)
		return errors.EncodingError(errMessage)
	}

	return nil
}
