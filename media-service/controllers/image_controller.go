package controllers

import (
	"context"
	"fmt"
	"net/http"

	"snake_api/models"
	"snake_api/services"

	// "snake_api/utils"

	"github.com/labstack/echo/v4"
	"io"
	"os"
)

type ImageController struct {
	service services.ImageService
}

func NewImageController(service services.ImageService) *ImageController {
	return &ImageController{service: service}
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

func (ctrl *ImageController) InsertImage(c echo.Context) error {
	// Lấy file từ form-data
	file, err := c.FormFile("imageFile")
	if err != nil {
		fmt.Println("Error getting form file:", err)
		return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid file"))
	}
	fmt.Println("Form file retrieved successfully")

	// Tạo đường dẫn lưu file
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)

	// Kiểm tra nếu file đã tồn tại
	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("File already exists at path:", filePath)
	} else if os.IsNotExist(err) {
		// Nếu file chưa tồn tại, thực hiện lưu file
		fmt.Println("File does not exist, saving new file")

		// Mở file
		src, err := file.Open()
		if err != nil {
			fmt.Println("Error opening file:", err)
			return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, "Cannot open file"))
		}
		defer src.Close()

		// Tạo file đích để lưu
		dst, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, "Cannot create file"))
		}
		defer dst.Close()

		// Ghi dữ liệu vào file
		if _, err = io.Copy(dst, src); err != nil {
			fmt.Println("Error saving file:", err)
			return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, "Cannot save file"))
		}
		fmt.Println("File saved successfully")
	} else {
		fmt.Println("Error checking file existence:", err)
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, "Error checking file existence"))
	}

	// Lấy thông tin từ form-data
	userID := c.FormValue("user_id")
	postID := c.FormValue("post_id")
	storyID := c.FormValue("story_id")
	fmt.Println("Form values retrieved - userID:", userID, "postID:", postID, "storyID:", storyID)

	// Lưu metadata vào model
	input := models.Image{
		Name:    file.Filename,
		Url:     fmt.Sprintf("/uploads/%s", file.Filename), // URL để truy cập
		UserID:  userID,
		PostID:  postID,
		StoryID: storyID,
	}
	if err := ctrl.service.InsertImage(context.Background(), input); err != nil {
		fmt.Println("Error inserting image metadata:", err)
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}
	fmt.Println("Image metadata inserted successfully")

	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Insert successful",
		"url":     input.Url,
	}, nil))
}



func (ctrl *ImageController) FindAllImages(c echo.Context) error {
	// Gọi hàm service để tìm tất cả ảnh
	Images, err := ctrl.service.FindAllImages(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	// Trả về phản hồi khi tìm thành công
	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, Images, nil))
}

func (ctrl *ImageController) DeleteAllImages(c echo.Context) error {
	// Gọi hàm service để xóa tất cả ảnh
	err := ctrl.service.DeleteAllImages(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	// Trả về phản hồi khi xóa thành công
	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Delete all images successful",
	}, nil))
}

func (ctrl *ImageController) EditImage(c echo.Context) error {
	var input models.Image
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, newAPIResponse(http.StatusBadRequest, nil, "Invalid input"))
	}

	// Gọi hàm service để edit
	err := ctrl.service.EditImage(context.Background(), input.UserID, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	// Trả về phản hồi khi edit thành công
	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Edit successful",
	}, nil))
}

func (ctrl *ImageController) GetImageByUserID(c echo.Context) error {
	id := c.Param("id")

	// Gọi hàm service để tìm ảnh theo ID
	Image, err := ctrl.service.GetImageByUserID(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	// Trả về phản hồi khi tìm thành công
	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, Image, nil))
}

func (ctrl *ImageController) GetImageByPostID(c echo.Context) error {
	id := c.Param("id")

	// Gọi hàm service để tìm ảnh theo ID
	Image, err := ctrl.service.GetImageByPostID(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	// Trả về phản hồi khi tìm thành công
	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, Image, nil))
}

func (ctrl *ImageController) DeleteAvatar(c echo.Context) error {
	id := c.Param("id")

	// Gọi hàm service để xóa ảnh
	err := ctrl.service.DeleteAvatar(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	// Trả về phản hồi khi xóa thành công
	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Delete avatar successful",
	}, nil))
}

func (ctrl *ImageController) DeleteImageByPostID(c echo.Context) error {
	id := c.Param("id")

	// Gọi hàm service để xóa ảnh
	err := ctrl.service.DeleteImageByPostID(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, newAPIResponse(http.StatusInternalServerError, nil, err.Error()))
	}

	// Trả về phản hồi khi xóa thành công
	return c.JSON(http.StatusOK, newAPIResponse(http.StatusOK, map[string]interface{}{
		"message": "Delete image successful",
	}, nil))
}
