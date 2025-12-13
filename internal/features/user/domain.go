package user

import (
	"cukurly-app/pkg/helper"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          string     `json:"id" gorm:"type:uuid;primaryKey"`
	Email       string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password    string     `json:"-" gorm:"type:varchar(255);not null"`
	FullName    string     `json:"full_name" gorm:"type:varchar(255);not null"`
	PhoneNumber string     `json:"phone_number" gorm:"type:varchar(20);uniqueIndex"`
	Address     string     `json:"address" gorm:"type:text"`
	Bio         string     `json:"bio" gorm:"type:text"`
	PhotoURL    string     `json:"photo_url" gorm:"type:varchar(255)"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	IsVerified  bool       `json:"is_verified" gorm:"default:false"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   time.Time  `json:"deleted_at" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	helper.SetUUIDIfEmpty(&u.ID)
	return
}
