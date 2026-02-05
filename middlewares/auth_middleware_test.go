package middlewares

import (
	"net/http"
	"net/http/httptest"
	"news/config"
	"news/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestAuthMiddleware_ValidToken kiểm tra auth middleware với token hợp lệ
func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware())
	
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if exists {
			c.JSON(http.StatusOK, gin.H{"userID": userID})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no userID"})
		}
	})

	// Generate token
	cfg := config.LoadConfig()
	token, err := utils.GenerateToken(123, cfg.JWTSecret)
	assert.NoError(t, err)

	// Test request với valid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Token "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "123")
}

// TestAuthMiddleware_InvalidToken kiểm tra auth middleware với token không hợp lệ
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware())
	
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if exists {
			c.JSON(http.StatusOK, gin.H{"userID": userID})
		} else {
			c.Status(http.StatusOK) // Không có userID nhưng vẫn OK vì auth là optional
		}
	})

	// Test request với invalid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Token invalid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Auth middleware sẽ abort với 401 khi token invalid
	// Vì trong AuthMiddleware, nếu ValidateToken fail thì sẽ abort
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddleware_NoToken kiểm tra auth middleware không có token
func TestAuthMiddleware_NoToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware())
	
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if exists {
			c.JSON(http.StatusOK, gin.H{"userID": userID})
		} else {
			c.Status(http.StatusOK)
		}
	})

	// Test request không có token
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestAuthMiddleware_InvalidFormat kiểm tra auth middleware với format header không đúng
func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware())
	
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Test request với format không đúng
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer token") // Sai format, phải là "Token token"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Auth middleware sẽ abort với error
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestRequireAuth_WithUserID kiểm tra RequireAuth middleware khi có userID
func TestRequireAuth_WithUserID(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware())
	router.Use(RequireAuth())
	
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Generate token
	cfg := config.LoadConfig()
	token, err := utils.GenerateToken(123, cfg.JWTSecret)
	assert.NoError(t, err)

	// Test request với valid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Token "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestRequireAuth_WithoutUserID kiểm tra RequireAuth middleware khi không có userID
func TestRequireAuth_WithoutUserID(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware())
	router.Use(RequireAuth())
	
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test request không có token
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// RequireAuth sẽ abort với 401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

