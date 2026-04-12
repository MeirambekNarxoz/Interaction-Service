package repository

import (
	"interaction-service/internal/models"

	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(report *models.Report) error {
	return r.db.Create(report).Error
}

func (r *ReportRepository) GetAllActive() ([]models.Report, error) {
	var reports []models.Report
	err := r.db.Where("status = ?", models.ReportStatusOpen).
		Order("created_at DESC").
		Find(&reports).Error
	return reports, err
}

func (r *ReportRepository) GetByID(id uint) (*models.Report, error) {
	var report models.Report
	err := r.db.First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) UpdateStatus(id uint, status models.ReportStatus) error {
	return r.db.Model(&models.Report{}).Where("id = ?", id).Update("status", status).Error
}
