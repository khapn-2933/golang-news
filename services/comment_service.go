package services

import (
	"errors"
	"news/dto"
	"news/repositories"
)

// CommentService chứa business logic cho comments
type CommentService struct {
	commentRepo *repositories.CommentRepository
	articleRepo *repositories.ArticleRepository
	userRepo    *repositories.UserRepository
	followRepo  *repositories.FollowRepository
}

// NewCommentService tạo instance mới của CommentService
func NewCommentService() *CommentService {
	return &CommentService{
		commentRepo: repositories.NewCommentRepository(),
		articleRepo: repositories.NewArticleRepository(),
		userRepo:    repositories.NewUserRepository(),
		followRepo:  repositories.NewFollowRepository(),
	}
}

// AddComment thêm comment vào article
func (s *CommentService) AddComment(slug string, authorID int, req dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	// Lấy article theo slug
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, errors.New("article not found")
	}

	// Tạo comment
	comment, err := s.commentRepo.Create(article.ID, authorID, req.Comment.Body)
	if err != nil {
		return nil, err
	}

	// Build response
	return s.buildCommentResponse(comment.ID, &authorID)
}

// GetComments lấy tất cả comments của article
func (s *CommentService) GetComments(slug string, currentUserID *int) (*dto.CommentListResponse, error) {
	// Lấy article theo slug
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, errors.New("article not found")
	}

	// Lấy comments
	comments, err := s.commentRepo.GetByArticleID(article.ID)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &dto.CommentListResponse{
		Comments: []dto.CommentResponse{},
	}

	for _, comment := range comments {
		commentResp, err := s.buildCommentResponse(comment.ID, currentUserID)
		if err != nil {
			return nil, err
		}
		response.Comments = append(response.Comments, *commentResp)
	}

	return response, nil
}

// DeleteComment xóa comment
func (s *CommentService) DeleteComment(slug string, commentID, userID int) error {
	// Lấy article theo slug
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return err
	}
	if article == nil {
		return errors.New("article not found")
	}

	// Lấy comment
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("comment not found")
	}

	// Kiểm tra comment thuộc về article
	if comment.ArticleID != article.ID {
		return errors.New("comment not found")
	}

	// Kiểm tra quyền sở hữu (chỉ author của comment mới được xóa)
	if comment.AuthorID != userID {
		return errors.New("permission denied")
	}

	// Xóa comment
	return s.commentRepo.Delete(commentID)
}

// buildCommentResponse build CommentResponse từ comment ID
func (s *CommentService) buildCommentResponse(commentID int, currentUserID *int) (*dto.CommentResponse, error) {
	// Lấy comment
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, errors.New("comment not found")
	}

	// Lấy author
	author, err := s.commentRepo.GetAuthorByID(comment.AuthorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.New("author not found")
	}

	// Kiểm tra following (nếu có currentUserID)
	following := false
	if currentUserID != nil && *currentUserID != comment.AuthorID {
		following, err = s.followRepo.IsFollowing(*currentUserID, comment.AuthorID)
		if err != nil {
			return nil, err
		}
	}

	// Build response
	response := &dto.CommentResponse{}
	response.Comment.ID = comment.ID
	response.Comment.Body = comment.Body
	response.Comment.CreatedAt = comment.CreatedAt.Format("2006-01-02T15:04:05.000Z")
	response.Comment.UpdatedAt = comment.UpdatedAt.Format("2006-01-02T15:04:05.000Z")
	response.Comment.Author.Username = author.Username
	response.Comment.Author.Bio = author.Bio
	response.Comment.Author.Image = author.Image
	response.Comment.Author.Following = following

	return response, nil
}
