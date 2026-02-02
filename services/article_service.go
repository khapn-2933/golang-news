package services

import (
	"errors"
	"news/dto"
	"news/repositories"
	"news/utils"
)

// ArticleService chứa business logic cho articles
type ArticleService struct {
	articleRepo *repositories.ArticleRepository
	tagRepo     *repositories.TagRepository
	userRepo    *repositories.UserRepository
	followRepo  *repositories.FollowRepository
}

// NewArticleService tạo instance mới của ArticleService
func NewArticleService() *ArticleService {
	return &ArticleService{
		articleRepo: repositories.NewArticleRepository(),
		tagRepo:     repositories.NewTagRepository(),
		userRepo:    repositories.NewUserRepository(),
		followRepo:  repositories.NewFollowRepository(),
	}
}

// CreateArticle tạo article mới
func (s *ArticleService) CreateArticle(authorID int, req dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
	// Tạo slug từ title
	baseSlug := utils.GenerateSlug(req.Article.Title)

	// Tạo slug unique
	slug := utils.GenerateUniqueSlug(baseSlug, func(slugStr string) bool {
		exists, _ := s.articleRepo.IsSlugExists(slugStr)
		return exists
	})

	// Tạo article
	article, err := s.articleRepo.Create(
		slug,
		req.Article.Title,
		req.Article.Description,
		req.Article.Body,
		authorID,
	)
	if err != nil {
		return nil, err
	}

	// Xử lý tags nếu có
	if len(req.Article.TagList) > 0 {
		tagIDs := []int{}
		for _, tagName := range req.Article.TagList {
			tag, err := s.tagRepo.GetOrCreate(tagName)
			if err != nil {
				return nil, err
			}
			tagIDs = append(tagIDs, tag.ID)
		}
		err = s.tagRepo.AddTagsToArticle(article.ID, tagIDs)
		if err != nil {
			return nil, err
		}
	}

	// Lấy article với đầy đủ thông tin để trả về
	return s.buildArticleResponse(article.ID, nil)
}

// GetArticle lấy article theo slug
func (s *ArticleService) GetArticle(slug string, currentUserID *int) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, errors.New("article not found")
	}

	return s.buildArticleResponse(article.ID, currentUserID)
}

// ListArticles lấy danh sách articles với filters và pagination
func (s *ArticleService) ListArticles(tag, author, favorited *string, limit, offset int, currentUserID *int) (*dto.ArticleListResponse, error) {
	// Lấy articles
	articles, err := s.articleRepo.List(tag, author, favorited, limit, offset)
	if err != nil {
		return nil, err
	}

	// Đếm tổng số
	count, err := s.articleRepo.Count(tag, author, favorited)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &dto.ArticleListResponse{
		Articles:      []dto.ArticleResponse{},
		ArticlesCount: count,
	}

	for _, article := range articles {
		articleResp, err := s.buildArticleResponse(article.ID, currentUserID)
		if err != nil {
			return nil, err
		}
		response.Articles = append(response.Articles, *articleResp)
	}

	return response, nil
}

// FeedArticles lấy articles từ users mà currentUser đang follow
func (s *ArticleService) FeedArticles(currentUserID int, limit, offset int) (*dto.ArticleListResponse, error) {
	// Lấy articles
	articles, err := s.articleRepo.Feed(currentUserID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Đếm tổng số
	count, err := s.articleRepo.FeedCount(currentUserID)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &dto.ArticleListResponse{
		Articles:      []dto.ArticleResponse{},
		ArticlesCount: count,
	}

	currentUserIDPtr := &currentUserID
	for _, article := range articles {
		articleResp, err := s.buildArticleResponse(article.ID, currentUserIDPtr)
		if err != nil {
			return nil, err
		}
		response.Articles = append(response.Articles, *articleResp)
	}

	return response, nil
}

// UpdateArticle cập nhật article
func (s *ArticleService) UpdateArticle(slug string, authorID int, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	// Lấy article hiện tại
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, errors.New("article not found")
	}

	// Kiểm tra quyền sở hữu
	if article.AuthorID != authorID {
		return nil, errors.New("permission denied")
	}

	// Chuẩn bị các giá trị để update
	var newSlug, title, description, body *string

	if req.Article.Title != nil {
		title = req.Article.Title
		// Nếu title thay đổi thì phải update slug
		if *title != article.Title {
			baseSlug := utils.GenerateSlug(*title)
			// Tạo slug unique (có thể trùng với slug hiện tại nếu không có article nào khác dùng)
			currentSlug := article.Slug
			newSlugValue := utils.GenerateUniqueSlug(baseSlug, func(slugStr string) bool {
				// Nếu slug trùng với slug hiện tại thì không tính là đã tồn tại
				if slugStr == currentSlug {
					return false
				}
				exists, _ := s.articleRepo.IsSlugExists(slugStr)
				return exists
			})
			newSlug = &newSlugValue
		}
	}

	if req.Article.Description != nil {
		description = req.Article.Description
	}

	if req.Article.Body != nil {
		body = req.Article.Body
	}

	// Update article
	updatedArticle, err := s.articleRepo.Update(article.ID, newSlug, title, description, body)
	if err != nil {
		return nil, err
	}

	// Build response
	return s.buildArticleResponse(updatedArticle.ID, &authorID)
}

