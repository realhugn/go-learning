package models

import "gorm.io/gorm"

type URL struct {
	gorm.Model
	ShortURL string `gorm:"uniqueIndex"`
	Original string
}

func (u *URL) Create(db *gorm.DB) error {
	return db.Create(u).Error
}

func GetURLByShortURL(db *gorm.DB, shortURL string) (*URL, error) {
	var url URL
	err := db.Where("short_url = ?", shortURL).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}
