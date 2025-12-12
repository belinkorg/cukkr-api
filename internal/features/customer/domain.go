package customer

import (
	"bLink-app/pkg/helper"
	"gorm.io/gorm"
	"time"
)

type Customer struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Email       string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(20);uniqueIndex"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"index"`
}

func (Customer) TableName() string {
	return "customers"
}

func (c *Customer) BeforeCreate(db *gorm.DB) (err error) {
	helper.SetUUIDIfEmpty(&c.ID)
	return
}
