package dto

// CreateCommentRequest định dạng request body cho tạo comment
// Theo RealWorld spec: {"comment": {"body": "..."}}
type CreateCommentRequest struct {
	Comment struct {
		Body string `json:"body" binding:"required"`
	} `json:"comment" binding:"required"`
}

// CommentResponse định dạng response theo RealWorld spec
// {"comment": {...}}
type CommentResponse struct {
	Comment struct {
		ID        int    `json:"id"`
		Body      string `json:"body"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
		Author    struct {
			Username  string  `json:"username"`
			Bio       *string `json:"bio"`
			Image     *string `json:"image"`
			Following bool    `json:"following"`
		} `json:"author"`
	} `json:"comment"`
}

// CommentListResponse định dạng response cho list comments
// {"comments": [...]}
type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
}
