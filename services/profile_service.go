package services

import (
	"errors"
	"news/dto"
	"news/repositories"
)

// ProfileService chứa business logic cho profiles
type ProfileService struct {
	userRepo   *repositories.UserRepository
	followRepo *repositories.FollowRepository
}

// NewProfileService tạo instance mới của ProfileService
func NewProfileService() *ProfileService {
	return &ProfileService{
		userRepo:   repositories.NewUserRepository(),
		followRepo: repositories.NewFollowRepository(),
	}
}

// GetProfile lấy thông tin profile của user theo username
// currentUserID có thể là nil nếu không có authentication
func (s *ProfileService) GetProfile(username string, currentUserID *int) (*dto.ProfileResponse, error) {
	// Lấy user theo username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Tạo response
	response := &dto.ProfileResponse{}
	response.Profile.Username = user.Username
	response.Profile.Bio = user.Bio
	response.Profile.Image = user.Image

	// Kiểm tra following status nếu có currentUserID
	if currentUserID != nil {
		isFollowing, err := s.followRepo.IsFollowing(*currentUserID, user.ID)
		if err != nil {
			return nil, err
		}
		response.Profile.Following = isFollowing
	} else {
		// Không có authentication, following = false
		response.Profile.Following = false
	}

	return response, nil
}

// FollowUser follow một user
// followerID là user đang thực hiện follow
// username là username của user được follow
func (s *ProfileService) FollowUser(followerID int, username string) (*dto.ProfileResponse, error) {
	// Lấy user được follow theo username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Không thể follow chính mình
	if followerID == user.ID {
		return nil, errors.New("cannot follow yourself")
	}

	// Kiểm tra xem đã follow chưa
	isFollowing, err := s.followRepo.IsFollowing(followerID, user.ID)
	if err != nil {
		return nil, err
	}

	// Nếu chưa follow thì tạo relationship
	if !isFollowing {
		err = s.followRepo.Follow(followerID, user.ID)
		if err != nil {
			return nil, err
		}
	}

	// Tạo response
	response := &dto.ProfileResponse{}
	response.Profile.Username = user.Username
	response.Profile.Bio = user.Bio
	response.Profile.Image = user.Image
	response.Profile.Following = true

	return response, nil
}

// UnfollowUser unfollow một user
func (s *ProfileService) UnfollowUser(followerID int, username string) (*dto.ProfileResponse, error) {
	// Lấy user được unfollow theo username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Xóa relationship follow
	err = s.followRepo.Unfollow(followerID, user.ID)
	if err != nil {
		return nil, err
	}

	// Tạo response
	response := &dto.ProfileResponse{}
	response.Profile.Username = user.Username
	response.Profile.Bio = user.Bio
	response.Profile.Image = user.Image
	response.Profile.Following = false

	return response, nil
}
