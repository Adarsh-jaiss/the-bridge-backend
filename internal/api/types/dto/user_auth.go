package dto

type EmailVerficationRequest struct {
	Email string `json:"email" binding:"safe_email"`
}

type OTPVerificationRequest struct {
	Email string `json:"email" binding:"safe_email"`
	OTP   string `json:"otp" binding:"safe_otp"`
}

type OTPVerificationResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
