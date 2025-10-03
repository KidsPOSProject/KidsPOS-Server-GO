package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/models"
)

// ApkVersionRepository handles APK version data access
type ApkVersionRepository struct {
	db *sql.DB
}

// FindLatest returns the latest active APK version
func (r *ApkVersionRepository) FindLatest() (*models.ApkVersion, error) {
	query := `SELECT id, version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt, createdAt, updatedAt
			  FROM apk_versions WHERE isActive = 1 ORDER BY versionCode DESC LIMIT 1`

	apk := &models.ApkVersion{}
	err := r.db.QueryRow(query).Scan(&apk.ID, &apk.Version, &apk.VersionCode, &apk.FileName,
		&apk.FileSize, &apk.FilePath, &apk.ReleaseNotes, &apk.IsActive, &apk.UploadedAt, &apk.CreatedAt, &apk.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return apk, nil
}

// FindByID finds an APK version by ID
func (r *ApkVersionRepository) FindByID(id int) (*models.ApkVersion, error) {
	query := `SELECT id, version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt, createdAt, updatedAt
			  FROM apk_versions WHERE id = ?`

	apk := &models.ApkVersion{}
	err := r.db.QueryRow(query, id).Scan(&apk.ID, &apk.Version, &apk.VersionCode, &apk.FileName,
		&apk.FileSize, &apk.FilePath, &apk.ReleaseNotes, &apk.IsActive, &apk.UploadedAt, &apk.CreatedAt, &apk.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("apk version not found")
	}
	if err != nil {
		return nil, err
	}
	return apk, nil
}

// FindAll returns all active APK versions
func (r *ApkVersionRepository) FindAll() ([]*models.ApkVersion, error) {
	query := `SELECT id, version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt, createdAt, updatedAt
			  FROM apk_versions WHERE isActive = 1 ORDER BY versionCode DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*models.ApkVersion
	for rows.Next() {
		apk := &models.ApkVersion{}
		err := rows.Scan(&apk.ID, &apk.Version, &apk.VersionCode, &apk.FileName,
			&apk.FileSize, &apk.FilePath, &apk.ReleaseNotes, &apk.IsActive, &apk.UploadedAt, &apk.CreatedAt, &apk.UpdatedAt)
		if err != nil {
			return nil, err
		}
		versions = append(versions, apk)
	}
	return versions, nil
}

// FindByVersionCode finds APK versions with version code greater than the given code
func (r *ApkVersionRepository) FindByVersionCode(code int) (*models.ApkVersion, error) {
	query := `SELECT id, version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt, createdAt, updatedAt
			  FROM apk_versions WHERE isActive = 1 AND versionCode > ? ORDER BY versionCode DESC LIMIT 1`

	apk := &models.ApkVersion{}
	err := r.db.QueryRow(query, code).Scan(&apk.ID, &apk.Version, &apk.VersionCode, &apk.FileName,
		&apk.FileSize, &apk.FilePath, &apk.ReleaseNotes, &apk.IsActive, &apk.UploadedAt, &apk.CreatedAt, &apk.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return apk, nil
}

// Create creates a new APK version
func (r *ApkVersionRepository) Create(apk *models.ApkVersion) error {
	query := `INSERT INTO apk_versions (version, versionCode, fileName, fileSize, filePath, releaseNotes, isActive, uploadedAt, createdAt, updatedAt)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, apk.Version, apk.VersionCode, apk.FileName, apk.FileSize,
		apk.FilePath, apk.ReleaseNotes, apk.IsActive, apk.UploadedAt, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	apk.ID = int(id)
	apk.CreatedAt = now
	apk.UpdatedAt = now

	return nil
}

// Update updates an APK version
func (r *ApkVersionRepository) Update(apk *models.ApkVersion) error {
	query := `UPDATE apk_versions SET version = ?, versionCode = ?, fileName = ?, fileSize = ?,
			  filePath = ?, releaseNotes = ?, isActive = ?, updatedAt = ? WHERE id = ?`

	now := time.Now()
	result, err := r.db.Exec(query, apk.Version, apk.VersionCode, apk.FileName, apk.FileSize,
		apk.FilePath, apk.ReleaseNotes, apk.IsActive, now, apk.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("apk version not found")
	}

	apk.UpdatedAt = now
	return nil
}

// Delete deletes an APK version
func (r *ApkVersionRepository) Delete(id int) error {
	query := `DELETE FROM apk_versions WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("apk version not found")
	}

	return nil
}

// Deactivate deactivates an APK version
func (r *ApkVersionRepository) Deactivate(id int) (*models.ApkVersion, error) {
	query := `UPDATE apk_versions SET isActive = 0, updatedAt = ? WHERE id = ?`

	now := time.Now()
	result, err := r.db.Exec(query, now, id)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, fmt.Errorf("apk version not found")
	}

	// Return the updated version
	return r.FindByID(id)
}
