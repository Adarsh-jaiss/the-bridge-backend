package domain

import "time"

type OTPRecord struct {
	Email     string
	OTP       string
	CreatedAt time.Time
	ExpiresAt time.Time
}
