package main

import (
	"interaction-service/internal/config"
	"interaction-service/internal/database"
	"interaction-service/internal/delivery/http"
	"interaction-service/internal/repository"
	"interaction-service/internal/routes"
	"interaction-service/internal/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)

	likeRepo := repository.NewLikeRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)

	likeService := services.NewLikeService(likeRepo)
	bookmarkService := services.NewBookmarkService(bookmarkRepo)

	likeHandler := http.NewLikeHandler(likeService)
	bookmarkHandler := http.NewBookmarkHandler(bookmarkService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:          12 * time.Hour,
	}))

	routes.RegisterInteractionRoutes(
		r,
		likeHandler,
		bookmarkHandler,
		cfg.JWTSecret,
	)

	if err := r.Run(":" + cfg.Port); err != nil {
		panic(err)
	}
}
