package auth

import (
	"context"
	"strconv"
	"time"

	"github.com/adarsh-jaiss/the-bridge/pkg/config"
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/adarsh-jaiss/the-bridge/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type CustomClaims struct {
	UserID int64     `json:"user_id"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64, ttype TokenType) (string, error) {
	cfg := config.Get()
	var secret string

	switch ttype {
	case AccessToken:
		secret = cfg.JWTAccessTokenSecret
	case RefreshToken:
		secret = cfg.JWTRefreshTokenSecret
	default:
		return "", utils.NewUnauthorizedError("invalid token type")
	}
	exp := time.Now().Add(time.Minute * 30)
	if ttype == RefreshToken {
		exp = time.Now().Add(time.Hour * 24 * 7)
	}

	claims := CustomClaims{
		UserID: userID,
		Type:   ttype,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    "the-bridge",
			Subject:   strconv.FormatInt(userID, 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(ctx context.Context,tokenStr string,ttype TokenType,) (*CustomClaims, error) {

	log := logger.FromContext(ctx)
	cfg := config.Get()

	var secret string
	switch ttype {
	case AccessToken:
		secret = cfg.JWTAccessTokenSecret
	case RefreshToken:
		secret = cfg.JWTRefreshTokenSecret
	default:
		return nil, utils.NewUnauthorizedError("invalid token type")
	}

	var claims CustomClaims

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&claims,
		func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Warn("unexpected signing method",
					zap.String("alg", t.Method.Alg()),
				)
				return nil, utils.NewUnauthorizedError("unexpected signing method")
			}

			return []byte(secret), nil
		},
	)
	if err != nil {
		log.Warn("token parsing failed", zap.Error(err))
		return nil, err
	}

	if !token.Valid {
		return nil, utils.NewUnauthorizedError("invalid token")
	}

	if claims.Type != ttype {
		log.Warn("token type mismatch",
			zap.String("expected", string(ttype)),
			zap.String("received", string(claims.Type)),
		)
		return nil, utils.NewUnauthorizedError("invalid token type")
	}

	return &claims, nil
}
