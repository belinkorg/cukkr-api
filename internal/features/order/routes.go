package order

import (
	"bLink-app/internal/middleware"
	"bLink-app/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all order routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtService *jwt.JWTService) {
	orders := router.Group("/orders")
	orders.Use(middleware.AuthMiddleware(jwtService)) // All order routes require auth
	{
		orders.POST("", handler.CreateOrder)
		orders.GET("", handler.GetUserOrders)
		orders.GET("/:id", handler.GetOrder)
		orders.PATCH("/:id/status", handler.UpdateOrderStatus)
	}
}
