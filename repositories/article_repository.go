package repositories

import (
	"database/sql"
	"news/database"
	"news/models"
	"strings"
	"time"
)

// ArticleRepository chứa các method để làm việc với bảng articles
type ArticleRepository struct{}

// NewArticleRepository tạo instance mới của ArticleRepository
func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{}
}

// Create tạo article mới trong database
func (r *ArticleRepository) Create(slug, title, description, body string, authorID int) (*models.Article, error) {
	query := `INSERT INTO articles (slug, title, description, body, author_id, favorites_count, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, 0, ?, ?)`

	now := time.Now()
	result, err := database.DB.Exec(query, slug, title, description, body, authorID, now, now)
	if err != nil {
		return nil, err
	}

	// Lấy ID vừa tạo
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Lấy article vừa tạo
	return r.GetByID(int(id))
}

// GetByID lấy article theo ID
func (r *ArticleRepository) GetByID(id int) (*models.Article, error) {
	query := `SELECT id, slug, title, description, body, author_id, favorites_count, created_at, updated_at 
	          FROM articles WHERE id = ?`

	article := &models.Article{}
	err := database.DB.QueryRow(query, id).Scan(
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Description,
		&article.Body,
		&article.AuthorID,
		&article.FavoritesCount,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return article, nil
}

// GetBySlug lấy article theo slug
func (r *ArticleRepository) GetBySlug(slug string) (*models.Article, error) {
	query := `SELECT id, slug, title, description, body, author_id, favorites_count, created_at, updated_at 
	          FROM articles WHERE slug = ?`

	article := &models.Article{}
	err := database.DB.QueryRow(query, slug).Scan(
		&article.ID,
		&article.Slug,
		&article.Title,
		&article.Description,
		&article.Body,
		&article.AuthorID,
		&article.FavoritesCount,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return article, nil
}

// List lấy danh sách articles với filters và pagination
// Filters: tag, author, favorited
func (r *ArticleRepository) List(tag, author, favorited *string, limit, offset int) ([]*models.Article, error) {
	// Build query với filters
	query := `SELECT DISTINCT a.id, a.slug, a.title, a.description, a.body, a.author_id, 
	          a.favorites_count, a.created_at, a.updated_at 
	          FROM articles a`

	joins := []string{}
	conditions := []string{}
	args := []interface{}{}

	// Filter by tag
	if tag != nil && *tag != "" {
		joins = append(joins, "INNER JOIN article_tags at ON a.id = at.article_id")
		joins = append(joins, "INNER JOIN tags t ON at.tag_id = t.id")
		conditions = append(conditions, "t.name = ?")
		args = append(args, *tag)
	}

	// Filter by author
	if author != nil && *author != "" {
		joins = append(joins, "INNER JOIN users u ON a.author_id = u.id")
		conditions = append(conditions, "u.username = ?")
		args = append(args, *author)
	}

	// Filter by favorited
	if favorited != nil && *favorited != "" {
		joins = append(joins, "INNER JOIN favorites f ON a.id = f.article_id")
		joins = append(joins, "INNER JOIN users u2 ON f.user_id = u2.id")
		conditions = append(conditions, "u2.username = ?")
		args = append(args, *favorited)
	}

	// Combine query
	if len(joins) > 0 {
		query += " " + strings.Join(joins, " ")
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Order by newest first và pagination
	query += " ORDER BY a.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		article := &models.Article{}
		err := rows.Scan(
			&article.ID,
			&article.Slug,
			&article.Title,
			&article.Description,
			&article.Body,
			&article.AuthorID,
			&article.FavoritesCount,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// Count đếm tổng số articles với filters
func (r *ArticleRepository) Count(tag, author, favorited *string) (int, error) {
	query := `SELECT COUNT(DISTINCT a.id) FROM articles a`

	joins := []string{}
	conditions := []string{}
	args := []interface{}{}

	// Filter by tag
	if tag != nil && *tag != "" {
		joins = append(joins, "INNER JOIN article_tags at ON a.id = at.article_id")
		joins = append(joins, "INNER JOIN tags t ON at.tag_id = t.id")
		conditions = append(conditions, "t.name = ?")
		args = append(args, *tag)
	}

	// Filter by author
	if author != nil && *author != "" {
		joins = append(joins, "INNER JOIN users u ON a.author_id = u.id")
		conditions = append(conditions, "u.username = ?")
		args = append(args, *author)
	}

	// Filter by favorited
	if favorited != nil && *favorited != "" {
		joins = append(joins, "INNER JOIN favorites f ON a.id = f.article_id")
		joins = append(joins, "INNER JOIN users u2 ON f.user_id = u2.id")
		conditions = append(conditions, "u2.username = ?")
		args = append(args, *favorited)
	}

	if len(joins) > 0 {
		query += " " + strings.Join(joins, " ")
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := database.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Feed lấy articles từ các users mà currentUser đang follow
func (r *ArticleRepository) Feed(currentUserID int, limit, offset int) ([]*models.Article, error) {
	query := `SELECT a.id, a.slug, a.title, a.description, a.body, a.author_id, 
	          a.favorites_count, a.created_at, a.updated_at 
	          FROM articles a
	          INNER JOIN follows f ON a.author_id = f.following_id
	          WHERE f.follower_id = ?
	          ORDER BY a.created_at DESC
	          LIMIT ? OFFSET ?`

	rows, err := database.DB.Query(query, currentUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		article := &models.Article{}
		err := rows.Scan(
			&article.ID,
			&article.Slug,
			&article.Title,
			&article.Description,
			&article.Body,
			&article.AuthorID,
			&article.FavoritesCount,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// FeedCount đếm tổng số articles trong feed
func (r *ArticleRepository) FeedCount(currentUserID int) (int, error) {
	query := `SELECT COUNT(*) 
	          FROM articles a
	          INNER JOIN follows f ON a.author_id = f.following_id
	          WHERE f.follower_id = ?`

	var count int
	err := database.DB.QueryRow(query, currentUserID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Update cập nhật article
func (r *ArticleRepository) Update(articleID int, slug, title, description, body *string) (*models.Article, error) {
	// Build query động
	query := "UPDATE articles SET updated_at = ?"
	args := []interface{}{time.Now()}

	if slug != nil {
		query += ", slug = ?"
		args = append(args, *slug)
	}
	if title != nil {
		query += ", title = ?"
		args = append(args, *title)
	}
	if description != nil {
		query += ", description = ?"
		args = append(args, *description)
	}
	if body != nil {
		query += ", body = ?"
		args = append(args, *body)
	}

	query += " WHERE id = ?"
	args = append(args, articleID)

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return r.GetByID(articleID)
}

// Delete xóa article
func (r *ArticleRepository) Delete(articleID int) error {
	query := `DELETE FROM articles WHERE id = ?`
	_, err := database.DB.Exec(query, articleID)
	return err
}

// IsSlugExists kiểm tra xem slug đã tồn tại chưa
func (r *ArticleRepository) IsSlugExists(slug string) (bool, error) {
	query := `SELECT COUNT(*) FROM articles WHERE slug = ?`
	var count int
	err := database.DB.QueryRow(query, slug).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Favorite thêm article vào favorites của user
func (r *ArticleRepository) Favorite(userID, articleID int) error {
	query := `INSERT INTO favorites (user_id, article_id) VALUES (?, ?)`
	_, err := database.DB.Exec(query, userID, articleID)
	if err != nil {
		// Nếu đã favorite rồi thì không báo lỗi
		return nil
	}

	// Tăng favorites_count
	updateQuery := `UPDATE articles SET favorites_count = favorites_count + 1 WHERE id = ?`
	_, err = database.DB.Exec(updateQuery, articleID)
	return err
}

// Unfavorite xóa article khỏi favorites của user
func (r *ArticleRepository) Unfavorite(userID, articleID int) error {
	query := `DELETE FROM favorites WHERE user_id = ? AND article_id = ?`
	result, err := database.DB.Exec(query, userID, articleID)
	if err != nil {
		return err
	}

	// Kiểm tra xem có xóa được không
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Nếu có xóa thì giảm favorites_count
	if rowsAffected > 0 {
		updateQuery := `UPDATE articles SET favorites_count = favorites_count - 1 WHERE id = ?`
		_, err = database.DB.Exec(updateQuery, articleID)
		return err
	}

	return nil
}

// IsFavorited kiểm tra xem user đã favorite article chưa
func (r *ArticleRepository) IsFavorited(userID, articleID int) (bool, error) {
	query := `SELECT COUNT(*) FROM favorites WHERE user_id = ? AND article_id = ?`
	var count int
	err := database.DB.QueryRow(query, userID, articleID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetTagsByArticleID lấy danh sách tags của article
func (r *ArticleRepository) GetTagsByArticleID(articleID int) ([]string, error) {
	query := `SELECT t.name 
	          FROM tags t
	          INNER JOIN article_tags at ON t.id = at.tag_id
	          WHERE at.article_id = ?
	          ORDER BY t.name`

	rows, err := database.DB.Query(query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetAuthorByID lấy thông tin author của article
func (r *ArticleRepository) GetAuthorByID(authorID int) (*models.User, error) {
	userRepo := NewUserRepository()
	return userRepo.GetByID(authorID)
}
