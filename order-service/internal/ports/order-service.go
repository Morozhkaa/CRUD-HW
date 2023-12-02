// The ports package contains the description of the interfaces.
package ports

import (
	"context"
	"order-service/internal/domain/models"

	"github.com/google/uuid"
)

type OrderService interface {
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (models.OrderInfo, error)
	GetAllOrders(ctx context.Context) ([]models.OrderInfo, error)
	GetUserOrders(ctx context.Context, username string) ([]models.OrderInfo, error)
	SaveOrder(ctx context.Context, username string, order models.Order, status string) (uuid.UUID, error)
	DeleteOrder(ctx context.Context, orderID uuid.UUID) error
}
