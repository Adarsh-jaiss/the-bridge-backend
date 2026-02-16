package middleware

import (
	"github.com/adarsh-jaiss/the-bridge/pkg/utils"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if appErr, ok := err.(*utils.AppError); ok {
				utils.JSONError(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
				return
			}

			// fallback
			utils.JSONError(c, 500, "INTERNAL_ERROR", "Something went wrong")
		}
	}
}
