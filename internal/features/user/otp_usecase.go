package user

import (
	"context"
	"cukkr-app/pkg/errors"
	"cukkr-app/pkg/logger"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type OTPUsecase interface {
	GenerateAndSendOTP(ctx context.Context, userID string, email string) error
	VerifyOTP(ctx context.Context, userID string, req *VerifyOTPRequest) error
}

type otpUsecase struct {
	repo         UserRepository
	logger       *logger.Logger
	emailUsecase EmailUsecase
}

func NewOTPUsecase(repo UserRepository, logger *logger.Logger, emailUsecase EmailUsecase) OTPUsecase {
	return &otpUsecase{
		repo:         repo,
		logger:       logger,
		emailUsecase: emailUsecase,
	}
}

func (o *otpUsecase) generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (o *otpUsecase) GenerateAndSendOTP(ctx context.Context, userID string, email string) error {
	otpCode := o.generateOTP()
	expiry := 5 * time.Minute // 5 minutes in seconds

	if err := o.repo.SaveOTPToRedis(ctx, userID, otpCode, expiry); err != nil {
		o.logger.WithError(err).Error("Failed to save OTP to Redis")
		return errors.New(http.StatusInternalServerError, "Failed to save OTP to Redis")
	}

	if err := o.emailUsecase.SendOTPEmail(email, otpCode); err != nil {
		o.logger.WithError(err).Error("Failed to send OTP email")
		return errors.New(http.StatusInternalServerError, "Failed to send OTP email")
	}

	o.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"email":   email,
	}).Info("OTP generated and sent successfully")

	return nil
}

func (o *otpUsecase) VerifyOTP(ctx context.Context, userID string, req *VerifyOTPRequest) error {
	storedOTP, err := o.repo.GetOTPFromRedis(ctx, userID)
	if err != nil {
		o.logger.WithError(err).Error("Failed to get OTP from Redis")
		return errors.New(http.StatusInternalServerError, "Internal server error")
	}

	if storedOTP == "" {
		return errors.New(http.StatusBadRequest, "OTP expired or not found")
	}

	if storedOTP != req.OTPCode {
		return errors.New(http.StatusBadRequest, "Invalid OTP")
	}

	if err := o.repo.DeleteOTPFromRedis(ctx, userID); err != nil {
		o.logger.WithField("user_id", userID).Warn("Failed to delete OTP from Redis")
	}

	o.logger.WithField("user_id", userID).Info("OTP verified successfully")
	return nil
}
