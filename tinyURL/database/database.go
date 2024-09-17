package database

import (
	"tinyURL/config"
	"tinyURL/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.URL{})
	return db, nil
}