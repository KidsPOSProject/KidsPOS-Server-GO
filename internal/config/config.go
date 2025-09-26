package config

import (
	"os"
)

type Config struct {
	DatabasePath       string
	Port               string
	ReceiptPrinterHost string
	ReceiptPrinterPort string
	QRCodeSize         int
	AllowedIPPrefix    string
	EncryptionKey      string
}

func New() *Config {
	return &Config{
		DatabasePath:       getEnv("DATABASE_PATH", "./kidspos.db"),
		Port:               getEnv("PORT", "8080"),
		ReceiptPrinterHost: getEnv("RECEIPT_PRINTER_HOST", "localhost"),
		ReceiptPrinterPort: getEnv("RECEIPT_PRINTER_PORT", "9100"),
		QRCodeSize:         getEnvAsInt("QR_CODE_SIZE", 200),
		AllowedIPPrefix:    getEnv("ALLOWED_IP_PREFIX", "192.168."),
		EncryptionKey:      getEnv("ENCRYPTION_KEY", "DefaultKidsPOSKey123!@#"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	// Simple implementation for now
	return defaultValue
}
