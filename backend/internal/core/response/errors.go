package response

import (
	"net/http"
)

// AppError represents a standardized error structure
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// New creates a new AppError with optional custom message
func New(code int, defaultMsg string, msg ...string) *AppError {
	message := defaultMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}

	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Predefined errors

func BadRequest(msg ...string) *AppError {
	return New(http.StatusBadRequest, "Bad request", msg...)
}

func Unauthorized(msg ...string) *AppError {
	return New(http.StatusUnauthorized, "Unauthorized", msg...)
}

func Forbidden(msg ...string) *AppError {
	return New(http.StatusForbidden, "Forbidden", msg...)
}

func NotFound(msg ...string) *AppError {
	return New(http.StatusNotFound, "Resource not found", msg...)
}

func Conflict(msg ...string) *AppError {
	return New(http.StatusConflict, "Conflict occurred", msg...)
}

func Internal(msg ...string) *AppError {
	return New(http.StatusInternalServerError, "Internal server error", msg...)
}
