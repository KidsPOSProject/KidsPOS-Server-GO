package handler

import (
	"log"
	"net/http"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/config"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/handlers"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/repository"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/service"
	"github.com/gin-gonic/gin"
)

var (
	engine *gin.Engine
)

// init関数でGinエンジンを初期化します
// Vercelのサーバーレス環境では、この初期化はコールドスタート時に一度だけ実行されます
func init() {
	// Initialize configuration
	cfg := config.New()

	// Initialize database
	db, err := repository.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Run migrations
	if err := repository.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	services := service.NewServices(repos)

	// Initialize Gin router
	engine = gin.Default()

	// Load HTML templates
	engine.LoadHTMLGlob("web/templates/*")

	// Serve static files
	engine.Static("/static", "./web/static")
	engine.Static("/css", "./web/static/css")
	engine.Static("/js", "./web/static/js")

	// Initialize handlers
	h := handlers.NewHandlers(services)

	// Setup routes
	handlers.SetupRoutes(engine, h)
}

// Handler はVercelがリクエストを処理するために呼び出す関数です
func Handler(w http.ResponseWriter, r *http.Request) {
	// Ginエンジンにリクエストを渡して処理させます
	engine.ServeHTTP(w, r)
}
