package testutil

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/repository"
	"github.com/gin-gonic/gin"
)

// SetupTestDB creates a test database
func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Create temp db file
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()

	// Initialize database
	db, err := repository.InitDB(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up on test completion
	t.Cleanup(func() {
		db.Close()
		os.Remove(dbPath)
	})

	return db
}

// SetupRouter creates a test router
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// MakeRequest makes a test HTTP request
func MakeRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	router.ServeHTTP(w, req)

	return w
}

// ParseJSONResponse parses JSON response body
func ParseJSONResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	t.Helper()

	if err := json.Unmarshal(w.Body.Bytes(), target); err != nil {
		t.Fatalf("Failed to parse JSON response: %v\nBody: %s", err, w.Body.String())
	}
}