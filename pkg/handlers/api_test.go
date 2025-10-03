package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/models"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/repository"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE staff (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			staffId TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE store (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			storeId TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE sale (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			staffId INTEGER NOT NULL,
			storeId INTEGER NOT NULL,
			totalPrice INTEGER NOT NULL,
			deposit INTEGER NOT NULL,
			saleAt DATETIME NOT NULL,
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (staffId) REFERENCES staff(id),
			FOREIGN KEY (storeId) REFERENCES store(id)
		)
	`)
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE apk_versions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version TEXT NOT NULL UNIQUE,
			versionCode INTEGER NOT NULL,
			fileName TEXT NOT NULL,
			fileSize INTEGER NOT NULL,
			filePath TEXT NOT NULL,
			releaseNotes TEXT,
			isActive INTEGER NOT NULL DEFAULT 1,
			uploadedAt DATETIME NOT NULL,
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	return db
}

func setupTestRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repos := repository.NewRepositories(db)
	services := service.NewServices(repos)
	handlers := NewHandlers(services)

	SetupRoutes(router, handlers)

	return router
}

func TestAPIStaffsList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	// Insert test data
	_, err := db.Exec("INSERT INTO staff (staffId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
		"STAFF-001", "Test Staff 1", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO staff (staffId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
		"STAFF-002", "Test Staff 2", time.Now(), time.Now())
	require.NoError(t, err)

	t.Run("get all staffs", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/staffs", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var staffs []*models.Staff
		err := json.Unmarshal(w.Body.Bytes(), &staffs)
		require.NoError(t, err)
		assert.Len(t, staffs, 2)
		assert.Equal(t, "Test Staff 2", staffs[0].Name) // DESC order
	})
}

func TestAPIStaffsGet(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	// Insert test data
	result, err := db.Exec("INSERT INTO staff (staffId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
		"STAFF-001", "Test Staff", time.Now(), time.Now())
	require.NoError(t, err)
	id, _ := result.LastInsertId()

	t.Run("get existing staff", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/staffs/%d", id), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var staff models.Staff
		err := json.Unmarshal(w.Body.Bytes(), &staff)
		require.NoError(t, err)
		assert.Equal(t, "Test Staff", staff.Name)
	})

	t.Run("get non-existent staff", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/staffs/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAPIStaffsCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("create staff successfully", func(t *testing.T) {
		staffData := map[string]interface{}{
			"name": "New Staff",
		}
		body, _ := json.Marshal(staffData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/staffs", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var staff models.Staff
		err := json.Unmarshal(w.Body.Bytes(), &staff)
		require.NoError(t, err)
		assert.Equal(t, "New Staff", staff.Name)
		assert.NotEmpty(t, staff.StaffID)
	})

	t.Run("create staff with empty name", func(t *testing.T) {
		staffData := map[string]interface{}{
			"name": "",
		}
		body, _ := json.Marshal(staffData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/staffs", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAPIStaffsUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	// Insert test data
	result, err := db.Exec("INSERT INTO staff (staffId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
		"STAFF-001", "Original Name", time.Now(), time.Now())
	require.NoError(t, err)
	id, _ := result.LastInsertId()

	t.Run("update staff successfully", func(t *testing.T) {
		staffData := map[string]interface{}{
			"name": "Updated Name",
		}
		body, _ := json.Marshal(staffData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/staffs/%d", id), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var staff models.Staff
		err := json.Unmarshal(w.Body.Bytes(), &staff)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", staff.Name)
	})

	t.Run("update non-existent staff", func(t *testing.T) {
		staffData := map[string]interface{}{
			"name": "Updated Name",
		}
		body, _ := json.Marshal(staffData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/api/staffs/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAPIStaffsDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("delete staff successfully", func(t *testing.T) {
		// Insert test data
		result, err := db.Exec("INSERT INTO staff (staffId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
			"STAFF-DEL", "To Delete", time.Now(), time.Now())
		require.NoError(t, err)
		id, _ := result.LastInsertId()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/staffs/%d", id), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify deletion
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM staff WHERE id = ?", id).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("delete staff with sales - should fail", func(t *testing.T) {
		// Insert staff
		result, err := db.Exec("INSERT INTO staff (staffId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
			"STAFF-REF", "Has Sales", time.Now(), time.Now())
		require.NoError(t, err)
		staffID, _ := result.LastInsertId()

		// Insert store
		storeResult, err := db.Exec("INSERT INTO store (storeId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
			"STORE-001", "Test Store", time.Now(), time.Now())
		require.NoError(t, err)
		storeID, _ := storeResult.LastInsertId()

		// Insert sale referencing staff
		_, err = db.Exec("INSERT INTO sale (staffId, storeId, totalPrice, deposit, saleAt) VALUES (?, ?, ?, ?, ?)",
			staffID, storeID, 1000, 1000, time.Now())
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/staffs/%d", staffID), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "referenced by")
	})

	t.Run("delete non-existent staff", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/staffs/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// Store API Tests
func TestAPIStoresDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("delete store successfully", func(t *testing.T) {
		// Insert test data
		result, err := db.Exec("INSERT INTO store (storeId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
			"STORE-DEL", "To Delete", time.Now(), time.Now())
		require.NoError(t, err)
		id, _ := result.LastInsertId()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/stores/%d", id), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify deletion
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM store WHERE id = ?", id).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("delete store with sales - should fail", func(t *testing.T) {
		// Insert store
		result, err := db.Exec("INSERT INTO store (storeId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
			"STORE-REF", "Has Sales", time.Now(), time.Now())
		require.NoError(t, err)
		storeID, _ := result.LastInsertId()

		// Insert staff
		staffResult, err := db.Exec("INSERT INTO staff (staffId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
			"STAFF-001", "Test Staff", time.Now(), time.Now())
		require.NoError(t, err)
		staffID, _ := staffResult.LastInsertId()

		// Insert sale referencing store
		_, err = db.Exec("INSERT INTO sale (staffId, storeId, totalPrice, deposit, saleAt) VALUES (?, ?, ?, ?, ?)",
			staffID, storeID, 1000, 1000, time.Now())
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/stores/%d", storeID), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "referenced by")
	})
}

// APK API Tests
func TestAPIApkLatest(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("no versions available", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/apk/version/latest", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("get latest version", func(t *testing.T) {
		// Insert test data
		_, err := db.Exec("INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			"1.0.0", 1, "test1.apk", 1000, "/tmp/test1.apk", "First release", 1, time.Now())
		require.NoError(t, err)

		_, err = db.Exec("INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			"1.1.0", 2, "test2.apk", 2000, "/tmp/test2.apk", "Second release", 1, time.Now())
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/apk/version/latest", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var apk models.ApkVersion
		err = json.Unmarshal(w.Body.Bytes(), &apk)
		require.NoError(t, err)
		assert.Equal(t, "1.1.0", apk.Version)
		assert.Equal(t, 2, apk.VersionCode)
	})
}

func TestAPIApkCheckUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	// Insert test data
	_, err := db.Exec("INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"1.0.0", 1, "test1.apk", 1000, "/tmp/test1.apk", "First release", 1, time.Now())
	require.NoError(t, err)

	_, err = db.Exec("INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"2.0.0", 2, "test2.apk", 2000, "/tmp/test2.apk", "Second release", 1, time.Now())
	require.NoError(t, err)

	t.Run("update available", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/apk/version/check?currentVersionCode=1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response["hasUpdate"].(bool))
		assert.NotNil(t, response["latestVersion"])
	})

	t.Run("no update available", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/apk/version/check?currentVersionCode=2", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.False(t, response["hasUpdate"].(bool))
		assert.Nil(t, response["latestVersion"])
	})

	t.Run("missing version code", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/apk/version/check", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAPIApkVersions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("get all versions", func(t *testing.T) {
		// Insert test data
		_, err := db.Exec("INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			"1.0.0", 1, "test1.apk", 1000, "/tmp/test1.apk", "First", 1, time.Now())
		require.NoError(t, err)

		_, err = db.Exec("INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			"2.0.0", 2, "test2.apk", 2000, "/tmp/test2.apk", "Second", 1, time.Now())
		require.NoError(t, err)

		_, err = db.Exec("INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			"3.0.0", 3, "test3.apk", 3000, "/tmp/test3.apk", "Third", 0, time.Now())
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/apk/version/all", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var versions []*models.ApkVersion
		err = json.Unmarshal(w.Body.Bytes(), &versions)
		require.NoError(t, err)
		// Only active versions should be returned
		assert.Len(t, versions, 2)
	})
}

func TestAPIApkDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("delete version successfully", func(t *testing.T) {
		// Insert test data
		result, err := db.Exec("INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			"1.0.0", 1, "test.apk", 1000, "/tmp/test.apk", "Test", 1, time.Now())
		require.NoError(t, err)
		id, _ := result.LastInsertId()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/apk/version/%d", id), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify deletion
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM apk_versions WHERE id = ?", id).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

func TestAPIApkDeactivate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("deactivate version successfully", func(t *testing.T) {
		// Insert test data
		result, err := db.Exec("INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			"1.0.0", 1, "test.apk", 1000, "/tmp/test.apk", "Test", 1, time.Now())
		require.NoError(t, err)
		id, _ := result.LastInsertId()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/apk/version/%d/deactivate", id), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify deactivation
		var isActive int
		err = db.QueryRow("SELECT isActive FROM apk_versions WHERE id = ?", id).Scan(&isActive)
		require.NoError(t, err)
		assert.Equal(t, 0, isActive)

		var apk models.ApkVersion
		err = json.Unmarshal(w.Body.Bytes(), &apk)
		require.NoError(t, err)
		assert.False(t, apk.IsActive)
	})
}

