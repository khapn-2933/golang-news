package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"news/controllers"
	"news/database"
	"news/dto"
	"news/middlewares"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRouter tạo router cho testing
func setupTestRouter() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Init database
	err := database.InitDB()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Setup router như trong main.go
	router := gin.New()
	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.AuthMiddleware())

	// Khởi tạo controllers
	authController := controllers.NewAuthController()
	profileController := controllers.NewProfileController()
	articleController := controllers.NewArticleController()
	commentController := controllers.NewCommentController()
	tagController := controllers.NewTagController()

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes
		api.POST("/users", authController.Register)
		api.POST("/users/login", authController.Login)
		api.GET("/user", middlewares.RequireAuth(), authController.GetCurrentUser)
		api.PUT("/user", middlewares.RequireAuth(), authController.UpdateCurrentUser)

		// Profile routes
		api.GET("/profiles/:username", profileController.GetProfile)
		api.POST("/profiles/:username/follow", middlewares.RequireAuth(), profileController.FollowUser)
		api.DELETE("/profiles/:username/follow", middlewares.RequireAuth(), profileController.UnfollowUser)

		// Article routes
		api.GET("/articles", articleController.ListArticles)
		api.GET("/articles/feed", middlewares.RequireAuth(), articleController.FeedArticles)
		api.GET("/articles/:slug", articleController.GetArticle)
		api.POST("/articles", middlewares.RequireAuth(), articleController.CreateArticle)
		api.PUT("/articles/:slug", middlewares.RequireAuth(), articleController.UpdateArticle)
		api.DELETE("/articles/:slug", middlewares.RequireAuth(), articleController.DeleteArticle)
		api.POST("/articles/:slug/favorite", middlewares.RequireAuth(), articleController.FavoriteArticle)
		api.DELETE("/articles/:slug/favorite", middlewares.RequireAuth(), articleController.UnfavoriteArticle)

		// Comment routes
		api.POST("/articles/:slug/comments", middlewares.RequireAuth(), commentController.AddComment)
		api.GET("/articles/:slug/comments", commentController.GetComments)
		api.DELETE("/articles/:slug/comments/:id", middlewares.RequireAuth(), commentController.DeleteComment)

		// Tag routes
		api.GET("/tags", tagController.GetTags)
	}

	return router
}

// TestRegisterAndLogin test flow register và login
func TestRegisterAndLogin(t *testing.T) {
	router := setupTestRouter()

	// Sử dụng unique username và email với timestamp để tránh conflict
	// Nếu user đã tồn tại, test sẽ skip register và chỉ test login
	username := "testuser_integration"
	email := "test_integration@example.com"

	// 1. Register user
	registerReq := dto.RegisterRequest{
		User: struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
		}{
			Username: username,
			Email:    email,
			Password: "password123",
		},
	}

	registerBody, _ := json.Marshal(registerReq)
	registerReqHTTP := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(registerBody))
	registerReqHTTP.Header.Set("Content-Type", "application/json")
	registerW := httptest.NewRecorder()
	router.ServeHTTP(registerW, registerReqHTTP)

	// Nếu user đã tồn tại (422), chỉ test login
	if registerW.Code == http.StatusUnprocessableEntity {
		t.Log("User already exists, testing login only")
	} else {
		// User mới được tạo, verify response
		assert.Equal(t, http.StatusOK, registerW.Code)
		var registerResp dto.UserResponse
		err := json.Unmarshal(registerW.Body.Bytes(), &registerResp)
		require.NoError(t, err)
		assert.NotEmpty(t, registerResp.User.Token)
		assert.Equal(t, username, registerResp.User.Username)
		assert.Equal(t, email, registerResp.User.Email)
	}

	// 2. Login với user vừa tạo
	loginReq := dto.LoginRequest{
		User: struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}{
			Email:    email, // Sử dụng email từ register
			Password: "password123",
		},
	}

	loginBody, _ := json.Marshal(loginReq)
	loginReqHTTP := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(loginBody))
	loginReqHTTP.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReqHTTP)

	assert.Equal(t, http.StatusOK, loginW.Code)
	var loginResp dto.UserResponse
	err := json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	require.NoError(t, err)
	assert.NotEmpty(t, loginResp.User.Token)
	assert.Equal(t, username, loginResp.User.Username)
	assert.Equal(t, email, loginResp.User.Email)
}

