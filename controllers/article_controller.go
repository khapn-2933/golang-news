package controllers

import (
	"net/http"
	"news/dto"
	"news/middlewares"
	"news/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ArticleController xử lý các HTTP request liên quan đến articles
type ArticleController struct {
	articleService *services.ArticleService
}

// NewArticleController tạo instance mới của ArticleController
func NewArticleController() *ArticleController {
	return &ArticleController{
		articleService: services.NewArticleService(),
	}
}

// CreateArticle tạo article mới
// POST /api/articles
// Authentication: required
func (c *ArticleController) CreateArticle(ctx *gin.Context) {
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

	var req dto.CreateArticleRequest

	// Bind request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Gọi service
	response, err := c.articleService.CreateArticle(userIDInt, req)
	if err != nil {
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to create article")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetArticle lấy article theo slug
// GET /api/articles/:slug
// Authentication: optional
func (c *ArticleController) GetArticle(ctx *gin.Context) {
	slug := ctx.Param("slug")

	// Lấy userID từ context nếu có
	var currentUserID *int
	if userID, exists := ctx.Get("userID"); exists {
		if userIDInt, ok := userID.(int); ok {
			currentUserID = &userIDInt
		}
	}

	// Gọi service
	response, err := c.articleService.GetArticle(slug, currentUserID)
	if err != nil {
		if err.Error() == "article not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to get article")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ListArticles lấy danh sách articles
// GET /api/articles
// Query params: tag, author, favorited, limit, offset
func (c *ArticleController) ListArticles(ctx *gin.Context) {
	// Lấy query params
	tag := ctx.Query("tag")
	author := ctx.Query("author")
	favorited := ctx.Query("favorited")

	// Parse limit và offset
	limit := 20 // default
	offset := 0 // default

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			if l > 100 {
				limit = 100 // max limit
			} else {
				limit = l
			}
		}
	}

	if offsetStr := ctx.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Lấy userID từ context nếu có
	var currentUserID *int
	if userID, exists := ctx.Get("userID"); exists {
		if userIDInt, ok := userID.(int); ok {
			currentUserID = &userIDInt
		}
	}

	// Convert string pointers
	var tagPtr, authorPtr, favoritedPtr *string
	if tag != "" {
		tagPtr = &tag
	}
	if author != "" {
		authorPtr = &author
	}
	if favorited != "" {
		favoritedPtr = &favorited
	}

	// Gọi service
	response, err := c.articleService.ListArticles(tagPtr, authorPtr, favoritedPtr, limit, offset, currentUserID)
	if err != nil {
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to list articles")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// FeedArticles lấy articles từ users mà currentUser đang follow
// GET /api/articles/feed
// Authentication: required
func (c *ArticleController) FeedArticles(ctx *gin.Context) {
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

	// Parse limit và offset
	limit := 20 // default
	offset := 0 // default

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			if l > 100 {
				limit = 100 // max limit
			} else {
				limit = l
			}
		}
	}

	if offsetStr := ctx.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Gọi service
	response, err := c.articleService.FeedArticles(userIDInt, limit, offset)
	if err != nil {
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to get feed")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateArticle cập nhật article
// PUT /api/articles/:slug
// Authentication: required
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
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

	var req dto.UpdateArticleRequest

	// Bind request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Gọi service
	response, err := c.articleService.UpdateArticle(slug, userIDInt, req)
	if err != nil {
		if err.Error() == "article not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "permission denied" {
			middlewares.AbortWithError(ctx, http.StatusForbidden, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to update article")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteArticle xóa article
// DELETE /api/articles/:slug
// Authentication: required
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {
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

	// Gọi service
	err := c.articleService.DeleteArticle(slug, userIDInt)
	if err != nil {
		if err.Error() == "article not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "permission denied" {
			middlewares.AbortWithError(ctx, http.StatusForbidden, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to delete article")
		return
	}

	ctx.Status(http.StatusOK)
}

// FavoriteArticle favorite article
// POST /api/articles/:slug/favorite
// Authentication: required
func (c *ArticleController) FavoriteArticle(ctx *gin.Context) {
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

	// Gọi service
	response, err := c.articleService.FavoriteArticle(slug, userIDInt)
	if err != nil {
		if err.Error() == "article not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to favorite article")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UnfavoriteArticle unfavorite article
// DELETE /api/articles/:slug/favorite
// Authentication: required
func (c *ArticleController) UnfavoriteArticle(ctx *gin.Context) {
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

	// Gọi service
	response, err := c.articleService.UnfavoriteArticle(slug, userIDInt)
	if err != nil {
		if err.Error() == "article not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to unfavorite article")
		return
	}

	ctx.JSON(http.StatusOK, response)
}
