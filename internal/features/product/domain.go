package product

import (
	"cukkr-app/pkg/helper"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock       int       `json:"stock" gorm:"default:0"`
	Category    string    `json:"category" gorm:"type:varchar(100)"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Product) TableName() string {
	return "products"
}

// BeforeCreate hook untuk generate UUID otomatis (sebagai string)
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	helper.SetUUIDIfEmpty(&p.ID)
	return
}
