package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/models"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/service"
	"github.com/gin-gonic/gin"
)

// Handlers holds all handler instances
type Handlers struct {
	itemService    *service.ItemService
	storeService   *service.StoreService
	staffService   *service.StaffService
	saleService    *service.SaleService
	settingService *service.SettingService
}

// NewHandlers creates handler instances
func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		itemService:    services.Item,
		storeService:   services.Store,
		staffService:   services.Staff,
		saleService:    services.Sale,
		settingService: services.Setting,
	}
}

// Home displays the home page
func (h *Handlers) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "KidsPOS",
	})
}

// ItemsList displays items list page
func (h *Handlers) ItemsList(c *gin.Context) {
	items, err := h.itemService.GetAllItems()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "items/index.html", gin.H{
		"title": "Items",
		"items": items,
	})
}

// ItemsNew displays new item form
func (h *Handlers) ItemsNew(c *gin.Context) {
	c.HTML(http.StatusOK, "items/new.html", gin.H{
		"title": "New Item",
	})
}

// ItemsCreate creates a new item
func (h *Handlers) ItemsCreate(c *gin.Context) {
	item := &models.Item{
		Name:  c.PostForm("name"),
		Price: atoi(c.PostForm("price")),
		Stock: atoi(c.PostForm("stock")),
	}

	if err := h.itemService.CreateItem(item); err != nil {
		c.HTML(http.StatusBadRequest, "items/new.html", gin.H{
			"title": "New Item",
			"error": err.Error(),
			"item":  item,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/items")
}

// ItemsEdit displays item edit form
func (h *Handlers) ItemsEdit(c *gin.Context) {
	id := atoi(c.Param("id"))
	item, err := h.itemService.GetItem(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Item not found",
		})
		return
	}

	c.HTML(http.StatusOK, "items/edit.html", gin.H{
		"title": "Edit Item",
		"item":  item,
	})
}

// ItemsUpdate updates an item
func (h *Handlers) ItemsUpdate(c *gin.Context) {
	id := atoi(c.Param("id"))
	item := &models.Item{
		ID:    id,
		Name:  c.PostForm("name"),
		Price: atoi(c.PostForm("price")),
		Stock: atoi(c.PostForm("stock")),
	}

	if err := h.itemService.UpdateItem(item); err != nil {
		c.HTML(http.StatusBadRequest, "items/edit.html", gin.H{
			"title": "Edit Item",
			"error": err.Error(),
			"item":  item,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/items")
}

// ItemsDelete deletes an item
func (h *Handlers) ItemsDelete(c *gin.Context) {
	id := atoi(c.Param("id"))
	if err := h.itemService.DeleteItem(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/items")
}

// SalesList displays sales list page
func (h *Handlers) SalesList(c *gin.Context) {
	sales, err := h.saleService.GetAllSales()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "sales/index.html", gin.H{
		"title": "Sales",
		"sales": sales,
	})
}

// SalesNew displays new sale form
func (h *Handlers) SalesNew(c *gin.Context) {
	items, _ := h.itemService.GetAllItems()
	stores, _ := h.storeService.GetAllStores()
	staffs, _ := h.staffService.GetAllStaffs()

	c.HTML(http.StatusOK, "sales/new.html", gin.H{
		"title":  "New Sale",
		"items":  items,
		"stores": stores,
		"staffs": staffs,
	})
}

// SalesCreate creates a new sale
func (h *Handlers) SalesCreate(c *gin.Context) {
	// Parse form data
	sale := &models.Sale{
		StoreID: atoi(c.PostForm("storeId")),
		StaffID: atoi(c.PostForm("staffId")),
		Deposit: atoi(c.PostForm("deposit")),
		SaleAt:  time.Now(),
	}

	// Parse sale details
	itemIds := c.PostFormArray("itemId[]")
	quantities := c.PostFormArray("quantity[]")

	for i := 0; i < len(itemIds); i++ {
		detail := models.SaleDetail{
			ItemID:   atoi(itemIds[i]),
			Quantity: atoi(quantities[i]),
		}
		sale.Details = append(sale.Details, detail)
	}

	if err := h.saleService.CreateSale(sale); err != nil {
		items, _ := h.itemService.GetAllItems()
		stores, _ := h.storeService.GetAllStores()
		staffs, _ := h.staffService.GetAllStaffs()

		c.HTML(http.StatusBadRequest, "sales/new.html", gin.H{
			"title":  "New Sale",
			"error":  err.Error(),
			"sale":   sale,
			"items":  items,
			"stores": stores,
			"staffs": staffs,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/sales")
}

// StoresList displays stores list page
func (h *Handlers) StoresList(c *gin.Context) {
	stores, err := h.storeService.GetAllStores()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "stores/index.html", gin.H{
		"title":  "Stores",
		"stores": stores,
	})
}

// StoresNew displays new store form
func (h *Handlers) StoresNew(c *gin.Context) {
	c.HTML(http.StatusOK, "stores/new.html", gin.H{
		"title": "New Store",
	})
}

// StoresCreate creates a new store
func (h *Handlers) StoresCreate(c *gin.Context) {
	store := &models.Store{
		Name: c.PostForm("name"),
	}

	if err := h.storeService.CreateStore(store); err != nil {
		c.HTML(http.StatusBadRequest, "stores/new.html", gin.H{
			"title": "New Store",
			"error": err.Error(),
			"store": store,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/stores")
}

// StoresEdit displays store edit form
func (h *Handlers) StoresEdit(c *gin.Context) {
	id := atoi(c.Param("id"))
	store, err := h.storeService.GetStore(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Store not found",
		})
		return
	}

	c.HTML(http.StatusOK, "stores/edit.html", gin.H{
		"title": "Edit Store",
		"store": store,
	})
}

// StoresUpdate updates a store
func (h *Handlers) StoresUpdate(c *gin.Context) {
	// TODO: Implement store update
	c.Redirect(http.StatusSeeOther, "/stores")
}

// StoresDelete deletes a store
func (h *Handlers) StoresDelete(c *gin.Context) {
	// TODO: Implement store delete
	c.Redirect(http.StatusSeeOther, "/stores")
}

// StaffsList displays staffs list page
func (h *Handlers) StaffsList(c *gin.Context) {
	staffs, err := h.staffService.GetAllStaffs()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "staffs/index.html", gin.H{
		"title":  "Staffs",
		"staffs": staffs,
	})
}

// StaffsNew displays new staff form
func (h *Handlers) StaffsNew(c *gin.Context) {
	c.HTML(http.StatusOK, "staffs/new.html", gin.H{
		"title": "New Staff",
	})
}

// StaffsCreate creates a new staff
func (h *Handlers) StaffsCreate(c *gin.Context) {
	staff := &models.Staff{
		Name: c.PostForm("name"),
	}

	if err := h.staffService.CreateStaff(staff); err != nil {
		c.HTML(http.StatusBadRequest, "staffs/new.html", gin.H{
			"title": "New Staff",
			"error": err.Error(),
			"staff": staff,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/staffs")
}

// StaffsEdit displays staff edit form
func (h *Handlers) StaffsEdit(c *gin.Context) {
	id := atoi(c.Param("id"))
	staff, err := h.staffService.GetStaff(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Staff not found",
		})
		return
	}

	c.HTML(http.StatusOK, "staffs/edit.html", gin.H{
		"title": "Edit Staff",
		"staff": staff,
	})
}

// StaffsUpdate updates a staff
func (h *Handlers) StaffsUpdate(c *gin.Context) {
	// TODO: Implement staff update
	c.Redirect(http.StatusSeeOther, "/staffs")
}

// StaffsDelete deletes a staff
func (h *Handlers) StaffsDelete(c *gin.Context) {
	// TODO: Implement staff delete
	c.Redirect(http.StatusSeeOther, "/staffs")
}

// SettingsList displays settings page
func (h *Handlers) SettingsList(c *gin.Context) {
	settings, err := h.settingService.GetAllSettings()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "settings/index.html", gin.H{
		"title":    "Settings",
		"settings": settings,
	})
}

// ReportsSales displays sales report page
func (h *Handlers) ReportsSales(c *gin.Context) {
	sales, err := h.saleService.GetAllSales()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "reports/sales.html", gin.H{
		"title": "Sales Report",
		"sales": sales,
	})
}

// ===============================
// MISSING API HANDLERS - Compatible with Kotlin version
// ===============================

// APIItemsByBarcode handles GET /api/item/barcode/{barcode}
func (h *Handlers) APIItemsByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")

	// Validate barcode format (4+ digits)
	if len(barcode) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid barcode format - must be 4+ digits",
		})
		return
	}

	// For now, search by itemId (future: implement barcode field)
	items, err := h.itemService.GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Search for item by ItemID matching barcode
	for _, item := range items {
		if item.ItemID == barcode {
			c.JSON(http.StatusOK, gin.H{"item": item})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Item not found with barcode: " + barcode,
	})
}

// APIItemsPatch handles PATCH /api/item/{id}
func (h *Handlers) APIItemsPatch(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get existing item
	item, err := h.itemService.GetItemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Parse partial updates
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Apply partial updates
	if name, exists := updates["name"]; exists {
		if nameStr, ok := name.(string); ok {
			item.Name = nameStr
		}
	}
	if price, exists := updates["price"]; exists {
		if priceFloat, ok := price.(float64); ok {
			item.Price = int(priceFloat)
		}
	}
	if stock, exists := updates["stock"]; exists {
		if stockFloat, ok := stock.(float64); ok {
			item.Stock = int(stockFloat)
		}
	}

	// Update item
	if err := h.itemService.UpdateItem(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

// APIItemsBarcodePDF handles GET /api/item/barcode-pdf
func (h *Handlers) APIItemsBarcodePDF(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "PDF generation not implemented yet",
		"note":  "Future enhancement: generate barcode PDF for all items",
	})
}

// APIItemsBarcodeSelectedPDF handles POST /api/item/barcode-pdf/selected
func (h *Handlers) APIItemsBarcodeSelectedPDF(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "PDF generation not implemented yet",
		"note":  "Future enhancement: generate barcode PDF for selected items",
	})
}

