package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/KidsPOSProject/KidsPOS-Server-GO/internal/models"
)

// Repositories holds all repository instances
type Repositories struct {
	Item    *ItemRepository
	Store   *StoreRepository
	Staff   *StaffRepository
	Sale    *SaleRepository
	Setting *SettingRepository
}

// NewRepositories creates all repository instances
func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Item:    &ItemRepository{db: db},
		Store:   &StoreRepository{db: db},
		Staff:   &StaffRepository{db: db},
		Sale:    &SaleRepository{db: db},
		Setting: &SettingRepository{db: db},
	}
}

// ItemRepository handles item data access
type ItemRepository struct {
	db *sql.DB
}

func (r *ItemRepository) FindAll() ([]*models.Item, error) {
	query := `SELECT id, itemId, name, price, stock, isDeleted, createdAt, updatedAt
			  FROM item WHERE isDeleted = 0 ORDER BY id DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		item := &models.Item{}
		err := rows.Scan(&item.ID, &item.ItemID, &item.Name, &item.Price,
			&item.Stock, &item.IsDeleted, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepository) FindByID(id int) (*models.Item, error) {
	query := `SELECT id, itemId, name, price, stock, isDeleted, createdAt, updatedAt
			  FROM item WHERE id = ? AND isDeleted = 0`

	item := &models.Item{}
	err := r.db.QueryRow(query, id).Scan(&item.ID, &item.ItemID, &item.Name,
		&item.Price, &item.Stock, &item.IsDeleted, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found")
	}
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *ItemRepository) Create(item *models.Item) error {
	query := `INSERT INTO item (itemId, name, price, stock, createdAt, updatedAt)
			  VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, item.ItemID, item.Name, item.Price,
		item.Stock, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	item.ID = int(id)
	item.CreatedAt = now
	item.UpdatedAt = now

	return nil
}

func (r *ItemRepository) Update(item *models.Item) error {
	query := `UPDATE item SET name = ?, price = ?, stock = ?, updatedAt = ?
			  WHERE id = ? AND isDeleted = 0`

	now := time.Now()
	_, err := r.db.Exec(query, item.Name, item.Price, item.Stock, now, item.ID)
	if err != nil {
		return err
	}
	item.UpdatedAt = now
	return nil
}

