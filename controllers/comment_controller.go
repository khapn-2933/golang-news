package controllers

import (
	"net/http"
	"news/dto"
	"news/middlewares"
	"news/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentController xử lý các HTTP request liên quan đến comments
type CommentController struct {
	commentService *services.CommentService
}

// NewCommentController tạo instance mới của CommentController
func NewCommentController() *CommentController {
	return &CommentController{
		commentService: services.NewCommentService(),
	}
}

// AddComment thêm comment vào article
// POST /api/articles/:slug/comments
// Authentication: required
func (c *CommentController) AddComment(ctx *gin.Context) {
	slug := ctx.Param("slug")

	// Lấy userID từ context
	userID, exists := ctx.Get("userID")
	if !exists {
		middlewares.AbortWithError(ctx, http.StatusUnauthorized, "Authentication required")
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Invalid user ID")
		return
	}

	var req dto.CreateCommentRequest

	// Bind request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Gọi service
	response, err := c.commentService.AddComment(slug, userIDInt, req)
	if err != nil {
		if err.Error() == "article not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to add comment")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetComments lấy tất cả comments của article
// GET /api/articles/:slug/comments
// Authentication: optional
func (c *CommentController) GetComments(ctx *gin.Context) {
	slug := ctx.Param("slug")

	// Lấy userID từ context nếu có
	var currentUserID *int
	if userID, exists := ctx.Get("userID"); exists {
		if userIDInt, ok := userID.(int); ok {
			currentUserID = &userIDInt
		}
	}

	// Gọi service
	response, err := c.commentService.GetComments(slug, currentUserID)
	if err != nil {
		if err.Error() == "article not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to get comments")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteComment xóa comment
// DELETE /api/articles/:slug/comments/:id
// Authentication: required
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	slug := ctx.Param("slug")
	commentIDStr := ctx.Param("id")

	// Parse comment ID
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		middlewares.AbortWithError(ctx, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	// Lấy userID từ context
	userID, exists := ctx.Get("userID")
	if !exists {
		middlewares.AbortWithError(ctx, http.StatusUnauthorized, "Authentication required")
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Invalid user ID")
		return
	}

	// Gọi service
	err = c.commentService.DeleteComment(slug, commentID, userIDInt)
	if err != nil {
		if err.Error() == "article not found" || err.Error() == "comment not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "permission denied" {
			middlewares.AbortWithError(ctx, http.StatusForbidden, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to delete comment")
		return
	}

	ctx.Status(http.StatusOK)
}
