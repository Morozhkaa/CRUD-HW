// The usecases package implements the application's business logic. Since the functions are simple
// and there is almost no preliminary preparation before working with data, we immediately call the storage methods.
package usecases

import (
	"context"
	"order-service/internal/domain/models"
	"order-service/internal/ports"

	"github.com/google/uuid"
)

type OrderSvc struct {
	storage ports.OrderStorage
}

var _ ports.OrderService = (*OrderSvc)(nil)

// New returns a new instance of OrderSvc.
func New(storage ports.OrderStorage) *OrderSvc {
	return &OrderSvc{
		storage: storage,
	}
}

func (s *OrderSvc) GetOrderByID(ctx context.Context, orderID uuid.UUID) (models.OrderInfo, error) {
	return s.storage.GetOrderByID(ctx, orderID)
}
func (s *OrderSvc) GetAllOrders(ctx context.Context) ([]models.OrderInfo, error) {
	return s.storage.GetAllOrders(ctx)
}

func (s *OrderSvc) GetUserOrders(ctx context.Context, username string) ([]models.OrderInfo, error) {
	return s.storage.GetUserOrders(ctx, username)
}

func (s *OrderSvc) SaveOrder(ctx context.Context, username string, order models.Order, status string) (uuid.UUID, error) {
	orderID := uuid.New()
	return orderID, s.storage.SaveOrder(ctx, username, orderID, order, status)
}

func (s *OrderSvc) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	return s.storage.DeleteOrder(ctx, orderID)
}
