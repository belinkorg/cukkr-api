package usecase

import (
	"bLink-app/internal/features/user"
	"bLink-app/internal/features/user/domain"
	"bLink-app/internal/features/user/model"
	"bLink-app/internal/features/user/repository"
	"bLink-app/pkg/errors"
	"bLink-app/pkg/logger"
	"context"
	"net/http"
	"time"
)

type Usecase interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.UserResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	GetProfile(ctx context.Context, userID string) (*model.UserResponse, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]model.UserResponse, int64, error)
	UpdateProfile(ctx context.Context, userID string, req *model.UpdateProfileRequest) (*model.UserResponse, error)
	ChangePassword(ctx context.Context, userID string, req *model.ChangePasswordRequest) error
	DeleteAccount(ctx context.Context, userID string) error
	SendOTPVerification(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, userID, otpCode string) (*model.UserResponse, error)
}

type usecase struct {
	repo         repository.Repository
	authUsecase  user.AuthUsecase
	otpUsecase   user.OTPUsecase
	emailUsecase user.EmailUsecase
	logger       *logger.Logger
}

func NewUsecase(
	repo repository.Repository,
	authService user.AuthUsecase,
	otpService user.OTPUsecase,
	emailService user.EmailUsecase,
	logger *logger.Logger,
) Usecase {
	return &usecase{
		repo:         repo,
		authUsecase:  authService,
		otpUsecase:   otpService,
		emailUsecase: emailService,
		logger:       logger,
	}
}

func (u *usecase) Register(ctx context.Context, req *model.RegisterRequest) (*model.UserResponse, error) {
	existingUser, err := u.repo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		u.logger.WithField("email", req.Email).Warn("Email already registered")
		return nil, errors.New(http.StatusConflict, "Email already registered")
	}

	existingPhoneUser, err := u.repo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err == nil && existingPhoneUser != nil {
		u.logger.WithField("phone_number", req.PhoneNumber).Warn("Phone number already registered")
		return nil, errors.New(http.StatusConflict, "Phone number already registered")
	}

	hashedPassword, err := u.authUsecase.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:       req.Email,
		Password:    hashedPassword,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		IsActive:    true,
		IsVerified:  false,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		u.logger.WithError(err).Error("Failed to create user")
		return nil, errors.New(http.StatusInternalServerError, "Internal server error")
	}

	u.logger.WithFields(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User registered successfully")

	return toUserResponse(user), nil
}

func (u *usecase) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		u.logger.WithField("email", req.Email).Warn("Email not found")
		return nil, errors.New(http.StatusUnauthorized, "Email not found")
	}

	if !user.IsActive {
		u.logger.WithField("user_id", user.ID).Warn("Login attempt for inactive user")
		return nil, errors.New(http.StatusForbidden, "User is inactive")
	}

	if err := u.authUsecase.VerifyPassword(user.Password, req.Password); err != nil {
		u.logger.WithField("user_id", user.ID).Warn("Invalid password")
		return nil, errors.New(http.StatusUnauthorized, "Invalid email or password")
	}

	token, err := u.authUsecase.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	if err := u.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		u.logger.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": user.ID,
		}).Warn("Failed to update last login")
	}

	u.logger.WithField("user_id", user.ID).Info("User logged in successfully")

	return &model.LoginResponse{
		Token: token,
		User:  *toUserResponse(user),
	}, nil
}

func (u *usecase) SendOTPVerification(ctx context.Context, email string) error {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		u.logger.WithField("email", email).Warn("User not found for OTP")
		return errors.New(http.StatusNotFound, "User not found")
	}

	if user.IsVerified {
		return errors.New(http.StatusConflict, "User already verified")
	}

	return u.otpUsecase.GenerateAndSendOTP(ctx, user.ID, email)
}

func (u *usecase) VerifyOTP(ctx context.Context, userID, otpCode string) (*model.UserResponse, error) {
	if err := u.otpUsecase.VerifyOTP(ctx, userID, otpCode); err != nil {
		return nil, err
	}

	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.IsVerified = true
	if err := u.repo.Update(ctx, user); err != nil {
		u.logger.WithError(err).Error("Failed to update user verification")
		return nil, err
	}

	u.logger.WithField("user_id", userID).Info("User email verified successfully")

	return toUserResponse(user), nil
}

func (u *usecase) GetProfile(ctx context.Context, userID string) (*model.UserResponse, error) {
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		u.logger.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		}).Error("Failed to get user profile")
		return nil, err
	}

	return toUserResponse(user), nil
}

func (u *usecase) GetAllUsers(ctx context.Context, limit, offset int) ([]model.UserResponse, int64, error) {
	users, err := u.repo.FindAll(ctx, limit, offset)
	if err != nil {
		u.logger.WithError(err).Error("Failed to get all users")
		return nil, 0, err
	}

	total, err := u.repo.CountAll(ctx)
	if err != nil {
		u.logger.WithError(err).Error("Failed to count users")
		return nil, 0, err
	}

	return toUserResponses(users), total, nil
}

func (u *usecase) UpdateProfile(ctx context.Context, userID string, req *model.UpdateProfileRequest) (*model.UserResponse, error) {
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New(http.StatusNotFound, "User not found")
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if err := u.repo.Update(ctx, user); err != nil {
		u.logger.WithError(err).Error("Failed to update user profile")
		return nil, err
	}

	u.logger.WithFields(map[string]interface{}{
		"user_id": user.ID,
	}).Info("User profile updated successfully")

	return toUserResponse(user), nil
}

func (u *usecase) ChangePassword(ctx context.Context, userID string, req *model.ChangePasswordRequest) error {
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return errors.New(http.StatusNotFound, "User not found")
	}

	if err := u.authUsecase.VerifyPassword(user.Password, req.OldPassword); err != nil {
		u.logger.WithField("user_id", userID).Warn("Change password attempt with wrong old password")
		return errors.New(http.StatusBadRequest, "Wrong old password")
	}

	hashedPassword, err := u.authUsecase.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	if err := u.repo.Update(ctx, user); err != nil {
		u.logger.WithError(err).Error("Failed to update password")
		return err
	}

	u.logger.WithField("user_id", user.ID).Info("Password changed successfully")

	return nil
}

func (u *usecase) DeleteAccount(ctx context.Context, userID string) error {
	if _, err := u.repo.FindByID(ctx, userID); err != nil {
		return errors.New(http.StatusNotFound, "User not found")
	}

	if err := u.repo.Delete(ctx, userID); err != nil {
		u.logger.WithError(err).Error("Failed to delete user account")
		return err
	}

	u.logger.WithField("user_id", userID).Info("User account deleted successfully")

	return nil
}

func toUserResponse(u *domain.User) *model.UserResponse {
	var lastLogin *string
	if u.LastLoginAt != nil {
		lastLoginStr := u.LastLoginAt.Format(time.RFC3339)
		lastLogin = &lastLoginStr
	}

	return &model.UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber,
		Address:     u.Address,
		Bio:         u.Bio,
		PhotoURL:    u.PhotoURL,
		IsActive:    u.IsActive,
		IsVerified:  u.IsVerified,
		LastLoginAt: lastLogin,
		CreatedAt:   u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   u.UpdatedAt.Format(time.RFC3339),
	}
}

func toUserResponses(users []domain.User) []model.UserResponse {
	responses := make([]model.UserResponse, len(users))
	for i, u := range users {
		responses[i] = *toUserResponse(&u)
	}
	return responses
}
