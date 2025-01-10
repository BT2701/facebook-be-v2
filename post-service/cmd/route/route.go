package route

import (
	"post-service/internal/adapters/inbound"
	"post-service/internal/adapters/outbound"
	"post-service/internal/app/service"
	"post-service/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"post-service/pkg/database"
	"os"
)

func SetupRouter() *echo.Echo {
	// Initialize MongoDB connection
	database.InitMongoDB()
	databaseName := os.Getenv("DB_NAME")
	postCollection := database.GetCollection(databaseName, "posts")
	commentCollection := database.GetCollection(databaseName, "comments")
	reactionCollection := database.GetCollection(databaseName, "reactions")
	storyCollection := database.GetCollection(databaseName, "stories")

	// Create repositories and services
	postRepo := outbound.NewPostRepository(postCollection)
	postService := service.NewPostService(postRepo)

	commentRepo := outbound.NewCommentRepository(commentCollection)
	commentService := service.NewCommentService(commentRepo)

	reactionRepo := outbound.NewReactionRepository(reactionCollection)
	reactionService := service.NewReactionService(reactionRepo)

	storyRepo := outbound.NewStoryRepository(storyCollection)
	storyService := service.NewStoryService(storyRepo)

	// Create handlers
	postHandler := inbound.NewPostHandler(postService)
	commentHandler := inbound.NewCommentHandler(commentService)
	reactionHandler := inbound.NewReactionHandler(reactionService)
	storyHandler := inbound.NewStoryHandler(storyService)

	// Set up Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			return next(c)
		}
	})
	e.Use(utils.CorsMiddleware())

	e.POST("/posts", postHandler.CreatePost)
	e.GET("/posts/:id", postHandler.GetPost)
	e.PUT("/posts/:id", postHandler.UpdatePost)
	e.DELETE("/posts/:id", postHandler.DeletePost)
	e.GET("/posts/user/:userID", postHandler.GetPostsByUserID)
	e.GET("/posts", postHandler.GetPosts)
	e.DELETE("/posts", postHandler.DeleteAllPosts)

	e.POST("/comments", commentHandler.CreateComment)
	e.GET("/comments/:id", commentHandler.GetComment)
	e.PUT("/comments/:id", commentHandler.UpdateComment)
	e.DELETE("/comments/:id", commentHandler.DeleteComment)

	e.POST("/reactions", reactionHandler.CreateReaction)
	e.GET("/reactions/:id", reactionHandler.GetReaction)
	e.PUT("/reactions/:id", reactionHandler.UpdateReaction)
	e.DELETE("/reactions/:id", reactionHandler.DeleteReaction)

	e.POST("/stories", storyHandler.CreateStory)
	e.GET("/stories/:id", storyHandler.GetStory)
	e.PUT("/stories/:id", storyHandler.UpdateStory)
	e.DELETE("/stories/:id", storyHandler.DeleteStory)
	e.GET("/stories/user/:userID", storyHandler.GetStoriesByUserID)
	e.GET("/stories", storyHandler.GetStories)
	e.DELETE("/stories", storyHandler.DeleteStories)

	return e
}
