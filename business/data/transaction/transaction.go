// Package transaction provides support for database transaction related functionality.
package transaction

import (
	"context"
)

// Transaction represents a value that can commit or rollback a transaction.
type Transaction interface {
	Commit() error
	Rollback() error
}

// Beginner represents a value that can begin a transaction.
type Beginner interface {
	Begin() (Transaction, error)
}

// =============================================================================

type ctxKey int

const trKey ctxKey = 2

// Set stores a value that can manage a transaction.
func Set(ctx context.Context, tx Transaction) context.Context {
	return context.WithValue(ctx, trKey, tx)
}

// Get retrieves the value that can manage a transaction.
func Get(ctx context.Context) (Transaction, bool) {
	v, ok := ctx.Value(trKey).(Transaction)
	return v, ok
}
