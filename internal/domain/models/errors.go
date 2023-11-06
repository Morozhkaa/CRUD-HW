// The models package contains a description of the main entities used in the service.
package models

import "fmt"

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrInvalidEmailFormat = fmt.Errorf("invalid format of parameter 'email'")    // 400
	ErrInvalidPhoneFormat = fmt.Errorf("invalid format of parameter 'phone'")    // 400
	ErrBadRequest         = fmt.Errorf("missing required parameters")            // 400
	ErrUserAlreadyExists  = fmt.Errorf("user with this username already exists") // 400
	ErrUserNotFound       = fmt.Errorf("user not found")                         // 404
)
