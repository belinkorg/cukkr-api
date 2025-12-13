package order

import (
	"cukurly-app/pkg/helper"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         string      `json:"id" gorm:"type:uuid;primaryKey"`
	UserID     string      `json:"user_id" gorm:"type:uuid;not null"`
	TotalPrice float64     `json:"total_price" gorm:"type:decimal(10,2);not null"`
	Status     string      `json:"status" gorm:"type:varchar(50);default:'pending'"`
	Items      []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt  time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

type OrderItem struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey"`
	OrderID   string    `json:"order_id" gorm:"type:uuid;not null"`
	ProductID string    `json:"product_id" gorm:"type:uuid;not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Order) TableName() string {
	return "orders"
}

func (OrderItem) TableName() string {
	return "order_items"
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	helper.SetUUIDIfEmpty(&o.ID)
	return
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	helper.SetUUIDIfEmpty(&oi.ID)
	return
}

// Status constants
const (
	StatusPending    = "Pending"
	StatusProcessing = "Processing"
	StatusCompleted  = "Completed"
	StatusCancelled  = "Cancelled"
)
