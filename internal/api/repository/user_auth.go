package repository

import (
	"context"
	"database/sql"

	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"go.uber.org/zap"
)

type UserAuth interface {
	FindUserByEmail(ctx context.Context, email string) (userId int64, err error)
	CreateUser(ctx context.Context, email string) (userId int64, err error)
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

func (u *UserAuthentication) FindUserByEmail(ctx context.Context, email string) (userId int64, err error) {
	log := logger.FromContext(ctx)
	query := `SELECT id FROM users where email=$1`

	if err = u.DB.QueryRowContext(ctx, query, email).Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			log.Debug("user not found", zap.String("email", email))
			return 0, sql.ErrNoRows
		}
		log.Error("failed to query user by email", zap.Error(err))
		return 0, err
	}

	return userId, nil
}

func (u *UserAuthentication) CreateUser(ctx context.Context, email string) (userId int64, err error) {
	log := logger.FromContext(ctx)
	query := `INSERT INTO users (email) VALUES ($1) RETURNING id`

	if err = u.DB.QueryRowContext(ctx, query, email).Scan(&userId); err != nil {
		log.Error("failed to create user", zap.Error(err))
		return 0, err
	}

	return userId, nil
}
