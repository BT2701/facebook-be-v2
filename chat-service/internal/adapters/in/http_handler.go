package in

import (
    "chat-service/internal/app/service"
    "chat-service/internal/model"
    "github.com/labstack/echo/v4"
    "net/http"
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
    // "log"
)

type ChatHandler struct {
    Service *service.ChatService
}

// SendMessage handles the HTTP request to send a message.
func (h *ChatHandler) SendMessage(c echo.Context) error {
    var message model.Message

    // Bind request body to message struct
    if err := c.Bind(&message); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "status": "error",
            "data":   nil,
            "error":  "Invalid input",
        })
    }

    message.ID = primitive.NewObjectID().Hex()
    message.CreatedAt = time.Now()

    // Call service to save the message
    if err := h.Service.SendMessage(&message); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "status": "error",
            "data":   nil,
            "error":  "Failed to send message",
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": "success",
        "data":   message,
        "error":  nil,
    })
}

// GetMessages handles the HTTP request to get messages between two users.
func (h *ChatHandler) GetMessages(c echo.Context) error {
    sender := c.QueryParam("sender")
    receiver := c.QueryParam("receiver")

    if sender == "" || receiver == "" {
        return c.JSON(http.StatusBadRequest, map[string]interface{}{
            "status": "error",
            "data":   nil,
            "error":  "Sender and receiver are required",
        })
    }

    messages, err := h.Service.GetMessages(sender, receiver)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]interface{}{
            "status": "error",
            "data":   nil,
            "error":  "Failed to get messages",
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": "success",
        "data":   messages,
        "error":  nil,
    })
}
func (h *ChatHandler) GetAllMessages(c echo.Context) error {
    messages, err := h.Service.GetAllMessages()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get messages"})
    }

    return c.JSON(http.StatusOK, messages)
}
func (h *ChatHandler) DeleteAllMessages(c echo.Context) error {
    err := h.Service.DeleteAllMessages()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete messages"})
    }

    return c.JSON(http.StatusOK, map[string]string{"status": "All messages deleted successfully"})
}