package service

import (
	"context"

	"github.com/adarsh-jaiss/the-bridge/pkg/config"
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/resend/resend-go/v3"
	"go.uber.org/zap"
)

type EmailRequest struct {
	To           []string
	Subject      string
	HTMLTemplate string
	ReplyTo      string
}

func SendEmail(ctx context.Context, req EmailRequest) (string, error) {
	log := logger.FromContext(ctx)
	cfg := config.Get()
	client := resend.NewClient(cfg.ResendAPIKey)

	params := &resend.SendEmailRequest{
		From:    "The Bridge <no-reply@support.bluviaglobal.com>",
		To:      req.To,
		Subject: req.Subject,
		Html:    req.HTMLTemplate,
		ReplyTo: req.ReplyTo,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Error("failed to send email via resend", zap.Error(err))
		return "", err
	}

	log.Info("email sent successfully",
		zap.String("email_sent_id", sent.Id),
		zap.String("toemail", req.To[0]),
	)
	return sent.Id, nil
}

func OTPTemplate(otp string) string {
	return `
		<h2>Your Verification Code</h2>
		<p>Your verification code is:</p>
		<h1>` + otp + `</h1>
		<p>This code expires in 10 minutes.</p>
	`
}
