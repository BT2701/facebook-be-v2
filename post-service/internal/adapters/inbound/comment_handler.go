package inbound

import (
	"net/http"
	"post-service/internal/app/service"
	"post-service/internal/model"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateComment handles the creation of a new comment
func (handler *CommentHandler) CreateComment(c echo.Context) error {
	var comment model.Comment
	if err := c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
	}

	if err := handler.commentService.CreateComment(&comment); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Comment created successfully"})
}

// GetComment retrieves a comment by its ID
func (handler *CommentHandler) GetComment(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing comment ID"})
	}

	comment, err := handler.commentService.GetComment(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if comment == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Comment not found"})
	}

	return c.JSON(http.StatusOK, comment)
}

// UpdateComment updates an existing comment
func (handler *CommentHandler) UpdateComment(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing comment ID"})
	}

	var comment model.Comment
	if err := c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
	}

	comment.ID = objectID // Đảm bảo ID trong request là ID cần cập nhật
	if err := handler.commentService.UpdateComment(&comment); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Comment updated successfully"})
}

// DeleteComment deletes a comment by its ID
func (handler *CommentHandler) DeleteComment(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing comment ID"})
	}

	if err := handler.commentService.DeleteComment(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Comment deleted successfully"})
}
