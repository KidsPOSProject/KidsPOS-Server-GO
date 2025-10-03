package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/models"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/pkg/service"
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

// Helper function to convert string to int
func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}