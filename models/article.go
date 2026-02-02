package models

import "time"

// Article model đại diện cho bảng articles trong database
type Article struct {
	ID             int       `json:"id"`
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	AuthorID       int       `json:"author_id"`
	FavoritesCount int       `json:"favorites_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ArticleWithAuthor chứa thông tin article kèm thông tin author
type ArticleWithAuthor struct {
	Article
	Author struct {
		Username string  `json:"username"`
		Bio      *string `json:"bio"`
		Image    *string `json:"image"`
	} `json:"author"`
	TagList   []string `json:"tagList"`
	Favorited bool     `json:"favorited"`
}
