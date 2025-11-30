package product_test

import (
	"bLink-app/internal/features/product"
	_ "bLink-app/pkg/errors"
	"bLink-app/pkg/validator"
	"bLink-app/tests/mocks"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

func TestHandler_CreateProduct_Success(t *testing.T) {
	mockUsecase := new(mocks.ProductUsecaseMock)
	v := validator.NewValidator()
	handler := product.NewHandlers(mockUsecase, v)
	router := setupTestRouter()
	router.POST("/products", handler.CreateProduct)

	req := product.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100000,
		Stock:       10,
		Category:    "Electronics",
	}

	expectedResponse := &product.ProductResponse{
		ID:          "1b47d260-3c9d-4b0a-90d3-b1b8a4a8c201", // contoh UUID
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}

	mockUsecase.
		On("CreateProduct", mock.Anything, mock.AnythingOfType("*product.CreateProductRequest")).
		Return(expectedResponse, nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUsecase.AssertExpectations(t)
}

//func TestHandler_CreateProduct_ValidationError(t *testing.T) {
//	mockUsecase := new(mocks.ProductUsecaseMock)
//	v := validator.NewValidator()
//	handler := product.NewHandlers(mockUsecase, v)
//	router := setupTestRouter()
//	router.POST("/products", handler.CreateProduct)
//
//	req := product.CreateProductRequest{}
//}
