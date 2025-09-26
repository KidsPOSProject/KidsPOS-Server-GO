package models

import (
	"database/sql/driver"
	"time"
)

// Item represents a product item
type Item struct {
	ID        int       `json:"id" db:"id"`
	ItemID    string    `json:"itemId" db:"itemId"`
	Name      string    `json:"name" db:"name"`
	Price     int       `json:"price" db:"price"`
	Stock     int       `json:"stock" db:"stock"`
	IsDeleted bool      `json:"isDeleted" db:"isDeleted"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}

// Store represents a store
type Store struct {
	ID        int       `json:"id" db:"id"`
	StoreID   string    `json:"storeId" db:"storeId"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}

// Staff represents a staff member
type Staff struct {
	ID        int       `json:"id" db:"id"`
	StaffID   string    `json:"staffId" db:"staffId"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}

// Sale represents a sale transaction
type Sale struct {
	ID         int          `json:"id" db:"id"`
	StoreID    int          `json:"storeId" db:"storeId"`
	StaffID    int          `json:"staffId" db:"staffId"`
	TotalPrice int          `json:"totalPrice" db:"totalPrice"`
	Deposit    int          `json:"deposit" db:"deposit"`
	SaleAt     time.Time    `json:"saleAt" db:"saleAt"`
	Details    []SaleDetail `json:"details,omitempty"`
	Store      *Store       `json:"store,omitempty"`
	Staff      *Staff       `json:"staff,omitempty"`
	CreatedAt  time.Time    `json:"createdAt" db:"createdAt"`
	UpdatedAt  time.Time    `json:"updatedAt" db:"updatedAt"`
}

// SaleDetail represents a sale detail item
type SaleDetail struct {
	ID        int       `json:"id" db:"id"`
	SaleID    int       `json:"saleId" db:"saleId"`
	ItemID    int       `json:"itemId" db:"itemId"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Price     int       `json:"price" db:"price"`
	Item      *Item     `json:"item,omitempty"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}

// Setting represents a configuration setting
type Setting struct {
	ID          int       `json:"id" db:"id"`
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	Type        string    `json:"type" db:"type"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updatedAt"`
}

// NullTime represents a nullable time.Time
type NullTime struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface
func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}
	nt.Valid = true
	switch v := value.(type) {
	case time.Time:
		nt.Time = v
	case string:
		var err error
		nt.Time, err = time.Parse(time.RFC3339, v)
		return err
	}
	return nil
}

// Value implements the driver Valuer interface
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
