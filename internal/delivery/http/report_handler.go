package http

import (
	"net/http"
	"strconv"

	"interaction-service/internal/middleware"
	"interaction-service/internal/models"
	"interaction-service/internal/services"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) SubmitReport(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	var req models.SubmitReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	err := h.service.SubmitReport(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "report submitted successfully"})
}

func (h *ReportHandler) GetReports(c *gin.Context) {
	status := c.Query("status")
	roomIDStr := c.Query("room_id")
	
	var roomIDPtr *uint
	if roomIDStr != "" {
		id, err := strconv.ParseUint(roomIDStr, 10, 32)
		if err == nil {
			val := uint(id)
			roomIDPtr = &val
		}
	}

	reports, err := h.service.GetReports(status, roomIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (h *ReportHandler) UpdateReportStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report id"})
		return
	}

	var req models.UpdateReportStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.service.UpdateReportStatus(uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "report status updated"})
}
