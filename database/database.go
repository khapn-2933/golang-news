package database

import (
	"database/sql"
	"news/config"

	_ "github.com/go-sql-driver/mysql"
)

// DB là global database connection
var DB *sql.DB

// InitDB khởi tạo kết nối database
func InitDB() error {
	cfg := config.LoadConfig()
	dsn := cfg.GetDSN()

	// Mở kết nối database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Kiểm tra kết nối
	if err := db.Ping(); err != nil {
		return err
	}

	// Cấu hình connection pool
	db.SetMaxOpenConns(25) // Số connection tối đa
	db.SetMaxIdleConns(5)  // Số connection idle tối đa

	DB = db
	return nil
}

// CloseDB đóng kết nối database
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
