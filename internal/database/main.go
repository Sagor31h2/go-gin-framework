package database

import (
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDb(dsn string) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// For production, might want to use a different logger or silence it
	if os.Getenv("GIN_MODE") == "release" {
		config.Logger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(sqlite.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database (%s): %w", dsn, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Basic connection pooling
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}