// TestArticleCRUD test flow create, read, update, delete article
func TestArticleCRUD(t *testing.T) {
	router := setupTestRouter()

	// 1. Register và login để lấy token
	token := registerAndLogin(t, router, "articleuser", "article@example.com")

	// 2. Create article
	createReq := dto.CreateArticleRequest{
		Article: struct {
			Title       string   `json:"title" binding:"required"`
			Description string   `json:"description" binding:"required"`
			Body        string   `json:"body" binding:"required"`
			TagList     []string `json:"tagList,omitempty"`
		}{
			Title:       "Test Article",
			Description: "Test Description",
			Body:        "Test Body Content",
			TagList:     []string{"test", "go"},
		},
	}

	createBody, _ := json.Marshal(createReq)
	createReqHTTP := httptest.NewRequest("POST", "/api/articles", bytes.NewBuffer(createBody))
	createReqHTTP.Header.Set("Content-Type", "application/json")
	createReqHTTP.Header.Set("Authorization", "Token "+token)
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReqHTTP)

	assert.Equal(t, http.StatusOK, createW.Code)
	var createResp dto.ArticleResponse
	err := json.Unmarshal(createW.Body.Bytes(), &createResp)
	require.NoError(t, err)
	assert.Equal(t, "test-article", createResp.Article.Slug)
	assert.Equal(t, "Test Article", createResp.Article.Title)
	articleSlug := createResp.Article.Slug

	// 3. Get article
	getReqHTTP := httptest.NewRequest("GET", "/api/articles/"+articleSlug, nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReqHTTP)

	assert.Equal(t, http.StatusOK, getW.Code)
	var getResp dto.ArticleResponse
	err = json.Unmarshal(getW.Body.Bytes(), &getResp)
	require.NoError(t, err)
	assert.Equal(t, "Test Article", getResp.Article.Title)

	// 4. Update article
	updateReq := dto.UpdateArticleRequest{
		Article: struct {
			Title       *string `json:"title,omitempty"`
			Description *string `json:"description,omitempty"`
			Body        *string `json:"body,omitempty"`
		}{
			Title: stringPtr("Updated Test Article"),
		},
	}

	updateBody, _ := json.Marshal(updateReq)
	updateReqHTTP := httptest.NewRequest("PUT", "/api/articles/"+articleSlug, bytes.NewBuffer(updateBody))
	updateReqHTTP.Header.Set("Content-Type", "application/json")
	updateReqHTTP.Header.Set("Authorization", "Token "+token)
	updateW := httptest.NewRecorder()
	router.ServeHTTP(updateW, updateReqHTTP)

	assert.Equal(t, http.StatusOK, updateW.Code)
	var updateResp dto.ArticleResponse
	err = json.Unmarshal(updateW.Body.Bytes(), &updateResp)
	require.NoError(t, err)
	assert.Equal(t, "Updated Test Article", updateResp.Article.Title)

	// 5. Delete article
	deleteReqHTTP := httptest.NewRequest("DELETE", "/api/articles/"+updateResp.Article.Slug, nil)
	deleteReqHTTP.Header.Set("Authorization", "Token "+token)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReqHTTP)

	assert.Equal(t, http.StatusOK, deleteW.Code)

	// 6. Verify article đã bị xóa
	getAfterDeleteHTTP := httptest.NewRequest("GET", "/api/articles/"+updateResp.Article.Slug, nil)
	getAfterDeleteW := httptest.NewRecorder()
	router.ServeHTTP(getAfterDeleteW, getAfterDeleteHTTP)

	assert.Equal(t, http.StatusNotFound, getAfterDeleteW.Code)
}

