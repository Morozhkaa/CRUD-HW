package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order-service/internal/domain/models"
	"order-service/internal/ports"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Success 200 {object} models.SuccessResponse
// @Router /health [get]
func (a *Adapter) health(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		models.SuccessResponse{Success: "service available"},
	)
}

func formBillingRequest(ctx *gin.Context, billingURL string, totalCount int64) (*http.Request, error) {
	// get the access and refresh tokens
	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader == "" {
		return nil, models.ErrForbidden
	}
	tokens := strings.Split(authorizationHeader, "Bearer ")[1]
	access, refresh := strings.Split(tokens, ";")[0], strings.Split(tokens, ";")[1]

	// form a request to the billing service
	body := models.UpdateBalance{
		Amount: totalCount,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}
	billingR, err := http.NewRequestWithContext(ctx.Request.Context(), "POST", billingURL+"balance", bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	billingR.Header.Set("Authorization", "Bearer "+fmt.Sprintf("%s;%s", access, refresh))
	return billingR, nil
}

// @ID createOrder
// @tags order
// @Summary Create order
// @Description Create a new order
// @Accept json
// @Param Authorization header string true "Authorization header: 'Bearer <access_token>;<refresh_token>'"
// @Param order body models.Order true "order"
// @Success 200 {object} models.CreateOrderResponse "New order was successfully created."
// @Failure 400 {object} models.ErrorResponse "Not enough funds to write off / required Authorization header or 'order' parameter is missing."
// @Failure 403 {object} models.ErrorResponse "Authentication failed"
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /createOrder [post]
func (a *Adapter) createOrder(ctx *gin.Context) {
	username := ctx.MustGet("login").(string)
	var in models.Order
	err := ctx.BindJSON(&in)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}
	// make a request to the billing service; if successful, we withdraw money for the order
	billingR, err := formBillingRequest(ctx, a.billingURL, -1*in.Price*in.Quantity)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	orderStatus := models.StatusSuccess
	resp, err := a.client.Do(billingR)
	if err != nil || resp.StatusCode != http.StatusOK {
		orderStatus = models.StatusFail
	}
	defer resp.Body.Close()

	orderID, err := a.orderSvc.SaveOrder(ctx, username, in, orderStatus)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	var msg string
	switch {
	case resp.StatusCode == http.StatusOK:
		msg = "order was successfully created"
	case resp.StatusCode == http.StatusBadRequest:
		msg = "order creation failed: not enough funds to write off"
	case resp.StatusCode == http.StatusNotFound:
		msg = "order creation failed: user account not found; to register you need to deposit money into your account"
	}
	ctx.JSON(
		http.StatusOK,
		models.CreateOrderResponse{Success: msg, OrderID: orderID},
	)
}

// @ID getOrderByID
// @tags order
// @Summary Get order
// @Description Returns the order by order_id.
// @Param Authorization header string true "Authorization header: 'Bearer <access_token>;<refresh_token>'"
// @Param orderID path string true "orderID in uuid format"
// @Success 200 {object} models.OrderInfo "The user's order have been successfully received."
// @Failure 400 {object} models.ErrorResponse "Missing required Authorization header."
// @Failure 403 {object} models.ErrorResponse "Authentication failed"
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /getOrderByID/{orderID} [get]
func (a *Adapter) getOrderByID(ctx *gin.Context) {
	oID := ctx.Param("orderID")
	if oID == "" || oID == ":orderID" {
		a.ErrorHandler(ctx, models.ErrBadRequest)
	}
	orderID, err := uuid.Parse(oID)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}
	orders, err := a.orderSvc.GetOrderByID(ctx, orderID)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

// @ID getUserOrders
// @tags order
// @Summary Get user orders
// @Description Returns the user's orders.
// @Param Authorization header string true "Authorization header: 'Bearer <access_token>;<refresh_token>'"
// @Success 200 {object} models.OrderInfo "The user's orders have been successfully received."
// @Failure 400 {object} models.ErrorResponse "Missing required Authorization header."
// @Failure 403 {object} models.ErrorResponse "Authentication failed"
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /getUserOrders [get]
func (a *Adapter) getUserOrders(ctx *gin.Context) {
	username := ctx.MustGet("login").(string)
	orders, err := a.orderSvc.GetUserOrders(ctx, username)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

// @ID getAllOrders
// @tags order
// @Summary Get all orders
// @Description Returns all orders.
// @Param Authorization header string true "Authorization header: 'Bearer <access_token>;<refresh_token>'"
// @Success 200 {object} []models.OrderInfo "All orders has been successfully received."
// @Failure 400 {object} models.ErrorResponse "Missing required Authorization header."
// @Failure 403 {object} models.ErrorResponse "Authentication failed"
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /getAllOrders [get]
func (a *Adapter) getAllOrders(ctx *gin.Context) {
	orders, err := a.orderSvc.GetAllOrders(ctx)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

// @ID deleteOrder
// @tags order
// @Summary Delete order
// @Description Deletes order by orderID.
// @Param Authorization header string true "Authorization header: 'Bearer <access_token>;<refresh_token>'"
// @Param orderID path string true "orderID in uuid format"
// @Success 200 {object} models.SuccessResponse "The order has been successfully deleted."
// @Failure 400 {object} models.ErrorResponse "The required Authorization header or 'orderID' field is missing."
// @Failure 403 {object} models.ErrorResponse "Authentication failed"
// @Failure 500 {object} models.ErrorResponse "Database error / Internal Server Error."
// @Router /deleteOrder/{orderID} [delete]
func (a *Adapter) deleteOrder(ctx *gin.Context) {
	oID := ctx.Param("orderID")
	if oID == "" || oID == ":orderID" {
		a.ErrorHandler(ctx, models.ErrBadRequest)
	}
	orderID, err := uuid.Parse(oID)
	if err != nil {
		a.ErrorHandler(ctx, models.ErrBadRequest)
		return
	}
	err = a.orderSvc.DeleteOrder(ctx, orderID)
	if err != nil {
		a.ErrorHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, models.SuccessResponse{Success: "order successfully deleted"})
}
