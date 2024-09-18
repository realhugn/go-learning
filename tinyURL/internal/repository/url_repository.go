package repository

import (
	"errors"
	"log"
	"tinyURL/internal/models"

	"gorm.io/gorm"
)

type URLRepository interface {
	Create(url *models.URL) error
	FindByOriginalURL(originalURL string) (*models.URL, error)
	FindByShortURL(shortURL string) (*models.URL, error)
	KeyExists(key string) bool
}

type postgresURLRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) URLRepository {
	return postgresURLRepository{
		db: db,
	}
}

func (r postgresURLRepository) Create(url *models.URL) error {
	return r.db.Create(url).Error
}

func (r postgresURLRepository) FindByOriginalURL(originalURL string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("original = ?", originalURL).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r postgresURLRepository) FindByShortURL(shortURL string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("short_url = ?", shortURL).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r postgresURLRepository) KeyExists(key string) bool {
	var url models.URL
	err := r.db.Where("key = ?", key).Scan(&url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		log.Printf("Error checking key existence: %v", err)
		return false
	}
	return true
}
