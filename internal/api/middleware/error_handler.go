package middleware

import (
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
		if appErr, ok := err.(*utils.AppError); ok {
			utils.JSONError(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
			return
		}
		// unknown error
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "oops!,Something went wrong")
	}
}
