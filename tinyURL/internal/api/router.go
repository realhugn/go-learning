package api

import (
	"tinyURL/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, urlHandler *handlers.URLHandler) {

	router.POST("/shorten", urlHandler.Shorten)
	router.GET("/:short_url", urlHandler.Original)
}
