package router

import (
	"cukurly-app/internal/features/order"
	"cukurly-app/internal/features/product"
	"cukurly-app/internal/features/user"
	"cukurly-app/internal/middleware"
	"cukurly-app/pkg/jwt"
	"cukurly-app/pkg/logger"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "cukurly-app/docs" // Import swagger docs
)

type Router struct {
	engine         *gin.Engine
	productHandler *product.Handler
	orderHandler   *order.Handler
	userHandler    *user.Handler
	jwtService     *jwt.JWTService
	logger         *logger.Logger
}

func NewRouter(
	productHandler *product.Handler,
	orderHandler *order.Handler,
	userHandler *user.Handler,
	jwtService *jwt.JWTService,
	logger *logger.Logger,
) *Router {
	return &Router{
		engine:         gin.Default(),
		productHandler: productHandler,
		orderHandler:   orderHandler,
		userHandler:    userHandler,
		jwtService:     jwtService,
		logger:         logger,
	}
}

func (r *Router) Setup() *gin.Engine {
	// Global middlewares
	r.engine.Use(middleware.CORSMiddleware())
	r.engine.Use(middleware.LoggerMiddleware(r.logger))

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger documentation
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	v1 := r.engine.Group("/api/v1")

	// Register feature routes
	product.RegisterRoutes(v1, r.productHandler, r.jwtService)
	order.RegisterRoutes(v1, r.orderHandler, r.jwtService)
	user.RegisterRoutes(v1, r.userHandler, r.jwtService)

	return r.engine
}
