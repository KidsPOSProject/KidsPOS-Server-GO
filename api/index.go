package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/config"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/handlers"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/repository"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/service"
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

	// Load HTML templates if they exist
	templatesPath := "web/templates"
	if _, err := os.Stat(templatesPath); err == nil {
		pattern := filepath.Join(templatesPath, "*")
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			engine.LoadHTMLGlob(pattern)
			log.Printf("Loaded %d HTML templates from %s", len(matches), pattern)
		} else {
			log.Printf("No HTML templates found in %s, skipping template loading", templatesPath)
		}
	} else {
		log.Printf("Templates directory %s not found, skipping template loading", templatesPath)
	}

	// Serve static files if they exist
	staticPath := "./web/static"
	if _, err := os.Stat(staticPath); err == nil {
		engine.Static("/static", staticPath)
		engine.Static("/css", filepath.Join(staticPath, "css"))
		engine.Static("/js", filepath.Join(staticPath, "js"))
	}

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
