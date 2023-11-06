package http

import (
	"errors"
	"net/http"
	"user-service/internal/domain/models"
	"user-service/pkg/infra/logger"

	"github.com/gin-gonic/gin"
)

func (a *Adapter) ErrorHandler(ctx *gin.Context, err error) {
	logger.Get().Warn("request failed: ", "desc", err.Error())

	switch {
	case errors.Is(err, models.ErrInvalidEmailFormat), errors.Is(err, models.ErrInvalidPhoneFormat),
		errors.Is(err, models.ErrUserAlreadyExists), errors.Is(err, models.ErrBadRequest):
		ctx.JSON(
			http.StatusBadRequest,
			models.ErrorResponse{Error: err.Error()},
		)
	case errors.Is(err, models.ErrUserNotFound):
		ctx.JSON(
			http.StatusNotFound,
			models.ErrorResponse{Error: err.Error()},
		)
	default:
		ctx.JSON(
			http.StatusInternalServerError,
			models.ErrorResponse{Error: err.Error()},
		)
	}
}
