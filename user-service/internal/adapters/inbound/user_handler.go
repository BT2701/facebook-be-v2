package inbound

import (
	"context"
	"fmt"
	"net/http"
	"user-service/internal/app/services"
	"user-service/internal/models"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// Chuẩn hóa phản hồi
type APIResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
}

func newAPIResponse(status int, data interface{}, err interface{}) *APIResponse {
	return &APIResponse{
		Status: status,
		Data:   data,
		Error:  err,
	}
}

func (ctrl *UserController) Login(c echo.Context) error {
    var input models.User
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
    }

    // Gọi hàm service để login
    token, err, user := ctrl.service.Login(context.Background(), input.Email, input.Password)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, newAPIResponse(http.StatusUnauthorized, nil, err.Error()))
    }

    // Trả về phản hồi khi login thành công
    return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
        "message": "Login successful",
        "token":   token,
        "user":    user, // Bao gồm thông tin user (nếu cần)
    }, nil))
}


func (ctrl *UserController) SignUp(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	err := ctrl.service.SignUp(context.Background(), user)
	if err != nil {
		return c.JSON(http.StatusConflict, newAPIResponse(http.StatusConflict, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "User created",
		"user":    user,
	}, nil))
}

func (ctrl *UserController) ForgotPassword(c echo.Context) error {
	var input struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid email"))
	}

	resetToken, err := ctrl.service.ForgotPassword(context.Background(), input.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message":     "Password reset email sent successfully",
		"reset_token": resetToken, // Optional: Include reset token for testing purposes
	}, nil))
}

func (ctrl *UserController) ResetPassword(c echo.Context) error {
	var input struct {
		Token    string `json:"token" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	err := ctrl.service.ResetPassword(context.Background(), input.Token, input.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Password reset successfully",
	}, nil))
}

func (ctrl *UserController) GetAllUsers(c echo.Context) error {
	users, err := ctrl.service.GetAllUsers(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, users, nil))
}

func (ctrl *UserController) DeleteAllUsers(c echo.Context) error {
	err := ctrl.service.DeleteAllUsers(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "All users deleted",
	}, nil))
}

func (ctrl *UserController) Logout(c echo.Context) error {
	var input struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid email"))
	}

	err := ctrl.service.Logout(context.Background(), input.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Logged out successfully",
	}, nil))
}

func (ctrl *UserController) EditUser(c echo.Context) error {
	var input struct {
		Email string `json:"email" validate:"required,email"`
		User  models.User
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	err := ctrl.service.EditUser(context.Background(), input.Email, input.User)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "User updated",
	}, nil))
}

func (ctrl *UserController) GetByID(c echo.Context) error {
	id := c.Param("id") // Lấy id từ query string
	fmt.Println("Received ID:", id)

	user, err := ctrl.service.GetByID(context.Background(), id)
	if err != nil {
		// Trả lỗi 404 nếu không tìm thấy user
		if err.Error() == "mongo: no documents in result" {
			fmt.Println("User not found for ID:", id)
			return c.JSON(http.StatusNotFound, newAPIResponse(http.StatusNotFound, nil, "User not found"))
		}
		// Trả lỗi 400 nếu id không hợp lệ
		if err.Error() == "id không hợp lệ" {
			fmt.Println("Invalid ID:", id)
			return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid ID"))
		}
		fmt.Println("Error retrieving user by ID:", err)
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	fmt.Println("User found:", user)
	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, user, nil))
}
func (ctrl *UserController) FindUserByEmail(c echo.Context) error {
    email := c.QueryParam("email")
    if email == "" {
        return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Email is required"))
    }

	fmt.Println("email: ",email)
    user, err := ctrl.service.FindUserByEmail(context.Background(), email)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
    }

    return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, user, nil))
}

func (ctrl *UserController) UpdateAvatar(c echo.Context) error {
	var input struct {
		Email  string `json:"email" validate:"required,email"`
		Avatar string `json:"avatar" validate:"required"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	err := ctrl.service.UpdateAvatar(context.Background(), input.Email, input.Avatar)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Avatar updated",
	}, nil))
}