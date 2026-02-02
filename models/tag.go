package models

// Tag model đại diện cho bảng tags trong database
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
