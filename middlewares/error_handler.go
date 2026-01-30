package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse định dạng error response theo RealWorld spec
type ErrorResponse struct {
	Errors map[string][]string `json:"errors"`
}

// ErrorHandler middleware xử lý lỗi và trả về format JSON thống nhất
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Kiểm tra xem có lỗi nào không
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Xử lý các loại lỗi khác nhau
			switch err.Type {
			case gin.ErrorTypeBind:
				// Lỗi validation từ binding
				c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
					Errors: map[string][]string{
						"body": {err.Error()},
					},
				})
			case gin.ErrorTypePublic:
				// Lỗi public (đã được xử lý)
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Errors: map[string][]string{
						"body": {err.Error()},
					},
				})
			default:
				// Lỗi internal server
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Errors: map[string][]string{
						"body": {"Internal server error"},
					},
				})
			}
		}
	}
}

// AbortWithError trả về error response và dừng request
func AbortWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Errors: map[string][]string{
			"body": {message},
		},
	})
	c.Abort()
}
