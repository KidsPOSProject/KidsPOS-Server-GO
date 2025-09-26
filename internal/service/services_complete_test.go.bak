package service

import (
	"testing"
	"time"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/models"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/repository"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemService_Complete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	itemRepo := repository.NewItemRepository(db)
	itemService := NewItemService(itemRepo)

	t.Run("CreateItem", func(t *testing.T) {
		item := &models.Item{
			ItemID: "TEST001",
			Name:   "Test Item",
			Price:  100,
			Stock:  10,
		}

		err := itemService.CreateItem(item)
		require.NoError(t, err)
		assert.NotZero(t, item.ID)
	})

	t.Run("GetItemByID", func(t *testing.T) {
		// Create item first
		item := &models.Item{
			ItemID: "GET001",
			Name:   "Get Test Item",
			Price:  200,
			Stock:  20,
		}
		err := itemService.CreateItem(item)
		require.NoError(t, err)

		// Get item
		retrieved, err := itemService.GetItemByID(item.ID)
		require.NoError(t, err)
		assert.Equal(t, item.Name, retrieved.Name)
		assert.Equal(t, item.Price, retrieved.Price)
	})

	t.Run("UpdateItem", func(t *testing.T) {
		// Create item first
		item := &models.Item{
			ItemID: "UPD001",
			Name:   "Update Test Item",
			Price:  300,
			Stock:  30,
		}
		err := itemService.CreateItem(item)
		require.NoError(t, err)

		// Update item
		item.Name = "Updated Item"
		item.Price = 400
		err = itemService.UpdateItem(item)
		require.NoError(t, err)

		// Verify update
		retrieved, err := itemService.GetItemByID(item.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Item", retrieved.Name)
		assert.Equal(t, 400, retrieved.Price)
	})

	t.Run("DeleteItem", func(t *testing.T) {
		// Create item first
		item := &models.Item{
			ItemID: "DEL001",
			Name:   "Delete Test Item",
			Price:  500,
			Stock:  50,
		}
		err := itemService.CreateItem(item)
		require.NoError(t, err)

		// Delete item
		err = itemService.DeleteItem(item.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = itemService.GetItemByID(item.ID)
		assert.Error(t, err)
	})

	t.Run("GetAllItems", func(t *testing.T) {
		// Create multiple items
		for i := 0; i < 3; i++ {
			item := &models.Item{
				ItemID: "LIST00" + string(rune('1'+i)),
				Name:   "List Item",
				Price:  100 * (i + 1),
				Stock:  10,
			}
			err := itemService.CreateItem(item)
			require.NoError(t, err)
		}

		// Get all items
		items, err := itemService.GetAllItems()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(items), 3)
	})
}

func TestStaffService_Complete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	staffRepo := repository.NewStaffRepository(db)
	staffService := NewStaffService(staffRepo)

	t.Run("CreateStaff", func(t *testing.T) {
		staff := &models.Staff{
			StaffID: "STAFF001",
			Name:    "Test Staff",
		}

		err := staffService.CreateStaff(staff)
		require.NoError(t, err)
		assert.NotZero(t, staff.ID)
	})

	t.Run("GetStaffByID", func(t *testing.T) {
		staff := &models.Staff{
			StaffID: "STAFF002",
			Name:    "Get Test Staff",
		}
		err := staffService.CreateStaff(staff)
		require.NoError(t, err)

		retrieved, err := staffService.GetStaffByID(staff.ID)
		require.NoError(t, err)
		assert.Equal(t, staff.Name, retrieved.Name)
	})

	t.Run("UpdateStaff", func(t *testing.T) {
		staff := &models.Staff{
			StaffID: "STAFF003",
			Name:    "Update Test Staff",
		}
		err := staffService.CreateStaff(staff)
		require.NoError(t, err)

		staff.Name = "Updated Staff"
		err = staffService.UpdateStaff(staff)
		require.NoError(t, err)

		retrieved, err := staffService.GetStaffByID(staff.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Staff", retrieved.Name)
	})

	t.Run("DeleteStaff", func(t *testing.T) {
		staff := &models.Staff{
			StaffID: "STAFF004",
			Name:    "Delete Test Staff",
		}
		err := staffService.CreateStaff(staff)
		require.NoError(t, err)

		err = staffService.DeleteStaff(staff.ID)
		require.NoError(t, err)

		_, err = staffService.GetStaffByID(staff.ID)
		assert.Error(t, err)
	})

	t.Run("GetAllStaff", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			staff := &models.Staff{
				StaffID: "STAFFLIST" + string(rune('1'+i)),
				Name:    "List Staff",
			}
			err := staffService.CreateStaff(staff)
			require.NoError(t, err)
		}

		staffList, err := staffService.GetAllStaff()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(staffList), 3)
	})
}