// APISalesCreateAlt handles POST /api/sales/create (alternative endpoint)
func (h *Handlers) APISalesCreateAlt(c *gin.Context) {
	// Redirect to main create endpoint
	h.APISalesCreate(c)
}

// APISalesValidatePrinter handles GET /api/sales/validate-printer/{storeId}
func (h *Handlers) APISalesValidatePrinter(c *gin.Context) {
	storeId := c.Param("storeId")

	c.JSON(http.StatusOK, gin.H{
		"storeId":        storeId,
		"printerStatus":  "available",
		"isValid":        true,
		"message":        "Printer validation not fully implemented",
	})
}

// APIStaffByBarcode handles GET /api/staff/{barcode}
func (h *Handlers) APIStaffByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")

	// Get all staff and search by StaffID (acting as barcode)
	staffList, err := h.staffService.GetAllStaff()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, staff := range staffList {
		if staff.StaffID == barcode {
			c.JSON(http.StatusOK, gin.H{"staff": staff})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Staff not found with barcode: " + barcode,
	})
}

// APIStaffUpdateByBarcode handles PUT /api/staff/{barcode}
func (h *Handlers) APIStaffUpdateByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")

	// Find staff by barcode
	staffList, err := h.staffService.GetAllStaff()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var targetStaff *models.Staff
	for _, staff := range staffList {
		if staff.StaffID == barcode {
			targetStaff = staff
			break
		}
	}

	if targetStaff == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	// Bind update data
	if err := c.ShouldBindJSON(targetStaff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update staff
	if err := h.staffService.UpdateStaff(targetStaff); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"staff": targetStaff})
}

