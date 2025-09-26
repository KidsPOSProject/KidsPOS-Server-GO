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
	router.LoadHTMLGlob("web/templates/*")

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

	// API routes - Compatible with original Kotlin version
	api := router.Group("/api")
	{
		// Item API - /api/item (singular, matching Kotlin version)
		api.GET("/item", h.APIItemsList)
		api.GET("/item/:id", h.APIItemsGet)
		api.GET("/item/barcode/:barcode", h.APIItemsByBarcode)
		api.POST("/item", h.APIItemsCreate)
		api.PUT("/item/:id", h.APIItemsUpdate)
		api.PATCH("/item/:id", h.APIItemsPatch)
		api.DELETE("/item/:id", h.APIItemsDelete)
		api.GET("/item/barcode-pdf", h.APIItemsBarcodePDF)
		api.POST("/item/barcode-pdf/selected", h.APIItemsBarcodeSelectedPDF)

		// Sales API - /api/sales (matching Kotlin version)
		api.POST("/sales", h.APISalesCreate)
		api.POST("/sales/create", h.APISalesCreateAlt)
		api.GET("/sales/:id", h.APISalesGet)
		api.GET("/sales", h.APISalesList)
		api.GET("/sales/validate-printer/:storeId", h.APISalesValidatePrinter)

		// Staff API - /api/staff (singular, matching Kotlin version)
		api.GET("/staff/:barcode", h.APIStaffByBarcode)
		api.GET("/staff", h.APIStaffsList)
		api.POST("/staff", h.APIStaffsCreate)
		api.PUT("/staff/:barcode", h.APIStaffUpdateByBarcode)
		api.DELETE("/staff/:barcode", h.APIStaffDeleteByBarcode)

		// Store API - /api/stores (matching Kotlin version)
		api.GET("/stores", h.APIStoresList)
		api.POST("/stores", h.APIStoresCreate)
		api.GET("/stores/:id", h.APIStoresGet)
		api.PUT("/stores/:id", h.APIStoresUpdate)
		api.DELETE("/stores/:id", h.APIStoresDelete)

		// Setting API - /api/setting (singular, matching Kotlin version)
		api.GET("/setting/status", h.APISettingsStatus)
		api.GET("/setting", h.APISettingsList)
		api.GET("/setting/:key", h.APISettingsGet)
		api.POST("/setting", h.APISettingsCreate)
		api.PUT("/setting/:key", h.APISettingsUpdate)
		api.DELETE("/setting/:key", h.APISettingsDelete)
		api.POST("/setting/printer/:storeId", h.APISettingsPrinterSet)
		api.GET("/setting/printer/:storeId", h.APISettingsPrinterGet)
		api.POST("/setting/application", h.APISettingsApplicationSet)
		api.GET("/setting/application", h.APISettingsApplicationGet)

		// User API - /api/users (matching Kotlin version)
		api.GET("/users", h.APIUsersList)
		api.GET("/users/:barcode", h.APIUsersByBarcode)
		api.POST("/users", h.APIUsersCreate)
		api.PUT("/users/:barcode", h.APIUsersUpdateByBarcode)
		api.DELETE("/users/:barcode", h.APIUsersDeleteByBarcode)

		// Reports API - /api/reports (matching Kotlin version)
		api.GET("/reports/sales/pdf", h.APIReportsSalesPDF)
		api.GET("/reports/sales/pdf/today", h.APIReportsSalesPDFToday)
		api.GET("/reports/sales/pdf/month", h.APIReportsSalesPDFMonth)
		api.GET("/reports/sales/excel", h.APIReportsSalesExcel)
		api.GET("/reports/sales/excel/today", h.APIReportsSalesExcelToday)
		api.GET("/reports/sales/excel/month", h.APIReportsSalesExcelMonth)
	}
}
