package ports

import (
	"context"
	"order-service/internal/domain/models"

	"github.com/google/uuid"
)

type OrderStorage interface {
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (models.OrderInfo, error)
	GetAllOrders(ctx context.Context) ([]models.OrderInfo, error)
	GetUserOrders(ctx context.Context, username string) ([]models.OrderInfo, error)
	SaveOrder(ctx context.Context, username string, orderID uuid.UUID, order models.Order, status string) error
	DeleteOrder(ctx context.Context, orderID uuid.UUID) error
}
