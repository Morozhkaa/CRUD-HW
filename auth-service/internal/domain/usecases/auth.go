// The usecases package implements the application's business logic. Since the functions are simple
// and there is almost no preliminary preparation before working with data, we immediately call the storage methods.
package usecases

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
	"user-service/internal/config"
	"user-service/internal/domain/models"
	"user-service/internal/ports"

	jwt4 "github.com/golang-jwt/jwt/v4"
	"github.com/xdg-go/pbkdf2"
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

const accessToken_expiration_time = 1 * time.Minute
const refreshToken_expiration_time = 60 * time.Minute

func EncodePassword(password string) string {
	cfg := config.Get()
	dk := pbkdf2.Key([]byte(password), []byte(cfg.Salt), 1000, 128, sha1.New)
	return base64.StdEncoding.EncodeToString([]byte(dk))
}

func (us *UserSvc) generateToken(login string, email string, expiredIn time.Duration) (token string, err error) {
	cfg, now := config.Get(), time.Now()
	claims := &jwt4.RegisteredClaims{
		ExpiresAt: jwt4.NewNumericDate(now.Add(expiredIn)),
		Issuer:    login,
		Subject:   email,
	}

	t := jwt4.NewWithClaims(jwt4.SigningMethodHS256, claims)
	token, err = t.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", models.ErrGenerateToken
	}
	return token, nil
}

func (us *UserSvc) verifyToken(tokenString string) (login, email string, err error) {
	cfg := config.Get()
	var claims jwt4.RegisteredClaims
	token, err := jwt4.ParseWithClaims(tokenString, &claims, func(token *jwt4.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if !token.Valid {
		return "", "", fmt.Errorf("parse token unexpected error: %w", err)
	}
	return claims.Issuer, claims.Subject, nil
}

func (us *UserSvc) Login(ctx context.Context, login, password string) (access, refresh string, err error) {
	user, err := us.storage.GetUserDetails(ctx, login)
	if err != nil {
		return "", "", models.ErrUserNotFound
	}
	if user.Password != EncodePassword(password) {
		return "", "", models.ErrForbidden
	}
	access, err = us.generateToken(login, user.Email, accessToken_expiration_time)
	if err != nil {
		return "", "", models.ErrGenerateToken
	}
	refresh, err = us.generateToken(login, user.Email, refreshToken_expiration_time)
	if err != nil {
		return "", "", models.ErrGenerateToken
	}
	return access, refresh, nil
}

func (us *UserSvc) Verify(ctx context.Context, access, refresh string) (r models.VerifyResponse, err error) {
	r.Login, r.Email, err = us.verifyToken(access)
	if err == nil {
		r.AccessToken = access
		r.RefreshToken = refresh
		return r, nil
	}

	r.Login, r.Email, err = us.verifyToken(refresh)
	if err == nil {
		r.AccessToken, err = us.generateToken(r.Login, r.Email, accessToken_expiration_time)
		if err != nil {
			return r, models.ErrGenerateToken
		}
		r.RefreshToken, err = us.generateToken(r.Login, r.Email, refreshToken_expiration_time)
		if err != nil {
			return r, models.ErrGenerateToken
		}
		return r, nil
	}
	return r, models.ErrTokenExpired
}

func (us *UserSvc) CreateUser(ctx context.Context, user models.User) error {
	if _, err := us.storage.GetUser(ctx, user.Username); err == nil {
		return models.ErrUserAlreadyExists
	}
	user.Password = EncodePassword(user.Password)
	return us.storage.SaveUser(ctx, user)
}

func (us *UserSvc) DeleteUser(ctx context.Context, username string) error {
	if _, err := us.storage.GetUser(ctx, username); err != nil {
		return models.ErrUserNotFound
	}
	return us.storage.DeleteUser(ctx, username)
}

func (us *UserSvc) UpdateUser(ctx context.Context, username string, user models.User) error {
	if _, err := us.storage.GetUser(ctx, username); err != nil {
		return models.ErrUserNotFound
	}
	return us.storage.UpdateUser(ctx, username, user)
}

func (us *UserSvc) GetUser(ctx context.Context, username string) (models.GetUserResponse, error) {
	return us.storage.GetUser(ctx, username)
}
