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

func (s *ReportService) GetReports(status string, roomID *uint, roles []string) ([]models.Report, error) {
	// If Admin, and no status specified, show both OPEN and ESCALATED
	if status == "" && hasRole(roles, "ADMIN") {
		// We pass empty status to repository, which normally defaults to OPEN.
		// Let's modify the repository to handle "actionable" reports for Admin.
		return s.repo.GetAllForAdmin(roomID)
	}

	// Default status logic for others
	if status == "" {
		if hasRole(roles, "MODERATOR") {
			status = string(models.ReportStatusOpen)
		}
	}
	return s.repo.GetAllActive(status, roomID)
}

func (s *ReportService) UpdateReportStatus(id uint, status models.ReportStatus, roles []string) error {
	isAdmin := hasRole(roles, "ADMIN")
	isMod := hasRole(roles, "MODERATOR")

	if !isAdmin && !isMod {
		return errors.New("unauthorized: insufficient permissions")
	}

	// Moderator logic
	if isMod && !isAdmin {
		if status != models.ReportStatusRejected && status != models.ReportStatusEscalated {
			return errors.New("moderators can only REJECT or ESCALATE reports to admins")
		}
	}

	// Check if exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("report not found")
	}

	err = s.repo.UpdateStatus(id, status)
	if err != nil {
		return err
	}

	return nil
}

func hasRole(roles []string, target string) bool {
	for _, r := range roles {
		if r == target {
			return true
		}
	}
	return false
}