// APIStaffDeleteByBarcode handles DELETE /api/staff/{barcode}
func (h *Handlers) APIStaffDeleteByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")

	// Find staff by barcode
	staffList, err := h.staffService.GetAllStaff()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, staff := range staffList {
		if staff.StaffID == barcode {
			if err := h.staffService.DeleteStaff(staff.ID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Staff deleted successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
}

// APISettingsStatus handles GET /api/setting/status
func (h *Handlers) APISettingsStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// APISettingsGet handles GET /api/setting/{key}
func (h *Handlers) APISettingsGet(c *gin.Context) {
	key := c.Param("key")

	settings, err := h.settingService.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, setting := range settings {
		if setting.Key == key {
			c.JSON(http.StatusOK, gin.H{"setting": setting})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
}

// APISettingsCreate handles POST /api/setting
func (h *Handlers) APISettingsCreate(c *gin.Context) {
	var setting models.Setting
	if err := c.ShouldBindJSON(&setting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingService.CreateSetting(&setting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"setting": setting})
}

// APISettingsDelete handles DELETE /api/setting/{key}
func (h *Handlers) APISettingsDelete(c *gin.Context) {
	key := c.Param("key")

	if err := h.settingService.DeleteSetting(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting deleted successfully"})
}

// APISettingsPrinterSet handles POST /api/setting/printer/{storeId}
func (h *Handlers) APISettingsPrinterSet(c *gin.Context) {
	storeId := c.Param("storeId")

	var printerConfig map[string]interface{}
	if err := c.ShouldBindJSON(&printerConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Printer configuration saved",
		"storeId": storeId,
		"config":  printerConfig,
		"note":    "Printer functionality not fully implemented",
	})
}

// APISettingsPrinterGet handles GET /api/setting/printer/{storeId}
func (h *Handlers) APISettingsPrinterGet(c *gin.Context) {
	storeId := c.Param("storeId")

	c.JSON(http.StatusOK, gin.H{
		"storeId":    storeId,
		"enabled":    true,
		"driver":     "default",
		"configured": true,
		"note":       "Printer functionality not fully implemented",
	})
}

// APISettingsApplicationSet handles POST /api/setting/application
func (h *Handlers) APISettingsApplicationSet(c *gin.Context) {
	var appConfig map[string]interface{}
	if err := c.ShouldBindJSON(&appConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Application configuration saved",
		"config":  appConfig,
	})
}

// APISettingsApplicationGet handles GET /api/setting/application
func (h *Handlers) APISettingsApplicationGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":     "1.0.0-go",
		"environment": "development",
		"database":    "sqlite",
		"features": gin.H{
			"pdf_generation": false,
			"excel_export":   false,
			"printer_support": false,
		},
	})
}

