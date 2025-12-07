package repository

import (
	"bLink-app/internal/features/user/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
	FindAll(ctx context.Context, limit, offset int) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	UpdateLastLogin(ctx context.Context, id string) error
	CountAll(ctx context.Context) (int64, error)

	SaveOTPToRedis(ctx context.Context, userID, otpCode string, expiry time.Duration) error
	GetOTPFromRedis(ctx context.Context, userID string) (string, error)
	DeleteOTPFromRedis(ctx context.Context, userID string) error
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

func (r *repository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	r.invalidateCache(ctx)
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindAll(ctx context.Context, limit, offset int) ([]domain.User, error) {
	cacheKey := fmt.Sprintf("users:list:%d:%d", limit, offset)

	// Try cache first
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var users []domain.User
		if err := json.Unmarshal([]byte(cached), &users); err == nil {
			return users, nil
		}
	}

	// If not in cache, query DB
	var users []domain.User
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

func (r *repository) Update(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	r.invalidateCache(ctx)
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	r.invalidateCache(ctx)
	return nil
}

func (r *repository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", id).
		Update("last_login_at", now).Error; err != nil {
		return err
	}
	r.invalidateCache(ctx)
	return nil
}

func (r *repository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) invalidateCache(ctx context.Context) {
	// Invalidate all user-related cache
	r.redis.Del(ctx, "users:list:*")
}

func (r *repository) SaveOTPToRedis(ctx context.Context, userID, otpCode string, expiry time.Duration) error {
	key := fmt.Sprintf("otp:%s", userID)
	return r.redis.Set(ctx, key, otpCode, expiry).Err()
}

func (r *repository) GetOTPFromRedis(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("otp:%s", userID)
	val, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *repository) DeleteOTPFromRedis(ctx context.Context, userID string) error {
	key := fmt.Sprintf("otp:%s", userID)
	return r.redis.Del(ctx, key).Err()
}
