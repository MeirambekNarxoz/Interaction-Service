package models

import "time"

type ReportStatus string

const (
	ReportStatusOpen     ReportStatus = "OPEN"
	ReportStatusRejected ReportStatus = "REJECTED"
	ReportStatusResolved ReportStatus = "RESOLVED"
)

type Report struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	ReporterID     uint         `json:"reporter_id" gorm:"not null"`
	TargetID       uint         `json:"target_id" gorm:"not null"`
	TargetAuthorID uint         `json:"target_author_id" gorm:"not null"`
	TargetType     TargetType   `json:"target_type" gorm:"not null"`
	Reason         string       `json:"reason" gorm:"not null"`
	Status         ReportStatus `json:"status" gorm:"default:OPEN"`
	CreatedAt      time.Time    `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type SubmitReportRequest struct {
	TargetID       uint       `json:"target_id" binding:"required"`
	TargetAuthorID uint       `json:"target_author_id" binding:"required"`
	TargetType     TargetType `json:"target_type" binding:"required"`
	Reason         string     `json:"reason" binding:"required"`
}

type UpdateReportStatusRequest struct {
	Status ReportStatus `json:"status" binding:"required"`
}
