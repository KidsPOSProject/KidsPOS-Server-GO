package handlers

import (
	"net/http"
	"testing"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/models"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/repository"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/service"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlers_APIItemsList(t *testing.T) {
	// Setup
	db := testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	// Create handlers
	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo)
	h := New(itemService, nil, nil, nil, nil)

	// Register route
	router.GET("/api/items", h.APIItemsList)

	// Create test data
	testItem := &models.Item{
		ItemID: "TEST001",
		Name:   "Test Item",
		Price:  100,
		Stock:  10,
	}
	err := itemRepo.Create(testItem)
	require.NoError(t, err)

	// Make request
	w := testutil.MakeRequest(router, "GET", "/api/items", nil)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Items []models.Item `json:"items"`
	}
	testutil.ParseJSONResponse(t, w, &response)

	assert.Len(t, response.Items, 1)
	assert.Equal(t, "Test Item", response.Items[0].Name)
	assert.Equal(t, 100, response.Items[0].Price)
}

func TestHandlers_APIItemsCreate(t *testing.T) {
	// Setup
	db := testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	// Create handlers
	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo)
	h := New(itemService, nil, nil, nil, nil)

	// Register route
	router.POST("/api/items", h.APIItemsCreate)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
		wantError  bool
	}{
		{
			name: "Valid item creation",
			payload: map[string]interface{}{
				"itemId": "NEW001",
				"name":   "New Item",
				"price":  200,
				"stock":  5,
			},
			wantStatus: http.StatusCreated,
			wantError:  false,
		},
		{
			name: "Invalid - empty name",
			payload: map[string]interface{}{
				"name":  "",
				"price": 100,
				"stock": 10,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "Invalid - negative price",
			payload: map[string]interface{}{
				"name":  "Invalid Item",
				"price": -100,
				"stock": 10,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name: "Invalid - negative stock",
			payload: map[string]interface{}{
				"name":  "Invalid Item",
				"price": 100,
				"stock": -10,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make request
			w := testutil.MakeRequest(router, "POST", "/api/items", tt.payload)

			// Assertions
			assert.Equal(t, tt.wantStatus, w.Code)

			if !tt.wantError {
				var response struct {
					Item models.Item `json:"item"`
				}
				testutil.ParseJSONResponse(t, w, &response)
				assert.Equal(t, tt.payload["name"], response.Item.Name)
				assert.Equal(t, tt.payload["price"].(int), response.Item.Price)
			}
		})
	}
}

func TestHandlers_APIItemsGet(t *testing.T) {
	// Setup
	db := testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	// Create handlers
	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo)
	h := New(itemService, nil, nil, nil, nil)

	// Register route
	router.GET("/api/items/:id", h.APIItemsGet)

	// Create test data
	testItem := &models.Item{
		ItemID: "TEST001",
		Name:   "Test Item",
		Price:  100,
		Stock:  10,
	}
	err := itemRepo.Create(testItem)
	require.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		wantStatus int
		wantError  bool
	}{
		{
			name:       "Existing item",
			id:         "1",
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name:       "Non-existent item",
			id:         "999",
			wantStatus: http.StatusNotFound,
			wantError:  true,
		},
		{
			name:       "Invalid ID",
			id:         "invalid",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make request
			w := testutil.MakeRequest(router, "GET", "/api/items/"+tt.id, nil)

			// Assertions
			assert.Equal(t, tt.wantStatus, w.Code)

			if !tt.wantError {
				var response struct {
					Item models.Item `json:"item"`
				}
				testutil.ParseJSONResponse(t, w, &response)
				assert.Equal(t, "Test Item", response.Item.Name)
			}
		})
	}
}

func TestHandlers_APIItemsUpdate(t *testing.T) {
	// Setup
	db := testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	// Create handlers
	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo)
	h := New(itemService, nil, nil, nil, nil)

	// Register route
	router.PUT("/api/items/:id", h.APIItemsUpdate)

	// Create test data
	testItem := &models.Item{
		ItemID: "TEST001",
		Name:   "Original Item",
		Price:  100,
		Stock:  10,
	}
	err := itemRepo.Create(testItem)
	require.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		payload    map[string]interface{}
		wantStatus int
		wantError  bool
	}{
		{
			name: "Valid update",
			id:   "1",
			payload: map[string]interface{}{
				"name":  "Updated Item",
				"price": 200,
				"stock": 20,
			},
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name: "Invalid - empty name",
			id:   "1",
			payload: map[string]interface{}{
				"name":  "",
				"price": 100,
				"stock": 10,
			},
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "Non-existent item",
			id:         "999",
			payload:    map[string]interface{}{"name": "Test"},
			wantStatus: http.StatusNotFound,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make request
			w := testutil.MakeRequest(router, "PUT", "/api/items/"+tt.id, tt.payload)

			// Assertions
			assert.Equal(t, tt.wantStatus, w.Code)

			if !tt.wantError {
				var response struct {
					Item models.Item `json:"item"`
				}
				testutil.ParseJSONResponse(t, w, &response)
				assert.Equal(t, tt.payload["name"], response.Item.Name)
			}
		})
	}
}

func TestHandlers_APIItemsDelete(t *testing.T) {
	// Setup
	db := testutil.SetupTestDB(t)
	router := testutil.SetupRouter()

	// Create handlers
	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo)
	h := New(itemService, nil, nil, nil, nil)

	// Register route
	router.DELETE("/api/items/:id", h.APIItemsDelete)

	// Create test data
	testItem := &models.Item{
		ItemID: "TEST001",
		Name:   "Test Item",
		Price:  100,
		Stock:  10,
	}
	err := itemRepo.Create(testItem)
	require.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "Existing item",
			id:         "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Non-existent item",
			id:         "999",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Invalid ID",
			id:         "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make request
			w := testutil.MakeRequest(router, "DELETE", "/api/items/"+tt.id, nil)

			// Assertions
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}