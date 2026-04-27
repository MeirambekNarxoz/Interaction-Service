package main

import (
	"fmt"
	"interaction-service/internal/config"
	"interaction-service/internal/database"
	"interaction-service/internal/delivery/http"
	"interaction-service/internal/rabbitmq"
	"interaction-service/internal/repository"
	"interaction-service/internal/routes"
	"interaction-service/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)

	// Create DB URL for golang-migrate
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
	database.RunMigrations(dbURL)

	// Initialize RabbitMQ Producer
	producer, err := rabbitmq.NewRabbitMQProducer(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to RabbitMQ: %v. Gamification events will be disabled.", err)
	} else {
		defer producer.Close()
	}

	likeRepo := repository.NewLikeRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)
	reportRepo := repository.NewReportRepository(db)

	likeService := services.NewLikeService(likeRepo, producer)
	bookmarkService := services.NewBookmarkService(bookmarkRepo, producer)
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
