package inbound

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"friend-service/pkg/utils"

	"github.com/BT2701/facebook-be-v2/shared/httpc"
	"github.com/labstack/echo/v4"
)

type userPayload struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
}

type usersAPIResponse struct {
	Status int           `json:"status"`
	Data   []userPayload `json:"data"`
}

func (handler *FriendHandler) GetNonFriends(c echo.Context) error {
	userID := c.Param("userID")
	friends, err := handler.friendService.GetFriendsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	excluded := []string{userID}
	seen := map[string]struct{}{userID: {}}
	for _, friend := range friends {
		for _, id := range []string{friend.UserID1, friend.UserID2} {
			if _, skip := seen[id]; skip {
				continue
			}
			seen[id] = struct{}{}
			excluded = append(excluded, id)
		}
	}

	users, err := fetchSuggestedUsers(c.Request().Context(), excluded)
	if err != nil {
		return c.JSON(http.StatusOK, []userPayload{})
	}

	return c.JSON(http.StatusOK, users)
}

func fetchSuggestedUsers(ctx context.Context, exclude []string) ([]userPayload, error) {
	base := os.Getenv("USER_SERVICE_URL")
	if base == "" {
		base = "http://user-service:8080"
	}

	query := url.Values{}
	query.Set("limit", "20")
	if len(exclude) > 0 {
		query.Set("exclude", strings.Join(exclude, ","))
	}

	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	var parsed usersAPIResponse
	if err := httpc.GetJSON(ctx, base+"/api/users?"+query.Encode(), &parsed); err != nil {
		return nil, err
	}
	if parsed.Data == nil {
		return []userPayload{}, nil
	}
	return parsed.Data, nil
}
