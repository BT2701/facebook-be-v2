package routes

import (
	"media-service/internal/adapters/inbound"
	"media-service/internal/adapters/outbound"
	"media-service/internal/app/services"
	"os"

	"github.com/BT2701/facebook-be-v2/shared/httpx"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(ImageCollection *mongo.Collection) *echo.Echo {
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URI"), // Địa chỉ Redis
	})
	// Khởi tạo repository, service, và controller
	ImageRepo := outbound.NewImageRepository(ImageCollection)
	ImageService := services.NewImageService(ImageRepo, redisClient)
	ImageController := inbound.NewImageController(ImageService)

	e := echo.New()
	httpx.Apply(e)

	e.POST("/image", ImageController.InsertImage)
	e.GET("/images", ImageController.FindAllImages)
	e.PUT("/image", ImageController.EditImage)
	e.GET("/image/user/:id", ImageController.GetImageByUserID)
	e.GET("/image/post/:id", ImageController.GetImageByPostID)
	e.DELETE("/image/avatar/:id", ImageController.DeleteAvatar)
	e.DELETE("/image/post/:id", ImageController.DeleteImageByPostID)
	// Phục vụ file tĩnh từ thư mục "uploads"
	e.Static("/uploads", "uploads")
	return e
}
