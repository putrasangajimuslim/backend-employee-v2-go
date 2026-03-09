// internal/infrastructure/database/gorm.go
package database

import (
	"fmt"
	"github.com/yourname/backend-employee-v2-go/internal/domain/entity"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/yourname/backend-employee-v2-go/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch cfg.DBDriver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.DBDsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres: %w", err)
		}
		log.Println("Connected to PostgreSQL database")
	default:
		db, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to sqlite: %w", err)
		}
		log.Println("Connected to SQLite database")
	}

	if err := RunMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

// cmd/api/main.go atau internal/infrastructure/database/gorm.go

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
	)
}
