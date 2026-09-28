package database

import (
	"fmt"
	"log"

	"github.com/abdallah-elngar/gormx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"products-api/internal/config"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var dialector gorm.Dialector
	switch cfg.Database.Driver {
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.DSN)
	default:
		return fmt.Errorf("unsupported driver: %s", cfg.Database.Driver)
	}

	logLevel := logger.Silent
	if cfg.IsDevelopment() {
		logLevel = logger.Info
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	// ✅ اربط بـ gormx
	gormx.SetDB(db)
	DB = db

	log.Println("✅ Database connected")
	return nil
}