// TestCommentFlow test flow add và get comments
func TestCommentFlow(t *testing.T) {
	router := setupTestRouter()

	// 1. Register, login và tạo article
	token := registerAndLogin(t, router, "commentuser", "comment@example.com")
	articleSlug := createArticle(t, router, token, "Comment Test Article", "Test", "Body")

	// 2. Add comment
	commentReq := dto.CreateCommentRequest{
		Comment: struct {
			Body string `json:"body" binding:"required"`
		}{
			Body: "This is a test comment",
		},
	}

	commentBody, _ := json.Marshal(commentReq)
	addCommentReqHTTP := httptest.NewRequest("POST", "/api/articles/"+articleSlug+"/comments", bytes.NewBuffer(commentBody))
	addCommentReqHTTP.Header.Set("Content-Type", "application/json")
	addCommentReqHTTP.Header.Set("Authorization", "Token "+token)
	addCommentW := httptest.NewRecorder()
	router.ServeHTTP(addCommentW, addCommentReqHTTP)

	assert.Equal(t, http.StatusOK, addCommentW.Code)
	var addCommentResp dto.CommentResponse
	err := json.Unmarshal(addCommentW.Body.Bytes(), &addCommentResp)
	require.NoError(t, err)
	assert.Equal(t, "This is a test comment", addCommentResp.Comment.Body)
	commentID := addCommentResp.Comment.ID

	// 3. Get comments
	getCommentsReqHTTP := httptest.NewRequest("GET", "/api/articles/"+articleSlug+"/comments", nil)
	getCommentsW := httptest.NewRecorder()
	router.ServeHTTP(getCommentsW, getCommentsReqHTTP)

	assert.Equal(t, http.StatusOK, getCommentsW.Code)
	var getCommentsResp dto.CommentListResponse
	err = json.Unmarshal(getCommentsW.Body.Bytes(), &getCommentsResp)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(getCommentsResp.Comments), 1)

	// 4. Delete comment
	deleteCommentReqHTTP := httptest.NewRequest("DELETE", "/api/articles/"+articleSlug+"/comments/"+strconv.Itoa(commentID), nil)
	deleteCommentReqHTTP.Header.Set("Authorization", "Token "+token)
	deleteCommentW := httptest.NewRecorder()
	router.ServeHTTP(deleteCommentW, deleteCommentReqHTTP)

	assert.Equal(t, http.StatusOK, deleteCommentW.Code)
}

// TestFavoriteFlow test flow favorite và unfavorite article
func TestFavoriteFlow(t *testing.T) {
	router := setupTestRouter()

	// 1. Register 2 users
	token1 := registerAndLogin(t, router, "favoriteuser1", "favorite1@example.com")
	token2 := registerAndLogin(t, router, "favoriteuser2", "favorite2@example.com")

	// 2. User1 tạo article
	articleSlug := createArticle(t, router, token1, "Favorite Test Article", "Test", "Body")

	// 3. User2 favorite article
	favoriteReqHTTP := httptest.NewRequest("POST", "/api/articles/"+articleSlug+"/favorite", nil)
	favoriteReqHTTP.Header.Set("Authorization", "Token "+token2)
	favoriteW := httptest.NewRecorder()
	router.ServeHTTP(favoriteW, favoriteReqHTTP)

	assert.Equal(t, http.StatusOK, favoriteW.Code)
	var favoriteResp dto.ArticleResponse
	err := json.Unmarshal(favoriteW.Body.Bytes(), &favoriteResp)
	require.NoError(t, err)
	assert.True(t, favoriteResp.Article.Favorited)
	assert.Equal(t, 1, favoriteResp.Article.FavoritesCount)

	// 4. User2 unfavorite article
	unfavoriteReqHTTP := httptest.NewRequest("DELETE", "/api/articles/"+articleSlug+"/favorite", nil)
	unfavoriteReqHTTP.Header.Set("Authorization", "Token "+token2)
	unfavoriteW := httptest.NewRecorder()
	router.ServeHTTP(unfavoriteW, unfavoriteReqHTTP)

	assert.Equal(t, http.StatusOK, unfavoriteW.Code)
	var unfavoriteResp dto.ArticleResponse
	err = json.Unmarshal(unfavoriteW.Body.Bytes(), &unfavoriteResp)
	require.NoError(t, err)
	assert.False(t, unfavoriteResp.Article.Favorited)
	assert.Equal(t, 0, unfavoriteResp.Article.FavoritesCount)
}

