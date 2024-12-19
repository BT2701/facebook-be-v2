package utils

// Định nghĩa enum cho notification types
type Notification int

const (
    Success Notification = iota + 1 // 1
    Error                           // 2
    NotFound                        // 3
    Unauthorized                    // 4
    InvalidInput                    // 5
)

// Hàm để chuyển Notification thành chuỗi dễ hiểu
func (n Notification) String() string {
    details := map[Notification]string{
        Success:      "Success: 200 - Operation completed successfully",
        Error:        "Error: 500 - Internal Server Error",
        NotFound:     "NotFound: 404 - Resource not found",
        Unauthorized: "Unauthorized: 401 - Unauthorized access",
        InvalidInput: "InvalidInput: 400 - Invalid input provided",
    }

    if msg, exists := details[n]; exists {
        return msg
    }
    return "Unknown notification type"
}
func GetSuccessMessage() string {
    return Success.String()
}

func GetErrorMessage() string {
    return Error.String()
}

func GetNotFoundMessage() string {
    return NotFound.String()
}

func GetUnauthorizedMessage() string {
    return Unauthorized.String()
}

func GetInvalidInputMessage() string {
    return InvalidInput.String()
}