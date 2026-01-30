package services

import (
	"errors"
	"news/config"
	"news/dto"
	"news/repositories"
	"news/utils"
)

// AuthService chứa business logic cho authentication
type AuthService struct {
	userRepo *repositories.UserRepository
}

// NewAuthService tạo instance mới của AuthService
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(),
	}
}

// Register đăng ký user mới
// Trả về UserResponse với JWT token
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Kiểm tra email đã tồn tại chưa
	existingUser, err := s.userRepo.GetByEmail(req.User.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Kiểm tra username đã tồn tại chưa
	existingUser, err = s.userRepo.GetByUsername(req.User.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.User.Password)
	if err != nil {
		return nil, err
	}

	// Tạo user mới
	user, err := s.userRepo.Create(req.User.Username, req.User.Email, passwordHash)
	if err != nil {
		return nil, err
	}

	// Tạo JWT token
	cfg := config.LoadConfig()
	token, err := utils.GenerateToken(user.ID, cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	// Tạo response
	response := &dto.UserResponse{}
	response.User.Email = user.Email
	response.User.Username = user.Username
	response.User.Token = token
	response.User.Bio = user.Bio
	response.User.Image = user.Image

	return response, nil
}

// Login đăng nhập user
// Trả về UserResponse với JWT token nếu email và password đúng
func (s *AuthService) Login(req dto.LoginRequest) (*dto.UserResponse, error) {
	// Tìm user theo email
	user, err := s.userRepo.GetByEmail(req.User.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Kiểm tra password
	if !utils.CheckPassword(req.User.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// Tạo JWT token
	cfg := config.LoadConfig()
	token, err := utils.GenerateToken(user.ID, cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	// Tạo response
	response := &dto.UserResponse{}
	response.User.Email = user.Email
	response.User.Username = user.Username
	response.User.Token = token
	response.User.Bio = user.Bio
	response.User.Image = user.Image

	return response, nil
}

// GetCurrentUser lấy thông tin user hiện tại từ userID
func (s *AuthService) GetCurrentUser(userID int) (*dto.UserResponse, error) {
	// Lấy user từ database
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Tạo JWT token mới (refresh token)
	cfg := config.LoadConfig()
	token, err := utils.GenerateToken(user.ID, cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	// Tạo response
	response := &dto.UserResponse{}
	response.User.Email = user.Email
	response.User.Username = user.Username
	response.User.Token = token
	response.User.Bio = user.Bio
	response.User.Image = user.Image

	return response, nil
}

// UpdateUser cập nhật thông tin user
func (s *AuthService) UpdateUser(userID int, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Kiểm tra user có tồn tại không
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Chuẩn bị các giá trị để update
	var email, username, passwordHash *string
	var bio, image *string

	if req.User.Email != nil {
		// Kiểm tra email mới có trùng với email của user khác không
		existingUser, err := s.userRepo.GetByEmail(*req.User.Email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("email already exists")
		}
		email = req.User.Email
	}

	if req.User.Username != nil {
		// Kiểm tra username mới có trùng với username của user khác không
		existingUser, err := s.userRepo.GetByUsername(*req.User.Username)
		if err != nil {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("username already exists")
		}
		username = req.User.Username
	}

	if req.User.Password != nil {
		// Hash password mới
		hash, err := utils.HashPassword(*req.User.Password)
		if err != nil {
			return nil, err
		}
		passwordHash = &hash
	}

	if req.User.Bio != nil {
		bio = req.User.Bio
	}

	if req.User.Image != nil {
		image = req.User.Image
	}

	// Update user
	updatedUser, err := s.userRepo.Update(userID, email, username, passwordHash, bio, image)
	if err != nil {
		return nil, err
	}

	// Tạo JWT token mới
	cfg := config.LoadConfig()
	token, err := utils.GenerateToken(updatedUser.ID, cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	// Tạo response
	response := &dto.UserResponse{}
	response.User.Email = updatedUser.Email
	response.User.Username = updatedUser.Username
	response.User.Token = token
	response.User.Bio = updatedUser.Bio
	response.User.Image = updatedUser.Image

	return response, nil
}
