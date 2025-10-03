package models

import "time"

// ApkVersion represents an APK version
type ApkVersion struct {
	ID           int       `json:"id" db:"id"`
	Version      string    `json:"version" db:"version"`
	VersionCode  int       `json:"versionCode" db:"versionCode"`
	FileName     string    `json:"fileName" db:"fileName"`
	FileSize     int64     `json:"fileSize" db:"fileSize"`
	FilePath     string    `json:"filePath" db:"filePath"`
	ReleaseNotes string    `json:"releaseNotes" db:"releaseNotes"`
	IsActive     bool      `json:"isActive" db:"isActive"`
	UploadedAt   time.Time `json:"uploadedAt" db:"uploadedAt"`
	CreatedAt    time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updatedAt"`
}
