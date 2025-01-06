package inbound

import (
	"net/http"
	"post-service/internal/app/service"
	"post-service/internal/model"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"post-service/pkg/utils"
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
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	if err := handler.storyService.CreateStory(&story); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusCreated, utils.NewAPIResponse(http.StatusCreated, map[string]interface{}{
		"message": "Story created successfully",
		"story":   story,
	}, nil))
}

// GetStory handles retrieving a story by ID
func (handler *StoryHandler) GetStory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Missing story ID"))
	}

	story, err := handler.storyService.GetStory(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	if story == nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, "Story not found"))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"story": story,
	}, nil))
}

// UpdateStory handles updating an existing story
func (handler *StoryHandler) UpdateStory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Missing story ID"))
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid story ID"))
	}

	var story model.Story
	if err := c.Bind(&story); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	story.ID = objectID

	if err := handler.storyService.UpdateStory(&story); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Story updated successfully",
		"story":   story,
	}, nil))
}

// DeleteStory handles deleting a story by ID
func (handler *StoryHandler) DeleteStory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Missing story ID"))
	}

	if err := handler.storyService.DeleteStory(id); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Story deleted successfully",
	}, nil))
	
}
