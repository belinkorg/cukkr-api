package mocks

import (
	"bLink-app/internal/features/product"
	"context"
	"gorm.io/gorm"

	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) Create(tx *gorm.DB, p *product.Product) error {
	args := m.Called(tx, p)
	return args.Error(0)
}

func (m *ProductRepositoryMock) FindByID(ctx context.Context, id string) (*product.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *ProductRepositoryMock) FindAll(ctx context.Context) ([]product.Product, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]product.Product), args.Error(1)
}

func (m *ProductRepositoryMock) Update(tx *gorm.DB, p *product.Product) error {
	args := m.Called(tx, p)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(tx *gorm.DB, id string) error {
	args := m.Called(tx, id)
	return args.Error(0)
}

func (m *ProductRepositoryMock) UpdateStock(tx *gorm.DB, id string, quantity int) (int64, error) {
	args := m.Called(tx, id, quantity)
	return 0, args.Error(0)
}
