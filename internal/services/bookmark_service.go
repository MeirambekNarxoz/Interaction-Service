package services

import (
	"errors"
	"interaction-service/internal/models"
	"interaction-service/internal/repository"
)

type BookmarkService struct {
	repo *repository.BookmarkRepository
}

func NewBookmarkService(repo *repository.BookmarkRepository) *BookmarkService {
	return &BookmarkService{repo: repo}
}

func (s *BookmarkService) AddBookmark(userID uint, targetType models.TargetType, targetID uint) error {
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

	return s.repo.Create(bookmark)
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
