package repositories

import (
	"database/sql"
	"news/database"
	"news/models"
	"time"
)

// CommentRepository chứa các method để làm việc với bảng comments
type CommentRepository struct{}

// NewCommentRepository tạo instance mới của CommentRepository
func NewCommentRepository() *CommentRepository {
	return &CommentRepository{}
}

// Create tạo comment mới trong database
func (r *CommentRepository) Create(articleID, authorID int, body string) (*models.Comment, error) {
	query := `INSERT INTO comments (article_id, author_id, body, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := database.DB.Exec(query, articleID, authorID, body, now, now)
	if err != nil {
		return nil, err
	}

	// Lấy ID vừa tạo
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Lấy comment vừa tạo
	return r.GetByID(int(id))
}

// GetByID lấy comment theo ID
func (r *CommentRepository) GetByID(id int) (*models.Comment, error) {
	query := `SELECT id, article_id, author_id, body, created_at, updated_at 
	          FROM comments WHERE id = ?`

	comment := &models.Comment{}
	err := database.DB.QueryRow(query, id).Scan(
		&comment.ID,
		&comment.ArticleID,
		&comment.AuthorID,
		&comment.Body,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return comment, nil
}

// GetByArticleID lấy tất cả comments của một article
func (r *CommentRepository) GetByArticleID(articleID int) ([]*models.Comment, error) {
	query := `SELECT id, article_id, author_id, body, created_at, updated_at 
	          FROM comments 
	          WHERE article_id = ? 
	          ORDER BY created_at DESC`

	rows, err := database.DB.Query(query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.ArticleID,
			&comment.AuthorID,
			&comment.Body,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// Delete xóa comment
func (r *CommentRepository) Delete(commentID int) error {
	query := `DELETE FROM comments WHERE id = ?`
	_, err := database.DB.Exec(query, commentID)
	return err
}

// GetAuthorByID lấy thông tin author của comment
func (r *CommentRepository) GetAuthorByID(authorID int) (*models.User, error) {
	userRepo := NewUserRepository()
	return userRepo.GetByID(authorID)
}
