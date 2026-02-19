package utils

import "net/http"

type AppError struct {
	HTTPStatus int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(status int, code, message string, err error) *AppError {
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func NewBadRequestError(message string) *AppError {
	return NewAppError(http.StatusBadRequest, "BAD_REQUEST", message, nil)
}

func NewInternalServerError(err error) *AppError {
	return NewAppError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Something went wrong", err)
}