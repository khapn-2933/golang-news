package dto

// CreateArticleRequest định dạng request body cho tạo article
// Theo RealWorld spec: {"article": {"title": "...", "description": "...", "body": "...", "tagList": [...]}}
type CreateArticleRequest struct {
	Article struct {
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description" binding:"required"`
		Body        string   `json:"body" binding:"required"`
		TagList     []string `json:"tagList,omitempty"`
	} `json:"article" binding:"required"`
}

// UpdateArticleRequest định dạng request body cho cập nhật article
// Tất cả các field đều optional
type UpdateArticleRequest struct {
	Article struct {
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
		Body        *string `json:"body,omitempty"`
	} `json:"article"`
}

// ArticleResponse định dạng response theo RealWorld spec
// {"article": {...}}
type ArticleResponse struct {
	Article struct {
		Slug           string   `json:"slug"`
		Title          string   `json:"title"`
		Description    string   `json:"description"`
		Body           string   `json:"body"`
		TagList        []string `json:"tagList"`
		CreatedAt      string   `json:"createdAt"`
		UpdatedAt      string   `json:"updatedAt"`
		Favorited      bool     `json:"favorited"`
		FavoritesCount int      `json:"favoritesCount"`
		Author         struct {
			Username  string  `json:"username"`
			Bio       *string `json:"bio"`
			Image     *string `json:"image"`
			Following bool    `json:"following"`
		} `json:"author"`
	} `json:"article"`
}

// ArticleListResponse định dạng response cho list articles
// {"articles": [...], "articlesCount": 10}
type ArticleListResponse struct {
	Articles      []ArticleResponse `json:"articles"`
	ArticlesCount int               `json:"articlesCount"`
}
