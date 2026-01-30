package config

import (
	"fmt"
	"os"
)

// Config chứa tất cả các cấu hình của ứng dụng
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Port       string
}

// LoadConfig đọc các biến môi trường và trả về Config
// Nếu không có biến môi trường, sử dụng giá trị mặc định
func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "news_user"),
		DBPassword: getEnv("DB_PASSWORD", "news_password"),
		DBName:     getEnv("DB_NAME", "news_db"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		Port:       getEnv("PORT", "8080"),
	}
}

// GetDSN trả về Data Source Name cho MySQL connection
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// getEnv lấy giá trị từ environment variable, nếu không có thì dùng defaultValue
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
