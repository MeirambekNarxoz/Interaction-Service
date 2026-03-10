package repository

import (
	"interaction-service/internal/models"

	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) Exists(userID uint, targetType models.TargetType, targetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Like{}).
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Count(&count).Error

	return count > 0, err
}

func (r *LikeRepository) Create(like *models.Like) error {
	return r.db.Create(like).Error
}

func (r *LikeRepository) Delete(userID uint, targetType models.TargetType, targetID uint) error {
	return r.db.
		Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).
		Delete(&models.Like{}).Error
}

func (r *LikeRepository) CountByTarget(targetType models.TargetType, targetID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Like{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count).Error

	return count, err
}
