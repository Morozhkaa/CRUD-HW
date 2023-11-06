package models

type User struct {
	Username  string `json:"username" example:"IvanIvanov2000"`
	Password  string `json:"password" example:"qwerty1234"`
	FirstName string `json:"first_name" example:"Ivan"`
	LastName  string `json:"last_name" example:"Ivanov"`
	Email     string `json:"email" example:"iivanov@gmail.com"`
	Phone     string `json:"phone" example:"+79999999999"`
}

// User without password
type GetUserResponse struct {
	Username  string `json:"username" example:"IvanIvanov2000"`
	FirstName string `json:"first_name" example:"Ivan"`
	LastName  string `json:"last_name" example:"Ivanov"`
	Email     string `json:"email" example:"iivanov@gmail.com"`
	Phone     string `json:"phone" example:"+79999999999"`
}

type SuccessResponse struct {
	Success string `json:"success"`
}
