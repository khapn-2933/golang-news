package controllers

import (
	"net/http"
	"news/dto"
	"news/middlewares"
	"news/services"

	"github.com/gin-gonic/gin"
)

// AuthController xử lý các HTTP request liên quan đến authentication
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController tạo instance mới của AuthController
func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// Register xử lý đăng ký user mới
// POST /api/users
func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	// Bind request body vào struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Gọi service để đăng ký
	response, err := c.authService.Register(req)
	if err != nil {
		// Kiểm tra loại lỗi
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			middlewares.AbortWithError(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to register user")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// Login xử lý đăng nhập
// POST /api/users/login
func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	// Bind request body vào struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Gọi service để đăng nhập
	response, err := c.authService.Login(req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			middlewares.AbortWithError(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to login")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetCurrentUser lấy thông tin user hiện tại
// GET /api/user
// Yêu cầu authentication
func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
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

	// Gọi service để lấy thông tin user
	response, err := c.authService.GetCurrentUser(userIDInt)
	if err != nil {
		if err.Error() == "user not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to get user")
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateCurrentUser cập nhật thông tin user hiện tại
// PUT /api/user
// Yêu cầu authentication
func (c *AuthController) UpdateCurrentUser(ctx *gin.Context) {
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

	var req dto.UpdateUserRequest

	// Bind request body vào struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Gọi service để cập nhật user
	response, err := c.authService.UpdateUser(userIDInt, req)
	if err != nil {
		if err.Error() == "user not found" {
			middlewares.AbortWithError(ctx, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			middlewares.AbortWithError(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}
		middlewares.AbortWithError(ctx, http.StatusInternalServerError, "Failed to update user")
		return
	}

	ctx.JSON(http.StatusOK, response)
}