// TestProfileFlow test flow get profile, follow và unfollow
func TestProfileFlow(t *testing.T) {
	router := setupTestRouter()

	// 1. Register 2 users
	token1 := registerAndLogin(t, router, "profileuser1", "profile1@example.com")
	_ = registerAndLogin(t, router, "profileuser2", "profile2@example.com") // User2 được tạo để User1 follow

	// 2. User1 get profile của User2 (chưa follow)
	getProfileReqHTTP := httptest.NewRequest("GET", "/api/profiles/profileuser2", nil)
	getProfileReqHTTP.Header.Set("Authorization", "Token "+token1)
	getProfileW := httptest.NewRecorder()
	router.ServeHTTP(getProfileW, getProfileReqHTTP)

	assert.Equal(t, http.StatusOK, getProfileW.Code)
	var getProfileResp dto.ProfileResponse
	err := json.Unmarshal(getProfileW.Body.Bytes(), &getProfileResp)
	require.NoError(t, err)
	assert.Equal(t, "profileuser2", getProfileResp.Profile.Username)
	assert.False(t, getProfileResp.Profile.Following)

	// 3. User1 follow User2
	followReqHTTP := httptest.NewRequest("POST", "/api/profiles/profileuser2/follow", nil)
	followReqHTTP.Header.Set("Authorization", "Token "+token1)
	followW := httptest.NewRecorder()
	router.ServeHTTP(followW, followReqHTTP)

	assert.Equal(t, http.StatusOK, followW.Code)
	var followResp dto.ProfileResponse
	err = json.Unmarshal(followW.Body.Bytes(), &followResp)
	require.NoError(t, err)
	assert.True(t, followResp.Profile.Following)

	// 4. User1 get profile của User2 (đã follow)
	getProfileAfterFollowReqHTTP := httptest.NewRequest("GET", "/api/profiles/profileuser2", nil)
	getProfileAfterFollowReqHTTP.Header.Set("Authorization", "Token "+token1)
	getProfileAfterFollowW := httptest.NewRecorder()
	router.ServeHTTP(getProfileAfterFollowW, getProfileAfterFollowReqHTTP)

	assert.Equal(t, http.StatusOK, getProfileAfterFollowW.Code)
	var getProfileAfterFollowResp dto.ProfileResponse
	err = json.Unmarshal(getProfileAfterFollowW.Body.Bytes(), &getProfileAfterFollowResp)
	require.NoError(t, err)
	assert.True(t, getProfileAfterFollowResp.Profile.Following)

	// 5. User1 unfollow User2
	unfollowReqHTTP := httptest.NewRequest("DELETE", "/api/profiles/profileuser2/follow", nil)
	unfollowReqHTTP.Header.Set("Authorization", "Token "+token1)
	unfollowW := httptest.NewRecorder()
	router.ServeHTTP(unfollowW, unfollowReqHTTP)

	assert.Equal(t, http.StatusOK, unfollowW.Code)
	var unfollowResp dto.ProfileResponse
	err = json.Unmarshal(unfollowW.Body.Bytes(), &unfollowResp)
	require.NoError(t, err)
	assert.False(t, unfollowResp.Profile.Following)
}

// Helper functions

// registerAndLogin helper để register và login, trả về token
func registerAndLogin(t *testing.T, router *gin.Engine, username, email string) string {
	// Register
	registerReq := dto.RegisterRequest{
		User: struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
		}{
			Username: username,
			Email:    email,
			Password: "password123",
		},
	}

	registerBody, _ := json.Marshal(registerReq)
	registerReqHTTP := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(registerBody))
	registerReqHTTP.Header.Set("Content-Type", "application/json")
	registerW := httptest.NewRecorder()
	router.ServeHTTP(registerW, registerReqHTTP)

	if registerW.Code != http.StatusOK {
		// Nếu đã tồn tại thì login
		loginReq := dto.LoginRequest{
			User: struct {
				Email    string `json:"email" binding:"required,email"`
				Password string `json:"password" binding:"required"`
			}{
				Email:    email,
				Password: "password123",
			},
		}

		loginBody, _ := json.Marshal(loginReq)
		loginReqHTTP := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(loginBody))
		loginReqHTTP.Header.Set("Content-Type", "application/json")
		loginW := httptest.NewRecorder()
		router.ServeHTTP(loginW, loginReqHTTP)

		var loginResp dto.UserResponse
		json.Unmarshal(loginW.Body.Bytes(), &loginResp)
		return loginResp.User.Token
	}

	var registerResp dto.UserResponse
	json.Unmarshal(registerW.Body.Bytes(), &registerResp)
	return registerResp.User.Token
}

// createArticle helper để tạo article, trả về slug
func createArticle(t *testing.T, router *gin.Engine, token, title, description, body string) string {
	createReq := dto.CreateArticleRequest{
		Article: struct {
			Title       string   `json:"title" binding:"required"`
			Description string   `json:"description" binding:"required"`
			Body        string   `json:"body" binding:"required"`
			TagList     []string `json:"tagList,omitempty"`
		}{
			Title:       title,
			Description: description,
			Body:        body,
		},
	}

	createBody, _ := json.Marshal(createReq)
	createReqHTTP := httptest.NewRequest("POST", "/api/articles", bytes.NewBuffer(createBody))
	createReqHTTP.Header.Set("Content-Type", "application/json")
	createReqHTTP.Header.Set("Authorization", "Token "+token)
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReqHTTP)

	var createResp dto.ArticleResponse
	json.Unmarshal(createW.Body.Bytes(), &createResp)
	return createResp.Article.Slug
}

// stringPtr helper để tạo pointer từ string
func stringPtr(s string) *string {
	return &s
}

