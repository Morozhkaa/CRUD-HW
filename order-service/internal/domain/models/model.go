package models

import "github.com/google/uuid"

type Order struct {
	ProductID int64 `json:"product_id" example:"312"`
	Quantity  int64 `json:"quantity" example:"2"`
	Price     int64 `json:"price" example:"560"`
}

var (
	StatusSuccess = "success"
	StatusFail    = "failed"
)

type OrderInfo struct {
	Username  string    `json:"username" example:"Maria"`
	OrderID   uuid.UUID `json:"order_id" example:"3f8f0d05-0c59-4e7b-a7b6-48e0d5c11f71"`
	ProductID int64     `json:"product_id" example:"312"`
	Quantity  int64     `json:"quantity" example:"2"`
	Price     int64     `json:"price" example:"560"`
	TotalCost int64     `json:"total_cost" example:"1120"`
	Status    string    `json:"status" example:"success"`
}

type UpdateBalance struct {
	Amount int64 `json:"amount" example:"100"`
}

type SuccessResponse struct {
	Success string `json:"success"`
}

type CreateOrderResponse struct {
	Success string    `json:"success"`
	OrderID uuid.UUID `json:"order_id" example:"3f8f0d05-0c59-4e7b-a7b6-48e0d5c11f71"`
}

type VerifyResponse struct {
	Email string `json:"email" example:"maria@example.com"`
	Login string `json:"login" example:"maria"`
}
