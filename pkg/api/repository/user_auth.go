package repository

import (
	"context"
	"database/sql"
)

type UserAuth interface {
	TriggerOTP(ctx context.Context, email string) (int, error)
}

var _ UserAuth = (*UserAuthentication)(nil)

type UserAuthentication struct {
	DB *sql.DB
}

func NewUserAuthentication(db *sql.DB) *UserAuthentication {
	return &UserAuthentication{
		DB: db,
	}
}

func (u *UserAuthentication) TriggerOTP(ctx context.Context, email string) (int, error) {
	
	return 0, nil
}
