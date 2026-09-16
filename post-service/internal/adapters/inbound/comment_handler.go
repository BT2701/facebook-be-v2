package inbound

import (
	"net/http"
	"post-service/internal/app/service"
	"post-service/internal/model"
	"post-service/pkg/utils"

	"github.com/BT2701/facebook-be-v2/shared/events"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
	commentService service.CommentService
	postService    service.PostService
	bus            *events.Bus
}

func NewCommentHandler(commentService service.CommentService, postService service.PostService, bus *events.Bus) *CommentHandler {
	return &CommentHandler{commentService: commentService, postService: postService, bus: bus}
}

// CreateComment handles the creation of a new comment
func (handler *CommentHandler) CreateComment(c echo.Context) error {
	var comment model.Comment
	if err := c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	if err := handler.commentService.CreateComment(&comment); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	handler.publishComment(c, &comment)
	return c.JSON(http.StatusCreated, utils.NewAPIResponse(http.StatusCreated, map[string]interface{}{
		"message": "Comment created successfully",
		"comment": comment,
	}, nil))
}

// GetComment retrieves a comment by its ID
func (handler *CommentHandler) GetComment(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Missing comment ID"))
	}

	comment, err := handler.commentService.GetComment(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	if comment == nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, "Comment not found"))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"comment": comment,
	}, nil))
}

// UpdateComment updates an existing comment
func (handler *CommentHandler) UpdateComment(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Missing comment ID"))
	}

	var comment model.Comment
	if err := c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid comment ID"))
	}

	comment.ID = objectID // Đảm bảo ID trong request là ID cần cập nhật
	if err := handler.commentService.UpdateComment(&comment); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Comment updated successfully",
		"comment": comment,
	}, nil))
}

// DeleteComment deletes a comment by its ID
func (handler *CommentHandler) DeleteComment(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Missing comment ID"))
	}

	if err := handler.commentService.DeleteComment(id); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Comment deleted successfully",
	}, nil))
}

func (handler *CommentHandler) publishComment(c echo.Context, comment *model.Comment) {
	if handler.bus == nil || handler.postService == nil || comment.PostID.IsZero() {
		return
	}
	post, err := handler.postService.GetPost(comment.PostID.Hex())
	if err != nil || post == nil || post.UserID == "" || post.UserID == comment.UserID {
		return
	}
	handler.bus.PublishCreate(c.Request().Context(), events.NotificationEvent{
		User:     comment.UserID,
		Receiver: post.UserID,
		Post:     comment.PostID.Hex(),
		Content:  "Commented your post",
		Action:   2,
	})
}

func (handler *CommentHandler) GetCommentsByPostID(c echo.Context) error {
	postID := c.Param("postID")
	comments, err := handler.commentService.GetCommentsByPostID(postID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	if comments == nil {
		comments = []model.Comment{}
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"comments": comments,
	}, nil))
}
