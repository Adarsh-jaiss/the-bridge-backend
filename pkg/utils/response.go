package utils

import "github.com/gin-gonic/gin"



type APIResponse struct {
	Success bool      `json:"success"`
	Data    any       `json:"data"`
	Error   *APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSONSuccess(c *gin.Context, status int, data any) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	})
}

func JSONError(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}
