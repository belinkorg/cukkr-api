package barbershop

import (
	"bLink-app/pkg/helper"
	"gorm.io/gorm"
	"time"
)

type Barbershop struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey"`
	Slug        string    `json:"slug" gorm:"type:varchar(50);uniqueIndex;not null"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Email       string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(20);uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	Address     string    `json:"address" gorm:"type:text"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"index"`
}

func (Barbershop) TableName() string {
	return "barbershops"
}

func (b *Barbershop) BeforeCreate(tx *gorm.DB) (err error) {
	helper.SetUUIDIfEmpty(&b.ID)
	return
}
