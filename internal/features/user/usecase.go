package user

import (
	"context"
	"cukkr-app/pkg/errors"
	"cukkr-app/pkg/logger"
	"net/http"
	"time"
)

type UserUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) error
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	GetProfile(ctx context.Context, userID string) (*UserResponse, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]UserResponse, int64, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*UserResponse, error)
	ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error
	DeleteAccount(ctx context.Context, userID string) error
	SendOTPVerification(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, userID string, req *VerifyOTPRequest) (*UserResponse, error)
}

type usecase struct {
	repo         UserRepository
	authUsecase  AuthUsecase
	otpUsecase   OTPUsecase
	emailUsecase EmailUsecase
	logger       *logger.Logger
}

func NewUserUsecase(
	repo UserRepository,
	authService AuthUsecase,
	otpService OTPUsecase,
	emailService EmailUsecase,
	logger *logger.Logger,
) UserUsecase {
	return &usecase{
		repo:         repo,
		authUsecase:  authService,
		otpUsecase:   otpService,
		emailUsecase: emailService,
		logger:       logger,
	}
}

func (u *usecase) Register(ctx context.Context, req *RegisterRequest) error {
	existingUser, err := u.repo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		u.logger.WithField("email", req.Email).Warn("Email already registered")
		return errors.New(http.StatusConflict, "Email already registered")
	}

	existingPhoneUser, err := u.repo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err == nil && existingPhoneUser != nil {
		u.logger.WithField("phone_number", req.PhoneNumber).Warn("Phone number already registered")
		return errors.New(http.StatusConflict, "Phone number already registered")
	}

	hashedPassword, err := u.authUsecase.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &User{
		Email:       req.Email,
		Password:    hashedPassword,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		IsActive:    true,
		IsVerified:  false,
	}

	// Save to Redis for a while (TTL 10 minutes)
	if err := u.repo.SavePendingUser(ctx, user, 10*time.Minute); err != nil {
		u.logger.WithError(err).Error("Failed to save pending user to Redis")
		return errors.New(http.StatusInternalServerError, "Registration failed")
	}

	// Generate and send OTP
	if err := u.otpUsecase.GenerateAndSendOTP(ctx, user.ID, user.Email); err != nil {
		u.logger.WithError(err).Error("Failed to send OTP")
		return errors.New(http.StatusInternalServerError, "Failed to send OTP verification")
	}

	u.logger.WithField("email", user.Email).Info("User registered, pending OTP verification")
	return nil
}

func (u *usecase) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
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

	return &LoginResponse{
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

func (u *usecase) VerifyOTP(ctx context.Context, userID string, req *VerifyOTPRequest) (*UserResponse, error) {
	if err := u.otpUsecase.VerifyOTP(ctx, userID, req); err != nil {
		return nil, err
	}

	// Get user from Redis
	user, err := u.repo.GetPendingUser(ctx, userID)
	if err != nil {
		u.logger.WithField("user_id", userID).Error("Pending user not found in Redis")
		return nil, errors.New(http.StatusNotFound, "Registration session expired")
	}

	// Save ke DB
	user.IsVerified = true
	if err := u.repo.Create(ctx, user); err != nil {
		u.logger.WithError(err).Error("Failed to create verified user")
		return nil, errors.New(http.StatusInternalServerError, "Failed to complete registration")
	}

	// Hapus dari Redis
	if err := u.repo.DeletePendingUser(ctx, userID); err != nil {
		u.logger.WithError(err).Warn("Failed to delete pending user from Redis")
	}

	u.logger.WithField("user_id", userID).Info("User email verified and registered successfully")
	return toUserResponse(user), nil
}

func (u *usecase) GetProfile(ctx context.Context, userID string) (*UserResponse, error) {
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

func (u *usecase) GetAllUsers(ctx context.Context, limit, offset int) ([]UserResponse, int64, error) {
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

func (u *usecase) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*UserResponse, error) {
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

func (u *usecase) ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error {
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

func toUserResponse(u *User) *UserResponse {
	var lastLogin *string
	if u.LastLoginAt != nil {
		lastLoginStr := u.LastLoginAt.Format(time.RFC3339)
		lastLogin = &lastLoginStr
	}

	return &UserResponse{
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

func toUserResponses(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = *toUserResponse(&u)
	}
	return responses
}
