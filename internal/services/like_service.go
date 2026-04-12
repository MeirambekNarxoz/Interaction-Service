package services

import (
	"errors"
	"fmt"
	"interaction-service/internal/models"
	"interaction-service/internal/rabbitmq"
	"interaction-service/internal/repository"

	"github.com/google/uuid"
)

type LikeService struct {
	repo     *repository.LikeRepository
	producer *rabbitmq.RabbitMQProducer
}

func NewLikeService(repo *repository.LikeRepository, producer *rabbitmq.RabbitMQProducer) *LikeService {
	return &LikeService{repo: repo, producer: producer}
}

func (s *LikeService) AddLike(userID uint, targetType models.TargetType, targetID uint, contentAuthorID uint, directionID uint) error {
	if !models.IsValidTargetType(targetType) {
		return errors.New("invalid target_type")
	}

	exists, err := s.repo.Exists(userID, targetType, targetID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("like already exists")
	}

	like := &models.Like{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   targetID,
	}

	err = s.repo.Create(like)
	if err != nil {
		return err
	}

	// Publish Gamification Event if producer is available and target is a post/article
	if s.producer != nil && contentAuthorID > 0 {
		event := map[string]interface{}{
			"eventId":     uuid.New().String(),
			"userId":      contentAuthorID, // Recipient of XP
			"type":        "REACTION_RECEIVED",
			"targetId":    targetID,
			"directionId": directionID,
		}

		routingKey := fmt.Sprintf("user.action.%s.received", targetType)
		_ = s.producer.PublishEvent(routingKey, event)
	}

	return nil
}

func (s *LikeService) RemoveLike(userID uint, targetType models.TargetType, targetID uint) error {
	if !models.IsValidTargetType(targetType) {
		return errors.New("invalid target_type")
	}

	return s.repo.Delete(userID, targetType, targetID)
}

func (s *LikeService) CountLikes(targetType models.TargetType, targetID uint) (int64, error) {
	if !models.IsValidTargetType(targetType) {
		return 0, errors.New("invalid target_type")
	}

	return s.repo.CountByTarget(targetType, targetID)
}
