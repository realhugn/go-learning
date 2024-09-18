package service

import (
	"errors"
	"tinyURL/internal/models"
	"tinyURL/internal/repository"
	"tinyURL/pkg/shortener"

	"github.com/go-playground/validator/v10"
)

type URLService struct {
	repo      *repository.URLRepository
	shortener *shortener.Shortener
	validator *validator.Validate
}

func NewURLService(repo *repository.URLRepository) *URLService {
	return &URLService{
		repo:      repo,
		shortener: shortener.New(),
		validator: validator.New(),
	}
}

func (s *URLService) Shorten(originalURL string) (string, error) {
	if err := s.validator.Var(originalURL, "required,url"); err != nil {
		return "", errors.New("invalid URL format")
	}

	shortURL, _ := s.repo.FindByOriginalURL(originalURL)

	if shortURL != nil {
		return shortURL.ShortURL, nil
	}

	generatedURL := s.shortener.Generate()

	url := &models.URL{
		ShortURL: generatedURL,
		Original: originalURL,
	}

	if err := s.validator.Struct(url); err != nil {
		return "", err
	}

	if err := s.repo.Create(url); err != nil {
		return "", err
	}

	return generatedURL, nil
}

func (s *URLService) Original(shortURL string) (string, error) {
	if err := s.validator.Var(shortURL, "required,min=6,max=10"); err != nil {
		return "", errors.New("invalid short URL format")
	}
	url, err := s.repo.FindByShortURL(shortURL)
	if err != nil {
		return "", err
	}

	return url.Original, nil
}
