package user

import (
	"cukkr-app/pkg/logger"
	"fmt"
	"gopkg.in/mail.v2"
	"os"
)

type EmailUsecase interface {
	SendOTPEmail(email, otp string) error
}

type emailUsecase struct {
	logger *logger.Logger
}

func NewEmailUsecase(logger *logger.Logger) EmailUsecase {
	return &emailUsecase{
		logger: logger,
	}
}

func (e *emailUsecase) SendOTPEmail(email, otp string) error {
	m := mail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_EMAIL"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Email Verification - OTP Code")
	m.SetBody("text/html", e.buildOTPEmailBody(otp))

	d := mail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_EMAIL"),
		os.Getenv("SMTP_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		e.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"email": email,
		}).Error("Failed to send OTP email")
		return err
	}

	e.logger.WithField("email", email).Info("OTP email sent successfully")
	return nil
}

func (e *emailUsecase) buildOTPEmailBody(otp string) string {
	return fmt.Sprintf(`
  <html>
  <body style="font-family: Arial, sans-serif;">
   <div style="max-width: 600px; margin: 0 auto; background-color: #f9f9f9; padding: 20px; border-radius: 5px;">
    <h2>Email Verification</h2>
    <p>Your OTP code is valid for 5 minutes</p>
    <div style="background-color: #f0f0f0; padding: 15px; text-align: center; font-size: 32px; font-weight: bold; letter-spacing: 5px; margin: 20px 0; border-radius: 5px;">%s</div>
    <p>If you didn't request this code, please ignore this email.</p>
   </div>
  </body>
  </html>
 `, otp)
}
