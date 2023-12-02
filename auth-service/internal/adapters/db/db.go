// The db package provides methods for working directly with the database.
package db

import (
	"context"
	"time"
	"user-service/internal/domain/models"
	"user-service/internal/ports"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	Pool *pgxpool.Pool
}

var _ ports.UserStorage = (*DBStorage)(nil)

// var ctx = context.Background()

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

func (db *DBStorage) SaveUser(ctx context.Context, user models.User) (err error) {
	const query = `
	INSERT INTO users (username, password, first_name, last_name, email, phone) VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err = db.Pool.Exec(ctx, query, user.Username, user.Password, user.FirstName, user.LastName, user.Email, user.Phone)
	return err
}

func (db *DBStorage) GetUser(ctx context.Context, username string) (user models.GetUserResponse, err error) {
	const query = `
	SELECT username, first_name, last_name, email, phone FROM users WHERE username = $1;
	`
	err = db.Pool.QueryRow(ctx, query, username).
		Scan(&user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Phone)
	return
}

func (db *DBStorage) GetUserDetails(ctx context.Context, username string) (user models.User, err error) {
	const query = `
	SELECT username, password, first_name, last_name, email, phone FROM users WHERE username = $1;
	`
	err = db.Pool.QueryRow(ctx, query, username).
		Scan(&user.Username, &user.Password, &user.FirstName, &user.LastName, &user.Email, &user.Phone)
	return
}

func (db *DBStorage) UpdateUser(ctx context.Context, username string, user models.User) (err error) {
	const query = `
	UPDATE users
	SET username = $1, password = $2, first_name = $3, last_name = $4, email = $5, phone = $6
	WHERE username = $7;
	`
	_, err = db.Pool.Exec(ctx, query, user.Username, user.Password, user.FirstName, user.LastName, user.Email, user.Phone, username)
	return
}

func (db *DBStorage) DeleteUser(ctx context.Context, username string) (err error) {
	const query = `	
	DELETE FROM users WHERE username = $1;
	`
	_, err = db.Pool.Exec(ctx, query, username)
	return
}
