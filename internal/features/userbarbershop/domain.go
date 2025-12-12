package userbarbershop

import (
	"time"
)

type UserBarbershopRole string

const (
	RoleOwner  UserBarbershopRole = "owner"
	RoleBarber UserBarbershopRole = "barber"
)

type UserBarbershop struct {
	UserID       string             `json:"user_id" gorm:"type:uuid;primaryKey"`
	BarbershopID string             `json:"barbershop_id" gorm:"type:uuid;primaryKey"`
	Role         UserBarbershopRole `json:"role" gorm:"type:varchar(30);not null"`
	CreatedAt    time.Time          `gorm:"autoCreateTime" json:"created_at"`
}

func (UserBarbershop) TableName() string {
	return "user_barbershops"
}

func (r UserBarbershopRole) IsValid() bool {
	switch r {
	case RoleOwner, RoleBarber:
		return true
	}
	return false
}
