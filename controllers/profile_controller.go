package controllers

import (
	"net/http"
	"news/middlewares"
	"news/services"

	"github.com/gin-gonic/gin"
)

// ProfileController xử lý các HTTP request liên quan đến profiles
type ProfileController struct {
	profileService *services.ProfileService
}

// NewProfileController tạo instance mới của ProfileController
func NewProfileController() *ProfileController {
	return &ProfileController{
		profileService: services.NewProfileService(),
	}
}

// GetProfile lấy thông tin profile của user
// GET /api/profiles/:username
// Authentication: optional
func (c *ProfileController) GetProfile(ctx *gin.Context) {
	username := ctx.Param("username")

	// Lấy userID từ context nếu có (từ auth middleware)
	var currentUserID *int
	if userID, exists := ctx.Get("userID"); exists {
		if userIDInt, ok := userID.(int); ok {
			currentUserID = &userIDInt
		}
	}

	// Gọi service để lấy profile
	response, err := c.profileService.GetProfile(username, currentUserID)
	if err != nil {
		if err.Error() == "user not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to get profile")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// FollowUser follow một user
// POST /api/profiles/:username/follow
// Authentication: required
func (c *ProfileController) FollowUser(ctx *gin.Context) {
	username := ctx.Param("username")

	// Lấy userID từ context (đã được set bởi auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		middlewares.AbortWithError(ctx, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Convert userID sang int
	userIDInt, ok := userID.(int)
	if !ok {
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Invalid user ID")
		return
	}

	// Gọi service để follow user
	response, err := c.profileService.FollowUser(userIDInt, username)
	if err != nil {
		if err.Error() == "user not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "cannot follow yourself" {
			middlewares.AbortWithError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to follow user")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UnfollowUser unfollow một user
// DELETE /api/profiles/:username/follow
// Authentication: required
func (c *ProfileController) UnfollowUser(ctx *gin.Context) {
	username := ctx.Param("username")

	// Lấy userID từ context
	userID, exists := ctx.Get("userID")
	if !exists {
		middlewares.AbortWithError(ctx, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Convert userID sang int
	userIDInt, ok := userID.(int)
	if !ok {
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Invalid user ID")
		return
	}

	// Gọi service để unfollow user
	response, err := c.profileService.UnfollowUser(userIDInt, username)
	if err != nil {
		if err.Error() == "user not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to unfollow user")
		return
	}

	ctx.JSON(http.StatusOK, response)
}
