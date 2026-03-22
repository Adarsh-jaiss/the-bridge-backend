package controllers

import (
	"strconv"
	"strings"

	"github.com/adarsh-jaiss/the-bridge/internal/api/types/dto"
	"github.com/adarsh-jaiss/the-bridge/internal/api/usecase"
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/adarsh-jaiss/the-bridge/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserProfileController struct {
	IUserProfileInteractor usecase.IUserProfileInteractor
}

func NewUserProfileController(u usecase.IUserProfileInteractor) *UserProfileController {
	return &UserProfileController{
		IUserProfileInteractor: u,
	}
}

// CreateProfile godoc
// @Summary Create User Profile
// @Description Creates a user profile for the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.User true "User Profile Payload"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body or missing required fields"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - user_id not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/user/onboarding [post]
func (u *UserProfileController) CreateProfile(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	userId, ok := c.Get("user_id")
	if !ok {
		log.Error("user_id not found")
		utils.JSONError(c, 401, "ERROR_UNAUTHORIZED", "user_id is not present")
		return
	}

	userIdInt, ok := userId.(int64)
	if !ok {
		log.Error("user_id is not int64", zap.Any("user_id", userId))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "invalid user_id type")
		return
	}

	var user dto.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("invalid request body", zap.Error(err))
		utils.JSONError(c, 400, "ERROR_BAD_REQUEST", "invalid request body")
		return
	}

	log.Info("create profile request received",
		zap.Any("request_body", user),
	)

	if user.FirstName == "" || user.LastName == "" || user.Rank == "" || user.Licence == "" || user.LicenceNumber == "" || user.LicenceType == "" || user.CompanyName == "" || user.SemanBook == "" {
		log.Error("request body is missing a field", zap.Any("request_body", user))
		utils.JSONError(c, 400, "ERROR_BAD_REQUEST", "invalid request body")
		return
	}

	if err := u.IUserProfileInteractor.CreateProfile(ctx, userIdInt, user); err != nil {
		log.Error("error creating profile:", zap.Error(err))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "failed to create profile")
		return
	}

	utils.JSONSuccess(c, 201, gin.H{"message": "user profile created"})
}

// Follow a user  godoc
// @Summary Follows a user
// @Description Follows a user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Followee user ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid or missing followee ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - user_id not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/user/{id}/follow [post]
func (u *UserProfileController) Follow(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	userId, ok := c.Get("user_id")
	if !ok {
		log.Error("user_id not found")
		utils.JSONError(c, 401, "ERROR_UNAUTHORIZED", "user_id is not present")
		return
	}

	userIdInt, ok := userId.(int64)
	if !ok {
		log.Error("user_id is not int64", zap.Any("user_id", userId))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "invalid user_id type")
		return
	}

	FolloweeID := c.Param("id")
	if strings.TrimSpace(FolloweeID) == "" {
		log.Error("followee_id is empty in path paramter")
		utils.JSONError(c, 400, "BAD_REQUEST", "followee_id is required")
		return
	}

	FolloweeIDInt, err := strconv.ParseInt(FolloweeID, 10, 64)
	if err != nil {
		log.Error("error converting followee_id to int64", zap.Any("followee_id:", FolloweeID))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "error converting followee_id to int64")
		return
	}

	log.Info("follow request received for",
		zap.Any("follower_id/user_id", userIdInt),
		zap.Any("followee_id", FolloweeIDInt),
	)

	if err := u.IUserProfileInteractor.FollowAndUnfollow(ctx, false, FolloweeIDInt, userIdInt); err != nil {
		log.Error("error sending follow request", zap.Error(err))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}
	utils.JSONSuccess(c, 201, gin.H{"message": "follow request succeeded"})

}

