package http

import (
	"billing-service/internal/domain/models"
	"billing-service/internal/ports"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func authMiddleware(a ports.AuthAdapter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := a.Verify(ctx)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrForbidden):
				ctx.JSON(http.StatusForbidden, gin.H{
					"error": err.Error(),
				})
			case errors.Is(err, models.ErrBadRequest):
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
			return
		}
		ctx.Next()
	}
}

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

// @ID updateBalance
// @tags balance
// @Summary Update balance
// @Description Updates user balance.
// @Accept json
// @Param Authorization header string true "Authorization header: 'Bearer <access_token>;<refresh_token>'"
// @Param user body models.UpdateBalance true "amount (may be negative)"
// @Success 200 {object} models.SuccessResponse "The user's balance information has been successfully updated."
// @Failure 400 {object} models.ErrorResponse "Not enough funds to write off / required Authorization header or 'amount' field is missing."
// @Failure 403 {object} models.ErrorResponse "Authentication failed"
// @Failure 404 {object} models.ErrorResponse "User account not found"
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /balance [post]
func (a *Adapter) updateBalance(ctx *gin.Context) {
	username := ctx.MustGet("login").(string)
	var in models.UpdateBalance

	err := ctx.BindJSON(&in)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}
	err = a.billingSvc.UpdateBalance(ctx, username, in.Amount)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK,
		models.SuccessResponse{Success: fmt.Sprintf("information for user with username '%s' updated", username)},
	)
}

// @ID getBalance
// @tags balance
// @Summary Get balance
// @Description Returns the user's balance.
// @Param Authorization header string true "Authorization header: 'Bearer <access_token>;<refresh_token>'"
// @Success 200 {object} int64 "The user's balance has been successfully received."
// @Failure 400 {object} models.ErrorResponse "Missing required Authorization header."
// @Failure 403 {object} models.ErrorResponse "Authentication failed"
// @Failure 404 {object} models.ErrorResponse "User account not found"
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /balance [get]
func (a *Adapter) getBalance(ctx *gin.Context) {
	username := ctx.MustGet("login").(string)
	balance, err := a.billingSvc.GetBalance(ctx, username)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, balance)
}
