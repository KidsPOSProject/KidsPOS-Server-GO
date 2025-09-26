package handlers

import (
	"net/http"
	"strconv"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/models"
	"github.com/gin-gonic/gin"
)

// APIItemsList returns items list as JSON
func (h *Handlers) APIItemsList(c *gin.Context) {
	items, err := h.itemService.GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// APIItemsGet returns a single item as JSON
func (h *Handlers) APIItemsGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	item, err := h.itemService.GetItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// APIItemsCreate creates a new item via API
func (h *Handlers) APIItemsCreate(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.itemService.CreateItem(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// APIItemsUpdate updates an item via API
func (h *Handlers) APIItemsUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.ID = id
	if err := h.itemService.UpdateItem(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// APIItemsDelete deletes an item via API
func (h *Handlers) APIItemsDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.itemService.DeleteItem(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

// APISalesList returns sales list as JSON
func (h *Handlers) APISalesList(c *gin.Context) {
	sales, err := h.saleService.GetAllSales()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sales)
}

// APISalesGet returns a single sale as JSON
func (h *Handlers) APISalesGet(c *gin.Context) {
	// TODO: Implement get single sale
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// APISalesCreate creates a new sale via API
func (h *Handlers) APISalesCreate(c *gin.Context) {
	var sale models.Sale
	if err := c.ShouldBindJSON(&sale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.saleService.CreateSale(&sale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sale)
}

// APIStoresList returns stores list as JSON
func (h *Handlers) APIStoresList(c *gin.Context) {
	stores, err := h.storeService.GetAllStores()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stores)
}

// APIStoresGet returns a single store as JSON
func (h *Handlers) APIStoresGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	store, err := h.storeService.GetStore(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	c.JSON(http.StatusOK, store)
}

// APIStoresCreate creates a new store via API
func (h *Handlers) APIStoresCreate(c *gin.Context) {
	var store models.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storeService.CreateStore(&store); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, store)
}

// APIStoresUpdate updates a store via API
func (h *Handlers) APIStoresUpdate(c *gin.Context) {
	// TODO: Implement store update
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// APIStoresDelete deletes a store via API
func (h *Handlers) APIStoresDelete(c *gin.Context) {
	// TODO: Implement store delete
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// APIStaffsList returns staffs list as JSON
func (h *Handlers) APIStaffsList(c *gin.Context) {
	staffs, err := h.staffService.GetAllStaffs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staffs)
}

// APIStaffsGet returns a single staff as JSON
func (h *Handlers) APIStaffsGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	staff, err := h.staffService.GetStaff(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	c.JSON(http.StatusOK, staff)
}

// APIStaffsCreate creates a new staff via API
func (h *Handlers) APIStaffsCreate(c *gin.Context) {
	var staff models.Staff
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.staffService.CreateStaff(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, staff)
}

// APIStaffsUpdate updates a staff via API
func (h *Handlers) APIStaffsUpdate(c *gin.Context) {
	// TODO: Implement staff update
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// APIStaffsDelete deletes a staff via API
func (h *Handlers) APIStaffsDelete(c *gin.Context) {
	// TODO: Implement staff delete
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// APISettingsList returns settings list as JSON
func (h *Handlers) APISettingsList(c *gin.Context) {
	settings, err := h.settingService.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// APISettingsUpdate updates a setting via API
func (h *Handlers) APISettingsUpdate(c *gin.Context) {
	key := c.Param("key")

	var payload struct {
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.settingService.UpdateSetting(key, payload.Value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting updated successfully"})
}

// APIReportsSales returns sales report as JSON
func (h *Handlers) APIReportsSales(c *gin.Context) {
	sales, err := h.saleService.GetAllSales()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate summary
	totalSales := 0
	totalAmount := 0
	for _, sale := range sales {
		totalSales++
		totalAmount += sale.TotalPrice
	}

	c.JSON(http.StatusOK, gin.H{
		"sales":       sales,
		"totalSales":  totalSales,
		"totalAmount": totalAmount,
	})
}

// APIReportsSalesExcel generates Excel report
func (h *Handlers) APIReportsSalesExcel(c *gin.Context) {
	// TODO: Implement Excel generation
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}
