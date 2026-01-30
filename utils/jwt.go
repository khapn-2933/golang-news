package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims chứa thông tin trong JWT token
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken tạo JWT token từ user ID và secret
// Token có thời hạn 24 giờ
func GenerateToken(userID int, secret string) (string, error) {
	// Tạo claims với user ID và thời gian hết hạn
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Hết hạn sau 24 giờ
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Tạo token với method HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ký token với secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken kiểm tra và parse JWT token
// Trả về user ID nếu token hợp lệ, error nếu không hợp lệ
func ValidateToken(tokenString, secret string) (int, error) {
	// Parse token với claims
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	// Lấy claims từ token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
