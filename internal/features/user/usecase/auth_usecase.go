package usecase

import (
	"bLink-app/internal/features/user"
	"bLink-app/pkg/errors"
	"bLink-app/pkg/jwt"
	"bLink-app/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type authUsecase struct {
	jwtService *jwt.JWTService
	logger     *logger.Logger
}

func NewAuthUsecase(jwtService *jwt.JWTService, logger *logger.Logger) user.AuthUsecase {
	return &authUsecase{
		jwtService: jwtService,
		logger:     logger,
	}
}

func (a *authUsecase) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.WithError(err).Error("Failed to hash password")
		return "", errors.New(http.StatusInternalServerError, "Internal server error")
	}
	return string(hashedPassword), nil
}

func (a *authUsecase) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (a *authUsecase) GenerateToken(userID, email string) (string, error) {
	token, err := a.jwtService.GenerateToken(userID, email)
	if err != nil {
		a.logger.WithError(err).Error("Failed to generate token")
		return "", err
	}
	return token, nil
}
