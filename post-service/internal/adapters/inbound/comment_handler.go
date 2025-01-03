package inbound

import (
	"encoding/json"
	"net/http"
	"post-service/internal/app/service"
	"post-service/internal/model"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// Helper function to write JSON response
func writeJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// Helper function to write JSON error response
func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	writeJSONResponse(w, statusCode, map[string]string{"error": message})
}

func (handler *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	err = handler.commentService.CreateComment(&comment)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]string{"message": "Comment created successfully"})
}

func (handler *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing comment ID")
		return
	}

	comment, err := handler.commentService.GetComment(id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if comment == nil {
		writeJSONError(w, http.StatusNotFound, "Comment not found")
		return
	}

	writeJSONResponse(w, http.StatusOK, comment)
}

func (handler *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing comment ID")
		return
	}

	var comment model.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	comment.ID = id // Đảm bảo ID trong request là ID cần cập nhật
	err = handler.commentService.UpdateComment(&comment)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Comment updated successfully"})
}

func (handler *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing comment ID")
		return
	}

	err := handler.commentService.DeleteComment(id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Comment deleted successfully"})
}
