package http

import (
	"errors"
	"net/http"
	"order-service/internal/domain/models"
	"order-service/pkg/infra/logger"

	"github.com/gin-gonic/gin"
)

func (a *Adapter) ErrorHandler(ctx *gin.Context, err error) {
	logger.Get().Warn("request failed: ", "desc", err.Error())

	switch {
	case errors.Is(err, models.ErrBadRequest), errors.Is(err, models.ErrNotEnoughFunds):
		ctx.JSON(
			http.StatusBadRequest,
			models.ErrorResponse{Error: err.Error()},
		)
	case errors.Is(err, models.ErrForbidden):
		ctx.JSON(
			http.StatusForbidden,
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
