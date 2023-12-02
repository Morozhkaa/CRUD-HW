// The ports package contains the description of the interfaces.
package ports

import (
	"context"
)

type BillingService interface {
	GetBalance(ctx context.Context, username string) (int64, error)
	UpdateBalance(ctx context.Context, username string, amount int64) error
}
