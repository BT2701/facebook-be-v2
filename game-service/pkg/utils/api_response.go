package utils

import (
	"encoding/json"
	"net/http"
)
// Chuẩn hóa phản hồi
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

func WriteJSONResponse(w http.ResponseWriter, response *APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}