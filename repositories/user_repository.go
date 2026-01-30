package repositories

import (
	"database/sql"
	"news/database"
	"news/models"
	"time"
)

// UserRepository chứa các method để làm việc với bảng users
type UserRepository struct{}

// NewUserRepository tạo instance mới của UserRepository
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create tạo user mới trong database
func (r *UserRepository) Create(username, email, passwordHash string) (*models.User, error) {
	query := `INSERT INTO users (username, email, password_hash, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := database.DB.Exec(query, username, email, passwordHash, now, now)
	if err != nil {
		return nil, err
	}

	// Lấy ID vừa tạo
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Lấy user vừa tạo
	return r.GetByID(int(id))
}

// GetByID lấy user theo ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, bio, image, created_at, updated_at 
	          FROM users WHERE id = ?`

	user := &models.User{}
	var bio, image sql.NullString

	err := database.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&bio,
		&image,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User không tồn tại
		}
		return nil, err
	}

	// Convert NullString sang *string
	if bio.Valid {
		user.Bio = &bio.String
	}
	if image.Valid {
		user.Image = &image.String
	}

	return user, nil
}

// GetByEmail lấy user theo email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, bio, image, created_at, updated_at 
	          FROM users WHERE email = ?`

	user := &models.User{}
	var bio, image sql.NullString

	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&bio,
		&image,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User không tồn tại
		}
		return nil, err
	}

	// Convert NullString sang *string
	if bio.Valid {
		user.Bio = &bio.String
	}
	if image.Valid {
		user.Image = &image.String
	}

	return user, nil
}

// GetByUsername lấy user theo username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, bio, image, created_at, updated_at 
	          FROM users WHERE username = ?`

	user := &models.User{}
	var bio, image sql.NullString

	err := database.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&bio,
		&image,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User không tồn tại
		}
		return nil, err
	}

	// Convert NullString sang *string
	if bio.Valid {
		user.Bio = &bio.String
	}
	if image.Valid {
		user.Image = &image.String
	}

	return user, nil
}

// Update cập nhật thông tin user
func (r *UserRepository) Update(userID int, email, username, passwordHash *string, bio, image *string) (*models.User, error) {
	// Build query động dựa trên các field cần update
	query := "UPDATE users SET updated_at = ?"
	args := []interface{}{time.Now()}

	if email != nil {
		query += ", email = ?"
		args = append(args, *email)
	}
	if username != nil {
		query += ", username = ?"
		args = append(args, *username)
	}
	if passwordHash != nil {
		query += ", password_hash = ?"
		args = append(args, *passwordHash)
	}
	if bio != nil {
		query += ", bio = ?"
		args = append(args, *bio)
	}
	if image != nil {
		query += ", image = ?"
		args = append(args, *image)
	}

	query += " WHERE id = ?"
	args = append(args, userID)

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// Lấy user đã được update
	return r.GetByID(userID)
}
