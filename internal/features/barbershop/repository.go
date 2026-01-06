package barbershop

import (
	"context"
	"cukkr-app/pkg/pagination"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BarbershopRepository interface {
	Create(ctx context.Context, tx *gorm.DB, barbershop *Barbershop) error
	FindByID(ctx context.Context, id string) (*Barbershop, error)
	FindBySlug(ctx context.Context, slug string) (*Barbershop, error)
	FindByUserID(ctx context.Context, userID string, params pagination.Pagination) (*pagination.PaginationResult[Barbershop], error)
	FindAll(ctx context.Context, params pagination.Pagination) (*pagination.PaginationResult[Barbershop], error)
	Update(ctx context.Context, barbershop *Barbershop) error
	Delete(ctx context.Context, id string) error
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
	CheckEmailExists(ctx context.Context, email string, excludeID string) (bool, error)
}

type repository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewRepository(db *gorm.DB, redis *redis.Client) BarbershopRepository {
	return &repository{
		db:    db,
		redis: redis,
	}
}

func (r *repository) Create(ctx context.Context, tx *gorm.DB, barbershop *Barbershop) error {
	db := r.db
	if tx != nil {
		tx = db
	}

	err := db.WithContext(ctx).Create(barbershop).Error
	if err != nil {
		return err
	}

	r.invalidateCache(context.Background())
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Barbershop, error) {
	var barbershop Barbershop
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&barbershop).Error
	if err != nil {
		return nil, err
	}
	return &barbershop, nil
}

func (r *repository) FindBySlug(ctx context.Context, slug string) (*Barbershop, error) {
	var barbershop Barbershop
	err := r.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", slug).First(&barbershop).Error
	if err != nil {
		return nil, err
	}
	return &barbershop, nil
}

func (r *repository) FindByUserID(ctx context.Context, userID string, params pagination.Pagination) (*pagination.PaginationResult[Barbershop], error) {
	var barbershops []Barbershop
	var total int64

	query := r.db.WithContext(ctx).
		Model(&Barbershop{}).
		Joins("JOIN user_barbershops ub ON ub.barbershop_id = barbershops.id").
		Where("ub.user_id = ? AND deleted_at IS NULL", userID)

	// Apply Search
	if params.Search != "" {
		pattern := "%" + params.Search + "%"
		query = query.Where("barbershops.name ILIKE ? OR barbershops.email ILIKE ?", pattern, pattern)
	}

	// Count Total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	orderClause := fmt.Sprintf("barbershops.%s %s", params.SortBy, params.SortOrder)
	query = query.Order(orderClause)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Find(&barbershops).Error; err != nil {
		return nil, err
	}

	return &pagination.PaginationResult[Barbershop]{
		Data:       barbershops,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: (total + int64(params.Limit) - 1) / int64(params.Limit),
	}, nil

}

func (r *repository) FindAll(ctx context.Context, params pagination.Pagination) (*pagination.PaginationResult[Barbershop], error) {
	var barbershops []Barbershop
	var total int64

	query := r.db.WithContext(ctx).
		Model(&Barbershop{}).
		Where("deleted_at IS NULL")

	// Apply search
	if params.Search != "" {
		pattern := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR address ILIKE ?", pattern, pattern, pattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	orderClause := fmt.Sprintf("%s %s", params.SortBy, params.SortOrder)
	query = query.Order(orderClause)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Find(&barbershops).Error; err != nil {
		return nil, err
	}

	return &pagination.PaginationResult[Barbershop]{
		Data:       barbershops,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: (total + int64(params.Limit) - 1) / int64(params.Limit),
	}, nil
}

func (r *repository) Update(ctx context.Context, barbershop *Barbershop) error {
	return r.db.WithContext(ctx).Save(barbershop).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&Barbershop{}).Error
}

func (r *repository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Barbershop{}).
		Where("slug = ? AND deleted_at IS NULL", slug).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) CheckEmailExists(ctx context.Context, email string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&Barbershop{}).
		Where("email = ? AND deleted_at IS NULL", email)

	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *repository) invalidateCache(ctx context.Context) {
	r.redis.Del(ctx, "products:all")
}
