package service

import (
	"fmt"
	"time"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/models"
	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/repository"
	"github.com/google/uuid"
)

// Services holds all service instances
type Services struct {
	Item    *ItemService
	Store   *StoreService
	Staff   *StaffService
	Sale    *SaleService
	Setting *SettingService
}

// NewServices creates all service instances
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Item:    &ItemService{repo: repos.Item},
		Store:   &StoreService{repo: repos.Store},
		Staff:   &StaffService{repo: repos.Staff},
		Sale:    &SaleService{repo: repos.Sale, itemRepo: repos.Item},
		Setting: &SettingService{repo: repos.Setting},
	}
}

// ItemService handles item business logic
type ItemService struct {
	repo *repository.ItemRepository
}

func (s *ItemService) GetAllItems() ([]*models.Item, error) {
	return s.repo.FindAll()
}

func (s *ItemService) GetItem(id int) (*models.Item, error) {
	return s.repo.FindByID(id)
}

func (s *ItemService) CreateItem(item *models.Item) error {
	// Generate item ID if not provided
	if item.ItemID == "" {
		item.ItemID = s.generateItemID()
	}

	// Validate
	if item.Name == "" {
		return fmt.Errorf("item name is required")
	}
	if item.Price < 0 {
		return fmt.Errorf("item price must be non-negative")
	}
	if item.Stock < 0 {
		return fmt.Errorf("item stock must be non-negative")
	}

	return s.repo.Create(item)
}

func (s *ItemService) UpdateItem(item *models.Item) error {
	// Validate
	if item.Name == "" {
		return fmt.Errorf("item name is required")
	}
	if item.Price < 0 {
		return fmt.Errorf("item price must be non-negative")
	}
	if item.Stock < 0 {
		return fmt.Errorf("item stock must be non-negative")
	}

	return s.repo.Update(item)
}

func (s *ItemService) DeleteItem(id int) error {
	return s.repo.Delete(id)
}

func (s *ItemService) GetItemByID(id int) (*models.Item, error) {
	return s.repo.FindByID(id)
}

func (s *ItemService) GetItemByBarcode(barcode string) (*models.Item, error) {
	return s.repo.FindByBarcode(barcode)
}

func (s *ItemService) generateItemID() string {
	// Generate unique item ID with prefix
	return fmt.Sprintf("ITEM-%s", uuid.New().String()[:8])
}

// StoreService handles store business logic
type StoreService struct {
	repo *repository.StoreRepository
}

func (s *StoreService) GetAllStores() ([]*models.Store, error) {
	return s.repo.FindAll()
}

func (s *StoreService) GetStore(id int) (*models.Store, error) {
	return s.repo.FindByID(id)
}

func (s *StoreService) CreateStore(store *models.Store) error {
	// Generate store ID if not provided
	if store.StoreID == "" {
		store.StoreID = s.generateStoreID()
	}

	// Validate
	if store.Name == "" {
		return fmt.Errorf("store name is required")
	}

	return s.repo.Create(store)
}

func (s *StoreService) generateStoreID() string {
	return fmt.Sprintf("STORE-%s", uuid.New().String()[:8])
}

// StaffService handles staff business logic
type StaffService struct {
	repo *repository.StaffRepository
}

func (s *StaffService) GetAllStaffs() ([]*models.Staff, error) {
	return s.repo.FindAll()
}

func (s *StaffService) GetStaff(id int) (*models.Staff, error) {
	return s.repo.FindByID(id)
}

func (s *StaffService) CreateStaff(staff *models.Staff) error {
	// Generate staff ID if not provided
	if staff.StaffID == "" {
		staff.StaffID = s.generateStaffID()
	}

	// Validate
	if staff.Name == "" {
		return fmt.Errorf("staff name is required")
	}

	return s.repo.Create(staff)
}

func (s *StaffService) GetAllStaff() ([]*models.Staff, error) {
	return s.repo.FindAll()
}

