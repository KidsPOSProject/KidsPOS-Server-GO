package handlers

import "github.com/gin-gonic/gin"

// SetupRoutes はGinのルーターを設定します
// cmd/server/main.goとapi/index.goの両方から呼び出されます
func SetupRoutes(router *gin.Engine, h *Handlers) {
	// Web routes
	router.GET("/", h.Home)
	router.GET("/items", h.ItemsList)
	router.GET("/items/new", h.ItemsNew)
	router.POST("/items", h.ItemsCreate)
	router.GET("/items/:id/edit", h.ItemsEdit)
	router.POST("/items/:id", h.ItemsUpdate)
	router.POST("/items/:id/delete", h.ItemsDelete)

	router.GET("/sales", h.SalesList)
	router.GET("/sales/new", h.SalesNew)
	router.POST("/sales", h.SalesCreate)

	router.GET("/stores", h.StoresList)
	router.GET("/stores/new", h.StoresNew)
	router.POST("/stores", h.StoresCreate)
	router.GET("/stores/:id/edit", h.StoresEdit)
	router.POST("/stores/:id", h.StoresUpdate)
	router.POST("/stores/:id/delete", h.StoresDelete)

	router.GET("/staffs", h.StaffsList)
	router.GET("/staffs/new", h.StaffsNew)
	router.POST("/staffs", h.StaffsCreate)
	router.GET("/staffs/:id/edit", h.StaffsEdit)
	router.POST("/staffs/:id", h.StaffsUpdate)
	router.POST("/staffs/:id/delete", h.StaffsDelete)

	router.GET("/settings", h.SettingsList)
	router.GET("/reports/sales", h.ReportsSales)

	router.GET("/apk", h.ApkList)
	router.GET("/apk/upload", h.ApkUploadPage)
	router.POST("/apk/upload", h.ApkUpload)

	// API routes
	api := router.Group("/api")
	{
		api.GET("/items", h.APIItemsList)
		api.GET("/items/:id", h.APIItemsGet)
		api.POST("/items", h.APIItemsCreate)
		api.PUT("/items/:id", h.APIItemsUpdate)
		api.DELETE("/items/:id", h.APIItemsDelete)

		api.GET("/sales", h.APISalesList)
		api.GET("/sales/:id", h.APISalesGet)
		api.POST("/sales", h.APISalesCreate)

		api.GET("/stores", h.APIStoresList)
		api.GET("/stores/:id", h.APIStoresGet)
		api.POST("/stores", h.APIStoresCreate)
		api.PUT("/stores/:id", h.APIStoresUpdate)
		api.DELETE("/stores/:id", h.APIStoresDelete)

		api.GET("/staffs", h.APIStaffsList)
		api.GET("/staffs/:id", h.APIStaffsGet)
		api.POST("/staffs", h.APIStaffsCreate)
		api.PUT("/staffs/:id", h.APIStaffsUpdate)
		api.DELETE("/staffs/:id", h.APIStaffsDelete)

		api.GET("/settings", h.APISettingsList)
		api.PUT("/settings/:key", h.APISettingsUpdate)

		api.GET("/reports/sales", h.APIReportsSales)
		api.GET("/reports/sales/excel", h.APIReportsSalesExcel)

		api.GET("/apk/version/latest", h.APIApkLatest)
		api.GET("/apk/version/check", h.APIApkCheckUpdate)
		api.GET("/apk/version/all", h.APIApkVersions)
		api.GET("/apk/download/:id", h.APIApkDownload)
		api.GET("/apk/download/latest", h.APIApkDownloadLatest)
		api.POST("/apk/upload", h.APIApkUpload)
		api.DELETE("/apk/version/:id", h.APIApkDelete)
		api.PUT("/apk/version/:id/deactivate", h.APIApkDeactivate)
	}
}
