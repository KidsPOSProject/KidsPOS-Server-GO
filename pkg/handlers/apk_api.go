package handlers

import (
	"net/http"
	"strconv"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/models"
	"github.com/gin-gonic/gin"
)

// APIApkLatest returns the latest APK version
func (h *Handlers) APIApkLatest(c *gin.Context) {
	version, err := h.apkVersionService.GetLatestVersion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if version == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No APK versions available"})
		return
	}

	c.JSON(http.StatusOK, version)
}

// APIApkCheckUpdate checks for updates
func (h *Handlers) APIApkCheckUpdate(c *gin.Context) {
	versionCodeStr := c.Query("currentVersionCode")
	if versionCodeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currentVersionCode is required"})
		return
	}

	currentVersionCode, err := strconv.Atoi(versionCodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version code"})
		return
	}

	newerVersion, err := h.apkVersionService.CheckForUpdate(currentVersionCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if newerVersion == nil {
		c.JSON(http.StatusOK, gin.H{
			"hasUpdate":     false,
			"latestVersion": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hasUpdate":     true,
		"latestVersion": newerVersion,
	})
}

// APIApkVersions returns all APK versions
func (h *Handlers) APIApkVersions(c *gin.Context) {
	versions, err := h.apkVersionService.GetAllVersions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if versions == nil {
		versions = make([]*models.ApkVersion, 0)
	}

	c.JSON(http.StatusOK, versions)
}

// APIApkDownload downloads an APK file by ID
func (h *Handlers) APIApkDownload(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get APK version
	version, err := h.apkVersionService.GetVersion(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "APK version not found"})
		return
	}

	// Get file path
	filePath, err := h.apkVersionService.GetApkFilePath(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Serve file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+version.FileName)
	c.Header("Content-Type", "application/vnd.android.package-archive")
	c.File(filePath)
}

// APIApkDownloadLatest downloads the latest APK file
func (h *Handlers) APIApkDownloadLatest(c *gin.Context) {
	version, err := h.apkVersionService.GetLatestVersion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if version == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No APK versions available"})
		return
	}

	// Get file path
	filePath, err := h.apkVersionService.GetApkFilePath(version.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Serve file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+version.FileName)
	c.Header("Content-Type", "application/vnd.android.package-archive")
	c.File(filePath)
}

// APIApkUpload uploads a new APK file
func (h *Handlers) APIApkUpload(c *gin.Context) {
	// Parse multipart form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	version := c.PostForm("version")
	versionCodeStr := c.PostForm("versionCode")
	releaseNotes := c.DefaultPostForm("releaseNotes", "")

	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Version is required"})
		return
	}
	if versionCodeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Version code is required"})
		return
	}

	versionCode, err := strconv.Atoi(versionCodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version code"})
		return
	}

	// Upload APK
	apk, err := h.apkVersionService.UploadApk(file, version, versionCode, releaseNotes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, apk)
}

// APIApkDelete deletes an APK version
func (h *Handlers) APIApkDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.apkVersionService.DeleteVersion(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "APK version deleted successfully"})
}

// APIApkDeactivate deactivates an APK version
func (h *Handlers) APIApkDeactivate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	apk, err := h.apkVersionService.DeactivateVersion(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, apk)
}
