package inbound

import (
	"net/http"
	"friend-service/internal/app/service"
	"friend-service/internal/model"
	"friend-service/pkg/utils"

	"github.com/BT2701/facebook-be-v2/shared/events"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestHandler struct {
	requestService service.RequestService
	bus            *events.Bus
}

func NewRequestHandler(requestService service.RequestService, bus *events.Bus) *RequestHandler {
	return &RequestHandler{requestService: requestService, bus: bus}
}

func (handler *RequestHandler) CreateRequest(c echo.Context) error {
	var request model.Request
	request.ID = primitive.NewObjectID()

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	createdRequest, err := handler.requestService.CreateRequest(request.Sender, request.Receiver)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	if handler.bus != nil {
		handler.bus.PublishCreate(c.Request().Context(), events.NotificationEvent{
			User:     request.Sender,
			Receiver: request.Receiver,
			Post:     "0",
			Content:  "sent you a friend request",
			Action:   3,
		})
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Request created successfully",
		"request": createdRequest,
	}, nil))
}

func (handler *RequestHandler) GetRequest(c echo.Context) error {
	sender := c.Param("sender")
	receiver := c.Param("receiver")

	request, err := handler.requestService.GetRequest(sender, receiver)
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

func (handler *RequestHandler) GetRequests(c echo.Context) error {
	receiver := c.Param("receiver")

	requests, err := handler.requestService.GetRequests(receiver)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"requests": requests,
	}, nil))
}

func (handler *RequestHandler) UpdateRequest(c echo.Context) error {
	sender := c.Param("sender")
	receiver := c.Param("receiver")

	var input struct {
		IsAccepted bool `json:"isAccepted"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	updatedRequest, err := handler.requestService.UpdateRequest(sender, receiver, input.IsAccepted)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Request updated successfully",
		"request": updatedRequest,
	}, nil))
}

func (handler *RequestHandler) DeleteRequest(c echo.Context) error {
	sender := c.Param("sender")
	receiver := c.Param("receiver")
	if sender == "" {
		sender = c.QueryParam("senderId")
	}
	if receiver == "" {
		receiver = c.QueryParam("receiverId")
	}

	if err := handler.requestService.DeleteRequest(sender, receiver); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	if handler.bus != nil {
		handler.bus.PublishDelete(c.Request().Context(), events.NotificationEvent{
			User:     sender,
			Receiver: receiver,
			Post:     "0",
			Action:   3,
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (handler *RequestHandler) DeleteRequestByID(c echo.Context) error {
	id := c.Param("id")
	if err := handler.requestService.DeleteRequestByID(id); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (handler *RequestHandler) GetRequestsByQuery(c echo.Context) error {
	receiver := c.QueryParam("id")
	if receiver == "" {
		receiver = c.Param("receiver")
	}
	requests, err := handler.requestService.GetRequests(receiver)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	if requests == nil {
		requests = []*model.Request{}
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"requests": requests,
	}, nil))
}
