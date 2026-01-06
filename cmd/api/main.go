package main

import (
	"cukkr-app/config"
	"cukkr-app/internal/database"
	"cukkr-app/internal/features/order"
	"cukkr-app/internal/features/product"
	"cukkr-app/internal/features/user"
	"cukkr-app/internal/router"
	"cukkr-app/pkg/helper"
	"cukkr-app/pkg/jwt"
	"cukkr-app/pkg/logger"
	"cukkr-app/pkg/validator"
	"fmt"
	"log"
)

//	@title			cukkr App API
//	@version		1.0
//	@description	Clean Architecture REST API with Go, Gin, PostgreSQL, Redis, and MongoDB
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.example.com/support
//	@contact.email	support@example.com

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@host		localhost:8080
//	@BasePath	/api/v1

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg.Log.Level)
	appLogger.Info("Starting application...")

	// Initialize databases
	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		appLogger.Fatal("Failed to connect to PostgreSQL:", err)
	}
	appLogger.Info("Connected to PostgreSQL")

	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		appLogger.Fatal("Failed to connect to Redis:", err)
	}
	appLogger.Info("Connected to Redis")

	//mongoClient, err := database.NewMongoClient(cfg.MongoDB)
	//if err != nil {
	//	appLogger.Fatal("Failed to connect to MongoDB:", err)
	//}
	//appLogger.Info("Connected to MongoDB")
	//_ = mongoClient // MongoDB for future use (logging, analytics, etc.)

	// Auto migrate
	//if err := db.AutoMigrate(&product.Product{}, &order.Order{}, &order.OrderItem{}, &user.User{}); err != nil {
	//	appLogger.Fatal("Failed to migrate database:", err)
	//}
	//appLogger.Info("Database migration completed")

	// Initialize validator
	v := validator.NewValidator()

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiredTime)

	// Initialize HTTPHandlerHelper
	httpHelper := helper.NewHTTPHandlerHelper(v, appLogger)

	// Initialize Product feature
	productRepo := product.NewRepository(db, redisClient)
	productUsecase := product.NewUsecase(db, productRepo, appLogger)
	productHandler := product.NewHandlers(productUsecase, httpHelper)

	// Initialize Order feature
	orderRepo := order.NewRepository(db)
	orderUsecase := order.NewUsecase(db, orderRepo, productRepo, appLogger)
	orderHandler := order.NewHandler(orderUsecase, httpHelper)

	// Initialize User feature
	userRepo := user.NewRepository(db, redisClient)
	emailUsecase := user.NewEmailUsecase(appLogger)
	authUsecase := user.NewAuthUsecase(jwtService, appLogger)
	otpUsecase := user.NewOTPUsecase(userRepo, appLogger, emailUsecase)
	userUsecase := user.NewUserUsecase(userRepo, authUsecase, otpUsecase, emailUsecase, appLogger)
	userHandler := user.NewHandler(userUsecase, httpHelper)

	// Setup router
	r := router.NewRouter(productHandler, orderHandler, userHandler, jwtService, appLogger)
	engine := r.Setup()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	appLogger.WithFields(map[string]interface{}{
		"port": cfg.App.Port,
		"env":  cfg.App.Env,
	}).Info("Server is running")

	if err := engine.Run(addr); err != nil {
		appLogger.Fatal("Failed to start server:", err)
	}
}
