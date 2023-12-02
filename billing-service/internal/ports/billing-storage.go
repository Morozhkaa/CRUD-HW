package ports

import (
	"context"
)

type BillingStorage interface {
	GetBalance(ctx context.Context, username string) (int64, error)
	UpdateBalance(ctx context.Context, username string, amount int64) error
}
