package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ApkList displays APK versions list page
func (h *Handlers) ApkList(c *gin.Context) {
	versions, err := h.apkVersionService.GetAllVersions()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "apk/index.html", gin.H{
		"title":    "APK Versions",
		"versions": versions,
	})
}

// ApkUploadPage displays APK upload form
func (h *Handlers) ApkUploadPage(c *gin.Context) {
	c.HTML(http.StatusOK, "apk/upload.html", gin.H{
		"title": "Upload APK",
	})
}

// ApkUpload handles APK file upload
func (h *Handlers) ApkUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusBadRequest, "apk/upload.html", gin.H{
			"title": "Upload APK",
			"error": "File is required",
		})
		return
	}

	version := c.PostForm("version")
	versionCodeStr := c.PostForm("versionCode")
	releaseNotes := c.DefaultPostForm("releaseNotes", "")

	if version == "" {
		c.HTML(http.StatusBadRequest, "apk/upload.html", gin.H{
			"title": "Upload APK",
			"error": "Version is required",
		})
		return
	}
	if versionCodeStr == "" {
		c.HTML(http.StatusBadRequest, "apk/upload.html", gin.H{
			"title": "Upload APK",
			"error": "Version code is required",
		})
		return
	}

	versionCode, err := strconv.Atoi(versionCodeStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "apk/upload.html", gin.H{
			"title": "Upload APK",
			"error": "Invalid version code",
		})
		return
	}

	// Upload APK
	_, err = h.apkVersionService.UploadApk(file, version, versionCode, releaseNotes)
	if err != nil {
		c.HTML(http.StatusBadRequest, "apk/upload.html", gin.H{
			"title": "Upload APK",
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/apk")
}
