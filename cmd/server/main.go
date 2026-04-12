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
	db.AutoMigrate(&models.Like{}, &models.Bookmark{}, &models.Report{})

	likeRepo := repository.NewLikeRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)
	reportRepo := repository.NewReportRepository(db)

	likeService := services.NewLikeService(likeRepo)
	bookmarkService := services.NewBookmarkService(bookmarkRepo)
	reportService := services.NewReportService(reportRepo)

	likeHandler := http.NewLikeHandler(likeService)
	bookmarkHandler := http.NewBookmarkHandler(bookmarkService)
	reportHandler := http.NewReportHandler(reportService)

	r := gin.Default()
	routes.RegisterInteractionRoutes(
		r,
		likeHandler,
		bookmarkHandler,
		reportHandler,
	)

	if err := r.Run(":" + cfg.Port); err != nil {
		panic(err)
	}
}
