package order

import (
	"cukurly-app/pkg/helper"
	"cukurly-app/pkg/jwt"
	"cukurly-app/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	usecase    Usecase
	httpHelper *helper.HTTPHandlerHelper
}

func NewHandler(usecase Usecase, httpHelper *helper.HTTPHandlerHelper) *Handler {
	return &Handler{
		usecase:    usecase,
		httpHelper: httpHelper,
	}
}

// CreateOrder godoc
//
//	@Summary		Create a new order
//	@Description	Create a new order with products (requires authentication)
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Param			order	body		CreateOrderRequest	true	"Order data"
//	@Success		201		{object}	response.Response{data=OrderResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Security		BearerAuth
//	@Router			/orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := h.httpHelper.BindAndValidate(c, &req); err != nil {
		return
	}

	// Get user ID from JWT claims
	claims := c.MustGet("claims").(*jwt.Claims)

	result, err := h.usecase.CreateOrder(c.Request.Context(), claims.UserID, &req)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to create order")
		return
	}

	response.Success(c, http.StatusCreated, "Order created successfully", result)
}

// GetOrder godoc
//
//	@Summary		Get order by ID
//	@Description	Get a single order by ID (requires authentication)
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Order ID"
//	@Success		200	{object}	response.Response{data=OrderResponse}
//	@Failure		401	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Security		BearerAuth
//	@Router			/orders/{id} [get]
func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	result, err := h.usecase.GetOrder(c.Request.Context(), id)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to get order")
		return
	}

	response.Success(c, http.StatusOK, "Order retrieved successfully", result)
}

// GetUserOrders godoc
//
//	@Summary		Get user orders
//	@Description	Get all orders for the authenticated user
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=[]OrderResponse}
//	@Failure		401	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Security		BearerAuth
//	@Router			/orders [get]
func (h *Handler) GetUserOrders(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.Claims)

	result, err := h.usecase.GetUserOrders(c.Request.Context(), claims.UserID)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to get user orders")
		return
	}
	response.Success(c, http.StatusOK, "Orders retrieved successfully", result)
}

// UpdateOrderStatus godoc
//
//	@Summary		Update order status
//	@Description	Update the status of an order (requires authentication)
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Order ID"
//	@Param			status	body		UpdateOrderStatusRequest	true	"Status data"
//	@Success		200		{object}	response.Response{data=OrderResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Security		BearerAuth
//	@Router			/orders/{id}/status [patch]
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateOrderStatusRequest
	if err := h.httpHelper.BindAndValidate(c, &req); err != nil {
		return
	}

	result, err := h.usecase.UpdateOrderStatus(c.Request.Context(), id, &req)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to update order status")
		return
	}

	response.Success(c, http.StatusOK, "Order status updated successfully", result)
}
