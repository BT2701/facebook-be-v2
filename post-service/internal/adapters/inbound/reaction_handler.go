package inbound

import (
	"net/http"
	"post-service/internal/app/service"
	"post-service/internal/model"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
	}

	if err := handler.reactionService.CreateReaction(&reaction); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Reaction created successfully"})
}

// GetReaction handles retrieving a reaction by ID
func (handler *ReactionHandler) GetReaction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing reaction ID"})
	}

	reaction, err := handler.reactionService.GetReaction(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if reaction == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Reaction not found"})
	}

	return c.JSON(http.StatusOK, reaction)
}

// UpdateReaction handles updating an existing reaction
func (handler *ReactionHandler) UpdateReaction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing reaction ID"})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid reaction ID format"})
	}

	var reaction model.Reaction
	if err := c.Bind(&reaction); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
	}

	reaction.ID = objectID

	if err := handler.reactionService.UpdateReaction(&reaction); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Reaction updated successfully"})
}

// DeleteReaction handles deleting a reaction by ID
func (handler *ReactionHandler) DeleteReaction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing reaction ID"})
	}

	if err := handler.reactionService.DeleteReaction(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Reaction deleted successfully"})
}
