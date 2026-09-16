package route

import (
	"os"
	"post-service/internal/adapters/inbound"
	"post-service/internal/adapters/outbound"
	"post-service/internal/app/service"
	"post-service/pkg/database"

	"github.com/BT2701/facebook-be-v2/shared/events"
	"github.com/BT2701/facebook-be-v2/shared/httpx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	bus := events.New(os.Getenv("REDIS_URI"))
	postHandler := inbound.NewPostHandler(postService)
	commentHandler := inbound.NewCommentHandler(commentService, postService, bus)
	reactionHandler := inbound.NewReactionHandler(reactionService, postService, bus)
	storyHandler := inbound.NewStoryHandler(storyService)

	// Set up Echo
	e := echo.New()
	e.Use(middleware.Logger())
	httpx.Apply(e)

	e.POST("/posts", postHandler.CreatePost)
	e.GET("/posts/:id", postHandler.GetPost)
	e.PUT("/posts/:id", postHandler.UpdatePost)
	e.DELETE("/posts/:id", postHandler.DeletePost)
	e.GET("/posts/user/:userID", postHandler.GetPostsByUserID)
	e.GET("/posts", postHandler.GetPosts)
	e.GET("/search", postHandler.SearchPosts)
	e.GET("/post-noti/:id/:currentUserId", postHandler.GetPostForNotification)

	e.POST("/comments", commentHandler.CreateComment)
	e.POST("/comment", commentHandler.CreateComment)
	e.GET("/comments/:id", commentHandler.GetComment)
	e.PUT("/comments/:id", commentHandler.UpdateComment)
	e.PUT("/comment/:id", commentHandler.UpdateComment)
	e.DELETE("/comments/:id", commentHandler.DeleteComment)
	e.DELETE("/comment/:id", commentHandler.DeleteComment)
	e.GET("/comments/post/:postID", commentHandler.GetCommentsByPostID)

	e.POST("/reactions", reactionHandler.CreateReaction)
	e.POST("/reaction", reactionHandler.CreateReaction)
	e.GET("/reactions/:id", reactionHandler.GetReaction)
	e.PUT("/reactions/:id", reactionHandler.UpdateReaction)
	e.DELETE("/reactions/:id", reactionHandler.DeleteReaction)
	e.GET("/reaction/:postId/:userId", reactionHandler.GetByPostAndUser)
	e.DELETE("/reaction/:postId/:userId", reactionHandler.DeleteByPostAndUser)

	e.POST("/stories", storyHandler.CreateStory)
	e.GET("/stories/:id", storyHandler.GetStory)
	e.PUT("/stories/:id", storyHandler.UpdateStory)
	e.DELETE("/stories/:id", storyHandler.DeleteStory)
	e.GET("/stories/user/:userID", storyHandler.GetStoriesByUserID)
	e.GET("/stories", storyHandler.GetStories)
	e.DELETE("/stories", storyHandler.DeleteStories)

	return e
}
