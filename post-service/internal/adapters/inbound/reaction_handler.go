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

type ReactionHandler struct {
	reactionService service.ReactionService
	postService     service.PostService
	bus             *events.Bus
}

func NewReactionHandler(reactionService service.ReactionService, postService service.PostService, bus *events.Bus) *ReactionHandler {
	return &ReactionHandler{reactionService: reactionService, postService: postService, bus: bus}
}

// CreateReaction handles creating a new reaction
func (handler *ReactionHandler) CreateReaction(c echo.Context) error {
	var reaction model.Reaction
	if err := c.Bind(&reaction); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	if err := handler.reactionService.CreateReaction(&reaction); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	handler.publishLike(c, &reaction)

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusCreated, map[string]interface{}{
		"message":  "Reaction created successfully",
		"reaction": reaction,
	}, nil))
}

// GetReaction handles retrieving a reaction by ID
func (handler *ReactionHandler) GetReaction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Missing reaction ID"))
	}

	reaction, err := handler.reactionService.GetReaction(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	if reaction == nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, "Reaction not found"))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"reaction": reaction,
	}, nil))
}

// UpdateReaction handles updating an existing reaction
func (handler *ReactionHandler) UpdateReaction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Missing reaction ID"))
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid reaction ID"))
	}

	var reaction model.Reaction
	if err := c.Bind(&reaction); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	reaction.ID = objectID

	if err := handler.reactionService.UpdateReaction(&reaction); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message":  "Reaction updated successfully",
		"reaction": reaction,
	}, nil))
}

// DeleteReaction handles deleting a reaction by ID
func (handler *ReactionHandler) DeleteReaction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Missing reaction ID"))
	}

	if err := handler.reactionService.DeleteReaction(id); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Reaction deleted successfully",
	}, nil))
}

func (handler *ReactionHandler) DeleteByPostAndUser(c echo.Context) error {
	reaction, err := handler.reactionService.GetByPostAndUser(c.Param("postId"), c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	if reaction == nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, "Reaction not found"))
	}
	if err := handler.reactionService.DeleteReaction(reaction.ID.Hex()); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	if handler.bus != nil && !reaction.PostID.IsZero() {
		post, err := handler.postService.GetPost(reaction.PostID.Hex())
		if err == nil && post != nil {
			handler.bus.PublishDelete(c.Request().Context(), events.NotificationEvent{
				User:     reaction.UserID,
				Receiver: post.UserID,
				Post:     reaction.PostID.Hex(),
				Action:   1,
			})
		}
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Reaction deleted successfully",
	}, nil))
}

func (handler *ReactionHandler) publishLike(c echo.Context, reaction *model.Reaction) {
	if handler.bus == nil || handler.postService == nil || reaction.PostID.IsZero() {
		return
	}
	post, err := handler.postService.GetPost(reaction.PostID.Hex())
	if err != nil || post == nil || post.UserID == "" || post.UserID == reaction.UserID {
		return
	}
	handler.bus.PublishCreate(c.Request().Context(), events.NotificationEvent{
		User:     reaction.UserID,
		Receiver: post.UserID,
		Post:     reaction.PostID.Hex(),
		Content:  "Liked your post",
		Action:   1,
	})
}

func (handler *ReactionHandler) GetByPostAndUser(c echo.Context) error {
	reaction, err := handler.reactionService.GetByPostAndUser(c.Param("postId"), c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	if reaction == nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, "Reaction not found"))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"reaction": reaction,
	}, nil))
}
