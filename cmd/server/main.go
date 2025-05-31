package main

import (
	"log"
	"os"
	"sales-analysis-system/internal/config"
	"sales-analysis-system/internal/database"
	"sales-analysis-system/internal/handlers"
	"sales-analysis-system/internal/middleware"
	"sales-analysis-system/internal/services"
	"sales-analysis-system/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize logger
	logger := utils.NewLogger()

	// Initialize configuration
	cfg := config.New()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		logger.Fatal("Failed to run migrations: ", err)
	}

	// Initialize services
	csvLoader := services.NewCSVLoader(db, logger)
	analyticsService := services.NewAnalyticsService(db, logger)
	refreshService := services.NewRefreshService(db, csvLoader, logger)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, logger)
	refreshHandler := handlers.NewRefreshHandler(refreshService, logger)

	// Setup cron for daily refresh
	c := cron.New()
	c.AddFunc("0 2 * * *", func() { // Daily at 2 AM
		logger.Info("Starting scheduled data refresh")
		refreshService.RefreshData("data/sales_data.csv")
	})
	c.Start()
	defer c.Stop()

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))

	// Health check
	router.GET("/health", healthHandler.Health)

	// API routes
	api := router.Group("/api/v1")
	{
		// Data refresh
		api.POST("/refresh", refreshHandler.TriggerRefresh)
		api.GET("/refresh/status", refreshHandler.GetRefreshStatus)

		// Analytics endpoints
		analytics := api.Group("/analytics")
		{
			analytics.GET("/revenue/total", analyticsHandler.GetTotalRevenue)
			analytics.GET("/revenue/by-product", analyticsHandler.GetRevenueByProduct)
			analytics.GET("/revenue/by-category", analyticsHandler.GetRevenueByCategory)
			analytics.GET("/revenue/by-region", analyticsHandler.GetRevenueByRegion)
			analytics.GET("/revenue/trends", analyticsHandler.GetRevenueTrends)

			analytics.GET("/products/top", analyticsHandler.GetTopProducts)
			analytics.GET("/products/top/by-category", analyticsHandler.GetTopProductsByCategory)
			analytics.GET("/products/top/by-region", analyticsHandler.GetTopProductsByRegion)

			analytics.GET("/customers/count", analyticsHandler.GetCustomerCount)
			analytics.GET("/orders/count", analyticsHandler.GetOrderCount)
			analytics.GET("/orders/average-value", analyticsHandler.GetAverageOrderValue)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting on port ", port)
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server: ", err)
	}
}
