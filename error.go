package app

import (
	"fmt"
)

const UNAUTHORIZED_ERR = "unauthorized"
const INVALID_ERR = "invalid"

type Error struct {
	// Machine-readable error code.
	Code string `json:"code"`

	// Human-readable error message.
	Message string `json:"message"`
}

// Error implements the error interface. Not used by the application otherwise.
func (e *Error) Error() string {
	return fmt.Sprintf("Error: code=%s message=%s", e.Code, e.Message)
}

// Errorf is a helper function to return an Error with a given code and formatted message.
func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
