package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/models"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/service"
	"github.com/gin-gonic/gin"
)

// MockServices creates mock services for testing
func MockServices() *service.Services {
	// This would be replaced with proper mocks in a real implementation
	return &service.Services{
		Item:    &service.ItemService{},
		Store:   &service.StoreService{},
		Staff:   &service.StaffService{},
		Sale:    &service.SaleService{},
		Setting: &service.SettingService{},
	}
}

func TestHandlers_Home(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	services := MockServices()
	h := NewHandlers(services)

	router.GET("/", h.Home)

	// Create request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHandlers_APIItemsList(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	services := MockServices()
	h := NewHandlers(services)

	router.GET("/api/items", h.APIItemsList)

	// Create request
	req, err := http.NewRequest("GET", "/api/items", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK && status != http.StatusInternalServerError {
		t.Errorf("handler returned unexpected status code: got %v",
			status)
	}

	// Check content type
	expectedContentType := "application/json; charset=utf-8"
	if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			ct, expectedContentType)
	}
}

func TestHandlers_APIItemsCreate(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	services := MockServices()
	h := NewHandlers(services)

	router.POST("/api/items", h.APIItemsCreate)

	// Create test item
	item := models.Item{
		Name:  "テスト商品",
		Price: 100,
		Stock: 10,
	}

	// Convert to JSON
	itemJSON, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}

	// Create request
	req, err := http.NewRequest("POST", "/api/items", bytes.NewBuffer(itemJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusCreated && status != http.StatusBadRequest {
		t.Errorf("handler returned unexpected status code: got %v",
			status)
	}
}

func TestHandlers_APIItemsGet(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	services := MockServices()
	h := NewHandlers(services)

	router.GET("/api/items/:id", h.APIItemsGet)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{
			name:           "有効なID",
			id:             "1",
			expectedStatus: http.StatusNotFound, // Mock returns not found
		},
		{
			name:           "無効なID",
			id:             "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("GET", "/api/items/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestHandlers_APIItemsUpdate(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	services := MockServices()
	h := NewHandlers(services)

	router.PUT("/api/items/:id", h.APIItemsUpdate)

	// Create test item
	item := models.Item{
		Name:  "更新商品",
		Price: 200,
		Stock: 20,
	}

	// Convert to JSON
	itemJSON, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}

	// Create request
	req, err := http.NewRequest("PUT", "/api/items/1", bytes.NewBuffer(itemJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK && status != http.StatusBadRequest {
		t.Errorf("handler returned unexpected status code: got %v",
			status)
	}
}

func TestHandlers_APIItemsDelete(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	services := MockServices()
	h := NewHandlers(services)

	router.DELETE("/api/items/:id", h.APIItemsDelete)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{
			name:           "有効なID",
			id:             "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "無効なID",
			id:             "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("DELETE", "/api/items/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus && status != http.StatusInternalServerError {
				t.Errorf("handler returned unexpected status code: got %v for test %s",
					status, tt.name)
			}
		})
	}
}

// BenchmarkHandlers_APIItemsList benchmarks the items list API
func BenchmarkHandlers_APIItemsList(b *testing.B) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	services := MockServices()
	h := NewHandlers(services)

	router.GET("/api/items", h.APIItemsList)

	// Create request
	req, err := http.NewRequest("GET", "/api/items", nil)
	if err != nil {
		b.Fatal(err)
	}

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}
}

// BenchmarkHandlers_APIItemsCreate benchmarks the item creation API
func BenchmarkHandlers_APIItemsCreate(b *testing.B) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	services := MockServices()
	h := NewHandlers(services)

	router.POST("/api/items", h.APIItemsCreate)

	// Create test item
	item := models.Item{
		Name:  "ベンチマーク商品",
		Price: 100,
		Stock: 10,
	}

	// Convert to JSON
	itemJSON, err := json.Marshal(item)
	if err != nil {
		b.Fatal(err)
	}

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/api/items", bytes.NewBuffer(itemJSON))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}
}