package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/models"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/repository"
	"github.com/google/uuid"
)

// ApkVersionService handles APK version business logic
type ApkVersionService struct {
	repo        *repository.ApkVersionRepository
	uploadDir   string
	maxFileSize int64
}

// NewApkVersionService creates a new APK version service
func NewApkVersionService(repo *repository.ApkVersionRepository) *ApkVersionService {
	uploadDir := "./uploads/apk"
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("Warning: failed to create upload directory: %v\n", err)
	}

	return &ApkVersionService{
		repo:        repo,
		uploadDir:   uploadDir,
		maxFileSize: 100 * 1024 * 1024, // 100MB
	}
}

// GetLatestVersion returns the latest active APK version
func (s *ApkVersionService) GetLatestVersion() (*models.ApkVersion, error) {
	return s.repo.FindLatest()
}

// CheckForUpdate checks if there's a newer version available
func (s *ApkVersionService) CheckForUpdate(currentVersionCode int) (*models.ApkVersion, error) {
	return s.repo.FindByVersionCode(currentVersionCode)
}

// GetAllVersions returns all active APK versions
func (s *ApkVersionService) GetAllVersions() ([]*models.ApkVersion, error) {
	return s.repo.FindAll()
}

// GetVersion returns an APK version by ID
func (s *ApkVersionService) GetVersion(id int) (*models.ApkVersion, error) {
	return s.repo.FindByID(id)
}

// UploadApk uploads a new APK file
func (s *ApkVersionService) UploadApk(file *multipart.FileHeader, version string, versionCode int, releaseNotes string) (*models.ApkVersion, error) {
	// Validate inputs
	if version == "" {
		return nil, fmt.Errorf("version is required")
	}
	if versionCode <= 0 {
		return nil, fmt.Errorf("version code must be positive")
	}

	// Validate file
	if file == nil {
		return nil, fmt.Errorf("file is required")
	}
	if file.Size > s.maxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum of %d bytes", s.maxFileSize)
	}

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".apk") {
		return nil, fmt.Errorf("file must be an APK file")
	}

	// Generate unique filename
	uniqueID := uuid.New().String()[:8]
	fileName := fmt.Sprintf("%s-%s.apk", version, uniqueID)
	filePath := filepath.Join(s.uploadDir, fileName)

	// Save file
	if err := s.saveFile(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create APK version record
	apk := &models.ApkVersion{
		Version:      version,
		VersionCode:  versionCode,
		FileName:     fileName,
		FileSize:     file.Size,
		FilePath:     filePath,
		ReleaseNotes: releaseNotes,
		IsActive:     true,
		UploadedAt:   time.Now(),
	}

	if err := s.repo.Create(apk); err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to create APK version: %w", err)
	}

	return apk, nil
}

// GetApkFilePath returns the file path for an APK version
func (s *ApkVersionService) GetApkFilePath(id int) (string, error) {
	apk, err := s.repo.FindByID(id)
	if err != nil {
		return "", err
	}

	// Check if file exists
	if _, err := os.Stat(apk.FilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("APK file not found")
	}

	return apk.FilePath, nil
}

// DeleteVersion deletes an APK version and its file
func (s *ApkVersionService) DeleteVersion(id int) error {
	apk, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete file
	if err := os.Remove(apk.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Delete database record
	return s.repo.Delete(id)
}

// DeactivateVersion deactivates an APK version
func (s *ApkVersionService) DeactivateVersion(id int) (*models.ApkVersion, error) {
	return s.repo.Deactivate(id)
}

// saveFile saves an uploaded file to the specified path
func (s *ApkVersionService) saveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
