package user

import "context"

type OTPUsecase interface {
	GenerateAndSendOTP(ctx context.Context, userID string, email string) error
	VerifyOTP(ctx context.Context, userID string, otp string) error
}

type EmailUsecase interface {
	SendOTPEmail(email, otp string) error
}

type AuthUsecase interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
	GenerateToken(userID, email string) (string, error)
}
