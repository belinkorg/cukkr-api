package product

import (
	"cukkr-app/pkg/helper"
	"cukkr-app/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase    Usecase
	httpHelper *helper.HTTPHandlerHelper
}

func NewHandlers(usecase Usecase, httpHelper *helper.HTTPHandlerHelper) *Handler {
	return &Handler{
		usecase:    usecase,
		httpHelper: httpHelper,
	}
}

// CreateProduct godoc
//
//	@Summary		Create a new product
//	@Description	Create a new product (requires authentication)
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product	body		CreateProductRequest	true	"Product data"
//	@Success		201		{object}	response.Response{data=ProductResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Security		BearerAuth
//	@Router			/products [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := h.httpHelper.BindAndValidate(c, &req); err != nil {
		return
	}

	result, err := h.usecase.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to create product")
		return
	}

	response.Success(c, http.StatusCreated, "Product created successfully", result)
}

// GetProduct godoc
//
//	@Summary		Get product by ID
//	@Description	Get a single product by ID with caching
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Product ID"
//	@Success		200	{object}	response.Response{data=ProductResponse}
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/products/{id} [get]
func (h *Handler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	result, err := h.usecase.GetProduct(c.Request.Context(), id)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to get product")
		return
	}

	response.Success(c, http.StatusOK, "Product retrieved successfully", result)
}

// GetAllProducts godoc
//
//	@Summary		Get all products
//	@Description	Get list of all products with caching
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=[]ProductResponse}
//	@Failure		500	{object}	response.Response
//	@Router			/products [get]
func (h *Handler) GetAllProducts(c *gin.Context) {
	result, err := h.usecase.GetAllProducts(c.Request.Context())
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to get products")
		return
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", result)
}

// UpdateProduct godoc
//
//	@Summary		Update a product
//	@Description	Update an existing product (requires authentication)
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Product ID"
//	@Param			product	body		UpdateProductRequest	true	"Product data"
//	@Success		200		{object}	response.Response{data=ProductResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Security		BearerAuth
//	@Router			/products/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.httpHelper.ValidateIDWithField(c, id, "product_id"); err != nil {
		return
	}

	var req UpdateProductRequest
	if err := h.httpHelper.BindAndValidate(c, &req); err != nil {
		return
	}

	result, err := h.usecase.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to update product")
		return
	}

	response.Success(c, http.StatusOK, "Product updated successfully", result)
}

// DeleteProduct godoc
//
//	@Summary		Delete a product
//	@Description	Delete a product by ID (requires authentication)
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Product ID"
//	@Success		200	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Security		BearerAuth
//	@Router			/products/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.httpHelper.ValidateIDWithField(c, id, "product_id"); err != nil {
		return
	}

	err := h.usecase.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		h.httpHelper.HandleError(c, err, "Failed to delete product")
		return
	}

	response.Success(c, http.StatusOK, "Product deleted successfully", nil)
}
