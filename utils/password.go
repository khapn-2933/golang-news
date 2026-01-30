package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hash password bằng bcrypt
// Cost factor = 10 (cân bằng giữa security và performance)
func HashPassword(password string) (string, error) {
	// Generate hash từ password với cost 10
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword so sánh password với hash
// Trả về true nếu password khớp, false nếu không khớp
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
