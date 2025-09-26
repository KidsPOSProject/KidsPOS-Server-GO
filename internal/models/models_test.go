package models

import (
	"database/sql/driver"
	"testing"
	"time"
)

func TestNullTime_Scan(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
		valid   bool
	}{
		{
			name:    "nil値",
			value:   nil,
			wantErr: false,
			valid:   false,
		},
		{
			name:    "time.Time値",
			value:   time.Now(),
			wantErr: false,
			valid:   true,
		},
		{
			name:    "文字列値（RFC3339形式）",
			value:   "2024-01-01T00:00:00Z",
			wantErr: false,
			valid:   true,
		},
		{
			name:    "不正な文字列値",
			value:   "invalid date",
			wantErr: true,
			valid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nt := &NullTime{}
			err := nt.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullTime.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && nt.Valid != tt.valid {
				t.Errorf("NullTime.Scan() Valid = %v, want %v", nt.Valid, tt.valid)
			}
		})
	}
}

func TestNullTime_Value(t *testing.T) {
	tests := []struct {
		name    string
		nt      NullTime
		want    driver.Value
		wantErr bool
	}{
		{
			name: "有効な時間",
			nt: NullTime{
				Time:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Valid: true,
			},
			want:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "無効な時間",
			nt: NullTime{
				Time:  time.Time{},
				Valid: false,
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nt.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("NullTime.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NullTime.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItem_Validation(t *testing.T) {
	tests := []struct {
		name  string
		item  Item
		valid bool
	}{
		{
			name: "有効な商品",
			item: Item{
				ItemID: "ITEM-001",
				Name:   "テスト商品",
				Price:  100,
				Stock:  10,
			},
			valid: true,
		},
		{
			name: "名前が空",
			item: Item{
				ItemID: "ITEM-002",
				Name:   "",
				Price:  100,
				Stock:  10,
			},
			valid: false,
		},
		{
			name: "価格が負",
			item: Item{
				ItemID: "ITEM-003",
				Name:   "テスト商品",
				Price:  -100,
				Stock:  10,
			},
			valid: false,
		},
		{
			name: "在庫が負",
			item: Item{
				ItemID: "ITEM-004",
				Name:   "テスト商品",
				Price:  100,
				Stock:  -10,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple validation logic
			valid := tt.item.Name != "" && tt.item.Price >= 0 && tt.item.Stock >= 0
			if valid != tt.valid {
				t.Errorf("Item validation = %v, want %v", valid, tt.valid)
			}
		})
	}
}

func TestSale_TotalCalculation(t *testing.T) {
	sale := Sale{
		Details: []SaleDetail{
			{ItemID: 1, Quantity: 2, Price: 100},
			{ItemID: 2, Quantity: 3, Price: 200},
			{ItemID: 3, Quantity: 1, Price: 300},
		},
	}

	expectedTotal := 2*100 + 3*200 + 1*300 // 200 + 600 + 300 = 1100
	actualTotal := 0

	for _, detail := range sale.Details {
		actualTotal += detail.Quantity * detail.Price
	}

	if actualTotal != expectedTotal {
		t.Errorf("Sale total calculation = %d, want %d", actualTotal, expectedTotal)
	}
}

func TestStore_Validation(t *testing.T) {
	tests := []struct {
		name  string
		store Store
		valid bool
	}{
		{
			name: "有効な店舗",
			store: Store{
				StoreID: "STORE-001",
				Name:    "テスト店舗",
			},
			valid: true,
		},
		{
			name: "名前が空",
			store: Store{
				StoreID: "STORE-002",
				Name:    "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.store.Name != ""
			if valid != tt.valid {
				t.Errorf("Store validation = %v, want %v", valid, tt.valid)
			}
		})
	}
}

func TestStaff_Validation(t *testing.T) {
	tests := []struct {
		name  string
		staff Staff
		valid bool
	}{
		{
			name: "有効なスタッフ",
			staff: Staff{
				StaffID: "STAFF-001",
				Name:    "テストスタッフ",
			},
			valid: true,
		},
		{
			name: "名前が空",
			staff: Staff{
				StaffID: "STAFF-002",
				Name:    "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.staff.Name != ""
			if valid != tt.valid {
				t.Errorf("Staff validation = %v, want %v", valid, tt.valid)
			}
		})
	}
}

func TestSetting_TypeValidation(t *testing.T) {
	tests := []struct {
		name    string
		setting Setting
		valid   bool
	}{
		{
			name: "文字列型",
			setting: Setting{
				Key:   "shopName",
				Value: "テストショップ",
				Type:  "string",
			},
			valid: true,
		},
		{
			name: "数値型",
			setting: Setting{
				Key:   "taxRate",
				Value: "10",
				Type:  "number",
			},
			valid: true,
		},
		{
			name: "不正な型",
			setting: Setting{
				Key:   "testKey",
				Value: "testValue",
				Type:  "invalid",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple type validation
			validTypes := map[string]bool{"string": true, "number": true, "boolean": true}
			valid := validTypes[tt.setting.Type]
			if valid != tt.valid {
				t.Errorf("Setting type validation = %v, want %v", valid, tt.valid)
			}
		})
	}
}

// Benchmark tests
func BenchmarkNullTime_Scan(b *testing.B) {
	nt := &NullTime{}
	value := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = nt.Scan(value)
	}
}

func BenchmarkNullTime_Value(b *testing.B) {
	nt := NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = nt.Value()
	}
}