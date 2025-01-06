package inbound

import (
	"net/http"
	"post-service/internal/model"
	"post-service/internal/app/service"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-service/pkg/utils"
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
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	if err := handler.postService.CreatePost(&post); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Post created successfully",
		"post":    post,
	}, nil))
}

func (handler *PostHandler) GetPost(c echo.Context) error {
	id := c.Param("id")

	post, err := handler.postService.GetPost(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	if post == nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, "Post not found"))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"post": post,
	}, nil))
}

func (handler *PostHandler) UpdatePost(c echo.Context) error {
	id := c.Param("id")

	var post model.Post
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid post ID"))
	}
	post.ID = objectID

	if err := handler.postService.UpdatePost(&post); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Post updated successfully",
		"post":    post,
	}, nil))
}

func (handler *PostHandler) DeletePost(c echo.Context) error {
	id := c.Param("id")

	if err := handler.postService.DeletePost(id); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Post deleted successfully",
	}, nil))
}

func (handler *PostHandler) GetPostsByUserID(c echo.Context) error {
	userID := c.Param("userID")

	posts, err := handler.postService.GetPostsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"userID": userID,
		"posts": posts,
	}, nil))
}

func (handler *PostHandler) GetPosts(c echo.Context) error {
	posts, err := handler.postService.GetPosts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"posts": posts,
	}, nil))
}