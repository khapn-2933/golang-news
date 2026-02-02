package models

import "time"

// Comment model đại diện cho bảng comments trong database
type Comment struct {
	ID        int       `json:"id"`
	ArticleID int       `json:"article_id"`
	AuthorID  int       `json:"author_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentWithAuthor chứa thông tin comment kèm thông tin author
type CommentWithAuthor struct {
	Comment
	Author struct {
		Username string  `json:"username"`
		Bio      *string `json:"bio"`
		Image    *string `json:"image"`
	} `json:"author"`
}
