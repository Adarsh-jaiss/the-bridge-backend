package routes

import (
	"database/sql"
	"net/http"

	_ "github.com/adarsh-jaiss/the-bridge/docs"
	"github.com/adarsh-jaiss/the-bridge/internal/api/controllers"
	"github.com/adarsh-jaiss/the-bridge/internal/api/middleware"
	"github.com/adarsh-jaiss/the-bridge/internal/api/repository"
	"github.com/adarsh-jaiss/the-bridge/internal/api/usecase"
	"github.com/adarsh-jaiss/the-bridge/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	customValidator "github.com/adarsh-jaiss/the-bridge/pkg/validator"
)

func Routes(db *sql.DB, log *zap.Logger, atomicLevel zap.AtomicLevel, cfg *config.Config) {
	authrepo := repository.NewUserAuthentication(db)
	authInteractor := usecase.NewUserAuthInteractor(authrepo)
	authController := controllers.NewUserAuthController(authInteractor)

	profileRepo := repository.NewUserProfile(db)
	profileInteractor := usecase.NewUserProfileInteractor(profileRepo)
	userProfileController := controllers.NewUserProfileController(profileInteractor)

	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		customValidator.RegisterOnValidator(v)
	}

	api := r.Group("api")
	api.Use(middleware.ErrorHandler())
	api.Use(middleware.RequestLogger(log))

	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Example admin route to change log level
	admin := api.Group("/admin")
	// admin.Use(middleware.AdminAuth()) // protect this
	admin.POST("/log-level", ChangeLogLevel(atomicLevel, log))

	v1 := api.Group("v1")

	auth := v1.Group("auth")
	auth.POST("/trigger-otp", authController.TriggerOtp)
	auth.POST("/verify-otp", authController.VerifyOTP)
	auth.GET("/dummy-tokens",authController.GenerateDummyTokens)

	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware())

	user := protected.Group("user")
	user.POST("/onboarding", userProfileController.CreateProfile)
	user.POST("/:id/follow", userProfileController.Follow)
	user.DELETE("/:id/unfollow", userProfileController.UnFollow)
	user.PATCH("/profile", userProfileController.UpdateProfile)
	user.GET("/profile", userProfileController.GetUserProfile)

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
