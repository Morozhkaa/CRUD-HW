// The ports package contains the description of the interfaces.
package ports

import (
	"user-service/internal/domain/models"
)

type UserService interface {
	CreateUser(user models.User) error
	GetUser(username string) (models.GetUserResponse, error)
	UpdateUser(username string, user models.User) error
	DeleteUser(username string) error
}
