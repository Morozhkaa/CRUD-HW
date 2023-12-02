package models

import "github.com/google/uuid"

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
}

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UpdateBalance struct {
	Amount int64 `json:"amount" example:"100"`
}

type CreateOrderResponse struct {
	Success string    `json:"success"`
	OrderID uuid.UUID `json:"order_id" example:"3f8f0d05-0c59-4e7b-a7b6-48e0d5c11f71"`
}

type Order struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderInfo struct {
	Username  string    `json:"username" example:"Maria"`
	OrderID   uuid.UUID `json:"order_id" example:"3f8f0d05-0c59-4e7b-a7b6-48e0d5c11f71"`
	ProductID int64     `json:"product_id" example:"312"`
	Quantity  int64     `json:"quantity" example:"2"`
	Price     int64     `json:"price" example:"560"`
	TotalCost int64     `json:"total_cost" example:"1120"`
	Status    string    `json:"status" example:"success"`
}