func TestSaleService_Complete(t *testing.T) {
	db := testutil.SetupTestDB(t)

	// Setup repositories
	saleRepo := repository.NewSaleRepository(db)
	itemRepo := repository.NewItemRepository(db)

	// Setup services
	saleService := NewSaleService(saleRepo)
	itemService := NewItemService(itemRepo)

	// Create test items
	item1 := &models.Item{
		ItemID: "SALE001",
		Name:   "Sale Item 1",
		Price:  100,
		Stock:  50,
	}
	err := itemService.CreateItem(item1)
	require.NoError(t, err)

	item2 := &models.Item{
		ItemID: "SALE002",
		Name:   "Sale Item 2",
		Price:  200,
		Stock:  30,
	}
	err = itemService.CreateItem(item2)
	require.NoError(t, err)

	t.Run("CreateSale", func(t *testing.T) {
		sale := &models.Sale{
			StoreID:    1,
			StaffID:    1,
			TotalPrice: 500,
			Deposit:    600,
			SaleAt:     time.Now(),
			Details: []models.SaleDetail{
				{ItemID: item1.ID, Price: 100, Quantity: 2},
				{ItemID: item2.ID, Price: 200, Quantity: 1},
			},
		}

		err := saleService.CreateSale(sale)
		require.NoError(t, err)
		assert.NotZero(t, sale.ID)

		// Verify stock reduction
		updatedItem1, err := itemService.GetItemByID(item1.ID)
		require.NoError(t, err)
		assert.Equal(t, 48, updatedItem1.Stock)

		updatedItem2, err := itemService.GetItemByID(item2.ID)
		require.NoError(t, err)
		assert.Equal(t, 29, updatedItem2.Stock)
	})

	t.Run("GetSaleByID", func(t *testing.T) {
		sale := &models.Sale{
			StoreID:    1,
			StaffID:    1,
			TotalPrice: 300,
			Deposit:    300,
			SaleAt:     time.Now(),
			Details: []models.SaleDetail{
				{ItemID: item1.ID, Price: 100, Quantity: 3},
			},
		}
		err := saleService.CreateSale(sale)
		require.NoError(t, err)

		retrieved, err := saleService.GetSaleByID(sale.ID)
		require.NoError(t, err)
		assert.Equal(t, sale.TotalPrice, retrieved.TotalPrice)
		assert.Len(t, retrieved.Details, 1)
	})

	t.Run("GetAllSales", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			sale := &models.Sale{
				StoreID:    1,
				StaffID:    1,
				TotalPrice: 100 * (i + 1),
				Deposit:    100 * (i + 1),
				SaleAt:     time.Now(),
				Details: []models.SaleDetail{
					{ItemID: item1.ID, Price: 100, Quantity: 1},
				},
			}
			err := saleService.CreateSale(sale)
			require.NoError(t, err)
		}

		sales, err := saleService.GetAllSales()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(sales), 3)
	})
}

func TestStoreService_Complete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	storeRepo := repository.NewStoreRepository(db)
	storeService := NewStoreService(storeRepo)

	t.Run("CreateStore", func(t *testing.T) {
		store := &models.Store{
			StoreID: "STORE001",
			Name:    "Test Store",
		}

		err := storeService.CreateStore(store)
		require.NoError(t, err)
		assert.NotZero(t, store.ID)
	})

	t.Run("GetStoreByID", func(t *testing.T) {
		store := &models.Store{
			StoreID: "STORE002",
			Name:    "Get Test Store",
		}
		err := storeService.CreateStore(store)
		require.NoError(t, err)

		retrieved, err := storeService.GetStoreByID(store.ID)
		require.NoError(t, err)
		assert.Equal(t, store.Name, retrieved.Name)
		assert.Equal(t, store.StoreID, retrieved.StoreID)
	})

	t.Run("UpdateStore", func(t *testing.T) {
		store := &models.Store{
			StoreID: "STORE003",
			Name:    "Update Test Store",
		}
		err := storeService.CreateStore(store)
		require.NoError(t, err)

		store.Name = "Updated Store"
		err = storeService.UpdateStore(store)
		require.NoError(t, err)

		retrieved, err := storeService.GetStoreByID(store.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Store", retrieved.Name)
	})

	t.Run("DeleteStore", func(t *testing.T) {
		store := &models.Store{
			StoreID: "STORE004",
			Name:    "Delete Test Store",
		}
		err := storeService.CreateStore(store)
		require.NoError(t, err)

		err = storeService.DeleteStore(store.ID)
		require.NoError(t, err)

		_, err = storeService.GetStoreByID(store.ID)
		assert.Error(t, err)
	})

	t.Run("GetAllStores", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			store := &models.Store{
				StoreID: "STORELIST" + string(rune('1'+i)),
				Name:    "List Store",
			}
			err := storeService.CreateStore(store)
			require.NoError(t, err)
		}

		stores, err := storeService.GetAllStores()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(stores), 3)
	})
}