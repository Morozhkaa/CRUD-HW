// The models package contains a description of the main entities used in the service.
package models

import "fmt"

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrInvalidEmailFormat    = fmt.Errorf("invalid format of parameter 'email'")                     // 400
	ErrInvalidPhoneFormat    = fmt.Errorf("invalid format of parameter 'phone'")                     // 400
	ErrBadRequest            = fmt.Errorf("missing required parameters")                             // 400
	ErrUserAlreadyExists     = fmt.Errorf("user with this username already exists")                  // 400
	ErrNoAuthorizationHeader = fmt.Errorf("authentication failed: no required Authorization header") // 400
	ErrGenerateToken         = fmt.Errorf("generate token failed")                                   // 400
	ErrForbidden             = fmt.Errorf("forbidden: wrong password")                               // 403
	ErrTokenExpired          = fmt.Errorf("token expired")                                           // 403
	ErrUserNotFound          = fmt.Errorf("user not found: make sure you are registered")            // 404
)
