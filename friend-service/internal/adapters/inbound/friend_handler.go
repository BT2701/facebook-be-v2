package inbound

import (
	"net/http"
	"friend-service/internal/app/service"
	"friend-service/internal/model"
	"friend-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FriendHandler struct {
	friendService service.FriendService
}

func NewFriendHandler(friendService service.FriendService) *FriendHandler {
	return &FriendHandler{friendService: friendService}
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


