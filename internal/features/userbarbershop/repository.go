package userbarbershop

import (
	"context"
	"cukkr-app/pkg/pagination"
	"fmt"
	"gorm.io/gorm"
)

type UserBarbershopRepository interface {
	Create(ctx context.Context, tx *gorm.DB, userBarbershop *UserBarbershop) error
	FindByUserID(ctx context.Context, userID string, params pagination.Pagination) (*pagination.PaginationResult[UserBarbershop], error)
	FindByBarbershopID(ctx context.Context, barbershopID string, params pagination.Pagination) (*pagination.PaginationResult[UserBarbershop], error)
	FindByUserAndBarbershop(ctx context.Context, userID, barbershopID string) (*UserBarbershop, error)
	IsOwner(ctx context.Context, userID, barbershopID string) (bool, error)
	IsMember(ctx context.Context, userID, barbershopID string) (bool, error)
	UpdateRole(ctx context.Context, tx *gorm.DB, userID, barbershopID string, role UserBarbershopRole) error
	Delete(ctx context.Context, userID, barbershopID string) error
	CheckMemberExists(ctx context.Context, userID, barbershopID string) (bool, error)
	CountOwners(ctx context.Context, barbershopID string) (int64, error)
}

type userBarbershopRepository struct {
	db *gorm.DB
}

func NewUserBarbershopRepository(db *gorm.DB) UserBarbershopRepository {
	return &userBarbershopRepository{db: db}
}

func (r *userBarbershopRepository) Create(ctx context.Context, tx *gorm.DB, userBarbershop *UserBarbershop) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(userBarbershop).Error
}

func (r *userBarbershopRepository) FindByUserID(ctx context.Context, userID string, params pagination.Pagination) (*pagination.PaginationResult[UserBarbershop], error) {
	var userBarbershops []UserBarbershop
	var total int64

	query := r.db.WithContext(ctx).
		Model(&UserBarbershop{}).
		Where("user_id = ?", userID)

	// Apply search (search by role)
	if params.Search != "" {
		query = query.Where("role ILIKE ?", "%"+params.Search+"%")
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
	if err := query.Offset(offset).Limit(params.Limit).Find(&userBarbershops).Error; err != nil {
		return nil, err
	}

	return &pagination.PaginationResult[UserBarbershop]{
		Data:       userBarbershops,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: (total + int64(params.Limit) - 1) / int64(params.Limit),
	}, nil
}

func (r *userBarbershopRepository) FindByBarbershopID(ctx context.Context, barbershopID string, params pagination.Pagination) (*pagination.PaginationResult[UserBarbershop], error) {
	var userBarbershops []UserBarbershop
	var total int64

	query := r.db.WithContext(ctx).
		Model(&UserBarbershop{}).
		Where("barbershop_id = ?", barbershopID)

	// Apply search (search by role)
	if params.Search != "" {
		query = query.Where("role ILIKE ?", "%"+params.Search+"%")
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
	if err := query.Offset(offset).Limit(params.Limit).Find(&userBarbershops).Error; err != nil {
		return nil, err
	}

	return &pagination.PaginationResult[UserBarbershop]{
		Data:       userBarbershops,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: (total + int64(params.Limit) - 1) / int64(params.Limit),
	}, nil
}

func (r *userBarbershopRepository) FindByUserAndBarbershop(ctx context.Context, userID, barbershopID string) (*UserBarbershop, error) {
	var userBarbershop UserBarbershop
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND barbershop_id = ?", userID, barbershopID).
		First(&userBarbershop).Error
	if err != nil {
		return nil, err
	}
	return &userBarbershop, nil
}

func (r *userBarbershopRepository) IsOwner(ctx context.Context, userID, barbershopID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&UserBarbershop{}).
		Where("user_id = ? AND barbershop_id = ? AND role = ?", userID, barbershopID, RoleOwner).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userBarbershopRepository) IsMember(ctx context.Context, userID, barbershopID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&UserBarbershop{}).
		Where("user_id = ? AND barbershop_id = ?", userID, barbershopID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userBarbershopRepository) UpdateRole(ctx context.Context, tx *gorm.DB, userID, barbershopID string, role UserBarbershopRole) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).
		Model(&UserBarbershop{}).
		Where("user_id = ? AND barbershop_id = ?", userID, barbershopID).
		Update("role", role).Error
}

func (r *userBarbershopRepository) Delete(ctx context.Context, userID, barbershopID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND barbershop_id = ?", userID, barbershopID).
		Delete(&UserBarbershop{}).Error
}

func (r *userBarbershopRepository) CheckMemberExists(ctx context.Context, userID, barbershopID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&UserBarbershop{}).
		Where("user_id = ? AND barbershop_id = ?", userID, barbershopID).
		Count(&count).Error
	return count > 0, err
}

func (r *userBarbershopRepository) CountOwners(ctx context.Context, barbershopID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&UserBarbershop{}).
		Where("barbershop_id = ? AND role = ?", barbershopID, RoleOwner).
		Count(&count).Error
	return count, err
}
