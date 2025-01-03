package inbound

import (
	"net/http"
	"post-service/internal/model"
	"post-service/internal/app/service"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (handler *PostHandler) CreatePost(c echo.Context) error {
	var post model.Post
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := handler.postService.CreatePost(&post); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

func (handler *PostHandler) GetPost(c echo.Context) error {
	id := c.Param("id")

	post, err := handler.postService.GetPost(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if post == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
	}

	return c.JSON(http.StatusOK, post)
}

func (handler *PostHandler) UpdatePost(c echo.Context) error {
	id := c.Param("id")

	var post model.Post
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	post.ID = objectID

	if err := handler.postService.UpdatePost(&post); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (handler *PostHandler) DeletePost(c echo.Context) error {
	id := c.Param("id")

	if err := handler.postService.DeletePost(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
