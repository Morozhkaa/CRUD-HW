// The db package provides methods for working directly with the database.
package db

import (
	"context"
	"fmt"
	"order-service/internal/domain/models"
	"order-service/internal/ports"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	Pool *pgxpool.Pool
}

var _ ports.OrderStorage = (*DBStorage)(nil)

// New establishes one connection and returns a new instance of DBStorage.
func New(ctx context.Context, conn string) (*DBStorage, error) {
	time.Sleep(time.Second)
	pool, err := pgxpool.Connect(ctx, conn)
	if err != nil {
		return nil, err
	}
	return &DBStorage{
		Pool: pool,
	}, nil
}

func (db *DBStorage) GetOrderByID(ctx context.Context, orderID uuid.UUID) (models.OrderInfo, error) {
	const query = `
	SELECT id, username, product_id, quantity, price, total_cost, status FROM orders WHERE id = $1
	`
	var order models.OrderInfo
	err := db.Pool.QueryRow(ctx, query, orderID).Scan(&order.OrderID, &order.Username, &order.ProductID, &order.Quantity, &order.Price, &order.TotalCost, &order.Status)
	if err != nil {
		return order, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

func (db *DBStorage) GetUserOrders(ctx context.Context, username string) ([]models.OrderInfo, error) {
	const query = `
	SELECT id, username, product_id, quantity, price, total_cost, status FROM orders WHERE username = $1
	`
	rows, err := db.Pool.Query(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	defer rows.Close()
	orders := []models.OrderInfo{}
	for rows.Next() {
		var order models.OrderInfo
		err = rows.Scan(&order.OrderID, &order.Username, &order.ProductID, &order.Quantity, &order.Price, &order.TotalCost, &order.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (db *DBStorage) GetAllOrders(ctx context.Context) ([]models.OrderInfo, error) {
	limit, offset := int64(10), int64(0)
	const query = `
	SELECT id, username, product_id, quantity, price, total_cost, status FROM orders limit $1 offset $2
	`
	rows, err := db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}
	defer rows.Close()
	orders := []models.OrderInfo{}
	for rows.Next() {
		var order models.OrderInfo
		err = rows.Scan(&order.OrderID, &order.Username, &order.ProductID, &order.Quantity, &order.Price, &order.TotalCost, &order.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (db *DBStorage) SaveOrder(ctx context.Context, username string, orderID uuid.UUID, order models.Order, status string) error {
	const query = `
	INSERT INTO orders (id, username, product_id, quantity, price, total_cost, status) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.Pool.Exec(ctx, query, orderID, username, order.ProductID, order.Quantity, order.Price, order.Quantity*order.Price, status)
	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}
	return nil
}

func (db *DBStorage) DeleteOrder(ctx context.Context, orderID uuid.UUID) error {
	const query = `
	DELETE FROM orders WHERE id = $1
	`
	_, err := db.Pool.Exec(context.Background(), query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}
