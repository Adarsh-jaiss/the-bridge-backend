package middleware

import (
	"errors"

	"github.com/adarsh-jaiss/the-bridge/pkg/utils"
	"github.com/gin-gonic/gin"
)
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) == 0 {
            return
        }

        err := c.Errors.Last().Err

        var appErr *utils.AppError
        if errors.As(err, &appErr) {  // ← handles wrapped errors
            utils.JSONError(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
            return
        }

        utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "oops!, Something went wrong")
    }
}