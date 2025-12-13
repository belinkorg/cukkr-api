package product_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"bLink-app/internal/features/product"
	"bLink-app/pkg/errors"
	"bLink-app/tests/mocks"
	"bLink-app/tests/testutil"
)

// --- UUID Helper lokal untuk test ---
func newUUID() uuid.UUID {
	return uuid.New()
}

func TestCreateProduct_Success(t *testing.T) {
	db, _ := testutil.SetupTestDBWithTransaction(t)
	mockRepo := new(mocks.ProductRepositoryMock)
	logger := testutil.GetTestLogger()
	uc := product.NewUsecase(db, mockRepo, logger)

	req := &product.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100000,
		Stock:       10,
		Category:    "Electronics",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).
		Return(nil).
		Run(func(args mock.Arguments) {
			p := args.Get(1).(*product.Product)
			p.ID = newUUID().String()
		})

	result, err := uc.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Product", result.Name)
	assert.Equal(t, float64(100000), result.Price)
	mockRepo.AssertExpectations(t)
}

func TestGetProduct_Success(t *testing.T) {
	db, _ := testutil.SetupTestDBWithTransaction(t)
	mockRepo := new(mocks.ProductRepositoryMock)
	logger := testutil.GetTestLogger()
	uc := product.NewUsecase(db, mockRepo, logger)

	expectedProduct := &product.Product{
		ID:          newUUID().String(),
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100000,
		Stock:       10,
		Category:    "Electronics",
	}

	mockRepo.On("FindByID", mock.Anything, expectedProduct.ID).Return(expectedProduct, nil)

	result, err := uc.GetProduct(context.Background(), expectedProduct.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProduct.ID, result.ID)
	assert.Equal(t, "Test Product", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetProduct_NotFound(t *testing.T) {
	db, _ := testutil.SetupTestDB(t)
	mockRepo := new(mocks.ProductRepositoryMock)
	logger := testutil.GetTestLogger()
	uc := product.NewUsecase(db, mockRepo, logger)

	nonExistingID := newUUID().String()

	mockRepo.On("FindByID", mock.Anything, nonExistingID).Return(nil, errors.ErrNotFound)

	result, err := uc.GetProduct(context.Background(), nonExistingID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errors.ErrNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts_Success(t *testing.T) {
	db, _ := testutil.SetupTestDB(t)
	mockRepo := new(mocks.ProductRepositoryMock)
	logger := testutil.GetTestLogger()
	uc := product.NewUsecase(db, mockRepo, logger)

	expectedProducts := []product.Product{
		{ID: newUUID().String(), Name: "Product 1", Price: 100000, Stock: 10, Category: "Electronics"},
		{ID: newUUID().String(), Name: "Product 2", Price: 200000, Stock: 5, Category: "Fashion"},
	}

	mockRepo.On("FindAll", mock.Anything).Return(expectedProducts, nil)

	result, err := uc.GetAllProducts(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Success(t *testing.T) {
	db, _ := testutil.SetupTestDBWithTransaction(t)
	mockRepo := new(mocks.ProductRepositoryMock)
	logger := testutil.GetTestLogger()
	uc := product.NewUsecase(db, mockRepo, logger)

	existingProduct := &product.Product{
		ID:          newUUID().String(),
		Name:        "Old Name",
		Description: "Old Description",
		Price:       100000,
		Stock:       10,
		Category:    "Electronics",
	}

	name := "New Name"
	price := 150000.0
	req := &product.UpdateProductRequest{
		Name:  &name,
		Price: &price,
	}

	mockRepo.On("FindByID", mock.Anything, existingProduct.ID).Return(existingProduct, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil)

	result, err := uc.UpdateProduct(context.Background(), existingProduct.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, float64(150000), result.Price)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	db, _ := testutil.SetupTestDBWithTransaction(t)
	mockRepo := new(mocks.ProductRepositoryMock)
	logger := testutil.GetTestLogger()
	uc := product.NewUsecase(db, mockRepo, logger)

	existingProduct := &product.Product{ID: newUUID().String(), Name: "Test Product"}

	mockRepo.On("FindByID", mock.Anything, existingProduct.ID).Return(existingProduct, nil)
	mockRepo.On("Delete", mock.Anything, existingProduct.ID).Return(nil)

	err := uc.DeleteProduct(context.Background(), existingProduct.ID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_NotFound(t *testing.T) {
	db, _ := testutil.SetupTestDB(t)
	mockRepo := new(mocks.ProductRepositoryMock)
	logger := testutil.GetTestLogger()
	uc := product.NewUsecase(db, mockRepo, logger)

	nonExistingID := newUUID().String()

	mockRepo.On("FindByID", mock.Anything, nonExistingID).Return(nil, errors.ErrNotFound)

	err := uc.DeleteProduct(context.Background(), nonExistingID)

	assert.Error(t, err)
	assert.Equal(t, errors.ErrNotFound, err)
	mockRepo.AssertExpectations(t)
}
