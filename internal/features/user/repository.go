package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context, limit, offset int) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	UpdateLastLogin(ctx context.Context, id string) error
	CountAll(ctx context.Context) (int64, error)
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

func (r *repository) Create(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	r.invalidateCache(ctx)
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) FindAll(ctx context.Context, limit, offset int) ([]User, error) {
	cacheKey := fmt.Sprintf("users:list:%d:%d", limit, offset)

	// Try cache first
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var users []User
		if err := json.Unmarshal([]byte(cached), &users); err == nil {
			return users, nil
		}
	}

	// If not in cache, query DB
	var users []User
	query := r.db.WithContext(ctx).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	// Cache the result
	data, _ := json.Marshal(users)
	r.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	return users, nil
}

func (r *repository) Update(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	r.invalidateCache(ctx)
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&User{}, "id = ?", id).Error; err != nil {
		return err
	}
	r.invalidateCache(ctx)
	return nil
}

func (r *repository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&User{}).
		Where("id = ?", id).
		Update("last_login_at", now).Error; err != nil {
		return err
	}
	r.invalidateCache(ctx)
	return nil
}

func (r *repository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) invalidateCache(ctx context.Context) {
	// Invalidate all user-related cache
	r.redis.Del(ctx, "users:list:*")
}
