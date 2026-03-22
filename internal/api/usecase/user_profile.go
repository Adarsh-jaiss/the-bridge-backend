package usecase

import (
	"context"

	"github.com/adarsh-jaiss/the-bridge/internal/api/repository"
	"github.com/adarsh-jaiss/the-bridge/internal/api/types/dto"
	"github.com/adarsh-jaiss/the-bridge/internal/api/types/models"
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/adarsh-jaiss/the-bridge/pkg/utils"
	"go.uber.org/zap"
)

type IUserProfileInteractor interface {
	CreateProfile(ctx context.Context, userId int64, user dto.User) error
	FollowAndUnfollow(ctx context.Context, IsUnfollowReq bool, followeeID, followerID int64) error
	UpdateBioAndProfilePic(ctx context.Context, req dto.UpdateBioAndProfilePic, userID int64) error
	FetchUserProfile(ctx context.Context, limit, cursor, userID int64) (*dto.UserProfileResponse, error)
}

type userProfileInteractor struct {
	repo repository.UserProfile
}

func NewUserProfileInteractor(r repository.UserProfile) IUserProfileInteractor {
	return &userProfileInteractor{
		repo: r,
	}
}

func (u *userProfileInteractor) CreateProfile(ctx context.Context, userID int64, user dto.User) error {
	log := logger.FromContext(ctx)

	userProfile := models.User{
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Rank:          user.Rank,
		SemanBook:     user.SemanBook,
		CompanyName:   user.CompanyName,
		Licence:       user.Licence,
		LicenceType:   user.LicenceType,
		LicenceNumber: user.LicenceNumber,
	}

	if err := u.repo.CreateProfile(ctx, userID, userProfile); err != nil {
		log.Error("error creating user profile in db", zap.Error(err))
		return err
	}
	return nil
}

func (u *userProfileInteractor) FollowAndUnfollow(ctx context.Context, IsUnfollowReq bool, followeeID, followerID int64) error {
	log := logger.FromContext(ctx)

	if IsUnfollowReq {
		if err := u.repo.UnFollow(ctx, followeeID, followerID); err != nil {
			log.Error("error adding user follow req in db", zap.Error(err))
			return err
		}
	} else {
		if followerID == followeeID {
			return utils.NewBadRequestError("cannot follow yourself")
		}
		if err := u.repo.Follow(ctx, followeeID, followerID); err != nil {
			log.Error("error adding user follow req in db", zap.Error(err))
			return err
		}
	}
	return nil
}

func (u *userProfileInteractor) UpdateBioAndProfilePic(ctx context.Context, req dto.UpdateBioAndProfilePic, userID int64) error {
	log := logger.FromContext(ctx)

	if len([]rune(req.Bio)) > 1000 {
		log.Error("bio cannot exceed 1000 characters")
		return utils.NewBadRequestError("bio cannot exceed 1000 characters")
	}

	if err := u.repo.UpdateBioAndProfilePic(ctx, req.Bio, req.ProfilePicture, userID); err != nil {
		log.Error("error updating user profile", zap.Error(err))
		return err
	}

	return nil
}

func (u *userProfileInteractor) FetchUserProfile(ctx context.Context, limit, cursor, userID int64) (*dto.UserProfileResponse, error) {
	log := logger.FromContext(ctx)
	if limit <= 0 {
		limit = 30
	}
	if cursor <= 0 {
		cursor = 0
	}

	// TODO : Cache it for 2 mins
	user, err := u.repo.FetchUserProfileSummary(ctx, userID)
	if err != nil {
		log.Error("failed to fetch user details", zap.Error(err))
		return nil, err
	}

	posts, cur, err := u.repo.FetchUserPosts(ctx, limit, cursor, userID)
	if err != nil {
		log.Error("failed to fetch posts", zap.Error(err))
		return nil, err
	}

	dtoPosts := make([]dto.Post, len(posts))
	for i, p := range posts {
		dtoPosts[i] = dto.Post{
			ID:          p.ID,
			Content:     p.Content,
			Attachments: p.Attachments,
			Likes:       p.Likes,
			Comments:    p.Comments,
			Reposts:     p.Reposts,
			Bookmarks:   p.Bookmarks,
			CreatedAt:   p.CreatedAt,
		}
	}
	profile := dto.UserProfileResponse{
		UserProfile: dto.UserProfile{
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			ProfilePicture: user.ProfilePicture,
			Rank:           user.Rank,
			Bio:            user.Bio,
			FollowersCount: user.FollowersCount,
			FollowingCount: user.FollowingCount,
			TotalPostCount: user.TotalPostCount,
			Posts:          dtoPosts,
		},
		Cursor: cur,
	}

	return &profile, nil
}
