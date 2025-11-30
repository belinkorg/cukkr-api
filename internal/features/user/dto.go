package user

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password    string `json:"password" validate:"required,min=8" example:"password"`
	FullName    string `json:"full_name" validate:"required,min=3,max=255" example:"John Doe"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,min=10,max=20" example:"+628123456789"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password string `json:"password" validate:"required" example:"password"`
}

type UpdateProfileRequest struct {
	FullName    string `json:"full_name" validate:"omitempty,min=3,max=255" example:"John Doe Updated"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,min=10,max=20" example:"+628123456789"`
	Address     string `json:"address" example:"Jl. Sudirman No. 123"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required" example:"OldPassword123!"`
	NewPassword string `json:"new_password" validate:"required,min=8" example:"NewPassword123!"`
}

type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID          string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email       string  `json:"email" example:"john.doe@example.com"`
	FullName    string  `json:"full_name" example:"John Doe"`
	PhoneNumber string  `json:"phone_number" example:"+628123456789"`
	Address     string  `json:"address" example:"Jl. Sudirman No. 123"`
	IsActive    bool    `json:"is_active" example:"true"`
	IsVerified  bool    `json:"is_verified" example:"false"`
	LastLoginAt *string `json:"last_login_at" example:"2024-01-15T10:30:00Z"`
	CreatedAt   string  `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt   string  `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}
