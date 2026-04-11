package main

import (
	"interaction-service/internal/config"
	"interaction-service/internal/database"
	"interaction-service/internal/delivery/http"
	"interaction-service/internal/models"
	"interaction-service/internal/repository"
	"interaction-service/internal/routes"
	"interaction-service/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)
	db.AutoMigrate(&models.Like{}, &models.Bookmark{})

	likeRepo := repository.NewLikeRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)

	likeService := services.NewLikeService(likeRepo)
	bookmarkService := services.NewBookmarkService(bookmarkRepo)

	likeHandler := http.NewLikeHandler(likeService)
	bookmarkHandler := http.NewBookmarkHandler(bookmarkService)

	r := gin.Default()
	routes.RegisterInteractionRoutes(
		r,
		likeHandler,
		bookmarkHandler,
	)

	if err := r.Run(":" + cfg.Port); err != nil {
		panic(err)
	}
}
