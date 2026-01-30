package middlewares

import (
	"news/config"
	"news/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware xác thực JWT token từ header
// Format header: Authorization: Token <jwt>
// Nếu token hợp lệ, lưu userID vào context với key "userID"
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy token từ header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Không có token, nhưng một số endpoint không bắt buộc auth
			// Nên không abort ở đây, để controller quyết định
			c.Next()
			return
		}

		// Parse header: "Token <jwt>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Token" {
			AbortWithError(c, 401, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// Validate token
		cfg := config.LoadConfig()
		userID, err := utils.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			AbortWithError(c, 401, "Invalid or expired token")
			return
		}

		// Lưu userID vào context để các handler sau có thể sử dụng
		c.Set("userID", userID)
		c.Next()
	}
}

// RequireAuth middleware bắt buộc phải có authentication
// Sử dụng sau AuthMiddleware để đảm bảo user đã đăng nhập
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Kiểm tra xem có userID trong context không
		userID, exists := c.Get("userID")
		if !exists || userID == nil {
			AbortWithError(c, 401, "Authentication required")
			return
		}
		c.Next()
	}
}
