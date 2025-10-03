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

func setupAPKTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	// Create apk_versions table
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

func createTestAPK(t *testing.T, db *sql.DB, version string, versionCode int, isActive bool) int {
	activeInt := 0
	if isActive {
		activeInt = 1
	}

	result, err := db.Exec(`
		INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, version, versionCode, "test.apk", 1000, "/path/to/test.apk", "Test release", activeInt, time.Now())
	require.NoError(t, err)

	id, err := result.LastInsertId()
	require.NoError(t, err)

	return int(id)
}

func TestApkRepository_Create(t *testing.T) {
	db := setupAPKTestDB(t)
	defer db.Close()

	repo := &ApkVersionRepository{db: db}

	t.Run("successful create", func(t *testing.T) {
		apk := &models.ApkVersion{
			Version:      "1.0.0",
			VersionCode:  1,
			FileName:     "app-1.0.0.apk",
			FileSize:     1024000,
			FilePath:     "/uploads/app-1.0.0.apk",
			ReleaseNotes: "Initial release",
			IsActive:     true,
			UploadedAt:   time.Now(),
		}

		err := repo.Create(apk)
		assert.NoError(t, err)
		assert.NotZero(t, apk.ID)
	})

	t.Run("create duplicate version - should fail", func(t *testing.T) {
		apk := &models.ApkVersion{
			Version:      "1.0.0", // Same as above
			VersionCode:  2,
			FileName:     "app-1.0.0-2.apk",
			FileSize:     1024000,
			FilePath:     "/uploads/app-1.0.0-2.apk",
			ReleaseNotes: "Duplicate version",
			IsActive:     true,
			UploadedAt:   time.Now(),
		}

		err := repo.Create(apk)
		assert.Error(t, err)
	})
}

func TestApkRepository_FindLatest(t *testing.T) {
	db := setupAPKTestDB(t)
	defer db.Close()

	repo := &ApkVersionRepository{db: db}

	t.Run("no versions available", func(t *testing.T) {
		version, err := repo.FindLatest()
		assert.NoError(t, err)
		assert.Nil(t, version)
	})

	t.Run("return latest active version", func(t *testing.T) {
		// Create multiple versions
		createTestAPK(t, db, "1.0.0", 1, true)
		createTestAPK(t, db, "1.1.0", 2, true)
		id3 := createTestAPK(t, db, "1.2.0", 3, true)
		createTestAPK(t, db, "2.0.0", 4, false) // Inactive

		version, err := repo.FindLatest()
		require.NoError(t, err)
		require.NotNil(t, version)
		assert.Equal(t, id3, version.ID)
		assert.Equal(t, "1.2.0", version.Version)
		assert.Equal(t, 3, version.VersionCode)
	})
}

func TestApkRepository_FindByID(t *testing.T) {
	db := setupAPKTestDB(t)
	defer db.Close()

	repo := &ApkVersionRepository{db: db}

	t.Run("find existing version", func(t *testing.T) {
		id := createTestAPK(t, db, "1.0.0", 1, true)

		version, err := repo.FindByID(id)
		require.NoError(t, err)
		require.NotNil(t, version)
		assert.Equal(t, id, version.ID)
		assert.Equal(t, "1.0.0", version.Version)
	})

	t.Run("find non-existent version", func(t *testing.T) {
		version, err := repo.FindByID(999)
		assert.Error(t, err)
		assert.Nil(t, version)
	})
}

func TestApkRepository_FindAll(t *testing.T) {
	db := setupAPKTestDB(t)
	defer db.Close()

	repo := &ApkVersionRepository{db: db}

	t.Run("no versions", func(t *testing.T) {
		versions, err := repo.FindAll()
		require.NoError(t, err)
		assert.Empty(t, versions)
	})

	t.Run("return only active versions", func(t *testing.T) {
		createTestAPK(t, db, "1.0.0", 1, true)
		createTestAPK(t, db, "1.1.0", 2, true)
		createTestAPK(t, db, "1.2.0", 3, false) // Inactive

		versions, err := repo.FindAll()
		require.NoError(t, err)
		assert.Len(t, versions, 2)
	})

	t.Run("ordered by versionCode desc", func(t *testing.T) {
		versions, err := repo.FindAll()
		require.NoError(t, err)
		require.Len(t, versions, 2)

		// Should be ordered newest first
		assert.Equal(t, "1.1.0", versions[0].Version)
		assert.Equal(t, "1.0.0", versions[1].Version)
	})
}

func TestApkRepository_FindByVersionCode(t *testing.T) {
	db := setupAPKTestDB(t)
	defer db.Close()

	repo := &ApkVersionRepository{db: db}

	t.Run("no newer versions", func(t *testing.T) {
		createTestAPK(t, db, "1.0.0", 1, true)

		version, err := repo.FindByVersionCode(1)
		assert.NoError(t, err)
		assert.Nil(t, version)
	})

	t.Run("find newer version", func(t *testing.T) {
		createTestAPK(t, db, "1.1.0", 2, true)
		id3 := createTestAPK(t, db, "1.2.0", 3, true)

		version, err := repo.FindByVersionCode(1)
		require.NoError(t, err)
		require.NotNil(t, version)
		assert.Equal(t, id3, version.ID)
		assert.Equal(t, "1.2.0", version.Version)
	})

	t.Run("ignore inactive versions", func(t *testing.T) {
		createTestAPK(t, db, "2.0.0", 4, false) // Inactive

		version, err := repo.FindByVersionCode(3)
		assert.NoError(t, err)
		assert.Nil(t, version)
	})
}

func TestApkRepository_Update(t *testing.T) {
	db := setupAPKTestDB(t)
	defer db.Close()

	repo := &ApkVersionRepository{db: db}

	t.Run("successful update", func(t *testing.T) {
		id := createTestAPK(t, db, "1.0.0", 1, true)

		apk := &models.ApkVersion{
			ID:           id,
			ReleaseNotes: "Updated release notes",
		}

		err := repo.Update(apk)
		assert.NoError(t, err)

		// Verify update
		var notes string
		err = db.QueryRow("SELECT releaseNotes FROM apk_versions WHERE id = ?", id).Scan(&notes)
		require.NoError(t, err)
		assert.Equal(t, "Updated release notes", notes)
	})

	t.Run("update non-existent version", func(t *testing.T) {
		apk := &models.ApkVersion{
			ID:           999,
			ReleaseNotes: "Will fail",
		}

		err := repo.Update(apk)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestApkRepository_Delete(t *testing.T) {
	db := setupAPKTestDB(t)
	defer db.Close()

	repo := &ApkVersionRepository{db: db}

	t.Run("successful delete", func(t *testing.T) {
		id := createTestAPK(t, db, "1.0.0", 1, true)

		err := repo.Delete(id)
		assert.NoError(t, err)

		// Verify deletion
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM apk_versions WHERE id = ?", id).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("delete non-existent version", func(t *testing.T) {
		err := repo.Delete(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestApkRepository_Deactivate(t *testing.T) {
	db := setupAPKTestDB(t)
	defer db.Close()

	repo := &ApkVersionRepository{db: db}

	t.Run("successful deactivate", func(t *testing.T) {
		id := createTestAPK(t, db, "1.0.0", 1, true)

		apk, err := repo.Deactivate(id)
		require.NoError(t, err)
		require.NotNil(t, apk)
		assert.False(t, apk.IsActive)

		// Verify in database
		var isActive bool
		err = db.QueryRow("SELECT isActive FROM apk_versions WHERE id = ?", id).Scan(&isActive)
		require.NoError(t, err)
		assert.False(t, isActive)
	})

	t.Run("deactivate non-existent version", func(t *testing.T) {
		apk, err := repo.Deactivate(999)
		assert.Error(t, err)
		assert.Nil(t, apk)
	})
}
