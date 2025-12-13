package product

import (
	"cukurly-app/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all product routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtService *jwt.JWTService) {
	products := router.Group("/products")
	{
		// Public routes
		products.GET("", handler.GetAllProducts)
		products.GET("/:id", handler.GetProduct)

		// Protected routes (require authentication)
		protected := products.Group("")
		//protected.Use(middleware.AuthMiddleware(jwtService))
		{
			protected.POST("", handler.CreateProduct)
			protected.PUT("/:id", handler.UpdateProduct)
			protected.DELETE("/:id", handler.DeleteProduct)
		}
	}
}
