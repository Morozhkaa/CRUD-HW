// The usecases package implements the application's business logic. Since the functions are simple
// and there is almost no preliminary preparation before working with data, we immediately call the storage methods.
package usecases

import (
	"billing-service/internal/ports"
	"context"
)

type BillingSvc struct {
	storage ports.BillingStorage
}

var _ ports.BillingService = (*BillingSvc)(nil)

// New returns a new instance of BillingSvc.
func New(storage ports.BillingStorage) *BillingSvc {
	return &BillingSvc{
		storage: storage,
	}
}

func (us *BillingSvc) UpdateBalance(ctx context.Context, username string, amount int64) error {
	return us.storage.UpdateBalance(ctx, username, amount)
}

func (us *BillingSvc) GetBalance(ctx context.Context, username string) (int64, error) {
	return us.storage.GetBalance(ctx, username)
}
