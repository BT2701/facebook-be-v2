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

type FriendHandler struct {
	friendService  service.FriendService
	requestService service.RequestService
	bus            *events.Bus
}

func NewFriendHandler(friendService service.FriendService) *FriendHandler {
	return &FriendHandler{friendService: friendService}
}

func NewFriendHandlerWithRequests(friendService service.FriendService, requestService service.RequestService, bus *events.Bus) *FriendHandler {
	return &FriendHandler{friendService: friendService, requestService: requestService, bus: bus}
}

func (handler *FriendHandler) CreateFriend(c echo.Context) error {
    var friend *model.Friend
    friend = &model.Friend{} // Khởi tạo con trỏ trước khi gán giá trị
    friend.ID = primitive.NewObjectID()

    if err := c.Bind(friend); err != nil {
        return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
    }

    createdFriend, err := handler.friendService.CreateFriend(friend.UserID1, friend.UserID2)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
    }

    return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
        "message": "Friend created successfully",
        "friend":  createdFriend,
    }, nil))
}


func (handler *FriendHandler) RemoveFriend(c echo.Context) error {
	userID := c.Param("userId")
	friendID := c.Param("friendId")
	if err := handler.friendService.DeleteFriend(userID, friendID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Friend removed successfully",
	}, nil))
}

func (handler *FriendHandler) GetFriend(c echo.Context) error {
	userID1 := c.Param("userID1")
	userID2 := c.Param("userID2")

	friend, err := handler.friendService.GetFriend(userID1, userID2)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	if friend == nil {
		return c.JSON(http.StatusNotFound, utils.NewAPIResponse(http.StatusNotFound, nil, "Friend not found"))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"friend": friend,
	}, nil))
}

func (handler *FriendHandler) GetFriends(c echo.Context) error {
	userID := c.Param("userID")
	if userID == "" {
		userID = c.QueryParam("id")
	}

	friends, err := handler.friendService.GetFriends(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"friends": friends,
	}, nil))
}

func (handler *FriendHandler) UpdateFriend(c echo.Context) error {
	userID1 := c.Param("userID1")
	userID2 := c.Param("userID2")

	var friend *model.Friend
	if err := c.Bind(&friend); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	friend, err := handler.friendService.UpdateFriend(userID1, userID2, friend.IsFriend)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Friend updated successfully",
		"friend":    friend,
	}, nil))
}

func (handler *FriendHandler) DeleteFriend(c echo.Context) error {
	userID1 := c.Param("userID1")
	userID2 := c.Param("userID2")

	if err := handler.friendService.DeleteFriend(userID1, userID2); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Friend deleted successfully",
	}, nil))
}

func (handler *FriendHandler) GetFriendsByUserID(c echo.Context) error {
	userID := c.Param("userID")

	friends, err := handler.friendService.GetFriendsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"friends": friends,
	}, nil))
}

func (handler *FriendHandler) IsFriend(c echo.Context) error {
	userID1 := c.Param("userID1")
	userID2 := c.Param("userID2")

	friend, err := handler.friendService.GetFriend(userID1, userID2)
	if err != nil || friend == nil {
		return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
			"isFriend": false,
			"friend":   nil,
		}, nil))
	}

	return c.JSON(http.StatusOK, utils.NewAPIResponse(http.StatusOK, map[string]interface{}{
		"isFriend": true,
		"friend":   friend,
	}, nil))
}

func (handler *FriendHandler) CreateAndDeleteRequest(c echo.Context) error {
	var input struct {
		UserID1   string `json:"userId1"`
		UserID2   string `json:"userId2"`
		RequestID string `json:"requestId"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	created, err := handler.friendService.CreateFriend(input.UserID1, input.UserID2)
	if err != nil {
		return c.JSON(http.StatusConflict, utils.NewAPIResponse(http.StatusConflict, nil, err.Error()))
	}
	if handler.requestService != nil {
		if input.RequestID != "" {
			_ = handler.requestService.DeleteRequestByID(input.RequestID)
		} else {
			_ = handler.requestService.DeleteRequest(input.UserID1, input.UserID2)
		}
	}
	if handler.bus != nil {
		handler.bus.PublishCreate(c.Request().Context(), events.NotificationEvent{
			User:     input.UserID1,
			Receiver: input.UserID2,
			Post:     "0",
			Content:  "accepted your friend request",
			Action:   4,
		})
		handler.bus.PublishDelete(c.Request().Context(), events.NotificationEvent{
			User:     input.UserID2,
			Receiver: input.UserID1,
			Post:     "0",
			Action:   3,
		})
	}

	return c.JSON(http.StatusCreated, utils.NewAPIResponse(http.StatusCreated, map[string]interface{}{
		"message": "Friend created and request removed",
		"friend":  created,
	}, nil))
}