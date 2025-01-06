package inbound

import (
	"net/http"
	"friend-service/internal/app/service"
	"friend-service/internal/model"
	"friend-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestHandler struct {
	friendService service.FriendService
}

func NewRequestHandler(friendService service.FriendService) *RequestHandler {
	return &RequestHandler{friendService: friendService}
}

func (handler *RequestHandler) CreateFriend(c echo.Context) error {
	var request model.Request
	request.ID = primitive.NewObjectID()
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	if err := handler.friendService.CreateFriend(request.Sender, request.Receiver); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Request created successfully",
		"request":    request,
	}, nil))
}

func (handler *RequestHandler) GetFriend(c echo.Context) error {
	sender := c.Param("sender")
	receiver := c.Param("receiver")

	request, err := handler.friendService.GetFriend(sender, receiver)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	if request == nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, "Request not found"))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"request": request,
	}, nil))
}

func (handler *RequestHandler) GetFriends(c echo.Context) error {
	receiver := c.Param("receiver")

	requests, err := handler.friendService.GetFriends(receiver)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"requests": requests,
	}, nil))
}

func (handler *RequestHandler) UpdateFriend(c echo.Context) error {
	sender := c.Param("sender")
	receiver := c.Param("receiver")

	var request model.Request
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	objectID, err := primitive.ObjectIDFromHex(sender)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid request ID"))
	}
	request.ID = objectID

	if err := handler.friendService.UpdateFriend(request); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Request updated successfully",
		"request":    request,
	}, nil))
}

func (handler *RequestHandler) DeleteFriend(c echo.Context) error {
	sender := c.Param("sender")
	receiver := c.Param("receiver")

	if err := handler.friendService.DeleteFriend(sender, receiver); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Request deleted successfully",
	}, nil))
}
