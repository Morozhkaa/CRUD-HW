// The ports package contains the description of the interfaces.
package ports

import (
	"context"
	"user-service/internal/domain/models"
)

type UserService interface {
	Login(ctx context.Context, login, password string) (string, string, error)
	Verify(ctx context.Context, access, refresh string) (models.VerifyResponse, error)
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, username string) (models.GetUserResponse, error)
	UpdateUser(ctx context.Context, username string, user models.User) error
	DeleteUser(ctx context.Context, username string) error
}
