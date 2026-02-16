package routes

import (
	"database/sql"
	"net/http"

	"github.com/adarsh-jaiss/the-bridge/pkg/api/controllers"
	"github.com/adarsh-jaiss/the-bridge/pkg/api/repository"
	"github.com/adarsh-jaiss/the-bridge/pkg/api/usecase"
	"github.com/adarsh-jaiss/the-bridge/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Routes(db *sql.DB, log *zap.Logger, atomicLevel zap.AtomicLevel) {
	authrepo := repository.NewUserAuthentication(db)
	authInteractor := usecase.NewUserAuthInteractor(authrepo)
	authController := controllers.NewUserAuthController(&authInteractor)

	r := gin.Default()

	api := r.Group("api")
	api.Use(middleware.ErrorHandler())
	api.Use(middleware.RequestLogger(log))

	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Example admin route to change log level
	admin := api.Group("/admin")
	// admin.Use(middleware.AdminAuth()) // protect this
	admin.POST("/log-level", ChangeLogLevel(atomicLevel, log))

	v1 := api.Group("v1")

	auth := v1.Group("auth")
	auth.POST("/trigger-otp", authController.TriggerOtp)

	r.Run()
}

func ChangeLogLevel(atomicLevel zap.AtomicLevel, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		type req struct {
			Level string `json:"level"`
		}

		var body req
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		switch body.Level {
		case "debug":
			atomicLevel.SetLevel(zap.DebugLevel)
		case "info":
			atomicLevel.SetLevel(zap.InfoLevel)
		case "warn":
			atomicLevel.SetLevel(zap.WarnLevel)
		case "error":
			atomicLevel.SetLevel(zap.ErrorLevel)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level"})
			return
		}

		log.Info("Log level changed",
			zap.String("new_level", body.Level),
		)

		c.JSON(http.StatusOK, gin.H{"message": "log level updated"})
	}
}
