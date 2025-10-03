package main

import (
	"log"
	"os"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/config"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/handlers"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/repository"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/service"
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
	handlers.SetupRoutes(router, h)

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

