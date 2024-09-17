package repository

import (
	"tinyURL/internal/models"

	"gorm.io/gorm"
)

type URLRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Create(url *models.URL) error {
	return r.db.Create(url).Error
}

func (r *URLRepository) FindByOriginalURL(originalURL string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("original = ?", originalURL).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *URLRepository) FindByShortURL(shortURL string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("short_url = ?", shortURL).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}
