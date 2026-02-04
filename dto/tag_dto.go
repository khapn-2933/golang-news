package dto

// TagListResponse định dạng response cho list tags
// Theo RealWorld spec: {"tags": ["tag1", "tag2", ...]}
type TagListResponse struct {
	Tags []string `json:"tags"`
}
