package inbound

import (
	"net/http"
	"post-service/internal/app/service"
	"post-service/internal/model"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-service/pkg/utils"
)

type ReactionHandler struct {
	reactionService service.ReactionService
}

func NewReactionHandler(reactionService service.ReactionService) *ReactionHandler {
	return &ReactionHandler{reactionService: reactionService}
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
