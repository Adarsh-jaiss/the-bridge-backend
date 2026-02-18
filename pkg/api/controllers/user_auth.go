package controllers

import (
	"github.com/adarsh-jaiss/the-bridge/pkg/api/types/dto"
	"github.com/adarsh-jaiss/the-bridge/pkg/api/usecase"
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/adarsh-jaiss/the-bridge/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserAuthController struct {
	IAuthInteractor usecase.IUserAuthInteractor
}

func NewUserAuthController(u usecase.IUserAuthInteractor) *UserAuthController {
	return &UserAuthController{
		IAuthInteractor: u,
	}
}

func (a *UserAuthController) TriggerOtp(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	var req dto.EmailVerficationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("invalid email payload", zap.Error(err))
		utils.JSONError(c, 400, "ERROR_BAD_REQUEST", "invalid request body")
		return
	}

	if req.Email == "" {
		log.Warn("email is empty in request")
		utils.JSONError(c, 400, "ERROR_BAD_REQUEST", "email is required")
		return
	}	
	
}
