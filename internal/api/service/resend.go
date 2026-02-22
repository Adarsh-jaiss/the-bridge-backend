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
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8" />
  <title>Verification Code</title>
</head>
<body style="margin:0; padding:0; background-color:#f2f2f2; font-family: Arial, sans-serif;">

  <table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:#f2f2f2; padding:40px 0;">
    <tr>
      <td align="center">

        <table width="600" cellpadding="0" cellspacing="0" border="0" 
          style="background-color:#ffffff; border-radius:8px; padding:40px 30px; text-align:center;">



          <tr>
            <td style="padding:20px 0;">
              <img src="https://assets.zocket.com/testing/media/agents/unnamed.png"
                   alt="Verification Illustration"
                   width="150"
                   style="display:block; margin:0 auto;" />
            </td>
          </tr>

          <tr>
            <td style="padding:10px 0;">
              <h1 style="margin:0; font-size:26px; color:#333333;">
                Verification Code
              </h1>
            </td>
          </tr>

          <tr>
            <td style="padding:10px 0;">
              <p style="margin:0; font-size:16px; color:#555555;">
                Please use the verification code below to sign in.
              </p>
            </td>
          </tr>

          <tr>
            <td style="padding:15px 0;">
              <p style="margin:0; font-size:28px; font-weight:bold; color:#000000; letter-spacing:4px;">
                ` + otp + `
              </p>
              <p style="margin-top:10px; font-size:14px; color:#777777;">
                This code expires in 10 minutes.
              </p>
            </td>
          </tr>

          <tr>
            <td style="padding-top:25px;">
              <p style="margin:0; font-size:14px; color:#777777; line-height:1.6;">
                If you didn't request an OTP, please contact our support team immediately.
              </p>
              <p style="margin:5px 0 0 0; font-size:14px;">
                <a href="mailto:support@thebridge.com" 
                   style="color:#1a73e8; text-decoration:none; font-weight:600;">
                   support@thebridge.com
                </a>
              </p>
            </td>
          </tr>

          <tr>
            <td style="padding-top:35px;">
              <p style="margin:0; font-size:13px; color:#999999;">
                Copyright © 2026 | The bridge | All rights reserved
              </p>
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`
}
