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
		CREATE TABLE sale (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			staffId INTEGER NOT NULL,
			storeId INTEGER NOT NULL,
			totalPrice INTEGER NOT NULL,
			deposit INTEGER NOT NULL,
			saleAt DATETIME NOT NULL,
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (staffId) REFERENCES staff(id)
		)
	`)
	require.NoError(t, err)

	return db
}

func TestStaffRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &StaffRepository{db: db}

	// Create a staff first
	staff := &models.Staff{
		StaffID: "STAFF-001",
		Name:    "Test Staff",
	}
	_, err := db.Exec("INSERT INTO staff (staffId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
		staff.StaffID, staff.Name, time.Now(), time.Now())
	require.NoError(t, err)

	// Get the created staff
	var id int
	err = db.QueryRow("SELECT id FROM staff WHERE staffId = ?", staff.StaffID).Scan(&id)
	require.NoError(t, err)
	staff.ID = id

	t.Run("successful update", func(t *testing.T) {
		staff.Name = "Updated Staff Name"
		err := repo.Update(staff)
		assert.NoError(t, err)

		// Verify the update
		var name string
		err = db.QueryRow("SELECT name FROM staff WHERE id = ?", staff.ID).Scan(&name)
		require.NoError(t, err)
		assert.Equal(t, "Updated Staff Name", name)
	})

	t.Run("update non-existent staff", func(t *testing.T) {
		nonExistent := &models.Staff{
			ID:   999,
			Name: "Non-existent",
		}
		err := repo.Update(nonExistent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestStaffRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &StaffRepository{db: db}

	t.Run("successful delete", func(t *testing.T) {
		// Create a staff
		result, err := db.Exec("INSERT INTO staff (staffId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
			"STAFF-002", "Test Staff 2", time.Now(), time.Now())
		require.NoError(t, err)
		id, _ := result.LastInsertId()

		// Delete the staff
		err = repo.Delete(int(id))
		assert.NoError(t, err)

		// Verify deletion
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM staff WHERE id = ?", id).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("delete staff with sales - should fail", func(t *testing.T) {
		// Create a staff
		result, err := db.Exec("INSERT INTO staff (staffId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)",
			"STAFF-003", "Test Staff 3", time.Now(), time.Now())
		require.NoError(t, err)
		staffID, _ := result.LastInsertId()

		// Create a sale referencing this staff
		_, err = db.Exec("INSERT INTO sale (staffId, storeId, totalPrice, deposit, saleAt) VALUES (?, ?, ?, ?, ?)",
			staffID, 1, 1000, 1000, time.Now())
		require.NoError(t, err)

		// Try to delete the staff - should fail
		err = repo.Delete(int(staffID))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "referenced by")
	})

	t.Run("delete non-existent staff", func(t *testing.T) {
		err := repo.Delete(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
