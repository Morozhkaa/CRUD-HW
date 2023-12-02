package models

type UpdateBalance struct {
	Amount int64 `json:"amount" example:"1000"`
}

type SuccessResponse struct {
	Success string `json:"success"`
}

type VerifyResponse struct {
	Email string `json:"email" example:"maria@example.com"`
	Login string `json:"login" example:"maria"`
}
