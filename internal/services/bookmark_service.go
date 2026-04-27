package services

import (
	"errors"
	"fmt"
	"interaction-service/internal/models"
	"interaction-service/internal/rabbitmq"
	"interaction-service/internal/repository"

	"github.com/google/uuid"
)

type BookmarkService struct {
	repo     *repository.BookmarkRepository
	producer *rabbitmq.RabbitMQProducer
}

func NewBookmarkService(repo *repository.BookmarkRepository, producer *rabbitmq.RabbitMQProducer) *BookmarkService {
	return &BookmarkService{repo: repo, producer: producer}
}

func (s *BookmarkService) AddBookmark(userID uint, targetType models.TargetType, targetID uint, contentAuthorID uint, directionID uint) error {
	if !models.IsValidTargetType(targetType) {
		return errors.New("invalid target_type")
	}

	exists, err := s.repo.Exists(userID, targetType, targetID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("bookmark already exists")
	}

	bookmark := &models.Bookmark{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   targetID,
	}

	err = s.repo.Create(bookmark)
	if err != nil {
		return err
	}

	// Publish Gamification Event for Bookmarks too (Reaction Received)
	if s.producer != nil && contentAuthorID > 0 {
		event := map[string]interface{}{
			"eventId":     uuid.New().String(),
			"userId":      contentAuthorID, // Recipient gets Reputation/XP
			"type":        "REACTION_RECEIVED",
			"targetId":    targetID,
			"directionId": directionID,
		}

		routingKey := fmt.Sprintf("user.action.%s.bookmarked", targetType)
		_ = s.producer.PublishEvent(routingKey, event)
	}

	return nil
}

func (s *BookmarkService) RemoveBookmark(userID uint, targetType models.TargetType, targetID uint) error {
	if !models.IsValidTargetType(targetType) {
		return errors.New("invalid target_type")
	}

	return s.repo.Delete(userID, targetType, targetID)
}

func (s *BookmarkService) GetMyBookmarks(userID uint) ([]models.Bookmark, error) {
	return s.repo.FindByUser(userID)
}
