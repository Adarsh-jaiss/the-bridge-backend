package usecase

import (
	"context"

	"github.com/adarsh-jaiss/the-bridge/pkg/api/repository"
)

type IUserAuthInteractor interface {
	TriggerOTP(ctx context.Context, email string) (int, error)
}

type userAuthInteractor struct {
	repo repository.UserAuth
}

func NewUserAuthInteractor(r repository.UserAuth) userAuthInteractor {
	return userAuthInteractor{
		repo: r,
	}
}

func (u *userAuthInteractor) TriggerOTP(ctx context.Context, email string) (int, error) {
	// call repository
	return u.repo.TriggerOTP(ctx, email)
}
