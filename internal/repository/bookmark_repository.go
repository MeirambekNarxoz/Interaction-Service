package repository

import (
	"interaction-service/internal/models"

	"gorm.io/gorm"
)

type BookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

func (r *BookmarkRepository) Exists(userID uint, targetType models.TargetType, targetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Bookmark{}).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Count(&count).Error

	return count > 0, err
}

func (r *BookmarkRepository) Create(bookmark *models.Bookmark) error {
	return r.db.Create(bookmark).Error
}

func (r *BookmarkRepository) Delete(userID uint, targetType models.TargetType, targetID uint) error {
	return r.db.
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Delete(&models.Bookmark{}).Error
}

func (r *BookmarkRepository) FindByUser(userID uint) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&bookmarks).Error
	return bookmarks, err
}
