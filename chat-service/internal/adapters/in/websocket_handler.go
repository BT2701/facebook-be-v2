package in

import (
    "chat-service/internal/app/service"
    "chat-service/internal/model"
    "github.com/gorilla/websocket"
    "github.com/labstack/echo/v4"
    "log"
    "net/http"
    "sync"
    "time"
)

// SocketHandler quản lý kết nối WebSocket.
type SocketHandler struct {
    Service      *service.ChatService
    Connections  map[string]*websocket.Conn // userID -> WebSocket connection
    ConnectionsMu sync.Mutex
}

// NewSocketHandler khởi tạo SocketHandler.
func NewSocketHandler(service *service.ChatService) *SocketHandler {
    return &SocketHandler{
        Service:     service,
        Connections: make(map[string]*websocket.Conn),
    }
}

// UpgradeWebSocket nâng cấp HTTP lên WebSocket.
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins (modify for production)
    },
}

// HandleConnection xử lý kết nối WebSocket.
func (h *SocketHandler) HandleConnection(c echo.Context) error {
    userID := c.QueryParam("userID")
    if userID == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "userID is required"})
    }

    ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return err
    }

    h.ConnectionsMu.Lock()
    h.Connections[userID] = ws
    h.ConnectionsMu.Unlock()

    log.Printf("User %s connected", userID)

    defer func() {
        h.ConnectionsMu.Lock()
        delete(h.Connections, userID)
        h.ConnectionsMu.Unlock()
        ws.Close()
        log.Printf("User %s disconnected", userID)
    }()

    for {
        var message model.Message
        if err := ws.ReadJSON(&message); err != nil {
            log.Println("Error reading WebSocket message:", err)
            break
        }

        message.ID = time.Now().Format("20060102150405") // Unique ID
        message.Timestamp = time.Now()

        log.Printf("Received message: %+v", message)

        // Lưu tin nhắn vào MongoDB
        if err := h.Service.SendMessage(&message); err != nil {
            log.Println("Error saving message:", err)
            continue
        }

        // Gửi tin nhắn đến người nhận nếu họ đang trực tuyến
        h.ConnectionsMu.Lock()
        if conn, ok := h.Connections[message.Receiver]; ok {
            if err := conn.WriteJSON(message); err != nil {
                log.Println("Error sending message to receiver:", err)
            }
        }
        h.ConnectionsMu.Unlock()
    }

    return nil
}
