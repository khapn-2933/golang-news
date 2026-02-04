package controllers

import (
	"net/http"
	"news/middlewares"
	"news/services"

	"github.com/gin-gonic/gin"
)

// TagController xử lý các HTTP request liên quan đến tags
type TagController struct {
	tagService *services.TagService
}

// NewTagController tạo instance mới của TagController
func NewTagController() *TagController {
	return &TagController{
		tagService: services.NewTagService(),
	}
}

// GetTags lấy tất cả tags
// GET /api/tags
// Authentication: not required
func (c *TagController) GetTags(ctx *gin.Context) {
	// Gọi service
	response, err := c.tagService.GetAllTags()
	if err != nil {
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to get tags")
		return
	}

	ctx.JSON(http.StatusOK, response)
}
