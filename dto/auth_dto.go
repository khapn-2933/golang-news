package dto

// RegisterRequest định dạng request body cho đăng ký
// Theo RealWorld spec: {"user": {"username": "...", "email": "...", "password": "..."}}
type RegisterRequest struct {
	User struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	} `json:"user" binding:"required"`
}

// LoginRequest định dạng request body cho đăng nhập
// Theo RealWorld spec: {"user": {"email": "...", "password": "..."}}
type LoginRequest struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	} `json:"user" binding:"required"`
}

// UpdateUserRequest định dạng request body cho cập nhật user
// Tất cả các field đều optional
type UpdateUserRequest struct {
	User struct {
		Email    *string `json:"email,omitempty"`
		Username *string `json:"username,omitempty"`
		Password *string `json:"password,omitempty"`
		Bio      *string `json:"bio,omitempty"`
		Image    *string `json:"image,omitempty"`
	} `json:"user"`
}

// UserResponse định dạng response theo RealWorld spec
// {"user": {"email": "...", "token": "...", "username": "...", "bio": "...", "image": "..."}}
type UserResponse struct {
	User struct {
		Email    string  `json:"email"`
		Token    string  `json:"token"`
		Username string  `json:"username"`
		Bio      *string `json:"bio"`
		Image    *string `json:"image"`
	} `json:"user"`
}
