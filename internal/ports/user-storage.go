package ports

import (
	"user-service/internal/domain/models"
)

type UserStorage interface {
	SaveUser(user models.User) error
	GetUser(username string) (models.GetUserResponse, error)
	UpdateUser(username string, user models.User) error
	DeleteUser(username string) error
}
