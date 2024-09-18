package main

import (
	"log"
	"tinyURL/config"
	"tinyURL/database"
	"tinyURL/internal/api"
	"tinyURL/internal/handlers"
	"tinyURL/internal/repository"
	"tinyURL/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo)
	urlHandler := handlers.NewURLHandler(urlService)

	router := gin.Default()
	api.SetupRoutes(router, urlHandler)

	log.Printf("Server running on :%s", cfg.ServerPort)
	err = router.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
