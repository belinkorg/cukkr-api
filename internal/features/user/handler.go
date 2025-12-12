package user

import (
	"bLink-app/pkg/helper"
	"bLink-app/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	usecase    UserUsecase
	httpHelper *helper.HTTPHandlerHelper
}

func NewHandler(usecase UserUsecase, httpHelper *helper.HTTPHandlerHelper) *Handler {
	return &Handler{
		usecase:    usecase,
		httpHelper: httpHelper,
	}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model.RegisterRequest	true	"User registration data"
//	@Success		201		{object}	response.Response{data=model.UserResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Security		BearerAuth
//	@Router			/users/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := h.httpHelper.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.usecase.Register(c.Request.Context(), &req)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to register user")
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully, please verify OTP sent to your email", nil)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body model.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=model.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := h.httpHelper.BindAndValidate(c, &req); err != nil {
		return
	}

	result, err := h.usecase.Login(c.Request.Context(), &req)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to login user")
		return
	}

	response.Success(c, http.StatusOK, "Login successful", result)
}

// VerifyEmailOTP godoc
// @Summary User Verify Email OTP
// @Description Verify email OTP and activate user account
// @Tags Users
// @Accept json
// @Produce json
// @Param	user body	model.VerifyOTPRequest true "User OTP verification data"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/verify-otp [post]
func (h *Handler) VerifyEmailOTP(c *gin.Context) {
	var req VerifyOTPRequest
	userID := c.GetString("user_id")
	if err := h.httpHelper.BindAndValidate(c, &req); err != nil {
		return
	}

	result, err := h.usecase.VerifyOTP(c.Request.Context(), userID, &req)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to verify otp")
		return
	}

	response.Success(c, http.StatusOK, "Verify OTP successful", result)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get profile of authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /users/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := c.GetString("user_id")

	result, err := h.usecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to get user profile")
		return
	}

	response.Success(c, http.StatusOK, "Profile retrieved successfully", result)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get list of all users with pagination (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{data=[]model.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	result, total, err := h.usecase.GetAllUsers(c.Request.Context(), limit, offset)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to get users")
		return
	}

	c.JSON(200, gin.H{
		"status_code": http.StatusOK,
		"message":     "Users retrieved successfully",
		"data":        result,
		"errors":      nil,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update profile of authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param profile body model.UpdateProfileRequest true "Profile data"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /users/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req UpdateProfileRequest
	if err := h.httpHelper.BindAndValidate(c, &req); err != nil {
		return
	}

	result, err := h.usecase.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to update user profile")
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", result)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password of authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param password body model.ChangePasswordRequest true "Password data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /users/change-password [post]
func (h *Handler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")

	var req ChangePasswordRequest
	if err := h.httpHelper.BindAndValidate(c, &req); err != nil {
		return
	}

	err := h.usecase.ChangePassword(c.Request.Context(), userID, &req)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to change password")
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}

// DeleteAccount godoc
// @Summary Delete user account
// @Description Delete account of authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /users/account [delete]
func (h *Handler) DeleteAccount(c *gin.Context) {
	userID := c.GetString("user_id")

	err := h.usecase.DeleteAccount(c.Request.Context(), userID)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to delete account")
		return
	}

	response.Success(c, http.StatusOK, "Account deleted successfully", nil)
}
