package customerr

import "errors"

const INVALID_TOKEN = "invalid token"
const NO_USER_IN_CONTEXT = "no user in context"
const NO_KEY_IN_CONTEXT = "no key in context"

// Custom error
type CustomError struct {
	Message          string `json:"message"`
	Code             int    `json:"-"`
	Err              error  `json:"-"`
	DeveloperMessage string `json:"-"`
}

// Convert error message to json
func ToJson(message string) *CustomError {
	return &CustomError{Message: message}
}

// Custom error
func (c *CustomError) Error() string {
	return c.Message
}

// Unwrap error
func (c *CustomError) Unwrap() error {
	return c.Err
}

// Add developer message
func AddDeveloperMessage(err error, message string) *CustomError {
	return &CustomError{DeveloperMessage: message, Err: err}
}

func Error(message string) *CustomError {
	return &CustomError{Err: errors.New(message), Message: message}
}
