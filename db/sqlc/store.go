package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

//Store providse all functions for db queries and transaction
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) 
}

//SQLStore provides all functions for sql queries and transaction
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

//NewStore creates a new Store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries: New(connPool),
	}
}

