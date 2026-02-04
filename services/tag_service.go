package services

import (
	"news/dto"
	"news/repositories"
)

// TagService chứa business logic cho tags
type TagService struct {
	tagRepo *repositories.TagRepository
}

// NewTagService tạo instance mới của TagService
func NewTagService() *TagService {
	return &TagService{
		tagRepo: repositories.NewTagRepository(),
	}
}

// GetAllTags lấy tất cả tags
func (s *TagService) GetAllTags() (*dto.TagListResponse, error) {
	// Lấy tất cả tags từ repository
	tags, err := s.tagRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Convert sang string array
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	// Build response
	response := &dto.TagListResponse{
		Tags: tagNames,
	}

	return response, nil
}
