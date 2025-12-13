package mocks

import (
	"context"
	"cukurly-app/internal/features/product"

	"github.com/stretchr/testify/mock"
)

type ProductUsecaseMock struct {
	mock.Mock
}

func (m *ProductUsecaseMock) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.ProductResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.ProductResponse), args.Error(1)
}

func (m *ProductUsecaseMock) GetProduct(ctx context.Context, id string) (*product.ProductResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.ProductResponse), args.Error(1)
}

func (m *ProductUsecaseMock) GetAllProducts(ctx context.Context) ([]product.ProductResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]product.ProductResponse), args.Error(1)
}

func (m *ProductUsecaseMock) UpdateProduct(ctx context.Context, id string, req *product.UpdateProductRequest) (*product.ProductResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.ProductResponse), args.Error(1)
}

func (m *ProductUsecaseMock) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
