package service

import (
	"testing"
	"time"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/models"
)

// MockItemRepository is a mock implementation of ItemRepository for testing
type MockItemRepository struct {
	items []models.Item
}

func (m *MockItemRepository) FindAll() ([]*models.Item, error) {
	result := make([]*models.Item, len(m.items))
	for i := range m.items {
		item := m.items[i]
		result[i] = &item
	}
	return result, nil
}

func (m *MockItemRepository) FindByID(id int) (*models.Item, error) {
	for _, item := range m.items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, nil
}

func (m *MockItemRepository) Create(item *models.Item) error {
	item.ID = len(m.items) + 1
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	m.items = append(m.items, *item)
	return nil
}

func (m *MockItemRepository) Update(item *models.Item) error {
	for i, existingItem := range m.items {
		if existingItem.ID == item.ID {
			item.UpdatedAt = time.Now()
			m.items[i] = *item
			return nil
		}
	}
	return nil
}

func (m *MockItemRepository) Delete(id int) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items[i].IsDeleted = true
			return nil
		}
	}
	return nil
}

func TestItemService_CreateItem(t *testing.T) {
	tests := []struct {
		name    string
		item    *models.Item
		wantErr bool
	}{
		{
			name: "正常な商品作成",
			item: &models.Item{
				Name:  "テスト商品",
				Price: 100,
				Stock: 10,
			},
			wantErr: false,
		},
		{
			name: "商品名が空",
			item: &models.Item{
				Name:  "",
				Price: 100,
				Stock: 10,
			},
			wantErr: true,
		},
		{
			name: "価格が負の値",
			item: &models.Item{
				Name:  "テスト商品",
				Price: -100,
				Stock: 10,
			},
			wantErr: true,
		},
		{
			name: "在庫が負の値",
			item: &models.Item{
				Name:  "テスト商品",
				Price: 100,
				Stock: -10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockItemRepository{items: []models.Item{}}
			service := &ItemService{repo: repo}

			err := service.CreateItem(tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemService.CreateItem() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.item.ItemID == "" {
				t.Error("ItemService.CreateItem() should generate ItemID")
			}
		})
	}
}

func TestItemService_UpdateItem(t *testing.T) {
	existingItem := models.Item{
		ID:     1,
		ItemID: "ITEM-001",
		Name:   "既存商品",
		Price:  100,
		Stock:  10,
	}

	tests := []struct {
		name    string
		item    *models.Item
		wantErr bool
	}{
		{
			name: "正常な商品更新",
			item: &models.Item{
				ID:    1,
				Name:  "更新商品",
				Price: 200,
				Stock: 20,
			},
			wantErr: false,
		},
		{
			name: "商品名が空",
			item: &models.Item{
				ID:    1,
				Name:  "",
				Price: 100,
				Stock: 10,
			},
			wantErr: true,
		},
		{
			name: "価格が負の値",
			item: &models.Item{
				ID:    1,
				Name:  "更新商品",
				Price: -100,
				Stock: 10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockItemRepository{items: []models.Item{existingItem}}
			service := &ItemService{repo: repo}

			err := service.UpdateItem(tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemService.UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestItemService_GetAllItems(t *testing.T) {
	items := []models.Item{
		{ID: 1, Name: "商品1", Price: 100, Stock: 10},
		{ID: 2, Name: "商品2", Price: 200, Stock: 20},
		{ID: 3, Name: "商品3", Price: 300, Stock: 30},
	}

	repo := &MockItemRepository{items: items}
	service := &ItemService{repo: repo}

	result, err := service.GetAllItems()
	if err != nil {
		t.Fatalf("ItemService.GetAllItems() error = %v", err)
	}

	if len(result) != len(items) {
		t.Errorf("ItemService.GetAllItems() returned %d items, want %d", len(result), len(items))
	}
}

func TestItemService_GetItem(t *testing.T) {
	items := []models.Item{
		{ID: 1, Name: "商品1", Price: 100, Stock: 10},
		{ID: 2, Name: "商品2", Price: 200, Stock: 20},
	}

	repo := &MockItemRepository{items: items}
	service := &ItemService{repo: repo}

	tests := []struct {
		name      string
		id        int
		wantName  string
		wantError bool
	}{
		{
			name:      "存在する商品",
			id:        1,
			wantName:  "商品1",
			wantError: false,
		},
		{
			name:      "存在しない商品",
			id:        999,
			wantName:  "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GetItem(tt.id)
			if (err != nil) != tt.wantError {
				t.Errorf("ItemService.GetItem() error = %v, wantError %v", err, tt.wantError)
			}

			if result != nil && result.Name != tt.wantName {
				t.Errorf("ItemService.GetItem() returned item with name %s, want %s", result.Name, tt.wantName)
			}
		})
	}
}

func TestItemService_DeleteItem(t *testing.T) {
	items := []models.Item{
		{ID: 1, Name: "商品1", Price: 100, Stock: 10},
		{ID: 2, Name: "商品2", Price: 200, Stock: 20},
	}

	repo := &MockItemRepository{items: items}
	service := &ItemService{repo: repo}

	err := service.DeleteItem(1)
	if err != nil {
		t.Fatalf("ItemService.DeleteItem() error = %v", err)
	}

	// 削除後の確認
	if !repo.items[0].IsDeleted {
		t.Error("ItemService.DeleteItem() should mark item as deleted")
	}
}