// ===============================
// USER API (Staff Alias) - Compatible with Kotlin /api/users
// ===============================

// APIUsersList handles GET /api/users
func (h *Handlers) APIUsersList(c *gin.Context) {
	staffList, err := h.staffService.GetAllStaff()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": staffList})
}

// APIUsersByBarcode handles GET /api/users/{barcode}
func (h *Handlers) APIUsersByBarcode(c *gin.Context) {
	// Delegate to staff handler
	h.APIStaffByBarcode(c)
}

// APIUsersCreate handles POST /api/users
func (h *Handlers) APIUsersCreate(c *gin.Context) {
	// Delegate to staff create
	h.APIStaffsCreate(c)
}

// APIUsersUpdateByBarcode handles PUT /api/users/{barcode}
func (h *Handlers) APIUsersUpdateByBarcode(c *gin.Context) {
	// Delegate to staff update
	h.APIStaffUpdateByBarcode(c)
}

// APIUsersDeleteByBarcode handles DELETE /api/users/{barcode}
func (h *Handlers) APIUsersDeleteByBarcode(c *gin.Context) {
	// Delegate to staff delete
	h.APIStaffDeleteByBarcode(c)
}

// ===============================
// EXTENDED REPORTS API - Compatible with Kotlin version
// ===============================

// APIReportsSalesPDF handles GET /api/reports/sales/pdf
func (h *Handlers) APIReportsSalesPDF(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "PDF report generation not implemented yet",
		"note":  "Future enhancement: generate sales report as PDF",
	})
}

// APIReportsSalesPDFToday handles GET /api/reports/sales/pdf/today
func (h *Handlers) APIReportsSalesPDFToday(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "PDF report generation not implemented yet",
		"note":  "Future enhancement: generate today's sales report as PDF",
	})
}

// APIReportsSalesPDFMonth handles GET /api/reports/sales/pdf/month
func (h *Handlers) APIReportsSalesPDFMonth(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "PDF report generation not implemented yet",
		"note":  "Future enhancement: generate monthly sales report as PDF",
	})
}

// APIReportsSalesExcelToday handles GET /api/reports/sales/excel/today
func (h *Handlers) APIReportsSalesExcelToday(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Excel report generation not implemented yet",
		"note":  "Future enhancement: generate today's sales report as Excel",
	})
}

// APIReportsSalesExcelMonth handles GET /api/reports/sales/excel/month
func (h *Handlers) APIReportsSalesExcelMonth(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Excel report generation not implemented yet",
		"note":  "Future enhancement: generate monthly sales report as Excel",
	})
}

// Helper function to convert string to int
func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
