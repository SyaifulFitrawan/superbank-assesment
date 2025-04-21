package database

import (
	"context"

	"gorm.io/gorm"
)

type ctxKey string

const TxKey ctxKey = "tx"

func NewContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func FromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(TxKey).(*gorm.DB)
	if ok {
		return tx
	}
	return fallback
}

func WithTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
