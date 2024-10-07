package memory

import (
	"context"
	"log"
)

type TransactionManager struct {
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{}
}

func (tm *TransactionManager) Begin(ctx context.Context) (context.Context, error) {
	log.Println("Begin transaction")
	return ctx, nil
}

func (tm *TransactionManager) Commit(ctx context.Context) error {
	log.Println("Commit transaction")
	return nil
}

func (tm *TransactionManager) Rollback(ctx context.Context) error {
	log.Println("Rollback transaction")
	return nil
}
