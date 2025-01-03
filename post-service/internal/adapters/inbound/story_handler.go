package inbound

import (
	"net/http"
	"post-service/internal/app/service"
	"post-service/internal/model"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryHandler struct {
	storyService service.StoryService
}

func NewStoryHandler(storyService service.StoryService) *StoryHandler {
	return &StoryHandler{storyService: storyService}
}

// CreateStory handles creating a new story
func (handler *StoryHandler) CreateStory(c echo.Context) error {
	var story model.Story
	if err := c.Bind(&story); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
	}

	if err := handler.storyService.CreateStory(&story); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Story created successfully"})
}

// GetStory handles retrieving a story by ID
func (handler *StoryHandler) GetStory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing story ID"})
	}

	story, err := handler.storyService.GetStory(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if story == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Story not found"})
	}

	return c.JSON(http.StatusOK, story)
}

// UpdateStory handles updating an existing story
func (handler *StoryHandler) UpdateStory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing story ID"})
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid story ID format"})
	}

	var story model.Story
	if err := c.Bind(&story); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
	}

	story.ID = objectID

	if err := handler.storyService.UpdateStory(&story); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Story updated successfully"})
}

// DeleteStory handles deleting a story by ID
func (handler *StoryHandler) DeleteStory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing story ID"})
	}

	if err := handler.storyService.DeleteStory(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Story deleted successfully"})
}
