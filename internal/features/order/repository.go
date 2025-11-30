package order

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(tx *gorm.DB, order *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
	FindByUserID(ctx context.Context, userID string) ([]Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(tx *gorm.DB, order *Order) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.Create(order).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Order, error) {
	var order Order
	if err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &order, nil
}

func (r *repository) FindByUserID(ctx context.Context, userID string) ([]Order, error) {
	var orders []Order
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id string, status string) error {
	if err := r.db.WithContext(ctx).
		Model(&Order{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
