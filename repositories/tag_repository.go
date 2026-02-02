package repositories

import (
	"database/sql"
	"news/database"
	"news/models"
	"strings"
)

// TagRepository chứa các method để làm việc với bảng tags
type TagRepository struct{}

// NewTagRepository tạo instance mới của TagRepository
func NewTagRepository() *TagRepository {
	return &TagRepository{}
}

// GetOrCreate lấy tag theo name, nếu không có thì tạo mới
func (r *TagRepository) GetOrCreate(name string) (*models.Tag, error) {
	// Tìm tag theo name
	query := `SELECT id, name FROM tags WHERE name = ?`
	tag := &models.Tag{}
	err := database.DB.QueryRow(query, name).Scan(&tag.ID, &tag.Name)

	if err == nil {
		// Tag đã tồn tại
		return tag, nil
	}

	if err != sql.ErrNoRows {
		// Lỗi khác
		return nil, err
	}

	// Tag chưa tồn tại, tạo mới
	insertQuery := `INSERT INTO tags (name) VALUES (?)`
	result, err := database.DB.Exec(insertQuery, name)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	tag.ID = int(id)
	tag.Name = name
	return tag, nil
}

// GetAll lấy tất cả tags
func (r *TagRepository) GetAll() ([]*models.Tag, error) {
	query := `SELECT id, name FROM tags ORDER BY name`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		tag := &models.Tag{}
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// AddTagsToArticle thêm tags vào article (many-to-many)
func (r *TagRepository) AddTagsToArticle(articleID int, tagIDs []int) error {
	// Xóa tags cũ của article trước
	deleteQuery := `DELETE FROM article_tags WHERE article_id = ?`
	_, err := database.DB.Exec(deleteQuery, articleID)
	if err != nil {
		return err
	}

	// Thêm tags mới
	if len(tagIDs) == 0 {
		return nil
	}

	insertQuery := `INSERT INTO article_tags (article_id, tag_id) VALUES `
	values := []interface{}{}
	placeholders := []string{}

	for _, tagID := range tagIDs {
		placeholders = append(placeholders, "(?, ?)")
		values = append(values, articleID, tagID)
	}

	insertQuery += strings.Join(placeholders, ", ")
	_, err = database.DB.Exec(insertQuery, values...)
	return err
}
