package main

import (
	"log"
	"news/config"
	"news/controllers"
	"news/database"
	"news/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config từ environment variables
	cfg := config.LoadConfig()

	// Khởi tạo database connection
	err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDB()

	// Tạo Gin router
	router := gin.Default()

	// Middleware xử lý lỗi
	router.Use(middlewares.ErrorHandler())

	// Middleware xác thực (không bắt buộc, để controller quyết định)
	router.Use(middlewares.AuthMiddleware())

	// Khởi tạo controllers
	authController := controllers.NewAuthController()
	profileController := controllers.NewProfileController()
	articleController := controllers.NewArticleController()
	commentController := controllers.NewCommentController()

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
	}

	// Chạy server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
