package routes

import (
	"interaction-service/internal/delivery/http"
	"interaction-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterInteractionRoutes(
	r *gin.Engine,
	likeHandler *http.LikeHandler,
	bookmarkHandler *http.BookmarkHandler,
) {
	api := r.Group("/api/interactions")

	api.GET("/likes/count", likeHandler.CountLikes)

	protected := api.Group("/")
	protected.Use(middleware.GatewayAuthMiddleware())
	{
		protected.POST("/likes", likeHandler.AddLike)
		protected.DELETE("/likes", likeHandler.RemoveLike)

		protected.POST("/bookmarks", bookmarkHandler.AddBookmark)
		protected.DELETE("/bookmarks", bookmarkHandler.RemoveBookmark)
		protected.GET("/bookmarks/my", bookmarkHandler.GetMyBookmarks)
	}
}
