package handlers

import (
	"net/http"
	"tinyURL/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type URLHandler interface {
	Shorten(c *gin.Context)
	Original(c *gin.Context)
}

type urlHandler struct {
	urlService service.URLService
	validator  validator.Validate
}

func NewURLHandler(urlService service.URLService) URLHandler {
	return &urlHandler{
		urlService: urlService,
		validator:  *validator.New(),
	}
}

type ShortenRequest struct {
	LongURL string `json:"long_url" validate:"required,url"`
}

func (h urlHandler) Shorten(c *gin.Context) {
	var input ShortenRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format " + err.Error()})
		return
	}

	if err := h.validator.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	shortURL, err := h.urlService.Shorten(input.LongURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8082/" + shortURL})
}

func (h urlHandler) Original(c *gin.Context) {
	shortURL := c.Param("short_url")

	if err := h.validator.Var(shortURL, "required,min=6,max=10"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid short URL format"})
		return
	}

	originalURL, err := h.urlService.Original(shortURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}
