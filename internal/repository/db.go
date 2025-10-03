package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// InitDB initializes the database connection
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set SQLite pragmas for better performance
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA cache_size = 10000;",
		"PRAGMA temp_store = MEMORY;",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			log.Printf("Warning: failed to set pragma: %s", pragma)
		}
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// RunMigrations runs database migrations
func RunMigrations(db *sql.DB) error {
	// Create tables if not exists
	schema := `
	CREATE TABLE IF NOT EXISTS item (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		itemId TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		price INTEGER NOT NULL,
		stock INTEGER NOT NULL DEFAULT 0,
		isDeleted INTEGER NOT NULL DEFAULT 0,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS store (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		storeId TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS staff (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		staffId TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sale (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		storeId INTEGER NOT NULL,
		staffId INTEGER NOT NULL,
		totalPrice INTEGER NOT NULL,
		deposit INTEGER NOT NULL,
		saleAt DATETIME NOT NULL,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (storeId) REFERENCES store(id),
		FOREIGN KEY (staffId) REFERENCES staff(id)
	);

	CREATE TABLE IF NOT EXISTS sale_detail (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		saleId INTEGER NOT NULL,
		itemId INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		price INTEGER NOT NULL,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (saleId) REFERENCES sale(id) ON DELETE CASCADE,
		FOREIGN KEY (itemId) REFERENCES item(id)
	);

	CREATE TABLE IF NOT EXISTS setting (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT NOT NULL UNIQUE,
		value TEXT NOT NULL,
		type TEXT NOT NULL DEFAULT 'string',
		description TEXT,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Create indexes
	CREATE INDEX IF NOT EXISTS idx_item_itemId ON item(itemId);
	CREATE INDEX IF NOT EXISTS idx_item_isDeleted ON item(isDeleted);
	CREATE INDEX IF NOT EXISTS idx_sale_storeId ON sale(storeId);
	CREATE INDEX IF NOT EXISTS idx_sale_staffId ON sale(staffId);
	CREATE INDEX IF NOT EXISTS idx_sale_saleAt ON sale(saleAt);
	CREATE INDEX IF NOT EXISTS idx_sale_detail_saleId ON sale_detail(saleId);
	CREATE INDEX IF NOT EXISTS idx_sale_detail_itemId ON sale_detail(itemId);

	-- Insert default settings if not exists
	INSERT OR IGNORE INTO setting (key, value, type, description) VALUES
		('shopName', 'KidsPOS Shop', 'string', 'Shop name'),
		('receiptFooter', 'Thank you!', 'string', 'Receipt footer message'),
		('taxRate', '10', 'number', 'Tax rate in percentage'),
		('currency', 'JPY', 'string', 'Currency code');

	-- Insert sample data if tables are empty
	INSERT OR IGNORE INTO store (storeId, name) VALUES
		('STORE001', 'Main Store');

	INSERT OR IGNORE INTO staff (staffId, name) VALUES
		('STAFF001', 'Admin');
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}