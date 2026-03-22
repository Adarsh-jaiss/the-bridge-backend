package controllers

import (
	"strconv"

	"github.com/adarsh-jaiss/the-bridge/internal/api/types/dto"
	"github.com/adarsh-jaiss/the-bridge/internal/api/usecase"

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

// TriggerOtp godoc
// @Summary Trigger OTP
// @Description Sends an OTP to the provided email address for authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.EmailVerficationRequest true "Email payload"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /v1/auth/trigger-otp [post]
func (a *UserAuthController) TriggerOtp(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	var req dto.EmailVerficationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("invalid request body", zap.Error(err))
		utils.JSONError(c, 400, "ERROR_BAD_REQUEST", "invalid request body")
		return
	}

	if req.Email == "" {
		log.Warn("email is empty in request")
		utils.JSONError(c, 400, "ERROR_BAD_REQUEST", "email is required")
		return
	}
	if err := a.IAuthInteractor.TriggerOTP(ctx, req.Email); err != nil {
		log.Error("error sending otp", zap.Error(err))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "failed to trigger otp")
		return
	}

	utils.JSONSuccess(c, 200, gin.H{"message": "OTP Sent Successfully"})
}

// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verifies OTP and returns access & refresh tokens on success
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.OTPVerificationRequest true "OTP verification payload"
// @Success 200 {object} utils.SuccessResponse{data=dto.OTPVerificationResponse}
// @Failure 400 {object} utils.ErrorResponse{error=utils.APIError}
// @Failure 401 {object} utils.ErrorResponse{error=utils.APIError}
// @Failure 500 {object} utils.ErrorResponse{error=utils.APIError}
// @Router /v1/auth/verify-otp [post]
func (a *UserAuthController) VerifyOTP(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	var req dto.OTPVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("invalid request body", zap.Error(err))
		utils.JSONError(c, 400, "ERROR_BAD_REQUEST", "invalid request body")
		return
	}
	if req.Email == "" {
		log.Warn("email is empty in request")
		utils.JSONError(c, 400, "ERROR_BAD_REQUEST", "email is required")
		return
	}
	if req.OTP == "" {
		log.Warn("otp is empty in request")
		utils.JSONError(c, 400, "ERROR_BAD_REQUEST", "otp is required")
		return
	}
	tokens, err := a.IAuthInteractor.VerifyOTP(ctx, req)
	if err != nil {
		log.Error("error verifying otp", zap.Error(err))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "failed to verify otp")
		return
	}

	utils.JSONSuccess(c, 200, tokens)
}


// GenerateDummyTokens godoc
// @Summary Generate dummy tokens
// @Description Generates dummy access and refresh tokens for a given user ID (for testing purposes)
// @Tags auth
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.OTPVerificationResponse}
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/auth/dummy-tokens [get]
func (a *UserAuthController) GenerateDummyTokens(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	tokens, err := a.IAuthInteractor.GenerateDummyTokens(int64(userId))
	if err != nil {
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "failed to verify otp")
		return
	}
	utils.JSONSuccess(c, 200, tokens)
}
