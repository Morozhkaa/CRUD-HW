// The db package provides methods for working directly with the database.
package db

import (
	"billing-service/internal/domain/models"
	"billing-service/internal/ports"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	Pool *pgxpool.Pool
}

var _ ports.BillingStorage = (*DBStorage)(nil)

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

func (db *DBStorage) GetBalance(ctx context.Context, username string) (int64, error) {
	const query = `
	SELECT balance FROM balances WHERE username = $1;
	`
	var balance int64
	row := db.Pool.QueryRow(ctx, query, username)
	err := row.Scan(&balance)
	if err != nil {
		switch {
		case err.Error() == "no rows in result set":
			return 0, models.ErrUserNotFound
		default:
			return 0, fmt.Errorf("failed to get balance: %w", err)
		}
	}
	return balance, nil
}

func (db *DBStorage) UpdateBalance(ctx context.Context, username string, amount int64) (err error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var balance int64
	row := tx.QueryRow(ctx, "SELECT balance FROM balances WHERE username = $1", username)
	err = row.Scan(&balance)
	if err != nil {
		switch {
		case err.Error() == "no rows in result set" && amount < 0: // проверка на errors.Is(err, pgx.ErrNoRows) почему-то не работала
			return models.ErrUserNotFound
		case err.Error() != "no rows in result set":
			return fmt.Errorf("failed to get balance: %w", err)
		}
	}
	total_count := balance + amount
	if total_count < 0 {
		return models.ErrNotEnoughFunds
	}

	const query = `
	INSERT INTO balances (username, balance)
	VALUES ($1, $2)
	ON CONFLICT (username)
	DO UPDATE SET balance = $2
	WHERE balances.username = $1;
	`
	_, err = tx.Exec(ctx, query, username, total_count)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
