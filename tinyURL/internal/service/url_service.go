package service

import (
	"errors"
	"time"
	"tinyURL/database"
	"tinyURL/internal/models"
	"tinyURL/internal/repository"
	"tinyURL/pkg/shortener"
	"tinyURL/pkg/uidgenerator"

	"github.com/go-playground/validator/v10"
)

type URLService interface {
	Shorten(originalURL string) (string, error)
	Original(shortURL string) (string, error)
}

type urlService struct {
	repo         repository.URLRepository
	shortener    shortener.Shortener
	validator    validator.Validate
	id_generator uidgenerator.IDGenerator
	cache        database.RedisCache
}

func NewURLService(repo repository.URLRepository, cache *database.RedisCache) URLService {
	return &urlService{
		repo:         repo,
		shortener:    *shortener.New(),
		validator:    *validator.New(),
		id_generator: *uidgenerator.NewIDGenerator(1),
		cache:        *cache,
	}
}

func (s *urlService) Shorten(originalURL string) (string, error) {
	if err := s.validator.Var(originalURL, "required,url"); err != nil {
		return "", errors.New("invalid URL format")
	}

	shortURL, _ := s.repo.FindByOriginalURL(originalURL)

	if shortURL != nil {
		return shortURL.ShortURL, nil
	}

	// This ensure the uniqueness of the generated ID
	uid := s.id_generator.GenerateID()
	generatedURL := s.shortener.ToBase62(int(uid))

	url := &models.URL{
		Id:       uid,
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

func (s *urlService) Original(shortURL string) (string, error) {
	if err := s.validator.Var(shortURL, "required,min=6,max=10"); err != nil {
		return "", errors.New("invalid short URL format")
	}

	longURL, err := s.cache.Get(shortURL)
	if err == nil {
		return longURL, nil
	}

	url, err := s.repo.FindByShortURL(shortURL)
	if err != nil {
		return "", err
	}

	err = s.cache.Set(shortURL, url.Original, time.Hour*48)
	if err != nil {
		return "", err
	}

	return url.Original, nil
}
