package services

import (
	"errors"
	"interaction-service/internal/models"
	"interaction-service/internal/repository"
)

type ReportService struct {
	repo *repository.ReportRepository
}

func NewReportService(repo *repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) SubmitReport(reporterID uint, req models.SubmitReportRequest) error {
	if !models.IsValidTargetType(req.TargetType) {
		return errors.New("invalid target_type")
	}

	report := &models.Report{
		ReporterID:     reporterID,
		TargetID:       req.TargetID,
		TargetAuthorID: req.TargetAuthorID,
		TargetType:     req.TargetType,
		Reason:         req.Reason,
		Status:         models.ReportStatusOpen,
		RoomID:         req.RoomID,
	}

	return s.repo.Create(report)
}

func (s *ReportService) GetReports(status string, roomID *uint) ([]models.Report, error) {
	return s.repo.GetAllActive(status, roomID)
}

func (s *ReportService) UpdateReportStatus(id uint, status models.ReportStatus) error {
	if status != models.ReportStatusRejected && status != models.ReportStatusResolved && status != models.ReportStatusEscalated {
		return errors.New("invalid status: must be REJECTED, RESOLVED or ESCALATED")
	}

	// Check if exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("report not found")
	}

	return s.repo.UpdateStatus(id, status)
}
