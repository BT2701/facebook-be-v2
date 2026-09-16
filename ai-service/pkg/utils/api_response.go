package utils

type APIResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
}

func NewAPIResponse(status int, data interface{}, err interface{}) *APIResponse {
	return &APIResponse{
		Status: status,
		Data:   data,
		Error:  err,
	}
}
