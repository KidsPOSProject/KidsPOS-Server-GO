package main

import (
	"log"
	"os"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/config"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/handlers"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/repository"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database
	db, err := repository.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := repository.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	services := service.NewServices(repos)

	// Initialize Gin router
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("web/templates/**/*")

	// Serve static files
	router.Static("/static", "./web/static")
	router.Static("/css", "./web/static/css")
	router.Static("/js", "./web/static/js")

	// Initialize handlers
	h := handlers.NewHandlers(services)

	// Setup routes
	setupRoutes(router, h)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(router *gin.Engine, h *handlers.Handlers) {
	// Web routes
	router.GET("/", h.Home)
	router.GET("/items", h.ItemsList)
	router.GET("/items/new", h.ItemsNew)
	router.POST("/items", h.ItemsCreate)
	router.GET("/items/:id/edit", h.ItemsEdit)
	router.POST("/items/:id", h.ItemsUpdate)
	router.POST("/items/:id/delete", h.ItemsDelete)

	router.GET("/sales", h.SalesList)
	router.GET("/sales/new", h.SalesNew)
	router.POST("/sales", h.SalesCreate)

	router.GET("/stores", h.StoresList)
	router.GET("/stores/new", h.StoresNew)
	router.POST("/stores", h.StoresCreate)
	router.GET("/stores/:id/edit", h.StoresEdit)
	router.POST("/stores/:id", h.StoresUpdate)
	router.POST("/stores/:id/delete", h.StoresDelete)

	router.GET("/staffs", h.StaffsList)
	router.GET("/staffs/new", h.StaffsNew)
	router.POST("/staffs", h.StaffsCreate)
	router.GET("/staffs/:id/edit", h.StaffsEdit)
	router.POST("/staffs/:id", h.StaffsUpdate)
	router.POST("/staffs/:id/delete", h.StaffsDelete)

	router.GET("/settings", h.SettingsList)
	router.GET("/reports/sales", h.ReportsSales)

	// API routes
	api := router.Group("/api")
	{
		api.GET("/items", h.APIItemsList)
		api.GET("/items/:id", h.APIItemsGet)
		api.POST("/items", h.APIItemsCreate)
		api.PUT("/items/:id", h.APIItemsUpdate)
		api.DELETE("/items/:id", h.APIItemsDelete)

		api.GET("/sales", h.APISalesList)
		api.GET("/sales/:id", h.APISalesGet)
		api.POST("/sales", h.APISalesCreate)

		api.GET("/stores", h.APIStoresList)
		api.GET("/stores/:id", h.APIStoresGet)
		api.POST("/stores", h.APIStoresCreate)
		api.PUT("/stores/:id", h.APIStoresUpdate)
		api.DELETE("/stores/:id", h.APIStoresDelete)

		api.GET("/staffs", h.APIStaffsList)
		api.GET("/staffs/:id", h.APIStaffsGet)
		api.POST("/staffs", h.APIStaffsCreate)
		api.PUT("/staffs/:id", h.APIStaffsUpdate)
		api.DELETE("/staffs/:id", h.APIStaffsDelete)

		api.GET("/settings", h.APISettingsList)
		api.PUT("/settings/:key", h.APISettingsUpdate)

		api.GET("/reports/sales", h.APIReportsSales)
		api.GET("/reports/sales/excel", h.APIReportsSalesExcel)
	}
}