// Unfollow a user godoc
// @Summary Unfollow a user
// @Description Unfollow a user by user ID
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Followee user ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid or missing followee ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - user_id not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/user/{id}/unfollow [DELETE]
func (u *UserProfileController) UnFollow(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	userId, ok := c.Get("user_id")
	if !ok {
		log.Error("user_id not found")
		utils.JSONError(c, 401, "ERROR_UNAUTHORIZED", "user_id is not present")
		return
	}

	userIdInt, ok := userId.(int64)
	if !ok {
		log.Error("user_id is not int64", zap.Any("user_id", userId))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "invalid user_id type")
		return
	}

	FolloweeID := c.Param("id")
	if strings.TrimSpace(FolloweeID) == "" {
		log.Error("followee_id is empty in path paramter")
		utils.NewBadRequestError("followee_id is required")
		return
	}

	FolloweeIDInt, err := strconv.ParseInt(FolloweeID, 10, 64)
	if err != nil {
		log.Error("error converting followee_id to int64", zap.Any("followee_id:", FolloweeID))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "error converting followee_id to int64")
		return
	}

	log.Info("unfollow request received for",
		zap.Any("follower_id/user_id", userIdInt),
		zap.Any("followee_id", FolloweeIDInt),
	)

	if err := u.IUserProfileInteractor.FollowAndUnfollow(ctx, true, FolloweeIDInt, userIdInt); err != nil {
		log.Error("error sending follow request", zap.Error(err))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "unfollow request failed")
		return
	}

	utils.JSONSuccess(c, 201, gin.H{"message": "unfollow request succeeded"})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user's bio and profile picture
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateBioAndProfilePic true "Update profile payload"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - user_id not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/user/profile [patch]
func (u *UserProfileController) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	userId, ok := c.Get("user_id")
	if !ok {
		log.Error("user_id not found")
		utils.JSONError(c, 401, "ERROR_UNAUTHORIZED", "user_id is not present")
		return
	}

	userIdInt, ok := userId.(int64)
	if !ok {
		log.Error("user_id is not int64", zap.Any("user_id", userId))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "invalid user_id type")
		return
	}

	var user dto.UpdateBioAndProfilePic
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("invalid request body", zap.Error(err))
		utils.JSONError(c, 400, "ERROR_BAD_REQUEST", "invalid request body")
		return
	}

	log.Info("update  profile request received",
		zap.Any("request_body", user),
	)

	if err := u.IUserProfileInteractor.UpdateBioAndProfilePic(ctx, user, userIdInt); err != nil {
		log.Error("error updating user profile picture and bio", zap.Error(err))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "failed to update profile")
		return
	}

	utils.JSONSuccess(c, 200, gin.H{
		"message": "profile updated successfully",
	})

}

// GetUserProfile godoc
// @Summary Get user profile
// @Description Fetches the authenticated user's profile along with paginated posts (cursor-based)
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int true "Number of items to fetch"
// @Param cursor query int true "Cursor for pagination (last seen ID)"
// @Success 200 {object} utils.SuccessResponse{data=dto.UserProfileResponse}
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - user_id not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /v1/user/profile [get]
func (u *UserProfileController) GetUserProfile(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	userId, ok := c.Get("user_id")
	if !ok {
		log.Error("user_id not found")
		utils.JSONError(c, 401, "ERROR_UNAUTHORIZED", "user_id is not present")
		return
	}
	userIdInt, ok := userId.(int64)
	if !ok {
		log.Error("user_id is not int64", zap.Any("user_id", userId))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "invalid user_id type")
		return
	}

	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		log.Error("error converting limit to int64", zap.Any("limit:", userId))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "error converting limit to int64")
		return
	}

	cursor, err := strconv.ParseInt(c.Query("cursor"), 10, 64)
	if err != nil {
		log.Error("error converting cursor to int64", zap.Any("cursor:", userId))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "error converting cursor to int64")
		return
	}

	profile, err := u.IUserProfileInteractor.FetchUserProfile(ctx, limit, cursor, userIdInt)
	if err != nil {
		log.Error("failed to fetch user profile", zap.Error(err))
		utils.JSONError(c, 500, "INTERNAL_SERVER_ERROR", "failed to fetch user profile")
		return
	}

	utils.JSONSuccess(c, 200, profile)

}
