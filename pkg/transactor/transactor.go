package transactor

import (
	"context"
)

type Transactor interface {
	// Do executes the given function within a transaction.
	// The transaction is added to the context, so it has to be retrieved
	// appropriately depending on the transactor implementation.
	Do(context.Context, func(context.Context) error) error
	// Skip shadows the transaction in the context.
	Skip(ctx context.Context) context.Context
}
