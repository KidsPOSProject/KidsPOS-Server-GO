package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupStoreTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	// Create tables
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
			FOREIGN KEY (storeId) REFERENCES store(id)
		)
	`)
	require.NoError(t, err)

	return db
}

func TestStoreRepository_Update(t *testing.T) {
	db := setupStoreTestDB(t)
	defer db.Close()

	repo := &StoreRepository{db: db}

	// Create a store first
	store := &models.Store{
		StoreID: "STORE-001",
		Name:    "Test Store",
	}
	_, err := db.Exec("INSERT INTO store (storeId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
		store.StoreID, store.Name, time.Now(), time.Now())
	require.NoError(t, err)

	// Get the created store
	var id int
	err = db.QueryRow("SELECT id FROM store WHERE storeId = ?", store.StoreID).Scan(&id)
	require.NoError(t, err)
	store.ID = id

	t.Run("successful update", func(t *testing.T) {
		store.Name = "Updated Store Name"
		err := repo.Update(store)
		assert.NoError(t, err)

		// Verify the update
		var name string
		err = db.QueryRow("SELECT name FROM store WHERE id = ?", store.ID).Scan(&name)
		require.NoError(t, err)
		assert.Equal(t, "Updated Store Name", name)
	})

	t.Run("update non-existent store", func(t *testing.T) {
		nonExistent := &models.Store{
			ID:   999,
			Name: "Non-existent",
		}
		err := repo.Update(nonExistent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestStoreRepository_Delete(t *testing.T) {
	db := setupStoreTestDB(t)
	defer db.Close()

	repo := &StoreRepository{db: db}

	t.Run("successful delete", func(t *testing.T) {
		// Create a store
		result, err := db.Exec("INSERT INTO store (storeId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
			"STORE-002", "Test Store 2", time.Now(), time.Now())
		require.NoError(t, err)
		id, _ := result.LastInsertId()

		// Delete the store
		err = repo.Delete(int(id))
		assert.NoError(t, err)

		// Verify deletion
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM store WHERE id = ?", id).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("delete store with sales - should fail", func(t *testing.T) {
		// Create a store
		result, err := db.Exec("INSERT INTO store (storeId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
			"STORE-003", "Test Store 3", time.Now(), time.Now())
		require.NoError(t, err)
		storeID, _ := result.LastInsertId()

		// Create a sale referencing this store
		_, err = db.Exec("INSERT INTO sale (staffId, storeId, totalPrice, deposit, saleAt) VALUES (?, ?, ?, ?, ?)",
			1, storeID, 1000, 1000, time.Now())
		require.NoError(t, err)

		// Try to delete the store - should fail
		err = repo.Delete(int(storeID))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "referenced by")
	})

	t.Run("delete non-existent store", func(t *testing.T) {
		err := repo.Delete(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
