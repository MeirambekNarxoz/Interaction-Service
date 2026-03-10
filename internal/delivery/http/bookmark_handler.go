package http

import (
	"net/http"

	"interaction-service/internal/middleware"
	"interaction-service/internal/models"
	"interaction-service/internal/services"

	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct {
	service *services.BookmarkService
}

func NewBookmarkHandler(service *services.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{service: service}
}

type bookmarkRequest struct {
	TargetType string `json:"target_type" binding:"required"`
	TargetID   uint   `json:"target_id" binding:"required"`
}

func (h *BookmarkHandler) AddBookmark(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	var req bookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.service.AddBookmark(userID, models.TargetType(req.TargetType), req.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "bookmark added",
	})
}

func (h *BookmarkHandler) RemoveBookmark(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	var req bookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.service.RemoveBookmark(userID, models.TargetType(req.TargetType), req.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "bookmark removed",
	})
}

func (h *BookmarkHandler) GetMyBookmarks(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	bookmarks, err := h.service.GetMyBookmarks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookmarks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bookmarks": bookmarks,
	})
}