func (r *ItemRepository) Delete(id int) error {
	query := `UPDATE item SET isDeleted = 1, updatedAt = ? WHERE id = ?`

	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

// StoreRepository handles store data access
type StoreRepository struct {
	db *sql.DB
}

func (r *StoreRepository) FindAll() ([]*models.Store, error) {
	query := `SELECT id, storeId, name, createdAt, updatedAt FROM store ORDER BY id DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []*models.Store
	for rows.Next() {
		store := &models.Store{}
		err := rows.Scan(&store.ID, &store.StoreID, &store.Name,
			&store.CreatedAt, &store.UpdatedAt)
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	return stores, nil
}

func (r *StoreRepository) FindByID(id int) (*models.Store, error) {
	query := `SELECT id, storeId, name, createdAt, updatedAt FROM store WHERE id = ?`

	store := &models.Store{}
	err := r.db.QueryRow(query, id).Scan(&store.ID, &store.StoreID,
		&store.Name, &store.CreatedAt, &store.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("store not found")
	}
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (r *StoreRepository) Create(store *models.Store) error {
	query := `INSERT INTO store (storeId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, store.StoreID, store.Name, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	store.ID = int(id)
	store.CreatedAt = now
	store.UpdatedAt = now

	return nil
}

// StaffRepository handles staff data access
type StaffRepository struct {
	db *sql.DB
}

func (r *StaffRepository) FindAll() ([]*models.Staff, error) {
	query := `SELECT id, staffId, name, createdAt, updatedAt FROM staff ORDER BY id DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staffs []*models.Staff
	for rows.Next() {
		staff := &models.Staff{}
		err := rows.Scan(&staff.ID, &staff.StaffID, &staff.Name,
			&staff.CreatedAt, &staff.UpdatedAt)
		if err != nil {
			return nil, err
		}
		staffs = append(staffs, staff)
	}
	return staffs, nil
}

func (r *StaffRepository) FindByID(id int) (*models.Staff, error) {
	query := `SELECT id, staffId, name, createdAt, updatedAt FROM staff WHERE id = ?`

	staff := &models.Staff{}
	err := r.db.QueryRow(query, id).Scan(&staff.ID, &staff.StaffID,
		&staff.Name, &staff.CreatedAt, &staff.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("staff not found")
	}
	if err != nil {
		return nil, err
	}
	return staff, nil
}

func (r *StaffRepository) Create(staff *models.Staff) error {
	query := `INSERT INTO staff (staffId, name, createdAt, updatedAt) VALUES (?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, staff.StaffID, staff.Name, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	staff.ID = int(id)
	staff.CreatedAt = now
	staff.UpdatedAt = now

	return nil
}

// SaleRepository handles sale data access
type SaleRepository struct {
	db *sql.DB
}

func (r *SaleRepository) FindAll() ([]*models.Sale, error) {
	query := `SELECT s.id, s.storeId, s.staffId, s.totalPrice, s.deposit, s.saleAt,
			  s.createdAt, s.updatedAt, st.storeId, st.name, sf.staffId, sf.name
			  FROM sale s
			  JOIN store st ON s.storeId = st.id
			  JOIN staff sf ON s.staffId = sf.id
			  ORDER BY s.id DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []*models.Sale
	for rows.Next() {
		sale := &models.Sale{
			Store: &models.Store{},
			Staff: &models.Staff{},
		}
		err := rows.Scan(&sale.ID, &sale.StoreID, &sale.StaffID, &sale.TotalPrice,
			&sale.Deposit, &sale.SaleAt, &sale.CreatedAt, &sale.UpdatedAt,
			&sale.Store.StoreID, &sale.Store.Name,
			&sale.Staff.StaffID, &sale.Staff.Name)
		if err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}
	return sales, nil
}

func (r *SaleRepository) Create(sale *models.Sale) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Insert sale
	query := `INSERT INTO sale (storeId, staffId, totalPrice, deposit, saleAt, createdAt, updatedAt)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := tx.Exec(query, sale.StoreID, sale.StaffID, sale.TotalPrice,
		sale.Deposit, sale.SaleAt, now, now)
	if err != nil {
		return err
	}

	saleID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	sale.ID = int(saleID)

	// Insert sale details
	for _, detail := range sale.Details {
		query := `INSERT INTO sale_detail (saleId, itemId, quantity, price, createdAt, updatedAt)
				  VALUES (?, ?, ?, ?, ?, ?)`
		_, err := tx.Exec(query, saleID, detail.ItemID, detail.Quantity, detail.Price, now, now)
		if err != nil {
			return err
		}

		// Update item stock
		updateQuery := `UPDATE item SET stock = stock - ? WHERE id = ?`
		_, err = tx.Exec(updateQuery, detail.Quantity, detail.ItemID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SettingRepository handles setting data access
type SettingRepository struct {
	db *sql.DB
}

func (r *SettingRepository) FindAll() ([]*models.Setting, error) {
	query := `SELECT id, key, value, type, description, createdAt, updatedAt
			  FROM setting ORDER BY key`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []*models.Setting
	for rows.Next() {
		setting := &models.Setting{}
		err := rows.Scan(&setting.ID, &setting.Key, &setting.Value, &setting.Type,
			&setting.Description, &setting.CreatedAt, &setting.UpdatedAt)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

func (r *SettingRepository) Update(key, value string) error {
	query := `UPDATE setting SET value = ?, updatedAt = ? WHERE key = ?`

	_, err := r.db.Exec(query, value, time.Now(), key)
	return err
}
