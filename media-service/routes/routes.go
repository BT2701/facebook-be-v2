package routes

import (
	"os"
	"snake_api/controllers"
	"snake_api/repositories"
	"snake_api/services"
	"snake_api/utils"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(ImageCollection *mongo.Collection) *echo.Echo {
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URI"), // Địa chỉ Redis
	})
	// Khởi tạo repository, service, và controller
	ImageRepo := repositories.NewImageRepository(ImageCollection)
	ImageService := services.NewImageService(ImageRepo, redisClient)
	ImageController := controllers.NewImageController(ImageService)

	e := echo.New()

	// Enable CORS
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			return next(c)
		}
	})
	e.Use(utils.CorsMiddleware())

	// Routes

	e.POST("/image", ImageController.InsertImage)
	e.GET("/images", ImageController.FindAllImages)
	e.DELETE("/images", ImageController.DeleteAllImages)
	e.PUT("/image", ImageController.EditImage)
	e.GET("/image/user/:id", ImageController.GetImageByUserID)
	e.GET("/image/post/:id", ImageController.GetImageByPostID)
	e.DELETE("/image/avatar/:id", ImageController.DeleteAvatar)
	e.DELETE("/image/post/:id", ImageController.DeleteImageByPostID)
	// Phục vụ file tĩnh từ thư mục "uploads"
	e.Static("/uploads", "uploads")
	return e
}
