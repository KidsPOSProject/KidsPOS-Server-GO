package service

import (
	"bytes"
	"database/sql"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
	"time"
	"unsafe"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupAPKServiceTestDB(t *testing.T) *sql.DB {
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

func createTestFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	require.NoError(t, err)

	_, err = io.Copy(part, bytes.NewReader(content))
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Parse the multipart form
	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(10 << 20) // 10MB
	require.NoError(t, err)

	files := form.File["file"]
	require.Len(t, files, 1)

	return files[0]
}

func TestApkVersionService_UploadApk(t *testing.T) {
	db := setupAPKServiceTestDB(t)
	defer db.Close()

	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "apk-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := &repository.ApkVersionRepository{}
	service := NewApkVersionService(repo)
	service.uploadDir = tmpDir

	// Set the repository's db
	repoStruct := (*struct {
		db *sql.DB
	})(unsafe.Pointer(repo))
	repoStruct.db = db

	t.Run("successful upload", func(t *testing.T) {
		fileHeader := createTestFileHeader(t, "test.apk", []byte("test content"))

		apk, err := service.UploadApk(fileHeader, "1.0.0", 1, "Test release")
		require.NoError(t, err)
		require.NotNil(t, apk)

		assert.Equal(t, "1.0.0", apk.Version)
		assert.Equal(t, 1, apk.VersionCode)
		assert.Equal(t, "Test release", apk.ReleaseNotes)
		assert.True(t, apk.IsActive)

		// Verify file was saved
		_, err = os.Stat(apk.FilePath)
		assert.NoError(t, err)
	})

	t.Run("missing version", func(t *testing.T) {
		fileHeader := createTestFileHeader(t, "test.apk", []byte("test content"))

		apk, err := service.UploadApk(fileHeader, "", 1, "Test release")
		assert.Error(t, err)
		assert.Nil(t, apk)
		assert.Contains(t, err.Error(), "version is required")
	})

	t.Run("invalid version code", func(t *testing.T) {
		fileHeader := createTestFileHeader(t, "test.apk", []byte("test content"))

		apk, err := service.UploadApk(fileHeader, "1.0.0", 0, "Test release")
		assert.Error(t, err)
		assert.Nil(t, apk)
		assert.Contains(t, err.Error(), "version code must be positive")
	})

	t.Run("file too large", func(t *testing.T) {
		// Create a large content (more than 100MB)
		largeContent := make([]byte, 101*1024*1024)
		fileHeader := createTestFileHeader(t, "large.apk", largeContent)

		apk, err := service.UploadApk(fileHeader, "1.0.0", 1, "Test release")
		assert.Error(t, err)
		assert.Nil(t, apk)
		assert.Contains(t, err.Error(), "file size exceeds")
	})

	t.Run("non-apk file", func(t *testing.T) {
		fileHeader := createTestFileHeader(t, "test.txt", []byte("test content"))

		apk, err := service.UploadApk(fileHeader, "1.0.0", 1, "Test release")
		assert.Error(t, err)
		assert.Nil(t, apk)
		assert.Contains(t, err.Error(), "must be an APK file")
	})
}

func TestApkVersionService_GetLatestVersion(t *testing.T) {
	db := setupAPKServiceTestDB(t)
	defer db.Close()

	repo := &repository.ApkVersionRepository{}
	service := NewApkVersionService(repo)

	// Set the repository's db using reflection
	repoStruct := (*struct {
		db *sql.DB
	})(unsafe.Pointer(repo))
	repoStruct.db = db

	t.Run("no versions", func(t *testing.T) {
		version, err := service.GetLatestVersion()
		assert.NoError(t, err)
		assert.Nil(t, version)
	})

	t.Run("get latest version", func(t *testing.T) {
		// Insert test data
		_, err := db.Exec(`INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, "1.0.0", 1, "test.apk", 1000, "/path", "notes", 1, time.Now())
		require.NoError(t, err)

		_, err = db.Exec(`INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, "1.1.0", 2, "test.apk", 1000, "/path", "notes", 1, time.Now())
		require.NoError(t, err)

		version, err := service.GetLatestVersion()
		require.NoError(t, err)
		require.NotNil(t, version)
		assert.Equal(t, "1.1.0", version.Version)
	})
}

func TestApkVersionService_CheckForUpdate(t *testing.T) {
	db := setupAPKServiceTestDB(t)
	defer db.Close()

	repo := &repository.ApkVersionRepository{}
	service := NewApkVersionService(repo)

	// Set the repository's db
	repoStruct := (*struct {
		db *sql.DB
	})(unsafe.Pointer(repo))
	repoStruct.db = db

	// Insert test data
	_, err := db.Exec(`INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, "1.0.0", 1, "test.apk", 1000, "/path", "notes", 1, time.Now())
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, "2.0.0", 2, "test.apk", 1000, "/path", "notes", 1, time.Now())
	require.NoError(t, err)

	t.Run("update available", func(t *testing.T) {
		version, err := service.CheckForUpdate(1)
		require.NoError(t, err)
		require.NotNil(t, version)
		assert.Equal(t, "2.0.0", version.Version)
	})

	t.Run("no update available", func(t *testing.T) {
		version, err := service.CheckForUpdate(2)
		require.NoError(t, err)
		assert.Nil(t, version)
	})
}

func TestApkVersionService_DeleteVersion(t *testing.T) {
	db := setupAPKServiceTestDB(t)
	defer db.Close()

	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "apk-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := &repository.ApkVersionRepository{}
	service := NewApkVersionService(repo)
	service.uploadDir = tmpDir

	// Set the repository's db
	repoStruct := (*struct {
		db *sql.DB
	})(unsafe.Pointer(repo))
	repoStruct.db = db

	t.Run("successful delete with file", func(t *testing.T) {
		// Create a test file
		testFilePath := filepath.Join(tmpDir, "test.apk")
		err := os.WriteFile(testFilePath, []byte("test content"), 0644)
		require.NoError(t, err)

		// Insert test data
		result, err := db.Exec(`INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, "1.0.0", 1, "test.apk", 1000, testFilePath, "notes", 1, time.Now())
		require.NoError(t, err)

		id, err := result.LastInsertId()
		require.NoError(t, err)

		err = service.DeleteVersion(int(id))
		assert.NoError(t, err)

		// Verify file was deleted
		_, err = os.Stat(testFilePath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("delete without physical file", func(t *testing.T) {
		// Insert test data with non-existent file path
		result, err := db.Exec(`INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, "2.0.0", 2, "test.apk", 1000, "/nonexistent/file.apk", "notes", 1, time.Now())
		require.NoError(t, err)

		id, err := result.LastInsertId()
		require.NoError(t, err)

		err = service.DeleteVersion(int(id))
		assert.NoError(t, err) // Should succeed even if file doesn't exist
	})
}

func TestApkVersionService_GetApkFilePath(t *testing.T) {
	db := setupAPKServiceTestDB(t)
	defer db.Close()

	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "apk-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := &repository.ApkVersionRepository{}
	service := NewApkVersionService(repo)

	// Set the repository's db
	repoStruct := (*struct {
		db *sql.DB
	})(unsafe.Pointer(repo))
	repoStruct.db = db

	t.Run("get existing file path", func(t *testing.T) {
		// Create a test file
		testFilePath := filepath.Join(tmpDir, "test.apk")
		err := os.WriteFile(testFilePath, []byte("test content"), 0644)
		require.NoError(t, err)

		// Insert test data
		result, err := db.Exec(`INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, "1.0.0", 1, "test.apk", 1000, testFilePath, "notes", 1, time.Now())
		require.NoError(t, err)

		id, err := result.LastInsertId()
		require.NoError(t, err)

		path, err := service.GetApkFilePath(int(id))
		require.NoError(t, err)
		assert.Equal(t, testFilePath, path)
	})

	t.Run("version not found", func(t *testing.T) {
		path, err := service.GetApkFilePath(999)
		assert.Error(t, err)
		assert.Empty(t, path)
		assert.Contains(t, err.Error(), "not found")
	})
}
