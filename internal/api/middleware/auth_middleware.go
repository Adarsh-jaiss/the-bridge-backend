package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/adarsh-jaiss/the-bridge/internal/auth"
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/adarsh-jaiss/the-bridge/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := logger.FromContext(ctx)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn("auth header is empty")
			c.Error(utils.NewUnauthorizedError("missing authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			log.Warn("invalid authorization format")
			c.Error(utils.NewUnauthorizedError("invalid authorization format"))
			return
		}

		claims, err := auth.ValidateToken(ctx, tokenStr, auth.AccessToken)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Next()
			return
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			refreshToken, err := c.Cookie("refresh_token")
			if err != nil {
				log.Error("refresh token missing in cookie")
				c.Error(utils.NewUnauthorizedError("refresh token is missing"))
				return
			}

			refreshClaims, err := auth.ValidateToken(ctx, refreshToken, auth.RefreshToken)
			if err != nil {
				log.Warn("invalid refresh token")
				c.Error(utils.NewUnauthorizedError("invalid refresh token"))
				return
			}
			if refreshClaims.ExpiresAt == nil {
				c.Error(utils.NewUnauthorizedError("invalid refresh token expiry"))
				return
			}

			userId := refreshClaims.UserID
			newAccessToken, err := auth.GenerateToken(userId, auth.AccessToken)
			if err != nil {
				log.Error("error generating new access token", zap.Error(err))
				c.Error(utils.NewInternalServerError(err))
				return
			}

			refreshExpiry := refreshClaims.ExpiresAt.Time
			timeLeft := time.Until(refreshExpiry)

			if timeLeft <= 15*time.Minute {
				newRefreshToken, err := auth.GenerateToken(userId, auth.RefreshToken)
				if err != nil {
					log.Error("error generating new refresh token", zap.Error(err))
					c.Error(utils.NewInternalServerError(err))
					return
				}

				c.SetCookie("refresh_token", newRefreshToken, int((7 * 24 * time.Hour).Seconds()), "/", "", true, true)
			}

			c.Header("X-New-Access-Token", newAccessToken)
			c.Set("user_id", userId)
			c.Next()
			return

		}

		log.Warn("invalid access token", zap.Error(err))
		c.Error(utils.NewUnauthorizedError("invalid access token"))

	}
}
