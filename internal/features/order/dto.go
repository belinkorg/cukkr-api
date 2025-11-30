package order

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required,gt=0"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending processing completed cancelled"`
}

type OrderResponse struct {
	ID         string              `json:"id"`
	UserID     string              `json:"user_id"`
	TotalPrice float64             `json:"total_price"`
	Status     string              `json:"status"`
	Items      []OrderItemResponse `json:"items"`
	CreatedAt  string              `json:"created_at"`
	UpdatedAt  string              `json:"updated_at"`
}

type OrderItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
