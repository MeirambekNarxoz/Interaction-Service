package services

import (
	"errors"
	"interaction-service/internal/models"
	"interaction-service/internal/repository"
)

type LikeService struct {
	repo *repository.LikeRepository
}

func NewLikeService(repo *repository.LikeRepository) *LikeService {
	return &LikeService{repo: repo}
}

func (s *LikeService) AddLike(userID uint, targetType models.TargetType, targetID uint) error {
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

	return s.repo.Create(like)
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
