package routes

import (
	"bLink-app/internal/features/user/handler"
	"bLink-app/internal/middleware"
	"bLink-app/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all user routes
func RegisterRoutes(router *gin.RouterGroup, handler *handler.Handler, jwtService *jwt.JWTService) {
	users := router.Group("/users")
	{
		// Public routes - Anyone can access
		users.POST("/register", handler.Register)
		users.POST("/login", handler.Login)

		// Protected routes - Require authentication
		protected := users.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			// Profile management
			protected.GET("/profile", handler.GetProfile)
			protected.PUT("/profile", handler.UpdateProfile)

			// Password management
			protected.POST("/change-password", handler.ChangePassword)

			// Account management
			protected.DELETE("/account", handler.DeleteAccount)

			// Admin routes - Get all users with pagination
			protected.GET("", handler.GetAllUsers)
		}
	}
}
