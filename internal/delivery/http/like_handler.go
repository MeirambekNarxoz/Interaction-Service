package http

import (
	"net/http"
	"strconv"

	"interaction-service/internal/middleware"
	"interaction-service/internal/models"
	"interaction-service/internal/services"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	service *services.LikeService
}

func NewLikeHandler(service *services.LikeService) *LikeHandler {
	return &LikeHandler{service: service}
}

type likeRequest struct {
	TargetType string `json:"target_type" binding:"required"`
	TargetID   uint   `json:"target_id" binding:"required"`
	AuthorID   uint   `json:"author_id"`    // Optional, for gamification
	DirectionID uint  `json:"direction_id"` // Optional, for gamification
}

func (h *LikeHandler) AddLike(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	var req likeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.service.AddLike(userID, models.TargetType(req.TargetType), req.TargetID, req.AuthorID, req.DirectionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "like added",
	})
}

func (h *LikeHandler) RemoveLike(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	var req likeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.service.RemoveLike(userID, models.TargetType(req.TargetType), req.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "like removed",
	})
}

func (h *LikeHandler) CountLikes(c *gin.Context) {
	targetType := c.Query("target_type")
	targetIDStr := c.Query("target_id")

	targetID, err := strconv.ParseUint(targetIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target_id"})
		return
	}

	count, err := h.service.CountLikes(models.TargetType(targetType), uint(targetID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}
