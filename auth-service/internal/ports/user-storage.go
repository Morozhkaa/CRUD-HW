package ports

import (
	"context"
	"user-service/internal/domain/models"
)

type UserStorage interface {
	SaveUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, username string) (models.GetUserResponse, error)
	UpdateUser(ctx context.Context, username string, user models.User) error
	DeleteUser(ctx context.Context, username string) error
	GetUserDetails(ctx context.Context, username string) (models.User, error)
}
