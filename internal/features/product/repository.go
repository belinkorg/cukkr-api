package product

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	Create(tx *gorm.DB, product *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context) ([]Product, error)
	Update(tx *gorm.DB, product *Product) error
	Delete(tx *gorm.DB, id string) error
	UpdateStock(tx *gorm.DB, id string, quantity int) (int64, error)
}

type repository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewRepository(db *gorm.DB, redis *redis.Client) Repository {
	return &repository{
		db:    db,
		redis: redis,
	}
}

func (r *repository) Create(tx *gorm.DB, product *Product) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.Create(product).Error; err != nil {
		return err
	}

	r.invalidateCache(context.Background())
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Product, error) {
	cacheKey := fmt.Sprintf("product:%s", id)

	// Try cache first
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var product Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return &product, nil
		}
	}

	// If not in cache, query DB
	var product Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	// Cache the result
	data, _ := json.Marshal(product)
	r.redis.Set(ctx, cacheKey, data, 10*time.Minute)

	return &product, nil
}

func (r *repository) FindAll(ctx context.Context) ([]Product, error) {
	cacheKey := "products:all"

	// Try cache first
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var products []Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			return products, nil
		}
	}

	// If not in cache, query DB
	var products []Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}

	// Cache the result
	data, _ := json.Marshal(products)
	r.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	return products, nil
}

func (r *repository) Update(tx *gorm.DB, product *Product) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.Save(product).Error; err != nil {
		return err
	}

	r.invalidateCache(context.Background())
	return nil
}

func (r *repository) Delete(tx *gorm.DB, id string) error {
	if tx == nil {
		tx = r.db
	}

	result := tx.Where("id = ?", id).Delete(&Product{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	r.invalidateCache(context.Background())
	return nil
}

func (r *repository) UpdateStock(tx *gorm.DB, id string, quantity int) (int64, error) {
	if tx == nil {
		tx = r.db
	}

	result := tx.Model(&Product{}).
		Where("id = ? AND stock >= ?", id, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))

	if result.Error != nil {
		return 0, result.Error
	}

	r.invalidateCache(context.Background())
	return result.RowsAffected, nil
}

func (r *repository) invalidateCache(ctx context.Context) {
	r.redis.Del(ctx, "products:all")
}
