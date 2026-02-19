package usecase

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"time"

	"github.com/adarsh-jaiss/the-bridge/internal/api/repository"
	"github.com/adarsh-jaiss/the-bridge/internal/api/service"
	"github.com/adarsh-jaiss/the-bridge/internal/api/types/domain"
	"github.com/adarsh-jaiss/the-bridge/internal/api/types/dto"
	"github.com/adarsh-jaiss/the-bridge/internal/auth"
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/adarsh-jaiss/the-bridge/pkg/utils"
	"go.uber.org/zap"
)

type IUserAuthInteractor interface {
	TriggerOTP(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, req dto.OTPVerificationRequest) (*dto.OTPVerificationResponse, error)
}

type userAuthInteractor struct {
	repo repository.UserAuth
}

func NewUserAuthInteractor(r repository.UserAuth) IUserAuthInteractor {
	return &userAuthInteractor{
		repo: r,
	}
}

var otpStore = make(map[string]domain.OTPRecord)

func generateOTP() (string, error) {
	otpBytes := make([]byte, 6)
	randomBytes := make([]byte, 6)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "error filling random bytes", err
	}

	for i := 0; i < len(otpBytes); i++ {
		otpBytes[i] = '0' + byte(randomBytes[i]%10)
	}
	otp := string(otpBytes)
	return otp, nil
}

func (u *userAuthInteractor) TriggerOTP(ctx context.Context, email string) error {
	log := logger.FromContext(ctx)
	otp, err := generateOTP()
	if err != nil {
		log.Error("error generating otp", zap.Error(err))
		return err
	}
	log.Debug("otp successfully generated", zap.String("email", email), zap.String("otp", otp))

	// In-memory store for OTPs (in production, use Redis or similar)
	otpStore[email] = domain.OTPRecord{
		Email:     email,
		OTP:       otp,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * 15),
	}

	emailConfig := service.EmailRequest{
		To:           []string{email},
		Subject:      "Your The Bridge verification code",
		HTMLTemplate: service.OTPTemplate(otp),
		ReplyTo:      "",
	}

	_, err = service.SendEmail(ctx, emailConfig)
	if err != nil {
		log.Error("error sending email", zap.Error(err))
		return err
	}

	return nil
}

func (u *userAuthInteractor) VerifyOTP(ctx context.Context, req dto.OTPVerificationRequest) (*dto.OTPVerificationResponse, error) {
	log := logger.FromContext(ctx)
	record, exists := otpStore[req.Email]
	if !exists {
		log.Warn("no otp found for this email", zap.String("email", req.Email))
		return nil, utils.NewUnauthorizedError("OTP not found for this email")
	}

	if time.Now().After(record.ExpiresAt) {
		log.Warn("otp exipred", zap.String("email", req.Email), zap.String("expires_at", record.ExpiresAt.String()))
		return nil, utils.NewUnauthorizedError("OTP is expired")
	}
	if record.OTP != req.OTP {
		log.Warn("no otp found for this email", zap.String("email", req.Email), zap.String("provided_otp", req.OTP), zap.String("generated_otp", record.OTP))
		return nil, utils.NewUnauthorizedError("invalid OTP")
	}

	delete(otpStore, req.Email)

	userId, err := u.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userId, err = u.repo.CreateUser(ctx, req.Email)
			if err != nil {
				log.Error("error creating user", zap.Error(err))
				return nil, err
			}
		}
		log.Error("database error while finding user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, err
	}

	accessToken, err := auth.GenerateToken(userId, auth.AccessToken)
	if err != nil {
		log.Error("error generating access token", zap.Error(err))
		return nil, err
	}

	refreshToken, err := auth.GenerateToken(userId, auth.AccessToken)
	if err != nil {
		log.Error("error generating refresh token", zap.Error(err))
		return nil, err
	}

	return &dto.OTPVerificationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
