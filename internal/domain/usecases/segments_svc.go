// The usecases package implements the application's business logic. Since the functions are simple
// and there is almost no preliminary preparation before working with data, we immediately call the storage methods.
package usecases

import (
	"user-service/internal/domain/models"
	"user-service/internal/ports"
)

type UserSvc struct {
	storage ports.UserStorage
}

var _ ports.UserService = (*UserSvc)(nil)

// New returns a new instance of UserSvc.
func New(storage ports.UserStorage) *UserSvc {
	return &UserSvc{
		storage: storage,
	}
}

func (us *UserSvc) CreateUser(user models.User) error {
	if _, err := us.storage.GetUser(user.Username); err == nil {
		return models.ErrUserAlreadyExists
	}
	return us.storage.SaveUser(user)
}

func (us *UserSvc) DeleteUser(username string) error {
	if _, err := us.storage.GetUser(username); err != nil {
		return models.ErrUserNotFound
	}
	return us.storage.DeleteUser(username)
}

func (us *UserSvc) UpdateUser(username string, user models.User) error {
	if _, err := us.storage.GetUser(username); err != nil {
		return models.ErrUserNotFound
	}
	return us.storage.UpdateUser(username, user)
}

func (us *UserSvc) GetUser(username string) (models.GetUserResponse, error) {
	return us.storage.GetUser(username)
}
