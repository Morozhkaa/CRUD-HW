// The models package contains a description of the main entities used in the service.
package models

import "fmt"

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrBadRequest     = fmt.Errorf("missing required Authorization header")                                           // 400
	ErrNotEnoughFunds = fmt.Errorf("not enough funds to write off")                                                   // 400
	ErrForbidden      = fmt.Errorf("forbidden: authentication failed")                                                // 403
	ErrUserNotFound   = fmt.Errorf("user account not found: to register you need to deposit money into your account") // 404
)
