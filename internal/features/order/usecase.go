package order

import (
	"context"
	"cukkr-app/internal/features/product"
	"cukkr-app/pkg/errors"
	"cukkr-app/pkg/logger"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

type Usecase interface {
	CreateOrder(ctx context.Context, userID string, req *CreateOrderRequest) (*OrderResponse, error)
	GetOrder(ctx context.Context, id string) (*OrderResponse, error)
	GetUserOrders(ctx context.Context, userID string) (*[]OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, id string, req *UpdateOrderStatusRequest) (*OrderResponse, error)
}

type usecase struct {
	db          *gorm.DB
	repo        Repository
	productRepo product.Repository
	logger      *logger.Logger
}

func NewUsecase(db *gorm.DB, repo Repository, productRepo product.Repository, logger *logger.Logger) Usecase {
	return &usecase{
		db:          db,
		repo:        repo,
		productRepo: productRepo,
		logger:      logger,
	}
}

func (u *usecase) CreateOrder(ctx context.Context, userID string, req *CreateOrderRequest) (*OrderResponse, error) {
	u.logger.WithFields(map[string]interface{}{
		"user_id":     userID,
		"items_count": len(req.Items),
	}).Info("Creating order")

	var createdOrder *Order

	// PROPER TRANSACTION with consistent context
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var totalPrice float64
		orderItems := make([]OrderItem, 0, len(req.Items))

		// Validate products
		for _, item := range req.Items {
			prod, err := u.productRepo.FindByID(ctx, item.ProductID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					u.logger.WithField("product_id", item.ProductID).Warn("Product not found")
					return errors.New(http.StatusNotFound, "product not found")
				}
				u.logger.WithError(err).Error("Failed to get product")
				return err
			}

			if prod.Stock < item.Quantity {
				u.logger.WithFields(map[string]interface{}{
					"product_id": item.ProductID,
					"available":  prod.Stock,
					"requested":  item.Quantity,
				}).Warn("Insufficient stock")
				return errors.New(
					http.StatusBadRequest, fmt.Sprintf("insufficient stock for product: %s (available: %d, requested: %d)",
						prod.Name, prod.Stock, item.Quantity),
				)
			}

			itemPrice := prod.Price * float64(item.Quantity)
			totalPrice += itemPrice

			orderItems = append(orderItems, OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     prod.Price,
			})
		}

		// Create order
		order := &Order{
			UserID:     userID,
			TotalPrice: totalPrice,
			Status:     "PENDING",
			Items:      orderItems,
		}

		if err := u.repo.Create(tx, order); err != nil {
			u.logger.WithError(err).Error("Failed to create order")
			return err
		}

		// Update stocks
		for _, item := range req.Items {
			affected, err := u.productRepo.UpdateStock(tx, item.ProductID, item.Quantity)
			if err != nil {
				u.logger.WithError(err).WithField("product_id", item.ProductID).Error("Failed to update stock")
				return errors.New(http.StatusBadRequest, "Failed to update product stock")
			}

			if affected == 0 {
				return errors.New(http.StatusBadRequest, "Insufficient stock or product not found")
			}
		}

		createdOrder = order
		return nil
	})

	if err != nil {
		return nil, err
	}

	u.logger.WithFields(map[string]interface{}{
		"order_id":    createdOrder.ID,
		"total_price": createdOrder.TotalPrice,
	}).Info("Order created")

	return toOrderResponse(createdOrder), nil
}

func (u *usecase) GetOrder(ctx context.Context, id string) (*OrderResponse, error) {
	// Note: GetOrder bisa dipanggil tanpa items, users load sesuai kebutuhan
	order, err := u.repo.FindByID(ctx, id)
	if err != nil {
		u.logger.WithError(err).Error("Failed to get order")
		return nil, errors.New(http.StatusNotFound, "Order not found")
	}

	return toOrderResponse(order), nil
}

func (u *usecase) GetUserOrders(ctx context.Context, userID string) (*[]OrderResponse, error) {
	orders, err := u.repo.FindByUserID(ctx, userID)
	if err != nil {
		u.logger.WithFields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		}).Error("Failed to get user orders")
		return nil, errors.New(http.StatusNotFound, "Orders for this user is not found")
	}

	responses := toOrderResponses(orders)
	return &responses, nil
}

func (u *usecase) UpdateOrderStatus(ctx context.Context, id string, req *UpdateOrderStatusRequest) (*OrderResponse, error) {
	var updatedOrder *Order

	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update order status
		result := tx.Model(&Order{}).Where("id = ?", id).Update("status", req.Status)

		if result.Error != nil {
			u.logger.WithError(result.Error).Error("Failed to update order status")
			return result.Error
		}

		// Check RowsAffected untuk validasi
		if result.RowsAffected == 0 {
			return errors.New(http.StatusNotFound, "Order not found")
		}

		// Get updated order for response
		if err := tx.First(&updatedOrder, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New(http.StatusNotFound, "Order not found")
			}
			u.logger.WithError(err).Error("Failed to fetch updated order")
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	u.logger.WithFields(map[string]interface{}{
		"order_id": updatedOrder.ID,
		"status":   req.Status,
	}).Info("Order status updated successfully")

	return toOrderResponse(updatedOrder), nil
}

func toOrderResponse(o *Order) *OrderResponse {
	items := make([]OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return &OrderResponse{
		ID:         o.ID,
		UserID:     o.UserID,
		TotalPrice: o.TotalPrice,
		Status:     o.Status,
		Items:      items,
		CreatedAt:  o.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  o.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toOrderResponses(orders []Order) []OrderResponse {
	responses := make([]OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = *toOrderResponse(&o)
	}
	return responses
}