func (s *StaffService) GetStaffByBarcode(barcode string) (*models.Staff, error) {
	return s.repo.FindByBarcode(barcode)
}

func (s *StaffService) UpdateStaff(staff *models.Staff) error {
	// Validate
	if staff.Name == "" {
		return fmt.Errorf("staff name is required")
	}

	return s.repo.Update(staff)
}

func (s *StaffService) UpdateStaffByBarcode(barcode string, staff *models.Staff) error {
	// Find existing staff by barcode
	existingStaff, err := s.repo.FindByBarcode(barcode)
	if err != nil {
		return err
	}

	// Update fields
	existingStaff.Name = staff.Name
	return s.repo.Update(existingStaff)
}

func (s *StaffService) DeleteStaff(id int) error {
	return s.repo.Delete(id)
}

func (s *StaffService) DeleteStaffByBarcode(barcode string) error {
	// Find staff by barcode first
	staff, err := s.repo.FindByBarcode(barcode)
	if err != nil {
		return err
	}

	return s.repo.Delete(staff.ID)
}

func (s *StaffService) generateStaffID() string {
	return fmt.Sprintf("STAFF-%s", uuid.New().String()[:8])
}

// SaleService handles sale business logic
type SaleService struct {
	repo     *repository.SaleRepository
	itemRepo *repository.ItemRepository
}

func (s *SaleService) GetAllSales() ([]*models.Sale, error) {
	return s.repo.FindAll()
}

func (s *SaleService) CreateSale(sale *models.Sale) error {
	// Set sale time if not provided
	if sale.SaleAt.IsZero() {
		sale.SaleAt = time.Now()
	}

	// Validate
	if sale.StoreID <= 0 {
		return fmt.Errorf("store is required")
	}
	if sale.StaffID <= 0 {
		return fmt.Errorf("staff is required")
	}
	if len(sale.Details) == 0 {
		return fmt.Errorf("sale must have at least one item")
	}

	// Calculate total price
	totalPrice := 0
	for _, detail := range sale.Details {
		// Validate stock availability
		item, err := s.itemRepo.FindByID(detail.ItemID)
		if err != nil {
			return fmt.Errorf("item not found: %d", detail.ItemID)
		}
		if item.Stock < detail.Quantity {
			return fmt.Errorf("insufficient stock for item: %s", item.Name)
		}

		// Set price from item if not provided
		if detail.Price == 0 {
			detail.Price = item.Price
		}

		totalPrice += detail.Price * detail.Quantity
	}

	sale.TotalPrice = totalPrice

	// Set deposit to total price if not provided
	if sale.Deposit == 0 {
		sale.Deposit = sale.TotalPrice
	}

	return s.repo.Create(sale)
}

func (s *SaleService) GetSalesReport(startDate, endDate time.Time) ([]*models.Sale, error) {
	// For now, return all sales
	// TODO: Implement date filtering in repository
	return s.repo.FindAll()
}

// SettingService handles setting business logic
type SettingService struct {
	repo *repository.SettingRepository
}

func (s *SettingService) GetAllSettings() ([]*models.Setting, error) {
	return s.repo.FindAll()
}

func (s *SettingService) GetSetting(key string) (*models.Setting, error) {
	return s.repo.FindByKey(key)
}

func (s *SettingService) CreateSetting(setting *models.Setting) error {
	// Validate
	if setting.Key == "" {
		return fmt.Errorf("setting key is required")
	}
	if setting.Value == "" {
		return fmt.Errorf("setting value is required")
	}

	return s.repo.Create(setting)
}

func (s *SettingService) UpdateSetting(key, value string) error {
	if key == "" {
		return fmt.Errorf("setting key is required")
	}
	if value == "" {
		return fmt.Errorf("setting value is required")
	}

	return s.repo.Update(key, value)
}

func (s *SettingService) DeleteSetting(key string) error {
	if key == "" {
		return fmt.Errorf("setting key is required")
	}

	return s.repo.Delete(key)
}
