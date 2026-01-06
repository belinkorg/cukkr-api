package product

import (
	"context"
	"cukkr-app/pkg/errors"
	"cukkr-app/pkg/logger"
	"gorm.io/gorm"
	"net/http"
)

type Usecase interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error)
	GetProduct(ctx context.Context, id string) (*ProductResponse, error)
	GetAllProducts(ctx context.Context) ([]ProductResponse, error)
	UpdateProduct(ctx context.Context, id string, req *UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, id string) error
}

type usecase struct {
	db     *gorm.DB
	repo   Repository
	logger *logger.Logger
}

func NewUsecase(db *gorm.DB, repo Repository, logger *logger.Logger) Usecase {
	return &usecase{
		db:     db,
		repo:   repo,
		logger: logger,
	}
}

func (u *usecase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error) {
	u.logger.WithField("name", req.Name).Info("Creating product")

	product := &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}

	// Transaction
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return u.repo.Create(tx, product)
	})

	if err != nil {
		return nil, err
	}

	u.logger.WithField("product_id", product.ID).Info("Product created")
	return toProductResponse(product), nil
}

func (u *usecase) GetProduct(ctx context.Context, id string) (*ProductResponse, error) {
	product, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			u.logger.WithField("product_id", id).Warn("Product not found")
			return nil, errors.New(http.StatusNotFound, "Product not found")
		}
		u.logger.WithError(err).Error("Failed to get product")
		return nil, err
	}

	return toProductResponse(product), nil
}

func (u *usecase) GetAllProducts(ctx context.Context) ([]ProductResponse, error) {
	products, err := u.repo.FindAll(ctx)
	if err != nil {
		u.logger.WithError(err).Error("Failed to get products")
		return nil, err
	}

	u.logger.WithField("count", len(products)).Info("Products retrieved")
	responses := toProductResponses(products)
	return responses, nil
}

func (u *usecase) UpdateProduct(ctx context.Context, id string, req *UpdateProductRequest) (*ProductResponse, error) {
	u.logger.WithField("product_id", id).Info("Updating product")

	var updatedProduct *Product

	// Transaction
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		product, err := u.repo.FindByID(ctx, id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				u.logger.WithField("product_id", id).Warn("Product not found")
				return errors.New(http.StatusNotFound, "Product not found")
			}
			u.logger.WithError(err).Error("Failed to get product")
			return err
		}

		// Update fields dengan pointer untuk distinguish nil vs 0
		if req.Name != nil && *req.Name != "" {
			product.Name = *req.Name
		}
		if req.Description != nil && *req.Description != "" {
			product.Description = *req.Description
		}
		if req.Price != nil && *req.Price > 0 {
			product.Price = *req.Price
		}
		if req.Stock != nil && *req.Stock >= 0 {
			product.Stock = *req.Stock
		}
		if req.Category != nil && *req.Category != "" {
			product.Category = *req.Category
		}

		if err := u.repo.Update(tx, product); err != nil {
			u.logger.WithError(err).Error("Failed to update product")
			return err
		}

		updatedProduct = product
		return nil
	})

	if err != nil {
		return nil, err
	}

	u.logger.WithField("product_id", id).Info("Product updated")
	return toProductResponse(updatedProduct), nil
}

func (u *usecase) DeleteProduct(ctx context.Context, id string) error {
	u.logger.WithField("product_id", id).Info("Deleting product")

	_, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return u.repo.Delete(tx, id)
	})

	if err != nil {
		u.logger.WithError(err).Error("Failed to delete product")
		return err
	}

	u.logger.WithField("product_id", id).Info("Product deleted")
	return nil
}

func toProductResponse(p *Product) *ProductResponse {
	return &ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toProductResponses(products []Product) []ProductResponse {
	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = *toProductResponse(&p)
	}
	return responses
}