// DeleteArticle xóa article
func (s *ArticleService) DeleteArticle(slug string, authorID int) error {
	// Lấy article
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return err
	}
	if article == nil {
		return errors.New("article not found")
	}

	// Kiểm tra quyền sở hữu
	if article.AuthorID != authorID {
		return errors.New("permission denied")
	}

	// Xóa article
	return s.articleRepo.Delete(article.ID)
}

// FavoriteArticle thêm article vào favorites
func (s *ArticleService) FavoriteArticle(slug string, userID int) (*dto.ArticleResponse, error) {
	// Lấy article
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, errors.New("article not found")
	}

	// Favorite article
	err = s.articleRepo.Favorite(userID, article.ID)
	if err != nil {
		return nil, err
	}

	// Build response
	return s.buildArticleResponse(article.ID, &userID)
}

// UnfavoriteArticle xóa article khỏi favorites
func (s *ArticleService) UnfavoriteArticle(slug string, userID int) (*dto.ArticleResponse, error) {
	// Lấy article
	article, err := s.articleRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, errors.New("article not found")
	}

	// Unfavorite article
	err = s.articleRepo.Unfavorite(userID, article.ID)
	if err != nil {
		return nil, err
	}

	// Build response
	return s.buildArticleResponse(article.ID, &userID)
}

// buildArticleResponse build ArticleResponse từ article ID
func (s *ArticleService) buildArticleResponse(articleID int, currentUserID *int) (*dto.ArticleResponse, error) {
	// Lấy article
	article, err := s.articleRepo.GetByID(articleID)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, errors.New("article not found")
	}

	// Lấy author
	author, err := s.articleRepo.GetAuthorByID(article.AuthorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.New("author not found")
	}

	// Lấy tags
	tags, err := s.articleRepo.GetTagsByArticleID(article.ID)
	if err != nil {
		return nil, err
	}

	// Kiểm tra favorited
	favorited := false
	if currentUserID != nil {
		favorited, err = s.articleRepo.IsFavorited(*currentUserID, article.ID)
		if err != nil {
			return nil, err
		}
	}

	// Kiểm tra following (nếu có currentUserID)
	following := false
	if currentUserID != nil && *currentUserID != article.AuthorID {
		following, err = s.followRepo.IsFollowing(*currentUserID, article.AuthorID)
		if err != nil {
			return nil, err
		}
	}

	// Build response
	response := &dto.ArticleResponse{}
	response.Article.Slug = article.Slug
	response.Article.Title = article.Title
	response.Article.Description = article.Description
	response.Article.Body = article.Body
	response.Article.TagList = tags
	response.Article.CreatedAt = article.CreatedAt.Format("2006-01-02T15:04:05.000Z")
	response.Article.UpdatedAt = article.UpdatedAt.Format("2006-01-02T15:04:05.000Z")
	response.Article.Favorited = favorited
	response.Article.FavoritesCount = article.FavoritesCount
	response.Article.Author.Username = author.Username
	response.Article.Author.Bio = author.Bio
	response.Article.Author.Image = author.Image
	response.Article.Author.Following = following

	return response, nil
}
