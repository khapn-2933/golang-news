package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestErrorHandler_BindError kiểm tra error handler với bind error
func TestErrorHandler_BindError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandler())
	
	router.POST("/test", func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			c.Error(err).SetType(gin.ErrorTypeBind)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test request với invalid JSON
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Error handler sẽ trả về 422
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// TestAbortWithError kiểm tra AbortWithError helper
func TestAbortWithError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.GET("/test", func(c *gin.Context) {
		AbortWithError(c, http.StatusBadRequest, "test error")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "test error")
	assert.Contains(t, w.Body.String(), "errors")
}

