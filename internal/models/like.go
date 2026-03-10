package models

import "time"

type Like struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index:idx_like_unique,unique" json:"user_id"`
	TargetType TargetType `gorm:"type:varchar(20);not null;index:idx_like_unique,unique" json:"target_type"`
	TargetID   uint       `gorm:"not null;index:idx_like_unique,unique" json:"target_id"`
	CreatedAt  time.Time  `json:"created_at"`
}
