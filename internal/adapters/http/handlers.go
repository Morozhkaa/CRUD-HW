package http

import (
	"fmt"
	"net/http"
	"regexp"
	"user-service/internal/domain/models"
	"user-service/pkg/infra/logger"

	"github.com/gin-gonic/gin"
)

// @ID health
// @Summary Check service status
// @Success 200 {object} string
// @Router /health [get]
func (a *Adapter) health(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		models.SuccessResponse{Success: "service available"},
	)
}

// @ID createUser
// @tags user
// @Summary Create user
// @Description Creates a new user with given data. Checks that email and phone are in the correct format, and that the user with given username is not yet in the database, otherwise it returns the BadRequest status.
// @Accept json
// @Param user body models.User true "user data"
// @Success 200 {object} models.SuccessResponse "User created successfully."
// @Failure 400 {object} models.ErrorResponse "User already exists / missing required 'user' parameter / invalid format of 'email' or 'phone' parameters."
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /user [post]
func (a *Adapter) createUser(ctx *gin.Context) {
	var user models.User
	err := ctx.BindJSON(&user)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}
	matchedEmail, err := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, user.Email)
	if err != nil || !matchedEmail {
		a.ErrorHandler(ctx, models.ErrInvalidEmailFormat)
		return
	}
	matchedPhone, err := regexp.MatchString(`^\+\d{11}$`, user.Phone)
	if err != nil || !matchedPhone {
		a.ErrorHandler(ctx, models.ErrInvalidPhoneFormat)
		return
	}
	err = a.userSvc.CreateUser(user)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK,
		models.SuccessResponse{Success: fmt.Sprintf("user with username '%s' created", user.Username)},
	)
}

// @ID deleteUser
// @tags user
// @Summary Delete user
// @Description Deletes the user with given username.
// @Param username path string true "username of the user to delete."
// @Success 200 {object} models.SuccessResponse "User deleted successfully."
// @Failure 400 {object} models.ErrorResponse "Missing required 'username' parameter."
// @Failure 404 {object} models.ErrorResponse "User with given username not found."
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /user/{username} [delete]
func (a *Adapter) deleteUser(ctx *gin.Context) {
	username, err := a.getUsernameFromPath(ctx)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}
	err = a.userSvc.DeleteUser(username)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK,
		models.SuccessResponse{Success: fmt.Sprintf("user with username '%s' deleted", username)},
	)
}

// @ID updateUser
// @tags user
// @Summary Update user
// @Description Updates user data with given username.
// @Accept json
// @Param username path string true "username of the user to update"
// @Param user body models.User true "user data"
// @Success 200 {object} models.SuccessResponse "User information updated successfully."
// @Failure 400 {object} models.ErrorResponse "Missing required 'username' or 'user' parameters."
// @Failure 404 {object} models.ErrorResponse "User with given username not found."
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /user/{username} [put]
func (a *Adapter) updateUser(ctx *gin.Context) {
	username, err := a.getUsernameFromPath(ctx)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	var user models.User
	err = ctx.BindJSON(&user)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}
	err = a.userSvc.UpdateUser(username, user)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK,
		models.SuccessResponse{Success: fmt.Sprintf("information for user with username '%s' updated", username)},
	)
}

// @ID getUser
// @tags user
// @Summary Get user
// @Description Returns information about the user with the given username.
// @Param username path string true "Username of the user to get"
// @Success 200 {object} models.GetUserResponse "User data received successfully."
// @Failure 400 {object} models.ErrorResponse "Missing required 'username' parameter."
// @Failure 404 {object} models.ErrorResponse "User with given username not found."
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /user/{username} [get]
func (a *Adapter) getUser(ctx *gin.Context) {
	username, err := a.getUsernameFromPath(ctx)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	user, err := a.userSvc.GetUser(username)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrUserNotFound)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (a *Adapter) getUsernameFromPath(ctx *gin.Context) (string, error) {
	username := ctx.Param("username")
	logger.Get().Debug("got parameter from path", "username", username)
	if username == ":username" {
		return "", models.ErrBadRequest
	}
	return username, nil
}
